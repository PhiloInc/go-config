package config

import (
	"fmt"
	"reflect"
	"strconv"
)

var intType = reflect.TypeOf(0)
var int8Type = reflect.TypeOf(int8(0))
var int16Type = reflect.TypeOf(int16(0))
var int32Type = reflect.TypeOf(int32(0))
var int64Type = reflect.TypeOf(int64(0))

type intSetter struct {
	val reflect.Value
	tag reflect.StructTag
}

func (is *intSetter) String() string {
	if is.val.Kind() == reflect.Invalid {
		return "0"
	}
	return strconv.FormatInt(is.val.Int(), 10)
}

func (is *intSetter) bitSize() int {
	switch is.val.Type() {
	case int8Type:
		return 8
	case int16Type:
		return 16
	case int32Type:
		return 32
	case int64Type:
		return 64
	}
	return 0
}

func (is *intSetter) Set(val string) error {
	ival, err := strconv.ParseInt(val, 0, is.bitSize())
	if err != nil {
		return &ConversionError{Value: val, ToType: is.val.Type()}
	}
	return is.set(ival)
}

func (is *intSetter) set(val int64) error {
	bitSize := is.bitSize()

	tag := is.tag.Get("le")
	if tag == "" {
		tag = is.tag.Get("max")
	}
	if tag != "" {
		n, err := strconv.ParseInt(tag, 0, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val <= n) {
			msg := fmt.Sprintf("%d is not less than or equal to %d", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	tag = is.tag.Get("ge")
	if tag == "" {
		tag = is.tag.Get("min")
	}
	if tag != "" {
		n, err := strconv.ParseInt(tag, 0, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val >= n) {
			msg := fmt.Sprintf("%d is not greater than or equal to %d", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	if tag = is.tag.Get("lt"); tag != "" {
		n, err := strconv.ParseInt(tag, 0, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val < n) {
			msg := fmt.Sprintf("%d is not less than %d", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	if tag = is.tag.Get("gt"); tag != "" {
		n, err := strconv.ParseInt(tag, 0, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val > n) {
			msg := fmt.Sprintf("%d is not greater than %d", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	is.val.SetInt(val)
	return nil
}

func (is *intSetter) SetInt(val int64) error {
	ival := reflect.New(is.val.Type()).Elem()
	ival.SetInt(val)
	if ival.Int() != val {
		return &ConversionError{Value: val, ToType: is.val.Type()}
	}
	return is.set(val)
}

func (is *intSetter) SetUint(val uint64) error {
	t := is.val.Type()
	ival := reflect.New(t).Elem()
	ival.Set(reflect.ValueOf(val).Convert(t))
	if i := ival.Int(); uint64(i) != val || i < 0 {
		return &ConversionError{Value: val, ToType: is.val.Type()}
	}
	return is.set(ival.Int())
}

func (is *intSetter) SetFloat(val float64) error {
	t := is.val.Type()
	ival := reflect.New(t).Elem()
	ival.Set(reflect.ValueOf(val).Convert(t))
	if float64(ival.Int()) != val {
		return &ConversionError{Value: val, ToType: is.val.Type()}
	}
	return is.set(ival.Int())
}

func (is *intSetter) SetBool(val bool) error {
	if val {
		return is.set(1)
	}
	return is.set(0)
}

func (is *intSetter) Get() interface{} {
	if is.val.Kind() == reflect.Invalid {
		return nil
	}
	return is.val.Interface()
}

type intSetterCreator struct {
	t reflect.Type
}

func (isc intSetterCreator) Type() reflect.Type {
	return isc.t
}

func (isc intSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	if val.Type() != isc.t {
		panic(fmt.Sprintf("value must be type %s", isc.t))
	}
	return &intSetter{val: val, tag: tag}
}

func init() {
	DefaultSetterRegistry.Add(intSetterCreator{t: intType})
	DefaultSetterRegistry.Add(intSetterCreator{t: int8Type})
	DefaultSetterRegistry.Add(intSetterCreator{t: int16Type})
	DefaultSetterRegistry.Add(intSetterCreator{t: int32Type})
	DefaultSetterRegistry.Add(intSetterCreator{t: int64Type})
}
