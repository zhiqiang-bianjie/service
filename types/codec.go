package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgDefineService{}, "irismod/service/MsgDefineService", nil)
	cdc.RegisterConcrete(MsgBindService{}, "irismod/service/MsgBindService", nil)
	cdc.RegisterConcrete(MsgUpdateServiceBinding{}, "irismod/service/MsgUpdateServiceBinding", nil)
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "irismod/service/MsgSetWithdrawAddress", nil)
	cdc.RegisterConcrete(MsgDisableServiceBinding{}, "irismod/service/MsgDisableServiceBinding", nil)
	cdc.RegisterConcrete(MsgEnableServiceBinding{}, "irismod/service/MsgEnableServiceBinding", nil)
	cdc.RegisterConcrete(MsgRefundServiceDeposit{}, "irismod/service/MsgRefundServiceDeposit", nil)

	cdc.RegisterConcrete(MsgCallService{}, "irismod/service/MsgCallService", nil)
	cdc.RegisterConcrete(MsgRespondService{}, "irismod/service/MsgRespondService", nil)
	cdc.RegisterConcrete(MsgPauseRequestContext{}, "irismod/service/MsgPauseRequestContext", nil)
	cdc.RegisterConcrete(MsgStartRequestContext{}, "irismod/service/MsgStartRequestContext", nil)
	cdc.RegisterConcrete(MsgKillRequestContext{}, "irismod/service/MsgKillRequestContext", nil)
	cdc.RegisterConcrete(MsgUpdateRequestContext{}, "irismod/service/MsgUpdateRequestContext", nil)
	cdc.RegisterConcrete(MsgWithdrawEarnedFees{}, "irismod/service/MsgWithdrawEarnedFees", nil)

	cdc.RegisterConcrete(ServiceDefinition{}, "irismod/service/ServiceDefinition", nil)
	cdc.RegisterConcrete(ServiceBinding{}, "irismod/service/ServiceBinding", nil)
	cdc.RegisterConcrete(RequestContext{}, "irismod/service/RequestContext", nil)
	cdc.RegisterConcrete(CompactRequest{}, "irismod/service/CompactRequest", nil)
	cdc.RegisterConcrete(Request{}, "irismod/service/Request", nil)
	cdc.RegisterConcrete(Response{}, "irismod/service/Response", nil)
	cdc.RegisterConcrete(EarnedFeesOutput{}, "irismod/service/EarnedFeesOutput", nil)

	cdc.RegisterConcrete(&Params{}, "irismod/service/Params", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
