package v3_gorm_visitor

import (
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"strings"
)

type SubQuery struct {
	Clause string
	Args   []interface{}
}

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

type genericGormVisitor struct{}

var _ SemanticReducer[SubQuery] = (*genericGormVisitor)(nil)

func (g genericGormVisitor) PostVisitAnd(sqs []*SubQuery) (*SubQuery, error) {
	clauses := make([]string, 0, len(sqs))
	args := make([]interface{}, 0)
	for _, sq := range sqs {
		clauses = append(clauses, sq.Clause)
		args = append(args, sq.Args)
	}

	return &SubQuery{
		Clause: strings.Join(clauses, " AND "),
		Args:   args,
	}, nil
}

func (g genericGormVisitor) VisitIn(args ...string) (*SubQuery, error) {
	s := make([]interface{}, len(args)-1)
	for i, v := range args[1:] {
		s[i] = v
	}

	return &SubQuery{
		Clause: fmt.Sprintf("%s IN ?", args[0]),
		Args:   []interface{}{s},
	}, nil
}

func (g genericGormVisitor) VisitEq(first, second string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s = ?", first),
		Args:   []interface{}{second},
	}, nil
}

func (g genericGormVisitor) VisitLe(first, second string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s <= ?", first),
		Args:   []interface{}{second},
	}, nil
}

func (g genericGormVisitor) VisitLt(first, second string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s < ?", first),
		Args:   []interface{}{second},
	}, nil
}

func (g genericGormVisitor) VisitGe(first, second string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s >= ?", first),
		Args:   []interface{}{second},
	}, nil
}

func (g genericGormVisitor) VisitGt(first, second string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s > ?", first),
		Args:   []interface{}{second},
	}, nil
}

func (g genericGormVisitor) VisitLike(first, second string) (*SubQuery, error) {

	return &SubQuery{
		Clause: fmt.Sprintf("%s ILIKE ?", first),
		Args:   []interface{}{processLikeWildcards(second)},
	}, nil
}

func processLikeWildcards(valString string) string {
	if valString == "*" {
		return "%"
	}
	var startsWithStar = strings.HasPrefix(valString, "*")
	var endsWithStar = strings.HasSuffix(valString, "*")
	if startsWithStar {
		valString = valString[1:]
	}
	if endsWithStar {
		valString = valString[:len(valString)-1]
	}
	valString = escapeWildcards(valString)
	if startsWithStar {
		valString = "%" + valString
	}
	if endsWithStar {
		valString += "%"
	}
	return valString
}

func escapeWildcards(valString string) string {
	valString = strings.ReplaceAll(valString, "%", "\\%")
	valString = strings.ReplaceAll(valString, "_", "\\_")
	return valString
}

//func GormVisitorReducer[R any](a *epsearchast_v3.AstNode, applyFn func(*epsearchast_v3.AstNode, []*R) (*R, error)) {
//
//}

func GormVisitorReducer[R any](v *SemanticReducer[R]) func(a *epsearchast_v3.AstNode, c []*R) (*R, error) {

	var foo = *v

	return func(a *epsearchast_v3.AstNode, subQueries []*R) (*R, error) {
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
			return foo.PostVisitAnd(subQueries)
		default:
			return nil, fmt.Errorf("unsupported node type: %s", a.NodeType)
		}
	}

}
