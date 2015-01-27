package main

import "errors"

type ErrValidation string

func (ev ErrValidation) Error() string {
	return string(ev)
}

type appError struct {
	Err  error
	Code int
}

func (ae appError) Error() string {
	return ae.Err.Error()
}

var ErrInput = errors.New("error: invalid input type")
var ErrInvalidLength = errors.New("error: invalid content length")
