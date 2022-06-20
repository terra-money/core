package types

const (
	// ModuleName is the name of the ante module
	ModuleName = "ante"

	// StoreKey is the store key string for ante
	StoreKey = ModuleName

	// RouterKey is the router key string for ante
	RouterKey = ModuleName
)

// Keys for ante store
// Items are stored with the following key: values
//
// - 0x00: MinimumCommission
var (
	MinimumCommissionKey = []byte{0x00} // key for minimum commission
)
