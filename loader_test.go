package config

import (
	"strings"
	"testing"
	"time"
)

func TestSplitName(t *testing.T) {
	slicesEqual := func(a, b []string) bool {
		for len(a) != 0 && len(b) != 0 {
			if a[0] != b[0] {
				return false
			}
			a, b = a[1:], b[1:]
		}
		return len(a) == len(b)
	}
	sliceStr := func(a []string) string {
		return strings.Join(a, ", ")
	}
	type value struct {
		in       string
		expected []string
	}
	values := []value{
		{"HTMLEntityID", []string{"HTML", "Entity", "ID"}},
		{"HTMLEntityID2", []string{"HTML", "Entity", "ID2"}},
		{"UUID", []string{"UUID"}},
		{"UUIDv2", []string{"UUID", "v2"}},
		{"userName", []string{"user", "Name"}},
		{"userNameD", []string{"user", "NameD"}},
		{"UserName", []string{"User", "Name"}},
		{"IPv6Net", []string{"IP", "v6", "Net"}},
		{"URLs", []string{"URLs"}},
	}
	for i := range values {
		out := SplitName(values[i].in)
		if !slicesEqual(out, values[i].expected) {
			t.Errorf("splitting name %s returned %s", values[i].in, sliceStr(out))
		}
	}
}

func TestFriendlyTypeName(t *testing.T) {
	type value struct {
		in       interface{}
		expected string
	}
	values := []value{
		{"", "string"},
		{time.Time{}, "time"},
		{new(**int), "int"},
		{make([][]*time.Duration, 1), "duration"},
		{struct{}{}, "value"},
	}
	for i := range values {
		out := FriendlyTypeName(values[i].in)
		if out != values[i].expected {
			t.Errorf("making friendly name for type %T returned %s", values[i].in, out)
		}
	}
}
