package config

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestFlagLoader(t *testing.T) {
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
	var verbose *bool
	settings := settings{
		{
			Path: node.AddPath(node.NewPath("duration")),
			Setter: DefaultSetterRegistry.GetSetter(
				reflect.ValueOf(&duration).Elem(), "",
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
		{
			Path: node.AddPath(node.NewPath("verbose")),
			Setter: DefaultSetterRegistry.GetSetter(
				reflect.ValueOf(&verbose).Elem(), "",
			),
		},
	}

	t.Run("load", func(t *testing.T) {
		loader := new(FlagLoader)
		loader.Init(settings)
		args := []string{
			"-test-name", "mytest", "-test-duration", "30s", "-test-verbose",
		}
		if err := loader.Parse(args); err != nil {
			t.Errorf("unexpected error parsing args %s: %s", args, err)
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
		if verbose == nil || !*verbose {
			t.Errorf("unexpected value %s for verbose", ptrStr(name))
		}
	})

	t.Run("default", func(t *testing.T) {
		loader := new(FlagLoader)
		loader.Init(settings)
		args := []string{"-test-iterations", "2"}
		if err := loader.Parse(args); err != nil {
			t.Errorf("unexpected error parsing args %s: %s", args, err)
		}
		if iterations == nil || *iterations != 2 {
			t.Errorf("unexpected value %s for iteration", ptrStr(iterations))
		}
	})

	t.Run("error", func(t *testing.T) {
		loader := new(FlagLoader)
		loader.Init(settings)
		args := []string{"-test-name", "mytest2"}
		if err := loader.Parse(args); err == nil {
			t.Errorf("parsing args %s did not return an error", args)
		}
		if name == nil || *name != "mytest" {
			t.Errorf("unexpected value %s for name", ptrStr(name))
		}
	})

	t.Run("multi", func(t *testing.T) {
		duration = nil
		loader := new(FlagLoader)
		loader.Init(settings)
		args := []string{"-test-duration", "15s", "-test-duration", "30s"}
		if err := loader.Parse(args); err != nil {
			t.Errorf("unexpected error parsing args %s: %s", args, err)
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
