package config

import (
	"reflect"
	"strconv"
)

var boolType = reflect.TypeOf(false)

type boolSetter struct {
	val *bool
}

func (bs *boolSetter) IsBoolFlag() bool {
	return true
}

func (bs *boolSetter) String() string {
	if bs.val == nil {
		return strconv.FormatBool(false)
	}
	return strconv.FormatBool(*bs.val)
}

func (bs *boolSetter) Set(val string) error {
	bval, err := strconv.ParseBool(val)
	if err != nil {
		return &ConversionError{Value: val, ToType: boolType}
	}
	*bs.val = bval
	return nil
}

func (bs *boolSetter) SetInt(val int64) error {
	*bs.val = val != 0
	return nil
}

func (bs *boolSetter) SetUint(val uint64) error {
	*bs.val = val != 0
	return nil
}

func (bs *boolSetter) SetFloat(val float64) error {
	*bs.val = val != 0
	return nil
}

func (bs *boolSetter) SetBool(val bool) error {
	*bs.val = val
	return nil
}

func (bs *boolSetter) Get() interface{} {
	if bs.val == nil {
		return false
	}
	return *bs.val
}

type boolSetterCreator struct{}

func (bsc boolSetterCreator) Type() reflect.Type {
	return boolType
}

func (bsc boolSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	return &boolSetter{val: val.Addr().Interface().(*bool)}
}

func init() {
	DefaultSetterRegistry.Add(boolSetterCreator{})
}
