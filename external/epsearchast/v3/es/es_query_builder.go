package astes

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
)

type JsonObject map[string]any

type DefaultEsQueryBuilder struct {
	// OpenSearch supports Multi-Fields (https://opensearch.org/docs/latest/field-types/supported-field-types/index/#multifields) which allows a single field to be encoded in different ways
	// If you have multiple mappings for a field, this can let that field be used in range, keyword, or text queries.
	// The keys here should be field names from the filter (after processing from NestedFieldToQuery), and for each type of filter query the resulting filter to use.
	OpTypeToFieldNames map[string]*OperatorTypeToMultiFieldName

	// https://opensearch.org/docs/latest/field-types/supported-field-types/nested/
	// https://opensearch.org/docs/latest/query-dsl/joining/nested/
	// NestedFieldToQuery is a keyed map that takes as a key a regular expression for an attribute that we should match (e.g., requested by the user, after aliases have been processed).
	// The value is information about how to replace it, and allows us to create a nested query (https://opensearch.org/docs/latest/query-dsl/joining/nested/) that contains the path nested.
	// The regular expression can have capture groups that will be used as replacements in the subquery keys and values.
	NestedFieldToQuery map[string]NestedReplacement

	// The default value for fuzziness
	// https://opensearch.org/docs/latest/query-dsl/term/fuzzy/
	// Default value is treated as zero
	DefaultFuzziness string
}

type NestedReplacement struct {
	// The path that will be used in the nested argument (See: https://opensearch.org/docs/latest/query-dsl/joining/nested/#parameters)
	Path string

	// A map which generates the set of subqueries queries that should be generated.
	// Named capture groups in the parent map will be replaced (e.g., a field ^foo\[(?P<id>\d+)\].bar$) can use $id as a replacement in this string.
	Subqueries map[string]Replacement
}

type Replacement struct {
	// The value we should search for, we can use the named capture groups from the parent regex as replacements, also the special value $value is available
	Value string

	// By default, we will use the existing search term as a replacement, if set to true, we will generate an equality match.
	ForceEQ bool
}

// Elasticsearch can encode data in multiple formats using multi fields
// https://www.elastic.co/guide/en/elasticsearch/reference/current/multi-fields.html

type OperatorTypeToMultiFieldName struct {
	// The field name to use for equality operators (eq, in)
	Equality string

	// The field name to use for relational operators (lt, gt, le, ge)
	Relational string

	// The field name to use for text fields (nothing yet)
	Text string

	// The field name for use with array fields
	Array string

	// The field name for wild card fields
	Wildcard string
}

// MustValidate will ensure that the configuration of the query builder is correct and if not, panics. It simplifies safe initialization of the variable.
func (d DefaultEsQueryBuilder) MustValidate() {
	for k := range d.NestedFieldToQuery {
		re := regexp.MustCompile(k)

		if k[0] != '^' {
			panic(fmt.Sprintf("All nested fields must be anchored to the start of the string (e.g., start with a ^), [%s] does not", k))
		}

		if k[len(k)-1] != '$' {
			panic(fmt.Sprintf("All nested fields must be anchored at the end of the string (e.g., end in an $), [%s] does not", k))
		}

		var groupKeys []string
		for _, v := range re.SubexpNames() {
			if v == "value" {
				panic(fmt.Sprintf("Named capture group 'value' is reserved for the replacement value, [%s] cannot use this", k))
			}
			groupKeys = append(groupKeys, v)
		}

		// We need to resolve keys in decreasing order of length
		// So that if we substitute templates with their replacement in consistent order.
		// E.g., if you have templates $user=foo and $username=bar, "$user and $username" needs to resolve to
		// "foo and bar" not "foo and fooname", which if you replace the string user first, is what you get.
		sortByDecreasingLength(groupKeys)

		if d.NestedFieldToQuery[k].Path == "" {
			panic(fmt.Sprintf("Path must be set for nested field [%s]", k))
		}

		if len(d.NestedFieldToQuery[k].Subqueries) < 1 {
			panic(fmt.Sprintf("Subqueries must be set for nested field [%s]", k))
		}

		for sK, sV := range d.NestedFieldToQuery[k].Subqueries {
			if strings.Contains(sK, "$value") {
				// This exists for 3 reasons:
				// 1. In the case of in, it's undefined what it would be since there are multiple values.
				// 2. It would in theory be a form of injection since the users could supply anything and look at any field.
				//   *  Other things shouldn't have this property because you should be validating fields match patterns.
				// 3. I didn't implement this so removing this panic only defers a problem from start up, to runtime.
				panic(fmt.Sprintf("You cannot use $value as replacement in a key in [%s]", sK))
			}

			sqField := sK
			sqValue := sV.Value

			for _, group := range groupKeys {
				if group == "" {
					continue
				}
				sqField = strings.ReplaceAll(sqField, "$"+group, "")
				sqValue = strings.ReplaceAll(sqValue, "$"+group, "")
			}

			sqValue = strings.ReplaceAll(sqValue, "$value", "")

			if strings.Contains(sqField, "$") {
				panic(fmt.Sprintf("Not all templates replaced in nested field [%s] key [%s], after replacement left over with: %s ", k, sK, sqField))
			}

			if strings.Contains(sqValue, "$") {
				panic(fmt.Sprintf("Not all templates replaced in nested field [%s] key [%s] with value [%s], after replacement left over with: %s", k, sK, sV.Value, sqValue))
			}

		}

	}
}

func sortByDecreasingLength(groupKeys []string) {
	// We need to sort the group keys in decreasing order of length
	// So that we resolve templates in the correct order.
	sort.Slice(groupKeys, func(i, j int) bool {
		// First sort by length in decreasing order
		if len(groupKeys[i]) != len(groupKeys[j]) {
			return len(groupKeys[i]) > len(groupKeys[j])
		}
		// Then sort alphabetically
		return groupKeys[i] < groupKeys[j]
	})
}

var _ epsearchast.SemanticReducer[JsonObject] = (*DefaultEsQueryBuilder)(nil)

func (d DefaultEsQueryBuilder) PostVisitAnd(rs []*JsonObject) (*JsonObject, error) {
	return &JsonObject{
		"bool": map[string]any{
			"must": rs,
		},
	}, nil
}

func (d DefaultEsQueryBuilder) PostVisitOr(rs []*JsonObject) (*JsonObject, error) {
	return &JsonObject{
		"bool": map[string]any{
			"should": rs,
			// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-minimum-should-match.html
			"minimum_should_match": 1,
		},
	}, nil
}

func (d DefaultEsQueryBuilder) VisitIn(args ...string) (*JsonObject, error) {
	b := d.GetTermsQueryBuilderForEqualityField()

	return d.buildQueryWithBuilder(b, args...)
}

func (d DefaultEsQueryBuilder) GetTermsQueryBuilderForEqualityField() func(args ...string) *JsonObject {
	return func(args ...string) *JsonObject {
		return &JsonObject{
			"terms": map[string]any{
				d.GetFieldMapping(args[0]).Equality: args[1:],
			},
		}
	}
}

func (d DefaultEsQueryBuilder) VisitEq(first, second string) (*JsonObject, error) {
	b := d.GetTermQueryBuilderForEqualityField()

	return d.buildQueryWithBuilder(b, first, second)
}

func (d DefaultEsQueryBuilder) GetTermQueryBuilderForEqualityField() func(args ...string) *JsonObject {
	return func(args ...string) *JsonObject {
		return &JsonObject{
			"term": map[string]any{
				d.GetFieldMapping(args[0]).Equality: args[1],
			},
		}
	}
}

func (d DefaultEsQueryBuilder) VisitContains(first, second string) (*JsonObject, error) {
	b := d.GetTermQueryBuilderForArrayField()

	return d.buildQueryWithBuilder(b, first, second)
}

func (d DefaultEsQueryBuilder) VisitContainsAny(args ...string) (*JsonObject, error) {
	b := d.GetTermsQueryBuilderForArrayField()

	return d.buildQueryWithBuilder(b, args...)
}

func (d DefaultEsQueryBuilder) VisitContainsAll(args ...string) (*JsonObject, error) {
	// Build individual term queries for each value
	b := d.GetTermQueryBuilderForArrayField()

	var termQueries []*JsonObject
	for _, value := range args[1:] {
		query, err := d.buildQueryWithBuilder(b, args[0], value)
		if err != nil {
			return nil, err
		}
		termQueries = append(termQueries, query)
	}

	// Wrap in a bool query with must clause
	return &JsonObject{
		"bool": map[string]any{
			"must": termQueries,
		},
	}, nil
}

func (d DefaultEsQueryBuilder) GetTermQueryBuilderForArrayField() func(args ...string) *JsonObject {
	return func(args ...string) *JsonObject {
		return &JsonObject{
			"term": map[string]any{
				d.GetFieldMapping(args[0]).Array: args[1],
			},
		}
	}
}

func (d DefaultEsQueryBuilder) GetTermsQueryBuilderForArrayField() func(args ...string) *JsonObject {
	return func(args ...string) *JsonObject {
		return &JsonObject{
			"terms": map[string]any{
				d.GetFieldMapping(args[0]).Array: args[1:],
			},
		}
	}
}

func (d DefaultEsQueryBuilder) VisitText(first, second string) (*JsonObject, error) {
	b := d.BuildMatchBoolPrefixQuery()

	return d.buildQueryWithBuilder(b, first, second)
}

func (d DefaultEsQueryBuilder) BuildMatchBoolPrefixQuery() func(args ...string) *JsonObject {
	return func(args ...string) *JsonObject {

		f := d.DefaultFuzziness

		if f == "" {
			f = "0"
		}

		return &JsonObject{
			"match_bool_prefix": map[string]any{
				d.GetFieldMapping(args[0]).Text: map[string]any{
					"query":     args[1],
					"operator":  "and",
					"fuzziness": f,
				},
			},
		}
	}
}

// Useful doc: https://www.elastic.co/guide/en/elasticsearch/reference/7.17/query-dsl-range-query.html
func (d DefaultEsQueryBuilder) VisitLe(first, second string) (*JsonObject, error) {
	b := d.GetLteRangeQueryBuilder()

	return d.buildQueryWithBuilder(b, first, second)
}

func (d DefaultEsQueryBuilder) GetLteRangeQueryBuilder() func(args ...string) *JsonObject {
	return d.GetRangeQueryBuilder("lte")
}

func (d DefaultEsQueryBuilder) VisitLt(first, second string) (*JsonObject, error) {
	b := d.GetLtRangeQueryBuilder()

	return d.buildQueryWithBuilder(b, first, second)
}

func (d DefaultEsQueryBuilder) GetLtRangeQueryBuilder() func(args ...string) *JsonObject {
	return d.GetRangeQueryBuilder("lt")
}

func (d DefaultEsQueryBuilder) VisitGe(first, second string) (*JsonObject, error) {
	b := d.GetGteRangeQueryBuilder()
	return d.buildQueryWithBuilder(b, first, second)
}

func (d DefaultEsQueryBuilder) GetGteRangeQueryBuilder() func(args ...string) *JsonObject {
	return d.GetRangeQueryBuilder("gte")
}

func (d DefaultEsQueryBuilder) VisitGt(first, second string) (*JsonObject, error) {
	b := d.GetGtRangeQueryBuilder()
	return d.buildQueryWithBuilder(b, first, second)
}

func (d DefaultEsQueryBuilder) GetGtRangeQueryBuilder() func(args ...string) *JsonObject {
	return d.GetRangeQueryBuilder("gt")
}

func (d DefaultEsQueryBuilder) GetRangeQueryBuilder(op string) func(args ...string) *JsonObject {
	return func(args ...string) *JsonObject {
		return &JsonObject{
			"range": map[string]any{
				d.GetFieldMapping(args[0]).Relational: map[string]any{
					op: args[1],
				},
			},
		}
	}
}

func (d DefaultEsQueryBuilder) VisitLike(first, second string) (*JsonObject, error) {
	b := d.GetCaseSensitiveWildcardQueryBuilder()
	return d.buildQueryWithBuilder(b, first, second)
}

func (d DefaultEsQueryBuilder) VisitILike(first, second string) (*JsonObject, error) {
	b := d.GetCaseInsensitiveWildcardQueryBuilder()
	return d.buildQueryWithBuilder(b, first, second)
}

func (d DefaultEsQueryBuilder) GetCaseInsensitiveWildcardQueryBuilder() func(args ...string) *JsonObject {
	return func(args ...string) *JsonObject {
		return &JsonObject{
			"wildcard": map[string]any{
				d.GetFieldMapping(args[0]).Wildcard: map[string]any{
					"value":            d.EscapeWildcardString(args[1]),
					"case_insensitive": true,
				},
			},
		}
	}
}

func (d DefaultEsQueryBuilder) VisitIsNull(first string) (*JsonObject, error) {
	b := d.GetMustNotExistQueryBuilder()
	return d.buildQueryWithBuilder(b, first)
}

func (d DefaultEsQueryBuilder) GetCaseSensitiveWildcardQueryBuilder() func(args ...string) *JsonObject {
	return func(args ...string) *JsonObject {
		return &JsonObject{
			"wildcard": map[string]any{
				d.GetFieldMapping(args[0]).Wildcard: map[string]any{
					"value":            d.EscapeWildcardString(args[1]),
					"case_insensitive": false,
				},
			},
		}
	}
}

func (d DefaultEsQueryBuilder) GetMustNotExistQueryBuilder() func(arg ...string) *JsonObject {
	return func(args ...string) *JsonObject {
		return &JsonObject{
			"bool": map[string]any{
				"must_not": map[string]any{
					"exists": map[string]any{
						"field": d.GetFieldMapping(args[0]).Equality,
					},
				},
			},
		}
	}
}

// GetFieldMapping returns the field name to use for a given operator type, the struct is always guaranteed to return f, if nothing was set.
func (d DefaultEsQueryBuilder) GetFieldMapping(f string) *OperatorTypeToMultiFieldName {
	var o *OperatorTypeToMultiFieldName

	if d.OpTypeToFieldNames[f] == nil {
		o = &OperatorTypeToMultiFieldName{
			Equality:   f,
			Relational: f,
			Text:       f,
			Array:      f,
			Wildcard:   f,
		}
	}

	if v, ok := d.OpTypeToFieldNames[f]; ok {
		o = &OperatorTypeToMultiFieldName{
			Equality:   v.Equality,
			Relational: v.Relational,
			Text:       v.Text,
			Array:      v.Array,
			Wildcard:   v.Wildcard,
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
	}

	return o
}

func (d DefaultEsQueryBuilder) EscapeWildcardString(s string) string {
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

func (d DefaultEsQueryBuilder) buildQueryWithBuilder(b func(args ...string) *JsonObject, args ...string) (*JsonObject, error) {
	nestedQuery, ok, err := d.processNestedFieldToQuery(b, args...)

	if err != nil {
		return nil, err
	}

	if ok {
		return nestedQuery, nil
	}

	return b(args...), nil
}

// processNestedFieldToQuery converts a request for a field that is embedded in an ES nested object into an AND query that indexes into the object by the ID, and then searches the field
// in a nutshell, we can't query eq(field[0].foo, bar), we need to do eq(field.foo, bar):eq(field.id,0). In ES we also need to wrap this in another nested object.
// builder essentially takes the arguments and returns the subquery, it changes whether or not we need to build a match, term, range or other ES query.
func (d DefaultEsQueryBuilder) processNestedFieldToQuery(builder func(args ...string) *JsonObject, args ...string) (*JsonObject, bool, error) {

	var nestedQuery *JsonObject = nil

	numMatches := 0
	for k, v := range d.NestedFieldToQuery {
		pattern := regexp.MustCompile(k)

		if len(args) < 1 {
			return nil, false, fmt.Errorf("no arguments provided")
		}

		searchField := args[0]
		if pattern.MatchString(searchField) {
			// The requested search field matches this NestedFieldToQuery block.
			numMatches++

			groupMap := extractNamedGroupsFromSearchField(pattern, searchField)

			var musts []*JsonObject

			subQueryNames := make([]string, 0, len(v.Subqueries))
			for sqFieldName := range v.Subqueries {
				subQueryNames = append(subQueryNames, sqFieldName)
			}

			sort.Strings(subQueryNames)

			for _, sqFieldName := range subQueryNames {
				sqFieldValue := v.Subqueries[sqFieldName]

				// So in an example where we are supporting queries like field[0].attr=6
				// We might have a regex with ^field\[(?P<id>\d+)\].(?P<attr>\w+)$
				// this saves the stuff between the [] as a named capture group `id`, and the stuff after the period as `attr`
				// We might have replacements of field.idx = $id, and field.$attr=$value
				// We need to replace $id with id in the search value, the $attr in the field in the second term, and the $value with the user argument.

				// This applies the named groups to the field and value.
				sqField, sqValue := applyPatternGroupsToFieldNameAndValue(sqFieldName, sqFieldValue.Value, groupMap)

				// We need to build a set of replacement arguments, substituting $value for the original, and sqField for the first.
				replacedArgs := buildReplacementArgs(args, sqField, sqValue)

				if sqFieldValue.ForceEQ {
					if len(replacedArgs) < 1 {
						return nil, false, fmt.Errorf("expected two values for equality match, got %d", len(replacedArgs))
					}
					if len(replacedArgs) < 2 {
						// It's kind of kludgey hack, but something like is_null(foo[0].bar),
						// doesn't have a second argument, for us to replace with $value.
						// So we just add one.

						replacedArgs = append(replacedArgs, sqValue)
					}

					musts = append(musts, d.GetTermQueryBuilderForEqualityField()(replacedArgs...))
				} else {
					musts = append(musts, builder(replacedArgs...))
				}

			}

			nestedQuery = &JsonObject{
				"nested": JsonObject{
					"path": v.Path,
					"query": JsonObject{
						"bool": JsonObject{
							"must": musts,
						},
					},
				},
			}

		}
	}

	if numMatches > 1 {
		return nil, false, fmt.Errorf("found more than one nested field for %s", args[0])
	}

	if numMatches == 0 {
		return nil, false, nil
	}

	return nestedQuery, true, nil
}

func buildReplacementArgs(args []string, sqField string, sqValue string) []string {
	replacedArgs := make([]string, len(args))

	// Don't allow the field name to be replaced with $value as it can open up injection attacks.
	// Just use the resulting field
	replacedArgs[0] = sqField

	for i := 1; i < len(args); i++ {
		replacedArgs[i] = strings.ReplaceAll(sqValue, `$value`, args[i])
	}
	return replacedArgs
}

func applyPatternGroupsToFieldNameAndValue(sqFieldName string, sqFieldValue string, groupMap map[string]string) (string, string) {
	sqValue := sqFieldValue
	sqField := sqFieldName

	// Sort groupmap keys in decreasing order of length
	groupKeys := make([]string, 0, len(groupMap))
	for k := range groupMap {
		groupKeys = append(groupKeys, k)
	}

	// We need to resolve keys in decreasing order of length
	// So that if we substitute templates with their replacement in consistent order.
	// E.g., if you have templates $user=foo and $username=bar, "$user and $username" needs to resolve to
	// "foo and bar" not "foo and fooname", which if you replace the string user first, is what you get.
	sortByDecreasingLength(groupKeys)

	for _, k := range groupKeys {
		if k == "" {
			continue
		}

		group := k
		replacement := groupMap[k]

		sqField = strings.ReplaceAll(sqField, "$"+group, replacement)
		sqValue = strings.ReplaceAll(sqValue, "$"+group, replacement)
	}
	return sqField, sqValue
}

func extractNamedGroupsFromSearchField(p *regexp.Regexp, first string) map[string]string {
	res := p.FindStringSubmatch(first)

	// Extract specific named group
	groupMap := make(map[string]string)
	for i, name := range p.SubexpNames() {
		if i != 0 && name != "" { // Skip empty or unnamed groups
			groupMap[name] = res[i]
		}
	}
	return groupMap
}
