package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateMinimumCommissionEnforced(t *testing.T) {
	require.Error(t, validateMiniumCommissionEnforced(int64(3)))
	require.NoError(t, validateMiniumCommissionEnforced(true))
	require.NoError(t, validateMiniumCommissionEnforced(false))
}

func TestValidateMinimumCommission(t *testing.T) {
	require.Error(t, validateMinimumCommission(int64(3)))
	require.Error(t, validateMinimumCommission(sdk.NewDec(2)))
	require.Error(t, validateMinimumCommission(sdk.NewDec(-2)))
	require.NoError(t, validateMinimumCommission(sdk.NewDecWithPrec(5, 2)))
}

func TestValidateParams(t *testing.T) {
	require.Error(t, NewParams(false, sdk.NewDec(2)).Validate())
	require.Error(t, NewParams(false, sdk.NewDec(-2)).Validate())
	require.NoError(t, NewParams(false, sdk.NewDecWithPrec(2, 2)).Validate())
	require.NoError(t, NewParams(false, sdk.ZeroDec()).Validate())
	require.NoError(t, NewParams(false, sdk.OneDec()).Validate())
}
