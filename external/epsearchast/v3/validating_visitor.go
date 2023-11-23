package epsearchast_v3

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"sort"
	"strings"
)

type validatingVisitor struct {
	AllowedOperators map[string][]string
	ColumnAliases    map[string]string
	ValueValidators  map[string]string
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

var _ AstVisitor = (*validatingVisitor)(nil)

// Returns a new validatingVisitor though you should use the helper functions (e.g., [ValidateAstFieldAndOperators]) instead of this method
func NewValidatingVisitor(allowedOps map[string][]string, aliases map[string]string, valueValidators map[string]string) (AstVisitor, error) {

	for k, v := range aliases {
		if _, ok := allowedOps[v]; !ok {
			return nil, fmt.Errorf("alias from `%s` to `%s` points to a field not in the allowed ops", k, v)
		}
	}

	for k, v := range valueValidators {
		if _, ok := allowedOps[k]; !ok {
			if target, ok := aliases[k]; ok {
				// Supporting aliases for validators would be messy because it could be many to one.
				return nil, fmt.Errorf("validator for field `%s` with type `%s` points to an alias of `%s` instead of the field", k, v, target)
			} else {
				return nil, fmt.Errorf("validator for field `%s` with type `%s` points to an unknown field", k, v)
			}

		}
	}

	return &validatingVisitor{
		AllowedOperators: allowedOps,
		ColumnAliases:    aliases,
		ValueValidators:  valueValidators,
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

func (v *validatingVisitor) VisitIsNull(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if err := v.validateFieldAndValue("is_null", fieldName); err != nil {
		return false, err
	}

	return false, nil
}

func (v *validatingVisitor) isOperatorValidForField(operator, requestField string) (bool, error) {

	canonicalField := requestField
	if realName, ok := v.ColumnAliases[requestField]; ok {
		canonicalField = realName
	}

	if _, ok := v.AllowedOperators[canonicalField]; !ok {
		allowedFields := reflect.ValueOf(v.AllowedOperators).MapKeys()
		// Sort the allowed fields to give consistent errors
		sortedAllowedFields := make([]string, len(allowedFields))
		for i := range allowedFields {
			sortedAllowedFields[i] = allowedFields[i].String()
		}
		sort.Strings(sortedAllowedFields)
		return false, fmt.Errorf("unknown field [%s] specified in search filter, allowed fields are %v", requestField, sortedAllowedFields)
	}

	for _, op := range v.AllowedOperators[canonicalField] {
		if strings.ToLower(operator) == strings.ToLower(op) {
			return true, nil
		}
	}

	return false, fmt.Errorf("unknown operator [%s] specified in search filter for field [%s], allowed operators are %v", strings.ToLower(operator), requestField, v.AllowedOperators[canonicalField])
}

func (v *validatingVisitor) validateFieldAndValue(operator, requestField string, values ...string) error {

	if _, err := v.isOperatorValidForField(operator, requestField); err != nil {
		return err
	}

	canonicalField := requestField
	if realName, ok := v.ColumnAliases[requestField]; ok {
		canonicalField = realName
	}

	if vv, ok := v.ValueValidators[canonicalField]; ok {
		for _, value := range values {
			err := validate.Var(value, vv)

			if err != nil {

				if verrors, ok := err.(validator.ValidationErrors); ok {
					if len(verrors) > 0 {
						verror := verrors[0]
						return fmt.Errorf("could not validate [%s] with [%s], value [%s] does not satisify requirement [%s]", requestField, operator, verror.Value(), verror.Tag())
					}
				}

				return fmt.Errorf("could not validate [%s] with [%s] validation error: %w", requestField, value, err)
			}

		}
	}

	return nil
}
