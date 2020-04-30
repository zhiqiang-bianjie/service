package keeper

import (
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irismod/service/types"
)

// NewQuerier creates a new service Querier instance
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryDefinition:
			return queryServiceDefinition(ctx, path[1:], req, k)

		case types.QueryBinding:
			return queryBinding(ctx, req, k)

		case types.QueryBindings:
			return queryBindings(ctx, req, k)

		case types.QueryWithdrawAddress:
			return queryWithdrawAddress(ctx, req, k)

		case types.QueryRequest:
			return queryRequest(ctx, req, k)

		case types.QueryRequests:
			return queryRequests(ctx, req, k)

		case types.QueryResponse:
			return queryResponse(ctx, req, k)

		case types.QueryRequestContext:
			return queryRequestContext(ctx, req, k)

		case types.QueryRequestsByReqCtx:
			return queryRequestsByReqCtx(ctx, req, k)

		case types.QueryResponses:
			return queryResponses(ctx, req, k)

		case types.QueryEarnedFees:
			return queryEarnedFees(ctx, req, k)

		case types.QuerySchema:
			return querySchema(ctx, req, k)

		case types.QueryParameters:
			return queryParams(ctx, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query path: %s", types.ModuleName, path[0])
		}
	}
}

func queryServiceDefinition(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDefinitionParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	definition, found := k.GetServiceDefinition(ctx, params.ServiceName)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownServiceDefinition, params.ServiceName)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, definition)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryBinding(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryBindingParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	svcBinding, found := k.GetServiceBinding(ctx, params.ServiceName, params.Provider)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownServiceBinding, "")
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, svcBinding)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryBindings(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryBindingsParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	bindings := make([]types.ServiceBinding, 0)

	if params.Owner.Empty() {
		iterator := k.ServiceBindingsIterator(ctx, params.ServiceName)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			var binding types.ServiceBinding
			k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)

			bindings = append(bindings, binding)
		}
	} else {
		bindings = k.GetOwnerServiceBindings(ctx, params.Owner, params.ServiceName)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, bindings)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryWithdrawAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryWithdrawAddressParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	withdrawAddr := k.GetWithdrawAddress(ctx, params.Owner)

	bz, err := codec.MarshalJSONIndent(k.cdc, withdrawAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRequest(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryRequestParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if len(params.RequestID) != types.RequestIDLen {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequestID, "invalid length, expected: %d, got: %d",
			types.RequestIDLen, len(params.RequestID))
	}

	request, _ := k.GetRequest(ctx, params.RequestID)

	bz, err := codec.MarshalJSONIndent(k.cdc, request)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRequests(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryRequestsParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	iterator := k.ActiveRequestsIterator(ctx, params.ServiceName, params.Provider)
	defer iterator.Close()

	requests := make([]types.Request, 0)

	for ; iterator.Valid(); iterator.Next() {
		var requestID tmbytes.HexBytes
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &requestID)

		request, _ := k.GetRequest(ctx, requestID)
		requests = append(requests, request)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, requests)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryResponse(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryResponseParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if len(params.RequestID) != types.RequestIDLen {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequestID, "invalid length, expected: %d, got: %d",
			types.RequestIDLen, len(params.RequestID))
	}

	response, _ := k.GetResponse(ctx, params.RequestID)

	bz, err := codec.MarshalJSONIndent(k.cdc, response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRequestContext(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryRequestContextParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	requestContext, _ := k.GetRequestContext(ctx, params.RequestContextID)
	bz, err := codec.MarshalJSONIndent(k.cdc, requestContext)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryRequestsByReqCtx(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryRequestsByReqCtxParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	iterator := k.RequestsIteratorByReqCtx(ctx, params.RequestContextID, params.BatchCounter)
	defer iterator.Close()

	requests := make([]types.Request, 0)

	for ; iterator.Valid(); iterator.Next() {
		requestID := iterator.Key()[1:]
		request, _ := k.GetRequest(ctx, requestID)

		requests = append(requests, request)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, requests)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryResponses(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryResponsesParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	iterator := k.ResponsesIteratorByReqCtx(ctx, params.RequestContextID, params.BatchCounter)
	defer iterator.Close()

	responses := make([]types.Response, 0)

	for ; iterator.Valid(); iterator.Next() {
		var response types.Response
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &response)

		responses = append(responses, response)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, responses)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryEarnedFees(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryEarnedFeesParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	fees, found := k.GetEarnedFees(ctx, params.Provider)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrNoEarnedFees, "no earned fees for %s",
			params.Provider.String())
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, fees)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func querySchema(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QuerySchemaParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var schemaName = strings.ToLower(params.SchemaName)
	var schema string

	if schemaName == "pricing" {
		schema = types.PricingSchema
	} else if schemaName == "result" {
		schema = types.ResultSchema
	} else {
		return nil, sdkerrors.Wrap(types.ErrInvalidSchemaName, schema)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, schema)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	bz, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
