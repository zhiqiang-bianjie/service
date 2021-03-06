package service

import (
	"encoding/json"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irismod/service/types"
)

// EndBlocker handles block ending logic for service
func EndBlocker(ctx sdk.Context, k Keeper) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "endBlock").With("module", "iris/service"))

	// handler for the active request on expired
	expiredRequestHandler := func(requestID tmbytes.HexBytes, request Request) {
		if !request.SuperMode {
			_ = k.Slash(ctx, requestID)
			_ = k.RefundServiceFee(ctx, request.Consumer, request.ServiceFee)
		}

		k.DeleteActiveRequest(ctx, request.ServiceName, request.Provider, request.ExpirationHeight, requestID)
	}

	// handler for the expired request batch
	expiredRequestBatchHandler := func(requestContextID tmbytes.HexBytes, requestContext RequestContext) {
		if requestContext.BatchState != BATCHCOMPLETED {
			k.IterateActiveRequests(ctx, requestContextID, requestContext.BatchCounter, expiredRequestHandler)
			resContext := k.CompleteBatch(ctx, requestContext, requestContextID)
			requestContext = resContext
		}

		k.DeleteRequestBatchExpiration(ctx, requestContextID, ctx.BlockHeight())
		k.SetRequestContext(ctx, requestContextID, requestContext)

		if requestContext.State == COMPLETED {
			k.CompleteServiceContext(ctx, requestContext, requestContextID)
		}

		if requestContext.State == RUNNING {
			if requestContext.Repeated && (requestContext.RepeatedTotal < 0 || int64(requestContext.BatchCounter) < requestContext.RepeatedTotal) {
				k.AddNewRequestBatch(ctx, requestContextID, ctx.BlockHeight()-requestContext.Timeout+int64(requestContext.RepeatedFrequency))
			} else {
				k.CompleteServiceContext(ctx, requestContext, requestContextID)
			}
		}

		k.CleanBatch(ctx, requestContext, requestContextID)
	}

	providerRequests := make(map[string][]string)

	// handler for the new request batch
	newRequestBatchHandler := func(requestContextID tmbytes.HexBytes, requestContext RequestContext) {
		if requestContext.State == RUNNING {
			providers, totalPrices := k.FilterServiceProviders(
				ctx, requestContext.ServiceName,
				requestContext.Providers,
				requestContext.Timeout,
				requestContext.ServiceFeeCap,
				requestContext.Consumer,
			)

			if len(providers) > 0 && len(providers) >= int(requestContext.ResponseThreshold) {
				if !requestContext.SuperMode {
					if err := k.DeductServiceFees(ctx, requestContext.Consumer, totalPrices); err != nil {
						k.OnRequestContextPaused(ctx, requestContext, requestContextID, "insufficient balances")
					}
				}

				if requestContext.State == RUNNING {
					k.InitiateRequests(ctx, requestContextID, providers, providerRequests)
					k.AddRequestBatchExpiration(ctx, requestContextID, ctx.BlockHeight()+requestContext.Timeout)
				}
			} else {
				k.SkipCurrentRequestBatch(ctx, requestContextID, requestContext)
			}

			requestContext, _ := k.GetRequestContext(ctx, requestContextID)
			batchState := types.BatchState{
				BatchCounter:           requestContext.BatchCounter,
				State:                  requestContext.BatchState,
				BatchResponseThreshold: requestContext.BatchResponseThreshold,
				BatchRequestCount:      requestContext.BatchRequestCount,
				BatchResponseCount:     requestContext.BatchResponseCount,
			}
			stateJSON, _ := json.Marshal(batchState)

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeNewBatch,
					sdk.NewAttribute(types.AttributeKeyRequestContextID, requestContextID.String()),
					sdk.NewAttribute(types.AttributeKeyRequestContextState, string(stateJSON)),
				),
			})
		}

		k.DeleteNewRequestBatch(ctx, requestContextID, ctx.BlockHeight())
	}

	// handle the expired request batch queue
	k.IterateExpiredRequestBatch(ctx, ctx.BlockHeight(), expiredRequestBatchHandler)

	// handle the new request batch queue
	k.IterateNewRequestBatch(ctx, ctx.BlockHeight(), newRequestBatchHandler)

	for provider, requests := range providerRequests {
		requestsJSON, _ := json.Marshal(requests)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeNewBatchRequest,
				sdk.NewAttribute(types.AttributeKeyProvider, provider),
				sdk.NewAttribute(types.AttributeKeyRequests, string(requestsJSON)),
			),
		})
	}
}
