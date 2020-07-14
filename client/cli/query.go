package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irismod/service/client/utils"
	"github.com/irismod/service/types"
)

// GetQueryCmd returns the cli query commands for the module.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	serviceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the service module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	serviceQueryCmd.AddCommand(
		GetCmdQueryServiceDefinition(cdc),
		GetCmdQueryServiceBinding(cdc),
		GetCmdQueryServiceBindings(cdc),
		GetCmdQueryWithdrawAddr(cdc),
		GetCmdQueryServiceRequest(cdc),
		GetCmdQueryServiceRequests(cdc),
		GetCmdQueryServiceResponse(cdc),
		GetCmdQueryRequestContext(cdc),
		GetCmdQueryServiceResponses(cdc),
		GetCmdQueryEarnedFees(cdc),
		GetCmdQuerySchema(cdc),
		GetCmdQueryParams(cdc),
	)

	return serviceQueryCmd
}

// GetCmdQueryServiceDefinition implements the query service definition command.
func GetCmdQueryServiceDefinition(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "definition [service-name]",
		Short: "Query a service definition",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a service definition.

Example:
$ %s query service definition <service-name>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			if err := types.ValidateServiceName(args[0]); err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.QueryDefinitionParams{ServiceName: args[0]})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDefinition)
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
func GetCmdQueryServiceBinding(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "binding [service-name] [provider-address]",
		Short: "Query a service binding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a service binding.

Example:
$ %s query service binding <service-name> <provider-address>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

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

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBinding)
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
func GetCmdQueryServiceBindings(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bindings [service-name]",
		Short: "Query all bindings of a service definition with an optional owner",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all bindings of a service definition with an optional owner.

Example:
$ %s query service bindings <service-name> --owner=<address>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

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

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBindings)
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
func GetCmdQueryWithdrawAddr(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-addr [address]",
		Short: "Query the withdrawal address of an owner",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the withdrawal address of an owner.

Example:
$ %s query service withdraw-addr <address>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.QueryWithdrawAddressParams{Owner: owner})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryWithdrawAddress)
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
func GetCmdQueryServiceRequest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request [request-id]",
		Short: "Query a request by the request ID",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a service request.

Example:
$ %s query service request <request-id>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

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

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRequest)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var request types.Request
			_ = cdc.UnmarshalJSON(res, &request)
			if request.Empty() {
				request, err = utils.QueryRequestByTxQuery(cliCtx, types.QuerierRoute, params)
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
func GetCmdQueryServiceRequests(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "requests [service-name] [provider] | [request-context-id] [batch-counter]",
		Short: "Query active requests by the service binding or request context ID",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query active requests by the service binding or request context ID.

Example:
$ %s query service requests <service-name> <provider> | <request-context-id> <batch-counter>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			queryByBinding := true

			provider, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				queryByBinding = false
			}

			var requests []types.Request

			if queryByBinding {
				requests, _, err = utils.QueryRequestsByBinding(cliCtx, types.QuerierRoute, args[0], provider)
			} else {
				requests, _, err = utils.QueryRequestsByReqCtx(cliCtx, types.QuerierRoute, args[0], args[1])
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
func GetCmdQueryServiceResponse(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "response [request-id]",
		Short: "Query a response by the request ID",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of a service response.

Example:
$ %s query service response <request-id>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

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

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryResponse)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var response types.Response
			_ = cdc.UnmarshalJSON(res, &response)
			if response.Empty() {
				response, err = utils.QueryResponseByTxQuery(cliCtx, types.QuerierRoute, params)
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
func GetCmdQueryServiceResponses(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "responses [request-context-id] [batch-counter]",
		Short: "Query active responses by the request context ID and batch counter",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query active responses by the request context ID and batch counter.

Example:
$ %s query service responses <request-context-id> <batch-counter>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

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

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryResponses)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var responses []types.Response
			if err := cdc.UnmarshalJSON(res, &responses); err != nil {
				return err
			}

			return cliCtx.PrintOutput(responses)
		},
	}

	return cmd
}

// GetCmdQueryRequestContext implements the query request context command
func GetCmdQueryRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-context [request-context-id]",
		Short: "Query a request context",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a request context.

Example:
$ %s query service request-context <request-context-id>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			params := types.QueryRequestContextParams{
				RequestContextID: requestContextID,
			}

			requestContext, err := utils.QueryRequestContext(cliCtx, types.QuerierRoute, params)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(requestContext)
		},
	}

	return cmd
}

// GetCmdQueryEarnedFees implements the query earned fees command
func GetCmdQueryEarnedFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fees [provider-address]",
		Short: "Query the earned fees of a provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the earned fees of a provider.

Example:
$ %s query service fees <provider-address>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			provider, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.QueryEarnedFeesParams{Provider: provider})
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryEarnedFees)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var fees sdk.Coins
			if err := cdc.UnmarshalJSON(res, &fees); err != nil {
				return err
			}

			return cliCtx.PrintOutput(fees)
		},
	}

	return cmd
}

// GetCmdQuerySchema implements the query schema command
func GetCmdQuerySchema(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema [schema-name]",
		Short: "Query the system schema by the schema name",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the system schema by the schema name, only pricing and result allowed.

Example:
$ %s query service schema <schema-name>
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			params := types.QuerySchemaParams{
				SchemaName: args[0],
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySchema)
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
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current service parameter values",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as service parameters.
Example:
$ %s query service params
`,
				version.AppName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.NewContext().WithCodec(cdc).WithJSONMarshaler(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
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
