package epsearchast_v3

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var binOps = []string{"le", "lt", "eq", "ge", "gt", "like"}

var varOps = []string{"in"}

func TestValidationCatchesInvalidOperatorForBinaryOperators(t *testing.T) {
	for idx, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"first_arg": "amount",
				"second_arg": "5"
			}
			`, strings.ToUpper(binOp))

			astNode, err := GetAst(jsonTxt)
			require.NoError(t, err)

			otherBinOp := binOps[(idx+1)%len(binOps)]

			mockObj := new(MyMockedVisitor)
			mockObj.On("PreVisit").Return(nil).
				On("PostVisit").Return(nil)

			visitor := NewValidatingVisitor(mockObj, map[string][]string{"amount": {otherBinOp}}, map[string]string{})

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("unknown operator [%s] specified in search filter for field [amount], allowed operators are [%s]", binOp, otherBinOp))
		})
	}

}

func TestValidationReturnsErrorForBinaryOperators(t *testing.T) {

	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"first_arg": "amount",
				"second_arg": "5"
			}
			`, strings.ToUpper(binOp))

			astNode, err := GetAst(jsonTxt)
			require.NoError(t, err)

			mockObj := new(MyMockedVisitor)
			mockObj.On("PreVisit").Return(nil).
				On("PostVisit").Return(nil).
				On(fmt.Sprintf("Visit%s", strings.Title(binOp)), mock.Anything).Return(true, fmt.Errorf("mocked error: %s", binOp))

			visitor := NewValidatingVisitor(mockObj, map[string][]string{"amount": {binOp}}, map[string]string{})

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("mocked error: %s", binOp))
		})
	}
}

func TestValidationCatchesInvalidOperatorForVariableOperators(t *testing.T) {
	for idx, varOp := range varOps {
		t.Run(fmt.Sprintf("%s", varOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": ["amount", "5"]
			}
			`, strings.ToUpper(varOp))

			astNode, err := GetAst(jsonTxt)
			require.NoError(t, err)

			otherBinOp := binOps[(idx+1)%len(binOps)]

			mockObj := new(MyMockedVisitor)
			mockObj.On("PreVisit").Return(nil).
				On("PostVisit").Return(nil)

			visitor := NewValidatingVisitor(mockObj, map[string][]string{"amount": {otherBinOp}}, map[string]string{})

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("unknown operator [%s] specified in search filter for field [amount], allowed operators are [%s]", varOp, otherBinOp))
		})
	}

}

func TestValidationReturnsErrorForVariableOperators(t *testing.T) {

	for _, varOp := range varOps {
		t.Run(fmt.Sprintf("%s", varOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": ["amount", "5"]
			}
			`, strings.ToUpper(varOp))

			astNode, err := GetAst(jsonTxt)
			require.NoError(t, err)

			mockObj := new(MyMockedVisitor)
			mockObj.On("PreVisit").Return(nil).
				On("PostVisit").Return(nil).
				On(fmt.Sprintf("Visit%s", strings.Title(varOp)), mock.Anything).Return(true, fmt.Errorf("mocked error: %s", varOp))

			visitor := NewValidatingVisitor(mockObj, map[string][]string{"amount": {varOp}}, map[string]string{})

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("mocked error: %s", varOp))
		})
	}
}

func TestColumnAliases(t *testing.T) {
	panic("not implemented")
}

