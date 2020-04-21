package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MockTokenKeeper defines a mock implementation for types.TokenKeeper
type MockTokenKeeper struct{}

// GetToken gets the specified token
func (token MockTokenKeeper) GetToken(ctx sdk.Context, denom string) (types.TokenI, error) {
	if denom == sdk.DefaultBondDenom {
		return MockToken{
			MinUnit: sdk.DefaultBondDenom,
			Scale:   6,
		}, nil
	}

	return nil, fmt.Errorf("token %s does not exist", denom)
}
