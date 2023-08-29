package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/terra-money/core/v2/x/tokenfactory/types"
)

func MigrateStore(ctx sdk.Context, subspace paramtypes.Subspace) error {
	var params types.Params
	subspace.GetParamSet(ctx, &params)
	params.DenomCreationGasConsume = types.DefaultCreationGasFee
	subspace.SetParamSet(ctx, &params)
	return nil
}
