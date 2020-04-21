package service

// nolint

import (
	"github.com/irismod/service/keeper"
	"github.com/irismod/service/types"
)

const (
	ModuleName             = types.ModuleName
	StoreKey               = types.StoreKey
	QuerierRoute           = types.QuerierRoute
	RouterKey              = types.RouterKey
	DefaultParamspace      = types.DefaultParamspace
	DepositAccName         = types.DepositAccName
	RequestAccName         = types.RequestAccName
	TaxAccName             = types.TaxAccName
	QueryDefinition        = types.QueryDefinition
	QueryBinding           = types.QueryBinding
	QueryBindings          = types.QueryBindings
	QueryWithdrawAddress   = types.QueryWithdrawAddress
	EventTypeDefineService = types.EventTypeDefineService
	AttributeKeyAuthor     = types.AttributeKeyAuthor
	AttributeValueCategory = types.AttributeValueCategory
)

var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
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
	QueryDefinitionParams      = types.QueryDefinitionParams
	QueryBindingParams         = types.QueryBindingParams
	QueryBindingsParams        = types.QueryBindingsParams
	QueryWithdrawAddressParams = types.QueryWithdrawAddressParams
	TokenI                     = types.TokenI
	MockToken                  = types.MockToken
	MockTokenKeeper            = keeper.MockTokenKeeper
)
