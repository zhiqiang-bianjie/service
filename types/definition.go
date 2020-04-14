package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ServiceDefinition defines a struct for the service definition
type ServiceDefinition struct {
	Name              string         `json:"name" yaml:"name"`
	Description       string         `json:"description" yaml:"description"`
	Tags              []string       `json:"tags" yaml:"tags"`
	Author            sdk.AccAddress `json:"author" yaml:"author"`
	AuthorDescription string         `json:"author_description" yaml:"author_description"`
	Schemas           string         `json:"schemas" yaml:"schemas"`
}

// NewServiceDefinition creates a new ServiceDefinition instance
func NewServiceDefinition(
	name string,
	description string,
	tags []string,
	author sdk.AccAddress,
	authorDescription,
	schemas string,
) ServiceDefinition {
	return ServiceDefinition{
		Name:              name,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		Schemas:           schemas,
	}
}

// Validate validates the service definition
func (svcDef ServiceDefinition) Validate() error {
	if err := ValidateAuthor(svcDef.Author); err != nil {
		return err
	}

	if err := ValidateServiceName(svcDef.Name); err != nil {
		return err
	}

	if err := ValidateTags(svcDef.Tags); err != nil {
		return err
	}

	if err := ValidateServiceDescription(svcDef.Description); err != nil {
		return err
	}

	if err := ValidateAuthorDescription(svcDef.AuthorDescription); err != nil {
		return err
	}

	return ValidateServiceSchemas(svcDef.Schemas)
}
