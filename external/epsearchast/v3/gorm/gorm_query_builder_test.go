package v3_gorm

import (
	"encoding/json"
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

var binOps = []testOp{
	{"LE", "<="},
	{"LT", "<"},
	{"EQ", "=="},
	{"GT", ">"},
	{"GE", ">="},
	{"LIKE", "LIKE"},
}

var varOps = []testOp{
	{"IN", "IN"},
}

type testOp struct {
	AstOp string
	SqlOp string
}

func TestSimpleBinaryOperatorFiltersGeneratesCorrectWhereClause(t *testing.T) {
	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp.AstOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			jsonTxt := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount",  "5"]
			}`, binOp.AstOp)

			astNode, err := epsearchast_v3.GetAst(jsonTxt)

			err = json.Unmarshal([]byte(jsonTxt), astNode)
			require.NoError(t, err)

			var sr epsearchast_v3.SemanticReducer[SubQuery] = DefaultGormQueryBuilder{}

			// Execute SUT
			query, err := epsearchast_v3.SemanticReduceAst(astNode, sr)

			// Verification

			require.NoError(t, err)

			require.Equal(t, fmt.Sprintf("amount %s ?", binOp.SqlOp), query.Clause)
			require.Equal(t, []interface{}{"5"}, query.Args)
		})
	}

}

func TestSimpleVariableOperatorFiltersGeneratesCorrectWhereClause(t *testing.T) {
	for _, varOp := range varOps {
		t.Run(fmt.Sprintf("%s", varOp.AstOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			jsonTxt := fmt.Sprintf(`
				{
				"type": "%s",
				"args": ["amount", "5", "6", "7"]
			}`, varOp.AstOp)

			astNode, err := epsearchast_v3.GetAst(jsonTxt)

			err = json.Unmarshal([]byte(jsonTxt), astNode)
			require.NoError(t, err)

			var sr epsearchast_v3.SemanticReducer[SubQuery] = DefaultGormQueryBuilder{}

			// Execute SUT
			query, err := epsearchast_v3.SemanticReduceAst(astNode, sr)

			// Verification

			require.NoError(t, err)

			require.Equal(t, fmt.Sprintf("amount %s ?", varOp.SqlOp), query.Clause)
			require.Equal(t, []interface{}{[]interface{}{"5", "6", "7"}}, query.Args)
		})
	}
}

func TestLikeFilterWildCards(t *testing.T) {
	genTest := func(astLiteral string, sqlLiteral string) func(t *testing.T) {
		return func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			jsonTxt := fmt.Sprintf(`
				{
				"type": "LIKE",
				"args": [ "email",  "%s"]
			}`, astLiteral)

			astNode, err := epsearchast_v3.GetAst(jsonTxt)

			err = json.Unmarshal([]byte(jsonTxt), astNode)
			require.NoError(t, err)

			var sr epsearchast_v3.SemanticReducer[SubQuery] = DefaultGormQueryBuilder{}

			// Execute SUT
			query, err := epsearchast_v3.SemanticReduceAst(astNode, sr)

			// Verification

			require.NoError(t, err)

			require.Equal(t, fmt.Sprintf("email LIKE ?"), query.Clause)
			require.Equal(t, []interface{}{sqlLiteral}, query.Args)
		}
	}

	t.Run("Wildcard Only", genTest("*", "%"))
	t.Run("Wildcard Prefix", genTest("*s", "%s"))
	t.Run("Wildcard Suffix", genTest("s*", "s%"))
	t.Run("Wildcard Prefix & Suffix", genTest("*s*", "%s%"))
	t.Run("No Wildcards", genTest("s", "s"))
}
func TestSimpleRecursiveStructure(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `
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

	astNode, err := epsearchast_v3.GetAst(jsonTxt)

	err = json.Unmarshal([]byte(jsonTxt), astNode)
	require.NoError(t, err)

	var sr epsearchast_v3.SemanticReducer[SubQuery] = DefaultGormQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, sr)

	// Verification

	require.NoError(t, err)

	require.Equal(t, "( status IN ? AND amount >= ? )", query.Clause)
	require.Equal(t, []interface{}{[]interface{}{"new", "paid"}, "5"}, query.Args)
}

func TestSimpleRecursiveWithOverrideStructure(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `
				{
					"type":  "AND",
					"children": [
					{
						"type": "IN",
						"args": ["status", "new", "paid"]
					},
					{
						"type": "EQ",
						"args": [ "email",  "ron@swanson.com"]
					}
					]
				}
				`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)

	err = json.Unmarshal([]byte(jsonTxt), astNode)
	require.NoError(t, err)

	var sr epsearchast_v3.SemanticReducer[SubQuery] = &LowerCaseEmail{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, sr)

	// Verification

	require.NoError(t, err)

	require.Equal(t, "( status IN ? AND LOWER(email::text) = LOWER(?) )", query.Clause)
}

type LowerCaseEmail struct {
	DefaultGormQueryBuilder
}

func (l *LowerCaseEmail) VisitEq(first, second string) (*SubQuery, error) {
	if first == "email" {
		return &SubQuery{
			Clause: fmt.Sprintf("LOWER(%s::text) = LOWER(?)", first),
			Args:   []interface{}{second},
		}, nil
	} else {
		return DefaultGormQueryBuilder.VisitEq(l.DefaultGormQueryBuilder, first, second)
	}
}
