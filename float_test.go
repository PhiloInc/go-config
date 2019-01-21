package config

import (
	"reflect"
	"testing"
)

func TestFloatSetter(t *testing.T) {
	creator := floatSetterCreator{t: float64Type}
	t.Run("String", func(t *testing.T) {
		val := 99.99
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "99.99" {
			t.Errorf("Returned string %s for value %f", text, val)
		}
		if text := (&floatSetter{}).String(); text != "0" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		var val float64
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("99.99"); err != nil {
			t.Errorf("Setting 99.99 failed with error %s", err)
		}
		if val != 99.99 {
			t.Errorf("Setting 99.99 resulted in value %f", val)
		}
		if err := s.Set("notafloat"); err == nil {
			t.Error("Setting notafloat did not fail with error")
		}
		if val != 99.99 {
			t.Errorf("Setting notafloat unexpectedly changed value to %f", val)
		}
		creator := floatSetterCreator{t: float32Type}
		var val32 float32
		s = creator.Setter(reflect.ValueOf(&val32).Elem(), "")
		floatStr := "1.797693134862315708145274237317043567981e+308"
		if err := s.Set(floatStr); err == nil {
			t.Errorf("Setting %s did not fail with error", floatStr)
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		var val float64
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetInt(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if val != 99 {
			t.Errorf("Setting 99 resulted in value %f", val)
		}
		bigInt := int64(1<<63 - 1)
		if err := s.SetInt(bigInt); err == nil {
			t.Errorf("Setting %d did not fail with error", bigInt)
		}
		if val != 99 {
			t.Errorf("Setting %d unexpectedly changed value to %f", bigInt, val)
		}
	})
	t.Run("SetUint", func(t *testing.T) {
		var val float64
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetUint(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if val != 99 {
			t.Errorf("Setting 99 resulted in value %f", val)
		}
		bigUint := uint64(1<<64 - 1)
		if err := s.SetUint(bigUint); err == nil {
			t.Errorf("Setting %d did not fail with error", bigUint)
		}
		if val != 99 {
			t.Errorf("Setting %d unexpectedly changed value to %f", bigUint, val)
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		var val float64
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(99.99); err != nil {
			t.Errorf("Setting 99.99 failed with error %s", err)
		}
		if val != 99.99 {
			t.Errorf("Setting 99.99 resulted in value %f", val)
		}
		creator := floatSetterCreator{t: float32Type}
		var val32 float32
		s = creator.Setter(reflect.ValueOf(&val32).Elem(), "")
		bigFloat := 1.797693134862315708145274237317043567981e+308
		if err := s.SetFloat(bigFloat); err == nil {
			t.Errorf("Setting %e did not fail with error", bigFloat)
		}
		if val != 99.99 {
			t.Errorf("Setting %e unexpectedly changed value to %f", bigFloat, val)
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		var val float64
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err != nil {
			t.Errorf("Setting true failed with error %s", err)
		}
		if val != 1 {
			t.Errorf("Setting true resulted in value %f", val)
		}
		if err := s.SetBool(false); err != nil {
			t.Errorf("Setting false failed with error %s", err)
		}
		if val != 0 {
			t.Errorf("Setting false resulted in value %f", val)
		}
	})
	t.Run("Get", func(t *testing.T) {
		val := 99.99
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch f := s.Get().(type) {
		case float64:
			if f != val {
				t.Errorf("Getting value %f returned %f", val, f)
			}
		default:
			t.Errorf("Getting value %f returned %v (type %T)", val, f, f)
		}
		creator := floatSetterCreator{t: float32Type}
		val32 := float32(99.99)
		s = creator.Setter(reflect.ValueOf(&val32).Elem(), "")
		switch f := s.Get().(type) {
		case float32:
			if f != val32 {
				t.Errorf("Getting value %f returned %f", val32, f)
			}
		default:
			t.Errorf("Getting value %f returned %v (type %T)", val32, f, f)
		}
	})
	t.Run("le", func(t *testing.T) {
		var val float64
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `le:"1"`)
		if err := s.Set("1"); err != nil {
			t.Error("validation <= 1 failed when setting 1")
		}
		if err := s.Set("2"); err == nil {
			t.Error("validation <= 1 did not fail when setting 2")
		}
		if val != 1 {
			t.Errorf("set invalid value %f", val)
		}
	})
	t.Run("ge", func(t *testing.T) {
		var val float64
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `ge:"1"`)
		if err := s.Set("1"); err != nil {
			t.Error("validation >= 1 failed when setting 1")
		}
		if err := s.Set("0.1"); err == nil {
			t.Error("validation >= 1 did not fail when setting 0.1")
		}
		if val != 1 {
			t.Errorf("set invalid value %f", val)
		}
	})
	t.Run("lt", func(t *testing.T) {
		var val float64
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `lt:"1"`)
		if err := s.Set("0.1"); err != nil {
			t.Error("validation < 1 failed when setting 0.1")
		}
		if err := s.Set("1"); err == nil {
			t.Error("validation < 1 did not fail when setting 1")
		}
		if val != 0.1 {
			t.Errorf("set invalid value %f", val)
		}
	})
	t.Run("gt", func(t *testing.T) {
		var val float64
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `gt:"1"`)
		if err := s.Set("1.1"); err != nil {
			t.Error("validation >= 1 failed when setting 1.1")
		}
		if err := s.Set("1"); err == nil {
			t.Error("validation >= 1 did not fail when setting 1")
		}
		if val != 1.1 {
			t.Errorf("set invalid value %f", val)
		}
	})
}
