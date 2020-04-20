package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the service type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgDefineService:
			return handleMsgDefineService(ctx, k, msg)

		case MsgBindService:
			return handleMsgBindService(ctx, k, msg)

		case MsgUpdateServiceBinding:
			return handleMsgUpdateServiceBinding(ctx, k, msg)

		case MsgSetWithdrawAddress:
			return handleMsgSetWithdrawAddress(ctx, k, msg)

		case MsgDisableServiceBinding:
			return handleMsgDisableServiceBinding(ctx, k, msg)

		case MsgEnableServiceBinding:
			return handleMsgEnableServiceBinding(ctx, k, msg)

		case MsgRefundServiceDeposit:
			return handleMsgRefundServiceDeposit(ctx, k, msg)

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

func handleMsgBindService(ctx sdk.Context, k Keeper, msg MsgBindService) (*sdk.Result, error) {
	err := k.AddServiceBinding(ctx, msg.ServiceName, msg.Provider, msg.Deposit, msg.Pricing, msg.MinRespTime)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUpdateServiceBinding(ctx sdk.Context, k Keeper, msg MsgUpdateServiceBinding) (*sdk.Result, error) {
	err := k.UpdateServiceBinding(ctx, msg.ServiceName, msg.Provider, msg.Deposit, msg.Pricing, msg.MinRespTime)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSetWithdrawAddress(ctx sdk.Context, k Keeper, msg MsgSetWithdrawAddress) (*sdk.Result, error) {
	k.SetWithdrawAddress(ctx, msg.Provider, msg.WithdrawAddress)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDisableServiceBinding(ctx sdk.Context, k Keeper, msg MsgDisableServiceBinding) (*sdk.Result, error) {
	err := k.DisableServiceBinding(ctx, msg.ServiceName, msg.Provider)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgEnableServiceBinding(ctx sdk.Context, k Keeper, msg MsgEnableServiceBinding) (*sdk.Result, error) {
	err := k.EnableServiceBinding(ctx, msg.ServiceName, msg.Provider, msg.Deposit)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRefundServiceDeposit(ctx sdk.Context, k Keeper, msg MsgRefundServiceDeposit) (*sdk.Result, error) {
	err := k.RefundDeposit(ctx, msg.ServiceName, msg.Provider)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
