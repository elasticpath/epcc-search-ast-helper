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

func TestValidationCatchesInvalidOperatorForBinaryOperatorsForKnownField(t *testing.T) {
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

func TestValidationCatchesInvalidOperatorForBinaryOperatorsForUnknownField(t *testing.T) {
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

			visitor := NewValidatingVisitor(mockObj, map[string][]string{"other_field": {otherBinOp}}, map[string]string{})

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("unknown field [amount] specified in search filter, allowed fields are [other_field]"))
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

func TestValidationReturnsErrorForBinaryOperatorsWithAlias(t *testing.T) {

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

			visitor := NewValidatingVisitor(mockObj, map[string][]string{"total": {binOp}}, map[string]string{"amount": "total"})

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("mocked error: %s", binOp))
		})
	}
}

func TestValidationCatchesInvalidOperatorForVariableOperatorsForKnownField(t *testing.T) {
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

func TestValidationCatchesInvalidOperatorForVariableOperatorsForUnknownField(t *testing.T) {
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

			visitor := NewValidatingVisitor(mockObj, map[string][]string{"other_field": {otherBinOp}}, map[string]string{})

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("unknown field [amount] specified in search filter, allowed fields are [other_field]"))
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

func TestValidationReturnsErrorForVariableOperatorsWithAlias(t *testing.T) {

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

			visitor := NewValidatingVisitor(mockObj, map[string][]string{"total": {varOp}}, map[string]string{"amount": "total"})

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("mocked error: %s", varOp))
		})
	}
}

func TestValidationReturnsErrorForPostVisit(t *testing.T) {

	// Fixture Setup
	// language=JSON
	jsonTxt := `
	{
		"type": "IN",
		"args": ["amount", "5"]
	}`

	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PostVisit").Return(fmt.Errorf("mocked error: PostVisit")).
		On("VisitIn", mock.Anything).Return(true, nil)

	visitor := NewValidatingVisitor(mockObj, map[string][]string{"amount": {"in"}}, map[string]string{})

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.ErrorContains(t, err, fmt.Sprintf("mocked error: PostVisit"))

}

func TestValidationReturnsErrorForPostVisitAnd(t *testing.T) {

	// Fixture Setup
	// language=JSON
	jsonTxt := `
	{ 
		"type": "AND",
		"children": [
		  {
		    "type": "IN",
		    "args": ["amount", "5"]
		  },
		  { 
			"type": "EQ",
			"first_arg": "status",
			"second_arg": "paid"
		  }
		 ]	
}`
	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("VisitIn", mock.Anything).Return(true, nil).
		On("VisitEq", mock.Anything).Return(true, nil).
		On("PreVisitAnd", mock.Anything).Return(true, nil).
		On("PostVisitAnd", mock.Anything).Return(false, fmt.Errorf("mocked error: PostVisitAnd"))

	visitor := NewValidatingVisitor(mockObj, map[string][]string{"amount": {"in"}, "status": {"eq"}}, map[string]string{})

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.ErrorContains(t, err, fmt.Sprintf("mocked error: PostVisitAnd"))

}

func TestValidationReturnsErrorForPreVisitAnd(t *testing.T) {

	// Fixture Setup
	// language=JSON
	jsonTxt := `
	{ 
		"type": "AND",
		"children": [
		  {
		    "type": "IN",
		    "args": ["amount", "5"]
		  },
		  { 
			"type": "EQ",
			"first_arg": "status",
			"second_arg": "paid"
		  }
		 ]	
}`
	astNode, err := GetAst(jsonTxt)
	require.NoError(t, err)

	mockObj := new(MyMockedVisitor)
	mockObj.On("PreVisit").Return(nil).
		On("PreVisitAnd", mock.Anything).Return(false, fmt.Errorf("mocked error: PreVisitAnd"))

	visitor := NewValidatingVisitor(mockObj, map[string][]string{"amount": {"in"}, "status": {"eq"}}, map[string]string{})

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.ErrorContains(t, err, fmt.Sprintf("mocked error: PreVisitAnd"))

}
