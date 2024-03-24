package types

const (
	// module name
	ModuleName = "feeburn"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// prefix bytes for the fees persistent store
const prefixParams = 1

// KVStore key prefixes
var (
	ParamsKey = []byte{prefixParams}
)
