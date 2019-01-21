# go-config

A Configuration Package for Go

## Summary

go-config provides a convenient way to scan configuration settings from command
line flags and environment variables into structs. Its aim is to provide logical
organization of configuration parameters with minimal developer overhead.

It also provides for basic validation of configuration values by using struct
tags, and parsing of single values.

## Quick Start

```go
type Auth struct {
    User     string
    Password string
}

type Options struct {
    Auth    *Auth
    URLs    []*url.URL `scheme:"^(http|https)$"`
    Verbose int        `min:"0" max:"8"`
}

opts := Options{Verbose: 1}
if err := config.Configure(&opts); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
if opts.Auth != nil {
    fmt.Printf("Auth: %+v\n", *opts.Auth)
}
fmt.Printf("URLs: %s\nVerbose: %d\n", opts.URLs, opts.Verbose)
```

Running the program with -h demonstrates the generated command-line flags and
environment variable settings:

```text
Command Line Flags:
  -auth-password string
  -auth-user string
  -urls url
  -verbose int
         (default 1)

Environment Variables:
  AUTH_PASSWORD=string
  AUTH_USER=string
  URLS=url
  VERBOSE=int
         (default 1)

Help requested
```

Running the program with supported flags demonstrates how values are set:

```text
prog -urls http://www.google.com/ -urls http://www.yahoo.com/ -auth-user user1
Auth: {User:user1 Password:}
URLs: [http://www.google.com/ http://www.yahoo.com/]
Verbose: 1
```

To set a default value, simply set the value in the struct before calling
`Configure`. Struct tags can be used to specify permitted values. See the
[Validation](#validation) section for details.

## Detailed Usage

### Supported Types

The following types are supported:

* int
* int8
* int16
* int32
* int64
* uint
* uint8
* uint16
* uint32
* uint64
* float32
* float64
* string
* bool
* net.IP
* net.IPNet
* url.URL
* time.Duration
* time.Time

#### Pointers and Slices

Types that are pointer and/or slice types of a supported type are also
supported. For example, *[]int*, *\*int*, *\*[]\*int*, or any other combination
of slice and pointer indirection can be set.

#### Numeric types

int, int8, int16, int32, and int64, uint, uint8, uint16, uint32 and uint64 types
are parsed as base 10, unless they have a leading zero or *0x*, in which case
they are parsed as, respectively, base 8 or base 16 values. Values which would
overflow the type or discard a sign are considered invalid. For example,
assigning a value larger than 255 to int8, or assigning a negative value to a
uint.

float32 and float64 types are parsed using strconv.ParseFloat using the correct
size (bits).

#### Booleans

bool values are parsed using strconv.ParseBool. *1*, *t*, *T*, *TRUE*, *true*,
and *True* are parsed as `true`. *0*, *f*, *F*, *FALSE*, *false*, and *False*
evaluate to false. Any other value is an error.

#### URLs

url.URL values are parsed using url.Parse.

#### IP Addresses

net.IP values are are parsed with net.ParseIP and accept, e.g., *192.168.1.1*
or *2600:1700:5fa0:ef90:85e:79dc:5ea7:c711*. net.IPNet values are parsed with
net.ParseCIDR as address/bits, such as *169.254.0.0/16* or
*2001:cdba:9abc:5678::/64*. A value provided as a network address need not
necessarily specify the 0 network address: *169.254.1.1/16* will be understood
as *169.254.0.0/16*.

#### Durations

time.Duration values are parsed using time.ParseDuration, e.g., *0.75s* or
*1h37m27s*.

#### Times

A best effort is made to parse a time based on a static list of layouts:

* "2006-01-02T15:04:05Z07:00"
* "2006-01-02 15:04:05Z07:00"
* "2006-01-02T15:04:05"
* "2006-01-02 15:04:05"
* "2006-01-02T15:04Z07:00"
* "2006-01-02 15:04Z07:00"
* "2006-01-02T15:04"
* "2006-01-02 15:04"
* "2006-01-02T15"
* "2006-01-02 15"
* "2006-01-02"
* "2006-01"

If the time can not be parsed using any of the given layouts, it is an error.

### Flag and Environment Variable Names

go-config attempts to generate human-friendly names by parsing the names of
struct fields in a hierarchical way. For example, given

```go
type UserAccount struct {
    Network struct {
        IPv6FriendlyName string
    }
}
```

go-config will generate the command line flag
*-user-account-network-ip-v6-friendly-name*. See the section
[Overriding Names](#overriding-names) for details on how to change this
behavior.

### Validation

Note than for types in the below section, the validations also apply to other
types that are pointer and/or slice types of the base type. For example, if a
validation applies to *int*, it also applies to *[]int*, *\*int*, *\*[]\*int*,
or any other combination of slice and pointer indirection.

#### Greater or Less Than Comparisons

Number types int, int8, int16, int32, int64, uint, uint8, uint16, uint32,
uint64, float32, float64 as well as time.Duration and time.Time support a
number of struct tags for validation.

*ge* and *min* have the same meaning, and indicate a value must be greater than
or equal to the argument. *le* and *max* similarly indicate a value must be
less than or equal to a given value. Tags *lt* and *gt* indicate a strict
less than or greater than comparison (not equal to).

All arguments to the above tags should be given in a valid string for the type.
For example, *1s* for a duration, *2001-01-01 11:59:59Z* for a time, etc. The
special value *now* may be used for time.Time values.

The above tags can be combined to require values within a range.

Examples:

```go
type Example struct {
    Ratio   float64       `min:"0" lt:"1"`
    Timeout time.Duration `min:"0s"`
}
```

#### Regular Expressions

Fields of string type support the *regexp* tag. Only values that match the
given regular expression will be accepted.

Similarly, url.URL fields support *scheme*, *host* and *path* struct tags to
require that the URL scheme, host or path, respectively, match the given
regular expression.

Go's regular expression matching will match any part of a string. To match the
entire string, anchor the expression with `^` and `$`.

Struct tag values use quoted strings. This means that double quotes and
backslahes must be backslash-escaped.

Examples:

```go
type Example struct {
    BaseURL  *url.URL `scheme:"^(http|https)$"`
    Username string   `regexp:"^\\pL"` // usernames must start with a letter.
}
```

#### IP Addresses

net.IP and net.IPNet types support tags to validate a particular address.

##### is

The *is* tag can be used with values *global unicast*,
*interface local multicast*, *link local multicast*, *link local unicast*,
*loopback*, *multicast*, and *unspecified*. Each value matches the corresponding
range of IP addresses. A negative match can be specified by prefixing with a
`!`, and multiple values can be separated by commas. Validation succeeds if any
non-negated term matches the value being set. Validation fails, however, if any
negated term is matched.

Example:

```go
type Example struct {
    IP net.IP `is:"multicast,!interface local multicast"`
}
```

The above field will accept a multicast address, but not if it is an interface
local multicast.

##### net

The *net* tag can be used to specify restrictions for which networks an IP
address can be part of. Network addresses are specified in address/bits
notation, e.g., *169.254.0.0/16* or *2001:cdba:9abc:5678::/64*. As with the
*net* tag, a `!` negates the match term, and multiple terms can be seprated by
commas. Validation succeeds if any non-negated term matches, and fails if any
negated term is matched.

Example:

```go
type Example struct {
    IP net.IP `net:"192.168.0.0/16,!192.168.254.0/24"`
}
```

##### version

The *version* tag accepts values *4* or *6* to indicate that an address must be
an IPv4 or IPv6 address.

### Controlling Where Values are Loaded From

By default, values are parsed first from environment variables and then from
command line flags. If a value is set by both an environment variable and a
command line flag, the command line flag will overwrite the value set from the
environment variable.

The *from* struct tag can be used to control where values are loaded from. The
default is `from:"flag,env"`. To prevent a field from being set from an
environment variable, one could use `from:"flag"` in its struct tag.

The default loaders can also be overridden, either to change the order or to
exclude defaults. For example, to cause environment variables to overwrite
command line flags:

```go
loaders := config.Loaders{new(config.FlagLoader), new(config.EnvLoader)}
config.DefaultConfig.SetLoaders(loaders)
```

### Scanning Multiple Structs or Setting Individual Values

The *Configuration* function is convenient for scanning values from a single
struct and then loading them. However, functions are also provided for loading
values into more than one struct, or for setting individual variables. The
functions *Scan*, *Var* and *Load* are available.

To scan more than one struct, simply call *Scan* for each struct. Call *Var* to
add an individual variable. When finished, call *Load* to parse and set values
from the command line and/or environment.

Note, however, that name collisions will cause a panic.

### Using Non-global Configurations

For some cases, it may be desirable to avoid using the package global
*DefaultConfig*. In order to support standalone or multiple configurations,
simply use the config.Config type. The zero value is ready to use.

### Overriding Names

The *config* struct tag can be used to skip fields or change the names that
are generated for them. `config:"-"` will cause a field to be skipped.
`config:"OtherName"` will cause go-config to use *OtherName* as the name for
this field, which will be used in generating command line flag and environment
variable names.

Using `prefix:"-"` will cause a child struct to use the same prefix as the
parent struct. For example:

```go
type Auth struct {
    User     string
    Password string
}

type Options struct {
    Auth    *Auth      `prefix:"-"`
    URLs    []*url.URL `scheme:"^(http|https)$"`
    Verbose int        `min:"0" max:"8"`
}
```

In this case, the command line flag for the User field will be simply *-user*,
instead of *-auth-user*.

For even more control of naming, a struct can be passed to the Var function.
See the documentation of that function for details.

## Dependencies

Supports Go >= 1.10. go-config does not currently rely on any external packages.

## License

go-config uses the [MIT License](LICENSE.txt).

## Contributing

Pull requests and bug requests are welcome.

Please note that this project is released with a
[Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this
project you agree to abide by its terms.