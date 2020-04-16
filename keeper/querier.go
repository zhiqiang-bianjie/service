package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irismod/service/types"
)

// NewQuerier creates a new service Querier instance
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryDefinition:
			return queryServiceDefinition(ctx, path[1:], req, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query path: %s", types.ModuleName, path[0])
		}
	}
}

func queryServiceDefinition(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDefinitionParams
	if err := k.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	definition, found := k.GetServiceDefinition(ctx, params.ServiceName)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUnknownServiceDefinition, params.ServiceName)
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, definition)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
