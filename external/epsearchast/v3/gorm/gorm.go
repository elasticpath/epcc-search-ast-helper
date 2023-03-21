package v3_gorm

import (
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"strings"
)

type GormVisitor struct {
	clause string
	args   []interface{}
}

var _ epsearchast_v3.SearchFilterVisitor = (*GormVisitor)(nil)

func (g *GormVisitor) PreVisit() error {
	g.clause = " 1=1 "
	g.args = make([]interface{}, 0)

	return nil
}

func (g *GormVisitor) PostVisit() error {
	//g.Query = g.Query.Where(g.clause, g.args...)
	return nil
}

func (g *GormVisitor) PreVisitAnd() error {
	g.clause += "AND ( 1=1 "
	return nil
}

func (g *GormVisitor) PostVisitAnd() error {
	g.clause += ") "
	return nil
}

func (g *GormVisitor) VisitIn(args ...string) error {
	s := make([]interface{}, len(args)-1)
	for i, v := range args {
		s[i] = v
	}

	g.clause += fmt.Sprintf("AND %s IN ? ", args[0])
	g.args = append(g.args, s)
	return nil
}

func (g *GormVisitor) VisitEq(first, second string) error {
	g.clause += fmt.Sprintf("AND LOWER(%s::text) = LOWER(?) ", first)
	g.args = append(g.args, second)

	return nil
}

func (g *GormVisitor) VisitLe(first, second string) error {
	g.clause += fmt.Sprintf("AND %s <= ? ", first)
	g.args = append(g.args, second)
	return nil
}

func (g *GormVisitor) VisitLt(first, second string) error {
	g.clause += fmt.Sprintf("AND %s < ? ", first)
	g.args = append(g.args, second)
	return nil
}

func (g *GormVisitor) VisitGe(first, second string) error {
	g.clause += fmt.Sprintf("AND %s >= ? ", first)
	g.args = append(g.args, second)
	return nil
}

func (g *GormVisitor) VisitGt(first, second string) error {
	g.clause += fmt.Sprintf("AND %s > ? ", first)
	g.args = append(g.args, second)
	return nil
}

func (g *GormVisitor) VisitLike(first, second string) error {
	g.clause += fmt.Sprintf("AND %s ILIKE ? ", first)
	g.args = append(g.args, processLikeWildcards(second))
	return nil
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
