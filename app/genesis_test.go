package app_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terra-money/core/v2/app"
)

func TestNewGenesis(t *testing.T) {
	encCfg := app.MakeEncodingConfig()
	genesisState := app.NewDefaultGenesisState(encCfg.Marshaler)

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
		"capability": {
			"index": "1",
			"owners": []
		},
		"consensus": null,
		"crisis": {
			"constant_fee": {
				"denom": "stake",
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
						"denom": "stake",
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
						"*"
					]
				}
			}
		},
		"mint": {
			"minter": {
				"inflation": "0.130000000000000000",
				"annual_provisions": "0.000000000000000000"
			},
			"params": {
				"mint_denom": "stake",
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
				"bond_denom": "stake",
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

func TestNewGenesisWithBondDenom(t *testing.T) {
	encCfg := app.MakeEncodingConfig()
	genesisState := app.NewDefaultGenesisState(encCfg.Marshaler)
	genesisState.ConfigureBondDenom(encCfg.Marshaler, "uluna")

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
						"*"
					]
				}
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

func TestNewGenesisConfigureICA(t *testing.T) {
	encCfg := app.MakeEncodingConfig()
	genesisState := app.NewDefaultGenesisState(encCfg.Marshaler)
	genesisState.ConfigureICA(encCfg.Marshaler)

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
		"capability": {
			"index": "1",
			"owners": []
		},
		"consensus": null,
		"crisis": {
			"constant_fee": {
				"denom": "stake",
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
						"denom": "stake",
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
		"mint": {
			"minter": {
				"inflation": "0.130000000000000000",
				"annual_provisions": "0.000000000000000000"
			},
			"params": {
				"mint_denom": "stake",
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
				"bond_denom": "stake",
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
	fmt.Print(string(jsonGenState))

	require.JSONEq(t, string(jsonGenState), expectedState)
}
