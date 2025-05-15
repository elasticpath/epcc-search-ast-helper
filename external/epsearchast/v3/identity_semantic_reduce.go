package epsearchast_v3

// IdentitySemanticReducer is a SemanticReducer that returns the same AstNode it is given.
type IdentitySemanticReducer struct{}

func (i IdentitySemanticReducer) PostVisitAnd(nodes []*AstNode) (*AstNode, error) {
	return &AstNode{
		NodeType: "AND",
		Children: nodes,
	}, nil
}

func (i IdentitySemanticReducer) PostVisitOr(nodes []*AstNode) (*AstNode, error) {
	return &AstNode{
		NodeType: "OR",
		Children: nodes,
	}, nil
}

func (i IdentitySemanticReducer) VisitIn(args ...string) (*AstNode, error) {
	return &AstNode{NodeType: "IN", Args: args}, nil
}

func (i IdentitySemanticReducer) VisitEq(first, second string) (*AstNode, error) {
	return &AstNode{NodeType: "EQ", Args: []string{first, second}}, nil
}

func (i IdentitySemanticReducer) VisitLe(first, second string) (*AstNode, error) {
	return &AstNode{NodeType: "LE", Args: []string{first, second}}, nil
}

func (i IdentitySemanticReducer) VisitLt(first, second string) (*AstNode, error) {
	return &AstNode{NodeType: "LT", Args: []string{first, second}}, nil
}

func (i IdentitySemanticReducer) VisitGe(first, second string) (*AstNode, error) {
	return &AstNode{NodeType: "GE", Args: []string{first, second}}, nil
}

func (i IdentitySemanticReducer) VisitGt(first, second string) (*AstNode, error) {
	return &AstNode{NodeType: "GT", Args: []string{first, second}}, nil
}

func (i IdentitySemanticReducer) VisitLike(first, second string) (*AstNode, error) {
	return &AstNode{NodeType: "LIKE", Args: []string{first, second}}, nil
}

func (i IdentitySemanticReducer) VisitILike(first, second string) (*AstNode, error) {
	return &AstNode{NodeType: "ILIKE", Args: []string{first, second}}, nil
}

func (i IdentitySemanticReducer) VisitContains(first, second string) (*AstNode, error) {
	return &AstNode{NodeType: "CONTAINS", Args: []string{first, second}}, nil
}

func (i IdentitySemanticReducer) VisitText(first, second string) (*AstNode, error) {
	return &AstNode{NodeType: "TEXT", Args: []string{first, second}}, nil
}

func (i IdentitySemanticReducer) VisitIsNull(first string) (*AstNode, error) {
	return &AstNode{NodeType: "IS_NULL", Args: []string{first}}, nil
}
