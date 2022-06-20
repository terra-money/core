package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateMinimumCommissionEnforced(t *testing.T) {
	require.Error(t, validateMiniumCommissionEnforced(int64(3)))
	require.NoError(t, validateMiniumCommissionEnforced(true))
	require.NoError(t, validateMiniumCommissionEnforced(false))
}
