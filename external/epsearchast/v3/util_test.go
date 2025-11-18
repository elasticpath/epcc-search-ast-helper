package epsearchast_v3

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDepthOfSingleElement(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := fmt.Sprintf(`
{
  "type": "EQ",
  "args": [
    "amount",
    "5"
  ]
}
`)

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT

	depth := GetAstDepth(ast)

	// Verification
	require.Equal(t, 1, depth)

}

func TestDepthOfNestedAndElement(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := fmt.Sprintf(`{
  "type": "AND",
  "children": [
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "status",
        "paid"
      ]
    }
  ]
}
`)

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT

	depth := GetAstDepth(ast)

	// Verification
	require.Equal(t, 2, depth)

}

func TestDepthOfNestedOrElement(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := fmt.Sprintf(`{
  "type": "OR",
  "children": [
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "status",
        "paid"
      ]
    }
  ]
}
`)

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT

	depth := GetAstDepth(ast)

	// Verification
	require.Equal(t, 2, depth)

}

func TestDepthOfNestedOrAndAndThatIsNotBalancedElement(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := fmt.Sprintf(`{
  "type": "OR",
  "children": [
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "OR",
      "children": [
        {
          "type": "EQ",
          "args": [
            "amount",
            "5"
          ]
        },
        {
          "type": "AND",
          "children": [
            {
              "type": "EQ",
              "args": [
                "status",
                "paid"
              ]
            },
            {
              "type": "EQ",
              "args": [
                "shipping",
                "unfulfilled"
              ]
            }
          ]
        }
      ]
    },
  	{
      "type": "EQ",
      "args": [
        "account_id",
        "67"
      ]
    }
  ]
}
  
`)

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT

	depth := GetAstDepth(ast)

	// Verification
	require.Equal(t, 4, depth)

}

func TestEffectiveIndexIntersectionCountOfSingleElement(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := fmt.Sprintf(`
{
  "type": "EQ",
  "args": [
    "amount",
    "5"
  ]
}
`)

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT

	indexIntersectionCount, err := GetEffectiveIndexIntersectionCount(ast)

	// Verification
	require.NoError(t, err)
	require.Equal(t, uint64(1), indexIntersectionCount)

}

func TestEffectiveIndexIntersectionCountOfAndElement(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := fmt.Sprintf(`{
  "type": "AND",
  "children": [
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "status",
        "paid"
      ]
    }
  ]
}
`)

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT

	indexIntersectionCount, err := GetEffectiveIndexIntersectionCount(ast)

	// Verification
	require.NoError(t, err)
	require.Equal(t, uint64(1), indexIntersectionCount)

}

func TestEffectiveIndexIntersectionCountOfNestedOrElement(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := fmt.Sprintf(`{
  "type": "AND",
  "children": [
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "OR",
      "children": [
        {
          "type": "EQ",
          "args": [
            "amount",
            "5"
          ]
        },
        {
          "type": "EQ",
          "args": [
            "status",
            "paid"
          ]
        },
        {
          "type": "EQ",
          "args": [
            "shipping_status",
            "unfulfilled"
          ]
        }
      ]
    }
  ]
}
`)

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT

	indexIntersectionCount, err := GetEffectiveIndexIntersectionCount(ast)

	// Verification

	require.NoError(t, err)
	require.Equal(t, uint64(3), indexIntersectionCount)

}

func TestGetAllFirstArgsWithSingleElement(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "EQ",
  "args": [
    "amount",
    "5"
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	fields := GetAllFirstArgs(ast)

	// Verification
	require.Equal(t, []string{"amount"}, fields)
}

func TestGetAllFirstArgsWithBinaryOperators(t *testing.T) {
	operators := []string{"EQ", "LT", "LE", "GT", "GE", "LIKE", "ILIKE", "CONTAINS", "TEXT"}

	for _, op := range operators {
		t.Run(op, func(t *testing.T) {
			// Fixture Setup
			jsonTxt := fmt.Sprintf(`
{
  "type": "%s",
  "args": [
    "status",
    "paid"
  ]
}
`, op)

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			fields := GetAllFirstArgs(ast)

			// Verification
			require.Equal(t, []string{"status"}, fields)
		})
	}
}

func TestGetAllFirstArgsWithUnaryOperator(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "IS_NULL",
  "args": [
    "deleted_at"
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	fields := GetAllFirstArgs(ast)

	// Verification
	require.Equal(t, []string{"deleted_at"}, fields)
}

func TestGetAllFirstArgsWithVarargOperators(t *testing.T) {
	operators := []string{"IN", "CONTAINS_ANY", "CONTAINS_ALL"}

	for _, op := range operators {
		t.Run(op, func(t *testing.T) {
			// Fixture Setup
			jsonTxt := fmt.Sprintf(`
{
  "type": "%s",
  "args": [
    "state",
    "CA",
    "NY",
    "TX"
  ]
}
`, op)

			ast, err := GetAst(jsonTxt)
			require.NoError(t, err)

			// Execute SUT
			fields := GetAllFirstArgs(ast)

			// Verification
			// Should only return the first arg (field name), not all the values
			require.Equal(t, []string{"state"}, fields)
		})
	}
}

func TestGetAllFirstArgsWithAndOperator(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "AND",
  "children": [
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "status",
        "paid"
      ]
    }
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	fields := GetAllFirstArgs(ast)

	// Verification
	require.ElementsMatch(t, []string{"amount", "status"}, fields)
}

func TestGetAllFirstArgsWithOrOperator(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "OR",
  "children": [
    {
      "type": "EQ",
      "args": [
        "payment",
        "credit_card"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "payment",
        "paypal"
      ]
    }
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	fields := GetAllFirstArgs(ast)

	// Verification
	// Should include duplicates
	require.Equal(t, []string{"payment", "payment"}, fields)
}

func TestGetAllFirstArgsWithNestedStructure(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "OR",
  "children": [
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "AND",
      "children": [
        {
          "type": "EQ",
          "args": [
            "status",
            "paid"
          ]
        },
        {
          "type": "GE",
          "args": [
            "created_at",
            "2024-01-01"
          ]
        }
      ]
    },
    {
      "type": "IN",
      "args": [
        "shipping",
        "express",
        "standard"
      ]
    }
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	fields := GetAllFirstArgs(ast)

	// Verification
	require.ElementsMatch(t, []string{"amount", "status", "created_at", "shipping"}, fields)
}

func TestGetAllFirstArgsWithNilAst(t *testing.T) {
	// Execute SUT
	fields := GetAllFirstArgs(nil)

	// Verification
	require.Equal(t, []string{}, fields)
}

func TestGetAllFirstArgsSortedWithMultipleFields(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "AND",
  "children": [
    {
      "type": "EQ",
      "args": [
        "zebra",
        "1"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "apple",
        "2"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "banana",
        "3"
      ]
    }
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	fields := GetAllFirstArgsSorted(ast)

	// Verification
	require.Equal(t, []string{"apple", "banana", "zebra"}, fields)
}

func TestGetAllFirstArgsSortedWithDuplicates(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "OR",
  "children": [
    {
      "type": "EQ",
      "args": [
        "status",
        "paid"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "status",
        "pending"
      ]
    }
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	fields := GetAllFirstArgsSorted(ast)

	// Verification
	// Should include duplicates, but sorted
	require.Equal(t, []string{"amount", "status", "status"}, fields)
}

func TestGetAllFirstArgsUniqueWithDuplicates(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "OR",
  "children": [
    {
      "type": "EQ",
      "args": [
        "status",
        "paid"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "status",
        "pending"
      ]
    }
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	fields := GetAllFirstArgsUnique(ast)

	// Verification
	expected := map[string]struct{}{
		"status": {},
		"amount": {},
	}
	require.Equal(t, expected, fields)
}

func TestGetAllFirstArgsUniqueWithNestedStructure(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "AND",
  "children": [
    {
      "type": "OR",
      "children": [
        {
          "type": "EQ",
          "args": [
            "amount",
            "5"
          ]
        },
        {
          "type": "GE",
          "args": [
            "amount",
            "10"
          ]
        }
      ]
    },
    {
      "type": "EQ",
      "args": [
        "status",
        "paid"
      ]
    },
    {
      "type": "IS_NULL",
      "args": [
        "deleted_at"
      ]
    }
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	fields := GetAllFirstArgsUnique(ast)

	// Verification
	expected := map[string]struct{}{
		"amount":     {},
		"status":     {},
		"deleted_at": {},
	}
	require.Equal(t, expected, fields)
}

func TestGetAllFirstArgsUniqueWithNilAst(t *testing.T) {
	// Execute SUT
	fields := GetAllFirstArgsUnique(nil)

	// Verification
	require.Equal(t, map[string]struct{}{}, fields)
}

func TestHasFirstArgWithSingleElementFound(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "EQ",
  "args": [
    "amount",
    "5"
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	found := HasFirstArg(ast, "amount")

	// Verification
	require.True(t, found)
}

func TestHasFirstArgWithSingleElementNotFound(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "EQ",
  "args": [
    "amount",
    "5"
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	found := HasFirstArg(ast, "status")

	// Verification
	require.False(t, found)
}

func TestHasFirstArgWithNestedStructureFound(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "AND",
  "children": [
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "OR",
      "children": [
        {
          "type": "EQ",
          "args": [
            "status",
            "paid"
          ]
        },
        {
          "type": "GE",
          "args": [
            "created_at",
            "2024-01-01"
          ]
        }
      ]
    }
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	foundAmount := HasFirstArg(ast, "amount")
	foundStatus := HasFirstArg(ast, "status")
	foundCreatedAt := HasFirstArg(ast, "created_at")

	// Verification
	require.True(t, foundAmount)
	require.True(t, foundStatus)
	require.True(t, foundCreatedAt)
}

func TestHasFirstArgWithNestedStructureNotFound(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "AND",
  "children": [
    {
      "type": "EQ",
      "args": [
        "amount",
        "5"
      ]
    },
    {
      "type": "EQ",
      "args": [
        "status",
        "paid"
      ]
    }
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	found := HasFirstArg(ast, "email")

	// Verification
	require.False(t, found)
}

func TestHasFirstArgWithNilAst(t *testing.T) {
	// Execute SUT
	found := HasFirstArg(nil, "amount")

	// Verification
	require.False(t, found)
}

func TestHasFirstArgWithVarargOperator(t *testing.T) {
	// Fixture Setup
	// language=JSON
	jsonTxt := `
{
  "type": "IN",
  "args": [
    "state",
    "CA",
    "NY",
    "TX"
  ]
}
`

	ast, err := GetAst(jsonTxt)
	require.NoError(t, err)

	// Execute SUT
	foundState := HasFirstArg(ast, "state")
	foundCA := HasFirstArg(ast, "CA")

	// Verification
	require.True(t, foundState)
	require.False(t, foundCA) // CA is not a first arg, it's a value
}
