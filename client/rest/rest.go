package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// Rest variable names
// nolint
const (
	RestServiceName = "service-name"
	RestRequestID   = "request-id"
	RestProvider    = "provider"
	RestConsumer    = "consumer"
	RestAddress     = "address"
)

// RegisterRoutes defines routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
