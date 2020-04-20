package types

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message types for the service module
const (
	TypeMsgDefineService         = "define_service"          // type for MsgDefineService
	TypeMsgBindService           = "bind_service"            // type for MsgBindService
	TypeMsgUpdateServiceBinding  = "update_service_binding"  // type for MsgUpdateServiceBinding
	TypeMsgSetWithdrawAddress    = "set_withdraw_address"    // type for MsgSetWithdrawAddress
	TypeMsgDisableServiceBinding = "disable_service_binding" // type for MsgDisableServiceBinding
	TypeMsgEnableServiceBinding  = "enable_service_binding"  // type for MsgEnableServiceBinding
	TypeMsgRefundServiceDeposit  = "refund_service_deposit"  // type for MsgRefundServiceDeposit

	MaxNameLength        = 70  // maximum length of the service name
	MaxDescriptionLength = 280 // maximum length of the service and author description
	MaxTagsNum           = 10  // maximum total number of the tags
	MaxTagLength         = 70  // maximum length of the tag
)

// the service name only accepts alphanumeric characters, _ and -, beginning with alpha character
var reServiceName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

var (
	_ sdk.Msg = MsgDefineService{}
	_ sdk.Msg = MsgBindService{}
	_ sdk.Msg = MsgUpdateServiceBinding{}
	_ sdk.Msg = MsgSetWithdrawAddress{}
	_ sdk.Msg = MsgDisableServiceBinding{}
	_ sdk.Msg = MsgEnableServiceBinding{}
	_ sdk.Msg = MsgRefundServiceDeposit{}
)

//______________________________________________________________________

// MsgDefineService defines a message to define a service
type MsgDefineService struct {
	Name              string         `json:"name" yaml:"name"`
	Description       string         `json:"description" yaml:"description"`
	Tags              []string       `json:"tags" yaml:"tags"`
	Author            sdk.AccAddress `json:"author" yaml:"author"`
	AuthorDescription string         `json:"author_description" yaml:"author_description"`
	Schemas           string         `json:"schemas" yaml:"schemas"`
}

// NewMsgDefineService creates a new MsgDefineService instance
func NewMsgDefineService(
	name,
	description string,
	tags []string,
	author sdk.AccAddress,
	authorDescription,
	schemas string,
) MsgDefineService {
	return MsgDefineService{
		Name:              name,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		Schemas:           schemas,
	}
}

// Route implements Msg
func (msg MsgDefineService) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgDefineService) Type() string { return TypeMsgDefineService }

// ValidateBasic implements Msg
func (msg MsgDefineService) ValidateBasic() error {
	if err := ValidateAuthor(msg.Author); err != nil {
		return err
	}

	if err := ValidateServiceName(msg.Name); err != nil {
		return err
	}

	if err := ValidateServiceDescription(msg.Description); err != nil {
		return err
	}

	if err := ValidateAuthorDescription(msg.AuthorDescription); err != nil {
		return err
	}

	if err := ValidateTags(msg.Tags); err != nil {
		return err
	}

	if err := ValidateServiceSchemas(msg.Schemas); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg
func (msg MsgDefineService) GetSignBytes() []byte {
	if len(msg.Tags) == 0 {
		msg.Tags = nil
	}

	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgDefineService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}

//______________________________________________________________________

// MsgBindService defines a message to bind a service
type MsgBindService struct {
	ServiceName string         `json:"service_name" yaml:"service_name"`
	Provider    sdk.AccAddress `json:"provider" yaml:"provider"`
	Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	Pricing     string         `json:"pricing" yaml:"pricing"`
	MinRespTime uint64         `json:"min_resp_time" yaml:"min_resp_time"`
}

// NewMsgBindService creates a new MsgBindService instance
func NewMsgBindService(serviceName string, provider sdk.AccAddress, deposit sdk.Coins, pricing string, minRespTime uint64) MsgBindService {
	return MsgBindService{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
		Pricing:     pricing,
		MinRespTime: minRespTime,
	}
}

// Route implements Msg.
func (msg MsgBindService) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgBindService) Type() string { return TypeMsgBindService }

// GetSignBytes implements Msg.
func (msg MsgBindService) GetSignBytes() []byte {
	if msg.Deposit.Empty() {
		msg.Deposit = nil
	}

	b := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgBindService) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if err := ValidateServiceDeposit(msg.Deposit); err != nil {
		return err
	}

	if err := ValidateMinRespTime(msg.MinRespTime); err != nil {
		return err
	}

	return ValidateBindingPricing(msg.Pricing)
}

// GetSigners implements Msg.
func (msg MsgBindService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgUpdateServiceBinding defines a message to update a service binding
type MsgUpdateServiceBinding struct {
	ServiceName string         `json:"service_name" yaml:"service_name"`
	Provider    sdk.AccAddress `json:"provider" yaml:"provider"`
	Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	Pricing     string         `json:"pricing" yaml:"pricing"`
	MinRespTime uint64         `json:"min_resp_time" yaml:"min_resp_time"`
}

// NewMsgUpdateServiceBinding creates a new MsgUpdateServiceBinding instance
func NewMsgUpdateServiceBinding(
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
	minRespTime uint64,
) MsgUpdateServiceBinding {
	return MsgUpdateServiceBinding{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
		Pricing:     pricing,
		MinRespTime: minRespTime,
	}
}

// Route implements Msg.
func (msg MsgUpdateServiceBinding) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgUpdateServiceBinding) Type() string { return TypeMsgUpdateServiceBinding }

// GetSignBytes implements Msg.
func (msg MsgUpdateServiceBinding) GetSignBytes() []byte {
	if msg.Deposit.Empty() {
		msg.Deposit = nil
	}

	b := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateServiceBinding) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if !msg.Deposit.Empty() {
		if err := ValidateServiceDeposit(msg.Deposit); err != nil {
			return err
		}
	}

	if len(msg.Pricing) != 0 {
		return ValidateBindingPricing(msg.Pricing)
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateServiceBinding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSetWithdrawAddress defines a message to set the withdrawal address for a provider
type MsgSetWithdrawAddress struct {
	Provider        sdk.AccAddress `json:"provider" yaml:"provider"`
	WithdrawAddress sdk.AccAddress `json:"withdraw_address" yaml:"withdraw_address"`
}

// NewMsgSetWithdrawAddress creates a new MsgSetWithdrawAddress instance
func NewMsgSetWithdrawAddress(provider sdk.AccAddress, withdrawAddr sdk.AccAddress) MsgSetWithdrawAddress {
	return MsgSetWithdrawAddress{
		Provider:        provider,
		WithdrawAddress: withdrawAddr,
	}
}

// Route implements Msg.
func (msg MsgSetWithdrawAddress) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgSetWithdrawAddress) Type() string { return TypeMsgSetWithdrawAddress }

// GetSignBytes implements Msg.
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSetWithdrawAddress) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	return ValidateWithdrawAddress(msg.WithdrawAddress)
}

// GetSigners implements Msg.
func (msg MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgDisableServiceBinding defines a message to disable a service binding
type MsgDisableServiceBinding struct {
	ServiceName string         `json:"service_name" yaml:"service_name"`
	Provider    sdk.AccAddress `json:"provider" yaml:"provider"`
}

// NewMsgDisableServiceBinding creates a new MsgDisableServiceBinding instance
func NewMsgDisableServiceBinding(serviceName string, provider sdk.AccAddress) MsgDisableServiceBinding {
	return MsgDisableServiceBinding{
		ServiceName: serviceName,
		Provider:    provider,
	}
}

// Route implements Msg.
func (msg MsgDisableServiceBinding) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgDisableServiceBinding) Type() string { return TypeMsgDisableServiceBinding }

// GetSignBytes implements Msg.
func (msg MsgDisableServiceBinding) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgDisableServiceBinding) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	return ValidateServiceName(msg.ServiceName)
}

// GetSigners implements Msg.
func (msg MsgDisableServiceBinding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgEnableServiceBinding defines a message to enable a service binding
type MsgEnableServiceBinding struct {
	ServiceName string         `json:"service_name" yaml:"service_name"`
	Provider    sdk.AccAddress `json:"provider" yaml:"provider"`
	Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
}

// NewMsgEnableServiceBinding creates a new MsgEnableServiceBinding instance
func NewMsgEnableServiceBinding(serviceName string, provider sdk.AccAddress, deposit sdk.Coins) MsgEnableServiceBinding {
	return MsgEnableServiceBinding{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
	}
}

// Route implements Msg.
func (msg MsgEnableServiceBinding) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgEnableServiceBinding) Type() string { return TypeMsgEnableServiceBinding }

// GetSignBytes implements Msg.
func (msg MsgEnableServiceBinding) GetSignBytes() []byte {
	if msg.Deposit.Empty() {
		msg.Deposit = nil
	}

	b := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgEnableServiceBinding) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if !msg.Deposit.Empty() {
		if err := ValidateServiceDeposit(msg.Deposit); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgEnableServiceBinding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgRefundServiceDeposit defines a message to refund deposit from a service binding
type MsgRefundServiceDeposit struct {
	ServiceName string         `json:"service_name" yaml:"service_name"`
	Provider    sdk.AccAddress `json:"provider" yaml:"provider"`
}

// NewMsgRefundServiceDeposit creates a new MsgRefundServiceDeposit instance
func NewMsgRefundServiceDeposit(serviceName string, provider sdk.AccAddress) MsgRefundServiceDeposit {
	return MsgRefundServiceDeposit{
		ServiceName: serviceName,
		Provider:    provider,
	}
}

// Route implements Msg.
func (msg MsgRefundServiceDeposit) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgRefundServiceDeposit) Type() string { return TypeMsgRefundServiceDeposit }

// GetSignBytes implements Msg.
func (msg MsgRefundServiceDeposit) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgRefundServiceDeposit) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	return ValidateServiceName(msg.ServiceName)
}

// GetSigners implements Msg.
func (msg MsgRefundServiceDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

func ValidateAuthor(author sdk.AccAddress) error {
	if author.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "author missing")
	}

	return nil
}

// ValidateServiceName validates the service name
func ValidateServiceName(name string) error {
	if !reServiceName.MatchString(name) || len(name) > MaxNameLength {
		return sdkerrors.Wrap(ErrInvalidServiceName, name)
	}

	return nil
}

func ValidateTags(tags []string) error {
	if len(tags) > MaxTagsNum {
		return sdkerrors.Wrap(ErrInvalidTags, fmt.Sprintf("invalid tags size; got: %d, max: %d", len(tags), MaxTagsNum))
	}

	if HasDuplicate(tags) {
		return sdkerrors.Wrap(ErrInvalidTags, "duplicate tag")
	}

	for i, tag := range tags {
		if len(tag) == 0 {
			return sdkerrors.Wrap(ErrInvalidTags, fmt.Sprintf("invalid tag[%d] length: tag must not be empty", i))
		}

		if len(tag) > MaxTagLength {
			return sdkerrors.Wrap(ErrInvalidTags, fmt.Sprintf("invalid tag[%d] length; got: %d, max: %d", i, len(tag), MaxTagLength))
		}
	}

	return nil
}

func ValidateServiceDescription(svcDescription string) error {
	if len(svcDescription) > MaxDescriptionLength {
		return sdkerrors.Wrap(ErrInvalidDescription, fmt.Sprintf("invalid service description length; got: %d, max: %d", len(svcDescription), MaxDescriptionLength))
	}

	return nil
}

func ValidateAuthorDescription(authorDescription string) error {
	if len(authorDescription) > MaxDescriptionLength {
		return sdkerrors.Wrap(ErrInvalidDescription, fmt.Sprintf("invalid author description length; got: %d, max: %d", len(authorDescription), MaxDescriptionLength))
	}

	return nil
}

func ValidateProvider(provider sdk.AccAddress) error {
	if len(provider) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "provider missing")
	}

	return nil
}

func ValidateServiceDeposit(deposit sdk.Coins) error {
	if deposit.Empty() || len(deposit) != 1 || deposit[0].Denom != sdk.DefaultBondDenom || !deposit[0].IsPositive() {
		return sdkerrors.Wrap(ErrInvalidDeposit, "")
	}

	return nil
}

func ValidateMinRespTime(minRespTime uint64) error {
	if minRespTime == 0 {
		return sdkerrors.Wrap(ErrInvalidMinRespTime, "minimum response time must be greater than 0")
	}

	return nil
}

func ValidateWithdrawAddress(withdrawAddress sdk.AccAddress) error {
	if len(withdrawAddress) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "withdrawal address missing")
	}

	return nil
}
