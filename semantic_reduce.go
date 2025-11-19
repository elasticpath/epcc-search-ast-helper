package epsearchast

import "fmt"

// A SemanticReducer is essentially collection of functions that make it easier to reduce things that working with [epsearchast.AstNode]'s directly.
//
// It provides an individual method for each allowed keyword in the AST, which can make some transforms easier. In particular
// only conjunction operators are required to handle the child arguments, and most other types have there arguments passed in the right type.
type SemanticReducer[R any] interface {
	PostVisitAnd([]*R) (*R, error)
	PostVisitOr([]*R) (*R, error)
	VisitIn(args ...string) (*R, error)
	VisitEq(first, second string) (*R, error)
	VisitLe(first, second string) (*R, error)
	VisitLt(first, second string) (*R, error)
	VisitGe(first, second string) (*R, error)
	VisitGt(first, second string) (*R, error)
	VisitLike(first, second string) (*R, error)
	VisitILike(first, second string) (*R, error)
	VisitContains(first, second string) (*R, error)
	VisitContainsAny(args ...string) (*R, error)
	VisitContainsAll(args ...string) (*R, error)
	VisitText(first, second string) (*R, error)
	VisitIsNull(first string) (*R, error)
}

// SemanticReduceAst adapts an epsearchast.SemanticReducer for use with the epsearchast.ReduceAst function.
func SemanticReduceAst[T any](a *AstNode, v SemanticReducer[T]) (*T, error) {
	f := func(a *AstNode, t []*T) (*T, error) {
		switch a.NodeType {
		case "LT":
			return v.VisitLt(a.Args[0], a.Args[1])
		case "LE":
			return v.VisitLe(a.Args[0], a.Args[1])
		case "EQ":
			return v.VisitEq(a.Args[0], a.Args[1])
		case "GE":
			return v.VisitGe(a.Args[0], a.Args[1])
		case "GT":
			return v.VisitGt(a.Args[0], a.Args[1])
		case "LIKE":
			return v.VisitLike(a.Args[0], a.Args[1])
		case "ILIKE":
			return v.VisitILike(a.Args[0], a.Args[1])
		case "CONTAINS":
			return v.VisitContains(a.Args[0], a.Args[1])
		case "CONTAINS_ANY":
			return v.VisitContainsAny(a.Args...)
		case "CONTAINS_ALL":
			return v.VisitContainsAll(a.Args...)
		case "TEXT":
			return v.VisitText(a.Args[0], a.Args[1])
		case "IN":
			return v.VisitIn(a.Args...)
		case "AND":
			return v.PostVisitAnd(t)
		case "OR":
			return v.PostVisitOr(t)
		case "IS_NULL":
			return v.VisitIsNull(a.Args[0])
		default:
			return nil, fmt.Errorf("unsupported node type: %s", a.NodeType)
		}
	}

	return ReduceAst(a, f)
}
