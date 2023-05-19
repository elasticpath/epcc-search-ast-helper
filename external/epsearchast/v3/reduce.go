package epsearchast_v3

// ReduceAst is a generic function that can be used to compute or build "something" about an AST.
//
// This function recursively calls the supplied f on each node of the tree, passing in the return value of all
// child nodes as an argument.
//
// Depending on what you are doing you may find that [epsearchast_v3.SemanticReduceAst] to be simpler.
func ReduceAst[T any](a *AstNode, f func(*AstNode, []*T) (*T, error)) (*T, error) {
	collector := make([]*T, 0, len(a.Children))
	for _, n := range a.Children {
		v, err := ReduceAst(n, f)
		if err != nil {
			return nil, err
		}
		collector = append(collector, v)
	}

	return f(a, collector)
}
