package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
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

// ConfigureBondDenom generates the default state for the application.
func (genState GenesisState) ConfigureBondDenom(cdc codec.JSONCodec, bondDenom string) GenesisState {
	// customize bond denom
	var stakingGenState stakingtypes.GenesisState
	cdc.MustUnmarshalJSON(genState[stakingtypes.ModuleName], &stakingGenState)
	stakingGenState.Params.BondDenom = bondDenom
	genState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingGenState)

	var crisisGenState crisistypes.GenesisState
	cdc.MustUnmarshalJSON(genState[crisistypes.ModuleName], &crisisGenState)
	crisisGenState.ConstantFee.Denom = bondDenom
	genState[crisistypes.ModuleName] = cdc.MustMarshalJSON(&crisisGenState)

	var govGenState govtypes.GenesisState
	cdc.MustUnmarshalJSON(genState[govtypes.ModuleName], &govGenState)
	govGenState.DepositParams.MinDeposit[0].Denom = bondDenom
	genState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenState)

	var mintGenState minttypes.GenesisState
	cdc.MustUnmarshalJSON(genState[minttypes.ModuleName], &mintGenState)
	mintGenState.Params.MintDenom = bondDenom
	genState[minttypes.ModuleName] = cdc.MustMarshalJSON(&mintGenState)

	return genState
}

func (genState GenesisState) ConfigureICA(cdc codec.JSONCodec) GenesisState {
	// create ICS27 Controller submodule params
	controllerParams := icacontrollertypes.Params{}

	// create ICS27 Host submodule params
	hostParams := icahosttypes.Params{
		HostEnabled: true,
		AllowMessages: []string{
			authzMsgExec,
			authzMsgGrant,
			authzMsgRevoke,
			bankMsgSend,
			bankMsgMultiSend,
			distrMsgSetWithdrawAddr,
			distrMsgWithdrawValidatorCommission,
			distrMsgFundCommunityPool,
			distrMsgWithdrawDelegatorReward,
			feegrantMsgGrantAllowance,
			feegrantMsgRevokeAllowance,
			govMsgVoteWeighted,
			govMsgSubmitProposal,
			govMsgDeposit,
			govMsgVote,
			stakingMsgEditValidator,
			stakingMsgDelegate,
			stakingMsgUndelegate,
			stakingMsgBeginRedelegate,
			stakingMsgCreateValidator,
			vestingMsgCreateVestingAccount,
			transferMsgTransfer,
			wasmMsgStoreCode,
			wasmMsgInstantiateContract,
			wasmMsgExecuteContract,
			wasmMsgMigrateContract,
		},
	}

	var icaGenState icatypes.GenesisState
	cdc.MustUnmarshalJSON(genState[icatypes.ModuleName], &icaGenState)
	icaGenState.ControllerGenesisState.Params = controllerParams
	icaGenState.HostGenesisState.Params = hostParams
	genState[icatypes.ModuleName] = cdc.MustMarshalJSON(&icaGenState)

	return genState
}
