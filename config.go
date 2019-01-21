package config

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

/*
Config represents a configuration for an application.

It can be used in lieu of the package-level functions or DefaultConfig if
an application needs to manage multiple, separate configurations.

The zero value is ready to use.
*/
type Config struct {
	settings settings
	root     NodePath
	ptrs     [][2]reflect.Value
	loaders  Loaders
	reg      SetterRegistry
}

/*
GetSetterRegistry returns the SetterRegistry used by the Config.

It returns the last value passed to c.SetSetterRegistry, or
DefaultSetterRegistry if none is set.
*/
func (c *Config) GetSetterRegistry() SetterRegistry {
	if c.reg == nil {
		return DefaultSetterRegistry
	}
	return c.reg
}

/*
SetSetterRegistry sets a new SetterRegistry to be used by the Config.

This can be useful for providing custom SetterCreator implementations for new or
existing types without changing the global DefaultSetterRegistry.

Passing nil will cause the Config to use the DefaultSetterRegistry.
*/
func (c *Config) SetSetterRegistry(reg SetterRegistry) {
	c.reg = reg
}

/*
GetLoaders returns the current Loaders used by the Config.

It returns the last loaders set by SetLoaders, or a set of default loaders if
none was set.
*/
func (c *Config) GetLoaders() Loaders {
	if c.loaders == nil {
		c.loaders = GetDefaultLoaders()
	}
	return c.loaders
}

/*
SetLoaders sets loaders to be used by the Config.

It makes a copy of the parameter, so modifying the parameter after the function
returns is safe. Similarly, subsequently calling c.AddLoader will not modify the
parameter passed to this method.

Ordering is important, as subsequent loaders can overwrite configuration
values from earlier ones.

Passing a nil value will cause Config to use a default set of loaders.
*/
func (c *Config) SetLoaders(loaders Loaders) {
	c.loaders = loaders.Copy()
}

/*
AddLoader appends a loader to the set of loaders used by c.
*/
func (c *Config) AddLoader(loader Loader) {
	if c.loaders == nil {
		c.loaders = append(c.GetLoaders(), loader)
	} else {
		c.loaders.Add(loader)
	}
}

/*
Args returns the command-line arguments left after parsing.

Returns nil if Load has not been called, or if no flag loader was included in
the config.
*/
func (c *Config) Args() []string {
	// We only return something useful after Load() has been called, which means
	// c.loaders should be populated.
	for i := range c.loaders {
		if l, ok := c.loaders[i].(interface{ Args() []string }); ok {
			return l.Args()
		}
	}
	return nil
}

/*
Scan uses reflection to populate configuration settings from a struct.

The strct parameter must be a non-nil pointer to a struct type. Scan will panic
if any other value is passed.

This method is implicitly called by calling c.Configure. This method can be used
if multiple struct types are to be scanned; if additional settings are to be
added with c.Var before calling c.Load; or if unsupported type errors are to be
ignored. Note that any names scanned from the struct must be unique. If any
name is duplicated, Scan panics. This can only happen when calling Scan multiple
times or mixing calls to Var and Scan.

If the return value is non-nil, it will be of type Errors. Each element will in
turn be of type *UnknownTypeError, one for each scanned field for which a
Setter could not be created.

Only exported fields can be set. Unexported fields will be ignored. See the
package-level documentation for information about which field types are
supported by default and the effects of various struct tags.
*/
func (c *Config) Scan(strct interface{}) error {
	v := reflect.ValueOf(strct)
	if t := v.Type(); t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("parameter must be a pointer to a struct, not %s", t))
	}
	if v.IsNil() {
		panic(fmt.Sprintf("parameter can not be nil"))
	}
	return c.scan(v.Elem(), &c.root)
}

/*
Var adds a single value to be set by Config.

The value parameter must be a non-nil pointer. Passing any other value will
panic. The tag parameter can be used to set most tags as documented in the
package-level documentation, with the exception of 'config' and 'prefix' tags.

If value points to a struct, the function is equivalent to calling Scan with
value nested at the level indicated by name. For example:
	opts := struct {
		Level1: struct {
			Level2: struct {
				X int
			}
		}
	}
	Scan(&opts)
is equivalent to
	Var(&opts.Level1.Level2, "", "Level1", "Level2")

The remaining parameters are variadic, but at least one must be provided. If
only one value is provided, it is a non-prefixed name. If multiple values are
provided, leading values are prefixes. For example,
	DefaultConfig.Var(&version, "", "version")
will parse command-line flag -version and environment variable VERSION, while
	DefaultConfig.Var(&version, "", "api", "version")
parses command-line flag -api-version and environment variable API_VERSION.

Any name, with or without prefix, must be unique to the configuration. If it
duplicates a name already set with Var or parsed by Scan, Var panics.
*/
func (c *Config) Var(value interface{}, tag reflect.StructTag, name ...string) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("parameter 'value' must be a pointer, not %s", v.Type()))
	}
	if v.IsNil() {
		panic(fmt.Sprintf("parameter 'value' can not be nil"))
	}

	v = v.Elem()

	if len(name) == 0 {
		panic("parameter 'path' can not be empty")
	}

	lastPath := &c.root
	for i := 0; i < len(name)-1; i++ {
		if np := lastPath.FindNodePath(name[i]); np == nil {
			lastPath = lastPath.AddNodePath(lastPath.NewNodePath(name[i]))
		} else {
			lastPath = np
		}
	}

	setter := c.findSetter(v, tag)
	if setter == nil && v.Kind() == reflect.Struct {
		if np := lastPath.FindNodePath(name[len(name)-1]); np == nil {
			lastPath = lastPath.AddNodePath(lastPath.NewNodePath(name[len(name)-1]))
		} else {
			lastPath = np
		}
		return c.scan(v, lastPath)
	}

	p := lastPath.NewPath(name[len(name)-1])
	if setter == nil {
		return &UnknownTypeError{Type: v.Type(), Path: p}
	}

	c.settings.add(Setting{
		Path:   lastPath.AddPath(p),
		Tag:    tag,
		Setter: setter,
	})

	return nil
}

/*
Load calls loaders to parse and validate configuration.

This method is implicitly called by c.Configure. It is only necessary to call
it if not calling c.Configure. This method will only do something useful if
c.Scan or c.Var has been called to populate a list of values to set.

Any validation or load errors will result in a non-nil return status.

If one of the loaders is a *FlagLoader and the -h flag has not been overridden,
calling Load with "-h" in the application's command line arguments will cause
a list of settings and their descriptions to be printed to stderr.
*/
func (c *Config) Load() error {
	loaders := c.GetLoaders()
	settingsMap := c.settingsByLoader(loaders)
	var errs Errors
	for i := range loaders {
		if settings := settingsMap[loaders[i].Name()]; len(settings) != 0 {
			loaders[i].Init(settings)
			if loader, ok := loaders[i].(interface{ SetUsageFn(func()) }); ok {
				loader.SetUsageFn(func() { c.Usage(nil) })
			}
		}
	}
	for i := range loaders {
		if err := loaders[i].Load(); err != nil {
			errs.Append(err)
		}
	}
	c.setPtrs()
	return errs.AsError()
}

/*
Usage dumps usage information to an io.Writer.

Usage writes a list of setting names and their descriptions to w. If w is nil,
the list is dumped to os.Stderr.
*/
func (c *Config) Usage(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	loaders := c.GetLoaders()
	for i := range loaders {
		io.WriteString(w, loaders[i].Usage())
		io.WriteString(w, "\n")
	}
}

/*
Configure scans the strct argument to generate a list of parameters and then
loads them.

The strct parameter must be a non-nil pointer to a struct type. Any other value
panics.

It is equivalent to calling c.Scan followed by c.Load. The return value, if not
nil, will be of type Errors, and will contain the concatenated errors from both
Scan and Load.
*/
func (c *Config) Configure(strct interface{}) error {
	var errs Errors
	if err := c.Scan(strct); err != nil {
		errs.Append(err)
	}
	if err := c.Load(); err != nil {
		errs.Append(err)
	}
	return errs.AsError()
}

func (c *Config) scan(structVal reflect.Value, lastPath *NodePath) error {
	var errs Errors
	structType := structVal.Type()
	for i := 0; i < structType.NumField(); i++ {
		structField, fieldVal := structType.Field(i), structVal.Field(i)
		if !fieldVal.CanInterface() {
			continue
		}
		name := structField.Tag.Get("config")
		if name == "" {
			name = structField.Name
		} else if name == "-" {
			continue
		}
		p := lastPath.NewPath(name)
		if setter := c.findSetter(fieldVal, structField.Tag); setter != nil {
			c.settings.add(Setting{
				Path:   lastPath.AddPath(p),
				Tag:    structField.Tag,
				Setter: setter,
			})
			continue
		}
		// Find the last value that's a pointer and that isn't nil.
		for fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() {
			fieldVal = fieldVal.Elem()
		}
		if fieldVal.Kind() == reflect.Ptr {
			// We have a nil pointer. Iterate the type until we find a
			// non-pointer type. If it's a struct, create a temporary value
			// for fieldVal, and save the original (ptr) and temporary value
			// to c.ptrs.
			t := fieldVal.Type().Elem()
			for ; t.Kind() == reflect.Ptr; t = t.Elem() {
			}
			if t.Kind() == reflect.Struct {
				// we have to create the value with new, otherwise it won't be
				// addressable.
				ptrs := [2]reflect.Value{fieldVal, reflect.New(t).Elem()}
				fieldVal = ptrs[1]
				c.ptrs = append(c.ptrs, ptrs)
			}
		}
		if fieldVal.Kind() == reflect.Struct {
			prefix := structField.Tag.Get("prefix")
			if prefix == "" {
				prefix = name
			}
			var np *NodePath
			if prefix == "-" {
				np = lastPath
			} else {
				np = lastPath.FindNodePath(prefix)
				if np == nil {
					np = lastPath.AddNodePath(lastPath.NewNodePath(prefix))
				}
			}
			if err := c.scan(fieldVal, np); err != nil {
				errs.Append(err)
			}
		} else {
			errs.Append(&UnknownTypeError{Type: structField.Type, Path: p})
		}
	}
	return errs.AsError()
}

func (c *Config) setPtrs() {
	// Iterate in reverse of the order that this list was created. This causes
	// us to process more deeply nested structs before outer ones. Otherwise,
	// we will evaluate the shallower structs for the zero value before we've
	// had a chance to set pointers for their members.
	for i := len(c.ptrs) - 1; i >= 0; i-- {
		ptr, val := c.ptrs[i][0], c.ptrs[i][1]
		if reflect.Zero(val.Type()).Interface() != val.Interface() {
			for t := ptr.Type().Elem(); t.Kind() == reflect.Ptr; t = t.Elem() {
				ptr.Set(reflect.New(t))
				ptr = ptr.Elem()
			}
			ptr.Set(val.Addr())
		}
	}
}

func (c *Config) findSetter(val reflect.Value, tag reflect.StructTag) Setter {
	reg := c.reg
	if reg == nil {
		reg = DefaultSetterRegistry
	}
	return reg.GetSetter(val, tag)
}

func (c *Config) settingsByLoader(loaders []Loader) map[string]settings {
	m := make(map[string]settings, len(loaders))
	loaderNames := make([]string, len(loaders))
	for i := range loaders {
		loaderNames[i] = loaders[i].Name()
		m[loaderNames[i]] = make(settings, 0, len(c.settings))
	}
	for i := range c.settings {
		var keys []string
		if tag := c.settings[i].Tag.Get("from"); tag == "" || tag == "*" {
			keys = loaderNames
		} else {
			keys = strings.Split(tag, ",")
		}
		for _, k := range keys {
			if settings, ok := m[k]; ok {
				m[k] = append(settings, c.settings[i])
			}
		}
	}
	return m
}

/*
DefaultConfig is a default, global Config.

Is is used by the package-level Scan, Var and Configure functions, or can be
used directly by an application that doesn't require multiple configurations.
*/
var DefaultConfig = new(Config)

// Scan calls DefaultConfig.Scan.
func Scan(strct interface{}) error {
	return DefaultConfig.Scan(strct)
}

//Var calls DefaultConfig.Var
func Var(value interface{}, tag reflect.StructTag, name ...string) error {
	return DefaultConfig.Var(value, tag, name...)
}

//Load calls DefaultConfig.Load
func Load() error {
	return DefaultConfig.Load()
}

// Configure calls DefaultConfig.Configure
func Configure(strct interface{}) error {
	return DefaultConfig.Configure(strct)
}
