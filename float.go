package config

import (
	"fmt"
	"reflect"
	"strconv"
)

var float32Type = reflect.TypeOf(float32(0))
var float64Type = reflect.TypeOf(float64(0))

type floatSetter struct {
	val reflect.Value
	tag reflect.StructTag
}

func (fs *floatSetter) String() string {
	if fs.val.Kind() == reflect.Invalid {
		return "0"
	}
	return strconv.FormatFloat(fs.val.Float(), 'f', -1, fs.bitSize())
}

func (fs *floatSetter) bitSize() int {
	if fs.val.Type() == float32Type {
		return 32
	}
	return 64
}

func (fs *floatSetter) Set(val string) error {
	fval, err := strconv.ParseFloat(val, fs.bitSize())
	if err != nil {
		return &ConversionError{Value: val, ToType: fs.val.Type()}
	}
	return fs.set(fval)
}

func (fs *floatSetter) set(val float64) error {
	bitSize := fs.bitSize()

	tag := fs.tag.Get("le")
	if tag == "" {
		tag = fs.tag.Get("max")
	}
	if tag != "" {
		n, err := strconv.ParseFloat(tag, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val <= n) {
			msg := fmt.Sprintf("%f is not less than or equal to %f", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	tag = fs.tag.Get("ge")
	if tag == "" {
		tag = fs.tag.Get("min")
	}
	if tag != "" {
		n, err := strconv.ParseFloat(tag, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val >= n) {
			msg := fmt.Sprintf("%f is not greater than or equal to %f", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	if tag = fs.tag.Get("lt"); tag != "" {
		n, err := strconv.ParseFloat(tag, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val < n) {
			msg := fmt.Sprintf("%f is not less than %f", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	if tag = fs.tag.Get("gt"); tag != "" {
		n, err := strconv.ParseFloat(tag, bitSize)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val > n) {
			msg := fmt.Sprintf("%f is not greater than %f", val, n)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	fs.val.SetFloat(val)
	return nil
}

func (fs *floatSetter) SetInt(val int64) error {
	t := fs.val.Type()
	fval := reflect.New(t).Elem()
	fval.Set(reflect.ValueOf(val).Convert(t))
	if int64(fval.Float()) != val {
		return &ConversionError{Value: val, ToType: fs.val.Type()}
	}
	return fs.set(fval.Float())
}

func (fs *floatSetter) SetUint(val uint64) error {
	t := fs.val.Type()
	fval := reflect.New(t).Elem()
	fval.Set(reflect.ValueOf(val).Convert(t))
	if uint64(fval.Float()) != val {
		return &ConversionError{Value: val, ToType: fs.val.Type()}
	}
	return fs.set(fval.Float())
}

func (fs *floatSetter) SetFloat(val float64) error {
	fval := reflect.New(fs.val.Type()).Elem()
	fval.SetFloat(val)
	if fval.Float() != val {
		return &ConversionError{Value: val, ToType: fs.val.Type()}
	}
	return fs.set(val)
}

func (fs *floatSetter) SetBool(val bool) error {
	if val {
		return fs.set(1)
	}
	return fs.set(0)
}

func (fs *floatSetter) Get() interface{} {
	if fs.val.Kind() == reflect.Invalid {
		return nil
	}
	return fs.val.Interface()
}

type floatSetterCreator struct {
	t reflect.Type
}

func (fsc floatSetterCreator) Type() reflect.Type {
	return fsc.t
}

func (fsc floatSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	if val.Type() != fsc.t {
		panic(fmt.Sprintf("value must be type %s", fsc.t))
	}
	return &floatSetter{val: val, tag: tag}
}

func init() {
	DefaultSetterRegistry.Add(floatSetterCreator{t: float32Type})
	DefaultSetterRegistry.Add(floatSetterCreator{t: float64Type})
}
