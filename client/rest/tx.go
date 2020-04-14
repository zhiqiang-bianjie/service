package rest

import (
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
