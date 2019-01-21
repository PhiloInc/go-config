package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPtrSetter(t *testing.T) {
	creator := newPtrSetterCreator(intSetterCreator{t: int32Type})
	ptrStr := func(val *int32) string {
		if val == nil {
			return "<nil>"
		}
		return fmt.Sprintf("&%d", *val)
	}
	t.Run("String", func(t *testing.T) {
		// We always start out with a nil pointer. If the pointer was non-nil,
		// we wouldn't need to use the ptrSetter in the first place...
		var val *int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "" {
			t.Errorf("Returned string %s for nil value", text)
		}
		s.SetInt(99)
		if text := s.String(); text != "99" {
			t.Errorf("Returned string %s for value %s", text, ptrStr(val))
		}
		if text := (&ptrSetter{}).String(); text != "" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		var val *int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("notanint"); err == nil {
			t.Error("Setting notanint did not fail with error")
		}
		if val != nil {
			t.Errorf("Setting notanint unexpectedly changed value to %s", ptrStr(val))
		}
		if err := s.Set("99"); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if val == nil && *val != 99 {
			t.Errorf("Setting 99 resulted in value %s", ptrStr(val))
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		var val *int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		bigInt := int64(1<<63 - 1)
		if err := s.SetInt(bigInt); err == nil {
			t.Errorf("Setting %d did not fail with error", bigInt)
		}
		if val != nil {
			t.Errorf("Setting %d unexpectedly changed value to %s", bigInt, ptrStr(val))
		}
		if err := s.SetInt(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if val == nil || *val != 99 {
			t.Errorf("Setting 99 resulted in value %s", ptrStr(val))
		}

	})
	t.Run("SetUint", func(t *testing.T) {
		var val *int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		bigUint := uint64(1<<64 - 1)
		if err := s.SetUint(bigUint); err == nil {
			t.Errorf("Setting %d did not fail with error", bigUint)
		}
		if val != nil {
			t.Errorf("Setting %d unexpectedly changed value to %s", bigUint, ptrStr(val))
		}
		if err := s.SetUint(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if val == nil || *val != 99 {
			t.Errorf("Setting 99 resulted in value %s", ptrStr(val))
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		var val *int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(99.99); err == nil {
			t.Error("Setting 99.99 did not fail with error")
		}
		bigFloat := 1.797693134862315708145274237317043567981e+308
		if err := s.SetFloat(bigFloat); err == nil {
			t.Errorf("Setting %e did not fail with error", bigFloat)
		}
		if val != nil {
			t.Errorf("Setting %e unexpectedly changed value to %s", bigFloat, ptrStr(val))
		}
		if err := s.SetFloat(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if val == nil || *val != 99 {
			t.Errorf("Setting 99 resulted in value %s", ptrStr(val))
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		var val *int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err != nil {
			t.Errorf("Setting true failed with error %s", err)
		}
		if val == nil || *val != 1 {
			t.Errorf("Setting true resulted in value %s", ptrStr(val))
		}
		if err := s.SetBool(false); err != nil {
			t.Errorf("Setting false failed with error %s", err)
		}
		if val == nil || *val != 0 {
			t.Errorf("Setting false resulted in value %s", ptrStr(val))
		}
	})
	t.Run("Get", func(t *testing.T) {
		var val *int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch d := s.Get().(type) {
		case *int32:
			if d != nil {
				t.Errorf("Getting nil value returned %s", ptrStr(d))
			}
		default:
			t.Errorf("Getting nil value returned %v (type %T)", d, d)
		}
		s.SetInt(99)
		switch d := s.Get().(type) {
		case *int32:
			if d == nil || *d != 99 {
				t.Errorf("Getting value &99 returned %s", ptrStr(d))
			}
		default:
			t.Errorf("Getting nil value returned %v (type %T)", d, d)
		}
	})
}
