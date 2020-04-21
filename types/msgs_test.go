package types

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	testCoin1 = sdk.NewCoin(types.ServiceDepositCoinDenom, sdk.NewInt(10000))
	testCoin2 = sdk.NewCoin(types.ServiceDepositCoinDenom, sdk.NewInt(100))

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

// TestMsgDefineServiceRoute tests Route for MsgDefineService
func TestMsgDefineServiceRoute(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgDefineServiceType tests Type for MsgDefineService
func TestMsgDefineServiceType(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)

	require.Equal(t, "define_service", msg.Type())
}

// TestMsgDefineServiceValidation tests ValidateBasic for MsgDefineService
func TestMsgDefineServiceValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidLongDesc := strings.Repeat("d", MaxDescriptionLength+1)
	invalidMoreTags := strings.Split("t1,t2,t3,t4,t5,t6,t7,t8,t9,t10,t11", ",")
	invalidLongTags := []string{strings.Repeat("t", MaxTagLength+1)}
	invalidEmptyTags := []string{"t1", ""}
	invalidDuplicateTags := []string{"t1", "t1"}

	invalidSchemas := `{"input":"nonobject","output":"nonobject"}`
	invalidSchemasNoInput := `{"output":{"type":"object"}}`
	invalidSchemasNoOutput := `{"input":{"type":"object"}}`

	testMsgs := []MsgDefineService{
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),            // valid msg
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, emptyAddress, testAuthorDesc, testSchemas),          // missing author address
		NewMsgDefineService(invalidName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),                // service name contains illegal characters
		NewMsgDefineService(invalidLongName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),            // too long service name
		NewMsgDefineService(testServiceName, invalidLongDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas),            // too long service description
		NewMsgDefineService(testServiceName, testServiceDesc, invalidMoreTags, testAuthor, testAuthorDesc, testSchemas),            // too many tags
		NewMsgDefineService(testServiceName, testServiceDesc, invalidLongTags, testAuthor, testAuthorDesc, testSchemas),            // too long tag
		NewMsgDefineService(testServiceName, testServiceDesc, invalidEmptyTags, testAuthor, testAuthorDesc, testSchemas),           // empty tag
		NewMsgDefineService(testServiceName, testServiceDesc, invalidDuplicateTags, testAuthor, testAuthorDesc, testSchemas),       // duplicate tags
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, invalidLongDesc, testSchemas),           // too long author description
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, invalidSchemas),         // invalid schemas
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, invalidSchemasNoInput),  // missing input schema
		NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, invalidSchemasNoOutput), // missing output schema
	}

	testCases := []struct {
		msg     MsgDefineService
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing author address"},
		{testMsgs[2], false, "service name contains illegal characters"},
		{testMsgs[3], false, "too long service name"},
		{testMsgs[4], false, "too long service description"},
		{testMsgs[5], false, "too many tags"},
		{testMsgs[6], false, "too long tag"},
		{testMsgs[7], false, "empty tag"},
		{testMsgs[8], false, "duplicate tags"},
		{testMsgs[9], false, "too long author description"},
		{testMsgs[10], false, "invalid schemas"},
		{testMsgs[11], false, "missing input schema"},
		{testMsgs[12], false, "missing output schema"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgDefineServiceGetSignBytes tests GetSignBytes for MsgDefineService
func TestMsgDefineServiceGetSignBytes(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgDefineService","value":{"author":"cosmos1w3jhxapdv96hg6r0wg0dldpe","author_description":"test-author-desc","description":"test-service-desc","name":"test-service","schemas":"{\"input\":{\"type\":\"object\"},\"output\":{\"type\":\"object\"}}","tags":["tag1","tag2"]}}`
	require.Equal(t, expected, string(res))
}

// TestMsgDefineServiceGetSigners tests GetSigners for MsgDefineService
func TestMsgDefineServiceGetSigners(t *testing.T) {
	msg := NewMsgDefineService(testServiceName, testServiceDesc, testServiceTags, testAuthor, testAuthorDesc, testSchemas)
	res := msg.GetSigners()

	expected := "[746573742D617574686F72]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgBindServiceRoute tests Route for MsgBindService
func TestMsgBindServiceRoute(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testMinRespTime)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgBindServiceType tests Type for MsgBindService
func TestMsgBindServiceType(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testMinRespTime)

	require.Equal(t, "bind_service", msg.Type())
}

// TestMsgBindServiceValidation tests ValidateBasic for MsgBindService
func TestMsgBindServiceValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidDeposit := sdk.Coins{}
	invalidDenomDeposit := sdk.NewCoins(sdk.NewCoin("nonstake", sdk.NewInt(1000)))
	invalidMinRespTime := uint64(0)

	invalidPricing := `{"price":"1stake","other":"notallowedfield"}`
	invalidSymbolPricing := `{"price":"0.5invalidsymbol"}`
	invalidPromotionTimePricing := `{"price":"1stake","promotions_by_time":` +
		`[{"start_time":"2018-10-10T13:30:30","end_time":"2019-10-10T13:30:30Z","discount":"0.8"}]}`
	invalidPromotionVolPricing := `{"price":"1stake","promotions_by_volume":` +
		`[{"volume":0,"discount":"0.7"}]}`

	testMsgs := []MsgBindService{
		NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testMinRespTime),                 // valid msg
		NewMsgBindService(testServiceName, emptyAddress, testDeposit, testPricing, testMinRespTime),                 // missing provider address
		NewMsgBindService(invalidName, testProvider, testDeposit, testPricing, testMinRespTime),                     // service name contains illegal characters
		NewMsgBindService(invalidLongName, testProvider, testDeposit, testPricing, testMinRespTime),                 // too long service name
		NewMsgBindService(testServiceName, testProvider, invalidDeposit, testPricing, testMinRespTime),              // invalid deposit
		NewMsgBindService(testServiceName, testProvider, invalidDenomDeposit, testPricing, testMinRespTime),         // invalid deposit denom
		NewMsgBindService(testServiceName, testProvider, testDeposit, "", testMinRespTime),                          // missing pricing
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidPricing, testMinRespTime),              // invalid Pricing JSON Schema instance
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidSymbolPricing, testMinRespTime),        // invalid pricing symbol
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidPromotionTimePricing, testMinRespTime), // invalid promotion time lack of time zone
		NewMsgBindService(testServiceName, testProvider, testDeposit, invalidPromotionVolPricing, testMinRespTime),  // invalid promotion volume
		NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, invalidMinRespTime),              // invalid minimum response time                               // invalid promotion volume
	}

	testCases := []struct {
		msg     MsgBindService
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "service name contains illegal characters"},
		{testMsgs[3], false, "too long service name"},
		{testMsgs[4], false, "invalid deposit"},
		{testMsgs[5], false, "invalid deposit denom"},
		{testMsgs[6], false, "missing pricing"},
		{testMsgs[7], false, "invalid Pricing JSON Schema instance"},
		{testMsgs[8], false, "invalid pricing symbol"},
		{testMsgs[9], false, "invalid promotion time lack of time zone"},
		{testMsgs[10], false, "invalid promotion volume"},
		{testMsgs[11], false, "invalid minimum response time"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgBindServiceGetSignBytes tests GetSignBytes for MsgBindService
func TestMsgBindServiceGetSignBytes(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testMinRespTime)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgBindService","value":{"deposit":[{"amount":"10000","denom":"stake"}],"min_resp_time":"50","pricing":"{\"price\":\"1stake\"}","provider":"cosmos1w3jhxapdwpex7anfv3jhy8anr90","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgBindServiceGetSigners tests GetSigners for MsgBindService
func TestMsgBindServiceGetSigners(t *testing.T) {
	msg := NewMsgBindService(testServiceName, testProvider, testDeposit, testPricing, testMinRespTime)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgUpdateServiceBindingRoute tests Route for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingRoute(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "", 0)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgUpdateServiceBindingType tests Type for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingType(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "", 0)

	require.Equal(t, "update_service_binding", msg.Type())
}

// TestMsgUpdateServiceBindingValidation tests ValidateBasic for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}
	emptyAddedDeposit := sdk.Coins{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidDenomDeposit := sdk.NewCoins(sdk.NewCoin("nonstake", sdk.NewInt(1000)))

	invalidPricing := `{"price":"1stake","other":"notallowedfield"}`
	invalidSymbolPricing := `{"price":"1invalidsymbol"}`
	invalidPromotionTimePricing := `{"price":"1stake","promotions_by_time":` +
		`[{"start_time":"2018-10-10T13:30:30","end_time":"2019-10-10T13:30:30Z","discount":"0.8"}]}`
	invalidPromotionVolPricing := `{"price":"1stake","promotions_by_volume":` +
		`[{"volume":0,"discount":"0.7"}]}`

	testMsgs := []MsgUpdateServiceBinding{
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, testPricing, testMinRespTime),                 // valid msg
		NewMsgUpdateServiceBinding(testServiceName, testProvider, emptyAddedDeposit, testPricing, testMinRespTime),                // empty deposit is allowed
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "", testMinRespTime),                          // empty pricing is allowed
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, testPricing, 0),                               // 0 is allowed for minimum response time
		NewMsgUpdateServiceBinding(testServiceName, testProvider, emptyAddedDeposit, "", 0),                                       // deposit, pricing and min response time can be empty at the same time
		NewMsgUpdateServiceBinding(testServiceName, emptyAddress, testAddedDeposit, testPricing, testMinRespTime),                 // missing provider address
		NewMsgUpdateServiceBinding(invalidName, testProvider, testAddedDeposit, testPricing, testMinRespTime),                     // service name contains illegal characters
		NewMsgUpdateServiceBinding(invalidLongName, testProvider, testAddedDeposit, testPricing, testMinRespTime),                 // too long service name
		NewMsgUpdateServiceBinding(testServiceName, testProvider, invalidDenomDeposit, testPricing, testMinRespTime),              // invalid deposit denom
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidPricing, testMinRespTime),              // invalid Pricing JSON Schema instance
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidSymbolPricing, testMinRespTime),        // invalid pricing symbol
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidPromotionTimePricing, testMinRespTime), // invalid promotion time lack of time zone
		NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, invalidPromotionVolPricing, testMinRespTime),  // invalid promotion volume
	}

	testCases := []struct {
		msg     MsgUpdateServiceBinding
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], true, ""},
		{testMsgs[2], true, ""},
		{testMsgs[3], true, ""},
		{testMsgs[4], true, ""},
		{testMsgs[5], false, "missing provider address"},
		{testMsgs[6], false, "service name contains illegal characters"},
		{testMsgs[7], false, "too long service name"},
		{testMsgs[8], false, "invalid deposit denom"},
		{testMsgs[9], false, "invalid Pricing JSON Schema instance"},
		{testMsgs[10], false, "invalid pricing symbol"},
		{testMsgs[11], false, "invalid promotion time lack of time zone"},
		{testMsgs[12], false, "invalid promotion volume"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgUpdateServiceBindingGetSignBytes tests GetSignBytes for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingGetSignBytes(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "", 0)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgUpdateServiceBinding","value":{"deposit":[{"amount":"100","denom":"stake"}],"min_resp_time":"0","pricing":"","provider":"cosmos1w3jhxapdwpex7anfv3jhy8anr90","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgUpdateServiceBindingGetSigners tests GetSigners for MsgUpdateServiceBinding
func TestMsgUpdateServiceBindingGetSigners(t *testing.T) {
	msg := NewMsgUpdateServiceBinding(testServiceName, testProvider, testAddedDeposit, "", 0)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgSetWithdrawAddressRoute tests Route for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressRoute(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testProvider, testWithdrawAddr)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgSetWithdrawAddressType tests Type for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressType(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testProvider, testWithdrawAddr)

	require.Equal(t, "set_withdraw_address", msg.Type())
}

// TestMsgSetWithdrawAddressValidation tests ValidateBasic for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	testMsgs := []MsgSetWithdrawAddress{
		NewMsgSetWithdrawAddress(testProvider, testWithdrawAddr), // valid msg
		NewMsgSetWithdrawAddress(emptyAddress, testWithdrawAddr), // missing provider address
		NewMsgSetWithdrawAddress(testProvider, emptyAddress),     // missing withdrawal address
	}

	testCases := []struct {
		msg     MsgSetWithdrawAddress
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "missing withdrawal address"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgSetWithdrawAddressGetSignBytes tests GetSignBytes for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressGetSignBytes(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testProvider, testWithdrawAddr)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgSetWithdrawAddress","value":{"provider":"cosmos1w3jhxapdwpex7anfv3jhy8anr90","withdraw_address":"cosmos1w3jhxapdwa5hg6rywfshwctv94skgerjv4ehxd4yrry"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgSetWithdrawAddressGetSigners tests GetSigners for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddressGetSigners(t *testing.T) {
	msg := NewMsgSetWithdrawAddress(testProvider, testWithdrawAddr)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgDisableServiceBindingRoute tests Route for MsgDisableServiceBinding
func TestMsgDisableServiceBindingRoute(t *testing.T) {
	msg := NewMsgDisableServiceBinding(testServiceName, testProvider)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgDisableServiceBindingType tests Type for MsgDisableServiceBinding
func TestMsgDisableServiceBindingType(t *testing.T) {
	msg := NewMsgDisableServiceBinding(testServiceName, testProvider)

	require.Equal(t, "disable_service_binding", msg.Type())
}

// TestMsgDisableServiceBindingValidation tests ValidateBasic for MsgDisableServiceBinding
func TestMsgDisableServiceBindingValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)

	testMsgs := []MsgDisableServiceBinding{
		NewMsgDisableServiceBinding(testServiceName, testProvider), // valid msg
		NewMsgDisableServiceBinding(testServiceName, emptyAddress), // missing provider address
		NewMsgDisableServiceBinding(invalidName, testProvider),     // service name contains illegal characters
		NewMsgDisableServiceBinding(invalidLongName, testProvider), // too long service name
	}

	testCases := []struct {
		msg     MsgDisableServiceBinding
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "service name contains illegal characters"},
		{testMsgs[3], false, "too long service name"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgDisableServiceBindingGetSignBytes tests GetSignBytes for MsgDisableServiceBinding
func TestMsgDisableServiceBindingGetSignBytes(t *testing.T) {
	msg := NewMsgDisableServiceBinding(testServiceName, testProvider)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgDisableServiceBinding","value":{"provider":"cosmos1w3jhxapdwpex7anfv3jhy8anr90","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgDisableServiceBindingGetSigners tests GetSigners for MsgDisableServiceBinding
func TestMsgDisableServiceBindingGetSigners(t *testing.T) {
	msg := NewMsgDisableServiceBinding(testServiceName, testProvider)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgEnableServiceBindingRoute tests Route for MsgEnableServiceBinding
func TestMsgEnableServiceBindingRoute(t *testing.T) {
	msg := NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgEnableServiceBindingType tests Type for MsgEnableServiceBinding
func TestMsgEnableServiceBindingType(t *testing.T) {
	msg := NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit)

	require.Equal(t, "enable_service_binding", msg.Type())
}

// TestMsgEnableServiceBindingValidation tests ValidateBasic for MsgEnableServiceBinding
func TestMsgEnableServiceBindingValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}
	emptyAddedDeposit := sdk.Coins{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)
	invalidDenomDeposit := sdk.NewCoins(sdk.NewCoin("nonstake", sdk.NewInt(1000)))

	testMsgs := []MsgEnableServiceBinding{
		NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit),    // valid msg
		NewMsgEnableServiceBinding(testServiceName, testProvider, emptyAddedDeposit),   // empty deposit is allowed
		NewMsgEnableServiceBinding(testServiceName, emptyAddress, testAddedDeposit),    // missing provider address
		NewMsgEnableServiceBinding(invalidName, testProvider, testAddedDeposit),        // service name contains illegal characters
		NewMsgEnableServiceBinding(invalidLongName, testProvider, testAddedDeposit),    // too long service name
		NewMsgEnableServiceBinding(testServiceName, testProvider, invalidDenomDeposit), // invalid deposit denom
	}

	testCases := []struct {
		msg     MsgEnableServiceBinding
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], true, ""},
		{testMsgs[2], false, "missing provider address"},
		{testMsgs[3], false, "service name contains illegal characters"},
		{testMsgs[4], false, "too long service name"},
		{testMsgs[5], false, "invalid deposit denom"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgEnableServiceBindingGetSignBytes tests GetSignBytes for MsgEnableServiceBinding
func TestMsgEnableServiceBindingGetSignBytes(t *testing.T) {
	msg := NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgEnableServiceBinding","value":{"deposit":[{"amount":"100","denom":"stake"}],"provider":"cosmos1w3jhxapdwpex7anfv3jhy8anr90","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgEnableServiceBindingGetSigners tests GetSigners for MsgEnableServiceBinding
func TestMsgEnableServiceBindingGetSigners(t *testing.T) {
	msg := NewMsgEnableServiceBinding(testServiceName, testProvider, testAddedDeposit)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}

// TestMsgRefundServiceDepositRoute tests Route for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositRoute(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider)

	require.Equal(t, RouterKey, msg.Route())
}

// TestMsgRefundServiceDepositType tests Type for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositType(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider)

	require.Equal(t, "refund_service_deposit", msg.Type())
}

// TestMsgRefundServiceDepositValidation tests ValidateBasic for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositValidation(t *testing.T) {
	emptyAddress := sdk.AccAddress{}

	invalidName := "invalid/service/name"
	invalidLongName := strings.Repeat("s", MaxNameLength+1)

	testMsgs := []MsgRefundServiceDeposit{
		NewMsgRefundServiceDeposit(testServiceName, testProvider), // valid msg
		NewMsgRefundServiceDeposit(testServiceName, emptyAddress), // missing provider address
		NewMsgRefundServiceDeposit(invalidName, testProvider),     // service name contains illegal characters
		NewMsgRefundServiceDeposit(invalidLongName, testProvider), // too long service name
	}

	testCases := []struct {
		msg     MsgRefundServiceDeposit
		expPass bool
		errMsg  string
	}{
		{testMsgs[0], true, ""},
		{testMsgs[1], false, "missing provider address"},
		{testMsgs[2], false, "service name contains illegal characters"},
		{testMsgs[3], false, "too long service name"},
	}

	for i, tc := range testCases {
		err := tc.msg.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, "Msg %d failed: %v", i, err)
		} else {
			require.Error(t, err, "Invalid Msg %d passed: %s", i, tc.errMsg)
		}
	}
}

// TestMsgRefundServiceDepositGetSignBytes tests GetSignBytes for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositGetSignBytes(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider)
	res := msg.GetSignBytes()

	expected := `{"type":"irismod/service/MsgRefundServiceDeposit","value":{"provider":"cosmos1w3jhxapdwpex7anfv3jhy8anr90","service_name":"test-service"}}`
	require.Equal(t, expected, string(res))
}

// TestMsgRefundServiceDepositGetSigners tests GetSigners for MsgRefundServiceDeposit
func TestMsgRefundServiceDepositGetSigners(t *testing.T) {
	msg := NewMsgRefundServiceDeposit(testServiceName, testProvider)
	res := msg.GetSigners()

	expected := "[746573742D70726F7669646572]"
	require.Equal(t, expected, fmt.Sprintf("%v", res))
}
