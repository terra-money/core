package custom_queriers

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	alliancebindings "github.com/terra-money/alliance/x/alliance/bindings"
	"github.com/terra-money/alliance/x/alliance/bindings/types"
	"github.com/terra-money/core/v2/x/tokenfactory/bindings"
	types2 "github.com/terra-money/core/v2/x/tokenfactory/bindings/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func AlwaysErrorQuerier(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return nil, fmt.Errorf("always error")
}

func AlwaysUnknownQuerier(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return nil, fmt.Errorf("unknown query")
}

func AlwaysGoodQuerier(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return []byte("good"), nil
}

func TestCustomQueriers(t *testing.T) {
	querier := CustomQueriers(AlwaysUnknownQuerier, AlwaysErrorQuerier, AlwaysGoodQuerier)
	_, err := querier(sdk.Context{}, nil)
	require.ErrorContainsf(t, err, "always error", "")

	querier = CustomQueriers(AlwaysUnknownQuerier, AlwaysGoodQuerier, AlwaysErrorQuerier)
	_, err = querier(sdk.Context{}, nil)
	require.NoError(t, err)
}

func TestWithTfAndAllianceButCallAlliance(t *testing.T) {
	tfQuerier := bindings.CustomQuerier(&bindings.QueryPlugin{})
	allianceQuerier := alliancebindings.CustomQuerier(&alliancebindings.QueryPlugin{})
	querier := CustomQueriers(tfQuerier, allianceQuerier)

	query := types.AllianceQuery{
		Alliance: &types.Alliance{Denom: "123"},
	}
	bz, err := json.Marshal(query)
	require.NoError(t, err)

	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 1024)
			runtime.Stack(stack, false)
			// We make sure alliance is called here
			require.Containsf(t, string(stack), "GetAlliance", "")
		}
	}()

	// We call querier but it will panic because we don't have a keeper
	_, err = querier(sdk.Context{}, bz)
	require.Fail(t, "should panic")
}

func TestWithTfAndAllianceButCallTf(t *testing.T) {
	tfQuerier := bindings.CustomQuerier(&bindings.QueryPlugin{})
	allianceQuerier := alliancebindings.CustomQuerier(&alliancebindings.QueryPlugin{})
	querier := CustomQueriers(tfQuerier, allianceQuerier)

	query := types2.TokenFactoryQuery{
		Token: &types2.TokenQuery{
			Params: &types2.GetParams{},
		},
	}
	bz, err := json.Marshal(query)
	require.NoError(t, err)

	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 1024)
			runtime.Stack(stack, false)
			// We make sure tf is called here
			require.Containsf(t, string(stack), "GetParams", "")
		}
	}()

	// We call querier but it will panic because we don't have a keeper
	_, err = querier(sdk.Context{}, bz)
	require.Fail(t, "should panic")
}

func TestWithTfAndAllianceButRandomCall(t *testing.T) {
	tfQuerier := bindings.CustomQuerier(&bindings.QueryPlugin{})
	allianceQuerier := alliancebindings.CustomQuerier(&alliancebindings.QueryPlugin{})
	querier := CustomQueriers(tfQuerier, allianceQuerier)

	query := sdk.NewCoin("denom", sdk.NewInt(1))
	bz, err := json.Marshal(query)
	require.NoError(t, err)

	// We call querier but it will panic because we don't have a keeper
	_, err = querier(sdk.Context{}, bz)
	require.Error(t, err)
}

func TestRegisterCustomPlugins(t *testing.T) {
	options := RegisterCustomPlugins(nil, nil, nil)
	require.Len(t, options, 2)
}
