package config

import (
	"reflect"
	"testing"
)

func TestStringSetter(t *testing.T) {
	creator := stringSetterCreator{}
	t.Run("String", func(t *testing.T) {
		var val string
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "" {
			t.Errorf("Returned string %s for value %s", text, val)
		}
		val = "abc"
		if text := s.String(); text != "abc" {
			t.Errorf("Returned string %s for value %s", text, val)
		}
		if text := (&stringSetter{}).String(); text != "" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		var val string
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("myString"); err != nil {
			t.Errorf("Setting myString failed with error %s", err)
		}
		if val != "myString" {
			t.Errorf("Setting %s resulted in value %s", s, val)
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		var val string
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetInt(1); err != nil {
			t.Errorf("Setting 1 failed with error %s", err)
		}
		if val != "1" {
			t.Errorf("Setting 1 resulted in value %s", val)
		}
	})
	t.Run("SetUint", func(t *testing.T) {
		var val string
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetUint(1); err != nil {
			t.Errorf("Setting 1 failed with error %s", err)
		}
		if val != "1" {
			t.Errorf("Setting 1 resulted in value %s", val)
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		var val string
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(99.99); err != nil {
			t.Errorf("Setting 99.99 failed with error %s", err)
		}
		if val != "99.99" {
			t.Errorf("Setting 99.99 resulted in value %s", val)
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		var val string
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err != nil {
			t.Errorf("Setting true failed with error %s", err)
		}
		if val != "true" {
			t.Errorf("Setting true resulted in value %s", val)
		}
	})
	t.Run("Get", func(t *testing.T) {
		var val string
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch x := s.Get().(type) {
		case string:
			if x != "" {
				t.Errorf("Getting empty value returned %s", x)
			}
		default:
			t.Errorf("Getting empty value returned %v (type %T)", x, x)
		}
		val = "abc"
		switch x := s.Get().(type) {
		case string:
			if x != "abc" {
				t.Errorf("Getting value abc returned %s", x)
			}
		default:
			t.Errorf("Getting value abc returned %v (type %T)", x, x)
		}
	})
	t.Run("regexp", func(t *testing.T) {
		var val string
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `regexp:"^\\pL*$"`)
		if err := s.Set("abc"); err != nil {
			t.Error(`validation regexp:"^\\pL*$" failed when setting abc`)
		}
		if err := s.Set("123"); err == nil {
			t.Error(`validation regexp:"^\\pL*$" did not fail when setting 123`)
		}
	})
}
