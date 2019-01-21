package config

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

var ipType = reflect.TypeOf(net.IP{})
var ipNetType = reflect.TypeOf(net.IPNet{})

func ipVersion(ip net.IP, version string) error {
	if version != "" {
		switch version {
		case "4":
			if ip.To4() == nil {
				msg := fmt.Sprintf("%s is not an IPv4 address", ip)
				return &ValidationError{Value: ip, Message: msg}
			}
		case "6":
			if ip.To4() != nil {
				msg := fmt.Sprintf("%s is not an IPv6 address", ip)
				return &ValidationError{Value: ip, Message: msg}
			}
		default:
			msg := fmt.Sprintf("invalid IP address version %s", version)
			return &ValidationError{Value: ip, Message: msg}
		}
	}
	return nil
}

var isFnMap = map[string]func(net.IP) bool{
	"global unicast":            func(ip net.IP) bool { return ip.IsGlobalUnicast() },
	"interface local multicast": func(ip net.IP) bool { return ip.IsInterfaceLocalMulticast() },
	"link local multicast":      func(ip net.IP) bool { return ip.IsLinkLocalMulticast() },
	"link local unicast":        func(ip net.IP) bool { return ip.IsLinkLocalUnicast() },
	"loopback":                  func(ip net.IP) bool { return ip.IsLoopback() },
	"multicast":                 func(ip net.IP) bool { return ip.IsMulticast() },
	"unspecified":               func(ip net.IP) bool { return ip.IsUnspecified() },
}

func ipIs(ip net.IP, spec string) error {
	var oneof, prohibited []string
	for _, s := range strings.Split(spec, ",") {
		if s = strings.TrimSpace(s); s == "" {
			continue
		}
		x := &oneof
		if s[0] == '!' {
			if s = s[1:]; s == "" {
				continue
			}
			x = &prohibited
		}
		if _, ok := isFnMap[s]; !ok {
			msg := fmt.Sprintf("unrecognized address class %s", s)
			return &ValidationError{Value: ip, Message: msg}
		}
		*x = append(*x, s)
	}
	for _, s := range prohibited {
		if isFnMap[s](ip) {
			msg := fmt.Sprintf("%s can not be %s", ip, s)
			return &ValidationError{Value: ip, Message: msg}
		}
	}
	if len(oneof) == 0 {
		return nil
	}
	for _, s := range oneof {
		if isFnMap[s](ip) {
			return nil
		}
	}
	msg := fmt.Sprintf(
		"%s did not match any allowed class (%s)",
		ip, strings.Join(oneof, ","),
	)
	return &ValidationError{Value: ip, Message: msg}
}

func ipNet(ip net.IP, spec string) error {
	var oneof, prohibited []*net.IPNet
	for _, s := range strings.Split(spec, ",") {
		if s = strings.TrimSpace(s); s == "" {
			continue
		}
		x := &oneof
		if s[0] == '!' {
			if s = s[1:]; s == "" {
				continue
			}
			x = &prohibited
		}
		_, n, err := net.ParseCIDR(s)
		if err != nil {
			msg := fmt.Sprintf("invalid IP network %s: %s", s, err)
			return &ValidationError{Value: ip, Message: msg}
		}
		*x = append(*x, n)
	}
	for _, n := range prohibited {
		if n.Contains(ip) {
			msg := fmt.Sprintf("%s can not be in network %s", ip, n)
			return &ValidationError{Value: ip, Message: msg}
		}
	}
	if len(oneof) == 0 {
		return nil
	}
	for _, n := range oneof {
		if n.Contains(ip) {
			return nil
		}
	}
	oneofStr := make([]string, len(oneof))
	for i := range oneof {
		oneofStr[i] = oneof[i].String()
	}
	msg := fmt.Sprintf(
		"%s did not match any allowed network (%s)",
		ip, strings.Join(oneofStr, ","),
	)
	return &ValidationError{Value: ip, Message: msg}
}

type ipSetter struct {
	val *net.IP
	tag reflect.StructTag
}

func (is *ipSetter) String() string {
	if is.val == nil {
		return net.IP{}.String()
	}
	return is.val.String()
}

func (is *ipSetter) Set(val string) error {
	ipval := net.ParseIP(val)
	if ipval == nil {
		return &ConversionError{Value: val, ToType: ipType}
	}

	if err := ipVersion(ipval, is.tag.Get("version")); err != nil {
		return err
	}
	if err := ipIs(ipval, is.tag.Get("is")); err != nil {
		return err
	}
	if err := ipNet(ipval, is.tag.Get("net")); err != nil {
		return err
	}

	*is.val = ipval
	return nil
}

func (*ipSetter) SetInt(val int64) error {
	return &ConversionError{Value: val, ToType: ipType}
}

func (*ipSetter) SetUint(val uint64) error {
	return &ConversionError{Value: val, ToType: ipType}
}

func (*ipSetter) SetFloat(val float64) error {
	return &ConversionError{Value: val, ToType: ipType}
}

func (*ipSetter) SetBool(val bool) error {
	return &ConversionError{Value: val, ToType: ipType}
}

func (is *ipSetter) Get() interface{} {
	if is.val == nil {
		return net.IP{}
	}
	return *is.val
}

type ipSetterCreator struct{}

func (ipSetterCreator) Type() reflect.Type {
	return ipType
}

func (isc ipSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	return &ipSetter{val: val.Addr().Interface().(*net.IP), tag: tag}
}

type ipNetSetter struct {
	val *net.IPNet
	tag reflect.StructTag
}

func (ins *ipNetSetter) String() string {
	if ins.val == nil {
		return (&net.IPNet{}).String()
	}
	return ins.val.String()
}

func (ins *ipNetSetter) Set(val string) error {
	_, nval, err := net.ParseCIDR(val)
	if err != nil {
		return &ConversionError{Value: val, ToType: ipNetType}
	}

	if err := ipVersion(nval.IP, ins.tag.Get("version")); err != nil {
		if err, ok := err.(*ValidationError); ok {
			err.Value = nval
		}
		return err
	}
	if err := ipIs(nval.IP, ins.tag.Get("is")); err != nil {
		if err, ok := err.(*ValidationError); ok {
			err.Value = nval
		}
		return err
	}
	if err := ipNet(nval.IP, ins.tag.Get("net")); err != nil {
		if err, ok := err.(*ValidationError); ok {
			err.Value = nval
		}
		return err
	}

	*ins.val = *nval
	return nil
}

func (*ipNetSetter) SetInt(val int64) error {
	return &ConversionError{Value: val, ToType: ipType}
}

func (*ipNetSetter) SetUint(val uint64) error {
	return &ConversionError{Value: val, ToType: ipType}
}

func (*ipNetSetter) SetFloat(val float64) error {
	return &ConversionError{Value: val, ToType: ipType}
}

func (*ipNetSetter) SetBool(val bool) error {
	return &ConversionError{Value: val, ToType: ipType}
}

func (ins *ipNetSetter) Get() interface{} {
	if ins.val == nil {
		return net.IPNet{}
	}
	return *ins.val
}

type ipNetSetterCreator struct{}

func (ipNetSetterCreator) Type() reflect.Type {
	return ipNetType
}

func (isc ipNetSetterCreator) Setter(val reflect.Value, tag reflect.StructTag) Setter {
	return &ipNetSetter{val: val.Addr().Interface().(*net.IPNet), tag: tag}
}

func init() {
	DefaultSetterRegistry.Add(ipSetterCreator{})
	DefaultSetterRegistry.Add(ipNetSetterCreator{})
}
