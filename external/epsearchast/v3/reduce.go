package epsearchast_v3

import "fmt"

type SemanticReducer[R any] interface {
	PostVisitAnd([]*R) (*R, error)
	VisitIn(args ...string) (*R, error)
	VisitEq(first, second string) (*R, error)
	VisitLe(first, second string) (*R, error)
	VisitLt(first, second string) (*R, error)
	VisitGe(first, second string) (*R, error)
	VisitGt(first, second string) (*R, error)
	VisitLike(first, second string) (*R, error)
}

func ReduceAst[T any](a *AstNode, applyFn func(*AstNode, []*T) (*T, error)) (*T, error) {
	collector := make([]*T, 0, len(a.Children))
	for _, n := range a.Children {
		v, err := ReduceAst(n, applyFn)
		if err != nil {
			return nil, err
		}
		collector = append(collector, v)
	}

	return applyFn(a, collector)
}

func SemanticReduceAst[T any](a *AstNode, v *SemanticReducer[T]) (*T, error) {

	// Why do I need this
	var foo = *v
	f := func(a *AstNode, t []*T) (*T, error) {
		switch a.NodeType {
		case "LT":
			return foo.VisitLt(a.Args[0], a.Args[1])
		case "LE":
			return foo.VisitLe(a.Args[0], a.Args[1])
		case "EQ":
			return foo.VisitEq(a.Args[0], a.Args[1])
		case "GE":
			return foo.VisitGe(a.Args[0], a.Args[1])
		case "GT":
			return foo.VisitGt(a.Args[0], a.Args[1])
		case "LIKE":
			return foo.VisitLike(a.Args[0], a.Args[1])
		case "IN":
			return foo.VisitIn(a.Args...)
		case "AND":
			return foo.PostVisitAnd(t)
		default:
			return nil, fmt.Errorf("unsupported node type: %s", a.NodeType)
		}
	}

	return ReduceAst(a, f)
}
