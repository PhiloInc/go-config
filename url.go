package config

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
)

var urlType = reflect.TypeOf(url.URL{})

type urlSetter struct {
	val *url.URL
	tag reflect.StructTag
}

func (us *urlSetter) String() string {
	if us.val == nil {
		return (&url.URL{}).String()
	}
	return us.val.String()
}

func (us *urlSetter) Set(val string) error {
	uval, err := url.Parse(val)
	if err != nil {
		return &ConversionError{Value: val, ToType: urlType}
	}

	pairs := [][2]string{
		[2]string{"scheme", uval.Scheme},
		[2]string{"host", uval.Host},
		[2]string{"path", uval.Path},
	}
	for i := range pairs {
		key, s := pairs[i][0], pairs[i][1]
		if tag := us.tag.Get(key); tag != "" {
			r, err := regexp.Compile(tag)
			if err != nil {
				return &ValidationError{Value: uval, Message: err.Error()}
			}
			if !r.MatchString(s) {
				msg := fmt.Sprintf(
					"'%s' did not match regular expression '%s'", s, r,
				)
				return &ValidationError{Value: val, Message: msg}
			}
		}
	}

	*us.val = *uval
	return nil
}

func (*urlSetter) SetInt(val int64) error {
	return &ConversionError{Value: val, ToType: urlType}
}

func (*urlSetter) SetUint(val uint64) error {
	return &ConversionError{Value: val, ToType: urlType}
}

func (*urlSetter) SetFloat(val float64) error {
	return &ConversionError{Value: val, ToType: urlType}
}

func (*urlSetter) SetBool(val bool) error {
	return &ConversionError{Value: val, ToType: urlType}
}

func (us *urlSetter) Get() interface{} {
	if us.val == nil {
		return new(*url.URL)
	}
	return *us.val
}

type urlSetterCreator struct{}

func (urlSetterCreator) Type() reflect.Type {
	return urlType
}

func (urlSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	return &urlSetter{val: val.Addr().Interface().(*url.URL), tag: tag}
}

func init() {
	DefaultSetterRegistry.Add(urlSetterCreator{})
}
