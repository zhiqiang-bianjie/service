package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/irismod/service/types"
)

// Keeper defines the service keeper
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	supplyKeeper types.SupplyKeeper
	tokenKeeper  types.TokenKeeper
	paramstore   params.Subspace

	feeCollectorName string // name of the FeeCollector ModuleAccount

	// used to map the module name to response callback
	respCallbacks map[string]types.ResponseCallback

	// used to map the module name to state callback
	stateCallbacks map[string]types.StateCallback
}

// NewKeeper creates a new service Keeper instance
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	supplyKeeper types.SupplyKeeper,
	tokenKeeper types.TokenKeeper,
	paramstore params.Subspace,
	feeCollectorName string,
) Keeper {
	// ensure service module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.DepositAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.DepositAccName))
	}

	if addr := supplyKeeper.GetModuleAddress(types.RequestAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.RequestAccName))
	}

	keeper := Keeper{
		storeKey:         key,
		cdc:              cdc,
		supplyKeeper:     supplyKeeper,
		tokenKeeper:      tokenKeeper,
		feeCollectorName: feeCollectorName,
		paramstore:       paramstore.WithKeyTable(ParamKeyTable()),
	}

	keeper.respCallbacks = make(map[string]types.ResponseCallback)
	keeper.stateCallbacks = make(map[string]types.StateCallback)
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// GetServiceDepositAccount returns the service depost ModuleAccount
func (k Keeper) GetServiceDepositAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.DepositAccName)
}

// GetServiceRequestAccount returns the service request ModuleAccount
func (k Keeper) GetServiceRequestAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.RequestAccName)
}
