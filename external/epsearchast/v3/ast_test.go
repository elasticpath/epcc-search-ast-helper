package epsearchast_v3

import (
	"github.com/stretchr/testify/require"
	"net/url"
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
	require.ErrorAs(t, err, &ValidationErr{})
	require.Nil(t, astNode)
}

func TestInvalidQueryReturnsParsingError(t *testing.T) {
	// Fixture Setup

	// Execute SUT
	astNode, err := GetAst("{!@")

	// Verify
	require.ErrorContains(t, err, "could not parse filter")
	require.ErrorAs(t, err, &ParsingErr{})
	require.Nil(t, astNode)
}

func TestInvalidObjectReturnsError(t *testing.T) {
	// Fixture Setup

	// Execute SUT
	astNode, err := GetAst(`{"type": "FOO"}`)

	// Verify
	require.Error(t, err)
	require.EqualError(t, err, "error validating filter: unsupported operator foo()")
	require.ErrorAs(t, err, &ValidationErr{})
	require.Nil(t, astNode)
}

func TestValidObjectWithEqReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "EQ",
	"args": [ "status",  "paid"]
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
	"args": [ "orders",  "5"]
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
	"args": [ "orders",  "5"]
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
	"args": [ "orders",  "5"]
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
	"args": [ "orders",  "5"]
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
	"args": [ "status",  "p*"]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestValidObjectWithILikeReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "ILIKE",
	"args": [ "status",  "p*"]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestValidObjectWithContainsReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "CONTAINS",
	"args": [ "status",  "paid"]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestValidObjectWithTextReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "TEXT",
	"args": [ "name",  "John"]
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
	"args": [ "status",  "paid"],
	"children": [{
		"type": "EQ",
		"args": [ "status",  "paid"]
	}]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "should not have any children")
	require.ErrorAs(t, err, &ValidationErr{})
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
	require.ErrorAs(t, err, &ValidationErr{})
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
	require.ErrorContains(t, err, "unsupported operator foo()")
	require.ErrorAs(t, err, &ValidationErr{})
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
		"args": [ "status",  "paid"]
	}]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "should not have any children")
	require.ErrorAs(t, err, &ValidationErr{})
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
		"args": [ "status",  "paid"]
	}]
}
`
	// Execute SUT
	astNode, err := GetAst(jsonTxt)

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "and should have at least two children")
	require.ErrorAs(t, err, &ValidationErr{})
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
		"args": [ "status",  "paid"]
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
	require.ErrorContains(t, err, "unsupported operator foo()")
	require.ErrorAs(t, err, &ValidationErr{})
	require.Nil(t, astNode)
}

func TestValidObjectThatIsUrlEncodedReturnsAst(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
	"type": "EQ",
	"args": [ "status",  "paid"]
}
`
	// Execute SUT
	astNode, err := GetAst(url.QueryEscape(jsonTxt))

	// Verify
	require.NoError(t, err)
	require.NotNil(t, astNode)
}

func TestInValidObjectThatIsUrlEncodedReturnsError(t *testing.T) {
	// Fixture Setup
	jsonTxt := `
{
	"type": "EQ",
	"args": [ "status",  "paid"]
`
	// Execute SUT
	astNode, err := GetAst(url.QueryEscape(jsonTxt))

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "unexpected end of JSON input")
	require.Nil(t, astNode)
}

func TestInValidUrlEncodingReturnsError(t *testing.T) {
	// Fixture Setup
	// Execute SUT
	astNode, err := GetAst("%4")

	// Verify
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid URL escape")
	require.Nil(t, astNode)
}
