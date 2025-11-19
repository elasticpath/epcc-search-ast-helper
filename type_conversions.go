package epsearchast

import (
	"fmt"
	"strconv"
)

type FieldType int

const (
	String FieldType = iota
	Int64
	Boolean
	Float64
)

func (f FieldType) String() string {
	switch f {
	case String:
		return "string"
	case Int64:
		return "int64"
	case Boolean:
		return "bool"
	case Float64:
		return "float64"
	default:
		return "unknown"
	}
}

func Convert(t FieldType, v string) (interface{}, error) {

	err := ValidateValue(t, v)
	if err != nil {
		return nil, err
	}

	var newV interface{}
	switch t {
	case String:
		return v, nil
	case Int64:
		newV, _ = strconv.ParseInt(v, 10, 64)
	case Boolean:
		newV, _ = strconv.ParseBool(v)
	case Float64:
		newV, _ = strconv.ParseFloat(v, 64)
	}

	return newV, nil
}

func ConvertAll(t FieldType, vAll ...string) ([]interface{}, error) {

	newVAll := make([]interface{}, 0, len(vAll))

	for idx, v := range vAll {

		newV, err := Convert(t, v)
		if err != nil {
			return nil, fmt.Errorf("error converting value at index %v: %w", idx, err)
		}

		newVAll = append(newVAll, newV)
	}

	return newVAll, nil
}

func ValidateValue(t FieldType, v string) error {

	switch t {
	case String:
		return nil
	case Int64:
		_, e := strconv.ParseInt(v, 10, 64)
		if e != nil {
			return fmt.Errorf("invalid value for int64: `%v`", v)
		}
		return nil
	case Float64:
		_, e := strconv.ParseFloat(v, 64)
		if e != nil {
			return fmt.Errorf("invalid value for float64: `%v`", v)
		}
		return nil

	case Boolean:
		_, e := strconv.ParseBool(v)
		if e != nil {
			return fmt.Errorf("invalid value for boolean: `%v`", v)
		}
		return nil
	default:
		return fmt.Errorf("Unsupported field type %v:%v", t, v)
	}
}

func ValidateAllValues(t FieldType, v ...string) error {
	for idx, value := range v {
		err := ValidateValue(t, value)
		if err != nil {
			return fmt.Errorf("could not validate position %d, the value [%s] could not be converted: %w", idx, value, err)
		}
	}
	return nil
}
