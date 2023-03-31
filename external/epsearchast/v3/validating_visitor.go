package epsearchast_v3

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type validatingVisitor struct {
	AllowedOperators map[string][]string
	ColumnAliases    map[string]string
	Visitor          AstVisitor
	ValueValidators  map[string]string
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

var _ AstVisitor = (*validatingVisitor)(nil)

func NewValidatingVisitor(visitor AstVisitor, allowedOps map[string][]string, aliases map[string]string, valueValidators map[string]string) (AstVisitor, error) {

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
		Visitor:          visitor,
		AllowedOperators: allowedOps,
		ColumnAliases:    aliases,
		ValueValidators:  valueValidators,
	}, nil
}

func (v *validatingVisitor) PreVisit() error {
	return v.Visitor.PreVisit()
}

func (v *validatingVisitor) PostVisit() error {
	return v.Visitor.PostVisit()
}

func (v *validatingVisitor) PreVisitAnd(astNode *AstNode) (bool, error) {
	return v.Visitor.PreVisitAnd(astNode)
}

func (v *validatingVisitor) PostVisitAnd(astNode *AstNode) (bool, error) {
	return v.Visitor.PostVisitAnd(astNode)
}

func (v *validatingVisitor) VisitIn(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if _, err := v.validateFieldAndValue("in", fieldName, astNode.Args[1:]...); err != nil {
		return false, err
	}

	return v.Visitor.VisitIn(astNode)
}

func (v *validatingVisitor) VisitEq(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.validateFieldAndValue("eq", fieldName, astNode.SecondArg); err != nil {
		return false, err
	}

	return v.Visitor.VisitEq(astNode)
}

func (v *validatingVisitor) VisitLe(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.validateFieldAndValue("le", fieldName, astNode.SecondArg); err != nil {
		return false, err
	}

	return v.Visitor.VisitLe(astNode)
}

func (v *validatingVisitor) VisitLt(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.validateFieldAndValue("lt", fieldName, astNode.SecondArg); err != nil {
		return false, err
	}

	return v.Visitor.VisitLt(astNode)
}

func (v *validatingVisitor) VisitGe(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.validateFieldAndValue("ge", fieldName, astNode.SecondArg); err != nil {
		return false, err
	}

	return v.Visitor.VisitGe(astNode)
}

func (v *validatingVisitor) VisitGt(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.validateFieldAndValue("gt", fieldName, astNode.SecondArg); err != nil {
		return false, err
	}

	return v.Visitor.VisitGt(astNode)
}

func (v *validatingVisitor) VisitLike(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.validateFieldAndValue("like", fieldName, astNode.SecondArg); err != nil {
		return false, err
	}

	return v.Visitor.VisitLike(astNode)
}

func (v *validatingVisitor) isOperatorValidForField(operator, requestField string) (bool, error) {

	canonicalField := requestField
	if realName, ok := v.ColumnAliases[requestField]; ok {
		canonicalField = realName
	}

	if _, ok := v.AllowedOperators[canonicalField]; !ok {
		return false, fmt.Errorf("unknown field [%s] specified in search filter, allowed fields are %v", requestField, reflect.ValueOf(v.AllowedOperators).MapKeys())
	}

	for _, op := range v.AllowedOperators[canonicalField] {
		if strings.ToLower(operator) == strings.ToLower(op) {
			return true, nil
		}
	}

	return false, fmt.Errorf("unknown operator [%s] specified in search filter for field [%s], allowed operators are %v", strings.ToLower(operator), requestField, v.AllowedOperators[canonicalField])
}

func (v *validatingVisitor) validateFieldAndValue(operator, requestField string, values ...string) (bool, error) {

	if _, err := v.isOperatorValidForField(operator, requestField); err != nil {
		return false, err
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
						return false, fmt.Errorf("could not validate [%s] with [%s], value [%s] does not satisify requirement [%s]", requestField, operator, verror.Value(), verror.Tag())
					}
				}

				return false, fmt.Errorf("could not validate [%s] with [%s] validation error: %w", requestField, value, err)
			}

		}
	}

	return true, nil
}
