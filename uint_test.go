package config

import (
	"reflect"
	"testing"
)

func TestUintSetter(t *testing.T) {
	creator := uintSetterCreator{t: uintType}
	t.Run("String", func(t *testing.T) {
		val := uint(99)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "99" {
			t.Errorf("Returned string %s for value %d", text, val)
		}
		if text := (&uintSetter{}).String(); text != "0" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		var val uint
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("99"); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if val != 99 {
			t.Errorf("Setting 99 resulted in value %d", val)
		}
		if err := s.Set("notanint"); err == nil {
			t.Error("Setting notanint did not fail with error")
		}
		if val != 99 {
			t.Errorf("Setting notanint unexpectedly changed value to %d", val)
		}
		creator := uintSetterCreator{t: uint8Type}
		var val8 uint8
		s = creator.Setter(reflect.ValueOf(&val8).Elem(), "")
		uintStr := "1024"
		if err := s.Set(uintStr); err == nil {
			t.Errorf("Setting %s did not fail with error", uintStr)
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		creator := uintSetterCreator{t: uint32Type}
		var val uint32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetInt(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if val != 99 {
			t.Errorf("Setting 99 resulted in value %d", val)
		}
		bigInt := int64(1<<63 - 1)
		if err := s.SetInt(bigInt); err == nil {
			t.Errorf("Setting %d did not fail with error", bigInt)
		}
		if val != 99 {
			t.Errorf("Setting %d unexpectedly changed value to %d", bigInt, val)
		}
	})
	t.Run("SetUint", func(t *testing.T) {
		creator := uintSetterCreator{t: uint32Type}
		var val uint32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetUint(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if val != 99 {
			t.Errorf("Setting 99 resulted in value %d", val)
		}
		bigUint := uint64(1<<64 - 1)
		if err := s.SetUint(bigUint); err == nil {
			t.Errorf("Setting %d did not fail with error", bigUint)
		}
		if val != 99 {
			t.Errorf("Setting %d unexpectedly changed value to %d", bigUint, val)
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		creator := uintSetterCreator{t: uint32Type}
		var val uint32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if val != 99 {
			t.Errorf("Setting 99 resulted in value %d", val)
		}
		if err := s.SetFloat(99.99); err == nil {
			t.Error("Setting 99.99 did not fail with error")
		}
		bigFloat := 1.797693134862315708145274237317043567981e+308
		if err := s.SetFloat(bigFloat); err == nil {
			t.Errorf("Setting %e did not fail with error", bigFloat)
		}
		if val != 99 {
			t.Errorf("Setting %e unexpectedly changed value to %d", bigFloat, val)
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		var val uint
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err != nil {
			t.Errorf("Setting true failed with error %s", err)
		}
		if val != 1 {
			t.Errorf("Setting true resulted in value %d", val)
		}
		if err := s.SetBool(false); err != nil {
			t.Errorf("Setting false failed with error %s", err)
		}
		if val != 0 {
			t.Errorf("Setting false resulted in value %d", val)
		}
	})
	t.Run("Get", func(t *testing.T) {
		val := uint(99)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch u := s.Get().(type) {
		case uint:
			if u != val {
				t.Errorf("Getting value %d returned %d", val, u)
			}
		default:
			t.Errorf("Getting value %d returned %v (type %T)", val, u, u)
		}
		creator := uintSetterCreator{t: uint16Type}
		val16 := uint16(99)
		s = creator.Setter(reflect.ValueOf(&val16).Elem(), "")
		switch u := s.Get().(type) {
		case uint16:
			if u != val16 {
				t.Errorf("Getting value %d returned %d", val16, u)
			}
		default:
			t.Errorf("Getting value %d returned %v (type %T)", val16, u, u)
		}
	})
	t.Run("le", func(t *testing.T) {
		var val uint
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `le:"1"`)
		if err := s.Set("1"); err != nil {
			t.Error("validation <= 1 failed when setting 1")
		}
		if err := s.Set("2"); err == nil {
			t.Error("validation <= 1 did not fail when setting 2")
		}
		if val != 1 {
			t.Errorf("set invalid value %d", val)
		}
	})
	t.Run("ge", func(t *testing.T) {
		var val uint
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `ge:"1"`)
		if err := s.Set("1"); err != nil {
			t.Error("validation >= 1 failed when setting 1")
		}
		if err := s.Set("0"); err == nil {
			t.Error("validation >= 1 did not fail when setting 0")
		}
		if val != 1 {
			t.Errorf("set invalid value %d", val)
		}
	})
	t.Run("lt", func(t *testing.T) {
		creator := uintSetterCreator{t: uint64Type}
		var val uint64
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `lt:"1"`)
		if err := s.Set("0"); err != nil {
			t.Error("validation < 1 failed when setting 0")
		}
		if err := s.Set("1"); err == nil {
			t.Error("validation < 1 did not fail when setting 1")
		}
		if val != 0 {
			t.Errorf("set invalid value %d", val)
		}
	})
	t.Run("gt", func(t *testing.T) {
		var val uint
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `gt:"1"`)
		if err := s.Set("2"); err != nil {
			t.Error("validation >= 1 failed when setting 2")
		}
		if err := s.Set("1"); err == nil {
			t.Error("validation >= 1 did not fail when setting 1")
		}
		if val != 2 {
			t.Errorf("set invalid value %d", val)
		}
	})
}
