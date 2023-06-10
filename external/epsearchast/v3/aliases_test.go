package epsearchast_v3

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApplyAliasesWithFieldAliasSpecified(t *testing.T) {
	// Fixture Setup
	// language=JSON
	inputAstJson := `
	{
		"type": "EQ",
		"args": [ "payment_status",  "paid"]
	}`

	// language=JSON
	expectedAstJson := `
	{
		"type": "EQ",
		"args": [ "status",  "paid"]
	}`

	inputAstNode, err := GetAst(inputAstJson)
	require.NoError(t, err)

	expectedAstNode, err := GetAst(expectedAstJson)
	require.NoError(t, err)

	// Execute SUT
	aliasedAst, err := ApplyAliases(inputAstNode, map[string]string{"payment_status": "status"})

	// Verify
	require.NoError(t, err)
	require.NotNil(t, aliasedAst)

	require.Equal(t, expectedAstNode, aliasedAst)
}

func TestApplyAliasesReturnsCorrectAstWhenNoAliasesApply(t *testing.T) {
	// Fixture Setup
	// language=JSON
	inputAstJson := `
	{
		"type": "EQ",
		"args": [ "status",  "paid"]
	}`

	// language=JSON
	expectedAstJson := `
	{
		"type": "EQ",
		"args": [ "status",  "paid"]
	}`

	inputAstNode, err := GetAst(inputAstJson)
	require.NoError(t, err)

	expectedAstNode, err := GetAst(expectedAstJson)
	require.NoError(t, err)

	// Execute SUT
	aliasedAst, err := ApplyAliases(inputAstNode, map[string]string{"customer_name": "customer.name"})

	// Verify
	require.NoError(t, err)
	require.NotNil(t, aliasedAst)

	require.Equal(t, expectedAstNode, aliasedAst)
}

func TestApplyAliasesReturnsCorrectAstWhenAliasMapIsEmpty(t *testing.T) {
	// Fixture Setup
	// language=JSON
	inputAstJson := `
	{
		"type": "EQ",
		"args": [ "status",  "paid"]
	}`

	// language=JSON
	expectedAstJson := `
	{
		"type": "EQ",
		"args": [ "status",  "paid"]
	}`

	inputAstNode, err := GetAst(inputAstJson)
	require.NoError(t, err)

	expectedAstNode, err := GetAst(expectedAstJson)
	require.NoError(t, err)

	// Execute SUT
	aliasedAst, err := ApplyAliases(inputAstNode, map[string]string{})

	// Verify
	require.NoError(t, err)
	require.NotNil(t, aliasedAst)

	require.Equal(t, expectedAstNode, aliasedAst)
}

func TestApplyAliasesReturnsCorrectAstWhenAliasTwoFieldsAreAliasedInAnAnd(t *testing.T) {
	// Fixture Setup
	// language=JSON
	inputAstJson := `
	{ 
		"type": "AND",
		"children": [
			{
				"type": "EQ",
				"args": [ "payment_status",  "paid"]
			},
			{
				"type": "LIKE",
				"args": [ "customer_name",  "Ron*"]
			},
			{
				"type": "EQ",
				"args": [ "customer.email",  "ron@swanson.com"]
			},
			{
				"type": "IS_NULL",
				"args": [ "billing-email"]
			}		
		]
	}
	`

	// language=JSON
	expectedAstJson := `
	{ 
		"type": "AND",
		"children": [
			{
				"type": "EQ",
				"args": [ "status",  "paid"]
			},
			{
				"type": "LIKE",
				"args": [ "customer.name",  "Ron*"]
			},
			{
				"type": "EQ",
				"args": [ "customer.email",  "ron@swanson.com"]
			},
			{
				"type": "IS_NULL",
				"args": [ "billing.email"]
			}	
		]
	}`

	inputAstNode, err := GetAst(inputAstJson)
	require.NoError(t, err)

	expectedAstNode, err := GetAst(expectedAstJson)
	require.NoError(t, err)

	// Execute SUT
	aliasedAst, err := ApplyAliases(inputAstNode, map[string]string{"payment_status": "status", "customer_name": "customer.name", "billing-email": "billing.email"})

	// Verify
	require.NoError(t, err)
	require.NotNil(t, aliasedAst)

	require.Equal(t, expectedAstNode, aliasedAst)
}
