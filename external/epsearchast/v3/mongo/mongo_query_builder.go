package epsearchast_v3_mongo

import (
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"regexp"
	"strings"
)
import "go.mongodb.org/mongo-driver/bson"

type DefaultMongoQueryBuilder struct{}

var _ epsearchast_v3.SemanticReducer[bson.D] = (*DefaultMongoQueryBuilder)(nil)

func (d DefaultMongoQueryBuilder) PostVisitAnd(rs []*bson.D) (*bson.D, error) {
	// https://www.mongodb.com/docs/manual/reference/operator/query/and/
	return &bson.D{
		{"$and",
			rs,
		},
	}, nil
}

func (d DefaultMongoQueryBuilder) VisitIn(args ...string) (*bson.D, error) {
	// https://www.mongodb.com/docs/manual/reference/operator/query/in/
	return &bson.D{{args[0], bson.D{{"$in", args[1:]}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitEq(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/manual/reference/operator/query/eq/#std-label-eq-usage-examples
	// This is equivalent to { key: value } but makes for easier tests.
	return &bson.D{{first, bson.D{{"$eq", second}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitLe(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/manual/reference/operator/query/lte/
	return &bson.D{{first, bson.D{{"$lte", second}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitLt(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/manual/reference/operator/query/lt/
	return &bson.D{{first, bson.D{{"$lt", second}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitGe(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/manual/reference/operator/query/gte/
	return &bson.D{{first, bson.D{{"$gte", second}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitGt(first, second string) (*bson.D, error) {
	// https://www.mongodb.com/docs/manual/reference/operator/query/gt/
	return &bson.D{{first, bson.D{{"$gt", second}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitLike(first, second string) (*bson.D, error) {
	return &bson.D{{first, bson.D{{"$regex", d.ProcessLikeWildcards(second)}}}}, nil
}

func (d DefaultMongoQueryBuilder) ProcessLikeWildcards(valString string) string {
	if valString == "*" {
		return "^.*$"
	}

	var startsWithStar = strings.HasPrefix(valString, "*")
	var endsWithStar = strings.HasSuffix(valString, "*")
	if startsWithStar {
		valString = valString[1:]
	}
	if endsWithStar {
		valString = valString[:len(valString)-1]
	}

	valString = regexp.QuoteMeta(valString)

	if startsWithStar {
		valString = ".*" + valString
	}
	if endsWithStar {
		valString += ".*"
	}
	return "^" + valString + "$"
}
