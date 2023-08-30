package interchaintest

import (
	"fmt"

	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"

	"github.com/cosmos/cosmos-sdk/types/module/testutil"
)

var (
	Denom            = "uluna"
	VotingPeriod     = "15s"
	MaxDepositPeriod = "10s"
	Image            = ibc.DockerImage{
		Repository: "terramoneycore",
		Version:    "latest",
		UidGid:     "1025:1025",
	}
	IBCRelayerImage   = "ghcr.io/cosmos/relayer"
	IBCRelayerVersion = "main"
	config            = ibc.ChainConfig{
		Type:                   "cosmos",
		Name:                   "terra",
		ChainID:                "phoenix-1",
		Images:                 []ibc.DockerImage{Image},
		Bin:                    "terrad",
		Bech32Prefix:           "terra",
		Denom:                  Denom,
		CoinType:               "330",
		GasPrices:              fmt.Sprintf("0%s", Denom),
		GasAdjustment:          2.0,
		TrustingPeriod:         "112h",
		NoHostMount:            false,
		ConfigFileOverrides:    nil,
		EncodingConfig:         encoding(),
		UsingNewGenesisCommand: true,
		ModifyGenesis:          cosmos.ModifyGenesis(defaultGenesisKV),
	}
	// SDK v47 Genesis
	defaultGenesisKV = []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: VotingPeriod,
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: MaxDepositPeriod,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: Denom,
		},
	}
)

func encoding() *testutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()
	return &cfg
}
