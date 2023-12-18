package fast_query

import (
	log "github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/terra-money/core/v2/app/fast_query/db/driver"
	"github.com/terra-money/core/v2/app/fast_query/db/height_driver"
	"github.com/terra-money/core/v2/app/fast_query/store"
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
	// and writing data in the database in paralell.
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

	// store.LoadLatestVersion()

	return &FastQueryService{
		Store:             store,
		safeBatchDBCloser: safeBatchDBCloser,
		fastQueryDb:       fastQueryDb,
		logger:            logger,
	}, err
}

func (fqs *FastQueryService) CommitChanges(blockHeight int64, changeSet []types.StoreKVPair) error {
	fqs.logger.Info("CommitChanges", "blockHeight", blockHeight, "changeSet", changeSet)

	fqs.fastQueryDb.SetWriteHeight(blockHeight)
	fqs.safeBatchDBCloser.Open()

	for _, kv := range changeSet {
		key := fqs.Store.StoreKeysByName()[kv.StoreKey]
		ckvs := fqs.Store.GetCommitKVStore(key)
		if kv.Delete {
			ckvs.Delete(kv.Key)
		} else {
			ckvs.Set(kv.Key, kv.Value)
		}
	}

	if _, err := fqs.safeBatchDBCloser.Flush(); err != nil {
		return err
	}
	return nil
}
