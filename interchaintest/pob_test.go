package interchaintest

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"testing"

	"github.com/skip-mev/pob/tests/integration"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/stretchr/testify/suite"
)

var (
	numVals          = 4
	numFull          = 0
	Denom            = "uluna"
	VotingPeriod     = "15s"
	MaxDepositPeriod = "10s"
	config           = ibc.ChainConfig{
		Type:    "cosmos",
		Name:    "terra",
		ChainID: "phoenix-1",
		Images: []ibc.DockerImage{
			{
				Repository: "terramoneycore",
				Version:    "latest",
				UidGid:     "1025:1025",
			},
		},
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

func GetInterchainSpecForPOB() *interchaintest.ChainSpec {
	// update the genesis kv for juno
	updatedChainConfig := config
	updatedChainConfig.ModifyGenesis = cosmos.ModifyGenesis(append(defaultGenesisKV, []cosmos.GenesisKV{
		{
			Key:   "app_state.builder.params.max_bundle_size",
			Value: 3,
		},
		{
			Key:   "app_state.builder.params.reserve_fee.denom",
			Value: "uluna",
		},
		{
			Key:   "app_state.builder.params.reserve_fee.amount",
			Value: "1",
		},
		{
			Key:   "app_state.builder.params.min_bid_increment.denom",
			Value: "uluna",
		},
		{
			Key:   "app_state.builder.params.min_bid_increment.amount",
			Value: "1",
		},
	}...))

	return &interchaintest.ChainSpec{
		Name:          "terra",
		ChainName:     "terra",
		Version:       "latest",
		ChainConfig:   updatedChainConfig,
		NumValidators: &numVals,
		NumFullNodes:  &numFull,
	}

}

func TestPOB(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount("terra", "terra")
	s := integration.NewPOBIntegrationTestSuiteFromSpec(GetInterchainSpecForPOB())
	s.WithDenom("uluna")

	suite.Run(t, s)
}
