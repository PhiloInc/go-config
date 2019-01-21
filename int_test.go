package config

import (
	"reflect"
	"testing"
)

func TestIntSetter(t *testing.T) {
	creator := intSetterCreator{t: intType}
	t.Run("String", func(t *testing.T) {
		val := 99
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "99" {
			t.Errorf("Returned string %s for value %d", text, val)
		}
		if text := (&intSetter{}).String(); text != "0" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		var val int
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
		var val8 int8
		creator := intSetterCreator{t: int8Type} // shadow creator
		s = creator.Setter(reflect.ValueOf(&val8).Elem(), "")
		intStr := "1024"
		if err := s.Set(intStr); err == nil {
			t.Errorf("Setting %s did not fail with error", intStr)
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		creator := intSetterCreator{t: int32Type}
		var val int32
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
		creator := intSetterCreator{t: int64Type}
		var val int64
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
		creator := intSetterCreator{t: int32Type}
		var val int32
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
		var val int
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
		val := 99
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch d := s.Get().(type) {
		case int:
			if d != val {
				t.Errorf("Getting value %d returned %d", val, d)
			}
		default:
			t.Errorf("Getting value %d returned %v (type %T)", val, d, d)
		}
		creator := intSetterCreator{t: int16Type}
		val16 := int16(99)
		s = creator.Setter(reflect.ValueOf(&val16).Elem(), "")
		switch d := s.Get().(type) {
		case int16:
			if d != val16 {
				t.Errorf("Getting value %d returned %d", val16, d)
			}
		default:
			t.Errorf("Getting value %d returned %v (type %T)", val16, d, d)
		}
	})
	t.Run("le", func(t *testing.T) {
		var val int
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
		var val int
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
		creator := intSetterCreator{t: int64Type}
		var val int64
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
		var val int
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
