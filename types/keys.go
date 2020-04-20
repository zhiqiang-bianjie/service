package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the service module
	ModuleName = "service"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the service module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the service module
	RouterKey string = ModuleName

	// DefaultParamspace is the default name for parameter store
	DefaultParamspace = ModuleName

	// DepositAccName is the root string for the service deposit account address
	DepositAccName = "service_deposit_account"

	// RequestAccName is the root string for the service request account address
	RequestAccName = "service_request_account"

	// TaxAccName is the root string for the service tax account address
	TaxAccName = "service_tax_account"
)

var (
	// Separator for string key
	emptyByte = []byte{0x00}

	// Keys for store prefixes
	ServiceDefinitionKey = []byte{0x01} // prefix for service definition
	ServiceBindingKey    = []byte{0x02} // prefix for service binding
	PricingKey           = []byte{0x03} // prefix for pricing
	WithdrawAddrKey      = []byte{0x04} // prefix for withdrawal address
)

// GetServiceDefinitionKey gets the key for the service definition with the specified service name
// VALUE: service/ServiceDefinition
func GetServiceDefinitionKey(serviceName string) []byte {
	return append(ServiceDefinitionKey, []byte(serviceName)...)
}

// GetServiceBindingKey gets the key for the service binding with the specified service name and provider
// VALUE: service/ServiceBinding
func GetServiceBindingKey(serviceName string, provider sdk.AccAddress) []byte {
	return append(ServiceBindingKey, getStringsKey([]string{serviceName, provider.String()})...)
}

// GetPricingKey gets the key for the pricing of the specified binding
// VALUE: service/Pricing
func GetPricingKey(serviceName string, provider sdk.AccAddress) []byte {
	return append(PricingKey, getStringsKey([]string{serviceName, provider.String()})...)
}

// GetWithdrawAddrKey gets the key for the withdrawal address of the specified provider
// VALUE: withdrawal address ([]byte)
func GetWithdrawAddrKey(provider sdk.AccAddress) []byte {
	return append(WithdrawAddrKey, provider.Bytes()...)
}

// GetBindingsSubspace gets the key for retrieving all bindings of the specified service
func GetBindingsSubspace(serviceName string) []byte {
	return append(append(ServiceBindingKey, []byte(serviceName)...), emptyByte...)
}

func getStringsKey(ss []string) (result []byte) {
	for _, s := range ss {
		result = append(append(result, []byte(s)...), emptyByte...)
	}

	if len(result) > 0 {
		return result[0 : len(result)-1]
	}

	return
}
