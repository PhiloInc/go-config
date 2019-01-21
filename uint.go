package config

import (
	"fmt"
	"reflect"
	"strconv"
)

var uintType = reflect.TypeOf(uint(0))
var uint8Type = reflect.TypeOf(uint8(0))
var uint16Type = reflect.TypeOf(uint16(0))
var uint32Type = reflect.TypeOf(uint32(0))
var uint64Type = reflect.TypeOf(uint64(0))

type uintSetter struct {
	val reflect.Value
	tag reflect.StructTag
}

func (us *uintSetter) String() string {
	if us.val.Kind() == reflect.Invalid {
		return "0"
	}
	return strconv.FormatUint(us.val.Uint(), 10)
}

func (us *uintSetter) bitSize() int {
	switch us.val.Type() {
	case uint8Type:
		return 8
	case uint16Type:
		return 16
	case uint32Type:
		return 32
	case uint64Type:
		return 64
	}
	return 0
}

func (us *uintSetter) Set(val string) error {
	uval, err := strconv.ParseUint(val, 0, us.bitSize())
	if err != nil {
		return &ConversionError{Value: val, ToType: us.val.Type()}
	}
	return us.set(uval)
}

func (us *uintSetter) set(val uint64) error {
	bitSize := us.bitSize()

	tag := us.tag.Get("le")
	if tag == "" {
		tag = us.tag.Get("max")
	}
	if tag != "" {
		n, err := strconv.ParseUint(tag, 0, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val <= n) {
			msg := fmt.Sprintf("%d is not less than or equal to %d", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	tag = us.tag.Get("ge")
	if tag == "" {
		tag = us.tag.Get("min")
	}
	if tag != "" {
		n, err := strconv.ParseUint(tag, 0, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val >= n) {
			msg := fmt.Sprintf("%d is not greater than or equal to %d", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	if tag = us.tag.Get("lt"); tag != "" {
		n, err := strconv.ParseUint(tag, 0, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val < n) {
			msg := fmt.Sprintf("%d is not less than %d", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	if tag = us.tag.Get("gt"); tag != "" {
		n, err := strconv.ParseUint(tag, 0, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val > n) {
			msg := fmt.Sprintf("%d is not greater than %d", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	us.val.SetUint(val)
	return nil
}

func (us *uintSetter) SetInt(val int64) error {
	t := us.val.Type()
	uval := reflect.New(t).Elem()
	uval.Set(reflect.ValueOf(val).Convert(t))
	if int64(uval.Uint()) != val {
		return &ConversionError{Value: val, ToType: us.val.Type()}
	}
	return us.set(uval.Uint())
}

func (us *uintSetter) SetUint(val uint64) error {
	uval := reflect.New(us.val.Type()).Elem()
	uval.SetUint(val)
	if uval.Uint() != val {
		return &ConversionError{Value: val, ToType: us.val.Type()}
	}
	return us.set(val)
}

func (us *uintSetter) SetFloat(val float64) error {
	t := us.val.Type()
	uval := reflect.New(t).Elem()
	uval.Set(reflect.ValueOf(val).Convert(t))
	if float64(uval.Uint()) != val {
		return &ConversionError{Value: val, ToType: us.val.Type()}
	}
	return us.set(uval.Uint())
}

func (us *uintSetter) SetBool(val bool) error {
	if val {
		return us.set(1)
	}
	return us.set(0)
}

func (us *uintSetter) Get() interface{} {
	if us.val.Kind() == reflect.Invalid {
		return nil
	}
	return us.val.Interface()
}

type uintSetterCreator struct {
	t reflect.Type
}

func (usc uintSetterCreator) Type() reflect.Type {
	return usc.t
}

func (usc uintSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	if val.Type() != usc.t {
		panic(fmt.Sprintf("value must be type %s", usc.t))
	}
	return &uintSetter{val: val, tag: tag}
}

func init() {
	DefaultSetterRegistry.Add(uintSetterCreator{t: uintType})
	DefaultSetterRegistry.Add(uintSetterCreator{t: uint8Type})
	DefaultSetterRegistry.Add(uintSetterCreator{t: uint16Type})
	DefaultSetterRegistry.Add(uintSetterCreator{t: uint32Type})
	DefaultSetterRegistry.Add(uintSetterCreator{t: uint64Type})
}
