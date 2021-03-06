package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/irismod/service/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/service/definitions", defineServiceHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/service/bindings", bindServiceHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}", RestServiceName, RestProvider), updateServiceBindingHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/service/providers/{%s}/withdraw-address", RestProvider), setWithdrawAddrHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/disable", RestServiceName, RestProvider), disableServiceBindingHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/enable", RestServiceName, RestProvider), enableServiceBindingHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/service/bindings/{%s}/{%s}/refund-deposit", RestServiceName, RestProvider), refundServiceDepositHandlerFn(cliCtx)).Methods("POST")
	// initiate a service call
	r.HandleFunc("/service/contexts", requestServiceHandlerFn(cliCtx)).Methods("POST")
	// respond to a service request
	r.HandleFunc("/service/responses", respondServiceHandlerFn(cliCtx)).Methods("POST")
	// pause a request context
	r.HandleFunc(fmt.Sprintf("/service/contexts/{%s}/pause", RestRequestContextID), pauseRequestContextHandlerFn(cliCtx)).Methods("POST")
	// start a paused request context
	r.HandleFunc(fmt.Sprintf("/service/contexts/{%s}/start", RestRequestContextID), startRequestContextHandlerFn(cliCtx)).Methods("POST")
	// kill a request context
	r.HandleFunc(fmt.Sprintf("/service/contexts/{%s}/kill", RestRequestContextID), killRequestContextHandlerFn(cliCtx)).Methods("POST")
	// update a request context
	r.HandleFunc(fmt.Sprintf("/service/contexts/{%s}", RestRequestContextID), updateRequestContextHandlerFn(cliCtx)).Methods("PUT")
	// withdraw the earned fees of a provider
	r.HandleFunc(fmt.Sprintf("/service/fees/{%s}/withdraw", RestProvider), withdrawEarnedFeesHandlerFn(cliCtx)).Methods("POST")
}

// DefineServiceReq defines the properties of a define service request's body.
type DefineServiceReq struct {
	BaseReq           rest.BaseReq `json:"base_req" yaml:"base_req"`
	Name              string       `json:"name" yaml:"name"`
	Description       string       `json:"description" yaml:"description"`
	Tags              []string     `json:"tags" yaml:"tags"`
	Author            string       `json:"author" yaml:"author"`
	AuthorDescription string       `json:"author_description" yaml:"author_description"`
	Schemas           string       `json:"schemas" yaml:"schemas"`
}

// BindServiceReq defines the properties of a bind service request's body.
type BindServiceReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	ServiceName string       `json:"service_name" yaml:"service_name"`
	Provider    string       `json:"provider" yaml:"provider"`
	Deposit     string       `json:"deposit" yaml:"deposit"`
	Pricing     string       `json:"pricing" yaml:"pricing"`
	MinRespTime uint64       `json:"min_resp_time" yaml:"min_resp_time"`
}

// UpdateServiceBindingReq defines the properties of an update service binding request's body.
type UpdateServiceBindingReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Deposit     string       `json:"deposit" yaml:"deposit"`
	Pricing     string       `json:"pricing" yaml:"pricing"`
	MinRespTime uint64       `json:"min_resp_time" yaml:"min_resp_time"`
}

// SetWithdrawAddrReq defines the properties of a set withdraw address request's body.
type SetWithdrawAddrReq struct {
	BaseReq         rest.BaseReq `json:"base_req" yaml:"base_req"`
	WithdrawAddress string       `json:"withdraw_address" yaml:"withdraw_address"`
}

// DisableServiceBindingReq defines the properties of a disable service binding request's body.
type DisableServiceBindingReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
}

// EnableServiceBindingReq defines the properties of an enable service binding request's body.
type EnableServiceBindingReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Deposit string       `json:"deposit" yaml:"deposit"`
}

// RefundServiceDepositReq defines the properties of a refund service deposit request's body.
type RefundServiceDepositReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
}

type callServiceReq struct {
	BaseReq           rest.BaseReq `json:"base_req"` // basic tx info
	ServiceName       string       `json:"service_name"`
	Providers         []string     `json:"providers"`
	Consumer          string       `json:"consumer"`
	Input             string       `json:"input"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	Timeout           int64        `json:"timeout"`
	SuperMode         bool         `json:"super_mode"`
	Repeated          bool         `json:"repeated"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
}

type respondServiceReq struct {
	BaseReq   rest.BaseReq `json:"base_req"` // basic tx info
	RequestID string       `json:"request_id"`
	Provider  string       `json:"provider"`
	Result    string       `json:"result"`
	Output    string       `json:"output"`
}

type pauseRequestContextReq struct {
	BaseReq  rest.BaseReq `json:"base_req"` // basic tx info
	Consumer string       `json:"consumer"`
}

type startRequestContextReq struct {
	BaseReq  rest.BaseReq `json:"base_req"` // basic tx info
	Consumer string       `json:"consumer"`
}

type killRequestContextReq struct {
	BaseReq  rest.BaseReq `json:"base_req"` // basic tx info
	Consumer string       `json:"consumer"`
}

type updateRequestContextReq struct {
	BaseReq           rest.BaseReq `json:"base_req"` // basic tx info
	Providers         []string     `json:"providers"`
	ServiceFeeCap     string       `json:"service_fee_cap"`
	Timeout           int64        `json:"timeout"`
	RepeatedFrequency uint64       `json:"repeated_frequency"`
	RepeatedTotal     int64        `json:"repeated_total"`
	Consumer          string       `json:"consumer"`
}

type withdrawEarnedFeesReq struct {
	BaseReq rest.BaseReq `json:"base_req"` // basic tx info
}

func defineServiceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DefineServiceReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		author, err := sdk.AccAddressFromBech32(req.Author)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgDefineService(req.Name, req.Description, req.Tags, author, req.AuthorDescription, req.Schemas)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func bindServiceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BindServiceReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		provider, err := sdk.AccAddressFromBech32(req.Provider)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		deposit, err := sdk.ParseCoins(req.Deposit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgBindService(req.ServiceName, provider, deposit, req.Pricing, req.MinRespTime)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func updateServiceBindingHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req UpdateServiceBindingReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var deposit sdk.Coins
		if req.Deposit != "" {
			deposit, err = sdk.ParseCoins(req.Deposit)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		msg := types.NewMsgUpdateServiceBinding(serviceName, provider, deposit, req.Pricing, req.MinRespTime)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func setWithdrawAddrHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req SetWithdrawAddrReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		withdrawAddr, err := sdk.AccAddressFromBech32(req.WithdrawAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSetWithdrawAddress(provider, withdrawAddr)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func disableServiceBindingHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req DisableServiceBindingReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgDisableServiceBinding(serviceName, provider)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func enableServiceBindingHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req EnableServiceBindingReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var deposit sdk.Coins
		if len(req.Deposit) != 0 {
			deposit, err = sdk.ParseCoins(req.Deposit)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		msg := types.NewMsgEnableServiceBinding(serviceName, provider, deposit)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func refundServiceDepositHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		serviceName := vars[RestServiceName]
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req RefundServiceDepositReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgRefundServiceDeposit(serviceName, provider)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func requestServiceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req callServiceReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		consumer, err := sdk.AccAddressFromBech32(req.Consumer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		serviceFeeCap, err := sdk.ParseCoins(req.ServiceFeeCap)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var providers []sdk.AccAddress
		for _, p := range req.Providers {
			provider, err := sdk.AccAddressFromBech32(p)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			providers = append(providers, provider)
		}

		msg := types.NewMsgCallService(
			req.ServiceName, providers, consumer, req.Input, serviceFeeCap,
			req.Timeout, req.SuperMode, req.Repeated, req.RepeatedFrequency, req.RepeatedTotal,
		)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func respondServiceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req respondServiceReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		requestIDStr := req.RequestID
		requestID, err := types.ConvertRequestID(requestIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		provider, err := sdk.AccAddressFromBech32(req.Provider)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgRespondService(requestID, provider, req.Result, req.Output)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func pauseRequestContextHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RestRequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req pauseRequestContextReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		consumer, err := sdk.AccAddressFromBech32(req.Consumer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgPauseRequestContext(requestContextID, consumer)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func startRequestContextHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RestRequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req startRequestContextReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		consumer, err := sdk.AccAddressFromBech32(req.Consumer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgStartRequestContext(requestContextID, consumer)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func killRequestContextHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RestRequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req killRequestContextReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		consumer, err := sdk.AccAddressFromBech32(req.Consumer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgKillRequestContext(requestContextID, consumer)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func updateRequestContextHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestContextIDStr := vars[RestRequestContextID]

		requestContextID, err := hex.DecodeString(requestContextIDStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req updateRequestContextReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		consumer, err := sdk.AccAddressFromBech32(req.Consumer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var serviceFeeCap sdk.Coins

		if len(req.ServiceFeeCap) != 0 {
			serviceFeeCap, err = sdk.ParseCoins(req.ServiceFeeCap)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		var providers []sdk.AccAddress
		for _, p := range req.Providers {
			provider, err := sdk.AccAddressFromBech32(p)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			providers = append(providers, provider)
		}

		msg := types.NewMsgUpdateRequestContext(
			requestContextID, providers, serviceFeeCap, req.Timeout,
			req.RepeatedFrequency, req.RepeatedTotal, consumer,
		)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func withdrawEarnedFeesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		providerStr := vars[RestProvider]

		provider, err := sdk.AccAddressFromBech32(providerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req withdrawEarnedFeesReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgWithdrawEarnedFees(provider)
		if err = msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
