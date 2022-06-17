package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateGenesis(t *testing.T) {
	require.Error(t, NewGenesisState(NewParams(false, sdk.NewDec(2))).Validate())
	require.Error(t, NewGenesisState(NewParams(false, sdk.NewDec(-2))).Validate())
	require.NoError(t, NewGenesisState(NewParams(false, sdk.NewDecWithPrec(2, 2))).Validate())
	require.NoError(t, NewGenesisState(NewParams(false, sdk.ZeroDec())).Validate())
	require.NoError(t, NewGenesisState(NewParams(false, sdk.OneDec())).Validate())
}
