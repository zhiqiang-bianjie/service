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
func (k Keeper) AddEarnedFee(ctx sdk.Context, provider sdk.AccAddress, fee sdk.Coins) error {
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

	fees, _ := k.GetEarnedFees(ctx, provider)
	k.SetEarnedFees(ctx, provider, fees.Coins.Add(earnedFee...))

	return nil
}

// SetEarnedFees sets the earned fees for the specified provider
func (k Keeper) SetEarnedFees(ctx sdk.Context, provider sdk.AccAddress, fees sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	earnedFees := types.NewEarnedFees(provider, fees)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(earnedFees)

	store.Set(types.GetEarnedFeesKey(provider), bz)
}

// DeleteEarnedFees removes the earned fees of the specified provider
func (k Keeper) DeleteEarnedFees(ctx sdk.Context, provider sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetEarnedFeesKey(provider))
}

// GetEarnedFees retrieves the earned fees of the specified provider
func (k Keeper) GetEarnedFees(ctx sdk.Context, provider sdk.AccAddress) (fees types.EarnedFees, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetEarnedFeesKey(provider))
	if bz == nil {
		return fees, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &fees)
	return fees, true
}

// WithdrawEarnedFees withdraws the earned fees of the specified provider
func (k Keeper) WithdrawEarnedFees(ctx sdk.Context, provider sdk.AccAddress) error {
	fees, found := k.GetEarnedFees(ctx, provider)
	if !found {
		return sdkerrors.Wrapf(types.ErrNoEarnedFees, provider.String())
	}

	withdrawAddr := k.GetWithdrawAddress(ctx, provider)

	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, withdrawAddr, fees.Coins)
	if err != nil {
		return err
	}

	k.DeleteEarnedFees(ctx, provider)

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
		var earnedFees types.EarnedFees
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &earnedFees)

		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.RequestAccName, earnedFees.Address, earnedFees.Coins)
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
