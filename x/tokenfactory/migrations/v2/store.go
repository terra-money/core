package v2

import (
	"github.com/terra-money/core/v2/app/config"
	"github.com/terra-money/core/v2/x/tokenfactory/exported"
	"github.com/terra-money/core/v2/x/tokenfactory/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(ctx sdk.Context, legacySubspace exported.Subspace, cdc codec.BinaryCodec) error {
	var params types.Params
	legacySubspace.GetParamSetIfExists(ctx, &params)

	// FIX: when token factory was implemented for the first time *denom creation fee* field was setup to
	// nil which makes this migration fails. This if statement will fix the issue:
	// https://github.com/terra-money/core/blob/a03a0657c7430d32e6329d86de78bb4aab9a9aa7/app/app.go#L1053
	if params.DenomCreationFee == nil {
		params.DenomCreationFee = sdk.NewCoins(sdk.NewCoin(config.BondDenom, sdk.NewInt(10_000_000)))
	}

	params.DenomCreationGasConsume = types.DefaultCreationGasFee
	legacySubspace.SetParamSet(ctx, &params)
	return nil
}
