package epsearchast_v3_els_test

import (
	"encoding/json"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	epsearchast_v3_els "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3/els"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleBinaryOperatorGeneratesCorrectQueryWithFieldOveride(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "GT",
  "args": [
    "updated_at",
    "2020-12-25"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "range": {
    "updated_at.date": {
      "gt": "2020-12-25"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_els.JsonObject] = epsearchast_v3_els.DefaultElsQueryBuilder{
		OpTypeToFieldNames: map[string]epsearchast_v3_els.OperatorTypeToMultiFieldName{
			"updated_at": {
				Relational: "updated_at.date",
			},
		},
	}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleBinaryOperatorGeneratesCorrectQuery(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "EQ",
  "args": [
    "email",
    "foo@test.com"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "term": {
    "email": "foo@test.com"
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_els.JsonObject] = epsearchast_v3_els.DefaultElsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleRecursiveStructure(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `
{
  "type": "AND",
  "children": [
    {
      "type": "IN",
      "args": [
        "status",
        "new",
        "paid"
      ]
    },
    {
      "type": "GE",
      "args": [
        "amount",
        "5"
      ]
    }
  ]
}`

	//language=JSON
	expectedJson := `{
  "bool": {
    "must": [
      {
        "terms": {
          "status": [
            "new",
            "paid"
          ]
        }
      },
      {
        "range": {
          "amount": {
            "gte": "5"
          }
        }
      }
    ]
  }
}`
	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_els.JsonObject] = epsearchast_v3_els.DefaultElsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}
