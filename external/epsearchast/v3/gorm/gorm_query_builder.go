package epsearchast_v3_gorm

import (
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"strconv"
	"strings"
)

type SubQuery struct {
	// The clause that can be passed to Where
	Clause string
	// An array that should be passed in using the ... operator to Where
	Args []interface{}
}

type DefaultGormQueryBuilder struct{}

var _ epsearchast_v3.SemanticReducer[SubQuery] = (*DefaultGormQueryBuilder)(nil)

func (g DefaultGormQueryBuilder) PostVisitAnd(sqs []*SubQuery) (*SubQuery, error) {
	clauses := make([]string, 0, len(sqs))
	args := make([]interface{}, 0)
	for _, sq := range sqs {
		clauses = append(clauses, sq.Clause)
		args = append(args, sq.Args...)
	}

	return &SubQuery{
		Clause: "( " + strings.Join(clauses, " AND ") + " )",
		Args:   args,
	}, nil
}

func (g DefaultGormQueryBuilder) VisitIn(args ...string) (*SubQuery, error) {
	s := make([]interface{}, len(args)-1)
	for i, v := range args[1:] {
		s[i] = v
	}

	return &SubQuery{
		Clause: fmt.Sprintf("%s IN ?", args[0]),
		Args:   []interface{}{s},
	}, nil
}

func (g DefaultGormQueryBuilder) VisitEq(first, second string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s = ?", first),
		Args:   []interface{}{second},
	}, nil
}

func (g DefaultGormQueryBuilder) VisitLe(first, second string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s <= ?", first),
		Args:   []interface{}{second},
	}, nil
}

func (g DefaultGormQueryBuilder) VisitLt(first, second string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s < ?", first),
		Args:   []interface{}{second},
	}, nil
}

func (g DefaultGormQueryBuilder) VisitGe(first, second string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s >= ?", first),
		Args:   []interface{}{second},
	}, nil
}

func (g DefaultGormQueryBuilder) VisitGt(first, second string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s > ?", first),
		Args:   []interface{}{second},
	}, nil
}

func (g DefaultGormQueryBuilder) VisitLike(first, second, booleanStr string) (*SubQuery, error) {
	keyword := "LIKE"
	caseInsensitive, err := strconv.ParseBool(booleanStr)
	if err != nil {
		return nil, err
	}
	if caseInsensitive {
		keyword = "ILIKE"
	}
	return &SubQuery{
		Clause: fmt.Sprintf("%s %s ?", first, keyword),
		Args:   []interface{}{g.ProcessLikeWildcards(second)},
	}, nil
}

func (g DefaultGormQueryBuilder) VisitIsNull(first string) (*SubQuery, error) {
	return &SubQuery{
		Clause: fmt.Sprintf("%s IS NULL", first),
	}, nil
}

func (g DefaultGormQueryBuilder) ProcessLikeWildcards(valString string) string {
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
	valString = g.EscapeWildcards(valString)
	if startsWithStar {
		valString = "%" + valString
	}
	if endsWithStar {
		valString += "%"
	}
	return valString
}

func (g DefaultGormQueryBuilder) EscapeWildcards(valString string) string {
	valString = strings.ReplaceAll(valString, "%", "\\%")
	valString = strings.ReplaceAll(valString, "_", "\\_")
	return valString
}
