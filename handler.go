package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the service type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgDefineService:
			return handleMsgDefineService(ctx, keeper, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgDefineService(ctx sdk.Context, k Keeper, msg MsgDefineService) (*sdk.Result, error) {
	err := k.AddServiceDefinition(ctx, msg.Name, msg.Description, msg.Tags, msg.Author, msg.AuthorDescription, msg.Schemas)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeDefineService,
			sdk.NewAttribute(AttributeKeyAuthor, msg.Author.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Author.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
