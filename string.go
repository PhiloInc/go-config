package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

var stringType = reflect.TypeOf("")

type stringSetter struct {
	val *string
	tag reflect.StructTag
}

func (ss *stringSetter) String() string {
	if ss.val == nil {
		return ""
	}
	return *ss.val
}

func (ss *stringSetter) Set(val string) error {
	if tag := ss.tag.Get("regexp"); tag != "" {
		r, err := regexp.Compile(tag)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !r.MatchString(val) {
			msg := fmt.Sprintf("'%s' did not match regular expression '%s'", val, r)
			return &ValidationError{Value: val, Message: msg}
		}
	}
	*ss.val = val
	return nil
}

func (ss *stringSetter) SetInt(val int64) error {
	return ss.Set(strconv.FormatInt(val, 10))
}

func (ss *stringSetter) SetUint(val uint64) error {
	return ss.Set(strconv.FormatUint(val, 10))
}

func (ss *stringSetter) SetFloat(val float64) error {
	return ss.Set(strconv.FormatFloat(val, 'f', -1, 64))
}

func (ss *stringSetter) SetBool(val bool) error {
	return ss.Set(strconv.FormatBool(val))
}

func (ss *stringSetter) Get() interface{} {
	if ss.val == nil {
		return ""
	}
	return *ss.val
}

type stringSetterCreator struct{}

func (ssc stringSetterCreator) Type() reflect.Type {
	return stringType
}

func (ssc stringSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	return &stringSetter{val: val.Addr().Interface().(*string), tag: tag}
}

func init() {
	DefaultSetterRegistry.Add(stringSetterCreator{})
}
