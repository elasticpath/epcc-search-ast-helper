// Package epsearchast_v3 implements structs and functions for working with the EP-Internal-Search-AST-v3 header.
package epsearchast_v3

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// An AstNode presents a particular level in the Abstract Syntax Tree.
type AstNode struct {
	NodeType string     `json:"type"`
	Children []*AstNode `json:"children"`
	Args     []string   `json:"args"`
}

func (a *AstNode) AsFilter() string {
	sb := strings.Builder{}
	switch a.NodeType {
	case "AND":
		for _, c := range a.Children {
			sb.WriteString(c.AsFilter())
			sb.WriteString(":")
		}
	default:
		sb.WriteString(strings.ToLower(a.NodeType))
		sb.WriteString("(")
		for i, arg := range a.Args {
			sb.WriteRune('"')
			sb.WriteString(strings.Replace(arg, `"`, `\"`, -1))
			sb.WriteRune('"')
			if i < len(a.Args)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")
	}

	return sb.String()
}

// GetAst converts the JSON to an AstNode if possible, returning an error otherwise.
// If the Error is a ParsingErr it largely means you should treat the error as a 5xx.
// If the Error is a ValidationErr it largely means you should treat the error as a 4xx.
func GetAst(jsonTxt string) (*AstNode, error) {
	astNode := &AstNode{}

	err := json.Unmarshal([]byte(jsonTxt), astNode)

	if err != nil {
		// url decode jsonTxt
		decoded, urlDecodingError := url.QueryUnescape(jsonTxt)

		if urlDecodingError == nil {
			urlDecodingError = json.Unmarshal([]byte(decoded), astNode)

			if urlDecodingError != nil {
				return nil, NewParsingErr(fmt.Errorf("error parsing decoded filter: %w %v", err, urlDecodingError))
			}
		} else {
			return nil, NewParsingErr(fmt.Errorf("%w, error decoding: %v", err, urlDecodingError))
		}
	}

	if err := astNode.checkValid(); err != nil {
		// It might not be obvious why an invalid AST should be a validation error, especially if
		// we receive something that doesn't make any sense like ge(a). The main argument case where we should
		// treat this as a validation is an unknown operator. However, in theory the upstream generator
		// passing us something likely means we should treat it as unsupported if we are out of date.
		return nil, NewValidationErr(fmt.Errorf("error validating filter (%s) :%w", astNode.AsFilter(), err))
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
	VisitILike(astNode *AstNode) (bool, error)
	VisitContains(astNode *AstNode) (bool, error)
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
	case "ILIKE":
		descend, err = v.VisitILike(a)
	case "CONTAINS":
		descend, err = v.VisitContains(a)
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
	case "EQ", "LE", "LT", "GT", "GE", "LIKE", "ILIKE", "CONTAINS", "TEXT":
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
		return fmt.Errorf("unsupported operator %s()", strings.ToLower(a.NodeType))
	}

	return nil

}
