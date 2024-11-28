package epsearchast_v3

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var binOps = []string{"le", "lt", "eq", "ge", "gt", "like", "text", "ilike", "contains"}

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
			err = ValidateAstFieldAndOperators(ast, map[string][]string{"other_field": {otherBinOp}, "another_field": {otherBinOp}})

			// Verification
			require.EqualError(t, err, fmt.Sprintf("unknown field [amount] specified in search filter, allowed fields are [another_field other_field]"))
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

func TestValidateAstFieldAndOperatorsAllowsRegularExpressionsWithAllowedOps(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "attributes.name",  "Foo"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperators(ast, map[string][]string{"^attributes.([^.]+)$": {"eq"}})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstFieldAndOperatorsReturnsErrorIfNoRegularExpressionMatchesTheValue(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "attributes.name",  "Foo"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperators(ast, map[string][]string{"^values.([^.]+)$": {"eq"}})

	// Verification
	require.ErrorContains(t, err, "unknown field [attributes.name]")
	require.ErrorContains(t, err, "allowed fields are [^values.([^.]+)$]")
}

func TestValidateAstFieldAndOperatorsWithAliasesAllowsAliasesInRegexesWithoutError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "attributes.name",  "Foo"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliases(ast, map[string][]string{"name": {"eq"}}, map[string]string{"^attributes\\.([^.]+)$": "$1"})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstFieldAndOperatorsWithAliasesAllowsAliasesInRegexesWithoutErrorAndMultipleValidators(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "LIKE",
				"args": [ "attributes.description",  "Foo"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliases(ast, map[string][]string{"name": {"eq"}, "description": {"like"}}, map[string]string{"^attributes\\.([^.]+)$": "$1"})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstFieldAndOperatorsWithAliasesAllowsAliasesInRegexesAndDetectsErrorsWhenFieldIsMissing(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "LIKE",
				"args": [ "attributes.description",  "Foo"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliases(ast, map[string][]string{"name": {"eq"}}, map[string]string{"^attributes\\.([^.]+)$": "$1"})

	// Verification
	require.ErrorContains(t, err, "unknown field [attributes.description]")
	require.ErrorContains(t, err, "allowed fields are [name]")
}

func TestValidateAstFieldAndOperatorsWithAliasesAllowsAliasesInRegexesAndValidatorsWithoutError(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "attributes.locales.fr-CA.name",  "Foo"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliases(ast, map[string][]string{"^locales\\.[^.]+\\.name$": {"eq"}}, map[string]string{"^attributes\\.(.+)$": "$1"})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstFieldAndOperatorsWithAliasesAllowsAliasesInRegexesAndValidatorsWithoutErrorWithConjunction(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
	{
		"type": "AND",
		"children": [
			{
			"type": "ILIKE",
			"args": [ "attributes.locales.fr-CA.description",  "Foo"]
			},
			{
			"type": "EQ",
			"args": [ "attributes.locales.en-CA.name",  "Bar"]
			}
		]
	}`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliases(ast, map[string][]string{
		"^locales\\.[^.]+\\.name$":        {"eq"},
		"^locales\\.[^.]+\\.description$": {"ilike"},
	}, map[string]string{"^attributes\\.(.+)$": "$1"})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstFieldAndOperatorsWithAliasesAllowsAliasesInRegexesAndValidatorsWithErrorWithConjunction(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
	{
		"type": "AND",
		"children": [
			{
			"type": "IN",
			"args": [ "attributes.locales.fr-CA.description",  "Foo"]
			},
			{
			"type": "EQ",
			"args": [ "attributes.locales.en-CA.name",  "Bar"]
			}
		]
	}`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliases(ast, map[string][]string{
		"^locales\\.[^.]+\\.name$":        {"eq"},
		"^locales\\.[^.]+\\.description$": {"like"},
	}, map[string]string{"^attributes\\.(.+)$": "$1"})

	// Verification
	require.ErrorContains(t, err, "unknown operator [in]")
	require.ErrorContains(t, err, "for field [attributes.locales.fr-CA.description]")
	require.ErrorContains(t, err, "allowed operators are [like]")
}

func TestValidateAstWithValueValidationUsingARegularExpressionReturnsErrorWhenTermIsNotAllowed(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "LIKE",
				"args": [ "name",  "poppycock"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithValueValidation(ast, map[string][]string{"name": {"like"}}, map[string]string{"^(name|description)$": "excludes=poppycock"})

	// Verification
	require.ErrorContains(t, err, "could not validate [name] with [like]")
	require.ErrorContains(t, err, "value [poppycock]")
	require.ErrorContains(t, err, "requirement [excludes]")
}

func TestValidateAstWithValueValidationUsingARegularExpressionReturnsNoErrorWhenRequestIsValid(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "status",  "paid"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithValueValidation(ast, map[string][]string{"status": {"eq"}}, map[string]string{"^(name|description)$": "oneOf=paid unpaid"})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstFieldAndOperatorsWithAliasesAndValueValidationReturnsNoErrorInConjunctionWhenAValueValidatorIsSatisfied(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
	{
		"type": "AND",
		"children": [
			{
			"type": "LIKE",
			"args": [ "attributes.locales.fr-CA.description",  "Foo"]
			},
			{
			"type": "EQ",
			"args": [ "attributes.locales.en-CA.name",  "Bar"]
			}
		]
	}`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(ast, map[string][]string{
		"^locales\\.[^.]+\\.name$":        {"eq"},
		"^locales\\.[^.]+\\.description$": {"like"},
	}, map[string]string{"^attributes\\.(.+)$": "$1"},
		map[string]string{"^locales\\.[^.]+\\.[a-zA-Z0-9_-]+$": "min=1"})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstFieldAndOperatorsWithAliasesAndValueValidationDetectsAnErrorInConjunctionWhenAValueValidatorIsNotSatisfied(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
	{
		"type": "AND",
		"children": [
			{
			"type": "LIKE",
			"args": [ "attributes.locales.fr-CA.description",  "Foo"]
			},
			{
			"type": "EQ",
			"args": [ "attributes.locales.en-CA.name",  ""]
			}
		]
	}`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(ast, map[string][]string{
		"^locales\\.[^.]+\\.name$":        {"eq"},
		"^locales\\.[^.]+\\.description$": {"like"},
	}, map[string]string{"^attributes\\.(.+)$": "$1"},
		map[string]string{"^locales\\.[^.]+\\.[a-zA-Z0-9_-]+$": "min=1"})

	// Verification
	require.ErrorContains(t, err, "could not validate [attributes.locales.en-CA.name]")
	require.ErrorContains(t, err, "with [eq]")
	require.ErrorContains(t, err, "value []")
	require.ErrorContains(t, err, "requirement [min]")
}

func TestValidateAstWithTypeValidationReturnsNoErrorWhenRequestIsValidAsString(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "status",  "paid"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithFieldTypes(ast, map[string][]string{"status": {"eq"}}, map[string]FieldType{"status": String})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstWithTypeValidationReturnsNoErrorWhenRequestIsValidAsBoolean(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "paid",  "true"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithFieldTypes(ast, map[string][]string{"paid": {"eq"}}, map[string]FieldType{"paid": Boolean})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstWithTypeValidationReturnsErrorWhenRequestIsNotValidAsBoolean(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "paid",  "Yes Sir!"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithFieldTypes(ast, map[string][]string{"paid": {"eq"}}, map[string]FieldType{"paid": Boolean})

	// Verification
	require.ErrorContains(t, err, "could not validate [paid]")
	require.ErrorContains(t, err, "the value [Yes Sir!]")
	require.ErrorContains(t, err, "invalid value for boolean")
}

func TestValidateAstWithTypeValidationReturnsNoErrorWhenRequestIsValidAsInt64(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "amount",  "16"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithFieldTypes(ast, map[string][]string{"amount": {"eq"}}, map[string]FieldType{"amount": Int64})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstWithTypeValidationReturnsErrorWhenRequestIsNotValidAsInt64(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "amount",  "Nothing"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithFieldTypes(ast, map[string][]string{"amount": {"eq"}}, map[string]FieldType{"amount": Int64})

	// Verification
	require.ErrorContains(t, err, "could not validate [amount]")
	require.ErrorContains(t, err, "the value [Nothing]")
	require.ErrorContains(t, err, "invalid value for int64")
}

func TestValidateAstWithTypeValidationReturnsNoErrorWhenRequestIsValidAsInt64AndPassesValidator(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "amount",  "16"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidationAndFieldTypes(ast, map[string][]string{"amount": {"eq"}}, map[string]string{}, map[string]string{"amount": "gt=12"}, map[string]FieldType{"amount": Int64})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstWithTypeValidationReturnsErrorWhenRequestIsAnInt64ButFailsValidator(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "amount",  "4572"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidationAndFieldTypes(ast, map[string][]string{"amount": {"eq"}}, map[string]string{}, map[string]string{"amount": "lt=128"}, map[string]FieldType{"amount": Int64})

	// Verification
	require.ErrorContains(t, err, "could not validate [amount]")
	require.ErrorContains(t, err, "with [eq]")
	require.ErrorContains(t, err, "value [4572] does not satisfy requirement [lt]")
}

func TestValidateAstWithTypeValidationReturnsErrorWhenRequestIsNotValidAsFloat64(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "amount",  "Nothing"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithFieldTypes(ast, map[string][]string{"amount": {"eq"}}, map[string]FieldType{"amount": Float64})

	// Verification
	require.ErrorContains(t, err, "could not validate [amount]")
	require.ErrorContains(t, err, "the value [Nothing]")
	require.ErrorContains(t, err, "invalid value for float64")
}

func TestValidateAstWithTypeValidationReturnsNoErrorWhenRequestIsValidAsFloat64AndIntegerAndPassesValidator(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "amount",  "16"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidationAndFieldTypes(ast, map[string][]string{"amount": {"eq"}}, map[string]string{}, map[string]string{"amount": "gt=12"}, map[string]FieldType{"amount": Float64})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstWithTypeValidationReturnsNoErrorWhenRequestIsValidAsFloat64AndNonIntegerAndPassesValidator(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "amount",  "16.57"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidationAndFieldTypes(ast, map[string][]string{"amount": {"eq"}}, map[string]string{}, map[string]string{"amount": "gt=12"}, map[string]FieldType{"amount": Float64})

	// Verification
	require.NoError(t, err)
}

func TestValidateAstWithTypeValidationReturnsErrorWhenRequestIsAnFloat64AndIntegerButFailsValidator(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "amount",  "4572"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidationAndFieldTypes(ast, map[string][]string{"amount": {"eq"}}, map[string]string{}, map[string]string{"amount": "lt=128"}, map[string]FieldType{"amount": Float64})

	// Verification
	require.ErrorContains(t, err, "could not validate [amount]")
	require.ErrorContains(t, err, "with [eq]")
	require.ErrorContains(t, err, "value [4572] does not satisfy requirement [lt]")
}

func TestValidateAstWithTypeValidationReturnsErrorWhenRequestIsAFloat64AndNonIntegerButFailsValidator(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"args": [ "amount",  "457.42"]
			}
			`
	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	err = ValidateAstFieldAndOperatorsWithAliasesAndValueValidationAndFieldTypes(ast, map[string][]string{"amount": {"eq"}}, map[string]string{}, map[string]string{"amount": "lt=128"}, map[string]FieldType{"amount": Float64})

	// Verification
	require.ErrorContains(t, err, "could not validate [amount]")
	require.ErrorContains(t, err, "with [eq]")
	require.ErrorContains(t, err, "value [457.42] does not satisfy requirement [lt]")
}
