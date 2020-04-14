package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
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
