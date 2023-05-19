package epsearchast_v3

// ValidateAstFieldAndOperators determines whether each field is using the allowed operators, a non-nil error is returned if and only if there is a problem.
// Validation of allowed fields is important because failing to do so could allow queries that are not performant against indexes.
func ValidateAstFieldAndOperators(astNode *AstNode, allowedOps map[string][]string) error {
	return ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(astNode, allowedOps, map[string]string{}, map[string]string{})
}

// ValidateAstFieldAndOperatorsWithAliases determines whether each field is using the allowed operators, a non-nil error is returned if and only if there is a problem.
// This version of the function unlike [ValidateAstFieldAndOperators] supports aliased names for fields, which enables the user to specify the same field in different ways, if say a column/field is renamed in the DB.
// Validation of allowed fields is important because failing to do so could allow queries that are not performant against indexes.
func ValidateAstFieldAndOperatorsWithAliases(astNode *AstNode, allowedOps map[string][]string, aliases map[string]string) error {
	return ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(astNode, allowedOps, aliases, map[string]string{})
}

// ValidateAstFieldAndOperatorsWithAliasesAndValueValidation determines whether each field is using the allowed operators, a non-nil error is returned if and only if there is a problem.
// This version of the function unlike [ValidateAstFieldAndOperators] supports validating individual values against a validation rule which can be important in some cases (e.g., if a column/field is an integer in the DB, and string values should be prohibited).
func ValidateAstFieldAndOperatorsWithValueValidation(astNode *AstNode, allowedOps map[string][]string, valueValidators map[string]string) error {
	return ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(astNode, allowedOps, map[string]string{}, valueValidators)
}

// ValidateAstFieldAndOperatorsWithAliasesAndValueValidation determines whether each field is using the allowed operators, a non-nil error is returned if and only if there is a problem.
// This version of the function unlike [ValidateAstFieldAndOperators] supports aliased names for fields which enables the user to specify the same field in different ways, if say a column is renamed in the DB. Validation of allowed fields is important because failing to do so could allow queries that are not performant against indexes.
// This version of the function unlike [ValidateAstFieldAndOperatorsWithAliases] also supports validating individual values against a validation rule which can be important in some cases (e.g., if a column/field is an integer in the DB, and string values should be prohibited).
func ValidateAstFieldAndOperatorsWithAliasesAndValueValidation(astNode *AstNode, allowedOps map[string][]string, aliases map[string]string, valueValidators map[string]string) error {
	visitor, err := NewValidatingVisitor(allowedOps, aliases, valueValidators)

	if err != nil {
		return err
	}

	return astNode.Accept(visitor)
}
