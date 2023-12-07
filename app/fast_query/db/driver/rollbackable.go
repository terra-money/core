package driver

import (
	mdb "github.com/cometbft/cometbft-db"
)

type HasRollbackBatch interface {
	RollbackBatch() mdb.Batch
}

var _ mdb.Batch = (*RollbackableBatch)(nil)

type RollbackableBatch struct {
	mdb.Batch

	db            mdb.DB
	RollbackBatch mdb.Batch
	RecordCount   int
}

func NewRollbackableBatch(db mdb.DB) *RollbackableBatch {
	return &RollbackableBatch{
		db:            db,
		Batch:         db.NewBatch(),
		RollbackBatch: db.NewBatch(),
	}
}

// revert value for key to previous state
func (b *RollbackableBatch) backup(key []byte) error {
	b.RecordCount++
	data, err := b.db.Get(key)
	if err != nil {
		return err
	}
	if data == nil {
		return b.RollbackBatch.Delete(key)
	} else {
		return b.RollbackBatch.Set(key, data)
	}
}

func (b *RollbackableBatch) Set(key, value []byte) error {
	if err := b.backup(key); err != nil {
		return err
	}
	return b.Batch.Set(key, value)
}

func (b *RollbackableBatch) Delete(key []byte) error {
	if err := b.backup(key); err != nil {
		return err
	}
	return b.Batch.Delete(key)
}
