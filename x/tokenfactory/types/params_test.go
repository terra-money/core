package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/terra-money/core/v2/x/tokenfactory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestCreateNewParms(t *testing.T) {
	// Creaate new params
	params := types.NewParams(sdk.NewCoins(sdk.NewCoin("uluna", math.NewInt(100000))), 10)
	new_expected_params := types.Params{
		DenomCreationFee:        sdk.NewCoins(sdk.NewCoin("uluna", math.NewInt(100000))),
		DenomCreationGasConsume: 10,
	}
	require.Equal(t, new_expected_params, params)

	// Validate params set creation and validate they are different than the default
	paramSetPairs := params.ParamSetPairs()
	require.Equal(t, 2, len(paramSetPairs))
	defaultKeyTable := types.ParamKeyTable()
	require.NotEqual(t, defaultKeyTable, paramSetPairs)
}
