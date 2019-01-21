/*
Package config provides types and functions for managing application
configuration. A set of command-line flags or environment variables to parse
can be inferred from a struct's type. The values are then subject to validation
and assigned to the struct fields.

Quick Start

The easiest way to use the package is to define a struct type and pass it to
Configure. Struct tags can be used to control how fields are scanned or
validated.
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

If the program was given the -h flag, the following help would be written to
stderr:
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

Running the program with no command line flags and none of the environment
variables set would output:
	URLs: []
	Verbose: 1

Note that the entire Auth field is left as nil; nil pointers are only set if
a value is configured for the field. Running the program again with flags
	-urls http://www.google.com/ -urls http://www.yahoo.com/ -auth-user user1
provides the following:
	Auth: {User:user1 Password:}
	URLs: [http://www.google.com/ http://www.yahoo.com/]
	Verbose: 1

Finally, some of the values have validation tags for their struct fields. If
the validations are not met, an error is given:
	invalid value "10" for flag -verbose: Validating 10 failed: 10 is not less than or equal to 8

Supported Types

The following types are supported by the package:
 * bool
 * float32, float64
 * int, int8, int16, int32, int64
 * net.IP, net.IPNet
 * string
 * time.Duration, time.Time
 * uint, uint8, uint16, uint32, uint64
 * url.URL

In addition, any types derived from pointers and slices to those types are also
supported.

bool, float, int and uint values are parsed by strconv.ParseBool,
strconv.ParseFloat, strconv.ParseInt and strconv.ParseUint. Trying to set a
value that would overflow the type results in an error.

net.IP values are parsed using net.ParseIP. net.IPNet values are parsed using
net.ParseCIDR. url.URL values are parsed using url.Parse. time.Duration values
are parsed using time.ParseDuration.

time.Time values will be parsed using the following layouts until one is
succesful or all have been tried:
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

Additional types can be supported either by implementing the Setter API on that
type, or by implementing and registering a SetterCreator for that type. The
latter option will enable the package to automatically wrap derived pointer and
struct types.

Validation

Struct tags can be used to impose validation on parsed values.

`min:"x"` or `ge:"x"` sets a minimum value for int, float, uint, time.Duration
and time.Time values. The value must be appropriate for the type, and is
parsed the same as values of that type. For example
	Timeout time.Duration `min:"30s"`
indicates that the Timeout value must be greater than or equal to 30 seconds.

`max:"x"`` or `le:"x"` indicates a maximum value for the above types. Similarly,
`gt:"x"` or `lt:"x"` will require values greater than or less than the values
specified.

Values can be combined to set ranges. Note, however, that it's possible to set
impossible validations this way, e.g., `gt:"10" lt:"5"`.

`regexp:"x"` sets a regular expression for validating string values. A match
can occur anywhere in the string. To match the whole string, anchor the
expression with "^$". Also note that the value in the struct tag will be subject
to string unquoting; backslashes and double-quotes must be escaped with a
backslash.

`scheme:"x"`, `host:"x"` or `path:"x"` tags can be applied to url.URL values.
They specify regular expressions that are used to validate the corresponding
parts of the URL.

The `is:"x"` tag can be used to specify required or disallowed classes of IP
addresses for the net.IP or net.IPNet types. Valid values are:
 * global unicast
 * interface local multicast
 * link local multicast
 * link local unicast
 * loopback
 * multicast
 * unspecified

The class name can be prefixed with an exclamation point to indicate that the
class is disallowed. Multiple values can be combined by separating with commas.
For example, `is:"!loopback,!link local unicast"` would allow addresses that are
neither loopback nor link local unicast. If any disallowed class is matched,
the validation will fail. The allowed classes are matched in an "or" fashion,
however, and any one match will cause the validation to succeed.

The `net:"x"` tag can be used with net.IP or net.IPNet values to specify
required or disallowed networks. Networks are specified in address/bits form,
e.g., 169.254.0.0/16 or 2001:cdba:9abc:5678::/64. Disallowed networks are
specified by prefixing with an exclamation point. Multiple values can be
combined by separating with commas. The semantics are the same as for the 'is'
tag: A match for any disallowed value fails validation; a match for any single
allowed value causes validation to succeed.

The `version:"4"` or `version:"6"` tags can be used with net.IP and net.IPNet
values to specify that the value must be an IPv4 or IPv6 address.

Other Struct Tags

`config:"X"` can be used to override the name of a struct field, instead of
using the reflected name. The name will still be subjected to parsing with
SplitName. Setting `config:"-"` will cause the struct field to be skipped, i.e.,
it will not be scanned or configured.

`prefix:"X"` can be used to override how nested structs are handled. If no
prefix is set, the child struct uses the name of its field in the parent struct.
If a name is given, that name will be used to group settings in the child
struct. (`config:"X"` can be used for the same purpose.) If `prefix:"-"` is
given, the child struct will use the same prefix as its parent, and names parsed
from the child's fields will be added to the parent's level as if the fields had
been parsed from the parent. For example:
	type Options struct {
		Timeout time.Duration
		Auth struct {
			User string
			Pass string
		} `prefix:"-"`
	}
	var opts Options
	config.Configure(&opts)
will result in command line flags being generated for -timeout, -user and
-pass. Without the prefix tag, the names would have been -auth-user and
-auth-pass.

`from:"X"` will cause the value to only be configured by the named loaders. For
example
	Interactive bool `from:"flag"`
will only let the Interactive value be set from a command line flag.

`append:"false"` or `append:"true"` (the default) can be used with slice types.
When true, values are appended to existing values in a slice. When false,
setting a new value will cause the slice to be overwritten. (Setting multiple
values will still result in second and subsequent values being appended after
the first new value).

`sep:"X" can be used with slice types. It indicates a string on which the
value should be split on in order to populate a slice.
*/
package config
