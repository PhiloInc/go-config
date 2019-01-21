package config

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

type env map[string]interface{}

func (e *env) Set(name, val string) error {
	if *e == nil {
		*e = make(map[string]interface{})
	}
	if _, ok := (*e)[name]; !ok {
		if old, ok := os.LookupEnv(name); ok {
			(*e)[name] = old
		} else {
			(*e)[name] = false
		}
	}
	return os.Setenv(name, val)
}

func (e *env) Restore() {
	for k, v := range *e {
		switch v := v.(type) {
		case string:
			os.Setenv(k, v)
		case bool:
			os.Unsetenv(k)
		}
	}
	*e = nil
}

func TestEnvLoader(t *testing.T) {
	ptrStr := func(x interface{}) string {
		v := reflect.ValueOf(x)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return "<nil>"
			}
			return fmt.Sprint(v.Elem().Interface())
		}
		return fmt.Sprint(v.Interface())
	}
	root := NewRootPath("")
	node := root.AddNodePath(root.NewNodePath("test"))
	var duration []*time.Duration
	iterations := new(int)
	*iterations = 1 // set a default value.
	var name *string
	settings := settings{
		{
			Path: node.AddPath(node.NewPath("duration")),
			Setter: DefaultSetterRegistry.GetSetter(
				reflect.ValueOf(&duration).Elem(), `sep:","`,
			),
		},
		{
			Path: node.AddPath(node.NewPath("iterations")),
			Setter: DefaultSetterRegistry.GetSetter(
				reflect.ValueOf(&iterations).Elem(), "",
			),
		},
		{
			Path: node.AddPath(node.NewPath("name")),
			Setter: DefaultSetterRegistry.GetSetter(
				reflect.ValueOf(&name).Elem(), `regexp:"^\\pL*$"`,
			),
		},
	}
	t.Run("load", func(t *testing.T) {
		loader := new(EnvLoader)
		loader.Init(settings)
		var env env
		defer env.Restore()
		env.Set("TEST_NAME", "mytest")
		env.Set("TEST_DURATION", "30s")
		if err := loader.Load(); err != nil {
			t.Errorf("unexpected error parsing env: %s", err)
		}
		if len(duration) != 1 || duration[0] == nil || *duration[0] != 30*time.Second {
			t.Errorf("unexpected value %s for duration", duration)
		}
		if iterations == nil || *iterations != 1 {
			t.Errorf("unexpected value %s for iteration", ptrStr(iterations))
		}
		if name == nil || *name != "mytest" {
			t.Errorf("unexpected value %s for name", ptrStr(name))
		}
	})

	t.Run("default", func(t *testing.T) {
		loader := new(EnvLoader)
		loader.Init(settings)
		var env env
		defer env.Restore()
		env.Set("TEST_ITERATIONS", "2")
		if err := loader.Load(); err != nil {
			t.Errorf("unexpected error parsing env: %s", err)
		}
		if iterations == nil || *iterations != 2 {
			t.Errorf("unexpected value %s for iteration", ptrStr(iterations))
		}
	})

	t.Run("error", func(t *testing.T) {
		loader := new(EnvLoader)
		loader.Init(settings)
		var env env
		defer env.Restore()
		env.Set("TEST_NAME", "mytest2")
		if err := loader.Load(); err == nil {
			t.Errorf("parsing env did not return an error")
		}
		if name == nil || *name != "mytest" {
			t.Errorf("unexpected value %s for name", ptrStr(name))
		}
	})

	t.Run("multi", func(t *testing.T) {
		duration = nil
		loader := new(EnvLoader)
		loader.Init(settings)
		var env env
		defer env.Restore()
		env.Set("TEST_DURATION", "15s,30s")
		if err := loader.Load(); err != nil {
			t.Errorf("unexpected error parsing env: %s", err)
		}
		failed := (len(duration) != 2 ||
			duration[0] == nil ||
			*duration[0] != 15*time.Second ||
			duration[1] == nil ||
			*duration[1] != 30*time.Second)
		if failed {
			t.Errorf("unexpected value %s for duration", duration)
		}
	})
}
