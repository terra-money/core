package bindings_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/app"
	bindings "github.com/terra-money/core/v2/x/tokenfactory/bindings/types"
)

func TestQuery(t *testing.T) {
	// Setup the environment and fund the user accounts
	user := RandomAccountAddress()
	app, ctx := SetupCustomApp(t, user)
	reflect := instantiateReflectContract(t, ctx, app, user)
	require.NotEmpty(t, reflect)
	fundAccount(t, ctx, app, reflect, sdk.Coins{sdk.NewInt64Coin("uluna", 100_000_000_000)})

	// Create ustart and ustart2 denoms thoguht the smart contract to
	// query and validate the query binding are working as expected
	msg := bindings.TokenMsg{CreateDenom: &bindings.CreateDenom{
		Subdenom: "ustart",
	}}
	err := executeCustom(t, ctx, app, reflect, user, msg, sdk.Coin{})
	require.NoError(t, err)
	msg = bindings.TokenMsg{CreateDenom: &bindings.CreateDenom{
		Subdenom: "ustart2",
	}}
	err = executeCustom(t, ctx, app, reflect, user, msg, sdk.Coin{})
	require.NoError(t, err)

	// Query params info
	query := bindings.TokenQuery{
		Params: &bindings.GetParams{},
	}
	paramsRes := bindings.ParamsResponse{}
	err = queryCustom(t, ctx, app, reflect, query, &paramsRes)
	require.NoError(t, err)

	require.EqualValues(t, bindings.ParamsResponse{
		Params: bindings.Params{
			DenomCreationFee: []wasmvmtypes.Coin{
				{
					Denom:  "uluna",
					Amount: "10000000",
				},
			},
		},
	}, paramsRes)

	// Query full denom name thought wasm binding
	query = bindings.TokenQuery{
		FullDenom: &bindings.FullDenom{
			CreatorAddr: reflect.String(),
			Subdenom:    "ustart",
		},
	}
	fulldenomresp := bindings.FullDenomResponse{}
	err = queryCustom(t, ctx, app, reflect, query, &fulldenomresp)
	require.NoError(t, err)

	require.EqualValues(t,
		fmt.Sprintf("factory/%s/ustart", reflect.String()),
		fulldenomresp.Denom,
	)

	// Query metadata thoguht wasm binding
	query = bindings.TokenQuery{
		Metadata: &bindings.GetMetadata{
			Denom: fmt.Sprintf("factory/%s/ustart", reflect.String()),
		},
	}
	metadataRes := bindings.MetadataResponse{}
	err = queryCustom(t, ctx, app, reflect, query, &metadataRes)
	require.NoError(t, err)

	require.EqualValues(t, bindings.MetadataResponse{
		Metadata: &bindings.Metadata{
			Description: "",
			Base:        fmt.Sprintf("factory/%s/ustart", reflect.String()),
			Display:     "",
			Name:        "",
			Symbol:      "",
			DenomUnits: []bindings.DenomUnit{
				{
					Denom:    fmt.Sprintf("factory/%s/ustart", reflect.String()),
					Exponent: 0,
					Aliases:  nil,
				},
			},
		},
	}, metadataRes)

	// Query denom admin thoguht wasm binding
	query = bindings.TokenQuery{
		Admin: &bindings.DenomAdmin{
			Denom: fmt.Sprintf("factory/%s/ustart", reflect.String()),
		},
	}
	adminresp := bindings.AdminResponse{}
	err = queryCustom(t, ctx, app, reflect, query, &adminresp)
	require.NoError(t, err)

	require.EqualValues(t, reflect.String(), adminresp.Admin)

	// Query all denoms by user thoguht wasm binding
	query = bindings.TokenQuery{
		DenomsByCreator: &bindings.DenomsByCreator{
			Creator: reflect.String(),
		},
	}
	denomsbycreator := bindings.DenomsByCreatorResponse{}
	err = queryCustom(t, ctx, app, reflect, query, &denomsbycreator)
	require.NoError(t, err)

	expected := []string{
		fmt.Sprintf("factory/%s/ustart", reflect.String()),
		fmt.Sprintf("factory/%s/ustart2", reflect.String()),
	}
	require.EqualValues(t, expected, denomsbycreator.Denoms)
}

type ReflectQuery struct {
	Chain *ChainRequest `json:"chain,omitempty"`
}

type ChainRequest struct {
	Request wasmvmtypes.QueryRequest `json:"request"`
}

type ChainResponse struct {
	Data []byte `json:"data"`
}

func queryCustom(t *testing.T, ctx sdk.Context, app *app.TerraApp, contract sdk.AccAddress, request bindings.TokenQuery, response interface{}) error {
	t.Helper()
	wrapped := bindings.TokenFactoryQuery{
		Token: &request,
	}
	msgBz, err := json.Marshal(wrapped)
	if err != nil {
		return err
	}

	query := ReflectQuery{
		Chain: &ChainRequest{
			Request: wasmvmtypes.QueryRequest{Custom: msgBz},
		},
	}
	queryBz, err := json.Marshal(query)
	if err != nil {
		return err
	}

	resBz, err := app.WasmKeeper.QuerySmart(ctx, contract, queryBz)
	if err != nil {
		return err
	}

	var resp ChainResponse
	err = json.Unmarshal(resBz, &resp)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp.Data, response)
	if err != nil {
		return err
	}

	return nil
}
