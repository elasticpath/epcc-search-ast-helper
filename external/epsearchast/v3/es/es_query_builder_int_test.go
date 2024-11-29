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

func TestSmokeTestElasticSearchWithFilters(t *testing.T) {
	documents := []map[string]interface{}{
		{
			"string_field":          "test1",
			"array_field":           []string{"a", "b"},
			"nullable_string_field": nil,
			"text_field":            "Developers like IDEs",
		},
		{
			"string_field":          "test2",
			"array_field":           []string{"c", "d"},
			"nullable_string_field": "yay",
			"text_field":            "I like Development Environments",
		},
		{
			"string_field": "test3",
			"array_field":  []string{"c"},
			"text_field":   "Vim is the best",
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
	}

	for _, tc := range testCases {

		t.Run(tc.filter, func(t *testing.T) {

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

			var qb epsearchast_v3.SemanticReducer[JsonObject] = DefaultEsQueryBuilder{}
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
				t.Errorf("Expected count %d, but got %d", tc.count, count)
			}
		})
	}
}

func insertDocuments(index string, documents []map[string]interface{}) error {
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
		resp.Body.Close()

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
	queryBody := map[string]interface{}{
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
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"array_field": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},
				"nullable_string_field": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},
				"string_field": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},
				"text_field": map[string]interface{}{
					"type":     "text",
					"analyzer": "english", // Enables stemming
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
