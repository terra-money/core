package driver

import mdb "github.com/cometbft/cometbft-db"

var _ mdb.Batch = (*SafeBatchNullified)(nil)

type SafeBatchNullified struct {
	batch mdb.Batch
}

func NewSafeBatchNullify(batch mdb.Batch) mdb.Batch {
	return &SafeBatchNullified{
		batch: batch,
	}
}

func (s SafeBatchNullified) Set(key, value []byte) error {
	return s.batch.Set(key, value)
}

func (s SafeBatchNullified) Delete(key []byte) error {
	return s.batch.Delete(key)
}

func (s SafeBatchNullified) Write() error {
	// noop
	return nil
}

func (s SafeBatchNullified) WriteSync() error {
	return s.Write()
}

func (s SafeBatchNullified) Close() error {
	// noop
	return nil
}
