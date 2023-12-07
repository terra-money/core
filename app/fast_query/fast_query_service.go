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
	fastQueryDriver   height_driver.HeightEnabledDB
	logger            log.Logger
}

func NewFastQueryService(homedir string, logger log.Logger) (*FastQueryService, error) {
	// Create a new instance of the Database Driver that uses LevelDB
	fastQueryDriver, err := driver.NewDBDriver(homedir)
	if err != nil {
		return nil, err
	}

	// Create HeightDB Driver that implements optimization for reading
	// and writing data in the database in paralell.
	fastQueryHeightDriver := height_driver.NewHeightDB(
		fastQueryDriver,
		&height_driver.HeightDBConfig{
			Debug: true,
		},
	)

	heightEnabledDB := driver.NewSafeBatchDB(fastQueryHeightDriver)
	safeBatchDBCloser := heightEnabledDB.(driver.SafeBatchDBCloser)
	store := store.NewStore(heightEnabledDB, fastQueryHeightDriver, logger)

	return &FastQueryService{
		Store:             store,
		safeBatchDBCloser: safeBatchDBCloser,
		fastQueryDriver:   fastQueryDriver,
		logger:            logger,
	}, err
}

func (fqs *FastQueryService) CommitChanges(blockHeight int64, changeSet []types.StoreKVPair) error {
	fqs.logger.Info("CommitChanges", "blockHeight", blockHeight, "changeSet", changeSet)
	return nil
}
