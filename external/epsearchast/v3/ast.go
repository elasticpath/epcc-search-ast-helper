// Package epsearchast_v3 implements structs and functions for working with the EP-Internal-Search-AST-v3 header.
package epsearchast_v3

import (
	"encoding/json"
	"fmt"
	"strings"
)

// An AstNode presents a particular level in the Abstract Syntax Tree.
type AstNode struct {
	NodeType string     `json:"type"`
	Children []*AstNode `json:"children"`
	Args     []string   `json:"args"`
}

// GetAst converts the JSON to an AstNode if possible, returning an error otherwise.
func GetAst(jsonTxt string) (*AstNode, error) {
	astNode := &AstNode{}

	err := json.Unmarshal([]byte(jsonTxt), astNode)

	if err != nil {
		return nil, fmt.Errorf("could not parse filter:%w", err)
	} else if err := astNode.checkValid(); err != nil {
		return nil, fmt.Errorf("error validating filter:%w", err)
	} else {
		return astNode, nil
	}
}

// The AstVisitor interface provides a way of specifying a [Visitor] for visiting an AST.
//
// This interface is clunky to use for conversions or when you need to return state, and you should use [epsearchast_v3.ReduceAst] instead.
// In particular because the return values are restricted to error, you need to manage and combine the state yourself, which can be more annoying than necessary.
//
// [Visitor]: https://en.wikipedia.org/wiki/Visitor_pattern
type AstVisitor interface {
	PreVisit() error
	PostVisit() error
	PreVisitAnd(astNode *AstNode) (bool, error)
	PostVisitAnd(astNode *AstNode) error
	VisitIn(astNode *AstNode) (bool, error)
	VisitEq(astNode *AstNode) (bool, error)
	VisitLe(astNode *AstNode) (bool, error)
	VisitLt(astNode *AstNode) (bool, error)
	VisitGe(astNode *AstNode) (bool, error)
	VisitGt(astNode *AstNode) (bool, error)
	VisitLike(astNode *AstNode) (bool, error)

	VisitText(astNode *AstNode) (bool, error)
	VisitIsNull(astNode *AstNode) (bool, error)
}

// Accept triggers a visit of the AST.
func (a *AstNode) Accept(v AstVisitor) error {
	err := v.PreVisit()

	if err != nil {
		return err
	}

	err = a.accept(v)

	if err != nil {
		return err
	}

	return v.PostVisit()
}

func (a *AstNode) accept(v AstVisitor) error {

	var descend = false
	var err error = nil

	switch a.NodeType {
	case "AND":
		descend, err = v.PreVisitAnd(a)
	case "IN":
		descend, err = v.VisitIn(a)
	case "EQ":
		descend, err = v.VisitEq(a)
	case "LE":
		descend, err = v.VisitLe(a)
	case "LT":
		descend, err = v.VisitLt(a)
	case "GT":
		descend, err = v.VisitGt(a)
	case "GE":
		descend, err = v.VisitGe(a)
	case "LIKE":
		descend, err = v.VisitLike(a)
	case "TEXT":
		descend, err = v.VisitText(a)
	case "IS_NULL":
		descend, err = v.VisitIsNull(a)
	default:
		return fmt.Errorf("unknown operator %s", a.NodeType)
	}

	if err != nil {
		return err
	}

	if descend {
		for _, c := range a.Children {
			err = c.accept(v)
			if err != nil {
				return err
			}
		}
	}

	switch a.NodeType {
	case "AND":
		err = v.PostVisitAnd(a)

		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AstNode) checkValid() error {
	switch a.NodeType {
	case "AND":
		for _, c := range a.Children {
			err := c.checkValid()
			if err != nil {
				return err
			}
		}
		if len(a.Children) < 2 {
			return fmt.Errorf("and should have at least two children")
		}
	case "IN":
		if len(a.Children) > 0 {
			return fmt.Errorf("operator %v should not have any children", strings.ToLower(a.NodeType))
		}

		if len(a.Args) < 2 {
			return fmt.Errorf("insufficient number of arguments to %s", strings.ToLower(a.NodeType))
		}
	case "EQ", "LE", "LT", "GT", "GE", "LIKE", "TEXT":
		if len(a.Children) > 0 {
			return fmt.Errorf("operator %v should not have any children", strings.ToLower(a.NodeType))
		}

		if len(a.Args) != 2 {
			return fmt.Errorf("operator %v should have exactly 2 arguments", strings.ToLower(a.NodeType))

		}
	case "IS_NULL":
		if len(a.Children) > 0 {
			return fmt.Errorf("operator %v should not have any children", strings.ToLower(a.NodeType))
		}

		if len(a.Args) != 1 {
			return fmt.Errorf("operator %v should have exactly 1 argument", strings.ToLower(a.NodeType))

		}
	default:
		return fmt.Errorf("unknown operator %s", a.NodeType)
	}

	return nil

}
