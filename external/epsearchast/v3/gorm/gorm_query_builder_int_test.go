package epsearchast_v3_gorm

import (
	"context"
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
	"os"
	"testing"
)

var postgresDB *gorm.DB

func TestMain(m *testing.M) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var err error
	dsn := "host=localhost user=admin password=admin dbname=test_db port=20001 sslmode=disable"

	postgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	os.Exit(m.Run())
}

var testName string = ""

type TestTable struct {
	ID                  int64          `gorm:"primaryKey"`
	StringField         string         `gorm:"type:varchar(255)"`
	NullableStringField *string        `gorm:"type:varchar(255)"`
	ArrayField          pq.StringArray `gorm:"type:text[]"`
	TextField           string         `gorm:"type:text"`
}

func (a *TestTable) TableName() string {
	return testName
}

func TestSmokeTestPostgresWithFilters(t *testing.T) {

	yay := "yay"
	documents := []TestTable{
		{
			StringField:         "test1",
			ArrayField:          []string{"a", "b"},
			NullableStringField: nil,
			TextField:           "Developers like IDEs",
		}, {
			StringField:         "test2",
			ArrayField:          []string{"c", "d"},
			NullableStringField: &yay,
			TextField:           "I like Development Environments",
		}, {
			StringField: "test3",
			ArrayField:  []string{"c"},
			// No "nullable_string_field"
			TextField: "Vim is the best",
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
						"type": "TEXT",
						"args": ["text_field", "developers"]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
					  "type": "OR",
					  "children": [
						{
						  "type": "CONTAINS",
						  "args": [
							"array_field",
							"d"
						  ]
						},
						{
						  "type": "CONTAINS",
						  "args": [
							"array_field",
							"a"
						  ]
						}
					  ]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
					  "type": "OR",
					  "children": [
						{
						  "type": "TEXT",
						  "args": [
							"text_field",
							"developers"
						  ]
						},
						{
						  "type": "CONTAINS",
						  "args": [
							"array_field",
							"c"
						  ]
						}
					  ]
					}`,
			count: 3,
		},
		{
			//language=JSON
			filter: `{
					  "type": "OR",
					  "children": [
						{
						  "type": "TEXT",
						  "args": [
							"text_field",
							"developers"
						  ]
						},
						{
						  "type": "AND",
						  "children": [
							{
							  "type": "EQ",
							  "args": [
								"nullable_string_field",
								"yay"
							  ]
							},
							{
							  "type": "CONTAINS",
							  "args": [
								"array_field",
								"c"
							  ]
							}
						  ]
						}
					  ]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
					  "type": "AND",
					  "children": [
						{
						  "type": "TEXT",
						  "args": [
							"text_field",
							"developers"
						  ]
						},
						{
						  "type": "OR",
						  "children": [
							{
							  "type": "EQ",
							  "args": [
								"nullable_string_field",
								"yay"
							  ]
							},
							{
							  "type": "CONTAINS",
							  "args": [
								"array_field",
								"b"
							  ]
							}
						  ]
						}
					  ]
					}`,
			count: 2,
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
			testName = t.Name()
			/*
				Fixture Setup
			*/
			ctx := context.Background()
			SetupDB(t, ctx, postgresDB)
			InsertDocumentsOrFail(t, postgresDB, documents)

			/*
			  Execute SUT
			*/

			// Perform a count query with a filter

			// Create query builder
			var qb epsearchast_v3.SemanticReducer[SubQuery] = DefaultGormQueryBuilder{}

			// Create Query Object
			ast, err := epsearchast_v3.GetAst(tc.filter)
			if err != nil {
				t.Fatalf("Failed to get filter: %v", err)
			}

			query, err := epsearchast_v3.SemanticReduceAst(ast, qb)

			if err != nil {
				t.Fatalf("Failed to convert filter: %v", err)
			}

			/*
				Verification
			*/

			var count int64
			err = postgresDB.Model(&TestTable{}).Where(query.Clause, query.Args...).Count(&count).Error
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

func InsertDocumentsOrFail(t *testing.T, db *gorm.DB, documents []TestTable) {

	for _, doc := range documents {
		if err := db.Create(&doc).Error; err != nil {
			t.Fatalf("Failed to insert test documents: %v", err)
		}
	}

}

func SetupDB(t *testing.T, ctx context.Context, db *gorm.DB) {

	tt := TestTable{}

	err := db.AutoMigrate(&tt)
	if err != nil {
		t.Fatalf("Failed to migrate table: %v", err)
	}

	err = db.Delete(&tt, "1 = 1").Error

	if err != nil {
		t.Fatalf("Failed to clear data: %v", err)
	}

}
