package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Errors represents a collection of one or more errors.
type Errors []error

/*
Append adds an error to the list of errors.

If err is of type *Errors, each element in it will be appended individually
instead of nesting the *Errors values.
*/
func (e *Errors) Append(err error) {
	switch err := err.(type) {
	case *Errors:
		*e = append(*e, (*err)...)
	case error:
		*e = append(*e, err)
	}
}

/*
Error returns a string by joining the Error() value of each element with
newlines.
*/
func (e *Errors) Error() string {
	ss := make([]string, len(*e))
	for i := range *e {
		ss[i] = (*e)[i].Error()
	}
	return strings.Join(ss, "\n")
}

/*
AsError returns e as an error type.

If *e is empty, returns error(nil). This is a convenience to avoid a common
pitfall:
	func foo() error {
		var *errs Errors
		return &errs
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
Values of an interface type only compare equal to nil if they have neither a
type nor a value. In this case, the return value from foo() will have a nil
value, but the type will be (*Errors), so err != nil is true.
*/
func (e *Errors) AsError() error {
	if len(*e) == 0 {
		return nil
	}
	return e
}

// UnknownTypeError represents an error when no Setter can be created for a type.
type UnknownTypeError struct {
	Type reflect.Type
	Path *Path
}

func (ute *UnknownTypeError) Error() string {
	return fmt.Sprintf("unknown type %s at path %s", ute.Type, ute.Path)
}

// ConversionError is returned for values which can not be convered by a Setter.
type ConversionError struct {
	Value  interface{}
	ToType reflect.Type
	Path   *Path
}

func (ce *ConversionError) Error() string {
	msg := fmt.Sprintf(
		"Cannot convert %v (type %T) to %s", ce.Value, ce.Value, ce.ToType,
	)
	if ce.Path != nil {
		msg += fmt.Sprintf(" at %s", ce.Path)
	}
	return msg
}

// ValidationError is returned by a Setter when values fail validation.
type ValidationError struct {
	Value   interface{}
	Message string
	Path    *Path
}

func (ve *ValidationError) Error() string {
	if ve.Path == nil {
		return fmt.Sprintf("Validating %v failed: %s", ve.Value, ve.Message)
	}
	return fmt.Sprintf("Validating %v failed at %s: %s", ve.Value, ve.Path, ve.Message)
}

// ErrHelp is returned when help is requested via the -h command line flag.
var ErrHelp = errors.New("Help requested")
