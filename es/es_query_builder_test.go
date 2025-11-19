package astes

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/elasticpath/epcc-search-ast-helper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
			"email": {
				Equality: "email.keyword",
			},
		},
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
			"amount": {
				Relational: "amount.range",
			},
		},
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
			"amount": {
				Relational: "amount.range",
			},
		},
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
			"amount": {
				Relational: "amount.range",
			},
		},
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
			"amount": {
				Relational: "amount.range",
			},
		},
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
			"email": {
				Wildcard: "email.keyword",
			},
		},
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

			astNode, err := epsearchast.GetAst(jsonTxt)
			require.NoError(t, err)

			var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{}

			// Execute SUT
			query, err := epsearchast.SemanticReduceAst(astNode, qb)
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
	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
			"status": {
				Equality: "status.keyword",
			},
		},
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleRecursiveWithStringOverrideStruct(t *testing.T) {
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
						"args": [ "email",  "RON@SWANSON.COM"]
					}
					]
				}
				`

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
        "term": {
          "email": "ron@swanson.com"
        }
      }
    ]
  }
}`

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb = &LowerCaseEmail{
		DefaultEsQueryBuilder{
			OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
				"status": {
					Equality: "status.keyword",
				},
			},
		},
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst[JsonObject](astNode, qb)
	require.NoError(t, err)
	// Verification

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
  "match_bool_prefix": {
    "description": {
      "fuzziness": "0",
      "operator": "and",
      "query": "Cars"
    }
  }
}`

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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
  "match_bool_prefix": {
    "description.text": {
      "fuzziness": "0",
      "operator": "and",
      "query": "Cars"
    }
  }
}`

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
			"description": {
				Text: "description.text",
			},
		},
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestSimpleBinaryTextOperatorGeneratesCorrectQueryWithFuzzinessSetting(t *testing.T) {
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
  "match_bool_prefix": {
    "description": {
      "fuzziness": "AUTO",
      "operator": "and",
      "query": "Cars"
    }
  }
}`

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
		DefaultFuzziness: "AUTO",
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
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

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
		OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
			"sort_order": {
				Equality: "sort_order.keyword",
			},
		},
	}

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst(astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))
}

func TestMustValidateDoesNotPanicOnEmptyObject(t *testing.T) {
	// Fixture Setup
	qb := DefaultEsQueryBuilder{}

	// Execute SUT & Verification
	assert.NotPanics(t, func() {
		qb.MustValidate()
	})
}

func TestMustValidatePanicsWhenRegexDoesNotCompile(t *testing.T) {
	// Fixture Setup
	qb := DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		"tes(": {
			Path: "foo",
		},
	}}

	// Execute SUT & Verification
	assert.PanicsWithValue(t, "regexp: Compile(`tes(`): error parsing regexp: missing closing ): `tes(`", func() {
		qb.MustValidate()
	})

}

func TestMustValidatePanicsWhenRegexDoesNotHaveStartAnchor(t *testing.T) {
	// Fixture Setup
	qb := DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		"test": {
			Path: "foo",
		},
	}}

	// Execute SUT & Verification
	assert.PanicsWithValue(t, "All nested fields must be anchored to the start of the string (e.g., start with a ^), [test] does not", func() {
		qb.MustValidate()
	})

}

func TestMustValidatePanicsWhenRegexDoesNotHaveEndAnchor(t *testing.T) {
	// Fixture Setup
	qb := DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		"^test": {
			Path: "foo",
		},
	}}

	// Execute SUT & Verification
	assert.PanicsWithValue(t, "All nested fields must be anchored at the end of the string (e.g., end in an $), [^test] does not", func() {
		qb.MustValidate()
	})

}

func TestMustValidatePanicsWhenNoPathIsSet(t *testing.T) {
	// Fixture Setup
	qb := DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		"^test.foo$": {},
	}}

	// Execute SUT & Verification
	assert.PanicsWithValue(t, "Path must be set for nested field [^test.foo$]", func() {
		qb.MustValidate()
	})

}

func TestMustValidatePanicsWhenRegexHasValueCaptureGroup(t *testing.T) {
	// Fixture Setup
	qb := DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		"^test.(?P<value>.+)value$": {
			Path: "foo",
		},
	}}

	// Execute SUT & Verification
	assert.PanicsWithValue(t, "Named capture group 'value' is reserved for the replacement value, [^test.(?P<value>.+)value$] cannot use this", func() {
		qb.MustValidate()
	})

}

func TestMustValidatePanicsWhenRegexKeyHasValueReplacement(t *testing.T) {
	// Fixture Setup
	qb := DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		"^test.(?P<id>.+).value$": {
			Path: "foo",
			Subqueries: map[string]Replacement{
				"foo.$value": {
					Value: "$value",
				},
			},
		},
	}}

	// Execute SUT & Verification
	assert.PanicsWithValue(t, "You cannot use $value as replacement in a key in [foo.$value]", func() {
		qb.MustValidate()
	})

}

func TestMustValidatePanicsNoSubqueriesSet(t *testing.T) {
	// Fixture Setup
	qb := DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		"^test.(?P<id>.+).value$": {
			Path:       "foo",
			Subqueries: map[string]Replacement{},
		},
	}}

	// Execute SUT & Verification
	assert.PanicsWithValue(t, "Subqueries must be set for nested field [^test.(?P<id>.+).value$]", func() {
		qb.MustValidate()
	})
}

func TestMustValidatePanicsWhenFieldHasTemplateNotInField(t *testing.T) {
	// Fixture Setup
	qb := DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		"^test.(?P<id>.+).value$": {
			Path: "foo",
			Subqueries: map[string]Replacement{
				"foo.$i.$id": {
					Value: "$id",
				},
			},
		},
	}}

	// Execute SUT & Verification
	assert.PanicsWithValue(t, "Not all templates replaced in nested field [^test.(?P<id>.+).value$] key [foo.$i.$id], after replacement left over with: foo.$i. ", func() {
		qb.MustValidate()
	})
}

func TestMustValidatePanicsWhenFieldValueHasTemplateNotInField(t *testing.T) {
	// Fixture Setup
	qb := DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		"^test.(?P<id>.+).value.(?P<bar>.+)$": {
			Path: "foo",
			Subqueries: map[string]Replacement{
				"foo.$id": {
					Value: "$yas.$bar",
				},
			},
		},
	}}

	// Execute SUT & Verification
	assert.PanicsWithValue(t, "Not all templates replaced in nested field [^test.(?P<id>.+).value.(?P<bar>.+)$] key [foo.$id] with value [$yas.$bar], after replacement left over with: $yas.", func() {
		qb.MustValidate()
	})
}

func TestMultipleRegexMatchesAreReplacedCorrectlyWhenCaptureGroupsOverlap(t *testing.T) {
	// We are kind of sloppy with how we do templates, replacing $key with value by string.Replace.
	// If two capture groups have a prefix then if you aren't careful you can get into weird states.
	// For instance "$user and $username" if you have templates "user=foo, and username=bar", then
	// "$user and $username" might get replaced with "foo and fooname" instead of "foo and bar".
	// We could use a non-alphanumeric suffix but that might be limiting.

	//language=JSON
	jsonTxt := `{
  "type": "EQ",
  "args": [
    "field.X.Y.Z","123"
  ]
}
`

	//language=JSON
	// You can see an example of what this tests, by disabling the sorting of keys in applyPatternGroupsToFieldNameAndValue
	// You'll get some left over as in the response.
	expectedJson := `{
  "nested": {
    "path": "foo",
    "query": {
      "bool": {
        "must": [
          {
            "term": {
              "foo.Y.Y.Y": "ZYX"
            }
          },
          {
            "term": {
              "foo.XZY": "YZX"
            }
          }
        ]
      }
    }
  }
}`

	qb := &DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		`^field\.(?P<aaa>[^.]+)\.(?P<a>[^.]+)\.(?P<aa>[^.]+)$`: {
			Path: "foo",
			Subqueries: map[string]Replacement{
				"foo.$aaa$aa$a": {
					Value: "$a$aa$aaa",
				},
				"foo.$a.$a.$a": {
					Value: "$aa$a$aaa",
				},
			},
		},
	}}

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst[JsonObject](astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))

}

func TestMultipleReplacementFieldsInNestedObjectAreSortedInDeterministicOrder(t *testing.T) {
	// For consumers who want to unit test things, it's nice if things are deterministic.

	//language=JSON
	jsonTxt := `{
  "type": "EQ",
  "args": [
    "field.X","123"
  ]
}
`

	//language=JSON
	// You can see an example of what this tests, by disabling the sorting of keys in processNestedFieldToQuery
	expectedJson := `{
  "nested": {
    "path": "foo",
    "query": {
      "bool": {
        "must": [
          {
            "term": {
              "foo.a": "123"
            }
          },
          {
            "term": {
              "foo.b": "123"
            }
          },
          {
            "term": {
              "foo.c": "123"
            }
          },
          {
            "term": {
              "foo.d": "123"
            }
          },
          {
            "term": {
              "foo.e": "123"
            }
          },
          {
            "term": {
              "foo.f": "123"
            }
          },
          {
            "term": {
              "foo.g": "123"
            }
          },
          {
            "term": {
              "foo.h": "123"
            }
          },
          {
            "term": {
              "foo.i": "123"
            }
          },
          {
            "term": {
              "foo.j": "123"
            }
          },
          {
            "term": {
              "foo.k": "123"
            }
          },
          {
            "term": {
              "foo.l": "123"
            }
          },
          {
            "term": {
              "foo.m": "123"
            }
          }
        ]
      }
    }
  }
}`

	qb := &DefaultEsQueryBuilder{NestedFieldToQuery: map[string]NestedReplacement{
		`^field\.(?P<aaa>[^.]+)$`: {
			Path: "foo",
			Subqueries: map[string]Replacement{
				"foo.j": {Value: "$value"},
				"foo.i": {Value: "$value"},
				"foo.l": {Value: "$value"},
				"foo.m": {Value: "$value"},
				"foo.h": {Value: "$value"},
				"foo.c": {Value: "$value"},
				"foo.f": {Value: "$value"},
				"foo.k": {Value: "$value"},
				"foo.a": {Value: "$value"},
				"foo.b": {Value: "$value"},
				"foo.e": {Value: "$value"},
				"foo.d": {Value: "$value"},
				"foo.g": {Value: "$value"},
			},
		},
	}}

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	query, err := epsearchast.SemanticReduceAst[JsonObject](astNode, qb)
	require.NoError(t, err)

	// Verification
	queryJson, err := json.MarshalIndent(query, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedJson, string(queryJson))

}

type LowerCaseEmail struct {
	DefaultEsQueryBuilder
}

func (l *LowerCaseEmail) VisitEq(first, second string) (*JsonObject, error) {
	if first == "email" {
		return DefaultEsQueryBuilder.VisitEq(l.DefaultEsQueryBuilder, first, strings.ToLower(second))
	} else {
		return DefaultEsQueryBuilder.VisitEq(l.DefaultEsQueryBuilder, first, second)
	}
}
