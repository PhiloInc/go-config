package config

import (
	"reflect"
)

type ptrSetter struct {
	ptr    reflect.Value // The pointer to create/set on success.
	tmp    reflect.Value // The value used by 'Value'.
	setter Setter
}

func (ps *ptrSetter) IsBoolFlag() bool {
	if setter, ok := ps.setter.(interface{ IsBoolFlag() bool }); ok {
		return setter.IsBoolFlag()
	}
	return false
}

func (ps *ptrSetter) String() string {
	if ps.ptr.Kind() == reflect.Invalid || ps.ptr.IsNil() {
		return ""
	}
	return ps.setter.String()
}

func (ps *ptrSetter) Set(s string) error {
	if err := ps.setter.Set(s); err != nil {
		return err
	}
	ps.ptr.Set(ps.tmp.Addr())
	return nil
}

func (ps *ptrSetter) SetInt(i int64) error {
	if err := ps.setter.SetInt(i); err != nil {
		return err
	}
	ps.ptr.Set(ps.tmp.Addr())
	return nil
}

func (ps *ptrSetter) SetUint(u uint64) error {
	if err := ps.setter.SetUint(u); err != nil {
		return err
	}
	ps.ptr.Set(ps.tmp.Addr())
	return nil
}

func (ps *ptrSetter) SetFloat(f float64) error {
	if err := ps.setter.SetFloat(f); err != nil {
		return err
	}
	ps.ptr.Set(ps.tmp.Addr())
	return nil
}

func (ps *ptrSetter) SetBool(b bool) error {
	if err := ps.setter.SetBool(b); err != nil {
		return err
	}
	ps.ptr.Set(ps.tmp.Addr())
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
	tmp := reflect.New(psc.setterCreator.Type()).Elem()
	return &ptrSetter{
		ptr:    val,
		tmp:    tmp,
		setter: psc.setterCreator.Setter(tmp, tag),
	}
}

func newPtrSetterCreator(sc SetterCreator) SetterCreator {
	return &ptrSetterCreator{setterCreator: sc}
}
