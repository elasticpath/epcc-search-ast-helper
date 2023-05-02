package epsearchast_v3

import (
	"fmt"
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
				"args": [ "amount",  "5"]
			}
			`, strings.ToUpper(binOp))

			astNode, err := GetAst(jsonTxt)
			require.NoError(t, err)

			otherBinOp := binOps[(idx+1)%len(binOps)]

			visitor, err := NewValidatingVisitor(map[string][]string{"amount": {otherBinOp}}, map[string]string{}, map[string]string{})
			require.NoError(t, err)

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
				"args": [ "amount",  "5"]
			}
			`, strings.ToUpper(binOp))

			astNode, err := GetAst(jsonTxt)
			require.NoError(t, err)

			otherBinOp := binOps[(idx+1)%len(binOps)]

			visitor, err := NewValidatingVisitor(map[string][]string{"other_field": {otherBinOp}}, map[string]string{}, map[string]string{})
			require.NoError(t, err)

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("unknown field [amount] specified in search filter, allowed fields are [other_field]"))
		})
	}

}

func TestValidationReturnsNoErrorForBinaryOperatorsWhenAstSatisfiesConstraints(t *testing.T) {

	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": [ "amount",  "5"]
			}
			`, strings.ToUpper(binOp))

			astNode, err := GetAst(jsonTxt)
			require.NoError(t, err)

			visitor, err := NewValidatingVisitor(map[string][]string{"amount": {binOp}}, map[string]string{}, map[string]string{})
			require.NoError(t, err)

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestValidationReturnsNoErrorForBinaryOperatorWhenAstUsesAliasAndSatisfiesContraints(t *testing.T) {

	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": [ "amount",  "5"]
			}
			`, strings.ToUpper(binOp))

			astNode, err := GetAst(jsonTxt)
			require.NoError(t, err)

			visitor, err := NewValidatingVisitor(map[string][]string{"total": {binOp}}, map[string]string{"amount": "total"}, map[string]string{})
			require.NoError(t, err)

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestValidationReturnsErrorForBinaryOperatorsValueValidation(t *testing.T) {

	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": [ "pkey",  "5"]
			}
			`, strings.ToUpper(binOp))

			astNode, err := GetAst(jsonTxt)
			require.NoError(t, err)

			visitor, err := NewValidatingVisitor(map[string][]string{"id": {binOp}}, map[string]string{"pkey": "id"}, map[string]string{"id": "uuid"})
			require.NoError(t, err)

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("could not validate [pkey] with [%s]", binOp))
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

			visitor, err := NewValidatingVisitor(map[string][]string{"amount": {otherBinOp}}, map[string]string{}, map[string]string{})
			require.NoError(t, err)
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

			visitor, err := NewValidatingVisitor(map[string][]string{"other_field": {otherBinOp}}, map[string]string{}, map[string]string{})
			require.NoError(t, err)

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("unknown field [amount] specified in search filter, allowed fields are [other_field]"))
		})
	}

}

func TestValidationReturnsNoErrorForVariableOperatorWhenAstSatisfiesConstraints(t *testing.T) {

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

			visitor, err := NewValidatingVisitor(map[string][]string{"amount": {varOp}}, map[string]string{}, map[string]string{})
			require.NoError(t, err)

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestValidationReturnsNoErrorForVariableOperatorWhenAstUsesAliasesAndSatisfiesConstraints(t *testing.T) {

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

			visitor, err := NewValidatingVisitor(map[string][]string{"total": {varOp}}, map[string]string{"amount": "total"}, map[string]string{})
			require.NoError(t, err)

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestValidationReturnsErrorForVariableOperatorsValueValidation(t *testing.T) {

	for _, varOp := range varOps {
		t.Run(fmt.Sprintf("%s", varOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": ["email", "foo@foo.com", "bar@bar.com", "5"]
			}
			`, strings.ToUpper(varOp))

			astNode, err := GetAst(jsonTxt)
			require.NoError(t, err)

			visitor, err := NewValidatingVisitor(map[string][]string{"email": {varOp}}, map[string]string{}, map[string]string{"email": "email"})
			require.NoError(t, err)

			// Execute SUT
			err = astNode.Accept(visitor)

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("could not validate [email] with [%s]", varOp))
		})
	}
}

//func TestValidationReturnsErrorForPostVisit(t *testing.T) {
//
//	// Fixture Setup
//	// language=JSON
//	jsonTxt := `
//	{
//		"type": "IN",
//		"args": ["amount", "5"]
//	}`
//
//	astNode, err := GetAst(jsonTxt)
//	require.NoError(t, err)
//
//	visitor, err := NewValidatingVisitor(map[string][]string{"amount": {"in"}}, map[string]string{}, map[string]string{})
//	require.NoError(t, err)
//
//	// Execute SUT
//	err = astNode.Accept(visitor)
//
//	// Verification
//	require.ErrorContains(t, err, fmt.Sprintf("mocked error: PostVisit"))
//
//}
//
//func TestValidationReturnsErrorForPostVisitAnd(t *testing.T) {
//
//	// Fixture Setup
//	// language=JSON
//	jsonTxt := `
//	{
//		"type": "AND",
//		"children": [
//		  {
//		    "type": "IN",
//		    "args": ["amount", "5"]
//		  },
//		  {
//			"type": "EQ",
//			"args": [ "status",  "paid"]
//		  }
//		 ]
//}`
//	astNode, err := GetAst(jsonTxt)
//	require.NoError(t, err)
//
//	visitor, err := NewValidatingVisitor(map[string][]string{"amount": {"in"}, "status": {"eq"}}, map[string]string{}, map[string]string{})
//	require.NoError(t, err)
//
//	// Execute SUT
//	err = astNode.Accept(visitor)
//
//	// Verification
//	require.ErrorContains(t, err, fmt.Sprintf("mocked error: PostVisitAnd"))
//
//}
//
//func TestValidationReturnsErrorForPreVisitAnd(t *testing.T) {
//
//	// Fixture Setup
//	// language=JSON
//	jsonTxt := `
//	{
//		"type": "AND",
//		"children": [
//		  {
//		    "type": "IN",
//		    "args": ["amount", "5"]
//		  },
//		  {
//			"type": "EQ",
//			"args": [ "status",  "paid"]
//		  }
//		 ]
//}`
//	astNode, err := GetAst(jsonTxt)
//	require.NoError(t, err)
//
//	visitor, err := NewValidatingVisitor(map[string][]string{"amount": {"in"}, "status": {"eq"}}, map[string]string{}, map[string]string{})
//	require.NoError(t, err)
//
//	// Execute SUT
//	err = astNode.Accept(visitor)
//
//	// Verification
//	require.ErrorContains(t, err, fmt.Sprintf("mocked error: PreVisitAnd"))
//
//}

func TestNewConstructorDetectsUnknownAliasTarget(t *testing.T) {
	// Fixture Setup

	// Execute SUT
	_, err := NewValidatingVisitor(map[string][]string{"status": {"eq"}}, map[string]string{"total": "amount"}, map[string]string{})

	// Verification
	require.ErrorContains(t, err, fmt.Sprintf("alias from `total` to `amount` points to a field not in the allowed ops"))
}

func TestNewConstructorDetectsUnknownValueValidatorTarget(t *testing.T) {
	// Fixture Setup
	// Execute SUT
	_, err := NewValidatingVisitor(map[string][]string{"status": {"eq"}}, map[string]string{}, map[string]string{"total": "int"})

	// Verification
	require.ErrorContains(t, err, fmt.Sprintf("validator for field `total` with type `int` points to an unknown field"))
}

func TestNewConstructorDetectsAliasedValueValidatorTarget(t *testing.T) {
	// Fixture Setup

	// Execute SUT
	_, err := NewValidatingVisitor(map[string][]string{"status": {"eq"}}, map[string]string{"state": "status"}, map[string]string{"state": "int"})

	// Verification
	require.ErrorContains(t, err, fmt.Sprintf("validator for field `state` with type `int` points to an alias of `status` instead of the field"))
}
