package snappy

import (
	"io/ioutil"
	"os"
	"testing"

	db "github.com/cometbft/cometbft-db"
	tmjson "github.com/cometbft/cometbft/libs/json"
	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/assert"
)

func TestSnappyDB(t *testing.T) {
	snappy := NewSnappyDB(db.NewMemDB(), CompatModeEnabled)

	assert.Nil(t, snappy.Set([]byte("test"), []byte("testValue")))

	var v []byte
	var err error

	// nil buffer test
	v, err = snappy.Get([]byte("non-existing"))
	assert.Nil(t, v)
	assert.Nil(t, err)

	v, err = snappy.Get([]byte("test"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("testValue"), v)

	assert.Nil(t, snappy.Delete([]byte("test")))
	v, err = snappy.Get([]byte("test"))
	assert.Nil(t, v)
	assert.Nil(t, err)

	// iterator is not supported
	var it db.Iterator
	it, err = snappy.Iterator([]byte("start"), []byte("end"))
	assert.Nil(t, it)
	assert.Equal(t, errIteratorNotSupported, err)

	it, err = snappy.ReverseIterator([]byte("start"), []byte("end"))
	assert.Nil(t, it)
	assert.Equal(t, errIteratorNotSupported, err)

	// batched store is compressed as well
	var batch db.Batch
	batch = snappy.NewBatch()

	assert.Nil(t, batch.Set([]byte("key"), []byte("batchedValue")))
	assert.Nil(t, batch.Write())
	assert.Nil(t, batch.Close())

	v, err = snappy.Get([]byte("key"))
	assert.Equal(t, []byte("batchedValue"), v)

	batch = snappy.NewBatch()
	assert.Nil(t, batch.Delete([]byte("key")))
	assert.Nil(t, batch.Write())
	assert.Nil(t, batch.Close())

	v, err = snappy.Get([]byte("key"))
	assert.Nil(t, v)
	assert.Nil(t, err)
}

func TestSnappyDBCompat(t *testing.T) {
	mdb := db.NewMemDB()
	testKey := []byte("testKey")

	nocompat := NewSnappyDB(mdb, CompatModeDisabled)
	indexSampleTx(nocompat, testKey)

	nocompatResult, _ := nocompat.Get(testKey)

	compat := NewSnappyDB(mdb, CompatModeEnabled)
	compatResult, _ := compat.Get(testKey)
	assert.Equal(t, nocompatResult, compatResult)

	nocompatResult2, _ := nocompat.Get(testKey)
	assert.Equal(t, compatResult, nocompatResult2)
}

func indexSampleTx(mdb db.DB, key []byte) {
	block := &cometbfttypes.Block{}
	blockFile, _ := os.Open("../../indexer/fixtures/block_4814775.json")
	blockJSON, _ := ioutil.ReadAll(blockFile)
	if err := tmjson.Unmarshal(blockJSON, block); err != nil {
		panic(err)
	}

	_ = mdb.Set(key, blockJSON)
}
