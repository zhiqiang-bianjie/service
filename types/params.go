package types

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Service params default values
var (
	DefaultMaxRequestTimeout    = int64(100)
	DefaultMinDepositMultiple   = int64(200)
	DefaultMinDeposit           = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(6000))) // 6000stake
	DefaultServiceFeeTax        = sdk.NewDecWithPrec(1, 1)                                          // 10%
	DefaultSlashFraction        = sdk.NewDecWithPrec(1, 3)                                          // 0.1%
	DefaultComplaintRetrospect  = 15 * 24 * time.Hour                                               // 15 days
	DefaultArbitrationTimeLimit = 5 * 24 * time.Hour                                                // 5 days
	DefaultTxSizeLimit          = uint64(4000)
	DefaultBaseDenom            = sdk.DefaultBondDenom
)

// Keys for parameter access
// nolint
var (
	KeyMaxRequestTimeout    = []byte("MaxRequestTimeout")
	KeyMinDepositMultiple   = []byte("MinDepositMultiple")
	KeyMinDeposit           = []byte("MinDeposit")
	KeyServiceFeeTax        = []byte("ServiceFeeTax")
	KeySlashFraction        = []byte("SlashFraction")
	KeyComplaintRetrospect  = []byte("ComplaintRetrospect")
	KeyArbitrationTimeLimit = []byte("ArbitrationTimeLimit")
	KeyTxSizeLimit          = []byte("TxSizeLimit")
	KeyBaseDenom            = []byte("BaseDenom")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for service
type Params struct {
	MaxRequestTimeout    int64         `json:"max_request_timeout" yaml:"max_request_timeout"`
	MinDepositMultiple   int64         `json:"min_deposit_multiple" yaml:"min_deposit_multiple"`
	MinDeposit           sdk.Coins     `json:"min_deposit" yaml:"min_deposit"`
	ServiceFeeTax        sdk.Dec       `json:"service_fee_tax" yaml:"service_fee_tax"`
	SlashFraction        sdk.Dec       `json:"slash_fraction" yaml:"slash_fraction"`
	ComplaintRetrospect  time.Duration `json:"complaint_retrospect" yaml:"complaint_retrospect"`
	ArbitrationTimeLimit time.Duration `json:"arbitration_time_limit" yaml:"arbitration_time_limit"`
	TxSizeLimit          uint64        `json:"tx_size_limit" yaml:"tx_size_limit"`
	BaseDenom            string        `json:"base_denom" yaml:"base_denom"`
}

// NewParams creates a new Params instance
func NewParams(
	maxRequestTimeout,
	minDepositMultiple int64,
	minDeposit sdk.Coins,
	serviceFeeTax,
	slashFraction sdk.Dec,
	complaintRetrospect,
	arbitrationTimeLimit time.Duration,
	txSizeLimit uint64,
	baseDenom string,
) Params {
	return Params{
		MaxRequestTimeout:    maxRequestTimeout,
		MinDepositMultiple:   minDepositMultiple,
		MinDeposit:           minDeposit,
		ServiceFeeTax:        serviceFeeTax,
		SlashFraction:        slashFraction,
		ComplaintRetrospect:  complaintRetrospect,
		ArbitrationTimeLimit: arbitrationTimeLimit,
		TxSizeLimit:          txSizeLimit,
		BaseDenom:            baseDenom,
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMaxRequestTimeout, &p.MaxRequestTimeout, validateMaxRequestTimeout),
		params.NewParamSetPair(KeyMinDepositMultiple, &p.MinDepositMultiple, validateMinDepositMultiple),
		params.NewParamSetPair(KeyMinDeposit, &p.MinDeposit, validateMinDeposit),
		params.NewParamSetPair(KeyServiceFeeTax, &p.ServiceFeeTax, validateServiceFeeTax),
		params.NewParamSetPair(KeySlashFraction, &p.SlashFraction, validateSlashFraction),
		params.NewParamSetPair(KeyComplaintRetrospect, &p.ComplaintRetrospect, validateComplaintRetrospect),
		params.NewParamSetPair(KeyArbitrationTimeLimit, &p.ArbitrationTimeLimit, validateArbitrationTimeLimit),
		params.NewParamSetPair(KeyTxSizeLimit, &p.TxSizeLimit, validateTxSizeLimit),
		params.NewParamSetPair(KeyBaseDenom, &p.BaseDenom, validateTxBaseDenom),
	}
}

// Equal returns a boolean determining if two Param types are identical.
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxRequestTimeout,
		DefaultMinDepositMultiple,
		DefaultMinDeposit,
		DefaultServiceFeeTax,
		DefaultSlashFraction,
		DefaultComplaintRetrospect,
		DefaultArbitrationTimeLimit,
		DefaultTxSizeLimit,
		DefaultBaseDenom,
	)
}

// String implements stringer
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Max Request Timeout:     %d
  Min Deposit Multiple:    %d
  Min Deposit:             %s
  Service Fee Tax:         %s
  Slash Fraction:          %s
  Complaint Retrospect:    %s
  Arbitration Time Limit:  %s
  Tx Size Limit:           %d
  Base Denom:              %s`,
		p.MaxRequestTimeout, p.MinDepositMultiple, p.MinDeposit.String(), p.ServiceFeeTax.String(), p.SlashFraction.String(),
		p.ComplaintRetrospect, p.ArbitrationTimeLimit, p.TxSizeLimit, p.BaseDenom)
}

// MustUnmarshalParams unmarshals the current service params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}

	return params
}

// UnmarshalParams unmarshals the current service params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	return
}

// Validate validates a set of params
func (p Params) Validate() error {
	if err := validateMaxRequestTimeout(p.MaxRequestTimeout); err != nil {
		return err
	}
	if err := validateMinDepositMultiple(p.MinDepositMultiple); err != nil {
		return err
	}
	if err := validateMinDeposit(p.MinDeposit); err != nil {
		return err
	}
	if err := validateSlashFraction(p.SlashFraction); err != nil {
		return err
	}
	if err := validateServiceFeeTax(p.ServiceFeeTax); err != nil {
		return err
	}
	if err := validateComplaintRetrospect(p.ComplaintRetrospect); err != nil {
		return err
	}
	if err := validateArbitrationTimeLimit(p.ArbitrationTimeLimit); err != nil {
		return err
	}
	if err := sdk.ValidateDenom(p.BaseDenom); err != nil {
		return err
	}

	return validateTxSizeLimit(p.TxSizeLimit)
}

func validateMaxRequestTimeout(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("maximum request timeout must be positive: %d", v)
	}

	return nil
}

func validateMinDepositMultiple(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("minimum deposit multiple must be positive: %d", v)
	}

	return nil
}

func validateMinDeposit(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsValid() {
		return fmt.Errorf("invalid minimum deposit: %s", v)
	}

	return nil
}

func validateSlashFraction(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroDec()) || v.GT(sdk.OneDec()) {
		return fmt.Errorf("slashing fraction must be between [0, 1]: %s", v)
	}

	return nil
}

func validateServiceFeeTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroDec()) || v.GTE(sdk.OneDec()) {
		return fmt.Errorf("service fee tax must be between [0, 1): %s", v)
	}

	return nil
}

func validateComplaintRetrospect(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("complaint retrospect must be positive: %d", v)
	}

	return nil
}

func validateArbitrationTimeLimit(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("arbitration time limit must be positive: %d", v)
	}

	return nil
}

func validateTxSizeLimit(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("tx size limit must be positive: %d", v)
	}

	return nil
}

func validateTxBaseDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return sdk.ValidateDenom(v)
}
