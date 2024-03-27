package wasm

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type Querier func(ctx sdk.Context, request json.RawMessage) ([]byte, error)

func CustomQueriers(queriers ...Querier) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		for _, querier := range queriers {
			res, err := querier(ctx, request)
			if err == nil || !strings.Contains(err.Error(), "unknown query") {
				return res, err
			}
		}
		return nil, fmt.Errorf("unknown query")
	}
}
