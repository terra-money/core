package height_driver

import (
	dbm "github.com/cometbft/cometbft-db"
)

type HeigthDB interface {
	dbm.DB
	SetReadHeight(int64)
	ClearReadHeight() int64
	SetWriteHeight(int64)
	ClearWriteHeight() int64
}

type HeightEnabledDB interface {
	// Get fetches the value of the given key, or nil if it does not exist.
	// CONTRACT: key, value readonly []byte
	Get(maxHeight int64, key []byte) ([]byte, error)

	// Has checks if a key exists.
	// CONTRACT: key, value readonly []byte
	Has(maxHeight int64, key []byte) (bool, error)

	// Set sets the value for the given key, replacing it if it already exists.
	// CONTRACT: key, value readonly []byte
	Set(atHeight int64, key, value []byte) error

	// SetSync sets the value for the given key, and flushes it to storage before returning.
	SetSync(atHeight int64, key, value []byte) error

	// Delete deletes the key, or does nothing if the key does not exist.
	// CONTRACT: key readonly []byte
	Delete(atHeight int64, key []byte) error

	// DeleteSync deletes the key, and flushes the delete to storage before returning.
	DeleteSync(atHeight int64, key []byte) error

	// Iterator returns an iterator over a domain of keys, in ascending order. The caller must call
	// Close when done. End is exclusive, and start must be less than end. A nil start iterates
	// from the first key, and a nil end iterates to the last key (inclusive).
	// CONTRACT: No writes may happen within a domain while an iterator exists over it.
	// CONTRACT: start, end readonly []byte
	Iterator(maxHeight int64, start, end []byte) (HeightEnabledIterator, error)

	// ReverseIterator returns an iterator over a domain of keys, in descending order. The caller
	// must call Close when done. End is exclusive, and start must be less than end. A nil end
	// iterates from the last key (inclusive), and a nil start iterates to the first key (inclusive).
	// CONTRACT: No writes may happen within a domain while an iterator exists over it.
	// CONTRACT: start, end readonly []byte
	ReverseIterator(maxHeight int64, start, end []byte) (HeightEnabledIterator, error)

	// Close closes the database connection.
	Close() error

	// NewBatch creates a batch for atomic updates. The caller must call Batch.Close.
	NewBatch(atHeight int64) HeightEnabledBatch

	// Print is used for debugging.
	Print() error

	// Stats returns a map of property values for all keys and the size of the cache.
	Stats() map[string]string
}

type HeightEnabledIterator interface {
	dbm.Iterator
}

type HeightEnabledBatch interface {
	dbm.Batch
}
