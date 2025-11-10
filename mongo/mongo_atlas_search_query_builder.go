package epsearchast_v3_mongo

import (
	"fmt"
	"strings"

	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type DefaultAtlasSearchQueryBuilder struct {
}

var _ epsearchast_v3.SemanticReducer[bson.D] = (*DefaultAtlasSearchQueryBuilder)(nil)

func (d DefaultAtlasSearchQueryBuilder) PostVisitAnd(rs []*bson.D) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/compound/
	return &bson.D{
		{"compound", bson.D{
			{"must", rs},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) PostVisitOr(rs []*bson.D) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/compound/
	return &bson.D{
		{"compound", bson.D{
			{"should", rs},
			// minimumShouldMatch: 1 means at least one should clause must match
			// https://www.mongodb.com/docs/atlas/atlas-search/compound/#std-label-compound-ref
			{"minimumShouldMatch", 1},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) VisitText(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/text/
	return &bson.D{
		{"text", bson.D{
			{"query", second},
			{"path", first},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) VisitIn(args ...string) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/in/
	if len(args) < 2 {
		return nil, fmt.Errorf("IN operator requires at least 2 arguments (field and at least one value)")
	}

	fieldName := args[0]
	values := args[1:]

	return &bson.D{
		{"in", bson.D{
			{"path", fieldName},
			{"value", values},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) VisitEq(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/equals/
	return &bson.D{
		{"equals", bson.D{
			{"path", first},
			{"value", second},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) VisitLe(_, _ string) (*bson.D, error) {
	return nil, fmt.Errorf("LE operator not yet implemented for Atlas Search")
}

func (d DefaultAtlasSearchQueryBuilder) VisitLt(_, _ string) (*bson.D, error) {
	return nil, fmt.Errorf("LT operator not yet implemented for Atlas Search")
}

func (d DefaultAtlasSearchQueryBuilder) VisitGe(_, _ string) (*bson.D, error) {
	return nil, fmt.Errorf("GE operator not yet implemented for Atlas Search")
}

func (d DefaultAtlasSearchQueryBuilder) VisitGt(_, _ string) (*bson.D, error) {
	return nil, fmt.Errorf("GT operator not yet implemented for Atlas Search")
}

func (d DefaultAtlasSearchQueryBuilder) VisitLike(_, _ string) (*bson.D, error) {
	return nil, fmt.Errorf("LIKE operator not yet implemented for Atlas Search")
}

func (d DefaultAtlasSearchQueryBuilder) VisitILike(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/wildcard/
	// allowAnalyzedField: true makes the search case-insensitive
	return &bson.D{
		{"wildcard", bson.D{
			{"path", first},
			{"query", d.ProcessWildcardString(second)},
			{"allowAnalyzedField", true},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) VisitContains(_, _ string) (*bson.D, error) {
	return nil, fmt.Errorf("CONTAINS operator not yet implemented for Atlas Search")
}

func (d DefaultAtlasSearchQueryBuilder) VisitContainsAny(_ ...string) (*bson.D, error) {
	return nil, fmt.Errorf("CONTAINS_ANY operator not yet implemented for Atlas Search")
}

func (d DefaultAtlasSearchQueryBuilder) VisitContainsAll(_ ...string) (*bson.D, error) {
	return nil, fmt.Errorf("CONTAINS_ALL operator not yet implemented for Atlas Search")
}

func (d DefaultAtlasSearchQueryBuilder) VisitIsNull(_ string) (*bson.D, error) {
	return nil, fmt.Errorf("IS_NULL operator not yet implemented for Atlas Search")
}

// ProcessWildcardString processes wildcard strings for Atlas Search wildcard queries
// Escapes special characters except * and ? at the beginning/end
func (d DefaultAtlasSearchQueryBuilder) ProcessWildcardString(s string) string {
	// Atlas Search wildcard uses * and ? as wildcards, similar to ES
	// Escape all wildcards first
	str := strings.ReplaceAll(s, "?", `\?`)
	str = strings.ReplaceAll(str, "*", `\*`)

	// Un-escape wildcards at the beginning
	if strings.HasPrefix(str, `\*`) {
		str = str[1:]
	}

	// Un-escape wildcards at the end
	if strings.HasSuffix(str, `\*`) {
		str = str[:len(str)-2] + "*"
	}

	return str
}
