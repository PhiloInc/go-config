package config

import (
	"reflect"
	"sort"
)

/*
Setting represents one configuration value.

It represents a value parsed from either a struct field or passed to Config.Var.
*/
type Setting struct {
	// Path is the given path within the configuration hierarchy. This should
	// be used by Loader implementations to create or identify the parameter.
	// it will never be nil when passed to Loader.Init().
	Path *Path
	// Tag is the struct tag parsed by Config.Scan or passed to Config.Var.
	Tag reflect.StructTag
	// Setter is an implementation of Setter. When passed to Loader.Init(), it
	// will never be nil.
	Setter
}

type settings []Setting

func (s *settings) add(setting Setting) {
	i := sort.Search(len(*s), func(i int) bool {
		return PathCmp(setting.Path, (*s)[i].Path) < 0
	})
	*s = append(*s, Setting{})
	copy((*s)[i+1:], (*s)[i:])
	(*s)[i] = setting
}
