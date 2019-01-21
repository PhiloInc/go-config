package config

import (
	"reflect"
	"testing"
)

func TestBoolSetter(t *testing.T) {
	creator := boolSetterCreator{}
	t.Run("String", func(t *testing.T) {
		val := false
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "false" {
			t.Errorf("Returned string %s for value %t", text, val)
		}
		val = true
		if text := s.String(); text != "true" {
			t.Errorf("Returned string %s for value %t", text, val)
		}
		if text := (&boolSetter{}).String(); text != "false" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		truthy := []string{"1", "t", "T", "TRUE", "true", "True"}
		for _, text := range truthy {
			val := false
			s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
			if err := s.Set(text); err != nil {
				t.Errorf("Setting %s failed with error %s", s, err)
			}
			if !val {
				t.Errorf("Setting %s resulted in value %t", s, val)
			}
		}
		falsy := []string{"0", "f", "F", "FALSE", "false", "False"}
		for _, text := range falsy {
			val := true
			s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
			if err := s.Set(text); err != nil {
				t.Errorf("Setting %s failed with error %s", s, err)
			}
			if val {
				t.Errorf("Setting %s resulted in value %t", s, val)
			}
		}
		val := true
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("notabool"); err == nil {
			t.Error("Setting notabool did not fail with error")
		}
		if !val {
			t.Error("Setting notabool unexpectedly changed value to false")
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		val := false
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetInt(1); err != nil {
			t.Errorf("Setting 1 failed with error %s", err)
		}
		if !val {
			t.Errorf("Setting 1 resulted in value %t", val)
		}
		val = true
		s = creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetInt(0); err != nil {
			t.Errorf("Setting 0 failed with error %s", err)
		}
		if val {
			t.Errorf("Setting 0 resulted in value %t", val)
		}
	})
	t.Run("SetUint", func(t *testing.T) {
		val := false
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetUint(1); err != nil {
			t.Errorf("Setting 1 failed with error %s", err)
		}
		if !val {
			t.Errorf("Setting 1 resulted in value %t", val)
		}
		val = true
		s = creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetUint(0); err != nil {
			t.Errorf("Setting 0 failed with error %s", err)
		}
		if val {
			t.Errorf("Setting 0 resulted in value %t", val)
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		val := false
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(0.1); err != nil {
			t.Errorf("Setting 0.1 failed with error %s", err)
		}
		if !val {
			t.Errorf("Setting 0.1 resulted in value %t", val)
		}
		val = true
		s = creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(0); err != nil {
			t.Errorf("Setting 0 failed with error %s", err)
		}
		if val {
			t.Errorf("Setting 0 resulted in value %t", val)
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		val := false
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err != nil {
			t.Errorf("Setting true failed with error %s", err)
		}
		if !val {
			t.Errorf("Setting true resulted in value %t", val)
		}
		val = true
		s = creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(false); err != nil {
			t.Errorf("Setting false failed with error %s", err)
		}
		if val {
			t.Errorf("Setting false resulted in value %t", val)
		}
	})
	t.Run("Get", func(t *testing.T) {
		val := false
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch b := s.Get().(type) {
		case bool:
			if b {
				t.Error("Getting false value returned true")
			}
		default:
			t.Errorf("Getting false value returned %v (type %T)", b, b)
		}
		val = true
		switch b := s.Get().(type) {
		case bool:
			if !b {
				t.Error("Getting true value returned false")
			}
		default:
			t.Errorf("Getting true value returned %v (type %T)", b, b)
		}
	})
}
