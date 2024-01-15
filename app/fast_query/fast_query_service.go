package fast_query

import (
	"fmt"

	log "github.com/cometbft/cometbft/libs/log"
	"github.com/terra-money/core/v2/app/fast_query/db/driver"
	"github.com/terra-money/core/v2/app/fast_query/db/height_driver"
	"github.com/terra-money/core/v2/app/fast_query/store"

	"github.com/cosmos/cosmos-sdk/store/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type FastQueryService struct {
	Store             *store.Store
	safeBatchDBCloser driver.SafeBatchDBCloser
	fastQueryDb       *height_driver.HeightDB
	logger            log.Logger
}

func NewFastQueryService(homedir string, logger log.Logger, storeKeys map[string]*types.KVStoreKey) (*FastQueryService, error) {
	// Create a copy of the store keys
	fastQueryDbDriver, err := driver.NewDBDriver(homedir)
	if err != nil {
		return nil, err
	}

	// Create HeightDB Driver that implements optimization for reading
	// and writing data in the database in parallel.
	fastQueryDb := height_driver.NewHeightDB(
		fastQueryDbDriver,
		&height_driver.HeightDBConfig{
			Debug: true,
		},
	)
	// Create the new BatchingDB with it's safe batch closer
	heightEnabledDB := driver.NewSafeBatchDB(fastQueryDb)
	safeBatchDBCloser := heightEnabledDB.(driver.SafeBatchDBCloser)
	store, err := store.NewStore(heightEnabledDB, fastQueryDb, logger, storeKeys)
	if err != nil {
		return nil, err
	}

	return &FastQueryService{
		Store:             store,
		safeBatchDBCloser: safeBatchDBCloser,
		fastQueryDb:       fastQueryDb,
		logger:            logger,
	}, err
}

func (fqs *FastQueryService) CommitChanges(blockHeight int64, changeSet []types.StoreKVPair) error {
	fqs.logger.Debug("CommitChanges", "blockHeight", blockHeight, "changeSet", changeSet)
	if blockHeight-fqs.Store.LatestVersion() != 1 {
		fmt.Println(fmt.Sprintf("invalid block height: %s vs %s", blockHeight, fqs.Store.LatestVersion()))
		panic("")
	}
	fqs.fastQueryDb.SetWriteHeight(blockHeight)
	fqs.safeBatchDBCloser.Open()

	for _, change := range changeSet {
		storeKey := storetypes.NewKVStoreKey(change.StoreKey)
		commitKVStore := fqs.Store.GetStoreByName(storeKey.Name()).(types.CommitKVStore)
		if change.Delete {
			commitKVStore.Delete(change.Key)
		} else {
			commitKVStore.Set(change.Key, change.Value)
		}
	}

	fqs.Store.Commit()
	if _, err := fqs.safeBatchDBCloser.Flush(); err != nil {
		return err
	}
	fqs.fastQueryDb.ClearWriteHeight()
	return nil
}
