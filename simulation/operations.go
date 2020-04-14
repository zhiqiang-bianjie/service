package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/irismod/service/keeper"
	"github.com/irismod/service/simapp/helpers"
	simappparams "github.com/irismod/service/simapp/params"
	"github.com/irismod/service/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgDefineService = "op_weight_msg_define_service"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams,
	cdc *codec.Codec,
	ak types.AccountKeeper,
	k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgDefineService int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgDefineService, &weightMsgDefineService, nil,
		func(_ *rand.Rand) {
			weightMsgDefineService = simappparams.DefaultWeightMsgDefineService
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgDefineService,
			SimulateMsgDefineService(ak, k),
		),
	}
}

// SimulateMsgDefineService generates a MsgDefineService with random values.
func SimulateMsgDefineService(ak types.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		simAccount, _ := simulation.RandomAcc(r, accs)

		serviceName := simulation.RandStringOfLength(r, 20)
		serviceDescription := simulation.RandStringOfLength(r, 50)
		authorDescription := simulation.RandStringOfLength(r, 50)
		tags := []string{simulation.RandStringOfLength(r, 20), simulation.RandStringOfLength(r, 20)}
		schemas := `{"input":{"type":"object"},"output":{"type":"object"}}`

		account := ak.GetAccount(ctx, simAccount.Address)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		msg := types.NewMsgDefineService(serviceName, serviceDescription, tags, simAccount.Address, authorDescription, schemas)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		if _, _, err := app.Deliver(tx); err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
