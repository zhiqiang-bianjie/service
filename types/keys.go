package types

const (
	// ModuleName is the name of the service module
	ModuleName = "service"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the service module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the service module
	RouterKey string = ModuleName
)

var (
	// Keys for store prefixes

	ServiceDefinitionKey = []byte{0x01} // prefix for service definition
)

// GetServiceDefinitionKey gets the key for the service definition with the specified service name
// VALUE: service/ServiceDefinition
func GetServiceDefinitionKey(serviceName string) []byte {
	return append(ServiceDefinitionKey, []byte(serviceName)...)
}
