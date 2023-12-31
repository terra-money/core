package app

import (
	"encoding/json"

	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icagenesistypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/genesis/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	"github.com/terra-money/core/v2/app/config"
	tokenfactorytypes "github.com/terra-money/core/v2/x/tokenfactory/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
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

// SetDefaultTerraConfig generates the default state for Terra's Core application.
func (genState GenesisState) SetDefaultTerraConfig(cdc codec.JSONCodec) GenesisState {
	// customize bond denom
	var stakingGenState stakingtypes.GenesisState
	cdc.MustUnmarshalJSON(genState[stakingtypes.ModuleName], &stakingGenState)
	stakingGenState.Params.BondDenom = config.BondDenom
	genState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingGenState)

	var crisisGenState crisistypes.GenesisState
	cdc.MustUnmarshalJSON(genState[crisistypes.ModuleName], &crisisGenState)
	crisisGenState.ConstantFee.Denom = config.BondDenom
	genState[crisistypes.ModuleName] = cdc.MustMarshalJSON(&crisisGenState)

	var govGenState govtypesv1.GenesisState
	cdc.MustUnmarshalJSON(genState[govtypes.ModuleName], &govGenState)
	govGenState.Params.MinDeposit[0].Denom = config.BondDenom
	genState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenState)

	var mintGenState minttypes.GenesisState
	cdc.MustUnmarshalJSON(genState[minttypes.ModuleName], &mintGenState)
	mintGenState.Params.MintDenom = config.BondDenom
	genState[minttypes.ModuleName] = cdc.MustMarshalJSON(&mintGenState)

	var tokenFactoryGenState tokenfactorytypes.GenesisState
	cdc.MustUnmarshalJSON(genState[tokenfactorytypes.ModuleName], &tokenFactoryGenState)
	tokenFactoryGenState.Params.DenomCreationFee = sdk.NewCoins(sdk.NewCoin(config.BondDenom, sdk.NewInt(10000000)))
	genState[tokenfactorytypes.ModuleName] = cdc.MustMarshalJSON(&tokenFactoryGenState)

	var icqGenState icqtypes.GenesisState
	cdc.MustUnmarshalJSON(genState[icqtypes.ModuleName], &icqGenState)
	icqGenState.Params.HostEnabled = true
	icqGenState.Params.AllowQueries = icqAllowedQueries()
	genState[icqtypes.ModuleName] = cdc.MustMarshalJSON(&icqGenState)

	var icaGenState icagenesistypes.GenesisState
	cdc.MustUnmarshalJSON(genState[icatypes.ModuleName], &icaGenState)
	icaGenState.ControllerGenesisState.Params = icacontrollertypes.Params{
		ControllerEnabled: true,
	}
	icaGenState.HostGenesisState.Params = icahosttypes.Params{
		HostEnabled:   true,
		AllowMessages: icaAllowedMsgs(),
	}
	genState[icatypes.ModuleName] = cdc.MustMarshalJSON(&icaGenState)

	return genState
}

func icaAllowedMsgs() []string {
	return []string{
		config.AuthzMsgExec,
		config.AuthzMsgGrant,
		config.AuthzMsgRevoke,
		config.BankMsgSend,
		config.BankMsgMultiSend,
		config.DistrMsgSetWithdrawAddr,
		config.DistrMsgWithdrawValidatorCommission,
		config.DistrMsgFundCommunityPool,
		config.DistrMsgWithdrawDelegatorReward,
		config.FeegrantMsgGrantAllowance,
		config.FeegrantMsgRevokeAllowance,
		config.GovMsgVoteWeighted,
		config.GovMsgSubmitProposal,
		config.GovMsgDeposit,
		config.GovMsgVote,
		config.StakingMsgEditValidator,
		config.StakingMsgDelegate,
		config.StakingMsgUndelegate,
		config.StakingMsgBeginRedelegate,
		config.StakingMsgCreateValidator,
		config.VestingMsgCreateVestingAccount,
		config.TransferMsgTransfer,
		config.WasmMsgStoreCode,
		config.WasmMsgInstantiateContract,
		config.WasmMsgExecuteContract,
		config.WasmMsgMigrateContract,
	}
}

func icqAllowedQueries() []string {
	return []string{
		config.QueryAllAllianceValidators,
		config.QueryAllAlliancesDelegations,
		config.QueryAlliance,
		config.QueryAllianceDelegation,
		config.QueryAllianceDelegationRewards,
		config.QueryAllianceRedelegations,
		config.QueryAllianceUnbondings,
		config.QueryAllianceUnbondingsByDenomAndDelegator,
		config.QueryAllianceValidator,
		config.QueryAlliances,
		config.QueryAlliancesDelegation,
		config.QueryAlliancesDelegationByValidator,
		config.QueryIBCAlliance,
		config.QueryIBCAllianceDelegation,
		config.QueryIBCAllianceDelegationRewards,
		config.QueryAllianceParams,
		config.QueryAccount,
		config.QueryAccountAddressByID,
		config.QueryAccountInfo,
		config.QueryAccounts,
		config.QueryAddressBytesToString,
		config.QueryAddressStringToBytes,
		config.QueryBech32Prefix,
		config.QueryModuleAccountByName,
		config.QueryModuleAccounts,
		config.QueryAuthParams,
		config.QueryGranteeGrants,
		config.QueryGranterGrants,
		config.QueryGrants,
		config.QueryAllBalances,
		config.QueryBalance,
		config.QueryDenomMetadata,
		config.QueryDenomOwners,
		config.QueryDenomsMetadata,
		config.QueryBankParams,
		config.QuerySendEnabled,
		config.QuerySpendableBalanceByDenom,
		config.QuerySpendableBalances,
		config.QuerySupplyOf,
		config.QueryTotalSupply,
		config.QueryConsensusParams,
		config.QueryCommunityPool,
		config.QueryDelegationRewards,
		config.QueryDelegationTotalRewards,
		config.QueryDistributionDelegatorValidators,
		config.QueryDelegatorWithdrawAddress,
		config.QueryDistributionParams,
		config.QueryValidatorCommission,
		config.QueryValidatorDistributionInfo,
		config.QueryValidatorOutstandingRewards,
		config.QueryValidatorSlashes,
		config.QueryAllEvidence,
		config.QueryEvidence,
		config.QueryAllowance,
		config.QueryAllowances,
		config.QueryAllowancesByGranter,
		config.QueryDeposit,
		config.QueryDeposits,
		config.QueryGovParams,
		config.QueryProposal,
		config.QueryProposals,
		config.QueryTallyResult,
		config.QueryVote,
		config.QueryVotes,
		config.QueryAnnualProvisions,
		config.QueryInflation,
		config.QueryMintParams,
		config.QueryParamsModuleParams,
		config.QuerySubspaces,
		config.QuerySlashingParams,
		config.QuerySigningInfo,
		config.QuerySigningInfos,
		config.QueryDelegation,
		config.QueryDelegatorDelegations,
		config.QueryDelegatorUnbondingDelegations,
		config.QueryDelegatorValidator,
		config.QueryDelegatorValidators,
		config.QueryHistoricalInfo,
		config.QueryStakingParams,
		config.QueryPool,
		config.QueryRedelegations,
		config.QueryUnbondingDelegation,
		config.QueryValidator,
		config.QueryValidatorDelegations,
		config.QueryValidatorUnbondingDelegations,
		config.QueryValidators,
		config.QueryAppliedPlan,
		config.QueryAuthority,
		config.QueryCurrentPlan,
		config.QueryModuleVersions,
		config.QueryUpgradedConsensusState,
		config.QueryAllContractState,
		config.QueryCode,
		config.QueryCodes,
		config.QueryContractHistory,
		config.QueryContractInfo,
		config.QueryContractsByCode,
		config.QueryContractsByCreator,
		config.QueryWasmParams,
		config.QueryPinnedCodes,
		config.QueryRawContractState,
		config.QuerySmartContractState,
		config.QueryCounterpartyPayee,
		config.QueryFeeEnabledChannel,
		config.QueryFeeEnabledChannels,
		config.QueryIncentivizedPacket,
		config.QueryIncentivizedPackets,
		config.QueryIncentivizedPacketsForChannel,
		config.QueryPayee,
		config.QueryTotalAckFees,
		config.QueryTotalRecvFees,
		config.QueryTotalTimeoutFees,
		config.QueryInterchainAccount,
		config.QueryInterchainAccControllerParams,
		config.QueryInterchainAccHostParams,
		config.QueryDenomHash,
		config.QueryDenomTrace,
		config.QueryDenomTraces,
		config.QueryEscrowAddress,
		config.QueryTransferParams,
		config.QueryTotalEscrowForDenom,
		config.QueryChannel,
		config.QueryChannelClientState,
		config.QueryChannelConsensusState,
		config.QueryChannels,
		config.QueryConnectionChannels,
		config.QueryNextSequenceReceive,
		config.QueryPacketAcknowledgement,
		config.QueryPacketAcknowledgements,
		config.QueryPacketCommitment,
		config.QueryPacketCommitments,
		config.QueryPacketReceipt,
		config.QueryUnreceivedAcks,
		config.QueryUnreceivedPackets,
		config.QueryClientParams,
		config.QueryClientState,
		config.QueryClientStates,
		config.QueryClientStatus,
		config.QueryConsensusState,
		config.QueryConsensusStateHeights,
		config.QueryConsensusStates,
		config.QueryUpgradedClientState,
		config.QueryUpgradedConsensusStateV1,
		config.QueryClientConnections,
		config.QueryConnection,
		config.QueryConnectionClientState,
		config.QueryConnectionConsensusState,
		config.QueryConnectionParams,
		config.QueryConnections,
		config.QueryICQParams,
		config.QueryDeployerFeeShares,
		config.QueryFeeShare,
		config.QueryFeeShares,
		config.QueryFeeshareParams,
		config.QueryWithdrawerFeeShares,
		config.QueryBeforeSendHookAddress,
		config.QueryDenomAuthorityMetadata,
		config.QueryDenomsFromCreator,
		config.QueryTokeFactoryParams,
		config.QueryPFMParams,
	}
}
