package types

const (
	// module name
	ModuleName = "smartaccount"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// prefix bytes for the fees persistent store
const (
	prefixParams = iota + 1
	prefixSetting
)

// KVStore key prefixes
var (
	ParamsKey        = []byte{prefixParams}
	KeyPrefixSetting = []byte{prefixSetting}
)

// GetKeyPrefixSetting returns the KVStore key prefix for storing
// registered smartaccount contract for a deployer
func GetKeyPrefixSetting(ownderAddr string) []byte {
	return append(KeyPrefixSetting, []byte(ownderAddr)...)
}
