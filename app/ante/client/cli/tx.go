package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/terra-money/core/v2/app/ante/types"
)

// NewCmdSubmitMinimumCommissionUpdateProposal implements a command handler for submitting a minimum commission proposal transaction.
func NewCmdSubmitMinimumCommissionUpdateProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "minimum-commission-update [minimum-commission] [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a minimum commission update proposal",
		Long:  "Submit a minimum commission update along with an initial deposit.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			minimumCommissionStr := args[0]
			minimumCommission, err := sdk.NewDecFromStr(minimumCommissionStr)
			if err != nil {
				return err
			}

			content, err := parseArgsToContent(cmd, minimumCommission)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := gov.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	return cmd
}

func parseArgsToContent(cmd *cobra.Command, minimumCommission sdk.Dec) (gov.Content, error) {
	title, err := cmd.Flags().GetString(cli.FlagTitle)
	if err != nil {
		return nil, err
	}

	description, err := cmd.Flags().GetString(cli.FlagDescription)
	if err != nil {
		return nil, err
	}

	content := types.NewMinimumCommissionUpdateProposal(title, description, minimumCommission)
	return content, nil
}
