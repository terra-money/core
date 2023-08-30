package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terra-money/core/v2/x/tokenfactory/types"
)

func TestAuthorityMetadata(t *testing.T) {
	data := types.DenomAuthorityMetadata{
		Admin: "satoshi",
	}

	require.Error(t, data.Validate())
}
