package epsearchast_v3

import "fmt"

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

func (e effectiveIndexIntersectionCount) VisitText(first, second string) (*uint64, error) {
	return ptr(1)
}

func (e effectiveIndexIntersectionCount) VisitIsNull(first string) (*uint64, error) {
	return ptr(1)
}

func ptr(i uint64) (*uint64, error) {
	return &i, nil
}
