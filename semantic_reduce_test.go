package epsearchast

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSemanticReduceNilAstDoesNotPanic(t *testing.T) {
	// Fixture Setup
	var ast *AstNode = nil

	// Execute SUT
	reducer := PanicyReducer{}

	result, err := SemanticReduceAst(ast, reducer)

	// Verification
	require.NoError(t, err)
	require.Nil(t, result)
}

type PanicyReducer struct {
}

func (p PanicyReducer) PostVisitAnd(rs []*string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) PostVisitOr(rs []*string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitIn(args ...string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitEq(first, second string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitLe(first, second string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitLt(first, second string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitGe(first, second string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitGt(first, second string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitLike(first, second string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitILike(first, second string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitContains(first, second string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitContainsAny(args ...string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitContainsAll(args ...string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitText(first, second string) (*string, error) {
	panic("not called")
}

func (p PanicyReducer) VisitIsNull(first string) (*string, error) {
	panic("not called")
}
