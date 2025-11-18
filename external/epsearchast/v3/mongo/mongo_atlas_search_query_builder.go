package epsearchast_v3_mongo

import (
	"fmt"
	"strings"

	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type DefaultAtlasSearchQueryBuilder struct {
	// Map from field name to multi-analyzer names for different operators
	// If a field is not in this map, or if the analyzer name is "",
	// the base path will be used without specifying a multi-analyzer
	FieldToMultiAnalyzers map[string]*StringMultiAnalyzers
}

type StringMultiAnalyzers struct {
	// Multi-analyzer name for case-insensitive wildcard (ILIKE)
	// If empty, will use: {"path": "field"}
	// If set, will use: {"path": {"value": "field", "multi": "this_value"}}
	WildcardCaseInsensitive string

	// Multi-analyzer name for case-sensitive wildcard (LIKE)
	// If empty, will use: {"path": "field"}
	// If set, will use: {"path": {"value": "field", "multi": "this_value"}}
	WildcardCaseSensitive string
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

func (d DefaultAtlasSearchQueryBuilder) VisitLe(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/range/
	return &bson.D{
		{"range", bson.D{
			{"path", first},
			{"lte", second},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) VisitLt(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/range/
	return &bson.D{
		{"range", bson.D{
			{"path", first},
			{"lt", second},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) VisitGe(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/range/
	return &bson.D{
		{"range", bson.D{
			{"path", first},
			{"gte", second},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) VisitGt(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/range/
	return &bson.D{
		{"range", bson.D{
			{"path", first},
			{"gt", second},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) VisitLike(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/wildcard/
	// Case-sensitive wildcard matching (unlike ILIKE which is case-insensitive)
	path := d.getWildcardPath(first, true)

	return &bson.D{
		{"wildcard", bson.D{
			{"path", path},
			{"query", d.ProcessWildcardString(second)},
			{"allowAnalyzedField", true},
		}},
	}, nil
}

func (d DefaultAtlasSearchQueryBuilder) VisitILike(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/atlas/atlas-search/wildcard/
	// Case-insensitive wildcard matching (uses allowAnalyzedField: true)
	path := d.getWildcardPath(first, false)

	return &bson.D{
		{"wildcard", bson.D{
			{"path", path},
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

// getWildcardPath returns the path configuration for wildcard queries (LIKE/ILIKE)
// If caseSensitive is true, uses WildcardCaseSensitive analyzer
// If caseSensitive is false, uses WildcardCaseInsensitive analyzer
// If no analyzer is configured (or is empty string), returns simple field name
func (d DefaultAtlasSearchQueryBuilder) getWildcardPath(fieldName string, caseSensitive bool) interface{} {
	// Check if field has multi-analyzer configuration
	if config, ok := d.FieldToMultiAnalyzers[fieldName]; ok && config != nil {
		var analyzerName string
		if caseSensitive {
			analyzerName = config.WildcardCaseSensitive
		} else {
			analyzerName = config.WildcardCaseInsensitive
		}

		// If analyzer name is specified, return path with multi
		if analyzerName != "" {
			return bson.D{
				{"value", fieldName},
				{"multi", analyzerName},
			}
		}
	}

	// Otherwise, return simple field name
	return fieldName
}
