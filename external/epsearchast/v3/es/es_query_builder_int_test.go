package epsearchast_v3_es

import (
	"bytes"
	"encoding/json"
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

const esBaseURL = "http://localhost:20003"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type testStruct struct {
	filter string
	count  int64
}

func (t *testStruct) String() string {

	ast, err := epsearchast_v3.GetAst(t.filter)

	if err != nil {
		panic(fmt.Sprintf("Failed to get filter: %s)", err))
	}

	return ast.AsFilter()
}

func TestSmokeTestElasticSearchWithFilters(t *testing.T) {
	documents := []map[string]any{
		{
			"string_field":          "test1",
			"array_field":           []string{"a", "b"},
			"nullable_string_field": nil,
			"text_field":            "Developers like IDEs",
			"key_value_field": []map[string]any{
				{
					"alpha":       "a",
					"num":         "1",
					"roman":       "I",
					"description": "multiplicative identity",
					"array":       []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				},
				{
					"alpha":       "c",
					"num":         "3",
					"roman":       "III",
					"description": "trinity of triads",
					"array":       []int{0, 3, 6, 9},
				},
			},
		},
		{
			"string_field":          "test2",
			"array_field":           []string{"c", "d"},
			"nullable_string_field": "yay",
			"text_field":            "I like Development Environments",
			"key_value_field": []map[string]any{
				{
					"alpha":       "b",
					"num":         "2",
					"roman":       "II",
					"description": "two is a crowd",
					"array":       []int{0, 2, 4, 6, 8},
				},
				{
					"alpha": "c",
					"num":   "3",
					// Note this is lower case
					"roman":       "iii",
					"description": "Trifecta of Triangles.",
					"array":       []int{0, 3, 6, 9},
				},
				{
					"alpha": "f",
					"num":   "6",
					// Note this is mixed case.
					"roman":       "vI",
					"description": "Half a dozen or hexagon.",
					"array":       []int{0, 6},
				},
			},
		},
		{
			"string_field": "test3",
			"array_field":  []string{"c"},
			"text_field":   "Vim is the best",
		},
	}

	var testCases = []testStruct{
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
		{
			//language=JSON
			filter: `{
						"type": "EQ",
						"args": ["key_value_field.c.roman", "III"]
					}`,
			count: 2,
		},
		{
			//language=JSON

			filter: `{
						"type": "EQ",
						"args": ["key_value_field[3].roman", "III"]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "TEXT",
						"args": ["key_value_field.c.description", "trifecta"]
					}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
						"type": "CONTAINS",
						"args": ["key_value_field.c.array", "6"]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "CONTAINS",
						"args": ["key_value_field[2].array", "8"]
					}`,
			count: 1,
		},
		{
			//language=JSON

			filter: `{
						"type": "IN",
						"args": ["key_value_field.c.roman", "III", "II"]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "IS_NULL",
						"args": ["key_value_field.c.foo"]
					}`,
			// No document with index c, has a foo attribute.
			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "IS_NULL",
						"args": ["key_value_field.a.roman"]
					}`,
			// All documents with index A have a roman attribute
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "GT",
						"args": ["key_value_field.c.num", "3"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "GT",
						"args": ["key_value_field.c.num", "2"]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "GE",
						"args": ["key_value_field.c.num", "3"]
					}`,
			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "GE",
						"args": ["key_value_field.c.num", "4"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "LT",
						"args": ["key_value_field.c.num", "3"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "LE",
						"args": ["key_value_field.c.num", "4"]
					}`,

			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "LT",
						"args": ["key_value_field[3].alpha", "c"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "LE",
						"args": ["key_value_field[3].alpha", "c"]
					}`,

			count: 2,
		},
		{
			//language=JSON
			filter: `{
						"type": "LIKE",
						"args": ["key_value_field[1].description", "multiplicative"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "ILIKE",
						"args": ["key_value_field[1].description", "Multiplicative"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "LIKE",
						"args": ["key_value_field[1].description", "multiplicative*"]
					}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
						"type": "ILIKE",
						"args": ["key_value_field[1].description", "Multiplicative*"]
					}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
						"type": "EQ",
						"args": ["key_value_field[1].description", "multiplicative identity"]
					}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
						"type": "EQ",
						"args": ["key_value_field[1].description", "multiplicative identity 2"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "EQ",
						"args": ["key_value_field[1].description", "Multiplicative Identity"]
					}`,
			count: 0,
		},
		{
			//language=JSON
			filter: `{
						"type": "AND",
						"children": [
							{
								"type": "EQ",
								"args": ["string_field", "test1"]
							},
							{
								"type": "EQ",
								"args": ["key_value_field[3].roman", "III"]
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
								"args": ["string_field", "test1"]
							},
							{
								"type": "EQ",
								"args": ["key_value_field[3].roman", "III"]
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
							"type": "AND",
							"children": [
								{
									"type": "EQ",
									"args": ["string_field", "test1"]
								},
								{
									"type": "EQ",
									"args": ["key_value_field[3].roman", "III"]
								}
							]
							},
							{
								"type": "IN",
								"args": ["string_field", "test1", "test3"]
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
							"type": "OR",
							"children": [
								{
									"type": "EQ",
									"args": ["string_field", "test3"]
								},
								{
									"type": "EQ",
									"args": ["key_value_field[3].roman", "III"]
								}
							]
							},
							{
								"type": "IN",
								"args": ["string_field", "test4", "test3"]
							}
						]
					}`,
			count: 1,
		},
	}

	for _, tc := range testCases {
		// table tests are limited to this:
		// https://www.jetbrains.com/help/go/performing-tests.html#productivity-tips
		t.Run(fmt.Sprintf("%s", tc.filter), func(t *testing.T) {

			var indexName = "test_index"
			err := deleteIndex(indexName)
			if err != nil {
				t.Fatalf("Failed to delete index: %v", err)
			}

			err = createIndex(indexName)
			if err != nil {
				t.Fatalf("Failed to delete index: %v", err)
			}

			// Insert documents into Elasticsearch
			err = insertDocuments(indexName, documents)
			if err != nil {
				t.Fatalf("Failed to insert documents: %v", err)
			}

			// Build Elasticsearch query
			ast, err := epsearchast_v3.GetAst(tc.filter)
			if err != nil {
				t.Fatalf("Failed to parse filter: %v", err)
			}

			var qb epsearchast_v3.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{
				OpTypeToFieldNames: map[string]*OperatorTypeToMultiFieldName{
					"key_value_field.description": {
						Wildcard: "key_value_field.description.wildcard",
						Equality: "key_value_field.description.keyword",
					},
				},
				NestedFieldToQuery: map[string]NestedReplacement{
					// Treats key value field as an associative array (JSON Object with string keys).
					// This will take the <key> named capture group and replace it with a look up on the alpha field
					// and then the attribute is adjacent.
					`^key_value_field\.(?P<key>[^.]+)\.(?P<attribute>[^.]+)$`: {
						Path: "key_value_field",
						Subqueries: map[string]Replacement{
							"key_value_field.alpha":      {"$key", true},
							"key_value_field.$attribute": {"$value", false},
						},
					},
					`^key_value_field\[(?P<key>[^.]+)\].(?P<attribute>[^.]+)$`: {
						Path: "key_value_field",
						Subqueries: map[string]Replacement{
							"key_value_field.num":        {"$key", true},
							"key_value_field.$attribute": {"$value", false},
						},
					},
				},
			}

			query, err := epsearchast_v3.SemanticReduceAst(ast, qb)
			if err != nil {
				t.Fatalf("Failed to reduce AST: %v", err)
			}

			// Execute search query
			count, err := countDocuments(indexName, query)
			if err != nil {
				t.Fatalf("Failed to query Elasticsearch: %v", err)
			}

			// Assert the expected count
			if count != tc.count {
				txt, _ := json.MarshalIndent(query, "", "  ")
				t.Errorf("Expected count %d, but got %d with query\n%s", tc.count, count, txt)
			}
		})
	}
}

func insertDocuments(index string, documents []map[string]any) error {
	for _, doc := range documents {
		body, err := json.Marshal(doc)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/_doc", esBaseURL, index), bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("failed to insert document: %s", body)
		}
	}

	// Refresh the index
	_, err := http.Post(fmt.Sprintf("%s/%s/_refresh", esBaseURL, index), "application/json", nil)
	return err
}

func countDocuments(index string, query *JsonObject) (int64, error) {
	queryBody := map[string]any{
		"query": query,
	}
	body, err := json.Marshal(queryBody)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(fmt.Sprintf("%s/%s/_search", esBaseURL, index), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("failed to query Elasticsearch: %s", body)
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
		} `json:"hits"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Hits.Total.Value, nil
}

func createIndex(indexName string) error {
	mapping := map[string]any{
		"mappings": map[string]any{
			"dynamic": "strict",
			"properties": map[string]any{
				"array_field": map[string]any{
					"type": "text",
					"fields": map[string]any{
						"keyword": map[string]any{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},
				"nullable_string_field": map[string]any{
					"type": "text",
					"fields": map[string]any{
						"keyword": map[string]any{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},
				"string_field": map[string]any{
					"type": "text",
					"fields": map[string]any{
						"keyword": map[string]any{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},
				"text_field": map[string]any{
					"type":     "text",
					"analyzer": "english", // Enables stemming
				},
				"key_value_field": map[string]any{
					"type": "nested",
					"properties": map[string]any{
						"alpha": map[string]any{
							"type": "keyword",
						},
						"num": map[string]any{
							"type": "long",
						},
						"roman": map[string]any{
							"type":       "keyword",
							"normalizer": "uppercase_normalizer",
						},
						"description": map[string]any{
							"type":     "text",
							"analyzer": "english", // Enables stemming
							"fields": map[string]any{
								"wildcard": map[string]any{
									"type": "wildcard",
								},
								// You probably could use wildcard for this.
								"keyword": map[string]any{
									"type": "keyword",
								},
							},
						},
						"array": map[string]any{
							"type": "integer",
						},
					},
				},
			},
		},
		"settings": map[string]any{
			"analysis": map[string]any{
				"normalizer": map[string]any{
					"uppercase_normalizer": map[string]any{
						"type":   "custom",
						"filter": []string{"uppercase"},
					},
				},
			},
		},
	}

	url := fmt.Sprintf("%s/%s", esBaseURL, indexName)

	body, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("failed to marshal mapping: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create index, status: %s", resp.Status)
	}

	fmt.Printf("Index %s created successfully\n", indexName)
	return nil
}

func deleteIndex(index string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", esBaseURL, index), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Ignore 404 status, as the index may not exist
	if resp.StatusCode != 404 && resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete index: %s", body)
	}

	return nil
}
