package config

import (
	"reflect"
)

type ptrSetter struct {
	ptr           reflect.Value // The pointer to create/set on success.
	setterCreator SetterCreator
	tag           reflect.StructTag
}

func (ps *ptrSetter) IsBoolFlag() bool {
	tmp := reflect.New(ps.setterCreator.Type()).Elem()
	setter := ps.setterCreator.Setter(tmp, ps.tag)
	if setter, ok := setter.(interface{ IsBoolFlag() bool }); ok {
		return setter.IsBoolFlag()
	}
	return false
}

func (ps *ptrSetter) String() string {
	if ps.ptr.Kind() == reflect.Invalid || ps.ptr.IsNil() {
		return ""
	}
	return ps.setterCreator.Setter(ps.ptr.Elem(), ps.tag).String()
}

func (ps *ptrSetter) Set(s string) error {
	tmp := reflect.New(ps.setterCreator.Type()).Elem()
	if err := ps.setterCreator.Setter(tmp, ps.tag).Set(s); err != nil {
		return err
	}
	ps.ptr.Set(tmp.Addr())
	return nil
}

func (ps *ptrSetter) SetInt(i int64) error {
	tmp := reflect.New(ps.setterCreator.Type()).Elem()
	if err := ps.setterCreator.Setter(tmp, ps.tag).SetInt(i); err != nil {
		return err
	}
	ps.ptr.Set(tmp.Addr())
	return nil
}

func (ps *ptrSetter) SetUint(u uint64) error {
	tmp := reflect.New(ps.setterCreator.Type()).Elem()
	if err := ps.setterCreator.Setter(tmp, ps.tag).SetUint(u); err != nil {
		return err
	}
	ps.ptr.Set(tmp.Addr())
	return nil
}

func (ps *ptrSetter) SetFloat(f float64) error {
	tmp := reflect.New(ps.setterCreator.Type()).Elem()
	if err := ps.setterCreator.Setter(tmp, ps.tag).SetFloat(f); err != nil {
		return err
	}
	ps.ptr.Set(tmp.Addr())
	return nil
}

func (ps *ptrSetter) SetBool(b bool) error {
	tmp := reflect.New(ps.setterCreator.Type()).Elem()
	if err := ps.setterCreator.Setter(tmp, ps.tag).SetBool(b); err != nil {
		return err
	}
	ps.ptr.Set(tmp.Addr())
	return nil
}

func (ps *ptrSetter) Get() interface{} {
	if ps.ptr.Kind() == reflect.Invalid {
		return nil
	}
	return ps.ptr.Interface()
}

type ptrSetterCreator struct {
	setterCreator SetterCreator
}

func (psc *ptrSetterCreator) Type() reflect.Type {
	return reflect.PtrTo(psc.setterCreator.Type())
}

func (psc *ptrSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	if val.Kind() != reflect.Ptr {
		panic("value must be a pointer")
	}
	return &ptrSetter{
		ptr:           val,
		setterCreator: psc.setterCreator,
		tag:           tag,
	}
}

func newPtrSetterCreator(sc SetterCreator) SetterCreator {
	return &ptrSetterCreator{setterCreator: sc}
}
