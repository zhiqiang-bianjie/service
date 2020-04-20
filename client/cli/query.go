package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irismod/service/types"
)

// GetQueryCmd returns the cli query commands for the module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	serviceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the service module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	serviceQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryServiceDefinition(queryRoute, cdc),
		GetCmdQueryServiceBinding(queryRoute, cdc),
		GetCmdQueryServiceBindings(queryRoute, cdc),
		GetCmdQueryWithdrawAddr(queryRoute, cdc),
	)...)

	return serviceQueryCmd
}

// GetCmdQueryServiceDefinition implements the query service definition command.
func GetCmdQueryServiceDefinition(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "definition [service-name]",
		Short: "Query a service definition",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a service definition.

Example:
$ %s query service definition <service-name>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			if err := types.ValidateServiceName(args[0]); err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.QueryDefinitionParams{ServiceName: args[0]})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDefinition)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var definition types.ServiceDefinition
			if err := cdc.UnmarshalJSON(res, &definition); err != nil {
				return err
			}

			return cliCtx.PrintOutput(definition)
		},
	}
}

// GetCmdQueryServiceBinding implements the query service binding command
func GetCmdQueryServiceBinding(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "binding [service-name] [provider]",
		Short: "Query a service binding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a service binding.

Example:
$ %s query service binding <service-name> <provider>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			if err := types.ValidateServiceName(args[0]); err != nil {
				return err
			}

			provider, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.QueryBindingParams{ServiceName: args[0], Provider: provider})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryBinding)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var binding types.ServiceBinding
			if err := cdc.UnmarshalJSON(res, &binding); err != nil {
				return err
			}

			return cliCtx.PrintOutput(binding)
		},
	}

	return cmd
}

// GetCmdQueryServiceBindings implements the query service bindings command
func GetCmdQueryServiceBindings(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bindings [service-name]",
		Short: "Query all bindings of a service definition",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all bindings of a service definition.

Example:
$ %s query service bindings <service-name>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			if err := types.ValidateServiceName(args[0]); err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.QueryBindingsParams{ServiceName: args[0]})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryBindings)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var bindings []types.ServiceBinding
			if err := cdc.UnmarshalJSON(res, &bindings); err != nil {
				return err
			}

			return cliCtx.PrintOutput(bindings)
		},
	}

	return cmd
}

// GetCmdQueryWithdrawAddr implements the query withdraw address command
func GetCmdQueryWithdrawAddr(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-addr [provider]",
		Short: "Query the withdrawal address of a provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the withdrawal address of a provider.

Example:
$ %s query service withdraw-addr <provider>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			provider, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.QueryWithdrawAddressParams{Provider: provider})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryWithdrawAddress)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var withdrawAddr sdk.AccAddress
			if err := cdc.UnmarshalJSON(res, &withdrawAddr); err != nil {
				return err
			}

			return cliCtx.PrintOutput(withdrawAddr)
		},
	}

	return cmd
}
