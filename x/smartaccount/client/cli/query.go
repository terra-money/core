package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/terra-money/core/v2/x/smartaccount/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	smartAccountsQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	smartAccountsQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQuerySetting(),
	)

	return smartAccountsQueryCmd
}

// GetCmdQueryParams implements a command to return the current SmartAccounts
// parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current smartaccount module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryParamsRequest{}

			res, err := queryClient.Params(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQuerySetting() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setting",
		Short: "Query the current smartaccount setting for address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QuerySettingRequest{
				Address: args[0],
			}

			res, err := queryClient.Setting(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Setting)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
