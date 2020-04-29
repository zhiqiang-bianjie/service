package types

// service module event types
const (
	EventTypeDefineService           = "define_service"
	EventTypeCreateContext           = "create-context"
	EventTypePauseContext            = "pause-context"
	EventTypeCompleteContext         = "complete-context"
	EventTypeNewBatch                = "new-batch"
	EventTypeNewBatchRequest         = "new-batch-request"
	EventTypeNewBatchRequestProvider = "new-batch-request-provider"
	EventTypeCompleteBatch           = "complete-batch"
	EventTypeServiceSlash            = "service-slash"

	AttributeValueCategory          = ModuleName
	AttributeKeyAuthor              = "author"
	AttributeKeyServiceName         = "service-name"
	AttributeKeyProvider            = "provider"
	AttributeKeyConsumer            = "consumer"
	AttributeKeyRequestContextID    = "request-context-id"
	AttributeKeyRequestContextState = "request-context-state"
	AttributeKeyRequests            = "requests"
	AttributeKeyRequestID           = "request-id"
	AttributeKeyServiceFee          = "service-fee"
	AttributeKeyRequestHeight       = "request-height"
	AttributeKeyExpirationHeight    = "expiration-height"
	AttributeKeySlashedCoins        = "slashed-coins"
)

type BatchState struct {
	BatchCounter           uint64                   `json:"batch_counter"`
	State                  RequestContextBatchState `json:"state"`
	BatchResponseThreshold uint16                   `json:"batch_response_threshold"`
	BatchRequestCount      uint16                   `json:"batch_request_count"`
	BatchResponseCount     uint16                   `json:"batch_response_count"`
}

// ActionTag appends action and all tagKeys
func ActionTag(action string, tagKeys ...string) string {
	tag := action
	for _, key := range tagKeys {
		tag = tag + "." + key
	}
	return tag
}
