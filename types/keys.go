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

	// ServiceDepositCoinDenom is the coin denom for service deposit
	ServiceDepositCoinDenom = sdk.DefaultBondDenom

	// ServiceDepositCoinDecimal is the coin decimal for service deposit
	ServiceDepositCoinDecimal = 6

	TxHash   = "tx_hash"
	MsgIndex = "msg_index"
)

var (
	// Separator for string key
	emptyByte = []byte{0x00}

	// Keys for store prefixes
	ServiceDefinitionKey         = []byte{0x01} // prefix for service definition
	ServiceBindingKey            = []byte{0x02} // prefix for service binding
	PricingKey                   = []byte{0x03} // prefix for pricing
	WithdrawAddrKey              = []byte{0x04} // prefix for withdrawal address
	RequestContextKey            = []byte{0x05}
	ExpiredRequestBatchKey       = []byte{0x06}
	NewRequestBatchKey           = []byte{0x07}
	ExpiredRequestBatchHeightKey = []byte{0x08}
	NewRequestBatchHeightKey     = []byte{0x09}
	RequestKey                   = []byte{0x10}
	ActiveRequestKey             = []byte{0x11}
	ActiveRequestByIDKey         = []byte{0x12}
	ResponseKey                  = []byte{0x13}
	RequestVolumeKey             = []byte{0x14}
	EarnedFeesKey                = []byte{0x15}
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

// GetRequestContextKey returns the key for the request context with the specified ID
func GetRequestContextKey(requestContextID []byte) []byte {
	return append(RequestContextKey, requestContextID...)
}

// GetExpiredRequestBatchKey returns the key for the request batch expiration of the specified request context
func GetExpiredRequestBatchKey(requestContextID []byte, batchExpirationHeight int64) []byte {
	reqBatchExpiration := append(sdk.Uint64ToBigEndian(uint64(batchExpirationHeight)), requestContextID...)
	return append(ExpiredRequestBatchKey, reqBatchExpiration...)
}

// GetNewRequestBatchKey returns the key for the new batch request of the specified request context in the given height
func GetNewRequestBatchKey(requestContextID []byte, requestBatchHeight int64) []byte {
	newBatchRequest := append(sdk.Uint64ToBigEndian(uint64(requestBatchHeight)), requestContextID...)
	return append(NewRequestBatchKey, newBatchRequest...)
}

// GetExpiredRequestBatchSubspace returns the key for iterating through the expired request batch queue in the specified height
func GetExpiredRequestBatchSubspace(batchExpirationHeight int64) []byte {
	return append(ExpiredRequestBatchKey, sdk.Uint64ToBigEndian(uint64(batchExpirationHeight))...)
}

// GetNewRequestBatchSubspace returns the key for iterating through the new request batch queue in the specified height
func GetNewRequestBatchSubspace(requestBatchHeight int64) []byte {
	return append(NewRequestBatchKey, sdk.Uint64ToBigEndian(uint64(requestBatchHeight))...)
}

// GetExpiredRequestBatchHeightKey returns the key for the current request batch expiration height of the specified request context
func GetExpiredRequestBatchHeightKey(requestContextID []byte) []byte {
	return append(ExpiredRequestBatchHeightKey, requestContextID...)
}

// GetNewRequestBatchHeightKey returns the key for the new request batch height of the specified request context
func GetNewRequestBatchHeightKey(requestContextID []byte) []byte {
	return append(NewRequestBatchHeightKey, requestContextID...)
}

// GetRequestKey returns the key for the request with the specified request ID
func GetRequestKey(requestID []byte) []byte {
	return append(RequestKey, requestID...)
}

// GetRequestSubspaceByReqCtx returns the key for the requests of the specified request context
func GetRequestSubspaceByReqCtx(requestContextID []byte, batchCounter uint64) []byte {
	return append(append(RequestKey, requestContextID...), sdk.Uint64ToBigEndian(batchCounter)...)
}

// GetActiveRequestKey returns the key for the active request with the specified request ID in the given height
func GetActiveRequestKey(serviceName string, provider sdk.AccAddress, expirationHeight int64, requestID []byte) []byte {
	activeRequest := append(append(append(getStringsKey([]string{serviceName, provider.String()}), emptyByte...), sdk.Uint64ToBigEndian(uint64(expirationHeight))...), requestID...)
	return append(ActiveRequestKey, activeRequest...)
}

// GetActiveRequestSubspace returns the key for the active requests for the specified provider
func GetActiveRequestSubspace(serviceName string, provider sdk.AccAddress) []byte {
	return append(append(ActiveRequestKey, getStringsKey([]string{serviceName, provider.String()})...), emptyByte...)
}

// GetActiveRequestKeyByID returns the key for the active request with the specified request ID
func GetActiveRequestKeyByID(requestID []byte) []byte {
	return append(ActiveRequestByIDKey, requestID...)
}

// GetActiveRequestSubspaceByReqCtx returns the key for the active requests for the specified request context
func GetActiveRequestSubspaceByReqCtx(requestContextID []byte, batchCounter uint64) []byte {
	return append(append(ActiveRequestByIDKey, requestContextID...), sdk.Uint64ToBigEndian(batchCounter)...)
}

// GetRequestVolumeKey returns the key for the request volume for the specified consumer and binding
func GetRequestVolumeKey(consumer sdk.AccAddress, serviceName string, provider sdk.AccAddress) []byte {
	return append(append(RequestVolumeKey, getStringsKey([]string{consumer.String(), serviceName, provider.String()})...), emptyByte...)
}

// GetResponseKey returns the key for the response for the given request ID
func GetResponseKey(requestID []byte) []byte {
	return append(ResponseKey, requestID...)
}

// GetResponseSubspaceByReqCtx returns the key for responses for the specified request context and batch counter
func GetResponseSubspaceByReqCtx(requestContextID []byte, batchCounter uint64) []byte {
	return append(append(ResponseKey, requestContextID...), sdk.Uint64ToBigEndian(batchCounter)...)
}

// GetEarnedFeesKey returns the key for the earned fees of the specified provider
func GetEarnedFeesKey(provider sdk.AccAddress) []byte {
	return append(EarnedFeesKey, provider.Bytes()...)
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
