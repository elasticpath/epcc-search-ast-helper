# EPCC Search AST Helper

## Introduction

This project is designed to help consume the `EP-Internal-Search-Ast-v*` headers. In particular, it provides functions for processing these headers in a variety of use cases.


### Retrieving an AST
The `GetAst()` function will convert the JSON header into a struct that can be then be processed by other functions:

```go
package example

import epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"

func Example(headerValue string) (*epsearchast_v3.AstNode, error) {
	
	ast, err := epsearchast_v3.GetAst(headerValue)
	
	if err != nil { 
		return nil, err
    } else { 
		return ast, nil
    }
	
}

```


### Aliases

This package provides a way to support aliases for fields, this will allow a user to specify multiple different names for a field, and still have it validated and converted properly:

```go
package example

import epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"

func Example(ast *epsearchast_v3.AstNode) error {
	
	//The ast from the user will be converted into a new one, and if the user specified a payment_status field, the new ast will have it recorded as status. 
	aliasedAst, err := ApplyAliases(ast, map[string]string{"payment_status": "status"})

	if err != nil { 
		return err
    }
	
	DoSomethingElse(aliasedAst)
	
	return err
}
```

### Validation

This package provides a concise way to validate that the operators and fields specified in the header are permitted:

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


### Generating Queries

#### GORM/SQL

The following examples shows how to generate a Gorm query with this library.

```go
package example

import epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
import epsearchast_v3_gorm "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3/gorm"
import "gorm.io/gorm"

func Example(ast *epsearchast_v3.AstNode, query *gorm.DB, tenantBoundaryId string) error {
	var err error
	
	// Not Shown: Validation
	
	// Create query builder
	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_gorm.SubQuery] = epsearchast_v3_gorm.DefaultGormQueryBuilder{}

	
	sq, err := epsearchast_v3.SemanticReduceAst(ast, qb)

	if err != nil {
		return err
	}

	// Don't forget to add additional filters 
	query.Where("tenant_boundary_id = ?", tenantBoundaryId)
	
	// Don't forget to expand the Args argument with ...
	query.Where(sq.Clause, sq.Args...)
}
```


##### Limitations

1. The GORM builder does not support aliases (easy MR to fix).
2. The GORM builder does not support joins (fixable in theory).
3. There is no way currently to specify the type of a field for SQL, which means everything gets written as a string today (fixable with MR).
4. The `text` operator implementation makes a number of assumptions, and you likely will want to override it's implementation:
   * English is hard coded as the language.
   * Postgres recommends using a [distinct tsvector column and using a stored generated column](https://www.postgresql.org/docs/current/textsearch-tables.html#TEXTSEARCH-TABLES-INDEX). The current implementation does not support this and, you would need to override the method to support it. A simple MR could be made to allow for the Gorm query builder to know if there is a tsvector column and use that.

##### Advanced Customization

In some cases you may want to change the behaviour of the generated SQL, the following example shows how to do that
in this case, we want all eq queries for emails to use the lower case, comparison, and for cart_items field to be numeric.

```go
package example

import (
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"strconv"
)
import epsearchast_v3_gorm "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3/gorm"
import "gorm.io/gorm"


func Example(ast *epsearchast_v3.AstNode, query *gorm.DB, tenantBoundaryId string) error {
	var err error

	// Not Shown: Validation

	// Create query builder
	var qb epsearchast_v3.SemanticReducer[epsearchast_v3_gorm.SubQuery] = &CustomQueryBuilder{}

	sq, err := epsearchast_v3.SemanticReduceAst(ast, qb)

	if err != nil {
		return err
	}

	// Don't forget to add additional filters 
	query.Where("tenant_boundary_id = ?", tenantBoundaryId)
	
	// Don't forget to expand the Args argument with ...
	query.Where(sq.Clause, sq.Args...)
}

type CustomQueryBuilder struct {
	epsearchast_v3_gorm.DefaultGormQueryBuilder
}

func (l *CustomQueryBuilder) VisitEq(first, second string) (*epsearchast_v3_gorm.SubQuery, error) {
	if first == "email" {
		return &epsearchast_v3_gorm.SubQuery{
			Clause: fmt.Sprintf("LOWER(%s::text) = LOWER(?)", first),
			Args:   []interface{}{second},
		}, nil
	} else if first == "cart_items" {
		n, err := strconv.Atoi(second)
		if err != nil {
			return nil, err
		}
		return &epsearchast_v3_gorm.SubQuery{
			Clause: fmt.Sprintf("%s = ?", first),
			Args:   []interface{}{n},
		}, nil
	} else {
		return DefaultGormQueryBuilder.VisitEq(l.DefaultGormQueryBuilder, first, second)
	}
}
```

#### Mongo

The following examples shows how to generate a Mongo Query with this library.

```go
package example

import (
	"context"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	epsearchast_v3_mongo "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Example(ast *epsearchast_v3.AstNode, collection *mongo.Collection, tenantBoundaryQuery bson.M)  (*mongo.Cursor, error) {
	// Not Shown: Validation

	// Create query builder
	var qb epsearchast_v3.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

	// Create Query Object
	queryObj, err := epsearchast_v3.SemanticReduceAst(ast, qb)

	if err != nil {
		return nil, err
	}

	mongoQuery := bson.D{
		{"$and",
			bson.A{
				tenantBoundaryQuery,
				queryObj,
			},
		}}
	
	
	return collection.Find(context.TODO(), mongoQuery)
}
```


##### Limitations

1. The Mongo Query builder is designed to produce filter compatible with the [filter argument in a Query](https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/query-document/#specify-a-query), if a field in the API is a projection that requires computation via the aggregation pipeline, then we would likely need code changes to support that.
2. The [$text](https://www.mongodb.com/docs/v7.0/reference/operator/query/text/#behavior) operator in Mongo has a number of limitations that make it unsuitable for arbitrary queries. In particular in mongo you can only search a collection, not fields for text data, and you must declare a text index. This means that any supplied field in the filter, is just dropped. It is recommended that when using `text` with Mongo, you only allow users to search `text(*,search)` , i.e., force them to use a wildcard as the field name. It is also recommended that you use a [Wildcard](https://www.mongodb.com/docs/manual/core/indexes/index-types/index-text/create-wildcard-text-index/) index to avoid the need of having to remove and modify it over time.

##### Advanced Customization

In some cases you may want to change the behaviour of the generated Mongo, the following example shows how to do that in this case we want to change emails because
we store them only in lower case in the db.

```go
package example

import (
	"context"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	epsearchast_v3_mongo "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

func Example(ast *epsearchast_v3.AstNode, collection *mongo.Collection, tenantBoundaryQuery *bson.M)  (*mongo.Cursor, error) {
	// Not Shown: Validation

	// Create query builder
	var qb epsearchast_v3.SemanticReducer[bson.D] = &LowerCaseEmail{}

	// Create Query Object
	queryObj, err := epsearchast_v3.SemanticReduceAst(ast, qb)

	if err != nil {
		return nil, err
	}

	mongoQuery := bson.D{
		{"$and",
			bson.A{
				tenantBoundaryQuery,
				queryObj,
			},
		}}
	
	return collection.Find(context.TODO(), mongoQuery)
}

type LowerCaseEmail struct {
	epsearchast_v3_mongo.DefaultMongoQueryBuilder
}

func (l *LowerCaseEmail) VisitEq(first, second string) (*bson.D, error) {
	if first == "email" {
		return &bson.D{{first, bson.D{{"$eq", strings.ToLower(second)}}}}, nil
	} else {
		return DefaultMongoQueryBuilder.VisitEq(l.DefaultMongoQueryBuilder, first, second)
	}
}
```


### FAQ

#### Design

##### Why does validation include alias resolution, why not process aliases first?

When validation errors occur, those errors go back to the user, so telling the user the error that occurred using the term they specified improves usability.
