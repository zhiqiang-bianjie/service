package simapp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irismod/service"
)

type MockToken struct {
	Symbol  string
	MinUnit string
	Scale   uint8
}

func (token MockToken) GetSymbol() string {
	return token.Symbol
}

func (token MockToken) GetName() string {
	return ""
}

func (token MockToken) GetMinUnit() string {
	return token.MinUnit
}

func (token MockToken) GetScale() uint8 {
	return token.Scale
}

func (token MockToken) GetMaxSupply() uint64 {
	return 0
}

func (token MockToken) GetInitialSupply() uint64 {
	return 0
}

func (token MockToken) GetMintable() bool {
	return true
}

func (token MockToken) GetOwner() sdk.AccAddress {
	return nil
}

func (token MockToken) ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error) {
	return sdk.DecCoin{}, nil
}

func (token MockToken) ToMinCoin(coin sdk.DecCoin) (sdk.Coin, error) {
	return sdk.Coin{}, nil
}

type MockTokenKeeper struct{}

func (token MockTokenKeeper) GetToken(ctx sdk.Context, denom string) (service.TokenI, error) {
	if denom == sdk.DefaultBondDenom {
		return MockToken{
			MinUnit: sdk.DefaultBondDenom,
			Scale:   6,
		}, nil
	}

	return nil, fmt.Errorf("token %s does not exist", denom)
}
