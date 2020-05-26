package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/irismod/service/types"
)

// QueryRequestContext queries a single request context
func QueryRequestContext(cliCtx context.CLIContext, queryRoute string, params types.QueryRequestContextParams) (
	requestContext types.RequestContext, err error) {
	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return requestContext, err
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequestContext)
	res, _, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return requestContext, err
	}

	_ = cliCtx.Codec.UnmarshalJSON(res, &requestContext)
	if requestContext.Empty() {
		requestContext, err = QueryRequestContextByTxQuery(cliCtx, queryRoute, params)
		if err != nil {
			return requestContext, err
		}
	}

	if requestContext.Empty() {
		return requestContext, fmt.Errorf("unknown request context: %s", hex.EncodeToString(params.RequestContextID))
	}
	return requestContext, nil
}

// QueryRequestContextByTxQuery will query for a single request context via a direct txs tags query.
func QueryRequestContextByTxQuery(cliCtx context.CLIContext, queryRoute string, params types.QueryRequestContextParams) (
	requestContext types.RequestContext, err error) {
	txHash, msgIndex, err := types.SplitRequestContextID(params.RequestContextID)
	if err != nil {
		return requestContext, err
	}

	// NOTE: QueryTx is used to facilitate the txs query which does not currently
	txInfo, err := authclient.QueryTx(cliCtx, txHash.String())
	if err != nil {
		return requestContext, err
	}

	if int64(len(txInfo.Tx.GetMsgs())) > msgIndex {
		msg := txInfo.Tx.GetMsgs()[msgIndex]
		if msg.Type() == types.TypeMsgCallService {
			requestMsg := msg.(types.MsgCallService)
			requestContext := types.NewRequestContext(
				requestMsg.ServiceName, requestMsg.Providers,
				requestMsg.Consumer, requestMsg.Input, requestMsg.ServiceFeeCap,
				requestMsg.Timeout, requestMsg.SuperMode, requestMsg.Repeated,
				requestMsg.RepeatedFrequency, requestMsg.RepeatedTotal,
				uint64(requestMsg.RepeatedTotal), 0, 0, 0,
				types.BATCHCOMPLETED, types.COMPLETED, 0, "",
			)

			return requestContext, nil
		}
	}

	return requestContext, nil
}

// QueryRequestByTxQuery will query for a single request via a direct txs tags query.
func QueryRequestByTxQuery(cliCtx context.CLIContext, queryRoute string, params types.QueryRequestParams) (
	request types.Request, err error) {
	requestID := params.RequestID
	if err != nil {
		return request, nil
	}

	contextID, _, requestHeight, batchRequestIndex, err := types.SplitRequestID(requestID)
	if err != nil {
		return request, err
	}

	// query request context
	requestContext, err := QueryRequestContext(cliCtx, queryRoute, types.QueryRequestContextParams{
		RequestContextID: contextID,
	})

	if err != nil {
		return request, err
	}

	// query batch request by requestHeight
	node, err := cliCtx.GetNode()
	if err != nil {
		return request, err
	}

	blockResult, err := node.BlockResults(&requestHeight)
	if err != nil {
		return request, err
	}

	for _, event := range blockResult.EndBlockEvents {
		if event.Type == types.EventTypeNewBatchRequest {
			var found bool
			var requests []types.CompactRequest
			var requestsBz []byte
			for _, attribute := range event.Attributes {
				if string(attribute.Key) == types.AttributeKeyRequests {
					requestsBz = attribute.GetValue()
				}
				if string(attribute.Key) == types.AttributeKeyRequestContextID &&
					string(attribute.GetValue()) == contextID.String() {
					found = true
				}
			}
			if found {
				err := json.Unmarshal(requestsBz, &requests)
				if err != nil {
					return request, err
				}

				if len(requests) > int(batchRequestIndex) {
					compactRequest := requests[batchRequestIndex]
					request = types.NewRequest(
						requestID,
						requestContext.ServiceName,
						compactRequest.Provider,
						requestContext.Consumer,
						requestContext.Input,
						compactRequest.ServiceFee,
						requestContext.SuperMode,
						compactRequest.RequestHeight,
						compactRequest.ExpirationHeight,
						compactRequest.RequestContextID,
						compactRequest.RequestContextBatchCounter,
					)

					return request, nil
				}
			}
		}
	}

	return request, nil
}

// QueryResponseByTxQuery will query for a single request via a direct txs tags query.
func QueryResponseByTxQuery(cliCtx context.CLIContext, queryRoute string, params types.QueryResponseParams) (
	response types.Response, err error) {

	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgRespondService),
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, types.AttributeKeyRequestID, []byte(fmt.Sprintf("%d", params.RequestID))),
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	result, err := authclient.QueryTxsByEvents(cliCtx, events, 1, 1, "")
	if err != nil {
		return response, err
	}

	if len(result.Txs) == 0 {
		return response, fmt.Errorf("unknown response: %s", params.RequestID)
	}

	requestID := params.RequestID

	contextID, batchCounter, _, _, err := types.SplitRequestID(requestID)
	if err != nil {
		return response, err
	}

	// query request context
	requestContext, err := QueryRequestContext(cliCtx, queryRoute, types.QueryRequestContextParams{
		RequestContextID: contextID,
	})

	if err != nil {
		return response, err
	}

	for _, msg := range result.Txs[0].Tx.GetMsgs() {
		if msg.Type() == types.TypeMsgRespondService {
			responseMsg := msg.(types.MsgRespondService)
			if responseMsg.RequestID.String() != params.RequestID.String() {
				continue
			}

			response := types.NewResponse(
				responseMsg.Provider, requestContext.Consumer,
				responseMsg.Result, responseMsg.Output,
				contextID, batchCounter,
			)

			return response, nil
		}
	}

	return response, nil
}

// QueryRequestsByBinding queries active requests by the service binding
func QueryRequestsByBinding(cliCtx context.CLIContext, queryRoute string, serviceName string, provider sdk.AccAddress) ([]types.Request, int64, error) {
	params := types.QueryRequestsParams{
		ServiceName: serviceName,
		Provider:    provider,
	}

	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return nil, 0, err
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequests)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return nil, 0, err
	}

	var requests []types.Request
	if err := cliCtx.Codec.UnmarshalJSON(res, &requests); err != nil {
		return nil, 0, err
	}

	return requests, height, nil
}

// QueryRequestsByReqCtx queries active requests by the request context ID
func QueryRequestsByReqCtx(cliCtx context.CLIContext, queryRoute, reqCtxIDStr, batchCounterStr string) ([]types.Request, int64, error) {
	requestContextID, err := hex.DecodeString(reqCtxIDStr)
	if err != nil {
		return nil, 0, err
	}

	batchCounter, err := strconv.ParseUint(batchCounterStr, 10, 64)
	if err != nil {
		return nil, 0, err
	}

	params := types.QueryRequestsByReqCtxParams{
		RequestContextID: requestContextID,
		BatchCounter:     batchCounter,
	}

	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return nil, 0, err
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequestsByReqCtx)
	res, height, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return nil, 0, err
	}

	var requests []types.Request
	if err := cliCtx.Codec.UnmarshalJSON(res, &requests); err != nil {
		return nil, 0, err
	}

	return requests, height, nil
}
