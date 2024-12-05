package epsearchast_v3

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReduceNilAstDoesNotPanic(t *testing.T) {
	// Fixture Setup
	var ast *AstNode = nil

	// Execute SUT
	result, err := ReduceAst(ast, func(*AstNode, []*string) (*string, error) {
		panic("should not be called")
	})

	// Verification
	require.NoError(t, err)
	require.Nil(t, result)
}
