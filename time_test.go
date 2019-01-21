package config

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeSetter(t *testing.T) {
	creator := timeSetterCreator{}
	loc := time.FixedZone("UTC-8", -8*60*60)
	t.Run("String", func(t *testing.T) {
		val := time.Date(2019, 1, 19, 10, 0, 0, 0, time.UTC)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "2019-01-19T10:00:00Z" {
			t.Errorf("Returned string %s for value %s", text, val)
		}
		if text := (&timeSetter{}).String(); text != "0001-01-01T00:00:00Z" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		type value struct {
			value string
			time  time.Time
		}
		values := []value{
			{"2019-01-19T10:04:34.51861-08:00", time.Date(2019, 1, 19, 10, 4, 34, 518610000, loc)},
			{"2019-01-19 10:04:34.51861-08:00", time.Date(2019, 1, 19, 10, 4, 34, 518610000, loc)},
			{"2019-01-19T10:04:34.51861", time.Date(2019, 1, 19, 10, 4, 34, 518610000, time.UTC)},
			{"2019-01-19 10:04:34.51861", time.Date(2019, 1, 19, 10, 4, 34, 518610000, time.UTC)},
			{"2019-01-19T10:04:34-08:00", time.Date(2019, 1, 19, 10, 4, 34, 0, loc)},
			{"2019-01-19 10:04:34-08:00", time.Date(2019, 1, 19, 10, 4, 34, 0, loc)},
			{"2019-01-19T10:04:34", time.Date(2019, 1, 19, 10, 4, 34, 0, time.UTC)},
			{"2019-01-19 10:04:34", time.Date(2019, 1, 19, 10, 4, 34, 0, time.UTC)},
			{"2019-01-19T10:04-08:00", time.Date(2019, 1, 19, 10, 4, 0, 0, loc)},
			{"2019-01-19 10:04-08:00", time.Date(2019, 1, 19, 10, 4, 0, 0, loc)},
			{"2019-01-19T10:04", time.Date(2019, 1, 19, 10, 4, 0, 0, time.UTC)},
			{"2019-01-19 10:04", time.Date(2019, 1, 19, 10, 4, 0, 0, time.UTC)},
			{"2019-01-19T10", time.Date(2019, 1, 19, 10, 0, 0, 0, time.UTC)},
			{"2019-01-19 10", time.Date(2019, 1, 19, 10, 0, 0, 0, time.UTC)},
			{"2019-01-19", time.Date(2019, 1, 19, 0, 0, 0, 0, time.UTC)},
			{"2019-01", time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)},
		}
		for i := range values {
			var val time.Time
			s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
			if err := s.Set(values[i].value); err != nil {
				t.Errorf("Setting %s failed with error %s", values[i].value, err)
			}
			if !val.Equal(values[i].time) {
				t.Errorf(
					"Setting %s resulted in value %s",
					values[i].value, val.Format(time.RFC3339Nano),
				)
			}
		}
		val := time.Date(2019, 1, 19, 10, 4, 34, 518610000, loc)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("notatime"); err == nil {
			t.Error("Setting notatime did not fail with error")
		}
		if !val.Equal(time.Date(2019, 1, 19, 10, 4, 34, 518610000, loc)) {
			t.Errorf(
				"Setting notaduration unexpectedly changed value to %s",
				val.Format(time.RFC3339Nano),
			)
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		var val time.Time
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetInt(1547921074); err != nil {
			t.Errorf("Setting 1547921074 failed with error %s", err)
		}
		if !val.Equal(time.Date(2019, 1, 19, 18, 4, 34, 0, time.UTC)) {
			t.Errorf("Setting 1547921074 resulted in value %d", val.Unix())
		}
	})
	t.Run("SetUint", func(t *testing.T) {
		var val time.Time
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetUint(1547921074); err != nil {
			t.Errorf("Setting 1547921074 failed with error %s", err)
		}
		if !val.Equal(time.Date(2019, 1, 19, 18, 4, 34, 0, time.UTC)) {
			t.Errorf("Setting 1547921074 resulted in value %d", val.Unix())
		}
		bigUint := uint64(1<<64 - 1)
		if err := s.SetUint(bigUint); err == nil {
			t.Errorf("Setting %d did not fail with error", bigUint)
		}
		if !val.Equal(time.Date(2019, 1, 19, 18, 4, 34, 0, time.UTC)) {
			t.Errorf("Setting %d unexpectedly changed value to %d", bigUint, val.Unix())
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		var val time.Time
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(1547921074.51861); err != nil {
			t.Errorf("Setting 1547921074.51861 failed with error %s", err)
		}
		if !val.Equal(time.Date(2019, 1, 19, 18, 4, 34, 518610000, time.UTC)) {
			t.Errorf("Setting 1547921074.51861 resulted in value %s", val)
		}
		bigFloat := 1.797693134862315708145274237317043567981e+308
		if err := s.SetFloat(bigFloat); err == nil {
			t.Errorf("Setting %e did not fail with error", bigFloat)
		}
		if !val.Equal(time.Date(2019, 1, 19, 18, 4, 34, 518610000, time.UTC)) {
			t.Errorf("Setting %e unexpectedly changed value to %d", bigFloat, val.Unix())
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		val := time.Date(2019, 1, 19, 10, 0, 0, 0, time.UTC)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err == nil {
			t.Error("Setting true did not fail with error")
		}
		if !val.Equal(time.Date(2019, 1, 19, 10, 0, 0, 0, time.UTC)) {
			t.Errorf("Setting true unexpectedly changed value to %s", val.Format(time.RFC3339Nano))
		}
	})
	t.Run("Get", func(t *testing.T) {
		val := time.Date(2019, 1, 19, 10, 0, 0, 0, time.UTC)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch tm := s.Get().(type) {
		case time.Time:
			if !tm.Equal(time.Date(2019, 1, 19, 10, 0, 0, 0, time.UTC)) {
				t.Errorf(
					"Getting value %s returned %s",
					val.Format(time.RFC3339Nano),
					tm.Format(time.RFC3339Nano),
				)
			}
		default:
			t.Errorf(
				"Getting value %s returned %v (type %T)",
				val.Format(time.RFC3339Nano), tm, tm,
			)
		}
	})
	t.Run("le", func(t *testing.T) {
		var val time.Time
		tag := reflect.StructTag(`le:"2019-01-19T10:00:00Z"`)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("2019-01-19T10:00:00Z"); err != nil {
			t.Errorf("validation %s failed when setting 2019-01-19T10:00:00Z", tag)
		}
		if err := s.Set("2019-01-19T10:00:01Z"); err == nil {
			t.Errorf("validation %s did not fail when setting 2019-01-19T10:00:01Z", tag)
		}
		if !val.Equal(time.Date(2019, 1, 19, 10, 0, 0, 0, time.UTC)) {
			t.Errorf("set invalid value %s", val.Format(time.RFC3339))
		}
	})
	t.Run("ge", func(t *testing.T) {
		var val time.Time
		tag := reflect.StructTag(`ge:"2019-01-19T10:00:00Z"`)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("2019-01-19T10:00:00Z"); err != nil {
			t.Errorf("validation %s failed when setting 2019-01-19T10:00:00Z", tag)
		}
		if err := s.Set("2019-01-19T09:59:59Z"); err == nil {
			t.Errorf("validation %s did not fail when setting 2019-01-19T09:59:59Z", tag)
		}
		if !val.Equal(time.Date(2019, 1, 19, 10, 0, 0, 0, time.UTC)) {
			t.Errorf("set invalid value %s", val.Format(time.RFC3339))
		}
	})
	t.Run("lt", func(t *testing.T) {
		var val time.Time
		tag := reflect.StructTag(`lt:"2019-01-19T10:00:00Z"`)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("2019-01-19T09:59:59Z"); err != nil {
			t.Errorf("validation %s failed when setting 2019-01-19T09:59:59Z", tag)
		}
		if err := s.Set("2019-01-19T10:00:00Z"); err == nil {
			t.Errorf("validation %s did not fail when setting 2019-01-19T10:00:00Z", tag)
		}
		if !val.Equal(time.Date(2019, 1, 19, 9, 59, 59, 0, time.UTC)) {
			t.Errorf("set invalid value %s", val.Format(time.RFC3339))
		}
	})
	t.Run("gt", func(t *testing.T) {
		var val time.Time
		tag := reflect.StructTag(`gt:"2019-01-19T10:00:00Z"`)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("2019-01-19T10:00:01Z"); err != nil {
			t.Errorf("validation %s failed when setting 2019-01-19T10:00:01Z", tag)
		}
		if err := s.Set("2019-01-19T10:00:00Z"); err == nil {
			t.Errorf("validation %s did not fail when setting 2019-01-19T10:00:00Z", tag)
		}
		if !val.Equal(time.Date(2019, 1, 19, 10, 0, 1, 0, time.UTC)) {
			t.Errorf("set invalid value %s", val.Format(time.RFC3339))
		}
	})
}

func TestDurationSetter(t *testing.T) {
	creator := durationSetterCreator{}
	t.Run("String", func(t *testing.T) {
		val := time.Second
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "1s" {
			t.Errorf("Returned string %s for value %d", text, val)
		}
		if text := (&durationSetter{}).String(); text != "0s" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		var val time.Duration
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("1s"); err != nil {
			t.Errorf("Setting 1s failed with error %s", err)
		}
		if val != time.Second {
			t.Errorf("Setting 1s resulted in value %s", val)
		}
		if err := s.Set("notaduration"); err == nil {
			t.Error("Setting notaduration did not fail with error")
		}
		if val != time.Second {
			t.Errorf("Setting notaduration unexpectedly changed value to %s", val)
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		var val time.Duration
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetInt(1); err != nil {
			t.Errorf("Setting 1 failed with error %s", err)
		}
		if val != time.Second {
			t.Errorf("Setting 1 resulted in value %s", val)
		}
	})
	t.Run("SetUint", func(t *testing.T) {
		var val time.Duration
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetUint(1); err != nil {
			t.Errorf("Setting 1 failed with error %s", err)
		}
		if val != time.Second {
			t.Errorf("Setting 1 resulted in value %s", val)
		}
		bigUint := uint64(1<<64 - 1)
		if err := s.SetUint(bigUint); err == nil {
			t.Errorf("Setting %d did not fail with error", bigUint)
		}
		if val != time.Second {
			t.Errorf("Setting %d unexpectedly changed value to %s", bigUint, val)
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		var val time.Duration
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(1.001); err != nil {
			t.Errorf("Setting 1.001 failed with error %s", err)
		}
		if val != time.Second+time.Millisecond {
			t.Errorf("Setting 1.001 resulted in value %s", val)
		}
		bigFloat := 1.797693134862315708145274237317043567981e+308
		if err := s.SetFloat(bigFloat); err == nil {
			t.Errorf("Setting %e did not fail with error", bigFloat)
		}
		if val != time.Second+time.Millisecond {
			t.Errorf("Setting %e unexpectedly changed value to %s", bigFloat, val)
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		val := time.Second
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err == nil {
			t.Error("Setting true did not fail with error")
		}
		if val != time.Second {
			t.Errorf("Setting true unexpectedly changed value to %s", val)
		}
	})
	t.Run("Get", func(t *testing.T) {
		val := time.Second
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch d := s.Get().(type) {
		case time.Duration:
			if d != time.Second {
				t.Errorf("Getting value %s returned %s", val, d)
			}
		default:
			t.Errorf("Getting value %s returned %v (type %T)", val, d, d)
		}
	})
	t.Run("le", func(t *testing.T) {
		var val time.Duration
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `le:"1s"`)
		if err := s.Set("1s"); err != nil {
			t.Error("validation <= 1s failed when setting 1s")
		}
		if err := s.Set("2s"); err == nil {
			t.Error("validation <= 1s did not fail when setting 2s")
		}
		if val != time.Second {
			t.Errorf("set invalid value %s", val)
		}
	})
	t.Run("ge", func(t *testing.T) {
		var val time.Duration
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `ge:"1s"`)
		if err := s.Set("1s"); err != nil {
			t.Error("validation >= 1s failed when setting 1s")
		}
		if err := s.Set("0.5s"); err == nil {
			t.Error("validation >= 1s did not fail when setting 0.5s")
		}
		if val != time.Second {
			t.Errorf("set invalid value %s", val)
		}
	})
	t.Run("lt", func(t *testing.T) {
		var val time.Duration
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `lt:"1s"`)
		if err := s.Set("0.5s"); err != nil {
			t.Error("validation < 1s failed when setting 0.5s")
		}
		if err := s.Set("1s"); err == nil {
			t.Error("validation < 1s did not fail when setting 1s")
		}
		if val != time.Second/2 {
			t.Errorf("set invalid value %s", val)
		}
	})
	t.Run("gt", func(t *testing.T) {
		var val time.Duration
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `gt:"1s"`)
		if err := s.Set("2s"); err != nil {
			t.Error("validation >= 1s failed when setting 2s")
		}
		if err := s.Set("1s"); err == nil {
			t.Error("validation >= 1s did not fail when setting 1s")
		}
		if val != time.Second*2 {
			t.Errorf("set invalid value %s", val)
		}
	})
}
