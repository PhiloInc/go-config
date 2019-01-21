package config

import (
	"flag"
	"reflect"
	"strings"
	"unicode"
)

/*
Loader is an interface for parsing and configuring settings using Config.
*/
type Loader interface {
	// Init is used to initialize a loader with a list of settings. Config
	// will always call it before Load. Implementations should support calling
	// it more than once.
	Init([]Setting)
	// Load will be called to parse and set variables. For each Setting passed
	// to Init, the implementation should look for an appropriate setting, and,
	// if found, call one of the 'Set*' methods on Setting.Setter.
	Load() error
	// Name should return a name for the loader. The name need not be unique.
	// However, the name will be used for matching against the `from:""` struct
	// tag.
	Name() string
	// Usage should return a string providing context-specific help for the
	// given loader. Config will always call Init before calling Usage, so
	// the implementation may use the Setting variables passed to Init to
	// provide setting-specific help.
	Usage() string
}

// Loaders provides a slice type for managing multiple loaders.
type Loaders []Loader

// Add appends a loader to the list of loaders.
func (l *Loaders) Add(loader Loader) {
	*l = append(*l, loader)
}

/*
Copy makes a copy of l.

This is primarily useful in order to modify the list of loaders based on an
existing template without modifying the original. Note that Loader
implementations may be pointer types, so while the slice is copied, individual
elements may point to shared data.
*/
func (l Loaders) Copy() Loaders {
	n := make(Loaders, len(l))
	copy(n, l)
	return n
}

/*
GetDefaultLoaders returns an ordered list of default loaders.

It currently returns {*EnvLoader, *FlagLoader}. This implies that command line
flags can override environment variables.

The returned values are always new values without other references, so can be
modified without affecting existing references.
*/
func GetDefaultLoaders() Loaders {
	return Loaders{
		new(EnvLoader),
		new(FlagLoader),
	}
}

/*
SplitName splits a name into logical words.

It attempts to parse names into words based on Go conventions and how humans
perceive them:
 * A capital letter following a lowercase one always starts a new word.
 * A series of capital letters are generally considered part of the same word.
   The last capital letter in the series will start a new word, except:
   * When it is followed by a lowercase letter and a digit, e.g., "v2".
   * When it is followed by a single lowercase letter and then the end of the
     string.
Examples:
	"HTMLEntityID" -> {"HTML", "Entity", "ID"}
	"UUID"         -> {"UUID"}
	"UUIDv2"       -> {"UUID", "v2"}
	"UUIDs"        -> {"UUIDs"}
	"IPv6Network"  -> {"IP", "v6", "Network"}

The return values are useful for creating parameter names that are readily
human-readable, such as for command line flags or environment variables.
*/
func SplitName(name string) []string {
	var b strings.Builder
	var ss []string
	runes := []rune(name)
	for i := range runes {
		split := false
		switch {
		case i == 0: // never split at beginning of string
		case i == len(runes)-1: // never split at end of string
		case unicode.IsUpper(runes[i]):
			switch {
			case !unicode.IsUpper(runes[i-1]):
				// first capital after a non-capital is always a new word.
				split = true
			case unicode.IsLower(runes[i+1]) && !(i == len(runes)-2 || unicode.IsDigit(runes[i+2])):
				// if we are a capital followed by a lowercase, unless that
				// lowercase is the end of the string or followed by a digit.
				split = true
			}
		case unicode.IsLower(runes[i]) && unicode.IsUpper(runes[i-1]) && unicode.IsDigit(runes[i+1]):
			// If the previous letter was a capital, and the current one is a
			// lowercase followed by a digit, then split.
			split = true
		}
		if split {
			ss = append(ss, b.String())
			b.Reset()
		}
		b.WriteRune(runes[i])
	}
	ss = append(ss, b.String())
	return ss
}

var ftnReplacer = strings.NewReplacer("*", "", " ", "", "[]", "")

/*
FriendlyTypeName returns a descriptive name for a value's type.

It is primarily useful for generating placeholder names for paremeters in help
output.

If types are element types, such as an array or pointer, they are dereferenced
until a non-element type is found. If that type has a name, it is lowercased
and returned. If the type is unnamed, the string "value" is returned.
*/
func FriendlyTypeName(x interface{}) string {
	t := reflect.TypeOf(x)
	for {
		switch t.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
			t = t.Elem()
			continue
		}
		break
	}
	if n := t.Name(); n != "" {
		return strings.ToLower(n)
	}
	return "value"
}

// The below is borrowed from Go's flag.go.
func isZeroValue(value flag.Value) bool {
	typ := reflect.TypeOf(value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return value.String() == z.Interface().(flag.Value).String()
}
