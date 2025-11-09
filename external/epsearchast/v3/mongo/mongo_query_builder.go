package epsearchast_v3_mongo

import (
	"fmt"
	epsearchast_v3 "github.com/elasticpath/epcc-search-ast-helper/external/epsearchast/v3"
	"regexp"
	"strings"
)
import "go.mongodb.org/mongo-driver/bson"

type DefaultMongoQueryBuilder struct {
	FieldTypes map[string]epsearchast_v3.FieldType
}

var _ epsearchast_v3.SemanticReducer[bson.D] = (*DefaultMongoQueryBuilder)(nil)

func (d DefaultMongoQueryBuilder) PostVisitAnd(rs []*bson.D) (*bson.D, error) {
	// https://www.mongodb.com/docs/manual/reference/operator/query/and/
	return &bson.D{
		{"$and",
			rs,
		},
	}, nil
}

func (d DefaultMongoQueryBuilder) PostVisitOr(rs []*bson.D) (*bson.D, error) {
	// https://www.mongodb.com/docs/manual/reference/operator/query/or/
	return &bson.D{
		{"$or",
			rs,
		},
	}, nil
}

func (d DefaultMongoQueryBuilder) VisitIn(args ...string) (*bson.D, error) {

	if err := d.ValidateValues(args[0], args[1:]...); err != nil {
		return nil, err
	}
	// https://www.mongodb.com/docs/manual/reference/operator/query/in/
	return &bson.D{{args[0], bson.D{{"$in", d.ConvertValues(args[0], args[1:]...)}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitEq(first, second string) (*bson.D, error) {
	if err := d.ValidateValue(first, second); err != nil {
		return nil, err
	}

	// https://www.mongodb.com/docs/manual/reference/operator/query/eq/#std-label-eq-usage-examples
	// This is equivalent to { key: value } but makes for easier tests.
	return &bson.D{{first, bson.D{{"$eq", d.ConvertValue(first, second)}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitLe(first, second string) (*bson.D, error) {
	if err := d.ValidateValue(first, second); err != nil {
		return nil, err
	}

	// https://www.mongodb.com/docs/manual/reference/operator/query/lte/
	return &bson.D{{first, bson.D{{"$lte", d.ConvertValue(first, second)}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitLt(first, second string) (*bson.D, error) {
	if err := d.ValidateValue(first, second); err != nil {
		return nil, err
	}

	// https://www.mongodb.com/docs/manual/reference/operator/query/lt/
	return &bson.D{{first, bson.D{{"$lt", d.ConvertValue(first, second)}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitGe(first, second string) (*bson.D, error) {
	if err := d.ValidateValue(first, second); err != nil {
		return nil, err
	}

	// https://www.mongodb.com/docs/manual/reference/operator/query/gte/
	return &bson.D{{first, bson.D{{"$gte", d.ConvertValue(first, second)}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitGt(first, second string) (*bson.D, error) {
	if err := d.ValidateValue(first, second); err != nil {
		return nil, err
	}

	// https://www.mongodb.com/docs/manual/reference/operator/query/gt/
	return &bson.D{{first, bson.D{{"$gt", d.ConvertValue(first, second)}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitLike(first, second string) (*bson.D, error) {
	if v, ok := d.FieldTypes[first]; ok {
		if v != epsearchast_v3.String {
			return nil, fmt.Errorf("like() operator is only supported for string fields, and [%s] is not a string", first)
		}
	}

	return &bson.D{{first, bson.D{{"$regex", d.ProcessLikeWildcards(second)}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitILike(first, second string) (*bson.D, error) {
	if v, ok := d.FieldTypes[first]; ok {
		if v != epsearchast_v3.String {
			return nil, fmt.Errorf("ilike() operator is only supported for string fields, and [%s] is not a string", first)
		}
	}

	return &bson.D{
		{first,
			bson.D{
				{"$regex", d.ProcessLikeWildcards(second)},
				{"$options", "i"},
			},
		},
	}, nil
}

func (d DefaultMongoQueryBuilder) VisitContains(first, second string) (*bson.D, error) {
	if err := d.ValidateValue(first, second); err != nil {
		return nil, err
	}

	// https://www.mongodb.com/docs/manual/reference/operator/query/elemMatch/
	// This is equivalent to { key: value } but makes for easier tests.
	return &bson.D{{first,
		bson.D{
			{"$elemMatch", bson.D{
				{"$eq", d.ConvertValue(first, second)},
			},
			},
		},
	}}, nil

}

func (d DefaultMongoQueryBuilder) VisitContainsAny(args ...string) (*bson.D, error) {

	if err := d.ValidateValues(args[0], args[1:]...); err != nil {
		return nil, err
	}
	// https://www.mongodb.com/docs/manual/reference/operator/query/in/
	// Matches arrays that contain at least one element that matches any of the provided values
	return &bson.D{{args[0], bson.D{{"$in", d.ConvertValues(args[0], args[1:]...)}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitContainsAll(args ...string) (*bson.D, error) {

	if err := d.ValidateValues(args[0], args[1:]...); err != nil {
		return nil, err
	}
	// https://www.mongodb.com/docs/manual/reference/operator/query/all/
	// Matches arrays that contain all the specified elements
	return &bson.D{{args[0], bson.D{{"$all", d.ConvertValues(args[0], args[1:]...)}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitText(first, second string) (*bson.D, error) {
	if v, ok := d.FieldTypes[first]; ok {
		if v != epsearchast_v3.String {
			return nil, fmt.Errorf("text() operator is only supported for string fields, and [%s] is not a string", first)
		}
	}

	// https://www.mongodb.com/docs/v7.0/reference/operator/query/text/#std-label-text-operator-phrases
	return &bson.D{{"$text", bson.D{{"$search", second}}}}, nil
}

func (d DefaultMongoQueryBuilder) VisitIsNull(first string) (*bson.D, error) {
	// https://www.mongodb.com/docs/manual/tutorial/query-for-null-fields/#equality-filter
	// This will match fields that either contain the item field whose value is nil or those that do not contain the field
	// Customize this method if you need different nil handling (i.e., explicit nil)
	return &bson.D{{first, bson.D{{"$eq", nil}}}}, nil
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

func (d DefaultMongoQueryBuilder) ValidateValue(fieldName string, v string) error {

	if fieldType, ok := d.FieldTypes[fieldName]; ok {
		return epsearchast_v3.ValidateValue(fieldType, v)
	}

	return nil
}

func (d DefaultMongoQueryBuilder) ValidateValues(fieldName string, v ...string) error {
	if fieldType, ok := d.FieldTypes[fieldName]; ok {
		return epsearchast_v3.ValidateAllValues(fieldType, v...)
	} else {
		return nil
	}
}

func (d DefaultMongoQueryBuilder) ConvertValue(fieldName string, v string) interface{} {

	if fieldType, ok := d.FieldTypes[fieldName]; ok {
		v, _ := epsearchast_v3.Convert(fieldType, v)
		return v
	}

	return v
}

func (d DefaultMongoQueryBuilder) ConvertValues(fieldName string, v ...string) []interface{} {

	if fieldType, ok := d.FieldTypes[fieldName]; ok {
		v, _ := epsearchast_v3.ConvertAll(fieldType, v...)
		return v
	} else {
		// We need to do the conversion to string, because we got a []string in, and need to
		// return a []interface{}
		v, _ := epsearchast_v3.ConvertAll(epsearchast_v3.String, v...)
		return v
	}
}
