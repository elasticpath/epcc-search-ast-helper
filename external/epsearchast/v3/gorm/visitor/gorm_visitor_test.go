package v3_gorm_visitor

import (
	"encoding/json"
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

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
	query, err := epsearchast_v3.SemanticReduceAst(astNode, &sr)

	// Verification

	require.NoError(t, err)

	require.Equal(t, "status IN ? AND amount >= ?", query.Clause)
	//require.Equal(t, []interface{}{[]interface{}{"new", "paid"}, "5"}, query.Args)
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
						"args": [ "email",  "Steve.Ramage@elasticpath.com"]
					}
					]
				}
				`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)

	err = json.Unmarshal([]byte(jsonTxt), astNode)
	require.NoError(t, err)

	var sr epsearchast_v3.SemanticReducer[SubQuery] = &LowerCaseEmail{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, &sr)

	// Verification

	require.NoError(t, err)

	require.Equal(t, "status IN ? AND LOWER(email::text) = LOWER(?)", query.Clause)
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
