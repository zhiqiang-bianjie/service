package types

// TokenI defines the interface for Token
type TokenI interface {
	GetMinUnit() string
	GetScale() uint8
}

// MockToken represents a mock implementation for TokenI
type MockToken struct {
	Symbol  string
	MinUnit string
	Scale   uint8
}

// GetSymbol gets the symbol
func (token MockToken) GetSymbol() string {
	return token.Symbol
}

// GetMinUnit gets the min unit
func (token MockToken) GetMinUnit() string {
	return token.MinUnit
}

// GetScale gets the scale
func (token MockToken) GetScale() uint8 {
	return token.Scale
}
