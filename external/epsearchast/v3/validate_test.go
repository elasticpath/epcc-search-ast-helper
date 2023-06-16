package epsearchast_v3

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var binOps = []string{"le", "lt", "eq", "ge", "gt", "like"}

var unaryOps = []string{"is_null"}

var varOps = []string{"in"}

func TestValidationReturnsErrorForBinaryOperatorsWhenAstUsesInvalidOperatorForKnownField(t *testing.T) {
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

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			otherBinOp := binOps[(idx+1)%len(binOps)]

			// Execute SUT
			err = ValidateAstFieldAndOperators(ast, map[string][]string{"amount": {otherBinOp}})

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("unknown operator [%s] specified in search filter for field [amount], allowed operators are [%s]", binOp, otherBinOp))
		})
	}

}

func TestValidationReturnsErrorForBinaryOperatorsWhenAstUsesUnknownField(t *testing.T) {
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

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			otherBinOp := binOps[(idx+1)%len(binOps)]

			// Execute SUT
			err = ValidateAstFieldAndOperators(ast, map[string][]string{"other_field": {otherBinOp}})

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

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			err = ValidateAstFieldAndOperators(ast, map[string][]string{"amount": {binOp}})

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

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			err = ValidateAstFieldAndOperatorsWithAliases(ast, map[string][]string{"total": {binOp}}, map[string]string{"amount": "total"})

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestValidationReturnsErrorForBinaryOperatorsFailedValueValidationWhenAstUseAliases(t *testing.T) {

	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": [ "years",  "ancient"]
			}
			`, strings.ToUpper(binOp))

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(ast, map[string][]string{"age": {binOp}}, map[string]string{"years": "age"}, map[string]string{"age": "number"})

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("could not validate [years] with [%s]", binOp))
		})
	}
}

func TestValidationReturnsNoErrorForBinaryOperatorsWhenAstUseAliasesAndValueValidationAndSatisfiesConstraints(t *testing.T) {

	for _, binOp := range binOps {
		t.Run(fmt.Sprintf("%s", binOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": [ "years",  "70"]
			}
			`, strings.ToUpper(binOp))

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(ast, map[string][]string{"age": {binOp}}, map[string]string{"years": "age"}, map[string]string{"age": "number"})

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestValidationReturnsNoErrorForUnaryOperatorWhenAstSatisfiesConstraints(t *testing.T) {

	for _, unaryOp := range unaryOps {
		t.Run(fmt.Sprintf("%s", unaryOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": ["amount"]
			}
			`, strings.ToUpper(unaryOp))

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			err = ValidateAstFieldAndOperators(ast, map[string][]string{"amount": {unaryOp}})

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestValidationReturnsNoErrorForUnaryOperatorWhenAstUsesAliasesAndSatisfiesConstraints(t *testing.T) {

	for _, unaryOp := range unaryOps {
		t.Run(fmt.Sprintf("%s", unaryOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": ["amount"]
			}
			`, strings.ToUpper(unaryOp))

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			err = ValidateAstFieldAndOperatorsWithAliases(ast, map[string][]string{"total": {unaryOp}}, map[string]string{"amount": "total"})

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestValidationReturnsNoErrorForUnaryOperatorsWhenAstUseAliasesAndValueValidationAndSatisfiesConstraints(t *testing.T) {

	for _, unaryOp := range unaryOps {
		t.Run(fmt.Sprintf("%s", unaryOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": [ "order_status"]
			}
			`, strings.ToUpper(unaryOp))

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			// Note: value validation doesn't do anything with is_null but importantly it doesn't crash which is what we test
			err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(ast, map[string][]string{"status": {unaryOp}}, map[string]string{"order_status": "status"}, map[string]string{"status": "oneof=incomplete complete processing cancelled"})

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestSmokeTestAndWithUnaryAndVariableReturnsErrorWhenBothAreInvalid(t *testing.T) {
	for _, varOp := range varOps {
		for _, unaryOp := range unaryOps {

			t.Run(fmt.Sprintf("%s/%s", varOp, unaryOp), func(t *testing.T) {
				// Fixture Setup
				// language=JSON
				jsonTxt := fmt.Sprintf(`
			{ 
				"type": "AND",
				"children": [
					{
					"type": "%s",
					"args": [ "status",  "complete", "cancelled"]
					},
					{
					"type": "%s",
					"args": [ "some_field"]
					}
				]
}`, strings.ToUpper(varOp), strings.ToUpper(unaryOp))

				ast, err := GetAst(jsonTxt)
				require.NoError(t, err)

				// Execute SUT
				err = ValidateAstFieldAndOperatorsWithValueValidation(ast, map[string][]string{"status": {varOp}, "other_field": {unaryOp}}, map[string]string{"status": "oneof=incomplete complete processing cancelled"})

				// Verification
				require.ErrorContains(t, err, fmt.Sprint("unknown field [some_field] specified in search filter"))
			})

		}
	}
}

func TestValidationReturnsErrorForVariableOperatorsWhenAstUsesInvalidOperatorForKnownField(t *testing.T) {
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

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			otherBinOp := binOps[(idx+1)%len(binOps)]

			// Execute SUT
			err = ValidateAstFieldAndOperators(ast, map[string][]string{"amount": {otherBinOp}})

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("unknown operator [%s] specified in search filter for field [amount], allowed operators are [%s]", varOp, otherBinOp))
		})
	}

}

func TestValidationReturnsErrorForVariableOperatorsWhenAstUsesUnknownField(t *testing.T) {
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

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			otherBinOp := binOps[(idx+1)%len(binOps)]

			// Execute SUT
			err = ValidateAstFieldAndOperators(ast, map[string][]string{"other_field": {otherBinOp}})

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

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			err = ValidateAstFieldAndOperators(ast, map[string][]string{"amount": {varOp}})

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

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			err = ValidateAstFieldAndOperatorsWithAliases(ast, map[string][]string{"total": {varOp}}, map[string]string{"amount": "total"})

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestValidationReturnsErrorForVariableOperatorsFailedValueValidationWhenAstUseAliases(t *testing.T) {

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

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			err = ValidateAstFieldAndOperatorsWithValueValidation(ast, map[string][]string{"email": {varOp}}, map[string]string{"email": "email"})

			// Verification
			require.ErrorContains(t, err, fmt.Sprintf("could not validate [email] with [%s]", varOp))
		})
	}
}

func TestValidationReturnsNoErrorForVariableOperatorsWhenAstUseAliasesAndValueValidationAndSatisfiesConstraints(t *testing.T) {

	for _, varOp := range varOps {
		t.Run(fmt.Sprintf("%s", varOp), func(t *testing.T) {
			// Fixture Setup
			// language=JSON
			jsonTxt := fmt.Sprintf(`
			{
				"type": "%s",
				"args": [ "order_status",  "complete", "cancelled"]
			}
			`, strings.ToUpper(varOp))

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(ast, map[string][]string{"status": {varOp}}, map[string]string{"order_status": "status"}, map[string]string{"status": "oneof=incomplete complete processing cancelled"})

			// Verification
			require.NoError(t, err)
		})
	}
}

func TestSmokeTestAndWithBinaryAndVariableReturnsErrorWhenBothAreInvalid(t *testing.T) {
	for _, varOp := range varOps {
		for _, binOp := range binOps {

			t.Run(fmt.Sprintf("%s/%s", varOp, binOp), func(t *testing.T) {
				// Fixture Setup
				// language=JSON
				jsonTxt := fmt.Sprintf(`
			{ 
				"type": "AND",
				"children": [
					{
					"type": "%s",
					"args": [ "status",  "complete", "cancelled"]
					},
					{
					"type": "%s",
					"args": [ "some_field",  "hello"]
					}
				]
}`, strings.ToUpper(varOp), strings.ToUpper(binOp))

				ast, err := GetAst(jsonTxt)
				require.NoError(t, err)

				// Execute SUT
				err = ValidateAstFieldAndOperatorsWithValueValidation(ast, map[string][]string{"status": {varOp}, "other_field": {binOp}}, map[string]string{"status": "oneof=incomplete complete processing cancelled"})

				// Verification
				require.ErrorContains(t, err, fmt.Sprint("unknown field [some_field] specified in search filter"))
			})

		}
	}
}

func TestSmokeTestAndWithBinaryAndVariableReturnsNoErrorWhenBothValid(t *testing.T) {
	for _, varOp := range varOps {
		for _, binOp := range binOps {

			t.Run(fmt.Sprintf("%s/%s", varOp, binOp), func(t *testing.T) {
				// Fixture Setup
				// language=JSON
				jsonTxt := fmt.Sprintf(`
			{ 
				"type": "AND",
				"children": [
					{
					"type": "%s",
					"args": [ "status",  "complete", "cancelled"]
					},
					{
					"type": "%s",
					"args": [ "some_field",  "hello"]
					}
				]
}`, strings.ToUpper(varOp), strings.ToUpper(binOp))

				ast, err := GetAst(jsonTxt)
				require.NoError(t, err)

				// Execute SUT
				err = ValidateAstFieldAndOperatorsWithValueValidation(ast, map[string][]string{"status": {varOp}, "some_field": {binOp}}, map[string]string{"status": "oneof=incomplete complete processing cancelled"})

				// Verification
				require.NoError(t, err)
			})

		}
	}
}

func TestNewConstructorDetectsUnknownAliasTarget(t *testing.T) {
	// Fixture Setup

	// Execute SUT
	err := ValidateAstFieldAndOperatorsWithAliases(nil, map[string][]string{"status": {"eq"}}, map[string]string{"total": "amount"})

	// Verification
	require.ErrorContains(t, err, fmt.Sprintf("alias from `total` to `amount` points to a field not in the allowed ops"))
}

func TestNewConstructorDetectsUnknownValueValidatorTarget(t *testing.T) {
	// Fixture Setup
	// Execute SUT
	err := ValidateAstFieldAndOperatorsWithValueValidation(nil, map[string][]string{"status": {"eq"}}, map[string]string{"total": "int"})

	// Verification
	require.ErrorContains(t, err, fmt.Sprintf("validator for field `total` with type `int` points to an unknown field"))
}

func TestNewConstructorDetectsAliasedValueValidatorTarget(t *testing.T) {
	// Fixture Setup

	// Execute SUT
	err := ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(nil, map[string][]string{"status": {"eq"}}, map[string]string{"state": "status"}, map[string]string{"state": "int"})

	// Verification
	require.ErrorContains(t, err, fmt.Sprintf("validator for field `state` with type `int` points to an alias of `status` instead of the field"))
}
