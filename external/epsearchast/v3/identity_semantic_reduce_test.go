package epsearchast_v3

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var binOpsForTest = []string{
	"LE",
	"LT",
	"EQ",
	"GT",
	"GE",
	"CONTAINS",
	"LIKE",
	"ILIKE",
	"TEXT",
}

var unaryOpsForTest = []string{
	"IS_NULL",
}

var varOpsForTest = []string{
	"IN",
	"CONTAINS_ANY",
	"CONTAINS_ALL",
}

func TestIdentitySemanticReducerWithBinaryOperators(t *testing.T) {
	for _, binOp := range binOpsForTest {
		t.Run(fmt.Sprintf("%s", binOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount",  "5"]
			}`, binOp)

			astNode, err := GetAst(astJson)

			var qb SemanticReducer[AstNode] = IdentitySemanticReducer{}

			expectedAstJson, err := json.MarshalIndent(astNode, "", "  ")
			require.NoError(t, err)

			// Execute SUT
			actualAst, err := SemanticReduceAst(astNode, qb)
			require.NoError(t, err)

			// Verification
			actualAstJson, err := json.MarshalIndent(actualAst, "", "  ")
			require.NoError(t, err)

			require.Equal(t, string(expectedAstJson), string(actualAstJson))
		})
	}
}

func TestIdentitySemanticReducerWithUnaryOperators(t *testing.T) {
	for _, unaryOp := range unaryOpsForTest {
		t.Run(fmt.Sprintf("%s", unaryOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount"]
			}`, unaryOp)

			astNode, err := GetAst(astJson)

			var qb SemanticReducer[AstNode] = IdentitySemanticReducer{}

			expectedAstJson, err := json.MarshalIndent(astNode, "", "  ")
			require.NoError(t, err)

			// Execute SUT
			actualAst, err := SemanticReduceAst(astNode, qb)
			require.NoError(t, err)

			// Verification
			actualAstJson, err := json.MarshalIndent(actualAst, "", "  ")
			require.NoError(t, err)

			require.Equal(t, string(expectedAstJson), string(actualAstJson))
		})
	}
}

func TestIdentitySemanticReducerWithVarargOperators(t *testing.T) {
	for _, varOp := range varOpsForTest {
		t.Run(fmt.Sprintf("%s", varOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "state", "CA", "NY", "TX"]
			}`, varOp)

			astNode, err := GetAst(astJson)

			var qb SemanticReducer[AstNode] = IdentitySemanticReducer{}

			expectedAstJson, err := json.MarshalIndent(astNode, "", "  ")
			require.NoError(t, err)

			// Execute SUT
			actualAst, err := SemanticReduceAst(astNode, qb)
			require.NoError(t, err)

			// Verification
			actualAstJson, err := json.MarshalIndent(actualAst, "", "  ")
			require.NoError(t, err)

			require.Equal(t, string(expectedAstJson), string(actualAstJson))
		})
	}
}

func TestSimpleRecursiveStructureWithAnd(t *testing.T) {
	//Fixture Setup
	//language=JSON
	astJson := `
				{
					"type":  "AND",
					"children": [
					{
						"type": "IN",
						"args": ["status", "new", "paid"]
					},
					{
						"type": "GE",
						"args": [ "amount",  "5"]
					}
					]
				}
				`

	astNode, err := GetAst(astJson)

	var qb SemanticReducer[AstNode] = IdentitySemanticReducer{}

	expectedAstJson, err := json.MarshalIndent(astNode, "", "  ")
	require.NoError(t, err)

	// Execute SUT
	actualAst, err := SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	actualAstJson, err := json.MarshalIndent(actualAst, "", "  ")
	require.NoError(t, err)

	require.Equal(t, string(expectedAstJson), string(actualAstJson))
}

func TestSimpleRecursiveStructureWithOr(t *testing.T) {
	//Fixture Setup
	//language=JSON
	astJson := `
				{
					"type":  "OR",
					"children": [
					{
						"type": "IN",
						"args": ["status", "new", "paid"]
					},
					{
						"type": "GE",
						"args": [ "amount",  "5"]
					}
					]
				}
				`

	astNode, err := GetAst(astJson)

	var qb SemanticReducer[AstNode] = IdentitySemanticReducer{}

	expectedAstJson, err := json.MarshalIndent(astNode, "", "  ")
	require.NoError(t, err)

	// Execute SUT
	actualAst, err := SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	actualAstJson, err := json.MarshalIndent(actualAst, "", "  ")
	require.NoError(t, err)

	require.Equal(t, string(expectedAstJson), string(actualAstJson))
}
