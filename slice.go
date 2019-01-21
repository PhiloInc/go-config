package config

import (
	"reflect"
	"strconv"
	"strings"
)

type sliceSetter struct {
	append        bool
	slice         reflect.Value // The slice to set or append to on success.
	setterCreator SetterCreator
	tag           reflect.StructTag
}

func (ss *sliceSetter) IsBoolFlag() bool {
	tmp := reflect.New(ss.setterCreator.Type()).Elem()
	setter := ss.setterCreator.Setter(tmp, ss.tag)
	if setter, ok := setter.(interface{ IsBoolFlag() bool }); ok {
		return setter.IsBoolFlag()
	}
	return false
}

func (ss *sliceSetter) String() string {
	if ss.slice.Kind() == reflect.Invalid {
		return ""
	}
	l := ss.slice.Len()
	s := make([]string, ss.slice.Len())
	for i := 0; i < l; i++ {
		s[i] = ss.setterCreator.Setter(ss.slice.Index(i), "").String()
	}
	return strings.Join(s, ", ")
}

func (ss *sliceSetter) set(tmp reflect.Value) {
	if ss.append {
		ss.slice.Set(reflect.Append(ss.slice, tmp))
	} else {
		if ss.slice.Len() == 0 {
			t := reflect.SliceOf(ss.setterCreator.Type())
			ss.slice.Set(reflect.MakeSlice(t, 1, 1))
		} else {
			ss.slice.Set(ss.slice.Slice(0, 1))
		}
		ss.slice.Index(0).Set(tmp)
		ss.append = true
	}
}

func (ss *sliceSetter) Set(s string) error {
	vals := []string{s}
	if sep := ss.tag.Get("sep"); sep != "" {
		vals = strings.Split(s, sep)
	}
	var errs Errors
	for _, v := range vals {
		tmp := reflect.New(ss.setterCreator.Type()).Elem()
		if err := ss.setterCreator.Setter(tmp, ss.tag).Set(v); err != nil {
			errs.Append(err)
		} else {
			ss.set(tmp)
		}
	}
	return errs.AsError()
}

func (ss *sliceSetter) SetInt(i int64) error {
	tmp := reflect.New(ss.setterCreator.Type()).Elem()
	if err := ss.setterCreator.Setter(tmp, ss.tag).SetInt(i); err != nil {
		return err
	}
	ss.set(tmp)
	return nil
}

func (ss *sliceSetter) SetUint(u uint64) error {
	tmp := reflect.New(ss.setterCreator.Type()).Elem()
	if err := ss.setterCreator.Setter(tmp, ss.tag).SetUint(u); err != nil {
		return err
	}
	ss.set(tmp)
	return nil
}

func (ss *sliceSetter) SetFloat(f float64) error {
	tmp := reflect.New(ss.setterCreator.Type()).Elem()
	if err := ss.setterCreator.Setter(tmp, ss.tag).SetFloat(f); err != nil {
		return err
	}
	ss.set(tmp)
	return nil
}

func (ss *sliceSetter) SetBool(b bool) error {
	tmp := reflect.New(ss.setterCreator.Type()).Elem()
	if err := ss.setterCreator.Setter(tmp, ss.tag).SetBool(b); err != nil {
		return err
	}
	ss.set(tmp)
	return nil
}

func (ss *sliceSetter) Get() interface{} {
	if ss.slice.Kind() == reflect.Invalid {
		return nil
	}
	return ss.slice.Interface()
}

type sliceSetterCreator struct {
	setterCreator SetterCreator
}

func (ssc *sliceSetterCreator) Type() reflect.Type {
	return reflect.SliceOf(ssc.setterCreator.Type())
}

func (ssc *sliceSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	if val.Kind() != reflect.Slice {
		panic("value must be a slice")
	}
	append := true
	if tagVal, ok := tag.Lookup("append"); ok {
		append, _ = strconv.ParseBool(tagVal)
	}
	return &sliceSetter{
		append:        append,
		slice:         val,
		setterCreator: ssc.setterCreator,
		tag:           tag,
	}
}

func newSliceSetterCreator(vm SetterCreator) SetterCreator {
	return &sliceSetterCreator{setterCreator: vm}
}
