package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/app/ante/types"
)

// GetQueryCmd returns the transaction commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the ante module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		QueryParamsCmd(),
		QueryMinimumCommissionCmd(),
	)

	return cmd
}

// QueryParamsCmd returns the command handler for ante parameter querying.
func QueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current ante parameters",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query the current ante parameters:

$ <appd> query ante params
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// QueryMinimumCommissionCmd returns the command handler for minimum commission querying.
func QueryMinimumCommissionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "minimum-commission",
		Short: "Query the current minimum commission",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query the current minimum commission:

$ <appd> query ante minimum-commission
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.MinimumCommission(cmd.Context(), &types.QueryMinimumCommissionRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&sdk.DecProto{Dec: res.MinimumCommission})
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
