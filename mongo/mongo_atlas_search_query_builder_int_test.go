package astmongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/elasticpath/epcc-search-ast-helper"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestSmokeTestAtlasSearchWithFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Connect to Atlas Search enabled MongoDB
	ctx := context.Background()
	atlasClient, err := mongo.Connect(options.Client().ApplyURI("mongodb://admin:admin@localhost:20004/?replicaSet=rs0&directConnection=true"))
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
			"string_field":          "test1 test1",
			"array_field":           []string{"a a", "b b"},
			"nullable_string_field": nil,
			"text_field":            "Developers like IDEs",
			"uuid_field":            "550e8400-e29b-41d4-a716-446655440001",
			"date_field":            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		bson.M{
			"string_field":          "test2 test2",
			"array_field":           []string{"c c", "d d"},
			"nullable_string_field": "yay yay",
			"text_field":            "I like Development Environments",
			"uuid_field":            "550e8400-e29b-41d4-a716-446655440002",
			"date_field":            time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
		},
		bson.M{
			"string_field": "test3 test3",
			"array_field":  []string{"c c"},
			"text_field":   "Vim is the best",
			"uuid_field":   "550e8400-e29b-41d4-a716-446655440003",
			"date_field":   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
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
						"args": ["string_field", "test1 test1"]
					}`,
			count: 1,
		},
		{
			//language=JSON
			filter: `{
						"type": "EQ",
						"args": ["string_field", "test2 test2"]
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
						"args": ["string_field", "*test1"]
					}`,
			count: 1,
		},
		{
			// Test ILIKE with wildcard matching across space - pattern with space and wildcard
			//language=JSON
			filter: `{
						"type": "ILIKE",
						"args": ["string_field", "test1 *"]
					}`,
			count: 1,
		},
		{
			// Test ILIKE with wildcard matching across space - wildcard before space
			//language=JSON
			filter: `{
						"type": "ILIKE",
						"args": ["string_field", "* test1"]
					}`,
			count: 1,
		},
		{
			// Test ILIKE with wildcard at both beginning and end - should match "test1 test1"
			//language=JSON
			filter: `{
						"type": "ILIKE",
						"args": ["string_field", "*1 test*"]
					}`,
			count: 1,
		},
		{
			// Test LIKE (case-sensitive) with exact case match
			//language=JSON
			filter: `{
						"type": "LIKE",
						"args": ["string_field", "test*"]
					}`,
			count: 3,
		},
		{
			// Test LIKE (case-sensitive) with wrong case - should NOT match
			//language=JSON
			filter: `{
						"type": "LIKE",
						"args": ["string_field", "Test*"]
					}`,
			count: 0,
		},
		{
			// Test LIKE (case-sensitive) with wildcard at end
			//language=JSON
			filter: `{
						"type": "LIKE",
						"args": ["string_field", "test1 *"]
					}`,
			count: 1,
		},
		{
			// Test LIKE (case-sensitive) with wildcard at beginning
			//language=JSON
			filter: `{
						"type": "LIKE",
						"args": ["string_field", "*test1"]
					}`,
			count: 1,
		},
		{
			// Test LIKE (case-sensitive) with wildcard at both ends
			//language=JSON
			filter: `{
						"type": "LIKE",
						"args": ["string_field", "*1 test*"]
					}`,
			count: 1,
		},
		{
			// Test IN operator - multiple values
			//language=JSON
			filter: `{
						"type": "IN",
						"args": ["string_field", "test1 test1", "test2 test2", "test4"]
					}`,
			count: 2,
		},
		{
			// Test IN operator - single value
			//language=JSON
			filter: `{
						"type": "IN",
						"args": ["string_field", "test3 test3"]
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
							{"type": "EQ", "args": ["string_field", "test1 test1"]},
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
							{"type": "EQ", "args": ["string_field", "test2 test2"]},
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
							{"type": "EQ", "args": ["string_field", "test1 test1"]},
							{"type": "EQ", "args": ["string_field", "test2 test2"]}
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
							{"type": "EQ", "args": ["string_field", "test1 test1"]},
							{"type": "EQ", "args": ["string_field", "test2 test2"]}
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
							{"type": "EQ", "args": ["string_field", "test1 test1"]},
							{"type": "EQ", "args": ["string_field", "test2 test2"]},
							{"type": "EQ", "args": ["string_field", "test3 test3"]}
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
									{"type": "EQ", "args": ["string_field", "test1 test1"]},
									{"type": "EQ", "args": ["string_field", "test2 test2"]}
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
							{"type": "IN", "args": ["string_field", "test1 test1", "test2 test2"]}
						]
					}`,
			count: 2,
		},
		{
			// Test GT on string field - lexicographic comparison
			//language=JSON
			filter: `{
						"type": "GT",
						"args": ["string_field", "test1 test1"]
					}`,
			count: 2,
		},
		{
			// Test GE on string field - lexicographic comparison
			//language=JSON
			filter: `{
						"type": "GE",
						"args": ["string_field", "test2 test2"]
					}`,
			count: 2,
		},
		{
			// Test LT on string field - lexicographic comparison
			//language=JSON
			filter: `{
						"type": "LT",
						"args": ["string_field", "test3 test3"]
					}`,
			count: 2,
		},
		{
			// Test LE on string field - lexicographic comparison
			//language=JSON
			filter: `{
						"type": "LE",
						"args": ["string_field", "test2 test2"]
					}`,
			count: 2,
		},
		{
			// Test GT on string field - no matches
			//language=JSON
			filter: `{
						"type": "GT",
						"args": ["string_field", "test3 test3"]
					}`,
			count: 0,
		},
		{
			// Test LT on string field - no matches
			//language=JSON
			filter: `{
						"type": "LT",
						"args": ["string_field", "test1 test1"]
					}`,
			count: 0,
		},
	}

	collection := SetupAtlasDB(t, ctx, atlasClient)
	InsertDocumentsOrFail(t, collection, ctx, documents)

	// Create search index with explicit field mappings
	// Index fields with multiple types (like ES multi-fields):
	// - type: "string" for text/wildcard search
	// - type: "token" for exact matching with equals operator
	searchIndexModel := mongo.SearchIndexModel{
		Definition: bson.D{
			// Define custom analyzer for case-insensitive keyword matching
			{"analyzers", bson.A{
				bson.D{
					{"name", "caseInsensitiveKeyword"},
					{"tokenizer", bson.D{
						{"type", "keyword"},
					}},
					{"tokenFilters", bson.A{
						bson.D{{"type", "lowercase"}},
					}},
				},
			}},
			{"mappings", bson.D{
				{"dynamic", false},
				{"fields", bson.D{
					// string_field: indexed as both string (for wildcard) and token (for equals)
					{"string_field", bson.A{
						// String type with standard analyzer (for TEXT queries) and keyword multi-analyzers (for LIKE/ILIKE)
						bson.D{
							{"type", "string"},
							{"analyzer", "lucene.standard"},
							{"multi", bson.D{
								{"keywordAnalyzer", bson.D{
									{"type", "string"},
									{"analyzer", "caseInsensitiveKeyword"},
								}},
								{"caseSensitiveKeywordAnalyzer", bson.D{
									{"type", "string"},
									{"analyzer", "lucene.keyword"},
								}},
							}},
						},
						// Token type for exact EQ/IN matching
						bson.D{{"type", "token"}},
					}},
					{"array_field", bson.D{
						// String supports (moreLikeThis, phrase, queryString, regex, span, text, wildcard)
						{"type", "string"},
					}},
					{"nullable_string_field", bson.A{
						// String type with standard analyzer (for TEXT queries) and keyword multi-analyzers (for LIKE/ILIKE)
						bson.D{
							{"type", "string"},
							{"analyzer", "lucene.standard"},
							{"multi", bson.D{
								{"keywordAnalyzer", bson.D{
									{"type", "string"},
									{"analyzer", "caseInsensitiveKeyword"},
								}},
								{"caseSensitiveKeywordAnalyzer", bson.D{
									{"type", "string"},
									{"analyzer", "lucene.keyword"},
								}},
							}},
						},
						// Token type for exact EQ/IN matching
						bson.D{{"type", "token"}},
					}},
					// text_field: indexed for both text search (with english analyzer) and exact matching
					{"text_field", bson.A{
						bson.D{
							// String supports (moreLikeThis, phrase, queryString, regex, span, text, wildcard)
							{"type", "string"},
							{"analyzer", "lucene.english"},
						},
						bson.D{
							// Token supports (equals, facet, in, range)
							{"type", "token"},
							{"normalizer", "none"}, // case-sensitive exact matching
						},
					}},
					// uuid_field: indexed as token for equals and in operations
					{"uuid_field", bson.D{
						{"type", "uuid"},
					}},
					// date_field: indexed as token for range and equality operations
					{"date_field", bson.D{
						{"type", "date"},
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

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.filter), func(t *testing.T) {
			/*
				Fixture Setup
			*/

			/*
			  Execute SUT
			*/

			// Create query builder
			// Configure multi-analyzers for fields that support LIKE/ILIKE
			var qb epsearchast.SemanticReducer[bson.D] = DefaultAtlasSearchQueryBuilder{
				FieldToMultiAnalyzers: map[string]*StringMultiAnalyzers{
					"string_field": {
						WildcardCaseInsensitive: "keywordAnalyzer",
						WildcardCaseSensitive:   "caseSensitiveKeywordAnalyzer",
					},
					"nullable_string_field": {
						WildcardCaseInsensitive: "keywordAnalyzer",
						WildcardCaseSensitive:   "caseSensitiveKeywordAnalyzer",
					},
				},
			}

			// Create Query Object
			ast, err := epsearchast.GetAst(tc.filter)
			if err != nil {
				t.Fatalf("Failed to get filter: %v", err)
			}

			query, err := epsearchast.SemanticReduceAst(ast, qb)

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
