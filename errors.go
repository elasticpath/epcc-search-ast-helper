package epsearchast

import "fmt"

type ParsingErr struct {
	err error
}

func (pe ParsingErr) Error() string {
	return pe.err.Error()
}

func NewParsingErr(err error) ParsingErr {
	return ParsingErr{
		err: fmt.Errorf("could not parse filter: %w", err),
	}
}

type ValidationErr struct {
	err error
}

func (ve ValidationErr) Error() string {
	return ve.err.Error()
}

func NewValidationErr(err error) ValidationErr {
	return ValidationErr{
		err: fmt.Errorf("error validating filter: %w", err),
	}
}
