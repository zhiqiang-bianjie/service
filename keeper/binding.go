package keeper

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irismod/service/types"
)

// AddServiceBinding creates a new service binding
func (k Keeper) AddServiceBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
	minRespTime uint64,
) error {
	if _, found := k.GetServiceDefinition(ctx, serviceName); !found {
		return sdkerrors.Wrap(types.ErrUnknownServiceDefinition, serviceName)
	}

	if _, found := k.GetServiceBinding(ctx, serviceName, provider); found {
		return sdkerrors.Wrap(types.ErrServiceBindingExists, "")
	}

	if err := k.validateDeposit(ctx, deposit); err != nil {
		return err
	}

	maxReqTimeout := k.MaxRequestTimeout(ctx)
	if minRespTime > uint64(maxReqTimeout) {
		return sdkerrors.Wrapf(types.ErrInvalidMinRespTime, "minimum response time [%d] must not be greater than maximum request timeout [%d]", minRespTime, maxReqTimeout)
	}

	parsedPricing, err := k.ParsePricing(ctx, pricing)
	if err != nil {
		return err
	}

	if err := types.ValidatePricing(parsedPricing); err != nil {
		return err
	}

	minDeposit := k.getMinDeposit(ctx, parsedPricing)
	if !deposit.IsAllGTE(minDeposit) {
		return sdkerrors.Wrapf(types.ErrInvalidDeposit, "insufficient deposit: minimum deposit %s, %s got", minDeposit, deposit)
	}

	// Send coins from provider's account to the deposit module account
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(
		ctx, provider, types.DepositAccName, deposit,
	); err != nil {
		return err
	}

	available := true
	disabledTime := time.Time{}

	svcBinding := types.NewServiceBinding(serviceName, provider, deposit, pricing, minRespTime, available, disabledTime)
	k.SetServiceBinding(ctx, svcBinding)

	k.SetPricing(ctx, serviceName, provider, parsedPricing)

	return nil
}

// UpdateServiceBinding updates the specified service binding
func (k Keeper) UpdateServiceBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
	minRespTime uint64,
) error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownServiceBinding, "")
	}

	updated := false

	if minRespTime != 0 {
		maxReqTimeout := k.MaxRequestTimeout(ctx)
		if minRespTime > uint64(maxReqTimeout) {
			return sdkerrors.Wrapf(types.ErrInvalidMinRespTime, "minimum response time [%d] must not be greater than maximum request timeout [%d]", minRespTime, maxReqTimeout)
		}

		binding.MinRespTime = minRespTime
		updated = true
	}

	// add the deposit
	if !deposit.Empty() {
		if err := k.validateDeposit(ctx, deposit); err != nil {
			return err
		}

		binding.Deposit = binding.Deposit.Add(deposit...)
		updated = true
	}

	parsedPricing := k.GetPricing(ctx, serviceName, provider)

	// update the pricing
	if len(pricing) != 0 {
		parsedPricing, err := k.ParsePricing(ctx, pricing)
		if err != nil {
			return err
		}

		if err := types.ValidatePricing(parsedPricing); err != nil {
			return err
		}

		binding.Pricing = pricing
		k.SetPricing(ctx, serviceName, provider, parsedPricing)

		updated = true
	}

	// only check deposit when the binding is available and updated
	if binding.Available && updated {
		minDeposit := k.getMinDeposit(ctx, parsedPricing)
		if !binding.Deposit.IsAllGTE(minDeposit) {
			return sdkerrors.Wrapf(types.ErrInvalidDeposit, "insufficient deposit: minimum deposit %s, %s got", minDeposit, binding.Deposit)
		}
	}

	if !deposit.Empty() {
		// Send coins from provider's account to the deposit module account
		if err := k.supplyKeeper.SendCoinsFromAccountToModule(
			ctx, provider, types.DepositAccName, deposit,
		); err != nil {
			return err
		}
	}

	if updated {
		k.SetServiceBinding(ctx, binding)
	}

	return nil
}

// DisableServiceBinding disables the specified service binding
func (k Keeper) DisableServiceBinding(ctx sdk.Context, serviceName string, provider sdk.AccAddress) error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownServiceBinding, "")
	}

	if !binding.Available {
		return sdkerrors.Wrap(types.ErrServiceBindingUnavailable, "")
	}

	binding.Available = false
	binding.DisabledTime = ctx.BlockHeader().Time

	k.SetServiceBinding(ctx, binding)

	return nil
}

// EnableServiceBinding enables the specified service binding
func (k Keeper) EnableServiceBinding(ctx sdk.Context, serviceName string, provider sdk.AccAddress, deposit sdk.Coins) error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownServiceBinding, "")
	}

	if binding.Available {
		return sdkerrors.Wrap(types.ErrServiceBindingAvailable, "")
	}

	// add the deposit
	if !deposit.Empty() {
		if err := k.validateDeposit(ctx, deposit); err != nil {
			return err
		}

		binding.Deposit = binding.Deposit.Add(deposit...)
	}

	minDeposit := k.getMinDeposit(ctx, k.GetPricing(ctx, serviceName, provider))
	if !binding.Deposit.IsAllGTE(minDeposit) {
		return sdkerrors.Wrapf(types.ErrInvalidDeposit, "insufficient deposit: minimum deposit %s, %s got", minDeposit, binding.Deposit)
	}

	if !deposit.Empty() {
		// Send coins from provider's account to the deposit module account
		if err := k.supplyKeeper.SendCoinsFromAccountToModule(
			ctx, provider, types.DepositAccName, deposit,
		); err != nil {
			return err
		}
	}

	binding.Available = true
	binding.DisabledTime = time.Time{}

	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundDeposit refunds the deposit from the specified service binding
func (k Keeper) RefundDeposit(ctx sdk.Context, serviceName string, provider sdk.AccAddress) error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return sdkerrors.Wrap(types.ErrUnknownServiceBinding, "")
	}

	if binding.Available {
		return sdkerrors.Wrap(types.ErrServiceBindingAvailable, "")
	}

	if binding.Deposit.IsZero() {
		return sdkerrors.Wrap(types.ErrInvalidDeposit, "the deposit of the service binding is zero")
	}

	refundableTime := binding.DisabledTime.Add(k.ArbitrationTimeLimit(ctx)).Add(k.ComplaintRetrospect(ctx))

	currentTime := ctx.BlockHeader().Time
	if currentTime.Before(refundableTime) {
		return sdkerrors.Wrapf(types.ErrIncorrectRefundTime, "%v", refundableTime)
	}

	// Send coins from the deposit module account to the provider's account
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx, types.DepositAccName, binding.Provider, binding.Deposit,
	); err != nil {
		return err
	}

	binding.Deposit = sdk.Coins{}
	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundDeposits refunds the deposits of all the service bindings
func (k Keeper) RefundDeposits(ctx sdk.Context) error {
	iterator := k.AllServiceBindingsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var binding types.ServiceBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)

		if err := k.supplyKeeper.SendCoinsFromModuleToAccount(
			ctx, types.DepositAccName, binding.Provider, binding.Deposit,
		); err != nil {
			return err
		}
	}

	return nil
}

// SetServiceBinding sets the service binding
func (k Keeper) SetServiceBinding(ctx sdk.Context, svcBinding types.ServiceBinding) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(svcBinding)
	store.Set(types.GetServiceBindingKey(svcBinding.ServiceName, svcBinding.Provider), bz)
}

// GetServiceBinding retrieves the specified service binding
func (k Keeper) GetServiceBinding(ctx sdk.Context, serviceName string, provider sdk.AccAddress) (svcBinding types.ServiceBinding, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetServiceBindingKey(serviceName, provider))
	if bz == nil {
		return svcBinding, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &svcBinding)
	return svcBinding, true
}

// ParsePricing parses the given string to Pricing
func (k Keeper) ParsePricing(ctx sdk.Context, pricing string) (p types.Pricing, err error) {
	var rawPricing types.RawPricing
	if err := json.Unmarshal([]byte(pricing), &rawPricing); err != nil {
		return p, sdkerrors.Wrapf(types.ErrInvalidPricing, "failed to unmarshal the pricing: %s", err.Error())
	}

	denom, amtStr, err := types.ParseCoinParts(rawPricing.Price)
	if err != nil {
		return p, sdkerrors.Wrapf(types.ErrInvalidPricing, "failed to parse the price: %s", err.Error())
	}

	amt, err := sdk.NewDecFromStr(amtStr)
	if err != nil {
		return p, sdkerrors.Wrapf(types.ErrInvalidPricing, fmt.Sprintf("failed to parse the price: %s", err))
	}

	token, err := k.tokenKeeper.GetToken(ctx, denom)
	if err != nil {
		return p, sdkerrors.Wrapf(types.ErrInvalidPricing, "invalid price: %s", err.Error())
	}

	price,err := token.ToMinCoin(sdk.NewDecCoinFromDec(denom,amt))
	if err != nil {
		return p, sdkerrors.Wrapf(types.ErrInvalidPricing, "invalid price: %s", err.Error())
	}
	p.Price = sdk.NewCoins(price)
	p.PromotionsByTime = rawPricing.PromotionsByTime
	p.PromotionsByVolume = rawPricing.PromotionsByVolume

	return p, nil
}

// SetPricing sets the pricing for the specified service binding
func (k Keeper) SetPricing(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	pricing types.Pricing,
) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pricing)
	store.Set(types.GetPricingKey(serviceName, provider), bz)
}

// GetPricing retrieves the pricing of the specified service binding
func (k Keeper) GetPricing(ctx sdk.Context, serviceName string, provider sdk.AccAddress) (pricing types.Pricing) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetPricingKey(serviceName, provider))
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pricing)
	return pricing
}

// SetWithdrawAddress sets the withdrawal address for the specified provider
func (k Keeper) SetWithdrawAddress(ctx sdk.Context, provider, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetWithdrawAddrKey(provider), withdrawAddr.Bytes())
}

// GetWithdrawAddress gets the withdrawal address of the specified provider
func (k Keeper) GetWithdrawAddress(ctx sdk.Context, provider sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetWithdrawAddrKey(provider))
	if bz == nil {
		return provider
	}

	return sdk.AccAddress(bz)
}

// IterateWithdrawAddresses iterates through all withdrawal addresses
func (k Keeper) IterateWithdrawAddresses(
	ctx sdk.Context,
	op func(provider sdk.AccAddress, withdrawAddress sdk.AccAddress) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.WithdrawAddrKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		providerAddress := sdk.AccAddress(iterator.Key()[1:])
		withdrawAddress := sdk.AccAddress(iterator.Value())

		if stop := op(providerAddress, withdrawAddress); stop {
			break
		}
	}
}

// ServiceBindingsIterator returns an iterator for all bindings of the specified service definition
func (k Keeper) ServiceBindingsIterator(ctx sdk.Context, serviceName string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GetBindingsSubspace(serviceName))
}

// AllServiceBindingsIterator returns an iterator for all bindings
func (k Keeper) AllServiceBindingsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ServiceBindingKey)
}

func (k Keeper) IterateServiceBindings(
	ctx sdk.Context,
	op func(binding types.ServiceBinding) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ServiceBindingKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var binding types.ServiceBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)

		if stop := op(binding); stop {
			break
		}
	}
}

// getMinDeposit gets the minimum deposit required for the service binding
func (k Keeper) getMinDeposit(ctx sdk.Context, pricing types.Pricing) sdk.Coins {
	minDepositMultiple := sdk.NewInt(k.MinDepositMultiple(ctx))
	minDepositParam := k.MinDeposit(ctx)
	baseDenom := k.BaseDenom(ctx)

	price := pricing.Price.AmountOf(baseDenom)

	// minimum deposit = max(price * minDepositMultiple, minDepositParam)
	minDeposit := sdk.NewCoins(sdk.NewCoin(baseDenom, price.Mul(minDepositMultiple)))
	if minDeposit.IsAllLT(minDepositParam) {
		minDeposit = minDepositParam
	}

	return minDeposit
}

// validateDeposit validates the given deposit
func (k Keeper) validateDeposit(ctx sdk.Context, deposit sdk.Coins) error {
	baseDenom := k.BaseDenom(ctx)

	if len(deposit) != 1 || deposit[0].Denom != baseDenom {
		return sdkerrors.Wrapf(types.ErrInvalidDeposit, "deposit only accepts %s", baseDenom)
	}

	return nil
}
