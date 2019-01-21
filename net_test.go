package config

import (
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestIPSetter(t *testing.T) {
	creator := ipSetterCreator{}
	t.Run("String", func(t *testing.T) {
		val := net.IP{127, 0, 0, 1}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "127.0.0.1" {
			t.Errorf("Returned string %s for value %s", text, val)
		}
		if text := (&ipSetter{}).String(); text != "<nil>" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		var val net.IP
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("127.0.0.1"); err != nil {
			t.Errorf("Setting 127.0.0.1 failed with error %s", err)
		}
		if !val.Equal(net.IP{127, 0, 0, 1}) {
			t.Errorf("Setting 127.0.0.1 resulted in value %s", val)
		}
		if err := s.Set("notanip"); err == nil {
			t.Error("Setting notanip did not fail with error")
		}
		if !val.Equal(net.IP{127, 0, 0, 1}) {
			t.Errorf("Setting notanip unexpectedly changed value to %s", val)
		}
		if err := s.Set("::1"); err != nil {
			t.Errorf("Setting ::1 failed with error %s", err)
		}
		if !val.Equal(net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}) {
			t.Errorf("Setting ::1 resulted in value %s", val)
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		val := net.IP{127, 0, 0, 1}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetInt(1); err == nil {
			t.Error("Setting 1 did not fail with error")
		}
		if !val.Equal(net.IP{127, 0, 0, 1}) {
			t.Errorf("Setting 1 unexpectedly changed value to %s", val)
		}
	})
	t.Run("SetUint", func(t *testing.T) {
		val := net.IP{127, 0, 0, 1}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetUint(1); err == nil {
			t.Error("Setting 1 did not fail with error")
		}
		if !val.Equal(net.IP{127, 0, 0, 1}) {
			t.Errorf("Setting 1 unexpectedly changed value to %s", val)
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		val := net.IP{127, 0, 0, 1}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(99.99); err == nil {
			t.Error("Setting 99.99 did not fail with error")
		}
		if !val.Equal(net.IP{127, 0, 0, 1}) {
			t.Errorf("Setting 99.99 unexpectedly changed value to %s", val)
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		val := net.IP{127, 0, 0, 1}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err == nil {
			t.Error("Setting true did not fail with error")
		}
		if !val.Equal(net.IP{127, 0, 0, 1}) {
			t.Errorf("Setting true unexpectedly changed value to %s", val)
		}
	})
	t.Run("Get", func(t *testing.T) {
		val := net.IP{127, 0, 0, 1}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch ip := s.Get().(type) {
		case net.IP:
			if !val.Equal(ip) {
				t.Errorf("Getting value %s returned %s", val, ip)
			}
		default:
			t.Errorf("Getting value %s returned %v (type %T)", val, ip, ip)
		}
	})
	t.Run("version", func(t *testing.T) {
		var val net.IP
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `version:"4"`)
		if err := s.Set("127.0.0.1"); err != nil {
			t.Error(`validation version:"4" failed when setting 127.0.0.1`)
		}
		if err := s.Set("::1"); err == nil {
			t.Error(`validation version:"4" did not fail when setting ::1`)
		}
		if !val.Equal(net.IP{127, 0, 0, 1}) {
			t.Errorf("set invalid value %s", val)
		}
		s = creator.Setter(reflect.ValueOf(&val).Elem(), `version:"6"`)
		if err := s.Set("::1"); err != nil {
			t.Error(`validation version:"6" failed when setting ::1`)
		}
		if err := s.Set("127.0.0.1"); err == nil {
			t.Error(`validation version:"6" did not fail when setting 127.0.0.1`)
		}
		if !val.Equal(net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}) {
			t.Errorf("set invalid value %s", val)
		}
	})
	t.Run("is", func(t *testing.T) {
		values := [][2]string{
			[2]string{"global unicast", "10.0.0.1"},
			[2]string{"interface local multicast", "ff01::1"},
			[2]string{"link local multicast", "ff02::1"},
			[2]string{"link local unicast", "169.254.0.1"},
			[2]string{"loopback", "127.0.0.1"},
			[2]string{"multicast", "224.0.0.1"},
			[2]string{"unspecified", "0.0.0.0"},
		}
		for i := range values {
			spec, addr := values[i][0], values[i][1]
			var val net.IP
			tag := reflect.StructTag(fmt.Sprintf(`is:"%s"`, spec))
			s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
			if err := s.Set(addr); err != nil {
				t.Errorf("validation %s failed when setting %s: %s", tag, addr, err)
			}
			tag = reflect.StructTag(fmt.Sprintf(`is:"!%s"`, spec))
			s = creator.Setter(reflect.ValueOf(&val).Elem(), tag)
			if err := s.Set(addr); err == nil {
				t.Errorf("validation %s did not fail when setting %s", tag, addr)
			}
		}
		tag := reflect.StructTag(`is:"multicast,!interface local multicast"`)
		var val net.IP
		s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("ff02::1"); err != nil {
			t.Errorf("validation %s failed when setting ff02::1: %s", tag, err)
		}
		if err := s.Set("ff01::1"); err == nil {
			t.Errorf("validation %s did not fail when setting ff01::1", tag)
		}
		tag = reflect.StructTag(`is:"loopback,link local unicast"`)
		s = creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("127.0.0.1"); err != nil {
			t.Errorf("validation %s failed when setting 127.0.0.1: %s", tag, err)
		}
		if err := s.Set("169.254.0.1"); err != nil {
			t.Errorf("validation %s failed when setting 169.254.0.1: %s", tag, err)
		}
		if err := s.Set("10.0.0.1"); err == nil {
			t.Errorf("validation %s did not fail when setting 10.0.0.1", tag)
		}
	})
	t.Run("net", func(t *testing.T) {
		tag := reflect.StructTag(`net:"!127.0.0.1/32,127.0.0.0/8"`)
		var val net.IP
		s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("127.0.0.2"); err != nil {
			t.Errorf("validation %s failed when setting 127.0.0.2: %s", tag, err)
		}
		if err := s.Set("127.0.0.1"); err == nil {
			t.Errorf("validation %s did not fail when setting 127.0.0.1", tag)
		}
		tag = reflect.StructTag(`net:"169.254.0.0/16,127.0.0.0/8,!127.0.0.0/24"`)
		s = creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("169.254.0.1"); err != nil {
			t.Errorf("validation %s failed when setting 169.254.0.1: %s", tag, err)
		}
		if err := s.Set("127.1.0.1"); err != nil {
			t.Errorf("validation %s failed when setting 127.1.0.1: %s", tag, err)
		}
		if err := s.Set("127.0.0.1"); err == nil {
			t.Errorf("validation %s did not fail when setting 127.0.0.1", tag)
		}
	})
}

func TestIPNetSetter(t *testing.T) {
	creator := ipNetSetterCreator{}
	t.Run("String", func(t *testing.T) {
		val := net.IPNet{
			IP:   net.IP{127, 0, 0, 0},
			Mask: net.IPMask{255, 0, 0, 0},
		}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if text := s.String(); text != "127.0.0.0/8" {
			t.Errorf("Returned string %s for value %s", text, &val)
		}
		if text := (&ipNetSetter{}).String(); text != "<nil>" {
			t.Errorf("Returning string %s for zero setter", text)
		}
	})
	t.Run("Set", func(t *testing.T) {
		var val net.IPNet
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.Set("127.0.0.0/8"); err != nil {
			t.Errorf("Setting 127.0.0.1 failed with error %s", err)
		}
		if !val.IP.Equal(net.IP{127, 0, 0, 0}) || !net.IP(val.Mask).Equal(net.IP{255, 0, 0, 0}) {
			t.Errorf("Setting 127.0.0.0/8 resulted in value %s", &val)
		}
		if err := s.Set("notanipnet"); err == nil {
			t.Error("Setting notanipnet did not fail with error")
		}
		if !val.IP.Equal(net.IP{127, 0, 0, 0}) || !net.IP(val.Mask).Equal(net.IP{255, 0, 0, 0}) {
			t.Errorf("Setting notanint unexpectedly changed value to %s", &val)
		}
	})
	t.Run("SetInt", func(t *testing.T) {
		val := net.IPNet{
			IP:   net.IP{127, 0, 0, 0},
			Mask: net.IPMask{255, 0, 0, 0},
		}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetInt(1); err == nil {
			t.Error("Setting 1 did not fail with error")
		}
		if !val.IP.Equal(net.IP{127, 0, 0, 0}) || !net.IP(val.Mask).Equal(net.IP{255, 0, 0, 0}) {
			t.Errorf("Setting 1 unexpectedly changed value to %s", &val)
		}
	})
	t.Run("SetUint", func(t *testing.T) {
		val := net.IPNet{
			IP:   net.IP{127, 0, 0, 0},
			Mask: net.IPMask{255, 0, 0, 0},
		}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetUint(1); err == nil {
			t.Error("Setting 1 did not fail with error")
		}
		if !val.IP.Equal(net.IP{127, 0, 0, 0}) || !net.IP(val.Mask).Equal(net.IP{255, 0, 0, 0}) {
			t.Errorf("Setting 1 unexpectedly changed value to %s", &val)
		}
	})
	t.Run("SetFloat", func(t *testing.T) {
		val := net.IPNet{
			IP:   net.IP{127, 0, 0, 0},
			Mask: net.IPMask{255, 0, 0, 0},
		}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetFloat(99.99); err == nil {
			t.Error("Setting 99.99 did not fail with error")
		}
		if !val.IP.Equal(net.IP{127, 0, 0, 0}) || !net.IP(val.Mask).Equal(net.IP{255, 0, 0, 0}) {
			t.Errorf("Setting 99.99 unexpectedly changed value to %s", &val)
		}
	})
	t.Run("SetBool", func(t *testing.T) {
		val := net.IPNet{
			IP:   net.IP{127, 0, 0, 0},
			Mask: net.IPMask{255, 0, 0, 0},
		}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		if err := s.SetBool(true); err == nil {
			t.Error("Setting true did not fail with error")
		}
		if !val.IP.Equal(net.IP{127, 0, 0, 0}) || !net.IP(val.Mask).Equal(net.IP{255, 0, 0, 0}) {
			t.Errorf("Setting true unexpectedly changed value to %s", &val)
		}
	})
	t.Run("Get", func(t *testing.T) {
		val := net.IPNet{
			IP:   net.IP{127, 0, 0, 0},
			Mask: net.IPMask{255, 0, 0, 0},
		}
		s := creator.Setter(reflect.ValueOf(&val).Elem(), "")
		switch n := s.Get().(type) {
		case net.IPNet:
			ip, mask := net.IP{127, 0, 0, 0}, net.IP{255, 0, 0, 0}
			if !val.IP.Equal(ip) || !net.IP(val.Mask).Equal(mask) {
				t.Errorf("Getting value %s returned %s", &val, &n)
			}
		default:
			t.Errorf("Getting value %s returned %v (type %T)", &val, n, n)
		}
	})
	t.Run("version", func(t *testing.T) {
		var val net.IPNet
		s := creator.Setter(reflect.ValueOf(&val).Elem(), `version:"4"`)
		if err := s.Set("127.0.0.0/8"); err != nil {
			t.Error(`validation version:"4" failed when setting 127.0.0.0/8`)
		}
		if err := s.Set("::1/128"); err == nil {
			t.Error(`validation version:"4" did not fail when setting ::1/128`)
		}
		if !val.IP.Equal(net.IP{127, 0, 0, 0}) || !net.IP(val.Mask).Equal(net.IP{255, 0, 0, 0}) {
			t.Errorf("set invalid value %s", &val)
		}
		s = creator.Setter(reflect.ValueOf(&val).Elem(), `version:"6"`)
		if err := s.Set("::1/128"); err != nil {
			t.Error(`validation version:"6" failed when setting ::1/128`)
		}
		if err := s.Set("127.0.0.0/8"); err == nil {
			t.Error(`validation version:"6" did not fail when setting 127.0.0.0/8`)
		}
		ip := net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
		mask := net.IP{
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		}
		if !val.IP.Equal(ip) || !net.IP(val.Mask).Equal(mask) {
			t.Errorf("set invalid value %s", &val)
		}
	})
	t.Run("is", func(t *testing.T) {
		values := [][2]string{
			[2]string{"global unicast", "10.0.0.0/8"},
			[2]string{"interface local multicast", "ff01::db8:0:0/96"},
			[2]string{"link local multicast", "ff02::/104"},
			[2]string{"link local unicast", "169.254.0.0/16"},
			[2]string{"loopback", "127.0.0.0/8"},
			[2]string{"multicast", "224.0.0.0/3"},
			[2]string{"unspecified", "0.0.0.0/0"},
		}
		for i := range values {
			spec, addr := values[i][0], values[i][1]
			var val net.IPNet
			tag := reflect.StructTag(fmt.Sprintf(`is:"%s"`, spec))
			s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
			if err := s.Set(addr); err != nil {
				t.Errorf("validation %s failed when setting %s", tag, addr)
			}
			tag = reflect.StructTag(fmt.Sprintf(`is:"!%s"`, spec))
			s = creator.Setter(reflect.ValueOf(&val).Elem(), tag)
			if err := s.Set(addr); err == nil {
				t.Errorf("validation %s did not fail when setting %s", tag, addr)
			}
		}
		tag := reflect.StructTag(`is:"multicast,!interface local multicast"`)
		var val net.IPNet
		s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("ff02::/104"); err != nil {
			t.Errorf("validation %s failed when setting ff02::/104", tag)
		}
		if err := s.Set("ff01::db8:0:0/96"); err == nil {
			t.Errorf("validation %s did not fail when setting ff01::db8:0:0/96", tag)
		}
	})
	t.Run("net", func(t *testing.T) {
		tag := reflect.StructTag(`net:"!127.0.0.1/32,127.0.0.0/8"`)
		var val net.IPNet
		s := creator.Setter(reflect.ValueOf(&val).Elem(), tag)
		if err := s.Set("127.0.0.0/16"); err != nil {
			t.Errorf("validation %s failed when setting 127.0.0.0/16", tag)
		}
		if err := s.Set("127.0.0.1/32"); err == nil {
			t.Errorf("validation %s did not fail when setting 127.0.0.1/32", tag)
		}
	})
}
