package v3_gorm

import (
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleEqFilterGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "EQ",
				"first_arg": "amount",
				"second_arg": "5"
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "LOWER(amount::text) = LOWER(?)", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{"5"})
}

func TestSimpleLeFilterGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "LE",
				"first_arg": "amount",
				"second_arg": "5"
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "amount <= ?", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{"5"})
}

func TestSimpleLtFilterGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "LT",
				"first_arg": "amount",
				"second_arg": "5"
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "amount < ?", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{"5"})
}

func TestSimpleGeFilterGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "GE",
				"first_arg": "amount",
				"second_arg": "5"
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "amount >= ?", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{"5"})
}

func TestSimpleGtFilterGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "GT",
				"first_arg": "amount",
				"second_arg": "5"
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "amount > ?", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{"5"})
}

func TestSimpleInFilterGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "IN",
				"args": ["status", "new", "paid"]
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "status IN ?", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{[]interface{}{"new", "paid"}})
}

func TestSimpleLikeFilterGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "LIKE",
				"first_arg": "amount",
				"second_arg": "5"
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "amount ILIKE ?", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{"5"})
}

func TestSimpleLikeFilterWithWildcardAtStartGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "LIKE",
				"first_arg": "amount",
				"second_arg": "*5"
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "amount ILIKE ?", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{"%5"})
}

func TestSimpleLikeFilterWithWildcardAtEndGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "LIKE",
				"first_arg": "amount",
				"second_arg": "5*"
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "amount ILIKE ?", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{"5%"})
}

func TestSimpleLikeFilterWithWildcardAtBothEndsGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "LIKE",
				"first_arg": "amount",
				"second_arg": "*5*"
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "amount ILIKE ?", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{"%5%"})
}

func TestSimpleLikeFilterWithWildcardOnlyCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
			{
				"type": "LIKE",
				"first_arg": "amount",
				"second_arg": "*"
			}
			`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "amount ILIKE ?", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{"%"})
}

func TestSimpleAndFilterGeneratesCorrectWhereClause(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
				{ 
					"type":  "AND",
					"children": [
					{
						"type": "IN",
						"args": ["status", "new", "paid"]
					},
					{
						"type": "GE",
						"first_arg": "amount",
						"second_arg": "5"
					}
					]
				}
				`

	astNode, err := epsearchast_v3.GetAst(jsonTxt)
	require.NoError(t, err)

	gormVisitor := NewGormVisitor()
	visitor := epsearchast_v3.NewSearchFilterVisitorAdapter(gormVisitor)

	// Execute SUT
	err = astNode.Accept(visitor)

	// Verification
	require.NoError(t, err)

	require.Equal(t, "(status IN ? AND amount >= ?)", gormVisitor.Clause)
	require.Equal(t, gormVisitor.Args, []interface{}{[]interface{}{"new", "paid"}, "5"})
}
