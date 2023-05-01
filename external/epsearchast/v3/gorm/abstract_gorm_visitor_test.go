package v3_gorm

import (
	"encoding/json"
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

	// Execute SUT
	query, err := epsearchast_v3.ReduceAst(astNode, Apply)

	// Verification

	require.NoError(t, err)

	require.Equal(t, "status IN ? AND amount >= ?", query.Clause)
	require.Equal(t, []interface{}{[]interface{}{"new", "paid"}, "5"}, query.Args)
}
