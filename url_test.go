package config

import (
	"net/url"
	"reflect"
	"testing"
)

func TestURLSetter(t *testing.T) {
	creator := urlSetterCreator{}
	t.Run("String", func(t *testing.T) {
		val := url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "https://www.philo.com/player/mytv" {
			t.Errorf("Returned string %s for value %s", text, &val)
		}
		if text := (&urlSetter{}).String(); text != "" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		var val url.URL
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("https://www.philo.com/player/mytv"); err != nil {
			t.Errorf("Setting https://www.philo.com/player/mytv failed with error %s", err)
		}
		if val != (url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}) {
			t.Errorf("Setting https://www.philo.com/player/mytv resulted in value %s", &val)
		}
		if err := s.Set(":badurl"); err == nil {
			t.Error("Setting :badurl did not fail with error")
		}
		if val != (url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}) {
			t.Errorf("Setting :badurl unexpectedly changed value to %s", &val)
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		val := url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetInt(1); err == nil {
			t.Error("Setting 1 did not fail with error")
		}
		if val != (url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}) {
			t.Errorf("Setting 1 unexpectedly changed value to %s", &val)
		}
	})
	t.Run("SetUint", func(t *testing.T) {
		val := url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetUint(1); err == nil {
			t.Error("Setting 1 did not fail with error")
		}
		if val != (url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}) {
			t.Errorf("Setting 1 unexpectedly changed value to %s", &val)
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		val := url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(99.99); err == nil {
			t.Error("Setting 99.99 did not fail with error")
		}
		if val != (url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}) {
			t.Errorf("Setting 99.99 unexpectedly changed value to %s", &val)
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		val := url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err == nil {
			t.Error("Setting true did not fail with error")
		}
		if val != (url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}) {
			t.Errorf("Setting true unexpectedly changed value to %s", &val)
		}
	})
	t.Run("Get", func(t *testing.T) {
		val := url.URL{Scheme: "https", Host: "www.philo.com", Path: "/player/mytv"}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch u := s.Get().(type) {
		case url.URL:
			if val != u {
				t.Errorf("Getting value %s returned %s", &val, &u)
			}
		default:
			t.Errorf("Getting value %s returned %v (type %T)", &val, u, u)
		}
	})
	t.Run("scheme", func(t *testing.T) {
		var val url.URL
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `scheme:"^(http|https)$"`)
		if err := s.Set("https://www.philo.com/player/mytv"); err != nil {
			t.Error(`validation scheme:"^(http|https)$" failed when setting https://www.philo.com/player/mytv`)
		}
		if err := s.Set("mailto:support@philo.com"); err == nil {
			t.Error(`validation scheme:"^(http|https)$" did not fail when setting mailto:support@philo.com`)
		}
	})
	t.Run("host", func(t *testing.T) {
		var val url.URL
		tag := reflect.StructTag(`host:"^[A-Za-z0-9]+(-[A-Za-z0-9]+)*\\.philo\\.com$"`)
		s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("https://www.philo.com/player/mytv"); err != nil {
			t.Errorf("validation %s failed when setting https://www.philo.com/player/mytv", tag)
		}
		if err := s.Set("ftp://ftp.kernel.org"); err == nil {
			t.Errorf("validation %s did not fail when setting ftp://ftp.kernel.org", tag)
		}
	})
	t.Run("path", func(t *testing.T) {
		var val url.URL
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `path:"^/player/"`)
		if err := s.Set("https://www.philo.com/player/mytv"); err != nil {
			t.Error(`validation path:"^/player/" failed when setting https://www.philo.com/player/mytv`)
		}
		if err := s.Set("https://www.philo.com/user/"); err == nil {
			t.Error(`validation path:"^/player/" did not fail when setting https://www.philo.com/user/`)
		}
	})
}
