package v3_gorm

import (
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"strings"
)

type AbstractGormVisitor struct {
	epsearchast_v3.SearchFilterVisitor
}

type GormNode struct {
	NodeType string      `json:"type"`
	Children []*GormNode `json:"children"`
	Args     []string    `json:"args"`
}

type FlatAttributes struct {
	NodeType string   `json:"type"`
	Args     []string `json:"args"`
}

//func Recruse[K any](node *GormNode) K {
//	return nil
//}
//
//func Apply[K any](node *Gorm) K {
//
//}

type SubQuery struct {
	Clause string
	Args   []interface{}
}

//func Recurse(node *epsearchast_v3.AstNode, applyFn func(astNode *epsearchast_v3.AstNode, childVal []*SubQuery) (*SubQuery, error)) (*SubQuery, error) {
//	collector := make([]*SubQuery, 0, len(node.Children))
//	for _, n := range node.Children {
//		v, err := Recurse(n, applyFn)
//		if err != nil {
//			return nil, err
//		}
//		collector = append(collector, v)
//	}
//
//	return Apply(node, collector)
//}

func Apply(a *epsearchast_v3.AstNode, subQueries []*SubQuery) (*SubQuery, error) {

	var validateNoChildVals = true
	var clause string
	args := make([]interface{}, 0)
	switch a.NodeType {

	case "IN":
		clause = fmt.Sprintf("%s IN ?", a.Args[0])
		s := make([]interface{}, len(a.Args)-1)
		for i, v := range a.Args[1:] {
			s[i] = v
		}

		args = append(args, s)
	case "LT":
		clause = fmt.Sprintf("%s < ?", a.Args[0])
		args = append(args, a.Args[1])
	case "LE":
		clause = fmt.Sprintf("%s <= ?", a.Args[0])
		args = append(args, a.Args[1])
	case "EQ":
		clause = fmt.Sprintf("%s = ?", a.Args[0])
		args = append(args, a.Args[1])
	case "GE":
		clause = fmt.Sprintf("%s >= ?", a.Args[0])
		args = append(args, a.Args[1])
	case "GT":
		clause = fmt.Sprintf("%s > ?", a.Args[0])
		args = append(args, a.Args[1])
	case "LIKE":
		clause = fmt.Sprintf("%s ILIKE ?", a.Args[0])
		args = append(args, processLikeWildcards(a.Args[1]))
	case "AND":
		validateNoChildVals = false

		clauses := make([]string, 0, len(subQueries))
		for _, s := range subQueries {
			clauses = append(clauses, s.Clause)
			args = append(args, s.Args)
		}

		clause = strings.Join(clauses, " AND ")

		// Not sure why, but not important now

	default:
		return nil, fmt.Errorf("unknown operator %s", a.NodeType)
	}

	if validateNoChildVals {
		if len(subQueries) > 0 {
			return nil, fmt.Errorf("expected filter argument %s to have no subqueries but got %d", a.NodeType, len(subQueries))
		}
	}

	return &SubQuery{
		Clause: clause,
		Args:   args,
	}, nil
}
