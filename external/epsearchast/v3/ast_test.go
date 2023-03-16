package epsearchast_v3

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmptyJsonReturnsError(t *testing.T) {
	// Fixture Setup

	// Execute SUT
	astNode, err := GetAst("")

	// Verify
	require.Error(t, err)
	require.Nil(t, astNode)
}

func TestEmptyObjectReturnsError(t *testing.T) {
	// Fixture Setup

	// Execute SUT
	astNode, err := GetAst("{}")

	// Verify
	require.ErrorContains(t, err, "error validating filter")
	require.Nil(t, astNode)
}

func TestInvalidObjectReturnsError(t *testing.T) {
	// Fixture Setup

	// Execute SUT
	astNode, err := GetAst(`{"type": "FOO"}`)

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "unknown operator FOO")
	require.Nil(t, astNode)
}

func TestValidObjectWithEqReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "EQ",
	"first_arg": "status",
	"second_arg": "paid"
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestValidObjectWithLeReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LE",
	"first_arg": "orders",
	"second_arg": "5"
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestValidObjectWithLtReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LT",
	"first_arg": "orders",
	"second_arg": "5"
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestValidObjectWithGeReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "GE",
	"first_arg": "orders",
	"second_arg": "5"
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestValidObjectWithGtReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "GT",
	"first_arg": "orders",
	"second_arg": "5"
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestValidObjectWithLikeReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "LIKE",
	"first_arg": "status",
	"second_arg": "p*"
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestEqWithChildReturnsError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "EQ",
	"first_arg": "status",
	"second_arg": "paid",
	"children": [{
		"type": "EQ",
		"first_arg": "status",
		"second_arg": "paid"
	}]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "should not have any children")
	require.Nil(t, astNode)
}

func TestValidObjectWithInReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "IN",
	"args": ["orders", "5"]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestOneArgumentToInReturnsError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "IN",
	"args": ["orders"]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient number of arguments to in")
	require.Nil(t, astNode)
}

func TestInvalidOperatorReturnsError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "FOO",
	"args": ["orders"]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "unknown operator FOO")
	require.Nil(t, astNode)
}

func TestInWithChildReturnsError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "IN",
	"args": ["status", "paid"],
	"children": [{
		"type": "EQ",
		"first_arg": "status",
		"second_arg": "paid"
	}]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "should not have any children")
	require.Nil(t, astNode)
}

func TestAndReturnsErrorWithOneChildren(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "AND",
	"children": [{
		"type": "EQ",
		"first_arg": "status",
		"second_arg": "paid"
	}]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "and should have at least two children")
	require.Nil(t, astNode)
}

func TestAndReturnsErrorWithAnInvalidChild(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "AND",
	"children": [{
		"type": "EQ",
		"first_arg": "status",
		"second_arg": "paid"
	},
	{
		"type": "FOO"
	}
	]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "unknown operator FOO")
	require.Nil(t, astNode)
}
