package astmongo

import (
	"encoding/json"
	"fmt"
	"github.com/elasticpath/epcc-search-ast-helper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"testing"
)

var binOps = []testOp{
	{"LE", "$lte"},
	{"LT", "$lt"},
	{"EQ", "$eq"},
	{"GT", "$gt"},
	{"GE", "$gte"},

	// The regex conversion is to complex, so we test that method distinctly.
	//{"LIKE", "$regex"},
	// Same with contains
}

var unaryOps = []testOp{
	{"IS_NULL", `"$eq":null`},
}

var varOps = []testOp{
	{"IN", "$in"},
}

type testOp struct {
	AstOp   string
	MongoOp string
}

func TestSimpleBinaryOperatorFiltersGeneratesCorrectFilter(t *testing.T) {
	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp.AstOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount",  "5"]
			}`, binOp.AstOp)

			astNode, err := epsearchast.GetAst(astJson)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

			expectedSearchJson := fmt.Sprintf(`{"amount":{"%s":"5"}}`, binOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
	}
}

func TestSimpleBinaryOperatorFiltersGeneratesCorrectFilterWithInt64TypeConversion(t *testing.T) {
	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp.AstOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount",  "5"]
			}`, binOp.AstOp)

			astNode, err := epsearchast.GetAst(astJson)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{"amount": epsearchast.Int64}}

			// https://www.mongodb.com/docs/manual/reference/mongodb-extended-json/#mongodb-bsontype-Int64
			expectedSearchJson := fmt.Sprintf(`{"amount":{"%s":{"$numberLong":"5"}}}`, binOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
	}
}

func TestSimpleBinaryOperatorFiltersGeneratesCorrectFilterWithFloat64TypeConversion(t *testing.T) {
	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp.AstOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount",  "5"]
			}`, binOp.AstOp)

			astNode, err := epsearchast.GetAst(astJson)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{"amount": epsearchast.Float64}}

			// https://www.mongodb.com/docs/manual/reference/mongodb-extended-json/#mongodb-bsontype-Int64
			expectedSearchJson := fmt.Sprintf(`{"amount":{"%s":{"$numberDouble":"5.0"}}}`, binOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
	}
}

func TestSimpleBinaryOperatorFiltersGeneratesCorrectFilterWithBooleanTypeConversion(t *testing.T) {
	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp.AstOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "paid",  "true"]
			}`, binOp.AstOp)

			astNode, err := epsearchast.GetAst(astJson)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{"paid": epsearchast.Boolean}}

			// https://www.mongodb.com/docs/manual/reference/mongodb-extended-json/#mongodb-bsontype-Int64
			expectedSearchJson := fmt.Sprintf(`{"paid":{"%s":true}}`, binOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
	}
}

func TestSimpleBinaryOperatorFiltersGeneratesErrorWhenValueCantBeConverted(t *testing.T) {
	for _, fieldType := range []epsearchast.FieldType{epsearchast.Int64, epsearchast.Float64, epsearchast.Boolean} {
		for _, binOp := range binOps {
			t.Run(fmt.Sprintf("%s %s", fieldType, binOp.AstOp), func(t *testing.T) {
				//Fixture Setup
				//language=JSON
				astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "foo",  "Hello World!"]
			}`, binOp.AstOp)

				astNode, err := epsearchast.GetAst(astJson)

				var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{"foo": fieldType}}

				// Execute SUT
				_, err = epsearchast.SemanticReduceAst(astNode, qb)

				// Verification
				errStr := fmt.Sprintf("invalid value for %s", fieldType)
				require.ErrorContains(t, err, errStr)
			})
		}
	}
}

func TestTextBinaryOperatorFiltersGeneratesCorrectFilter(t *testing.T) {
	//Fixture Setup
	//language=JSON
	astJson := fmt.Sprintf(`
		{
		"type": "%s",
		"args": [ "*",  "computer"]
	}`, "TEXT")

	astNode, err := epsearchast.GetAst(astJson)

	var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

	expectedSearchJson := `{"$text":{"$search":"computer"}}`

	// Execute SUT
	queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

	// Verification

	require.NoError(t, err)

	doc, err := bson.MarshalExtJSON(queryObj, true, false)
	require.NoError(t, err)

	require.Equal(t, expectedSearchJson, string(doc))

}

func TestTextBinaryOperatorFiltersGeneratesErrorWhenNotAStringType(t *testing.T) {
	//Fixture Setup
	//language=JSON
	astJson := fmt.Sprintf(`
		{
		"type": "%s",
		"args": [ "*",  "computer"]
	}`, "TEXT")

	astNode, err := epsearchast.GetAst(astJson)

	var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{
		"*": epsearchast.Int64,
	}}

	// Execute SUT
	_, err = epsearchast.SemanticReduceAst(astNode, qb)

	// Verification
	require.ErrorContains(t, err, "text() operator is only supported for string fields")
	require.ErrorContains(t, err, "[*] is not a string")
}

func TestLikeBinaryOperatorFiltersGeneratesErrorWhenNotAStringType(t *testing.T) {
	//Fixture Setup
	//language=JSON
	astJson := fmt.Sprintf(`
		{
		"type": "%s",
		"args": [ "foo",  "52"]
	}`, "LIKE")

	astNode, err := epsearchast.GetAst(astJson)

	var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{
		"foo": epsearchast.Int64,
	}}

	// Execute SUT
	_, err = epsearchast.SemanticReduceAst(astNode, qb)

	// Verification
	require.ErrorContains(t, err, "like() operator is only supported for string fields")
	require.ErrorContains(t, err, "[foo] is not a string")
}

func TestILikeBinaryOperatorFiltersGeneratesErrorWhenNotAStringType(t *testing.T) {
	//Fixture Setup
	//language=JSON
	astJson := fmt.Sprintf(`
		{
		"type": "%s",
		"args": [ "foo",  "52"]
	}`, "ILIKE")

	astNode, err := epsearchast.GetAst(astJson)

	var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{
		"foo": epsearchast.Int64,
	}}

	// Execute SUT
	_, err = epsearchast.SemanticReduceAst(astNode, qb)

	// Verification
	require.ErrorContains(t, err, "like() operator is only supported for string fields")
	require.ErrorContains(t, err, "[foo] is not a string")
}

func TestILikeOperatorFiltersGeneratesCorrectFilter(t *testing.T) {

	//Fixture Setup
	astOp := "ILIKE"
	mongoOp := "$regex"
	//language=JSON
	astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount",  "5"]
			}`, astOp)

	astNode, err := epsearchast.GetAst(astJson)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

	expectedSearchJson := fmt.Sprintf(`{"amount":{"%s":"^5$","$options":"i"}}`, mongoOp)

	// Execute SUT
	queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

	// Verification

	require.NoError(t, err)

	doc, err := bson.MarshalExtJSON(queryObj, true, false)
	require.NoError(t, err)

	require.Equal(t, expectedSearchJson, string(doc))

}

func TestContainsOperatorFiltersGeneratesCorrectFilter(t *testing.T) {

	//Fixture Setup
	//language=JSON
	astJson := `{
				"type": "CONTAINS",
				"args": [ "favourite_colors",  "red"]
			}`

	astNode, err := epsearchast.GetAst(astJson)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

	expectedSearchJson :=
		`{"favourite_colors":{"$elemMatch":{"$eq":"red"}}}`

	// Execute SUT
	queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

	// Verification

	require.NoError(t, err)

	doc, err := bson.MarshalExtJSON(queryObj, true, false)
	require.NoError(t, err)

	require.Equal(t, expectedSearchJson, string(doc))
}

func TestSimpleUnaryOperatorFiltersGeneratesCorrectFilter(t *testing.T) {
	for _, unaryOp := range unaryOps {
		t.Run(fmt.Sprintf("%s", unaryOp.AstOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount"]
			}`, unaryOp.AstOp)

			astNode, err := epsearchast.GetAst(astJson)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

			expectedSearchJson := fmt.Sprintf(`{"amount":{%s}}`, unaryOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
	}
}

func TestSimpleVariableOperatorFiltersGeneratesCorrectFilter(t *testing.T) {
	for _, varOp := range varOps {
		t.Run(fmt.Sprintf("%s", varOp.AstOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount",  "5", "6", "7"]
			}`, varOp.AstOp)

			astNode, err := epsearchast.GetAst(astJson)

			require.NoError(t, err)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

			expectedSearchJson := fmt.Sprintf(`{"amount":{"%s":["5","6","7"]}}`, varOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
	}
}

func TestSimpleVariableOperatorFiltersGeneratesCorrectFilterWithInt64(t *testing.T) {
	for _, varOp := range varOps {
		t.Run(fmt.Sprintf("%s", varOp.AstOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount",  "5", "6", "7"]
			}`, varOp.AstOp)

			astNode, err := epsearchast.GetAst(astJson)

			require.NoError(t, err)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{"amount": epsearchast.Int64}}

			expectedSearchJson := fmt.Sprintf(`{"amount":{"%s":[{"$numberLong":"5"},{"$numberLong":"6"},{"$numberLong":"7"}]}}`, varOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
	}
}

func TestSimpleVariableOperatorFiltersGeneratesCorrectFilterWithFloat64(t *testing.T) {
	for _, varOp := range varOps {
		t.Run(fmt.Sprintf("%s", varOp.AstOp), func(t *testing.T) {
			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount",  "5", "6", "7"]
			}`, varOp.AstOp)

			astNode, err := epsearchast.GetAst(astJson)

			require.NoError(t, err)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{"amount": epsearchast.Float64}}

			expectedSearchJson := fmt.Sprintf(`{"amount":{"%s":[{"$numberDouble":"5.0"},{"$numberDouble":"6.0"},{"$numberDouble":"7.0"}]}}`, varOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
	}
}

func TestSimpleVariableOperatorFiltersGeneratesCorrectFilterWithBoolean(t *testing.T) {
	for _, varOp := range varOps {
		t.Run(fmt.Sprintf("%s", varOp.AstOp), func(t *testing.T) {
			// Yes this test is kind of silly the query is stupid,
			// but we should support it.

			//Fixture Setup
			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "paid",  "true", "false", "true", "true", "false"]
			}`, varOp.AstOp)

			astNode, err := epsearchast.GetAst(astJson)

			require.NoError(t, err)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{"paid": epsearchast.Boolean}}

			expectedSearchJson := fmt.Sprintf(`{"paid":{"%s":[true,false,true,true,false]}}`, varOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
	}
}

func TestSimpleVariableOperatorFiltersGeneratesErrorIfInvalidValue(t *testing.T) {
	for _, fieldType := range []epsearchast.FieldType{epsearchast.Int64, epsearchast.Float64, epsearchast.Boolean} {
		for _, varOp := range varOps {
			t.Run(fmt.Sprintf("%s %s", fieldType, varOp.AstOp), func(t *testing.T) {
				// Yes also this test case is kind of silly, for booleans.

				//Fixture Setup
				//language=JSON
				astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount", "Nothing!", "5", "7"]
			}`, varOp.AstOp)

				astNode, err := epsearchast.GetAst(astJson)

				require.NoError(t, err)

				var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{FieldTypes: map[string]epsearchast.FieldType{"amount": fieldType}}

				// Execute SUT
				_, err = epsearchast.SemanticReduceAst(astNode, qb)

				// Verification

				require.ErrorContains(t, err, "could not validate position")
				require.ErrorContains(t, err, "Nothing!")
				require.ErrorContains(t, err, fieldType.String())
			})
		}
	}
}

func TestLikeFilterWildCards(t *testing.T) {
	astOp := "LIKE"
	mongoOp := "$regex"

	genTest := func(astLiteral string, mongoRegexLiteral string) func(t *testing.T) {
		return func(t *testing.T) {

			//Fixture Setup

			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "status",  "%s"]
			}`, astOp, astLiteral)

			astNode, err := epsearchast.GetAst(astJson)
			require.NoError(t, err)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

			jsonMongoRegexLiteral, err := json.Marshal(mongoRegexLiteral)
			require.NoError(t, err)

			expectedSearchJson := fmt.Sprintf(`{"status":{"%s":%s}}`, mongoOp, jsonMongoRegexLiteral)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))

		}
	}

	t.Run("Wildcard Only", genTest("*", "^.*$"))
	t.Run("Wildcard Prefix", genTest("*aid", "^.*aid$"))
	t.Run("Wildcard Suffix", genTest("pai*", "^pai.*$"))
	t.Run("Wildcard Prefix & Suffix", genTest("*ai*", "^.*ai.*$"))
	t.Run("No Wildcards", genTest("paid", "^paid$"))
	t.Run("Middle wildcards escaped", genTest("p*d", `^p\*d$`))
	t.Run("Only Middle wildcards escaped", genTest("*p*d*", `^.*p\*d.*$`))
	t.Run("Middle dot escaped", genTest("p..d", `^p\.\.d$`))
}

func TestILikeFilterWildCards(t *testing.T) {
	astOp := "ILIKE"
	mongoOp := "$regex"

	genTest := func(astLiteral string, mongoRegexLiteral string) func(t *testing.T) {
		return func(t *testing.T) {

			//Fixture Setup

			//language=JSON
			astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "status",  "%s"]
			}`, astOp, astLiteral)

			astNode, err := epsearchast.GetAst(astJson)
			require.NoError(t, err)

			var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

			jsonMongoRegexLiteral, err := json.Marshal(mongoRegexLiteral)
			require.NoError(t, err)

			expectedSearchJson := fmt.Sprintf(`{"status":{"%s":%s,"$options":"i"}}`, mongoOp, jsonMongoRegexLiteral)

			// Execute SUT
			queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))

		}
	}

	t.Run("Wildcard Only", genTest("*", "^.*$"))
	t.Run("Wildcard Prefix", genTest("*aid", "^.*aid$"))
	t.Run("Wildcard Suffix", genTest("pai*", "^pai.*$"))
	t.Run("Wildcard Prefix & Suffix", genTest("*ai*", "^.*ai.*$"))
	t.Run("No Wildcards", genTest("paid", "^paid$"))
	t.Run("Middle wildcards escaped", genTest("p*d", `^p\*d$`))
	t.Run("Only Middle wildcards escaped", genTest("*p*d*", `^.*p\*d.*$`))
	t.Run("Middle dot escaped", genTest("p..d", `^p\.\.d$`))
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

	expectedMongoJSON := strings.Trim(
		//language=JSON
		`
{
  "$and": [
    {
      "status": {
        "$in": [
          "new",
          "paid"
        ]
      }
    },
    {
      "amount": {
        "$gte": "5"
      }
    }
  ]
}
`, "\n ")

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

	// Execute SUT
	queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

	// Verification

	require.NoError(t, err)

	doc, err := bson.MarshalExtJSONIndent(queryObj, true, false, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedMongoJSON, string(doc))

}

func TestSimpleRecursiveStructureWithOverrideStruct(t *testing.T) {
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
						"args": [ "email",  "Ron@Swanson.com"]
					}
					]
				}
				`

	expectedMongoJSON := strings.Trim(
		//language=JSON
		`
{
  "$and": [
    {
      "status": {
        "$in": [
          "new",
          "paid"
        ]
      }
    },
    {
      "email": {
        "$eq": "ron@swanson.com"
      }
    }
  ]
}
`, "\n ")

	astNode, err := epsearchast.GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	var qb epsearchast.SemanticReducer[bson.D] = &LowerCaseEmail{}
	queryObj, err := epsearchast.SemanticReduceAst(astNode, qb)

	// Verification

	require.NoError(t, err)

	doc, err := bson.MarshalExtJSONIndent(queryObj, true, false, "", "  ")
	require.NoError(t, err)

	require.Equal(t, expectedMongoJSON, string(doc))

}

type LowerCaseEmail struct {
	DefaultMongoQueryBuilder
}

func (l *LowerCaseEmail) VisitEq(first, second string) (*bson.D, error) {
	if first == "email" {
		return &bson.D{{first, bson.D{{"$eq", strings.ToLower(second)}}}}, nil
	} else {
		return DefaultMongoQueryBuilder.VisitEq(l.DefaultMongoQueryBuilder, first, second)
	}
}
