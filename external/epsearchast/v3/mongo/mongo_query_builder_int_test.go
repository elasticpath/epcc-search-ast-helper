package epsearchast_v3_mongo

import (
	"context"
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
	"time"
)

var client *mongo.Client

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:20002"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	os.Exit(m.Run())
}

func TestSmokeTestMongoWithFilters(t *testing.T) {

	documents := []interface{}{
		bson.M{
			"string_field":          "test1",
			"array_field":           []string{"a", "b"},
			"nullable_string_field": nil,
		},
		bson.M{
			"string_field": "test2",
			"array_field":  []string{"c", "d"},

			"nullable_string_field": "yay",
		},
		bson.M{
			"string_field": "test3",
			"array_field":  []string{"c"},
			// No "nullable_string_field"
		},
	}

	var testCases = []struct {
		filter string
		count  int64
	}{
		{
			//language=JSON
			filter: `{
						"type": "EQ",
						"args": ["string_field", "test1"]
					}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
						"type": "GT",
						"args": ["string_field", "test1"]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "GE",
						"args": ["string_field", "test1"]
					}`,
			count: 3,
		},
		{
			//language=JSON
			filter: `{
						"type": "LE",
						"args": ["string_field", "test1"]
					}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
						"type": "LT",
						"args": ["string_field", "test1"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "LIKE",
						"args": ["string_field", "test"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "LIKE",
						"args": ["string_field", "test*"]
					}`,
			count: 3,
		},
		{
			//language=JSON
			filter: `{
						"type": "LIKE",
						"args": ["string_field", "Test*"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "ILIKE",
						"args": ["string_field", "test"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "ILIKE",
						"args": ["string_field", "test*"]
					}`,
			count: 3,
		},
		{
			//language=JSON
			filter: `{
						"type": "ILIKE",
						"args": ["string_field", "Test*"]
					}`,
			count: 3,
		},
		{
			//language=JSON
			filter: `{
						"type": "ILIKE",
						"args": ["string_field", "*EST3"]
					}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
						"type": "ILIKE",
						"args": ["string_field", "*Est*"]
					}`,
			count: 3,
		},
		{
			//language=JSON
			filter: `{
						"type": "IN",
						"args": ["string_field", "test1", "test2", "test4"]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "IS_NULL",
						"args": ["string_field"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "IS_NULL",
						"args": ["nullable_string_field"]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "CONTAINS",
						"args": ["array_field", "c"]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "CONTAINS",
						"args": ["array_field", "a"]
					}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
						"type": "CONTAINS",
						"args": ["array_field", "z"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "AND",
						"children": [
							{
								"type": "CONTAINS",
								"args": ["array_field", "c"]
							},
							{
								"type": "EQ",
								"args": ["string_field", "test2"]	
							}]
					}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
						"type": "AND",
						"children": [
							{
								"type": "CONTAINS",
								"args": ["array_field", "c"]
							},
							{
								"type": "EQ",
								"args": ["string_field", "test1"]	
							}]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "OR",
						"children": [
							{
								"type": "CONTAINS",
								"args": ["array_field", "c"]
							},
							{
								"type": "EQ",
								"args": ["string_field", "test1"]	
							}]
					}`,
			count: 3,
		},
		{
			//language=JSON
			filter: `{
			  "type": "AND",
			  "children": [
				{
				  "type": "OR",
				  "children": [
					{
					  "type": "CONTAINS",
					  "args": [
						"array_field",
						"c"
					  ]
					},
					{
					  "type": "EQ",
					  "args": [
						"string_field",
						"test1"
					  ]
					}
				  ]
				},
				{
				  "type": "EQ",
				  "args": [
					"nullable_string_field",
					"yay"
				  ]
				}
			  ]
			}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
  "type": "AND",
  "children": [
    {
      "type": "EQ",
      "args": [
        "string_field",
        "test1"
      ]
    },
    {
      "type": "AND",
      "children": [
        {
          "type": "EQ",
          "args": [
            "string_field",
            "test1"
          ]
        },
        {
          "type": "AND",
          "children": [
            {
              "type": "EQ",
              "args": [
                "string_field",
                "test1"
              ]
            },
            {
              "type": "AND",
              "children": [
                {
                  "type": "EQ",
                  "args": [
                    "string_field",
                    "test1"
                  ]
                },
                {
                  "type": "AND",
                  "children": [
                    {
                      "type": "EQ",
                      "args": [
                        "string_field",
                        "test1"
                      ]
                    },
                    {
                      "type": "AND",
                      "children": [
                        {
                          "type": "EQ",
                          "args": [
                            "string_field",
                            "test1"
                          ]
                        },
                        {
                          "type": "AND",
                          "children": [
                            {
                              "type": "EQ",
                              "args": [
                                "string_field",
                                "test1"
                              ]
                            },
                            {
                              "type": "AND",
                              "children": [
                                {
                                  "type": "EQ",
                                  "args": [
                                    "string_field",
                                    "test1"
                                  ]
                                },
                                {
                                  "type": "AND",
                                  "children": [
                                    {
                                      "type": "EQ",
                                      "args": [
                                        "string_field",
                                        "test1"
                                      ]
                                    },
                                    {
                                      "type": "AND",
                                      "children": [
                                        {
                                          "type": "EQ",
                                          "args": [
                                            "string_field",
                                            "test1"
                                          ]
                                        },
                                        {
                                          "type": "EQ",
                                          "args": [
                                            "string_field",
                                            "test1"
                                          ]
                                        }
                                      ]
                                    }
                                  ]
                                }
                              ]
                            }
                          ]
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
  "type": "OR",
  "children": [
    {
      "type": "EQ",
      "args": [
        "string_field",
        "test1"
      ]
    },
    {
      "type": "OR",
      "children": [
        {
          "type": "EQ",
          "args": [
            "string_field",
            "test1"
          ]
        },
        {
          "type": "OR",
          "children": [
            {
              "type": "EQ",
              "args": [
                "string_field",
                "test1"
              ]
            },
            {
              "type": "OR",
              "children": [
                {
                  "type": "EQ",
                  "args": [
                    "string_field",
                    "test1"
                  ]
                },
                {
                  "type": "OR",
                  "children": [
                    {
                      "type": "EQ",
                      "args": [
                        "string_field",
                        "test1"
                      ]
                    },
                    {
                      "type": "OR",
                      "children": [
                        {
                          "type": "EQ",
                          "args": [
                            "string_field",
                            "test1"
                          ]
                        },
                        {
                          "type": "OR",
                          "children": [
                            {
                              "type": "EQ",
                              "args": [
                                "string_field",
                                "test1"
                              ]
                            },
                            {
                              "type": "OR",
                              "children": [
                                {
                                  "type": "EQ",
                                  "args": [
                                    "string_field",
                                    "test1"
                                  ]
                                },
                                {
                                  "type": "OR",
                                  "children": [
                                    {
                                      "type": "EQ",
                                      "args": [
                                        "string_field",
                                        "test1"
                                      ]
                                    },
                                    {
                                      "type": "OR",
                                      "children": [
                                        {
                                          "type": "EQ",
                                          "args": [
                                            "string_field",
                                            "test1"
                                          ]
                                        },
                                        {
                                          "type": "EQ",
                                          "args": [
                                            "string_field",
                                            "test2"
                                          ]
                                        }
                                      ]
                                    }
                                  ]
                                }
                              ]
                            }
                          ]
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}`,
			count: 2,
		},
	}

	for _, tc := range testCases {
		ast, err := epsearchast_v3.GetAst(tc.filter)
		if err != nil {
			t.Fatalf("Failed to get filter: %v", err)
		}

		t.Run(fmt.Sprintf("%s", ast.AsFilter()), func(t *testing.T) {
			/*
				Fixture Setup
			*/
			ctx := context.Background()
			collection := SetupDB(t, ctx)
			InsertDocumentsOrFail(t, collection, ctx, documents)

			/*
			  Execute SUT
			*/

			// Perform a count query with a filter

			// Create query builder
			var qb epsearchast_v3.SemanticReducer[bson.D] = DefaultMongoQueryBuilder{}

			// Create Query Object
			ast, err := epsearchast_v3.GetAst(tc.filter)
			if err != nil {
				t.Fatalf("Failed to get filter: %v", err)
			}

			query, err := epsearchast_v3.SemanticReduceAst(ast, qb)

			if err != nil {
				t.Fatalf("Failed to get filter: %v", err)
			}

			/*
				Verification
			*/

			count, err := collection.CountDocuments(ctx, query)
			if err != nil {
				t.Fatalf("Failed to count documents: %v", err)
			}

			// Assert the expected count
			expectedCount := tc.count
			if count != expectedCount {
				t.Errorf("Expected count %d, but got %d", expectedCount, count)
			}

			fmt.Printf("Test passed. Documents matching filter: %d\n", count)

		})
		// Verification
	}

}

func InsertDocumentsOrFail(t *testing.T, collection *mongo.Collection, ctx context.Context, documents []interface{}) {
	_, err := collection.InsertMany(ctx, documents)
	if err != nil {
		t.Fatalf("Failed to insert test documents: %v", err)
	}
}

func SetupDB(t *testing.T, ctx context.Context) *mongo.Collection {
	db := client.Database("testdb")

	collName := t.Name()

	if len(collName) > 64 {
		collName = collName[0:64]
	}

	collection := db.Collection(collName)

	collection.Drop(ctx)

	return collection
}
