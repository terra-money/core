package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/terra-money/core/v2/x/smartaccount/types"
)

// NewTxCmd returns a root CLI command handler for certain modules/SmartAccounts
// transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "SmartAccounts subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand()
	return txCmd
}
