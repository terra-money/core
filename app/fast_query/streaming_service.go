package fast_query

import (
	"context"
	"sort"
	"strings"
	"sync"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store/types"
)

var _ baseapp.StreamingService = &StreamingService{}

// StreamingService is a concrete implementation of StreamingService that accumulate the state changes in current block,
// writes the ordered changeset out to version storage.
type StreamingService struct {
	listeners          []*types.MemoryListener // the listeners that will be initialized with BaseApp
	fastQueryService   *FastQueryService
	currentBlockNumber int64 // the current block number
}

// NewStreamingService creates a new StreamingService for the provided writeDir, (optional) filePrefix, and storeKeys
func NewStreamingService(fastQueryService *FastQueryService, appKeys map[string]*types.KVStoreKey) *StreamingService {
	// Create a copy of the store keys to avoid mutating the app.keys
	storeKeys := make([]types.StoreKey, 0, len(appKeys))
	for _, storeKey := range appKeys {
		storeKeys = append(storeKeys, storeKey)
	}

	// sort by the storeKeys first to avoid indeterministic order
	sort.SliceStable(storeKeys, func(i, j int) bool {
		return strings.Compare(storeKeys[i].Name(), storeKeys[j].Name()) < 0
	})

	listeners := make([]*types.MemoryListener, len(storeKeys))
	for i, key := range storeKeys {
		listeners[i] = types.NewMemoryListener(key)
	}
	return &StreamingService{listeners, fastQueryService, 0}
}

// Listeners satisfies the baseapp.StreamingService interface
func (fss *StreamingService) Listeners() map[types.StoreKey][]types.WriteListener {
	listeners := make(map[types.StoreKey][]types.WriteListener, len(fss.listeners))
	for _, listener := range fss.listeners {
		listeners[listener.StoreKey()] = []types.WriteListener{listener}
	}
	return listeners
}

// ListenBeginBlock satisfies the baseapp.ABCIListener interface
func (fss *StreamingService) ListenBeginBlock(ctx context.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock) error {
	fss.currentBlockNumber = req.GetHeader().Height
	return nil
}

// ListenDeliverTx satisfies the baseapp.ABCIListener interface
func (fss *StreamingService) ListenDeliverTx(ctx context.Context, req abci.RequestDeliverTx, res abci.ResponseDeliverTx) error {
	return nil
}

// ListenEndBlock satisfies the baseapp.ABCIListener interface
func (fss *StreamingService) ListenEndBlock(ctx context.Context, req abci.RequestEndBlock, res abci.ResponseEndBlock) error {
	return nil
}

func (fss *StreamingService) ListenCommit(ctx context.Context, res abci.ResponseCommit) error {
	// concat the state caches
	var changeSet []types.StoreKVPair
	for _, listener := range fss.listeners {
		changeSet = append(changeSet, listener.PopStateCache()...)
	}

	return fss.fastQueryService.CommitChanges(fss.currentBlockNumber, changeSet)
}

// Stream satisfies the baseapp.StreamingService interface
func (fss *StreamingService) Stream(wg *sync.WaitGroup) error {
	return nil
}

// Close satisfies the io.Closer interface, which satisfies the baseapp.StreamingService interface
func (fss *StreamingService) Close() error {
	return nil
}
