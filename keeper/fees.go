package keeper

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irismod/service/types"
)

// RefundServiceFee refunds the service fee to the specified consumer
func (k Keeper) RefundServiceFee(ctx sdk.Context, consumer sdk.AccAddress, serviceFee sdk.Coins) error {
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, consumer, serviceFee)

	if err != nil {
		return err
	}

	return nil
}

// AddEarnedFee adds the earned fee for the given provider
func (k Keeper) AddEarnedFee(ctx sdk.Context, owner, provider sdk.AccAddress, fee sdk.Coins) error {
	taxRate := k.ServiceFeeTax(ctx)

	taxCoins := sdk.Coins{}
	for _, coin := range fee {
		taxAmount := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).TruncateInt()
		taxCoins = taxCoins.Add(sdk.NewCoin(coin.Denom, taxAmount))
	}

	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.RequestAccName, k.feeCollectorName, taxCoins)
	if err != nil {
		return err
	}

	earnedFee, hasNeg := fee.SafeSub(taxCoins)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "%s is less than %s", fee, taxCoins)
	}

	// add the provider's earned fees
	earnedFees, _ := k.GetEarnedFees(ctx, provider)
	k.SetEarnedFees(ctx, provider, earnedFees.Add(earnedFee...))

	// add the owner's earned fees
	ownerEarnedFees, _ := k.GetOwnerEarnedFees(ctx, owner)
	k.SetOwnerEarnedFees(ctx, owner, ownerEarnedFees.Add(earnedFee...))

	return nil
}

// SetEarnedFees sets the earned fees for the specified provider
func (k Keeper) SetEarnedFees(ctx sdk.Context, provider sdk.AccAddress, fees sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(fees)
	store.Set(types.GetEarnedFeesKey(provider), bz)
}

// GetEarnedFees retrieves the earned fees of the specified provider
func (k Keeper) GetEarnedFees(ctx sdk.Context, provider sdk.AccAddress) (fees sdk.Coins, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetEarnedFeesKey(provider))
	if bz == nil {
		return fees, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &fees)
	return fees, true
}

// DeleteEarnedFees removes the earned fees of the specified provider
func (k Keeper) DeleteEarnedFees(ctx sdk.Context, provider sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetEarnedFeesKey(provider))
}

// SetOwnerEarnedFees sets the earned fees for the specified owner
func (k Keeper) SetOwnerEarnedFees(ctx sdk.Context, owner sdk.AccAddress, fees sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(fees)
	store.Set(types.GetOwnerEarnedFeesKey(owner), bz)
}

// GetOwnerEarnedFees retrieves the earned fees of the specified owner
func (k Keeper) GetOwnerEarnedFees(ctx sdk.Context, owner sdk.AccAddress) (fees sdk.Coins, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetOwnerEarnedFeesKey(owner))
	if bz == nil {
		return fees, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &fees)
	return fees, true
}

// DeleteOwnerEarnedFees removes the earned fees of the specified owner
func (k Keeper) DeleteOwnerEarnedFees(ctx sdk.Context, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOwnerEarnedFeesKey(owner))
}

// WithdrawEarnedFees withdraws the earned fees of the specified provider or owner
func (k Keeper) WithdrawEarnedFees(ctx sdk.Context, owner, provider sdk.AccAddress) error {
	ownerEarnedFees, found := k.GetOwnerEarnedFees(ctx, owner)
	if !found || ownerEarnedFees.IsZero() {
		return sdkerrors.Wrap(types.ErrNoEarnedFees, owner.String())
	}

	var withdrawFees sdk.Coins

	if !provider.Empty() {
		earnedFees, found := k.GetEarnedFees(ctx, provider)
		if !found {
			return sdkerrors.Wrap(types.ErrNoEarnedFees, provider.String())
		}

		k.DeleteEarnedFees(ctx, provider)
		k.SetOwnerEarnedFees(ctx, owner, ownerEarnedFees.Sub(earnedFees))

		withdrawFees = earnedFees
	} else {
		iterator := k.OwnerProvidersIterator(ctx, owner)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			provider := sdk.AccAddress(iterator.Key()[sdk.AddrLen+1:])
			k.DeleteEarnedFees(ctx, provider)
		}

		k.DeleteOwnerEarnedFees(ctx, owner)
		withdrawFees = ownerEarnedFees
	}

	withdrawAddr := k.GetWithdrawAddress(ctx, owner)

	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, withdrawAddr, withdrawFees)
	if err != nil {
		return err
	}

	return nil
}

// AllEarnedFeesIterator returns an iterator for all the earned fees
func (k Keeper) AllEarnedFeesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.EarnedFeesKey)
}

// RefundEarnedFees refunds all the earned fees
func (k Keeper) RefundEarnedFees(ctx sdk.Context) error {
	iterator := k.AllEarnedFeesIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		provider := iterator.Key()[1:]

		var earnedFees sdk.Coins
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &earnedFees)

		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, provider, earnedFees)
		if err != nil {
			return err
		}
	}

	return nil
}

// RefundServiceFees refunds the service fees of all the active requests
func (k Keeper) RefundServiceFees(ctx sdk.Context) error {
	iterator := k.AllActiveRequestsIterator(ctx.KVStore(k.storeKey))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var requestID tmbytes.HexBytes
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &requestID)

		request, _ := k.GetRequest(ctx, requestID)

		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, request.Consumer, request.ServiceFee)
		if err != nil {
			return err
		}
	}

	return nil
}
