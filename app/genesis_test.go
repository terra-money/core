package app_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/v2/app"
)

func TestGenesis(t *testing.T) {
	encCfg := app.MakeEncodingConfig()
	genesisState := app.NewDefaultGenesisState(encCfg.Marshaler)
	genesisState.SetDefaultTerraConfig(encCfg.Marshaler)

	jsonGenState, err := json.Marshal(genesisState)
	require.Nil(t, err)

	expectedState := `{
		"06-solomachine": null,
		"07-tendermint": null,
		"alliance": {
			"params": {
				"reward_delay_time": "604800s",
				"take_rate_claim_interval": "300s",
				"last_take_rate_claim_time": "0001-01-01T00:00:00Z"
			},
			"assets": [],
			"validator_infos": [],
			"reward_weight_change_snaphots": [],
			"delegations": [],
			"redelegations": [],
			"undelegations": []
		},
		"auth": {
			"params": {
				"max_memo_characters": "256",
				"tx_sig_limit": "7",
				"tx_size_cost_per_byte": "10",
				"sig_verify_cost_ed25519": "590",
				"sig_verify_cost_secp256k1": "1000"
			},
			"accounts": []
		},
		"authz": {
			"authorization": []
		},
		"bank": {
			"params": {
				"send_enabled": [],
				"default_send_enabled": true
			},
			"balances": [],
			"supply": [],
			"denom_metadata": [],
			"send_enabled": []
		},
		"builder": {
			"params": {
				"max_bundle_size": 2,
				"escrow_account_address": "32sHF2qbF8xMmvwle9QEcy59Cbc=",
				"reserve_fee": {
					"denom": "uluna",
					"amount": "1"
				},
				"min_bid_increment": {
					"denom": "uluna",
					"amount": "1"
				},
				"front_running_protection": true,
				"proposer_fee": "0.000000000000000000"
			}
		},
		"capability": {
			"index": "1",
			"owners": []
		},
		"consensus": null,
		"crisis": {
			"constant_fee": {
				"denom": "uluna",
				"amount": "1000"
			}
		},
		"distribution": {
			"params": {
				"community_tax": "0.020000000000000000",
				"base_proposer_reward": "0.000000000000000000",
				"bonus_proposer_reward": "0.000000000000000000",
				"withdraw_addr_enabled": true
			},
			"fee_pool": {
				"community_pool": []
			},
			"delegator_withdraw_infos": [],
			"previous_proposer": "",
			"outstanding_rewards": [],
			"validator_accumulated_commissions": [],
			"validator_historical_rewards": [],
			"validator_current_rewards": [],
			"delegator_starting_infos": [],
			"validator_slash_events": []
		},
		"evidence": {
			"evidence": []
		},
		"feegrant": {
			"allowances": []
		},
		"feeibc": {
			"identified_fees": [],
			"fee_enabled_channels": [],
			"registered_payees": [],
			"registered_counterparty_payees": [],
			"forward_relayers": []
		},
		"feeshare": {
			"params": {
				"enable_fee_share": true,
				"developer_shares": "0.500000000000000000",
				"allowed_denoms": []
			},
			"fee_share": []
		},
		"genutil": {
			"gen_txs": []
		},
		"gov": {
			"starting_proposal_id": "1",
			"deposits": [],
			"votes": [],
			"proposals": [],
			"deposit_params": null,
			"voting_params": null,
			"tally_params": null,
			"params": {
				"min_deposit": [
					{
						"denom": "uluna",
						"amount": "10000000"
					}
				],
				"max_deposit_period": "172800s",
				"voting_period": "172800s",
				"quorum": "0.334000000000000000",
				"threshold": "0.500000000000000000",
				"veto_threshold": "0.334000000000000000",
				"min_initial_deposit_ratio": "0.000000000000000000",
				"burn_vote_quorum": false,
				"burn_proposal_deposit_prevote": false,
				"burn_vote_veto": true
			}
		},
		"ibc": {
			"client_genesis": {
				"clients": [],
				"clients_consensus": [],
				"clients_metadata": [],
				"params": {
					"allowed_clients": [
						"06-solomachine",
						"07-tendermint",
						"09-localhost"
					]
				},
				"create_localhost": false,
				"next_client_sequence": "0"
			},
			"connection_genesis": {
				"connections": [],
				"client_connection_paths": [],
				"next_connection_sequence": "0",
				"params": {
					"max_expected_time_per_block": "30000000000"
				}
			},
			"channel_genesis": {
				"channels": [],
				"acknowledgements": [],
				"commitments": [],
				"receipts": [],
				"send_sequences": [],
				"recv_sequences": [],
				"ack_sequences": [],
				"next_channel_sequence": "0"
			}
		},
		"ibchooks": {},
		"interchainaccounts": {
			"controller_genesis_state": {
				"active_channels": [],
				"interchain_accounts": [],
				"ports": [],
				"params": {
					"controller_enabled": true
				}
			},
			"host_genesis_state": {
				"active_channels": [],
				"interchain_accounts": [],
				"port": "icahost",
				"params": {
					"host_enabled": true,
					"allow_messages": [
						"/cosmos.authz.v1beta1.MsgExec",
						"/cosmos.authz.v1beta1.MsgGrant",
						"/cosmos.authz.v1beta1.MsgRevoke",
						"/cosmos.bank.v1beta1.MsgSend",
						"/cosmos.bank.v1beta1.MsgMultiSend",
						"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress",
						"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
						"/cosmos.distribution.v1beta1.MsgFundCommunityPool",
						"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
						"/cosmos.feegrant.v1beta1.MsgGrantAllowance",
						"/cosmos.feegrant.v1beta1.MsgRevokeAllowance",
						"/cosmos.gov.v1beta1.MsgVoteWeighted",
						"/cosmos.gov.v1beta1.MsgSubmitProposal",
						"/cosmos.gov.v1beta1.MsgDeposit",
						"/cosmos.gov.v1beta1.MsgVote",
						"/cosmos.staking.v1beta1.MsgEditValidator",
						"/cosmos.staking.v1beta1.MsgDelegate",
						"/cosmos.staking.v1beta1.MsgUndelegate",
						"/cosmos.staking.v1beta1.MsgBeginRedelegate",
						"/cosmos.staking.v1beta1.MsgCreateValidator",
						"/cosmos.vesting.v1beta1.MsgCreateVestingAccount",
						"/ibc.applications.transfer.v1.MsgTransfer",
						"/cosmwasm.wasm.v1.MsgStoreCode",
						"/cosmwasm.wasm.v1.MsgInstantiateContract",
						"/cosmwasm.wasm.v1.MsgExecuteContract",
						"/cosmwasm.wasm.v1.MsgMigrateContract"
					]
				}
			}
		},
		"interchainquery": {
			"host_port": "icqhost",
			"params": {
				"host_enabled": true,
				"allow_queries": [
					"/alliance.alliance.Query/AllAllianceValidators",
					"/alliance.alliance.Query/AllAlliancesDelegations",
					"/alliance.alliance.Query/Alliance",
					"/alliance.alliance.Query/AllianceDelegation",
					"/alliance.alliance.Query/AllianceDelegationRewards",
					"/alliance.alliance.Query/AllianceRedelegations",
					"/alliance.alliance.Query/AllianceUnbondings",
					"/alliance.alliance.Query/AllianceUnbondingsByDenomAndDelegator",
					"/alliance.alliance.Query/AllianceValidator",
					"/alliance.alliance.Query/Alliances",
					"/alliance.alliance.Query/AlliancesDelegation",
					"/alliance.alliance.Query/AlliancesDelegationByValidator",
					"/alliance.alliance.Query/IBCAlliance",
					"/alliance.alliance.Query/IBCAllianceDelegation",
					"/alliance.alliance.Query/IBCAllianceDelegationRewards",
					"/alliance.alliance.Query/Params",
					"/cosmos.auth.v1beta1.Query/Account",
					"/cosmos.auth.v1beta1.Query/AccountAddressByID",
					"/cosmos.auth.v1beta1.Query/AccountInfo",
					"/cosmos.auth.v1beta1.Query/Accounts",
					"/cosmos.auth.v1beta1.Query/AddressBytesToString",
					"/cosmos.auth.v1beta1.Query/AddressStringToBytes",
					"/cosmos.auth.v1beta1.Query/Bech32Prefix",
					"/cosmos.auth.v1beta1.Query/ModuleAccountByName",
					"/cosmos.auth.v1beta1.Query/ModuleAccounts",
					"/cosmos.auth.v1beta1.Query/Params",
					"/cosmos.authz.v1beta1.Query/GranteeGrants",
					"/cosmos.authz.v1beta1.Query/GranterGrants",
					"/cosmos.authz.v1beta1.Query/Grants",
					"/cosmos.bank.v1beta1.Query/AllBalances",
					"/cosmos.bank.v1beta1.Query/Balance",
					"/cosmos.bank.v1beta1.Query/DenomMetadata",
					"/cosmos.bank.v1beta1.Query/DenomOwners",
					"/cosmos.bank.v1beta1.Query/DenomsMetadata",
					"/cosmos.bank.v1beta1.Query/Params",
					"/cosmos.bank.v1beta1.Query/SendEnabled",
					"/cosmos.bank.v1beta1.Query/SpendableBalanceByDenom",
					"/cosmos.bank.v1beta1.Query/SpendableBalances",
					"/cosmos.bank.v1beta1.Query/SupplyOf",
					"/cosmos.bank.v1beta1.Query/TotalSupply",
					"/cosmos.consensus.v1.Query/Params",
					"/cosmos.distribution.v1beta1.Query/CommunityPool",
					"/cosmos.distribution.v1beta1.Query/DelegationRewards",
					"/cosmos.distribution.v1beta1.Query/DelegationTotalRewards",
					"/cosmos.distribution.v1beta1.Query/DelegatorValidators",
					"/cosmos.distribution.v1beta1.Query/DelegatorWithdrawAddress",
					"/cosmos.distribution.v1beta1.Query/Params",
					"/cosmos.distribution.v1beta1.Query/ValidatorCommission",
					"/cosmos.distribution.v1beta1.Query/ValidatorDistributionInfo",
					"/cosmos.distribution.v1beta1.Query/ValidatorOutstandingRewards",
					"/cosmos.distribution.v1beta1.Query/ValidatorSlashes",
					"/cosmos.evidence.v1beta1.Query/AllEvidence",
					"/cosmos.evidence.v1beta1.Query/Evidence",
					"/cosmos.feegrant.v1beta1.Query/Allowance",
					"/cosmos.feegrant.v1beta1.Query/Allowances",
					"/cosmos.feegrant.v1beta1.Query/AllowancesByGranter",
					"/cosmos.gov.v1.Query/Deposit",
					"/cosmos.gov.v1.Query/Deposits",
					"/cosmos.gov.v1.Query/Params",
					"/cosmos.gov.v1.Query/Proposal",
					"/cosmos.gov.v1.Query/Proposals",
					"/cosmos.gov.v1.Query/TallyResult",
					"/cosmos.gov.v1.Query/Vote",
					"/cosmos.gov.v1.Query/Votes",
					"/cosmos.mint.v1beta1.Query/AnnualProvisions",
					"/cosmos.mint.v1beta1.Query/Inflation",
					"/cosmos.mint.v1beta1.Query/Params",
					"/cosmos.params.v1beta1.Query/Params",
					"/cosmos.params.v1beta1.Query/Subspaces",
					"/cosmos.slashing.v1beta1.Query/Params",
					"/cosmos.slashing.v1beta1.Query/SigningInfo",
					"/cosmos.slashing.v1beta1.Query/SigningInfos",
					"/cosmos.staking.v1beta1.Query/Delegation",
					"/cosmos.staking.v1beta1.Query/DelegatorDelegations",
					"/cosmos.staking.v1beta1.Query/DelegatorUnbondingDelegations",
					"/cosmos.staking.v1beta1.Query/DelegatorValidator",
					"/cosmos.staking.v1beta1.Query/DelegatorValidators",
					"/cosmos.staking.v1beta1.Query/HistoricalInfo",
					"/cosmos.staking.v1beta1.Query/Params",
					"/cosmos.staking.v1beta1.Query/Pool",
					"/cosmos.staking.v1beta1.Query/Redelegations",
					"/cosmos.staking.v1beta1.Query/UnbondingDelegation",
					"/cosmos.staking.v1beta1.Query/Validator",
					"/cosmos.staking.v1beta1.Query/ValidatorDelegations",
					"/cosmos.staking.v1beta1.Query/ValidatorUnbondingDelegations",
					"/cosmos.staking.v1beta1.Query/Validators",
					"/cosmos.upgrade.v1beta1.Query/AppliedPlan",
					"/cosmos.upgrade.v1beta1.Query/Authority",
					"/cosmos.upgrade.v1beta1.Query/CurrentPlan",
					"/cosmos.upgrade.v1beta1.Query/ModuleVersions",
					"/cosmos.upgrade.v1beta1.Query/UpgradedConsensusState",
					"/cosmwasm.wasm.v1.Query/AllContractState",
					"/cosmwasm.wasm.v1.Query/Code",
					"/cosmwasm.wasm.v1.Query/Codes",
					"/cosmwasm.wasm.v1.Query/ContractHistory",
					"/cosmwasm.wasm.v1.Query/ContractInfo",
					"/cosmwasm.wasm.v1.Query/ContractsByCode",
					"/cosmwasm.wasm.v1.Query/ContractsByCreator",
					"/cosmwasm.wasm.v1.Query/Params",
					"/cosmwasm.wasm.v1.Query/PinnedCodes",
					"/cosmwasm.wasm.v1.Query/RawContractState",
					"/cosmwasm.wasm.v1.Query/SmartContractState",
					"/ibc.applications.fee.v1.Query/CounterpartyPayee",
					"/ibc.applications.fee.v1.Query/FeeEnabledChannel",
					"/ibc.applications.fee.v1.Query/FeeEnabledChannels",
					"/ibc.applications.fee.v1.Query/IncentivizedPacket",
					"/ibc.applications.fee.v1.Query/IncentivizedPackets",
					"/ibc.applications.fee.v1.Query/IncentivizedPacketsForChannel",
					"/ibc.applications.fee.v1.Query/Payee",
					"/ibc.applications.fee.v1.Query/TotalAckFees",
					"/ibc.applications.fee.v1.Query/TotalRecvFees",
					"/ibc.applications.fee.v1.Query/TotalTimeoutFees",
					"/ibc.applications.interchain_accounts.controller.v1.Query/InterchainAccount",
					"/ibc.applications.interchain_accounts.controller.v1.Query/Params",
					"/ibc.applications.interchain_accounts.host.v1.Query/Params",
					"/ibc.applications.transfer.v1.Query/DenomHash",
					"/ibc.applications.transfer.v1.Query/DenomTrace",
					"/ibc.applications.transfer.v1.Query/DenomTraces",
					"/ibc.applications.transfer.v1.Query/EscrowAddress",
					"/ibc.applications.transfer.v1.Query/Params",
					"/ibc.applications.transfer.v1.Query/TotalEscrowForDenom",
					"/ibc.core.channel.v1.Query/Channel",
					"/ibc.core.channel.v1.Query/ChannelClientState",
					"/ibc.core.channel.v1.Query/ChannelConsensusState",
					"/ibc.core.channel.v1.Query/Channels",
					"/ibc.core.channel.v1.Query/ConnectionChannels",
					"/ibc.core.channel.v1.Query/NextSequenceReceive",
					"/ibc.core.channel.v1.Query/PacketAcknowledgement",
					"/ibc.core.channel.v1.Query/PacketAcknowledgements",
					"/ibc.core.channel.v1.Query/PacketCommitment",
					"/ibc.core.channel.v1.Query/PacketCommitments",
					"/ibc.core.channel.v1.Query/PacketReceipt",
					"/ibc.core.channel.v1.Query/UnreceivedAcks",
					"/ibc.core.channel.v1.Query/UnreceivedPackets",
					"/ibc.core.client.v1.Query/ClientParams",
					"/ibc.core.client.v1.Query/ClientState",
					"/ibc.core.client.v1.Query/ClientStates",
					"/ibc.core.client.v1.Query/ClientStatus",
					"/ibc.core.client.v1.Query/ConsensusState",
					"/ibc.core.client.v1.Query/ConsensusStateHeights",
					"/ibc.core.client.v1.Query/ConsensusStates",
					"/ibc.core.client.v1.Query/UpgradedClientState",
					"/ibc.core.client.v1.Query/UpgradedConsensusState",
					"/ibc.core.connection.v1.Query/ClientConnections",
					"/ibc.core.connection.v1.Query/Connection",
					"/ibc.core.connection.v1.Query/ConnectionClientState",
					"/ibc.core.connection.v1.Query/ConnectionConsensusState",
					"/ibc.core.connection.v1.Query/ConnectionParams",
					"/ibc.core.connection.v1.Query/Connections",
					"/icq.v1.Query/Params",
					"/juno.feeshare.v1.Query/DeployerFeeShares",
					"/juno.feeshare.v1.Query/FeeShare",
					"/juno.feeshare.v1.Query/FeeShares",
					"/juno.feeshare.v1.Query/Params",
					"/juno.feeshare.v1.Query/WithdrawerFeeShares",
					"/osmosis.tokenfactory.v1beta1.Query/BeforeSendHookAddress",
					"/osmosis.tokenfactory.v1beta1.Query/DenomAuthorityMetadata",
					"/osmosis.tokenfactory.v1beta1.Query/DenomsFromCreator",
					"/osmosis.tokenfactory.v1beta1.Query/Params",
					"/pob.builder.v1.Query/Params",
					"/router.v1.Query/Params"
				]
			}
		},
		"mint": {
			"minter": {
				"inflation": "0.130000000000000000",
				"annual_provisions": "0.000000000000000000"
			},
			"params": {
				"mint_denom": "uluna",
				"inflation_rate_change": "0.130000000000000000",
				"inflation_max": "0.200000000000000000",
				"inflation_min": "0.070000000000000000",
				"goal_bonded": "0.670000000000000000",
				"blocks_per_year": "6311520"
			}
		},
		"packetfowardmiddleware": {
			"params": {
				"fee_percentage": "0.000000000000000000"
			},
			"in_flight_packets": {}
		},
		"params": null,
		"slashing": {
			"params": {
				"signed_blocks_window": "100",
				"min_signed_per_window": "0.500000000000000000",
				"downtime_jail_duration": "600s",
				"slash_fraction_double_sign": "0.050000000000000000",
				"slash_fraction_downtime": "0.010000000000000000"
			},
			"signing_infos": [],
			"missed_blocks": []
		},
		"staking": {
			"params": {
				"unbonding_time": "1814400s",
				"max_validators": 100,
				"max_entries": 7,
				"historical_entries": 10000,
				"bond_denom": "uluna",
				"min_commission_rate": "0.000000000000000000"
			},
			"last_total_power": "0",
			"last_validator_powers": [],
			"validators": [],
			"delegations": [],
			"unbonding_delegations": [],
			"redelegations": [],
			"exported": false
		},
		"tokenfactory": {
			"params": {
				"denom_creation_fee": [
					{
						"denom": "uluna",
						"amount": "10000000"
					}
				],
				"denom_creation_gas_consume": "1000000"
			},
			"factory_denoms": []
		},
		"transfer": {
			"port_id": "transfer",
			"denom_traces": [],
			"params": {
				"send_enabled": true,
				"receive_enabled": true
			},
			"total_escrowed": []
		},
		"upgrade": {},
		"vesting": {},
		"wasm": {
			"params": {
				"code_upload_access": {
					"permission": "Everybody",
					"addresses": []
				},
				"instantiate_default_permission": "Everybody"
			},
			"codes": [],
			"contracts": [],
			"sequences": []
		}
	}`
	require.JSONEq(t, string(jsonGenState), expectedState)
}
