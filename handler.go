package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irismod/service/types"
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

		case MsgCallService:
			return handleMsgCallService(ctx, k, msg)

		case MsgRespondService:
			return handleMsgRespondService(ctx, k, msg)

		case MsgPauseRequestContext:
			return handleMsgPauseRequestContext(ctx, k, msg)

		case MsgStartRequestContext:
			return handleMsgStartRequestContext(ctx, k, msg)

		case MsgKillRequestContext:
			return handleMsgKillRequestContext(ctx, k, msg)

		case MsgUpdateRequestContext:
			return handleMsgUpdateRequestContext(ctx, k, msg)

		case MsgWithdrawEarnedFees:
			return handleMsgWithdrawEarnedFees(ctx, k, msg)

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

// handleMsgCallService handles MsgCallService
func handleMsgCallService(ctx sdk.Context, k Keeper, msg MsgCallService) (*sdk.Result, error) {
	reqContextID, err := k.CreateRequestContext(
		ctx, msg.ServiceName, msg.Providers, msg.Consumer, msg.Input, msg.ServiceFeeCap, msg.Timeout,
		msg.SuperMode, msg.Repeated, msg.RepeatedFrequency, msg.RepeatedTotal, RUNNING, 0, "")
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
			sdk.NewAttribute(types.AttributeKeyRequestContextID, reqContextID.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgRespondService handles MsgRespondService
func handleMsgRespondService(ctx sdk.Context, k Keeper, msg MsgRespondService) (*sdk.Result, error) {
	request, _, err := k.AddResponse(ctx, msg.RequestID, msg.Provider, msg.Result, msg.Output)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Provider.String()),
			sdk.NewAttribute(types.AttributeKeyRequestContextID, request.RequestContextID.String()),
			sdk.NewAttribute(types.AttributeKeyRequestID, msg.RequestID.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgPauseRequestContext handles MsgPauseRequestContext
func handleMsgPauseRequestContext(ctx sdk.Context, k Keeper, msg MsgPauseRequestContext) (*sdk.Result, error) {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextID, true); err != nil {
		return nil, err
	}

	if err := k.PauseRequestContext(ctx, msg.RequestContextID, msg.Consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgStartRequestContext handles MsgStartRequestContext
func handleMsgStartRequestContext(ctx sdk.Context, k Keeper, msg MsgStartRequestContext) (*sdk.Result, error) {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextID, true); err != nil {
		return nil, err
	}

	if err := k.StartRequestContext(ctx, msg.RequestContextID, msg.Consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgKillRequestContext handles MsgKillRequestContext
func handleMsgKillRequestContext(ctx sdk.Context, k Keeper, msg MsgKillRequestContext) (*sdk.Result, error) {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextID, true); err != nil {
		return nil, err
	}

	if err := k.KillRequestContext(ctx, msg.RequestContextID, msg.Consumer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgUpdateRequestContext handles MsgUpdateRequestContext
func handleMsgUpdateRequestContext(ctx sdk.Context, k Keeper, msg MsgUpdateRequestContext) (*sdk.Result, error) {
	if err := k.CheckAuthority(ctx, msg.Consumer, msg.RequestContextID, true); err != nil {
		return nil, err
	}

	if err := k.UpdateRequestContext(
		ctx, msg.RequestContextID, msg.Providers, 0, msg.ServiceFeeCap,
		msg.Timeout, msg.RepeatedFrequency, msg.RepeatedTotal, msg.Consumer,
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Consumer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgWithdrawEarnedFees handles MsgWithdrawEarnedFees
func handleMsgWithdrawEarnedFees(ctx sdk.Context, k Keeper, msg MsgWithdrawEarnedFees) (*sdk.Result, error) {
	if err := k.WithdrawEarnedFees(ctx, msg.Provider); err != nil {
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
