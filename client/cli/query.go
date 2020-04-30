package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irismod/service/client/utils"
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
		GetCmdQueryServiceRequest(queryRoute, cdc),
		GetCmdQueryServiceRequests(queryRoute, cdc),
		GetCmdQueryServiceResponse(queryRoute, cdc),
		GetCmdQueryRequestContext(queryRoute, cdc),
		GetCmdQueryServiceResponses(queryRoute, cdc),
		GetCmdQueryEarnedFees(queryRoute, cdc),
		GetCmdQuerySchema(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
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
		Use:   "binding [service-name] [provider-address]",
		Short: "Query a service binding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a service binding.

Example:
$ %s query service binding <service-name> <provider-address>
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
		Short: "Query all bindings of a service definition with an optional owner",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all bindings of a service definition with an optional owner.

Example:
$ %s query service bindings <service-name> --owner=<address>
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

			var err error
			var owner sdk.AccAddress

			ownerStr := viper.GetString(FlagOwner)
			if len(ownerStr) > 0 {
				owner, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			params := types.QueryBindingsParams{
				ServiceName: args[0],
				Owner:       owner,
			}

			bz, err := cdc.MarshalJSON(params)
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

	cmd.Flags().AddFlagSet(FsQueryServiceBindings)

	return cmd
}

// GetCmdQueryWithdrawAddr implements the query withdraw address command
func GetCmdQueryWithdrawAddr(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-addr [address]",
		Short: "Query the withdrawal address of an owner",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the withdrawal address of an owner.

Example:
$ %s query service withdraw-addr <address>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.QueryWithdrawAddressParams{Owner: owner})
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

// GetCmdQueryServiceRequest implements the query service request command
func GetCmdQueryServiceRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "request [request-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a service request.

Example:
$ %s query service request <request-id>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			requestID, err := types.ConvertRequestID(args[0])
			if err != nil {
				return err
			}
			params := types.QueryRequestParams{
				RequestID: requestID,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequest)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var request types.Request
			_ = cdc.UnmarshalJSON(res, &request)
			if request.Empty() {
				request, err = utils.QueryRequestByTxQuery(cliCtx, queryRoute, params)
				if err != nil {
					return err
				}
			}

			if request.Empty() {
				return fmt.Errorf("unknown request: %s", params.RequestID)
			}

			return cliCtx.PrintOutput(request)
		},
	}

	return cmd
}

// GetCmdQueryServiceRequests implements the query service requests command
func GetCmdQueryServiceRequests(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "requests [service-name] [provider] | [request-context-id] [batch-counter]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query active requests by the service binding or request context ID.

Example:
$ %s query service requests <service-name> <provider> | <request-context-id> <batch-counter>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			queryByBinding := true

			provider, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				queryByBinding = false
			}

			var requests types.Requests

			if queryByBinding {
				requests, _, err = utils.QueryRequestsByBinding(cliCtx, queryRoute, args[0], provider)
			} else {
				requests, _, err = utils.QueryRequestsByReqCtx(cliCtx, queryRoute, args[0], args[1])
			}

			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(requests)
		},
	}

	return cmd
}

// GetCmdQueryServiceResponse implements the query service response command
func GetCmdQueryServiceResponse(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "response [request-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a service response.

Example:
$ %s query service response <request-id>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			requestID, err := types.ConvertRequestID(args[0])
			if err != nil {
				return err
			}
			params := types.QueryResponseParams{
				RequestID: requestID,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryResponse)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var response types.Response
			_ = cdc.UnmarshalJSON(res, &response)
			if response.Empty() {
				response, err = utils.QueryResponseByTxQuery(cliCtx, queryRoute, params)
				if err != nil {
					return err
				}
			}

			if response.Empty() {
				return fmt.Errorf("unknown response: %s", params.RequestID)
			}

			return cliCtx.PrintOutput(response)
		},
	}

	return cmd
}

// GetCmdQueryServiceResponses implements the query service responses command
func GetCmdQueryServiceResponses(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "responses [request-context-id] [batch-counter]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query active responses by the request context ID and batch counter.

Example:
$ %s query service responses <request-context-id> <batch-counter>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			batchCounter, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			params := types.QueryResponsesParams{
				RequestContextID: requestContextID,
				BatchCounter:     batchCounter,
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryResponses)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var responses types.Responses
			if err := cdc.UnmarshalJSON(res, &responses); err != nil {
				return err
			}

			return cliCtx.PrintOutput(responses)
		},
	}

	return cmd
}

// GetCmdQueryRequestContext implements the query request context command
func GetCmdQueryRequestContext(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "request-context [request-context-id]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a request context.

Example:
$ %s query service request-context <request-context-id>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			params := types.QueryRequestContextParams{
				RequestContextID: requestContextID,
			}

			requestContext, err := utils.QueryRequestContext(cliCtx, queryRoute, params)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(requestContext)
		},
	}

	return cmd
}

// GetCmdQueryEarnedFees implements the query earned fees command
func GetCmdQueryEarnedFees(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fees [provider-address]",
		Short: "Query the earned fees of a provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the earned fees of a provider.

Example:
$ %s query service fees <provider-address>
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

			bz, err := cdc.MarshalJSON(types.QueryEarnedFeesParams{Provider: provider})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryEarnedFees)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var feesOut types.EarnedFeesOutput
			if err := cdc.UnmarshalJSON(res, &feesOut); err != nil {
				return err
			}

			return cliCtx.PrintOutput(feesOut)
		},
	}

	return cmd
}

// GetCmdQuerySchema implements the query schema command
func GetCmdQuerySchema(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "schema [schema-name]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the system schema by the schema name, only pricing and result allowed.

Example:
$ %s query service schema <schema-name>
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			params := types.QuerySchemaParams{
				SchemaName: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySchema)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var schema utils.SchemaType
			if err := cdc.UnmarshalJSON(res, &schema); err != nil {
				return err
			}

			return cliCtx.PrintOutput(schema)
		},
	}

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current service parameter values",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as service parameters.
Example:
$ %s query service params
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)

			return cliCtx.PrintOutput(params)
		},
	}
}
