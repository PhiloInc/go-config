package config

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})
var durationType = reflect.TypeOf(time.Duration(0))

func parseTime(value string) (time.Time, error) {
	var layouts = []string{
		time.RFC3339,
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04Z07:00",
		"2006-01-02 15:04Z07:00",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02T15",
		"2006-01-02 15",
		"2006-01-02",
		"2006-01",
	}
	var t time.Time
	var err error
	for i := range layouts {
		if t, err = time.Parse(layouts[i], value); err == nil {
			break
		}
	}
	if err != nil {
		return t, &ConversionError{Value: value, ToType: timeType}
	}
	return t, nil
}

type timeSetter struct {
	val *time.Time
	tag reflect.StructTag
}

func (ts *timeSetter) String() string {
	if ts.val == nil {
		return time.Time{}.Format(time.RFC3339)
	}
	return ts.val.Format(time.RFC3339)
}

func (ts *timeSetter) set(tval time.Time) error {
	var tag string
	var err error

	if tag = ts.tag.Get("le"); tag == "" {
		tag = ts.tag.Get("max")
	}
	if tag != "" {
		var t time.Time
		if tag == "now" {
			t = time.Now()
		} else if t, err = parseTime(tag); err != nil {
			msg := fmt.Sprintf("invalid time %s", tag)
			return &ValidationError{Value: tval, Message: msg}
		}
		if tval.After(t) {
			msg := fmt.Sprintf("%s is after %s", tval.Format(time.RFC3339), t.Format(time.RFC3339))
			return &ValidationError{Value: tval, Message: msg}
		}
	}

	if tag = ts.tag.Get("ge"); tag == "" {
		tag = ts.tag.Get("min")
	}
	if tag != "" {
		var t time.Time
		if tag == "now" {
			t = time.Now()
		} else if t, err = parseTime(tag); err != nil {
			msg := fmt.Sprintf("invalid time %s", tag)
			return &ValidationError{Value: tval, Message: msg}
		}
		if tval.Before(t) {
			msg := fmt.Sprintf("%s is before %s", tval.Format(time.RFC3339), t.Format(time.RFC3339))
			return &ValidationError{Value: tval, Message: msg}
		}
	}

	if tag = ts.tag.Get("lt"); tag != "" {
		var t time.Time
		if tag == "now" {
			t = time.Now()
		} else if t, err = parseTime(tag); err != nil {
			msg := fmt.Sprintf("invalid time %s", tag)
			return &ValidationError{Value: tval, Message: msg}
		}
		if !tval.Before(t) {
			msg := fmt.Sprintf("%s is not before %s", tval.Format(time.RFC3339), t.Format(time.RFC3339))
			return &ValidationError{Value: tval, Message: msg}
		}
	}

	if tag = ts.tag.Get("gt"); tag != "" {
		var t time.Time
		if tag == "now" {
			t = time.Now()
		} else if t, err = parseTime(tag); err != nil {
			msg := fmt.Sprintf("invalid time %s", tag)
			return &ValidationError{Value: tval, Message: msg}
		}
		if !tval.After(t) {
			msg := fmt.Sprintf("%s is not after %s", tval.Format(time.RFC3339), t.Format(time.RFC3339))
			return &ValidationError{Value: tval, Message: msg}
		}
	}

	*ts.val = tval
	return nil
}

func (ts *timeSetter) Set(val string) error {
	tval, err := parseTime(val)
	if err != nil {
		return err
	}
	return ts.set(tval)
}

func (ts *timeSetter) SetInt(val int64) error {
	tval := time.Unix(val, 0)
	return ts.set(tval)
}

func (ts *timeSetter) SetUint(val uint64) error {
	tval := time.Unix(int64(val), 0)
	if u := tval.Unix(); uint64(u) != val || u < 0 {
		return &ConversionError{Value: val, ToType: timeType}
	}
	return ts.set(tval)
}

func (ts *timeSetter) SetFloat(val float64) error {
	d, err := time.ParseDuration(strconv.FormatFloat(val, 'f', -1, 64) + "s")
	if err != nil {
		return &ConversionError{Value: val, ToType: durationType}
	}
	s, ns := d.Seconds(), d%time.Second
	tval := time.Unix(int64(s), int64(ns))
	return ts.set(tval)
}

func (*timeSetter) SetBool(val bool) error {
	return &ConversionError{Value: val, ToType: timeType}
}

func (ts *timeSetter) Get() interface{} {
	if ts.val == nil {
		return time.Time{}
	}
	return *ts.val
}

type timeSetterCreator struct{}

func (timeSetterCreator) Type() reflect.Type {
	return timeType
}

func (timeSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	return &timeSetter{val: val.Addr().Interface().(*time.Time), tag: tag}
}

type durationSetter struct {
	val *time.Duration
	tag reflect.StructTag
}

func (ds *durationSetter) String() string {
	if ds.val == nil {
		return time.Duration(0).String()
	}
	return ds.val.String()
}

func (ds *durationSetter) set(val time.Duration) error {
	tag := ds.tag.Get("le")
	if tag == "" {
		tag = ds.tag.Get("max")
	}
	if tag != "" {
		d, err := time.ParseDuration(tag)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val <= d) {
			msg := fmt.Sprintf("%s is not less than or equal to %s", val, d)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	tag = ds.tag.Get("ge")
	if tag == "" {
		tag = ds.tag.Get("min")
	}
	if tag != "" {
		d, err := time.ParseDuration(tag)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val >= d) {
			msg := fmt.Sprintf("%s is not greater than or equal to %s", val, d)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	if tag = ds.tag.Get("lt"); tag != "" {
		d, err := time.ParseDuration(tag)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val < d) {
			msg := fmt.Sprintf("%s is not less than %s", val, d)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	if tag = ds.tag.Get("gt"); tag != "" {
		d, err := time.ParseDuration(tag)
		if err != nil {
			return &ValidationError{Value: val, Message: err.Error()}
		}
		if !(val > d) {
			msg := fmt.Sprintf("%s is not greater than %s", val, d)
			return &ValidationError{Value: val, Message: msg}
		}
	}

	*ds.val = val
	return nil
}

func (ds *durationSetter) Set(val string) error {
	dval, err := time.ParseDuration(val)
	if err != nil {
		return &ConversionError{Value: val, ToType: durationType}
	}
	return ds.set(dval)
}

func (ds *durationSetter) SetInt(val int64) error {
	dval := time.Duration(val) * time.Second
	if int64(dval/time.Second) != val {
		return &ConversionError{Value: val, ToType: durationType}
	}
	return ds.set(dval)
}

func (ds *durationSetter) SetUint(val uint64) error {
	dval := time.Duration(val) * time.Second
	if uint64(dval/time.Second) != val || dval < 0 {
		return &ConversionError{Value: val, ToType: durationType}
	}
	return ds.set(dval)
}

func (ds *durationSetter) SetFloat(val float64) error {
	dval, err := time.ParseDuration(strconv.FormatFloat(val, 'f', -1, 64) + "s")
	if err != nil {
		return &ConversionError{Value: val, ToType: durationType}
	}
	return ds.set(dval)
}

func (ds *durationSetter) SetBool(val bool) error {
	return &ConversionError{Value: val, ToType: durationType}
}

func (ds *durationSetter) Get() interface{} {
	if ds.val == nil {
		return time.Duration(0)
	}
	return *ds.val
}

type durationSetterCreator struct{}

func (durationSetterCreator) Type() reflect.Type {
	return durationType
}

func (durationSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	return &durationSetter{val: val.Addr().Interface().(*time.Duration), tag: tag}
}

func init() {
	DefaultSetterRegistry.Add(timeSetterCreator{})
	DefaultSetterRegistry.Add(durationSetterCreator{})
}
