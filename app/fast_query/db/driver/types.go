package driver

import (
	"encoding/binary"
	"math"

	"github.com/terra-money/core/v2/app/fast_query/db/utils"
)

const (
	DriverModeKeySuffixAsc  = iota
	DriverModeKeySuffixDesc = 1
	DBName                  = "fast_query"
)

var (
	cCurrentDataPrefix     = []byte{0}
	cKeysForIteratorPrefix = []byte{1}
	cDataWithHeightPrefix  = []byte{2}
)

func prefixCurrentDataKey(key []byte) []byte {
	return append(cCurrentDataPrefix, key...)
}

func prefixKeysForIteratorKey(key []byte) []byte {
	return append(cKeysForIteratorPrefix, key...)
}

func prefixDataWithHeightKey(key []byte) []byte {
	result := make([]byte, 0, len(cDataWithHeightPrefix)+len(key))
	result = append(result, cDataWithHeightPrefix...)
	result = append(result, key...)
	return result
}

func serializeHeight(mode int, height int64) []byte {
	if mode == DriverModeKeySuffixAsc {
		return utils.UintToBigEndian(uint64(height))
	} else {
		return utils.UintToBigEndian(math.MaxUint64 - uint64(height))
	}
}

func deserializeHeight(mode int, data []byte) int64 {
	if mode == DriverModeKeySuffixAsc {
		return int64(binary.BigEndian.Uint64(data))
	} else {
		return int64(math.MaxUint64 - binary.BigEndian.Uint64(data))
	}
}
