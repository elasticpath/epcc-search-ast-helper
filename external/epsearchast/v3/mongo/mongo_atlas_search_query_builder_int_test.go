package epsearchast_v3_mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestSmokeTestAtlasSearchWithFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Connect to Atlas Search enabled MongoDB
	ctx := context.Background()
	atlasClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:20004/?replicaSet=rs0&directConnection=true"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB Atlas Search: %v", err)
	}
	defer func() {
		if err := atlasClient.Disconnect(ctx); err != nil {
			t.Logf("Failed to disconnect: %v", err)
		}
	}()

	documents := []interface{}{
		bson.M{
			"string_field":          "test1",
			"array_field":           []string{"a", "b"},
			"nullable_string_field": nil,
			"text_field":            "Developers like IDEs",
		},
		bson.M{
			"string_field":          "test2",
			"array_field":           []string{"c", "d"},
			"nullable_string_field": "yay",
			"text_field":            "I like Development Environments",
		},
		bson.M{
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
						"type": "TEXT",
						"args": ["text_field", "like"]
					}`,
			count: 2,
		},
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
						"type": "EQ",
						"args": ["string_field", "test2"]
					}`,
			count: 1,
		},
		{
			// Test that EQ does exact matching, not partial matching
			//language=JSON
			filter: `{
						"type": "EQ",
						"args": ["string_field", "test"]
					}`,
			count: 0,
		},
		{
			// Test EQ on text_field - should match exact string, no stemming
			//language=JSON
			filter: `{
						"type": "EQ",
						"args": ["text_field", "Developers like IDEs"]
					}`,
			count: 1,
		},
		{
			// Test EQ on text_field - should be case sensitive
			//language=JSON
			filter: `{
						"type": "EQ",
						"args": ["text_field", "developers like ides"]
					}`,
			count: 0,
		},
		{
			// Test EQ on text_field - should not match partial/stemmed
			//language=JSON
			filter: `{
						"type": "EQ",
						"args": ["text_field", "Developers"]
					}`,
			count: 0,
		},
		{
			// Test TEXT on text_field - should use stemming and match
			//language=JSON
			filter: `{
						"type": "TEXT",
						"args": ["text_field", "developer"]
					}`,
			count: 2,
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
						"args": ["string_field", "*est1"]
					}`,
			count: 1,
		},
		{
			// Test IN operator - multiple values
			//language=JSON
			filter: `{
						"type": "IN",
						"args": ["string_field", "test1", "test2", "test4"]
					}`,
			count: 2,
		},
		{
			// Test IN operator - single value
			//language=JSON
			filter: `{
						"type": "IN",
						"args": ["string_field", "test3"]
					}`,
			count: 1,
		},
		{
			// Test IN operator - no matches
			//language=JSON
			filter: `{
						"type": "IN",
						"args": ["string_field", "test4", "test5"]
					}`,
			count: 0,
		},
		{
			// Test AND operator - two EQ conditions
			//language=JSON
			filter: `{
						"type": "AND",
						"children": [
							{"type": "EQ", "args": ["string_field", "test1"]},
							{"type": "EQ", "args": ["text_field", "Developers like IDEs"]}
						]
					}`,
			count: 1,
		},
		{
			// Test AND operator - EQ and TEXT
			//language=JSON
			filter: `{
						"type": "AND",
						"children": [
							{"type": "EQ", "args": ["string_field", "test2"]},
							{"type": "TEXT", "args": ["text_field", "Development"]}
						]
					}`,
			count: 1,
		},
		{
			// Test AND operator - no matches (impossible condition)
			//language=JSON
			filter: `{
						"type": "AND",
						"children": [
							{"type": "EQ", "args": ["string_field", "test1"]},
							{"type": "EQ", "args": ["string_field", "test2"]}
						]
					}`,
			count: 0,
		},
		{
			// Test OR operator - two EQ conditions
			//language=JSON
			filter: `{
						"type": "OR",
						"children": [
							{"type": "EQ", "args": ["string_field", "test1"]},
							{"type": "EQ", "args": ["string_field", "test2"]}
						]
					}`,
			count: 2,
		},
		{
			// Test OR operator - three conditions
			//language=JSON
			filter: `{
						"type": "OR",
						"children": [
							{"type": "EQ", "args": ["string_field", "test1"]},
							{"type": "EQ", "args": ["string_field", "test2"]},
							{"type": "EQ", "args": ["string_field", "test3"]}
						]
					}`,
			count: 3,
		},
		{
			// Test complex: (test1 OR test2) AND has "like" in text
			//language=JSON
			filter: `{
						"type": "AND",
						"children": [
							{
								"type": "OR",
								"children": [
									{"type": "EQ", "args": ["string_field", "test1"]},
									{"type": "EQ", "args": ["string_field", "test2"]}
								]
							},
							{"type": "TEXT", "args": ["text_field", "like"]}
						]
					}`,
			count: 2,
		},
		{
			// Test complex: (ILIKE wildcard) AND (IN multiple values)
			//language=JSON
			filter: `{
						"type": "AND",
						"children": [
							{"type": "ILIKE", "args": ["string_field", "test*"]},
							{"type": "IN", "args": ["string_field", "test1", "test2"]}
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
			collection := SetupAtlasDB(t, ctx, atlasClient)
			InsertDocumentsOrFail(t, collection, ctx, documents)

			// Create search index with explicit field mappings
			// Index fields with multiple types (like ES multi-fields):
			// - type: "string" for text/wildcard search
			// - type: "token" for exact matching with equals operator
			searchIndexModel := mongo.SearchIndexModel{
				Definition: bson.D{
					{"mappings", bson.D{
						{"dynamic", false},
						{"fields", bson.D{
							// string_field: indexed as both string (for wildcard) and token (for equals)
							{"string_field", bson.A{
								bson.D{{"type", "string"}},
								bson.D{{"type", "token"}},
							}},
							{"array_field", bson.D{
								{"type", "string"},
							}},
							{"nullable_string_field", bson.A{
								bson.D{{"type", "string"}},
								bson.D{{"type", "token"}},
							}},
							// text_field: indexed for both text search (with english analyzer) and exact matching
							{"text_field", bson.A{
								bson.D{
									{"type", "string"},
									{"analyzer", "lucene.english"},
								},
								bson.D{
									{"type", "token"},
									{"normalizer", "none"}, // case-sensitive exact matching
								},
							}},
						}},
					}},
				},
				Options: nil,
			}

			indexName, err := collection.SearchIndexes().CreateOne(ctx, searchIndexModel)
			if err != nil {
				t.Fatalf("Failed to create search index: %v", err)
			}
			t.Logf("Created search index: %s", indexName)

			// Wait for index to be ready by polling
			err = waitForAtlasSearchIndex(ctx, collection, indexName, 30*time.Second)
			if err != nil {
				t.Fatalf("Failed to wait for search index: %v", err)
			}
			t.Logf("Search index %s is ready", indexName)

			/*
			  Execute SUT
			*/

			// Create query builder
			// Atlas Search automatically routes operators to the correct index type
			var qb epsearchast_v3.SemanticReducer[bson.D] = DefaultAtlasSearchQueryBuilder{}

			// Create Query Object
			ast, err := epsearchast_v3.GetAst(tc.filter)
			if err != nil {
				t.Fatalf("Failed to get filter: %v", err)
			}

			query, err := epsearchast_v3.SemanticReduceAst(ast, qb)

			if err != nil {
				t.Fatalf("Failed to reduce AST: %v", err)
			}

			/*
				Verification
			*/

			// Execute the search using aggregation pipeline
			pipeline := mongo.Pipeline{
				{{Key: "$search", Value: query}},
			}

			cursor, err := collection.Aggregate(ctx, pipeline)
			if err != nil {
				t.Fatalf("Failed to execute search: %v", err)
			}
			defer cursor.Close(ctx)

			// Count results
			var results []bson.M
			err = cursor.All(ctx, &results)
			if err != nil {
				t.Fatalf("Failed to read results: %v", err)
			}

			count := int64(len(results))

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

func SetupAtlasDB(t *testing.T, ctx context.Context, atlasClient *mongo.Client) *mongo.Collection {
	db := atlasClient.Database("testdb")

	collName := t.Name()

	if len(collName) > 64 {
		collName = collName[0:64]
	}

	collection := db.Collection(collName)

	collection.Drop(ctx)

	return collection
}

// waitForAtlasSearchIndex polls until the search index appears in the list
// Note: MongoDB Community Search (used in docker) doesn't provide status like Atlas Cloud,
// so we just wait for the index to appear in the list and give it a moment to stabilize
func waitForAtlasSearchIndex(ctx context.Context, collection *mongo.Collection, indexName string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	pollInterval := 100 * time.Millisecond
	indexFound := false

	for time.Now().Before(deadline) {
		// List all search indexes
		cursor, err := collection.SearchIndexes().List(ctx, options.SearchIndexes())
		if err != nil {
			return fmt.Errorf("failed to list search indexes: %w", err)
		}

		// Find the index we're waiting for
		var indexes []bson.M
		if err := cursor.All(ctx, &indexes); err != nil {
			return fmt.Errorf("failed to read search indexes: %w", err)
		}

		for _, index := range indexes {
			name, ok := index["name"].(string)
			if !ok {
				continue
			}

			if name == indexName {
				// Check for status field (Atlas Cloud has this)
				if status, ok := index["status"].(string); ok {
					// If we have status, use it
					if status == "READY" {
						return nil
					}
					if status == "FAILED" {
						return fmt.Errorf("search index %s failed to build", indexName)
					}
					// Index is still building (INITIAL or BUILDING)
					indexFound = true
					break
				} else {
					// No status field (MongoDB Community Search)
					// Index exists in list, so it should be ready
					// Wait a bit longer to let it stabilize
					if !indexFound {
						indexFound = true
						time.Sleep(2 * time.Second)
					}
					return nil
				}
			}
		}

		// Wait before polling again
		time.Sleep(pollInterval)
	}

	if indexFound {
		return fmt.Errorf("search index %s found but never became ready", indexName)
	}
	return fmt.Errorf("timeout waiting for search index %s to appear", indexName)
}
