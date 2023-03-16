package epsearchast_v3

import (
	"fmt"
	"reflect"
	"strings"
)

type ValidatingVisitor struct {
	AllowedOperators map[string][]string
	ColumnAliases    map[string]string
	Visitor          AstVisitor
}

var _ AstVisitor = (*ValidatingVisitor)(nil)

func NewValidatingVisitor(visitor AstVisitor, allowedOps map[string][]string, aliases map[string]string) AstVisitor {
	return &ValidatingVisitor{
		Visitor:          visitor,
		AllowedOperators: allowedOps,
		ColumnAliases:    aliases,
	}
}

func (v *ValidatingVisitor) PreVisit() error {
	return v.Visitor.PreVisit()
}

func (v *ValidatingVisitor) PostVisit() error {
	return v.Visitor.PostVisit()
}

func (v *ValidatingVisitor) PreVisitAnd(astNode *AstNode) (bool, error) {
	return v.Visitor.PreVisitAnd(astNode)
}

func (v *ValidatingVisitor) PostVisitAnd(astNode *AstNode) (bool, error) {
	return v.Visitor.PostVisitAnd(astNode)
}

func (v *ValidatingVisitor) VisitIn(astNode *AstNode) (bool, error) {
	fieldName := astNode.Args[0]

	if _, err := v.isOperatorValidForField("in", fieldName); err != nil {
		return false, err
	}

	return v.Visitor.VisitIn(astNode)
}

func (v *ValidatingVisitor) VisitEq(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.isOperatorValidForField("eq", fieldName); err != nil {
		return false, err
	}

	return v.Visitor.VisitEq(astNode)
}

func (v *ValidatingVisitor) VisitLe(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.isOperatorValidForField("le", fieldName); err != nil {
		return false, err
	}

	return v.Visitor.VisitLe(astNode)
}

func (v *ValidatingVisitor) VisitLt(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.isOperatorValidForField("lt", fieldName); err != nil {
		return false, err
	}

	return v.Visitor.VisitLt(astNode)
}

func (v *ValidatingVisitor) VisitGe(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.isOperatorValidForField("ge", fieldName); err != nil {
		return false, err
	}

	return v.Visitor.VisitGe(astNode)
}

func (v *ValidatingVisitor) VisitGt(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.isOperatorValidForField("gt", fieldName); err != nil {
		return false, err
	}

	return v.Visitor.VisitGt(astNode)
}

func (v *ValidatingVisitor) VisitLike(astNode *AstNode) (bool, error) {
	fieldName := astNode.FirstArg

	if _, err := v.isOperatorValidForField("like", fieldName); err != nil {
		return false, err
	}

	return v.Visitor.VisitLike(astNode)
}

func (v *ValidatingVisitor) isOperatorValidForField(operator, requestField string) (bool, error) {

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
