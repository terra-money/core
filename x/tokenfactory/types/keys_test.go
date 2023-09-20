package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terra-money/core/v2/x/tokenfactory/types"
)

func TestKeys(t *testing.T) {
	denom_key := types.GetDenomPrefixStore("denom")
	creator_key := types.GetCreatorPrefix("creator")
	creator_prefix_key := types.GetCreatorsPrefix()

	require.Equal(t, []byte("denoms|denom|"), denom_key)
	require.Equal(t, []byte("creator|creator|"), creator_key)
	require.Equal(t, []byte("creator|"), creator_prefix_key)
}
