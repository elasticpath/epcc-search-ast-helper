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
