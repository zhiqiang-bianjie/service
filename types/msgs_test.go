package types

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	testServiceName = "test-service"
	testServiceDesc = "test-service-desc"
	testServiceTags = []string{"tag1", "tag2"}
	testAuthor      = sdk.AccAddress([]byte("test-author"))
	testAuthorDesc  = "test-author-desc"
	testSchemas     = `{"input":{"type":"object"},"output":{"type":"object"}}`
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
