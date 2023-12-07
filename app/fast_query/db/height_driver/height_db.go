package height_driver

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/terra-money/core/v2/app/fast_query/db/utils"

	tmdb "github.com/cometbft/cometbft-db"
)

const (
	LatestHeight  = 0
	InvalidHeight = 0

	debugKeyGet = iota
	debugKeySet
	debugKeyIterator
	debugKeyReverseIterator
	debugKeyGetResult
)

var LatestHeightBuf = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

var _ HeigthDB = (*HeightDB)(nil)

type HeightDB struct {
	odb         HeightEnabledDB
	readMutex   *sync.RWMutex
	writeMutex  *sync.RWMutex
	writeHeight int64
	readHeight  int64
	config      *HeightDBConfig

	// writeBatch HeightEnabledBatch
}

type HeightDBConfig struct {
	Debug bool
}

func NewHeightDB(db HeightEnabledDB, config *HeightDBConfig) *HeightDB {
	return &HeightDB{
		writeHeight: 0,
		readHeight:  0,
		readMutex:   new(sync.RWMutex),
		writeMutex:  new(sync.RWMutex),
		odb:         db,
		config:      config,
		// writeBatch:  nil,
	}
}

func (hd *HeightDB) BranchHeightDB(height int64) *HeightDB {
	newOne := NewHeightDB(hd.odb, hd.config)
	newOne.SetReadHeight(height)
	return newOne
}

// SetReadHeight sets a target read height in the db driver.
// It acts differently if the db mode is writer or reader:
//   - Reader uses readHeight as the max height at which the retrieved key/value pair is limited to,
//     allowing full block snapshot history
func (hd *HeightDB) SetReadHeight(height int64) {
	hd.readHeight = height
}

// ClearReadHeight sets internal readHeight to LatestHeight
func (hd *HeightDB) ClearReadHeight() int64 {
	lastKnownReadHeight := hd.readHeight
	hd.readHeight = LatestHeight
	return lastKnownReadHeight
}

// GetCurrentReadHeight gets the current readHeight
func (hd *HeightDB) GetCurrentReadHeight() int64 {
	return hd.readHeight
}

// SetWriteHeight sets a target write height in the db driver.
// - Writer uses writeHeight to append along with the key, so later when fetching with the driver
// you can find the latest known key/value pair before the writeHeight
func (hd *HeightDB) SetWriteHeight(height int64) {
	if height != 0 {
		hd.writeHeight = height
		// hd.writeBatch = hd.NewBatch()
	}
}

// ClearWriteHeight sets the next target write Height
// NOTE: evaluate the actual usage of it
func (hd *HeightDB) ClearWriteHeight() int64 {
	fmt.Println("!!! clearing write height...")
	lastKnownWriteHeight := hd.writeHeight
	hd.writeHeight = InvalidHeight
	// if batchErr := hd.writeBatch.Write(); batchErr != nil {
	// 	panic(batchErr)
	// }
	// hd.writeBatch = nil
	return lastKnownWriteHeight
}

// GetCurrentWriteHeight gets the current write height
func (hd *HeightDB) GetCurrentWriteHeight() int64 {
	return hd.writeHeight
}

// Get fetches the value of the given key, or nil if it does not exist.
// CONTRACT: key, value readonly []byte
func (hd *HeightDB) Get(key []byte) ([]byte, error) {
	return hd.odb.Get(hd.GetCurrentReadHeight(), key)
}

// Has checks if a key exists.
// CONTRACT: key, value readonly []byte
func (hd *HeightDB) Has(key []byte) (bool, error) {
	return hd.odb.Has(hd.GetCurrentReadHeight(), key)
}

// Set sets the value for the given key, replacing it if it already exists.
// CONTRACT: key, value readonly []byte
func (hd *HeightDB) Set(key []byte, value []byte) error {
	return hd.odb.Set(hd.writeHeight, key, value)
}

// SetSync sets the value for the given key, and flushes it to storage before returning.
func (hd *HeightDB) SetSync(key []byte, value []byte) error {
	return hd.Set(key, value)
}

// Delete deletes the key, or does nothing if the key does not exist.
// CONTRACT: key readonly []byte
// NOTE(mantlemint): delete should be marked?
func (hd *HeightDB) Delete(key []byte) error {
	return hd.odb.Delete(hd.writeHeight, key)
}

// DeleteSync deletes the key, and flushes the delete to storage before returning.
func (hd *HeightDB) DeleteSync(key []byte) error {
	return hd.Delete(key)
}

// Iterator returns an iterator over a domain of keys, in ascending order. The caller must call
// Close when done. End is exclusive, and start must be less than end. A nil start iterates
// from the first key, and a nil end iterates to the last key (inclusive).
// CONTRACT: No writes may happen within a domain while an iterator exists over it.
// CONTRACT: start, end readonly []byte
func (hd *HeightDB) Iterator(start, end []byte) (tmdb.Iterator, error) {
	return hd.odb.Iterator(hd.GetCurrentReadHeight(), start, end)
}

// ReverseIterator returns an iterator over a domain of keys, in descending order. The caller
// must call Close when done. End is exclusive, and start must be less than end. A nil end
// iterates from the last key (inclusive), and a nil start iterates to the first key (inclusive).
// CONTRACT: No writes may happen within a domain while an iterator exists over it.
// CONTRACT: start, end readonly []byte
func (hd *HeightDB) ReverseIterator(start, end []byte) (tmdb.Iterator, error) {
	return hd.odb.ReverseIterator(hd.GetCurrentReadHeight(), start, end)
}

// Close closes the database connection.
func (hd *HeightDB) Close() error {
	return hd.odb.Close()
}

// NewBatch creates a batch for atomic updates. The caller must call Batch.Close.
func (hd *HeightDB) NewBatch() tmdb.Batch {
	// if hld.writeBatch != nil {
	// 	// TODO: fix me
	// 	return hld.writeBatch
	// } else {
	// 	fmt.Println("!!! opening hld.batch", hld.GetCurrentWriteHeight())
	// 	hld.writeBatch = hld.odb.NewBatch(hld.GetCurrentWriteHeight())
	// 	return hld.writeBatch
	// }
	//
	return hd.odb.NewBatch(hd.GetCurrentWriteHeight())
}

//
// func (hd *HeightDB) FlushBatch() error {
// 	hd.writeBatch
// }

// Print is used for debugging.
func (hd *HeightDB) Print() error {
	return hd.odb.Print()
}

// Stats returns a map of property values for all keys and the size of the cache.
func (hd *HeightDB) Stats() map[string]string {
	return hd.odb.Stats()
}

func (hd *HeightDB) Debug(debugType int, key []byte, value []byte) {
	if !hd.config.Debug {
		return
	}

	keyFamily := key[:len(key)-9]
	keyHeight := key[len(key)-8:]

	var debugPrefix string
	switch debugType {
	case debugKeyGet:
		debugPrefix = "get"
	case debugKeySet:
		debugPrefix = "set"
	case debugKeyIterator:
		debugPrefix = "get/it"
	case debugKeyReverseIterator:
		debugPrefix = "get/rit"

	case debugKeyGetResult:
		debugPrefix = "get/response"
	}

	var actualKeyHeight int64
	if bytes.Compare(keyHeight, LatestHeightBuf) == 0 {
		actualKeyHeight = -1
	} else {
		actualKeyHeight = int64(utils.BigEndianToUint(keyHeight))
	}

	fmt.Printf("<%s @ %d> %s", debugPrefix, actualKeyHeight, keyFamily)
	fmt.Printf("\n")
}
