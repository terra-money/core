package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	appparams "github.com/terra-money/core/v2/app/params"
	"github.com/terra-money/core/v2/x/tokenfactory/types"
)

func TestGenesisState_Validate(t *testing.T) {
	types.DefaultGenesis()
	appparams.RegisterAddressesConfig()

	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "different admin from creator",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "empty admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "no admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin",
					},
				},
			},
			valid: true,
		},
		{
			desc: "invalid admin",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "moose",
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "multiple denoms",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
					{
						Denom: "factory/terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/litecoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
			},
			valid: true,
		},
		{
			desc: "duplicate denoms",
			genState: &types.GenesisState{
				FactoryDenoms: []types.GenesisDenom{
					{
						Denom: "factory/terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
					{
						Denom: "factory/terra19hukvr8hppdwqnx7tkaslarz5s449qahu5kp2g/bitcoin",
						AuthorityMetadata: types.DenomAuthorityMetadata{
							Admin: "",
						},
					},
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
