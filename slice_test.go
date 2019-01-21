package config

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestSliceSetter(t *testing.T) {
	creator := newSliceSetterCreator(intSetterCreator{t: int32Type})
	sliceStr := func(val []int32) string {
		s := make([]string, len(val))
		for i := range val {
			s[i] = strconv.FormatInt(int64(val[i]), 10)
		}
		return strings.Join(s, ",")
	}
	t.Run("String", func(t *testing.T) {
		var val []int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "" {
			t.Errorf("Returned string %s for nil value", text)
		}
		val = []int32{1, 2, 3}
		if text := s.String(); text != "1, 2, 3" {
			t.Errorf("Returned string %s for value %s", text, sliceStr(val))
		}
		if text := (&sliceSetter{}).String(); text != "" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		var val []int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("notanint"); err == nil {
			t.Error("Setting notanint did not fail with error")
		}
		if val != nil {
			t.Errorf("Setting notanint unexpectedly changed value to %s", sliceStr(val))
		}
		if err := s.Set("99"); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if len(val) != 1 || val[0] != 99 {
			t.Errorf("Setting 99 resulted in value %s", sliceStr(val))
		}
		if err := s.Set("100"); err != nil {
			t.Errorf("Setting 100 failed with error %s", err)
		}
		if len(val) != 2 || val[0] != 99 || val[1] != 100 {
			t.Errorf("Setting 99, 100 resulted in value %s", sliceStr(val))
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		var val []int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		bigInt := int64(1<<63 - 1)
		if err := s.SetInt(bigInt); err == nil {
			t.Errorf("Setting %d did not fail with error", bigInt)
		}
		if val != nil {
			t.Errorf("Setting %d unexpectedly changed value to %s", bigInt, sliceStr(val))
		}
		if err := s.SetInt(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if len(val) != 1 || val[0] != 99 {
			t.Errorf("Setting 99 resulted in value %s", sliceStr(val))
		}
		if err := s.SetInt(100); err != nil {
			t.Errorf("Setting 100 failed with error %s", err)
		}
		if len(val) != 2 || val[0] != 99 || val[1] != 100 {
			t.Errorf("Setting 99, 100 resulted in value %s", sliceStr(val))
		}

	})
	t.Run("SetUint", func(t *testing.T) {
		var val []int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		bigUint := uint64(1<<64 - 1)
		if err := s.SetUint(bigUint); err == nil {
			t.Errorf("Setting %d did not fail with error", bigUint)
		}
		if val != nil {
			t.Errorf("Setting %d unexpectedly changed value to %s", bigUint, sliceStr(val))
		}
		if err := s.SetUint(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if len(val) != 1 || val[0] != 99 {
			t.Errorf("Setting 99 resulted in value %s", sliceStr(val))
		}
		if err := s.SetUint(100); err != nil {
			t.Errorf("Setting 100 failed with error %s", err)
		}
		if len(val) != 2 || val[0] != 99 || val[1] != 100 {
			t.Errorf("Setting 99, 100 resulted in value %s", sliceStr(val))
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		var val []int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(99.99); err == nil {
			t.Error("Setting 99.99 did not fail with error")
		}
		bigFloat := 1.797693134862315708145274237317043567981e+308
		if err := s.SetFloat(bigFloat); err == nil {
			t.Errorf("Setting %e did not fail with error", bigFloat)
		}
		if val != nil {
			t.Errorf("Setting %e unexpectedly changed value to %s", bigFloat, sliceStr(val))
		}
		if err := s.SetFloat(99); err != nil {
			t.Errorf("Setting 99 failed with error %s", err)
		}
		if len(val) != 1 || val[0] != 99 {
			t.Errorf("Setting 99 resulted in value %s", sliceStr(val))
		}
		if err := s.SetFloat(100); err != nil {
			t.Errorf("Setting 100 failed with error %s", err)
		}
		if len(val) != 2 || val[0] != 99 || val[1] != 100 {
			t.Errorf("Setting 99, 100 resulted in value %s", sliceStr(val))
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		var val []int32
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err != nil {
			t.Errorf("Setting true failed with error %s", err)
		}
		if len(val) != 1 || val[0] != 1 {
			t.Errorf("Setting true resulted in value %s", sliceStr(val))
		}
		if err := s.SetBool(false); err != nil {
			t.Errorf("Setting false failed with error %s", err)
		}
		if len(val) != 2 || val[0] != 1 || val[1] != 0 {
			t.Errorf("Setting true, false resulted in value %s", sliceStr(val))
		}
	})
	t.Run("Get", func(t *testing.T) {
		val := []int32{1, 2, 3}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch d := s.Get().(type) {
		case []int32:
			if len(d) != 3 || d[0] != 1 || d[1] != 2 || d[2] != 3 {
				t.Errorf("Getting value 1, 2, 3 returned %s", sliceStr(val))
			}
		default:
			t.Errorf("Getting nil value returned %v (type %T)", d, d)
		}
		val = nil
		switch d := s.Get().(type) {
		case []int32:
			if d != nil {
				t.Errorf("Getting nil value returned %s", sliceStr(val))
			}
		default:
			t.Errorf("Getting nil value returned %v (type %T)", d, d)
		}
	})
}
