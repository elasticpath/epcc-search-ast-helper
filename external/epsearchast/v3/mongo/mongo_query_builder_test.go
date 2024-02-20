package epsearchast_v3_mongo

import (
	"encoding/json"
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
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

			astNode, err := epsearchast_v3.GetAst(astJson)

			var qb epsearchast_v3.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

			expectedSearchJson := fmt.Sprintf(`{"amount":{"%s":"5"}}`, binOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast_v3.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
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

	astNode, err := epsearchast_v3.GetAst(astJson)

	var qb epsearchast_v3.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

	expectedSearchJson := `{"$text":{"$search":"computer"}}`

	// Execute SUT
	queryObj, err := epsearchast_v3.SemanticReduceAst(astNode, qb)

	// Verification

	require.NoError(t, err)

	doc, err := bson.MarshalExtJSON(queryObj, true, false)
	require.NoError(t, err)

	require.Equal(t, expectedSearchJson, string(doc))

}

func TestLikeOperatorFiltersGeneratesCorrectFilter(t *testing.T) {

	//Fixture Setup
	astOp := "LIKE"
	mongoOp := "$regex"
	//language=JSON
	astJson := fmt.Sprintf(`
				{
				"type": "%s",
				"args": [ "amount",  "5"]
			}`, astOp)

	astNode, err := epsearchast_v3.GetAst(astJson)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

	expectedSearchJson := fmt.Sprintf(`{"amount":{"%s":"^5$"}}`, mongoOp)

	// Execute SUT
	queryObj, err := epsearchast_v3.SemanticReduceAst(astNode, qb)

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

			astNode, err := epsearchast_v3.GetAst(astJson)

			var qb epsearchast_v3.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

			expectedSearchJson := fmt.Sprintf(`{"amount":{%s}}`, unaryOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast_v3.SemanticReduceAst(astNode, qb)

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

			astNode, err := epsearchast_v3.GetAst(astJson)

			require.NoError(t, err)

			var qb epsearchast_v3.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

			expectedSearchJson := fmt.Sprintf(`{"amount":{"%s":["5","6","7"]}}`, varOp.MongoOp)

			// Execute SUT
			queryObj, err := epsearchast_v3.SemanticReduceAst(astNode, qb)

			// Verification

			require.NoError(t, err)

			doc, err := bson.MarshalExtJSON(queryObj, true, false)
			require.NoError(t, err)

			require.Equal(t, expectedSearchJson, string(doc))
		})
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

			astNode, err := epsearchast_v3.GetAst(astJson)
			require.NoError(t, err)

			var qb epsearchast_v3.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

			jsonMongoRegexLiteral, err := json.Marshal(mongoRegexLiteral)
			require.NoError(t, err)

			expectedSearchJson := fmt.Sprintf(`{"status":{"%s":%s}}`, mongoOp, jsonMongoRegexLiteral)

			// Execute SUT
			queryObj, err := epsearchast_v3.SemanticReduceAst(astNode, qb)

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

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	var qb epsearchast_v3.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

	// Execute SUT
	queryObj, err := epsearchast_v3.SemanticReduceAst(astNode, qb)

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

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	var qb epsearchast_v3.SemanticReducer[bson.D] = &LowerCaseEmail{}
	queryObj, err := epsearchast_v3.SemanticReduceAst(astNode, qb)

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
