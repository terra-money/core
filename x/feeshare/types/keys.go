package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// module name
	ModuleName = "feeshare"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// prefix bytes for the fees persistent store
const (
	prefixFeeShare = iota + 1
	prefixDeployer
	prefixWithdrawer
	prefixParams
)

// KVStore key prefixes
var (
	KeyPrefixFeeShare   = []byte{prefixFeeShare}
	KeyPrefixDeployer   = []byte{prefixDeployer}
	KeyPrefixWithdrawer = []byte{prefixWithdrawer}
	ParamsKey           = []byte{prefixParams}
)

// GetKeyPrefixDeployer returns the KVStore key prefix for storing
// registered feeshare contract for a deployer
func GetKeyPrefixDeployer(deployerAddress sdk.AccAddress) []byte {
	return append(KeyPrefixDeployer, deployerAddress.Bytes()...)
}

// GetKeyPrefixWithdrawer returns the KVStore key prefix for storing
// registered feeshare contract for a withdrawer
func GetKeyPrefixWithdrawer(withdrawerAddress sdk.AccAddress) []byte {
	return append(KeyPrefixWithdrawer, withdrawerAddress.Bytes()...)
}
