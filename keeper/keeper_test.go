package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irismod/service/keeper"
	"github.com/irismod/service/simapp"
	"github.com/irismod/service/types"
)

var (
	testCoin1 = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000))
	testCoin2 = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))

	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthor      = sdk.AccAddress([]byte("test-author"))
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"}}`

	testProvider     = sdk.AccAddress([]byte("test-provider"))
	testDeposit      = sdk.NewCoins(testCoin1)
	testPricing      = `{"price":"1stake"}`
	testMinRespTime  = uint64(50)
	testWithdrawAddr = sdk.AccAddress([]byte("test-withdrawal-address"))
	testAddedDeposit = sdk.NewCoins(testCoin2)
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper *keeper.Keeper
	app    *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.app = app
	suite.keeper = &app.ServiceKeeper

	suite.app.InitChainer(suite.ctx, abci.RequestInitChain{})
}

func (suite *KeeperTestSuite) setServiceDefinition() {
	svcDef := types.NewServiceDefinition(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	suite.keeper.SetServiceDefinition(suite.ctx, svcDef)
}

func (suite *KeeperTestSuite) setServiceBinding(available bool, disabledTime time.Time) {
	svcBinding := types.NewServiceBinding(testServiceName, testProvider, testDeposit, testPricing, testMinRespTime, available, disabledTime)
	suite.keeper.SetServiceBinding(suite.ctx, svcBinding)

	pricing, _ := suite.keeper.ParsePricing(suite.ctx, testPricing)
	suite.keeper.SetPricing(suite.ctx, testServiceName, testProvider, pricing)
}

func (suite *KeeperTestSuite) TestDefineService() {
	err := suite.keeper.AddServiceDefinition(suite.ctx, testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	suite.NoError(err)

	svcDef, found := suite.keeper.GetServiceDefinition(suite.ctx, testServiceName)
	suite.True(found)

	suite.Equal(testServiceName, svcDef.Name)
	suite.Equal(testServiceDesc, svcDef.Description)
	suite.Equal(testServiceTags, svcDef.Tags)
	suite.Equal(testAuthor, svcDef.Author)
	suite.Equal(testAuthorDesc, svcDef.AuthorDescription)
	suite.Equal(testSchemas, svcDef.Schemas)
}

func (suite *KeeperTestSuite) TestBindService() {
	suite.setServiceDefinition()

	err := suite.keeper.AddServiceBinding(suite.ctx, testServiceName, testProvider, testDeposit, testPricing, testMinRespTime)
	suite.NoError(err)

	svcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, testServiceName, testProvider)
	suite.True(found)

	suite.Equal(testServiceName, svcBinding.ServiceName)
	suite.Equal(testProvider, svcBinding.Provider)
	suite.Equal(testDeposit, svcBinding.Deposit)
	suite.Equal(testPricing, svcBinding.Pricing)
	suite.Equal(testMinRespTime, svcBinding.MinRespTime)
	suite.True(svcBinding.Available)
	suite.True(svcBinding.DisabledTime.IsZero())

	// update binding
	newPricing := `{"price":"2stake"}`
	newMinRespTime := uint64(80)

	err = suite.keeper.UpdateServiceBinding(suite.ctx, svcBinding.ServiceName, svcBinding.Provider, testAddedDeposit, newPricing, newMinRespTime)
	suite.NoError(err)

	updatedSvcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, svcBinding.ServiceName, svcBinding.Provider)
	suite.True(found)

	suite.True(updatedSvcBinding.Deposit.IsEqual(svcBinding.Deposit.Add(testAddedDeposit...)))
	suite.Equal(newPricing, updatedSvcBinding.Pricing)
	suite.Equal(newMinRespTime, updatedSvcBinding.MinRespTime)
}

func (suite *KeeperTestSuite) TestSetWithdrawAddress() {
	suite.setServiceBinding(true, time.Time{})

	withdrawAddr := suite.keeper.GetWithdrawAddress(suite.ctx, testProvider)
	suite.Equal(testProvider, withdrawAddr)

	suite.keeper.SetWithdrawAddress(suite.ctx, testProvider, testWithdrawAddr)

	withdrawAddr = suite.keeper.GetWithdrawAddress(suite.ctx, testProvider)
	suite.Equal(testWithdrawAddr, withdrawAddr)
}

func (suite *KeeperTestSuite) TestDisableServiceBinding() {
	suite.setServiceBinding(true, time.Time{})

	currentTime := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(currentTime)

	err := suite.keeper.DisableServiceBinding(suite.ctx, testServiceName, testProvider)
	suite.NoError(err)

	svcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, testServiceName, testProvider)
	suite.True(found)

	suite.False(svcBinding.Available)
	suite.Equal(currentTime, svcBinding.DisabledTime)
}

func (suite *KeeperTestSuite) TestEnableServiceBinding() {
	disabledTime := time.Now().UTC()
	suite.setServiceBinding(false, disabledTime)

	err := suite.keeper.EnableServiceBinding(suite.ctx, testServiceName, testProvider, nil)
	suite.NoError(err)

	svcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, testServiceName, testProvider)
	suite.True(found)

	suite.True(svcBinding.Available)
	suite.True(svcBinding.DisabledTime.IsZero())
}

func (suite *KeeperTestSuite) TestRefundDeposit() {
	disabledTime := time.Now().UTC()
	suite.setServiceBinding(false, disabledTime)

	_, err := suite.app.BankKeeper.AddCoins(suite.ctx, suite.keeper.GetServiceDepositAccount(suite.ctx).GetAddress(), testDeposit)
	suite.NoError(err)

	params := suite.keeper.GetParams(suite.ctx)
	blockTime := disabledTime.Add(params.ArbitrationTimeLimit).Add(params.ComplaintRetrospect)
	suite.ctx = suite.ctx.WithBlockTime(blockTime)

	err = suite.keeper.RefundDeposit(suite.ctx, testServiceName, testProvider)
	suite.NoError(err)

	svcBinding, found := suite.keeper.GetServiceBinding(suite.ctx, testServiceName, testProvider)
	suite.True(found)

	suite.Equal(sdk.Coins(nil), svcBinding.Deposit)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
