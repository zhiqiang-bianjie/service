package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/irismod/service/types"
)

// Keeper defines the service keeper
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	supplyKeeper types.SupplyKeeper
	paramstore   params.Subspace
}

// NewKeeper creates a new service Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, supplyKeeper types.SupplyKeeper, paramstore params.Subspace,
) Keeper {
	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
		paramstore:   paramstore,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}
