package epsearchast_v3_es_test

import (
	"encoding/json"
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	epsearchast_v3_es "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3/es"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestSimpleBinaryEqOperatorGeneratesCorrectQuery(t *testing.T) {
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

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleBinaryEqOperatorGeneratesCorrectQueryWithFieldOverride(t *testing.T) {
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
    "email.keyword": "foo@test.com"
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*epsearchast_v3_es.OperatorTypeToMultiFieldName{
			"email": {
				Equality: "email.keyword",
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

func TestSimpleBinaryLeOperatorGeneratesCorrectQuery(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "LE",
  "args": [
    "amount",
    "5"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "range": {
    "amount": {
      "lte": "5"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleBinaryLeOperatorGeneratesCorrectQueryWithFieldOverride(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "LE",
  "args": [
    "amount",
    "5"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "range": {
    "amount.range": {
      "lte": "5"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*epsearchast_v3_es.OperatorTypeToMultiFieldName{
			"amount": {
				Relational: "amount.range",
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

func TestSimpleBinaryLtOperatorGeneratesCorrectQuery(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "LT",
  "args": [
    "amount",
    "5"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "range": {
    "amount": {
      "lt": "5"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleBinaryLtOperatorGeneratesCorrectQueryWithFieldOverride(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "LT",
  "args": [
    "amount",
    "5"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "range": {
    "amount.range": {
      "lt": "5"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*epsearchast_v3_es.OperatorTypeToMultiFieldName{
			"amount": {
				Relational: "amount.range",
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

func TestSimpleBinaryGtOperatorGeneratesCorrectQuery(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "GT",
  "args": [
    "amount",
    "5"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "range": {
    "amount": {
      "gt": "5"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleBinaryGtOperatorGeneratesCorrectQueryWithFieldOverride(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "GT",
  "args": [
    "amount",
    "5"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "range": {
    "amount.range": {
      "gt": "5"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*epsearchast_v3_es.OperatorTypeToMultiFieldName{
			"amount": {
				Relational: "amount.range",
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

func TestSimpleBinaryGEOperatorGeneratesCorrectQuery(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "GE",
  "args": [
    "amount",
    "5"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "range": {
    "amount": {
      "gte": "5"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleBinaryGEOperatorGeneratesCorrectQueryWithFieldOverride(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "GE",
  "args": [
    "amount",
    "5"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "range": {
    "amount.range": {
      "gte": "5"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*epsearchast_v3_es.OperatorTypeToMultiFieldName{
			"amount": {
				Relational: "amount.range",
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

func TestSimpleBinaryLikeOperatorGeneratesCorrectQuery(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "LIKE",
  "args": [
    "email",
    "@test.com"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "wildcard": {
    "email": {
      "case_insensitive": false,
      "value": "@test.com"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleBinaryLikeOperatorGeneratesCorrectQueryWithFieldOverride(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "LIKE",
  "args": [
    "email",
    "@test.com"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "wildcard": {
    "email.keyword": {
      "case_insensitive": false,
      "value": "@test.com"
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*epsearchast_v3_es.OperatorTypeToMultiFieldName{
			"email": {
				Wildcard: "email.keyword",
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

func TestSimpleBinaryLikeOperatorGeneratesCorrectQueryWithWildcards(t *testing.T) {

	for _, tc := range []struct {
		searchTerm           string
		expectedWildcardTerm string
	}{
		{"@test.com", `@test.com`},
		{"*@test.com", `*@test.com`},
		{"*@*.com", `*@\*.com`},
		{"*@*", `*@*`},
		{"*", `*`},
		{"**", `**`},
		{"?", `\?`},
		{"user@*", `user@*`},
		{"user??@*", `user\?\?@*`},
		{"@**", `@\**`},
	} {
		t.Run(tc.searchTerm, func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			jsonTxt := fmt.Sprintf(`{
  "type": "LIKE",
  "args": [
    "email",
    "%s"
  ]
}
`, tc.searchTerm)

			//language=JSON
			expectedJson := fmt.Sprintf(`{
  "wildcard": {
    "email": {
      "case_insensitive": false,
      "value": "%s"
    }
  }
}`, strings.ReplaceAll(tc.expectedWildcardTerm, `\`, `\\`))

			astNode, err := epsearchast_v3.GetAst(jsonTxt)
			require.NoError(t, err)

			var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{}

			// Execute SUT
			query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
			require.NoError(t, err)

			// Verification
			queryJson, err := json.MarshalIndent(query, "", "  ")
			require.NoError(t, err)

			require.Equal(t, expectedJson, string(queryJson))
		})
	}

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
          "status.keyword": [
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

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*epsearchast_v3_es.OperatorTypeToMultiFieldName{
			"status": {
				Equality: "status.keyword",
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

func TestSimpleBinaryTextOperatorGeneratesCorrectQuery(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "TEXT",
  "args": [
    "description",
    "Cars"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "match": {
    "description": "Cars"
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleBinaryTextOperatorGeneratesCorrectQueryWithFieldOverride(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "TEXT",
  "args": [
    "description",
    "Cars"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "match": {
    "description.text": "Cars"
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*epsearchast_v3_es.OperatorTypeToMultiFieldName{
			"description": {
				Text: "description.text",
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

func TestSimpleUnaryIsNullOperatorGeneratesCorrectQuery(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "IS_NULL",
  "args": [
    "sort_order"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "bool": {
    "must_not": {
      "exists": {
        "field": "sort_order"
      }
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast_v3.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleUnaryIsNullOperatorGeneratesCorrectQueryWithFieldOverride(t *testing.T) {
	//Fixture Setup
	//language=JSON
	jsonTxt := `{
  "type": "IS_NULL",
  "args": [
    "sort_order"
  ]
}
`

	//language=JSON
	expectedJson := `{
  "bool": {
    "must_not": {
      "exists": {
        "field": "sort_order.keyword"
      }
    }
  }
}`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_es.JsonObject] = epsearchast_v3_es.DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*epsearchast_v3_es.OperatorTypeToMultiFieldName{
			"sort_order": {
				Equality: "sort_order.keyword",
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
