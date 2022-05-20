package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GenesisState - The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	return ModuleBasics.DefaultGenesis(cdc)
}

// ConvertBondDenom generates the default state for the application.
func (genState GenesisState) ConvertBondDenom(cdc codec.JSONCodec) GenesisState {
	// customize bond denom
	var stakingGenState stakingtypes.GenesisState
	cdc.MustUnmarshalJSON(genState[stakingtypes.ModuleName], &stakingGenState)
	stakingGenState.Params.BondDenom = BondDenom
	genState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingGenState)

	var crisisGenState crisistypes.GenesisState
	cdc.MustUnmarshalJSON(genState[crisistypes.ModuleName], &crisisGenState)
	crisisGenState.ConstantFee.Denom = BondDenom
	genState[crisistypes.ModuleName] = cdc.MustMarshalJSON(&crisisGenState)

	var govGenState govtypes.GenesisState
	cdc.MustUnmarshalJSON(genState[govtypes.ModuleName], &govGenState)
	govGenState.DepositParams.MinDeposit[0].Denom = BondDenom
	genState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenState)

	var mintGenState minttypes.GenesisState
	cdc.MustUnmarshalJSON(genState[minttypes.ModuleName], &mintGenState)
	mintGenState.Params.MintDenom = BondDenom
	genState[minttypes.ModuleName] = cdc.MustMarshalJSON(&mintGenState)

	return genState
}
