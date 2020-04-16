package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irismod/service/keeper"
	"github.com/irismod/service/simapp"
)

var (
	testAddr = sdk.AccAddress([]byte("testaddr"))

	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"}}`
)

type KeeperTestSuite struct {
	suite.Suite

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper *keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.keeper = &app.ServiceKeeper
}

func (suite *KeeperTestSuite) TestDefineService() {
	err := suite.keeper.AddServiceDefinition(suite.ctx, testServiceName, testServiceDesc, testServiceTags, testAddr, testAuthorDesc, testSchemas)
	suite.NoError(err)

	svcDef, found := suite.keeper.GetServiceDefinition(suite.ctx, testServiceName)
	suite.True(found)

	suite.Equal(testServiceName, svcDef.Name)
	suite.Equal(testServiceTags, svcDef.Tags)
	suite.Equal(testAddr, svcDef.Author)
	suite.Equal(testSchemas, svcDef.Schemas)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
