package types

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message types for the service module
const (
	TypeMsgDefineService = "define_service" // type for MsgDefineService

	MaxNameLength        = 70  // maximum length of the service name
	MaxDescriptionLength = 280 // maximum length of the service and author description
	MaxTagsNum           = 10  // maximum total number of the tags
	MaxTagLength         = 70  // maximum length of the tag
)

// the service name only accepts alphanumeric characters, _ and -, beginning with alpha character
var reServiceName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

var (
	_ sdk.Msg = MsgDefineService{}
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
