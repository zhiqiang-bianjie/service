package service

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the service module.
type AppModuleBasic struct{}

// Name returns the service module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the service module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the service
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

// ValidateGenesis performs genesis state validation for the service module.
func (AppModuleBasic) ValidateGenesis(_ json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers the REST routes for the service module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
}

// GetTxCmd returns the root tx command for the service module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

// GetQueryCmd returns the root query command for the service module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

// AppModule implements an application module for the service module.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k Keeper) AppModule {
	return AppModule{
		keeper: k,
	}
}

// Name returns the service module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the service module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

// Route returns the message routing key for the service module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the service module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the service module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the service module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the service module. It returns
// no validator updates.
func (am AppModule) InitGenesis(_ sdk.Context, _ json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the service
// module.
func (am AppModule) ExportGenesis(_ sdk.Context) json.RawMessage {
	return nil
}

// BeginBlock returns the begin blocker for the service module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

// EndBlock returns the end blocker for the service module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
