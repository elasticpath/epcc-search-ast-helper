package epsearchast

import "regexp"

// ApplyAliases will return a new AST where all aliases have been resolved to their new value.
// This function should be called after validating it.
func ApplyAliases(a *AstNode, aliases map[string]string) (*AstNode, error) {
	aliasFunc := func(a *AstNode, children []*AstNode) (*AstNode, error) {

		newArgs := make([]string, len(a.Args))
		copy(newArgs, a.Args)

		if len(newArgs) > 0 {
			if v, ok := aliases[newArgs[0]]; ok {
				newArgs[0] = v
			} else {
				for k, v := range aliases {
					if len(k) > 0 && k[0] == '^' && k[len(k)-1] == '$' {
						r := regexp.MustCompile(k)

						newArgs[0] = string(r.ReplaceAll([]byte(newArgs[0]), []byte(v)))
					}
				}
			}
		} else {
			newArgs = nil
		}

		// When we unmarshal the JSON AST a node with no children has nil for the field.
		// Reduce would get messy if you could pass in a nil.
		// if we want to do equality testing in Tests we need to not set empty children.
		// Or maybe make it a non pointer type or something.
		var childrenNodes []*AstNode = nil

		if len(children) > 0 {
			childrenNodes = children
		}

		return &AstNode{
			NodeType: a.NodeType,
			Children: childrenNodes,
			Args:     newArgs,
		}, nil
	}

	return ReduceAst(a, aliasFunc)
}
