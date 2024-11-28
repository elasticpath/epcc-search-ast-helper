package epsearchast_v3_els

import (
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"strings"
)

type DefaultElsQueryBuilder struct {
	OpTypeToFieldNames map[string]*OperatorTypeToMultiFieldName
}

type JsonObject map[string]interface{}

// Elastic Search can encode data in multiple formats using multi fields
// https://www.elastic.co/guide/en/elasticsearch/reference/current/multi-fields.html

type OperatorTypeToMultiFieldName struct {
	// The field name to use for equality operators (eq, in)
	Equality string

	// The field name to use for relational operators (lt, gt, le, ge)
	Relational string

	// The field name to use for text fields (nothing yet)
	Text string

	// The field name for use with array fields (not yet)
	Array string

	// The field name for wild card fields
	Wildcard string

	// The field name for lower case fields
	Lowercase string
}

var _ epsearchast_v3.SemanticReducer[JsonObject] = (*DefaultElsQueryBuilder)(nil)

func (d DefaultElsQueryBuilder) PostVisitAnd(rs []*JsonObject) (*JsonObject, error) {
	return (*JsonObject)(&map[string]interface{}{
		"bool": map[string]interface{}{
			"must": rs,
		},
	}), nil
}

func (d DefaultElsQueryBuilder) VisitIn(args ...string) (*JsonObject, error) {
	return (*JsonObject)(&map[string]interface{}{
		"terms": map[string]interface{}{
			d.getFieldMapping(args[0]).Equality: args[1:],
		},
	}), nil
}

func (d DefaultElsQueryBuilder) VisitEq(first, second string) (*JsonObject, error) {
	return (*JsonObject)(&map[string]interface{}{
		"term": map[string]interface{}{
			d.getFieldMapping(first).Equality: second,
		},
	}), nil
}

func (d DefaultElsQueryBuilder) VisitContains(first, second string) (*JsonObject, error) {
	// Please submit an MR I think this is just a term query, but I didn't have a reference example.
	return nil, fmt.Errorf("contains() is not supported")
}

func (d DefaultElsQueryBuilder) VisitText(first, second string) (*JsonObject, error) {
	return (*JsonObject)(&map[string]interface{}{
		"match_phrase": map[string]interface{}{
			d.getFieldMapping(first).Text: second,
		},
	}), nil
}

// Useful doc: https://www.elastic.co/guide/en/elasticsearch/reference/7.17/query-dsl-range-query.html

func (d DefaultElsQueryBuilder) VisitLe(first, second string) (*JsonObject, error) {
	return (*JsonObject)(&map[string]interface{}{
		"range": map[string]interface{}{
			d.getFieldMapping(first).Relational: map[string]interface{}{
				"lte": second,
			},
		},
	}), nil
}

func (d DefaultElsQueryBuilder) VisitLt(first, second string) (*JsonObject, error) {
	return (*JsonObject)(&map[string]interface{}{
		"range": map[string]interface{}{
			d.getFieldMapping(first).Relational: map[string]interface{}{
				"lt": second,
			},
		},
	}), nil
}

func (d DefaultElsQueryBuilder) VisitGe(first, second string) (*JsonObject, error) {
	return (*JsonObject)(&map[string]interface{}{
		"range": map[string]interface{}{
			d.getFieldMapping(first).Relational: map[string]interface{}{
				"gte": second,
			},
		},
	}), nil
}

func (d DefaultElsQueryBuilder) VisitGt(first, second string) (*JsonObject, error) {
	return (*JsonObject)(&map[string]interface{}{
		"range": map[string]interface{}{
			d.getFieldMapping(first).Relational: map[string]interface{}{
				"gt": second,
			},
		},
	}), nil
}

func (d DefaultElsQueryBuilder) VisitLike(first, second string) (*JsonObject, error) {
	return (*JsonObject)(&map[string]interface{}{
		"wildcard": map[string]interface{}{
			d.getFieldMapping(first).Wildcard: d.EscapeWildcardString(second),
		},
	}), nil
}

func (d DefaultElsQueryBuilder) VisitILike(first, second string) (*JsonObject, error) {
	// Please submit an MR, it should be easy to override, but it depends on how you have mapped the data.
	// I didn't want to implement anything without a reference example
	return nil, fmt.Errorf("ilike() is not supported")
}

func (d DefaultElsQueryBuilder) VisitIsNull(first string) (*JsonObject, error) {
	return (*JsonObject)(&map[string]interface{}{
		"bool": map[string]interface{}{
			"must_not": map[string]interface{}{
				"exists": map[string]interface{}{
					"field": d.getFieldMapping(first).Equality,
				},
			},
		},
	}), nil
}

// getFieldMapping returns the field name to use for a given operator type, the struct is always guaranteed to return f, if nothing was set.
func (d DefaultElsQueryBuilder) getFieldMapping(f string) *OperatorTypeToMultiFieldName {
	var o *OperatorTypeToMultiFieldName

	if d.OpTypeToFieldNames[f] == nil {
		o = &OperatorTypeToMultiFieldName{
			Equality:   f,
			Relational: f,
			Text:       f,
			Array:      f,
			Wildcard:   f,
			Lowercase:  f,
		}
	}

	if v, ok := d.OpTypeToFieldNames[f]; ok {
		o = &OperatorTypeToMultiFieldName{
			Equality:   v.Equality,
			Relational: v.Relational,
			Text:       v.Text,
			Array:      v.Array,
			Wildcard:   v.Wildcard,
			Lowercase:  v.Lowercase,
		}

		if o.Equality == "" {
			o.Equality = f
		}

		if o.Relational == "" {
			o.Relational = f
		}

		if o.Text == "" {
			o.Text = f
		}

		if o.Array == "" {
			o.Array = f
		}

		if o.Wildcard == "" {
			o.Wildcard = f
		}

		if o.Lowercase == "" {
			o.Lowercase = f
		}
	}

	return o
}

func (d DefaultElsQueryBuilder) EscapeWildcardString(s string) string {
	str := strings.ReplaceAll(s, "?", `\?`)
	str = strings.ReplaceAll(str, "*", `\*`)

	if strings.HasPrefix(str, `\*`) {
		str = str[1:]
	}

	if strings.HasSuffix(str, `\*`) {
		str = str[:len(str)-2] + "*"
	}

	return str
}

// Generate an implementation of SemanticReducer[JsonObject] for the Elasticsearch query builder.
