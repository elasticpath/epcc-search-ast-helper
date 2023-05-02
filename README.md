# EPCC Search AST Helper

## Introduction

This project is designed to help consume the `EP-Internal-Search-Ast-v*` headers




### Validation

This package provides a concise way to validate that the operators and specified in the header are permitted:

```go
package example

import epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"

func Example(ast *epsearchast_v3.AstNode) error {
	var err error
	// The following is an implementation of all the filter operators for orders https://elasticpath.dev/docs/orders/orders-api/orders-api-overview#filtering
	err = epsearchast_v3.ValidateAstFieldAndOperators(ast, map[string][]string {
		"status": {"eq"},
		"payment": {"eq"},
		"shipping": {"eq"},
		"name": {"eq", "like"},
		"email": {"eq", "like"},
		"customer_id": {"eq", "like"},
		"account_id": {"eq", "like"},
		"account_member_id": {"eq", "like"},
		"contact.name": {"eq", "like"},
		"contact.email": {"eq", "like"},
		"shipping_postcode": {"eq", "like"},
		"billing_postcode": {"eq", "like"},
		"with_tax": {"gt", "ge", "lt", "le"},
		"without_tax": {"gt", "ge", "lt", "le"},
		"currency": {"eq"},
		"product_id": {"eq"},
		"product_sku": {"eq"},
		"created_at": {"eq", "gt", "ge", "lt", "le"},
		"updated_at": {"eq", "gt", "ge", "lt", "le"},
    })
	
	if err != nil { 
		return err
    }
	
	// You can additionally create aliases which allows for one field to reference another:
	// In this case any headers that search for a field of `order_status` will be mapped to `status` and use those rules instead. 
	err = epsearchast_v3.ValidateAstFieldAndOperatorsWithAliases(ast, map[string][]string {"status": {"eq"}}, map[string]string {"order_status": "status"})
	if err != nil {
		return err
	}
	
	// Finally you can also supply validators on fields, which may be necessary in some cases depending on your data model or to improve user experience.
	// Validation is provided by the go-playground/validator package https://github.com/go-playground/validator#usage-and-documentation
	err = epsearchast_v3.ValidateAstFieldAndOperatorsWithValueValidation(ast, map[string][]string {"status": {"eq"}}, map[string]string {"status": "oneof=incomplete complete processing cancelled"})
	
	return err
}
```

#### Limitations

At present, you can only use string validators when validating a field, a simple pull request can be created to fix this issue if you need it.
