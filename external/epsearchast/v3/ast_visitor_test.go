package epsearchast_v3

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPreAndPostAndLeCalledOnAccept(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LE",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PostVisit").Return(nil).
		On("VisitLe", mock.Anything).Return(true, nil)

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.NoError(t, err)
}

func TestPreAndPostAndLtCalledOnAccept(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LT",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PostVisit").Return(nil).
		On("VisitLt", mock.Anything).Return(true, nil)

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.NoError(t, err)

}

func TestPreAndPostAndGeCalledOnAccept(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "GE",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PostVisit").Return(nil).
		On("VisitGe", mock.Anything).Return(true, nil)

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.NoError(t, err)

}

func TestPreAndPostAndGtCalledOnAccept(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "GT",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PostVisit").Return(nil).
		On("VisitGt", mock.Anything).Return(true, nil)

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.NoError(t, err)

}

func TestPreAndPostAndEqCalledOnAccept(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "EQ",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PostVisit").Return(nil).
		On("VisitEq", mock.Anything).Return(true, nil)

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.NoError(t, err)

}

func TestPreAndPostAndLikeCalledOnAccept(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LIKE",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PostVisit").Return(nil).
		On("VisitLike", mock.Anything).Return(true, nil)

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.NoError(t, err)

}

func TestPreAndPostAndInCalledOnAccept(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "IN",
	"args": [
		"status",
		"paid",
		"pending"
	]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PostVisit").Return(nil).
		On("VisitIn", mock.Anything).Return(true, nil)

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.NoError(t, err)

}

func TestPreOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LE",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")
}

func TestPreVisitAndLeAndPostVisitCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LE",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("VisitLe", mock.Anything).Return(true, nil).
		On("PostVisit").Return(fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")
}
func TestPreAndLeCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LE",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("VisitLe", mock.Anything).Return(true, fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")
}

func TestPreAndLtCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LT",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("VisitLt", mock.Anything).Return(true, fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")

}

func TestPreAndGeCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "GE",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("VisitGe", mock.Anything).Return(true, fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")

}

func TestPreAndGtCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "GT",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("VisitGt", mock.Anything).Return(true, fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")

}

func TestPreAndEqCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "EQ",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("VisitEq", mock.Anything).Return(true, fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")

}

func TestPreAndLikeCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LIKE",
	"args": [ "amount",  "5"]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("VisitLike", mock.Anything).Return(true, fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")

}

func TestPreAndInCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "IN",
	"args": [
		"status",
		"paid",
		"pending"
	]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("VisitIn", mock.Anything).Return(true, fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")
}

func TestPreAndPostAndEqAndAndCalledOnAccept(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "AND",
	"children": [
	{
		"type": "EQ",
		"args": [ "amount",  "5"]
	},{
		"type": "EQ",
		"args": [ "amount",  "5"]
	}]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PostVisit").Return(nil).
		On("VisitEq", mock.Anything).Return(true, nil).
		On("PreVisitAnd", mock.Anything).Return(true, nil).
		On("PostVisitAnd", mock.Anything).Return(nil)

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.NoError(t, err)

}

func TestPreAndPreVisitAndCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "AND",
	"children": [
	{
		"type": "EQ",
		"args": [ "amount",  "5"]
	},{
		"type": "EQ",
		"args": [ "amount",  "5"]
	}]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PreVisitAnd", mock.Anything).Return(true, fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")
}

func TestPreAndPreVisitAndEqCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "AND",
	"children": [
	{
		"type": "EQ",
		"args": [ "amount",  "5"]
	},{
		"type": "EQ",
		"args": [ "amount",  "5"]
	}]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PreVisitAnd", mock.Anything).Return(true, nil).
		On("VisitEq", mock.Anything).Return(true, fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")
}

func TestPreAndPreVisitAndEqAndPostVisitCalledOnAcceptWithError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "AND",
	"children": [
	{
		"type": "EQ",
		"args": [ "amount",  "5"]
	},{
		"type": "EQ",
		"args": [ "amount",  "5"]
	}]
}
`

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PreVisitAnd", mock.Anything).Return(true, nil).
		On("VisitEq", mock.Anything).Return(true, nil).
		On("PostVisitAnd", mock.Anything).Return(fmt.Errorf("foo"))

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = astNode.Accept(mockObj)

	// Verification
	require.ErrorContains(t, err, "foo")
}

type MyMockedVisitor struct {
	mock.Mock
}

func (m *MyMockedVisitor) PreVisit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MyMockedVisitor) PostVisit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MyMockedVisitor) PreVisitAnd(astNode *AstNode) (bool, error) {
	args := m.Called(astNode)
	return args.Bool(0), args.Error(1)
}

func (m *MyMockedVisitor) PostVisitAnd(astNode *AstNode) error {
	args := m.Called(astNode)
	return args.Error(0)
}

func (m *MyMockedVisitor) VisitIn(astNode *AstNode) (bool, error) {
	args := m.Called(astNode)
	return args.Bool(0), args.Error(1)
}

func (m *MyMockedVisitor) VisitEq(astNode *AstNode) (bool, error) {
	args := m.Called(astNode)
	return args.Bool(0), args.Error(1)
}

func (m *MyMockedVisitor) VisitLe(astNode *AstNode) (bool, error) {
	args := m.Called(astNode)
	return args.Bool(0), args.Error(1)
}

func (m *MyMockedVisitor) VisitLt(astNode *AstNode) (bool, error) {
	args := m.Called(astNode)
	return args.Bool(0), args.Error(1)
}

func (m *MyMockedVisitor) VisitGe(astNode *AstNode) (bool, error) {
	args := m.Called(astNode)
	return args.Bool(0), args.Error(1)
}

func (m *MyMockedVisitor) VisitGt(astNode *AstNode) (bool, error) {
	args := m.Called(astNode)
	return args.Bool(0), args.Error(1)
}

func (m *MyMockedVisitor) VisitLike(astNode *AstNode) (bool, error) {
	args := m.Called(astNode)
	return args.Bool(0), args.Error(1)
}

var _ AstVisitor = (*MyMockedVisitor)(nil)
