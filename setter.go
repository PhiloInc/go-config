package config

import (
	"reflect"
)

var setterType = reflect.TypeOf((*Setter)(nil)).Elem()

/*
Setter is an interface which is used to parse, validate and set values of a
specific type.

Note that SetInt, SetUint, SetFloat and SetBool are not used for the default
loaders, but may be used by custom Loader implementations or in future
features.
*/
type Setter interface {
	// String should return a descriptive string for the current value
	// associated with the Setter. Note that this method may be called on a
	// zero value of the implementation type and should not panic in that case.
	String() string
	// Set is called to parse, validate and set the associated value from a
	// string. The implementation must return an error (such as
	// *ConversionError) if the value cannot be parsed from the given string, or
	// if string is not a valid type to represent the parameter.
	Set(string) error
	// SetInt is called to interpret, validate and set the associated value from
	// an int64. The implementation must return an error (such as
	// *ConversionError) if the value cannot be interpreted from the given
	// int64, or if int64 is not a valid type to represent the parameter.
	SetInt(int64) error
	// SetUint is called to interpret, validate and set the associated value
	// from a uint64. The implementation must return an error (such as
	// *ConversionError) if the value cannot be interpreted from the given
	// uint64, or if uint64 is not a valid type to represent the parameter.
	SetUint(uint64) error
	// SetFloat is called to interpret, validate and set the associated value
	// from a float64. The implementation must return an error (such as
	// *ConversionError) if the value cannot be interpreted from the given
	// float64, or if float64 is not a valid type to represent the parameter.
	SetFloat(float64) error
	// SetBool is called to interpret, validate and set the associated value
	// from a bool. The implementation must return an error (such as
	// *ConversionError) if the value cannot be interpreted from the given
	// bool, or if bool is not a valid type to represent the parameter.
	SetBool(bool) error
	// Get is called to get the current value set on the Setter. The type of
	// the returned value should be the same as the value being set, even if
	// the value is currently unset or zero.
	Get() interface{}
}

// SetterCreator is an interface type for creating a Setter for a given value.
type SetterCreator interface {
	// Type must return the reflect.Type that the Setter returned from Setter
	// will support.
	Type() reflect.Type
	// Setter must return an implementation of Setter to set a value of the
	// type returned by Type(). The value passed to Setter will always be of
	// the type returned by Type(), and the Setter must be able to accept
	// and set that type.
	Setter(reflect.Value, reflect.StructTag) Setter
}

/*
SetterRegistry is a map of reflect.Type to SetterCreator.

It provides methods that can be used to find or create a Setter for a given
value.
*/
type SetterRegistry map[reflect.Type]SetterCreator

/*
Add adds a SetterCreator to the registry.

If there is already a SetterCreator registered for the same type returned by
sc.Type(), it will be replaced.
*/
func (sr *SetterRegistry) Add(sc SetterCreator) {
	if *sr == nil {
		*sr = make(map[reflect.Type]SetterCreator)
	}
	(*sr)[sc.Type()] = sc
}

/*
GetSetterCreator returns the registered SetterCreator for the given type, or nil
if none is registered.
*/
func (sr *SetterRegistry) GetSetterCreator(t reflect.Type) SetterCreator {
	return (*sr)[t]
}

/*
GetSetter returns a Setter that wraps/sets the given val.

If the value and type represented by val implements the Setter type,
val.Interface() is simply returned.

If val.Type() has an existing entry in the registry, the registered
SetterCreator will be used to create the Setter. Otherwise, if val.Type() is a
slice or pointer type, it will be dereferenced until either a type is found in
the registry, or a non-element type is found. If a Setter is returned for a
pointer or slice type, it may be created using wrappers to handle the
indirections, so may not be a value returned by one of the registered
SetterCreators.

If no SetterCreator can be found or created, GetSetter returns nil.
*/
func (sr *SetterRegistry) GetSetter(val reflect.Value, tag reflect.StructTag) Setter {
	// If a value implements the Setter interface it overrides everything else.
	if val.Type().Implements(setterType) {
		return val.Interface().(Setter)
	}

	if !val.CanSet() {
		panic("val must be settable")
	}

	var wrapperFns []func(SetterCreator) SetterCreator
	var sc SetterCreator
	t := val.Type()
	for sc = (*sr).GetSetterCreator(t); sc == nil; sc = (*sr).GetSetterCreator(t) {
		switch t.Kind() {
		case reflect.Slice:
			wrapperFns = append(wrapperFns, newSliceSetterCreator)
		case reflect.Ptr:
			wrapperFns = append(wrapperFns, newPtrSetterCreator)
		default:
			return nil
		}
		t = t.Elem()
	}

	for i := len(wrapperFns) - 1; i >= 0; i-- {
		sc = wrapperFns[i](sc)
	}

	return sc.Setter(val, tag)
}

/*
Copy creates a copy of the registry.

It is useful for modifying the existing registry without affecting other
references to it.
*/
func (sr *SetterRegistry) Copy() *SetterRegistry {
	n := make(SetterRegistry, len(*sr))
	for k, v := range *sr {
		n[k] = v
	}
	return &n
}

/*
DefaultSetterRegistry provides a default registry.

It supports SetterCreator instances for the built-in types, and is used by
Config if no other registry is specified.
*/
var DefaultSetterRegistry = SetterRegistry{}
