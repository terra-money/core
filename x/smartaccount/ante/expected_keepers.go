package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	smartaccounttypes "github.com/terra-money/core/v2/x/smartaccount/types"
)

type SmartAccountKeeper interface {
	GetSetting(ctx sdk.Context, ownerAddr string) (*smartaccounttypes.Setting, error)
}

type WasmKeeper interface {
	Sudo(ctx sdk.Context, contractAddress sdk.AccAddress, msg []byte) ([]byte, error)
}
