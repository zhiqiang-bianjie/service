package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ServiceBinding defines a struct for the service binding
type ServiceBinding struct {
	ServiceName  string         `json:"service_name" yaml:"service_name"`
	Provider     sdk.AccAddress `json:"provider" yaml:"provider"`
	Deposit      sdk.Coins      `json:"deposit" yaml:"deposit"`
	Pricing      string         `json:"pricing" yaml:"pricing"`
	MinRespTime  uint64         `json:"min_resp_time" yaml:"min_resp_time"`
	Available    bool           `json:"available" yaml:"available"`
	DisabledTime time.Time      `json:"disabled_time" yaml:"disabled_time"`
}

// NewServiceBinding creates a new ServiceBinding instance
func NewServiceBinding(
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
	minRespTime uint64,
	available bool,
	disabledTime time.Time,
) ServiceBinding {
	return ServiceBinding{
		ServiceName:  serviceName,
		Provider:     provider,
		Deposit:      deposit,
		Pricing:      pricing,
		MinRespTime:  minRespTime,
		Available:    available,
		DisabledTime: disabledTime,
	}
}

// RawPricing represents the raw pricing of a service binding
type RawPricing struct {
	Price              string              `json:"price"`                // base price string
	PromotionsByTime   []PromotionByTime   `json:"promotions_by_time"`   // promotions by time
	PromotionsByVolume []PromotionByVolume `json:"promotions_by_volume"` // promotions by volume
}

// Pricing represents the pricing of a service binding
type Pricing struct {
	Price              sdk.Coins           `json:"price"`                // base price
	PromotionsByTime   []PromotionByTime   `json:"promotions_by_time"`   // promotions by time
	PromotionsByVolume []PromotionByVolume `json:"promotions_by_volume"` // promotions by volume
}

// PromotionByTime defines the promotion by time
type PromotionByTime struct {
	StartTime time.Time `json:"start_time"` // starting time of the promotion
	EndTime   time.Time `json:"end_time"`   // ending time of the promotion
	Discount  sdk.Dec   `json:"discount"`   // discount during the promotion
}

// PromotionByVolume defines the promotion by volume
type PromotionByVolume struct {
	Volume   uint64  `json:"volume"`   // minimum volume for the promotion
	Discount sdk.Dec `json:"discount"` // discount for the promotion
}

// GetDiscountByTime gets the discount level by the specified time
func GetDiscountByTime(pricing Pricing, time time.Time) sdk.Dec {
	for _, p := range pricing.PromotionsByTime {
		if !time.Before(p.StartTime) && time.Before(p.EndTime) {
			return p.Discount
		}
	}

	return sdk.OneDec()
}

// GetDiscountByVolume gets the discount level by the specified volume
// Note: make sure that the promotions by volume are sorted in ascending order
func GetDiscountByVolume(pricing Pricing, volume uint64) sdk.Dec {
	promotionsByVol := pricing.PromotionsByVolume

	for i, p := range promotionsByVol {
		if volume < p.Volume {
			if i == 0 {
				return sdk.OneDec()
			}

			return promotionsByVol[i-1].Discount
		}

		if i == len(promotionsByVol)-1 {
			return p.Discount
		}
	}

	return sdk.OneDec()
}

// ValidatePricing validates the given pricing
func ValidatePricing(pricing Pricing) error {
	// CONTRACT:
	// p.EndTime > p.StartTime
	// p[i].StartTime >= p[i-1].EndTime
	for i, p := range pricing.PromotionsByTime {
		if !p.EndTime.After(p.StartTime) ||
			(i > 0 && p.StartTime.Before(pricing.PromotionsByTime[i-1].EndTime)) {
			return sdkerrors.Wrapf(ErrInvalidPricing, "invalid timing promotion %d", i)
		}
	}

	// CONTRACT:
	// p[i].Volume > p[i-1].Volume
	for i, p := range pricing.PromotionsByVolume {
		if i > 0 && p.Volume < pricing.PromotionsByVolume[i-1].Volume {
			return sdkerrors.Wrapf(ErrInvalidPricing, "invalid volume promotion %d", i)
		}
	}

	return nil
}

// Validate validates the service binding
func (binding ServiceBinding) Validate() error {
	if err := ValidateProvider(binding.Provider); err != nil {
		return err
	}

	if err := ValidateServiceName(binding.ServiceName); err != nil {
		return err
	}

	if err := ValidateServiceDeposit(binding.Deposit); err != nil {
		return err
	}

	if err := ValidateMinRespTime(binding.MinRespTime); err != nil {
		return err
	}

	return ValidateBindingPricing(binding.Pricing)
}
