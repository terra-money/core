package driver

import (
	"fmt"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/terra-money/core/v2/app/fast_query/db/height_driver"
)

var _ height_driver.HeightEnabledBatch = (*DriverBatch)(nil)
var _ HasRollbackBatch = (*DriverBatch)(nil)

type DriverBatch struct {
	height int64
	batch  *RollbackableBatch
	mode   int
}

func (b *DriverBatch) keyBytesWithHeight(key []byte) []byte {
	return append(prefixDataWithHeightKey(key), serializeHeight(b.mode, b.height)...)
}

func NewLevelDBBatch(atHeight int64, dbDriver *DBDriver) *DriverBatch {
	return &DriverBatch{
		height: atHeight,
		batch:  NewRollbackableBatch(dbDriver.session),
		mode:   dbDriver.mode,
	}
}

func (b *DriverBatch) Set(key, value []byte) error {
	newKey := b.keyBytesWithHeight(key)

	// make fixed size byte slice for performance
	buf := make([]byte, 0, len(value)+1)
	buf = append(buf, byte(0)) // 0 => not deleted
	buf = append(buf, value...)

	if err := b.batch.Set(prefixCurrentDataKey(key), buf[1:]); err != nil {
		return err
	}
	if err := b.batch.Set(prefixKeysForIteratorKey(key), []byte{}); err != nil {
		return err
	}
	return b.batch.Set(newKey, buf)
}

func (b *DriverBatch) Delete(key []byte) error {
	newKey := b.keyBytesWithHeight(key)

	buf := []byte{1}

	if err := b.batch.Delete(prefixCurrentDataKey(key)); err != nil {
		return err
	}
	if err := b.batch.Set(prefixKeysForIteratorKey(key), buf); err != nil {
		return err
	}
	return b.batch.Set(newKey, buf)
}

func (b *DriverBatch) Write() error {
	return b.batch.Write()
}

func (b *DriverBatch) WriteSync() error {
	return b.batch.WriteSync()
}

func (b *DriverBatch) Close() error {
	return b.batch.Close()
}

func (b *DriverBatch) RollbackBatch() tmdb.Batch {
	return b.batch.RollbackBatch
}

func (b *DriverBatch) Metric() {
	fmt.Printf("[rollback-batch] rollback batch for height %d's record length %d\n",
		b.height,
		b.batch.RecordCount,
	)
}
