package epsearchast_v3

import (
	"fmt"
	"sort"
)

func GetAstDepth(a *AstNode) int {
	if a == nil {
		return 0
	}

	maxDepth := 0
	for _, n := range a.Children {
		depth := GetAstDepth(n)
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth + 1
}

func GetEffectiveIndexIntersectionCount(a *AstNode) (uint64, error) {
	count, err := SemanticReduceAst(a, effectiveIndexIntersectionCount{})

	if err != nil {
		return 0, err
	}

	return *count, nil

}

var _ SemanticReducer[uint64] = (*effectiveIndexIntersectionCount)(nil)

type effectiveIndexIntersectionCount struct {
}

func (e effectiveIndexIntersectionCount) PostVisitAnd(rs []*uint64) (*uint64, error) {
	if len(rs) == 0 {
		return nil, fmt.Errorf("AND node has no children")
	}
	var product uint64 = 1

	for _, r := range rs {
		product *= *r
	}
	return ptr(product)
}

func (e effectiveIndexIntersectionCount) PostVisitOr(rs []*uint64) (*uint64, error) {
	if len(rs) == 0 {
		return nil, fmt.Errorf("OR node has no children")
	}
	var sum uint64

	for _, r := range rs {
		sum += *r
	}
	return ptr(sum)
}

func (e effectiveIndexIntersectionCount) VisitIn(args ...string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitEq(first, second string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitLe(first, second string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitLt(first, second string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitGe(first, second string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitGt(first, second string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitLike(first, second string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitILike(first, second string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitContains(first, second string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitContainsAny(args ...string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitContainsAll(args ...string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitText(first, second string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitIsNull(first string) (*uint64, error) {
	return ptr(1)
}

func ptr(i uint64) (*uint64, error) {
	return &i, nil
}

// GetAllFirstArgs returns all first arguments from the AST nodes.
// For operators like EQ, LT, GE, etc., this returns the field name being queried.
// For operators with multiple arguments like IN or CONTAINS_ANY, only the first argument (the field name) is returned.
// Duplicate field names are included in the result.
func GetAllFirstArgs(a *AstNode) []string {
	result, _ := ReduceAst(a, func(node *AstNode, children []*[]string) (*[]string, error) {
		fields := []string{}

		// Collect the first arg from this node (if it has args)
		if len(node.Args) > 0 {
			fields = append(fields, node.Args[0])
		}

		// Merge results from all children
		for _, child := range children {
			if child != nil {
				fields = append(fields, *child...)
			}
		}

		return &fields, nil
	})

	if result == nil {
		return []string{}
	}
	return *result
}

// GetAllFirstArgsSorted returns all first arguments from the AST nodes in sorted order.
// This is a convenience function that calls GetAllFirstArgs and sorts the result.
// Duplicate field names are included in the result.
func GetAllFirstArgsSorted(a *AstNode) []string {
	fields := GetAllFirstArgs(a)
	sort.Strings(fields)
	return fields
}

// GetAllFirstArgsUnique returns a set of unique first arguments from the AST nodes.
// This is a convenience function that calls GetAllFirstArgs and builds a unique set.
func GetAllFirstArgsUnique(a *AstNode) map[string]struct{} {
	fields := GetAllFirstArgs(a)
	unique := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		unique[field] = struct{}{}
	}
	return unique
}

// HasFirstArg returns true if the specified field name appears as a first argument anywhere in the AST.
// This is useful for quickly checking if a specific field is referenced in the query.
func HasFirstArg(a *AstNode, fieldName string) bool {
	result, _ := ReduceAst(a, func(node *AstNode, children []*bool) (*bool, error) {
		// Check if this node has the field we're looking for
		if len(node.Args) > 0 && node.Args[0] == fieldName {
			found := true
			return &found, nil
		}

		// Check if any children found it
		for _, child := range children {
			if child != nil && *child {
				found := true
				return &found, nil
			}
		}

		notFound := false
		return &notFound, nil
	})

	if result == nil {
		return false
	}
	return *result
}
