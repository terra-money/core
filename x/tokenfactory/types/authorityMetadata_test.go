package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/terra-money/core/v2/x/tokenfactory/types"
)

func TestAuthorityMetadataError(t *testing.T) {
	data := types.DenomAuthorityMetadata{
		Admin: "satoshi",
	}

	require.Error(t, data.Validate())
}

func TestAuthorityMetadata(t *testing.T) {
	data := types.DenomAuthorityMetadata{
		Admin: "terra1zdpgj8am5nqqvht927k3etljyl6a52kwqup0je",
	}

	require.Error(t, data.Validate())
}
