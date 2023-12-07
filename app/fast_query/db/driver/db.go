package driver

import (
	"fmt"
	"math"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/terra-money/core/v2/app/fast_query/db/height_driver"
	"github.com/terra-money/core/v2/app/fast_query/db/utils"
)

type DBDriver struct {
	session *dbm.GoLevelDB
	mode    int
}

func NewDBDriver(dir string) (*DBDriver, error) {
	ldb, err := dbm.NewGoLevelDB(DBName, dir)
	if err != nil {
		return nil, err
	}

	return &DBDriver{
		session: ldb,
		mode:    DriverModeKeySuffixDesc,
	}, nil
}

func (dbDriver *DBDriver) newInnerIterator(requestHeight int64, pdb *dbm.PrefixDB) (dbm.Iterator, error) {
	if dbDriver.mode == DriverModeKeySuffixAsc {
		heightEnd := utils.UintToBigEndian(uint64(requestHeight + 1))
		return pdb.ReverseIterator(nil, heightEnd)
	} else {
		heightStart := utils.UintToBigEndian(math.MaxUint64 - uint64(requestHeight))
		return pdb.Iterator(heightStart, nil)
	}
}

func (dbDriver *DBDriver) Get(maxHeight int64, key []byte) ([]byte, error) {
	if maxHeight == 0 {
		return dbDriver.session.Get(prefixCurrentDataKey(key))
	}
	var requestHeight = height_driver.Height(maxHeight).CurrentOrLatest().ToInt64()
	var requestHeightMin = height_driver.Height(0).CurrentOrNever().ToInt64()

	// check if requestHeightMin is
	if requestHeightMin > requestHeight {
		return nil, fmt.Errorf("invalid height")
	}

	pdb := dbm.NewPrefixDB(dbDriver.session, prefixDataWithHeightKey(key))

	iter, _ := dbDriver.newInnerIterator(requestHeight, pdb)
	defer iter.Close()

	// in tm-db@v0.6.4, key not found is NOT an error
	if !iter.Valid() {
		return nil, nil
	}

	value := iter.Value()
	deleted := value[0]
	if deleted == 1 {
		return nil, nil
	} else {
		if len(value) > 1 {
			return value[1:], nil
		}
		return []byte{}, nil
	}
}

func (dbDriver *DBDriver) Has(maxHeight int64, key []byte) (bool, error) {
	if maxHeight == 0 {
		return dbDriver.session.Has(prefixCurrentDataKey(key))
	}
	var requestHeight = height_driver.Height(maxHeight).CurrentOrLatest().ToInt64()
	var requestHeightMin = height_driver.Height(0).CurrentOrNever().ToInt64()

	// check if requestHeightMin is
	if requestHeightMin > requestHeight {
		return false, fmt.Errorf("invalid height")
	}

	pdb := dbm.NewPrefixDB(dbDriver.session, prefixDataWithHeightKey(key))

	iter, _ := dbDriver.newInnerIterator(requestHeight, pdb)
	defer iter.Close()

	// in tm-db@v0.6.4, key not found is NOT an error
	if !iter.Valid() {
		return false, nil
	}

	deleted := iter.Value()[0]

	if deleted == 1 {
		return false, nil
	} else {
		return true, nil
	}
}

func (dbDriver *DBDriver) Set(atHeight int64, key, value []byte) error {
	// should never reach here, all should be batched in tiered+hld
	panic("should never reach here")
}

func (dbDriver *DBDriver) SetSync(atHeight int64, key, value []byte) error {
	// should never reach here, all should be batched in tiered+hld
	panic("should never reach here")
}

func (dbDriver *DBDriver) Delete(atHeight int64, key []byte) error {
	// should never reach here, all should be batched in tiered+hld
	panic("should never reach here")
}

func (dbDriver *DBDriver) DeleteSync(atHeight int64, key []byte) error {
	return dbDriver.Delete(atHeight, key)
}

func (dbDriver *DBDriver) Iterator(maxHeight int64, start, end []byte) (height_driver.HeightEnabledIterator, error) {
	if maxHeight == 0 {
		pdb := dbm.NewPrefixDB(dbDriver.session, cCurrentDataPrefix)
		return pdb.Iterator(start, end)
	}
	return NewLevelDBIterator(dbDriver, maxHeight, start, end)
}

func (dbDriver *DBDriver) ReverseIterator(maxHeight int64, start, end []byte) (height_driver.HeightEnabledIterator, error) {
	if maxHeight == 0 {
		pdb := dbm.NewPrefixDB(dbDriver.session, cCurrentDataPrefix)
		return pdb.ReverseIterator(start, end)
	}
	return NewLevelDBReverseIterator(dbDriver, maxHeight, start, end)
}

func (dbDriver *DBDriver) Close() error {
	dbDriver.session.Close()
	return nil
}

func (dbDriver *DBDriver) NewBatch(atHeight int64) height_driver.HeightEnabledBatch {
	return NewLevelDBBatch(atHeight, dbDriver)
}

// TODO: Implement me
func (dbDriver *DBDriver) Print() error {
	return nil
}

func (dbDriver *DBDriver) Stats() map[string]string {
	return nil
}
