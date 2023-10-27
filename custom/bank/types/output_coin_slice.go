package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Define a custom type to hold Output and Coin together for sorting
type OutputCoin struct {
	Output banktypes.Output
	Coin   sdk.Coin
}

// Define a slice of OutputCoin and implement sort.Interface
type OutputCoinSlice []OutputCoin

func (o OutputCoinSlice) Len() int           { return len(o) }
func (o OutputCoinSlice) Less(i, j int) bool { return o[i].Coin.Denom < o[j].Coin.Denom }
func (o OutputCoinSlice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
