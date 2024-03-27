package wasm

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	alliancebindings "github.com/terra-money/alliance/x/alliance/bindings"
	"github.com/terra-money/alliance/x/alliance/bindings/types"
	"github.com/terra-money/core/v2/x/tokenfactory/bindings"
	types2 "github.com/terra-money/core/v2/x/tokenfactory/bindings/types"
	"runtime"
	"testing"
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
			require.Containsf(t, string(stack), "keeper.Keeper.GetAssetByDenom", "")
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
			require.Containsf(t, string(stack), "bindings.QueryPlugin.GetParams", "")
		}
	}()

	// We call querier but it will panic because we don't have a keeper
	_, err = querier(sdk.Context{}, bz)
	require.Fail(t, "should panic")
}
