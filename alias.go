package service

// nolint

import (
	"github.com/irismod/service/keeper"
	"github.com/irismod/service/types"
)

const (
	ModuleName                   = types.ModuleName
	StoreKey                     = types.StoreKey
	QuerierRoute                 = types.QuerierRoute
	RouterKey                    = types.RouterKey
	DefaultParamspace            = types.DefaultParamspace
	DepositAccName               = types.DepositAccName
	RequestAccName               = types.RequestAccName
	QueryDefinition              = types.QueryDefinition
	QueryBinding                 = types.QueryBinding
	QueryBindings                = types.QueryBindings
	QueryWithdrawAddress         = types.QueryWithdrawAddress
	EventTypeDefineService       = types.EventTypeDefineService
	EventTypeCreateContext       = types.EventTypeCreateContext
	EventTypePauseContext        = types.EventTypePauseContext
	EventTypeCompleteContext     = types.EventTypeCompleteContext
	EventTypeNewBatch            = types.EventTypeNewBatch
	EventTypeNewBatchRequest     = types.EventTypeNewBatchRequest
	EventTypeCompleteBatch       = types.EventTypeCompleteBatch
	AttributeValueCategory       = types.AttributeValueCategory
	AttributeKeyAuthor           = types.AttributeKeyAuthor
	AttributeKeyServiceName      = types.AttributeKeyServiceName
	AttributeKeyProvider         = types.AttributeKeyProvider
	AttributeKeyConsumer         = types.AttributeKeyConsumer
	AttributeKeyRequestContextID = types.AttributeKeyRequestContextID
	AttributeKeyRequestID        = types.AttributeKeyRequestID
	AttributeKeyServiceFee       = types.AttributeKeyServiceFee
	AttributeKeyRequestHeight    = types.AttributeKeyRequestHeight
	AttributeKeyExpirationHeight = types.AttributeKeyExpirationHeight
	AttributeKeySlashedCoins     = types.AttributeKeySlashedCoins

	RUNNING        = types.RUNNING
	PAUSED         = types.PAUSED
	COMPLETED      = types.COMPLETED
	BATCHRUNNING   = types.BATCHRUNNING
	BATCHCOMPLETED = types.BATCHCOMPLETED
)

var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewGenesisState     = types.NewGenesisState
)

type (
	Keeper                     = keeper.Keeper
	ServiceDefinition          = types.ServiceDefinition
	ServiceBinding             = types.ServiceBinding
	GenesisState               = types.GenesisState
	MsgDefineService           = types.MsgDefineService
	MsgBindService             = types.MsgBindService
	MsgUpdateServiceBinding    = types.MsgUpdateServiceBinding
	MsgSetWithdrawAddress      = types.MsgSetWithdrawAddress
	MsgDisableServiceBinding   = types.MsgDisableServiceBinding
	MsgEnableServiceBinding    = types.MsgEnableServiceBinding
	MsgRefundServiceDeposit    = types.MsgRefundServiceDeposit
	MsgCallService             = types.MsgCallService
	MsgRespondService          = types.MsgRespondService
	MsgPauseRequestContext     = types.MsgPauseRequestContext
	MsgStartRequestContext     = types.MsgStartRequestContext
	MsgKillRequestContext      = types.MsgKillRequestContext
	MsgUpdateRequestContext    = types.MsgUpdateRequestContext
	MsgWithdrawEarnedFees      = types.MsgWithdrawEarnedFees
	QueryDefinitionParams      = types.QueryDefinitionParams
	QueryBindingParams         = types.QueryBindingParams
	QueryBindingsParams        = types.QueryBindingsParams
	QueryWithdrawAddressParams = types.QueryWithdrawAddressParams
	TokenI                     = types.TokenI
	MockToken                  = types.MockToken
	MockTokenKeeper            = keeper.MockTokenKeeper
	Request                    = types.Request
	Response                   = types.Response
	RequestContext             = types.RequestContext
	EarnedFeesOutput           = types.EarnedFeesOutput
)
