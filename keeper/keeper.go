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
}

// NewKeeper creates a new service Keeper instance
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	supplyKeeper types.SupplyKeeper,
	tokenKeeper types.TokenKeeper,
	paramstore params.Subspace,
) Keeper {
	// ensure service module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.DepositAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.DepositAccName))
	}

	if addr := supplyKeeper.GetModuleAddress(types.RequestAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.RequestAccName))
	}

	if addr := supplyKeeper.GetModuleAddress(types.TaxAccName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.TaxAccName))
	}

	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
		tokenKeeper:  tokenKeeper,
		paramstore:   paramstore.WithKeyTable(ParamKeyTable()),
	}
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

// GetServiceTaxAccount returns the service tax ModuleAccount
func (k Keeper) GetServiceTaxAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.TaxAccName)
}
