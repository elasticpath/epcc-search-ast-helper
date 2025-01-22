package epsearchast_v3

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type validatingVisitor struct {
	AllowedOperators map[string][]string
	ColumnAliases    map[string]string
	ValueValidators  map[string]string
	FieldTypes       map[string]FieldType
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

var _ AstVisitor = (*validatingVisitor)(nil)

// Returns a new validatingVisitor though you should use the helper functions (e.g., [ValidateAstFieldAndOperators]) instead of this method
func NewValidatingVisitor(allowedOps map[string][]string, aliases map[string]string, valueValidators map[string]string, fieldTypeMap map[string]FieldType) (AstVisitor, error) {

	for k, v := range aliases {
		if len(k) > 0 && k[0] == '^' && k[len(k)-1] == '$' {
			// We can't validate regular expression based aliases without being too rigid, and having a lot of validation complexity for an edge case.
			// For example, you could declare an alias of `t.(a|b)` to `$1` (i.e., a or b) and then specify validators on just a or b.
			continue
		}

		if _, ok := allowedOps[v]; !ok {
			return nil, fmt.Errorf("alias from `%s` to `%s` points to a field not in the allowed ops", k, v)
		}
	}

	for k, v := range valueValidators {
		if len(k) > 0 && k[0] == '^' && k[len(k)-1] == '$' {
			// We can't validate regular expression based aliases without being too rigid, and having a lot of validation complexity for an edge case.
			// For example, you could declare an alias of `t.(a|b)` to `$1` (i.e., a or b) and then specify validators on just a or b.
			continue
		}

		if _, ok := allowedOps[k]; !ok {
			if target, ok := aliases[k]; ok {
				// Supporting aliases for validators would be messy because it could be many to one.
				return nil, fmt.Errorf("validator for field `%s` with type `%s` points to an alias of `%s` instead of the field", k, v, target)
			} else {
				return nil, fmt.Errorf("validator for field `%s` with type `%s` points to an unknown field", k, v)
			}

		}
	}

	ftMap := map[string]FieldType{}

	for k, v := range fieldTypeMap {
		ftMap[k] = v

		if allowedOps[k] == nil {
			return nil, fmt.Errorf("field type map for field `%s` is specified  but this is an unknown field", k)
		}
	}

	for k := range allowedOps {
		if _, ok := ftMap[k]; !ok {
			ftMap[k] = String
		}
	}

	return &validatingVisitor{
		AllowedOperators: allowedOps,
		ColumnAliases:    aliases,
		ValueValidators:  valueValidators,
		FieldTypes:       ftMap,
	}, nil
}

func (v *validatingVisitor) PreVisit() error {
	return nil
}

func (v *validatingVisitor) PostVisit() error {
	return nil
}

func (v *validatingVisitor) PreVisitAnd(astNode *AstNode) (bool, error) {
	return true, nil
}

func (v *validatingVisitor) PostVisitAnd(astNode *AstNode) error {
	return nil
}

func (v *validatingVisitor) PreVisitOr(astNode *AstNode) (bool, error) {
	return true, nil
}

func (v *validatingVisitor) PostVisitOr(astNode *AstNode) error {
	return nil
}

func (v *validatingVisitor) VisitIn(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("in", fieldName, astNode.Args[1:]...); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) VisitEq(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("eq", fieldName, astNode.Args[1]); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) VisitLe(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("le", fieldName, astNode.Args[1]); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) VisitLt(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("lt", fieldName, astNode.Args[1]); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) VisitGe(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("ge", fieldName, astNode.Args[1]); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) VisitGt(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("gt", fieldName, astNode.Args[1]); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) VisitLike(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("like", fieldName, astNode.Args[1]); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) VisitILike(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("ilike", fieldName, astNode.Args[1]); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) VisitContains(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("contains", fieldName, astNode.Args[1]); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) VisitText(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("text", fieldName, astNode.Args[1]); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) VisitIsNull(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("is_null", fieldName); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) isOperatorValidForField(operator, requestField string) (bool, error) {
	canonicalField := v.resolveFieldName(requestField)

	allowedOperatorsForField, ok := findMatchInMap(canonicalField, v.AllowedOperators)

	if !ok {
		allowedFields := reflect.ValueOf(v.AllowedOperators).MapKeys()
		// Sort the allowed fields to give consistent errors
		sortedAllowedFields := make([]string, len(allowedFields))
		for i := range allowedFields {
			sortedAllowedFields[i] = allowedFields[i].String()
		}
		sort.Strings(sortedAllowedFields)
		return false, fmt.Errorf("unknown field [%s] specified in search filter, allowed fields are %v", requestField, sortedAllowedFields)
	}

	for _, op := range allowedOperatorsForField {
		if strings.ToLower(operator) == strings.ToLower(op) {
			return true, nil
		}
	}

	return false, fmt.Errorf("unknown operator [%s] specified in search filter for field [%s], allowed operators are %v", strings.ToLower(operator), requestField, allowedOperatorsForField)
}

func (v *validatingVisitor) validateFieldAndValue(operator, requestField string, values ...string) error {

	if _, err := v.isOperatorValidForField(operator, requestField); err != nil {
		return err
	}

	canonicalField := v.resolveFieldName(requestField)

	fieldType, ok := findMatchInMap(canonicalField, v.FieldTypes)

	if !ok {
		// This is almost certainly a bug, we should always get a string back if something wasn't set.
		return fmt.Errorf("unknown field type for field [%s]", requestField)
	}

	for _, value := range values {
		err := ValidateValue(fieldType, value)

		if err != nil {
			return fmt.Errorf("could not validate [%s], the value [%s] could not be converted to %s: %w", requestField, value, fieldType, err)
		}
	}

	valueValidatorsForField, ok := findMatchInMap(canonicalField, v.ValueValidators)

	if ok {
		for _, value := range values {

			vt, _ := Convert(fieldType, value)
			err := validate.Var(vt, valueValidatorsForField)

			if err != nil {

				if verrors, ok := err.(validator.ValidationErrors); ok {
					if len(verrors) > 0 {
						verror := verrors[0]
						return fmt.Errorf("could not validate [%s] with [%s], value [%v] does not satisfy requirement [%s]", requestField, operator, verror.Value(), verror.Tag())
					}
				}

				return fmt.Errorf("could not validate [%s] with [%s] validation error: %w", requestField, value, err)
			}

		}
	}

	return nil
}

func (v *validatingVisitor) resolveFieldName(requestField string) string {
	canonicalField := requestField
	if realName, ok := v.ColumnAliases[requestField]; ok {
		canonicalField = realName
	} else {
		for k, v := range v.ColumnAliases {
			if len(k) > 0 && k[0] == '^' && k[len(k)-1] == '$' {
				r := regexp.MustCompile(k)
				canonicalField = string(r.ReplaceAll([]byte(canonicalField), []byte(v)))
			}
		}
	}

	return canonicalField
}

func findMatchInMap[T any](key string, m map[string]T) (T, bool) {

	if v, ok := m[key]; ok {
		return v, true
	}

	for k, v := range m {
		if len(k) > 0 && k[0] == '^' && k[len(k)-1] == '$' {
			r := regexp.MustCompile(k)
			if r.MatchString(key) {
				return v, true
			}
		}
	}
	var zero T

	return zero, false
}
