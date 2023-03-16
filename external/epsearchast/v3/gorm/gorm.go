package v3_gorm

import (
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"strings"
)

type GormVisitor struct {
	Clause string
	Args   []interface{}
}

func NewGormVisitor() *GormVisitor {
	return &GormVisitor{}
}

var _ epsearchast_v3.SearchFilterVisitor = (*GormVisitor)(nil)

func (g *GormVisitor) PreVisit() error {
	g.Clause = " 1=1 "
	g.Args = make([]interface{}, 0)

	return nil
}

func (g *GormVisitor) PostVisit() error {
	// Get rid of leading 1=1
	g.Clause = strings.ReplaceAll(g.Clause, " 1=1 AND ", "")
	// Remove whitespace from either side
	g.Clause = strings.Trim(g.Clause, " ")
	// Remove whitespace before closing parenthesis
	g.Clause = strings.ReplaceAll(g.Clause, " )", ")")

	// Remove whitespace after opening parenthesis
	g.Clause = strings.ReplaceAll(g.Clause, "( ", "(")
	return nil
}

func (g *GormVisitor) PreVisitAnd() error {
	g.Clause += "AND ( 1=1 "
	return nil
}

func (g *GormVisitor) PostVisitAnd() error {
	g.Clause += ") "
	return nil
}

func (g *GormVisitor) VisitIn(args ...string) error {
	s := make([]interface{}, len(args)-1)
	for i, v := range args[1:] {
		s[i] = v
	}

	g.Clause += fmt.Sprintf("AND %s IN ? ", args[0])
	g.Args = append(g.Args, s)
	return nil
}

func (g *GormVisitor) VisitEq(first, second string) error {
	g.Clause += fmt.Sprintf("AND LOWER(%s::text) = LOWER(?) ", first)
	g.Args = append(g.Args, second)

	return nil
}

func (g *GormVisitor) VisitLe(first, second string) error {
	g.Clause += fmt.Sprintf("AND %s <= ? ", first)
	g.Args = append(g.Args, second)
	return nil
}

func (g *GormVisitor) VisitLt(first, second string) error {
	g.Clause += fmt.Sprintf("AND %s < ? ", first)
	g.Args = append(g.Args, second)
	return nil
}

func (g *GormVisitor) VisitGe(first, second string) error {
	g.Clause += fmt.Sprintf("AND %s >= ? ", first)
	g.Args = append(g.Args, second)
	return nil
}

func (g *GormVisitor) VisitGt(first, second string) error {
	g.Clause += fmt.Sprintf("AND %s > ? ", first)
	g.Args = append(g.Args, second)
	return nil
}

func (g *GormVisitor) VisitLike(first, second string) error {
	g.Clause += fmt.Sprintf("AND %s ILIKE ? ", first)
	g.Args = append(g.Args, processLikeWildcards(second))
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
