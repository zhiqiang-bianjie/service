package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	serviceutils "github.com/irismod/service/client/utils"
	"github.com/irismod/service/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// query definition
	r.HandleFunc(fmt.Sprintf("/service/definitions/{%s}", RestServiceName), queryDefinitionHandlerFn(cliCtx)).Methods("GET")
	// query binding
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}", RestServiceName, RestProvider), queryBindingHandlerFn(cliCtx)).Methods("GET")
	// query bindings
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}", RestServiceName), queryBindingsHandlerFn(cliCtx)).Methods("GET")
	// query the withdrawal address
	r.HandleFunc(fmt.Sprintf("/service/providers/{%s}/withdraw-address", RestProvider), queryWithdrawAddrHandlerFn(cliCtx)).Methods("GET")
	// query a request by ID
	r.HandleFunc(fmt.Sprintf("/service/requests/{%s}", RestRequestID), queryRequestHandlerFn(cliCtx)).Methods("GET")
	// query active requests by the service binding or request context ID
	r.HandleFunc(fmt.Sprintf("/service/requests/{%s}/{%s}", RestArg1, RestArg2), queryRequestsHandlerFn(cliCtx)).Methods("GET")
	// query a response
	r.HandleFunc(fmt.Sprintf("/service/responses/{%s}", RestRequestID), queryResponseHandlerFn(cliCtx)).Methods("GET")
	// query a request context
	r.HandleFunc(fmt.Sprintf("/service/contexts/{%s}", RestRequestContextID), queryRequestContextHandlerFn(cliCtx)).Methods("GET")
	// query active responses by the request context ID and batch counter
	r.HandleFunc(fmt.Sprintf("/service/responses/{%s}/{%s}", RestRequestContextID, RestBatchCounter), queryResponsesHandlerFn(cliCtx)).Methods("GET")
	// query the earned fees of a provider
	r.HandleFunc(fmt.Sprintf("/service/fees/{%s}", RestProvider), queryEarnedFeesHandlerFn(cliCtx)).Methods("GET")
	// query the system schema by the schema name
	r.HandleFunc(fmt.Sprintf("/service/schemas/{%s}", RestSchemaName), querySchemaHandlerFn(cliCtx)).Methods("GET")
	// query the current service parameter values
	r.HandleFunc(fmt.Sprintf("/service/parameters"), queryParamsHandlerFn(cliCtx)).Methods("GET")
}

func queryDefinitionHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]

		if err := types.ValidateServiceName(serviceName); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryDefinitionParams{
			ServiceName: serviceName,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDefinition)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryBindingHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		providerStr := vars[RestProvider]

		if err := types.ValidateServiceName(serviceName); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryBindingParams{
			ServiceName: serviceName,
			Provider:    provider,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryBinding)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryBindingsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]

		if err := types.ValidateServiceName(serviceName); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryBindingsParams{
			ServiceName: serviceName,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryBindings)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryWithdrawAddrHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryWithdrawAddressParams{
			Provider: provider,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryWithdrawAddress)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestIDStr := vars[RestRequestID]

		requestID, err := types.ConvertRequestID(requestIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryRequestParams{
			RequestID: requestID,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryRequest)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var request types.Request
		_ = cliCtx.Codec.UnmarshalJSON(res, &request)
		if request.Empty() {
			request, err = serviceutils.QueryRequestByTxQuery(cliCtx, types.RouterKey, params)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if request.Empty() {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("unknown request: %s", params.RequestID))
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryRequestsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		arg1 := vars[RestArg1]
		arg2 := vars[RestArg2]

		queryByBinding := true

		provider, err := sdk.AccAddressFromBech32(arg2)
		if err != nil {
			queryByBinding = false
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		var res types.Requests
		var height int64

		if queryByBinding {
			res, height, err = serviceutils.QueryRequestsByBinding(cliCtx, types.RouterKey, arg1, provider)
		} else {
			res, height, err = serviceutils.QueryRequestsByReqCtx(cliCtx, types.RouterKey, arg1, arg2)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryResponseHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestIDStr := vars[RestRequestID]

		requestID, err := types.ConvertRequestID(requestIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryResponseParams{
			RequestID: requestID,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryResponse)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var response types.Response
		_ = cliCtx.Codec.UnmarshalJSON(res, &response)
		if response.Empty() {
			response, err = serviceutils.QueryResponseByTxQuery(cliCtx, types.RouterKey, params)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if response.Empty() {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("unknown request: %s", params.RequestID))
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, response)
	}
}

func queryRequestContextHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RestRequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryRequestContextParams{
			RequestContextID: requestContextID,
		}

		requestContext, err := serviceutils.QueryRequestContext(cliCtx, types.RouterKey, params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, requestContext)
	}
}

func queryResponsesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RestRequestContextID]
		batchCounterStr := vars[RestBatchCounter]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		batchCounter, err := strconv.ParseUint(batchCounterStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryResponsesParams{
			RequestContextID: requestContextID,
			BatchCounter:     batchCounter,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryResponses)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryEarnedFeesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryEarnedFeesParams{
			Provider: provider,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryEarnedFees)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func querySchemaHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QuerySchemaParams{
			SchemaName: vars[RestSchemaName],
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QuerySchema)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var schema serviceutils.SchemaType
		if err := cliCtx.Codec.UnmarshalJSON(res, &schema); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, schema)
	}
}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
