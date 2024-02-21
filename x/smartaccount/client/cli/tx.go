package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/terra-money/core/v2/x/smartaccount/types"
)

// NewTxCmd returns a root CLI command handler for certain modules/SmartAccounts
// transaction commands.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "SmartAccounts subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewCreateSmartAccount(),
		NewDisableSmartAccount(),
	)
	return txCmd
}

// NewRegisterFeeShare returns a CLI command handler for registering a
// contract for fee distribution
func NewCreateSmartAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-smart-account",
		Short: "Create a smart account for the caller.",
		Long:  "Create a smart account for the caller.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			msg := &types.MsgCreateSmartAccount{
				Account: account.String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewRegisterFeeShare returns a CLI command handler for registering a
// contract for fee distribution
func NewDisableSmartAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable-smart-account",
		Short: "Disable smart account of the caller.",
		Long:  "Disable a smart account of the caller.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			msg := &types.MsgDisableSmartAccount{
				Account: account.String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
