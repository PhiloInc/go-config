package config

import (
	"testing"
)

type ninetyNineLoader struct {
	settings settings
}

func (nnl *ninetyNineLoader) Init(settings []Setting) {
	nnl.settings = settings
}

func (nnl *ninetyNineLoader) Load() error {
	var errs Errors
	for i := range nnl.settings {
		if err := nnl.settings[i].Set("99"); err != nil {
			errs.Append(err)
		}
	}
	return errs.AsError()
}

func (nnl *ninetyNineLoader) Name() string {
	return "99"
}

func (nnl *ninetyNineLoader) Usage() string {
	return ""
}

func TestConfig(t *testing.T) {
	t.Run("Var", testConfigVar)
	t.Run("Scan", testConfigScan)
	t.Run("from", testConfigFrom)
}

func testConfigVar(t *testing.T) {
	var c Config
	c.SetLoaders(Loaders{new(ninetyNineLoader)})
	var x int
	c.Var(&x, "", "x")
	if err := c.Load(); err != nil {
		t.Errorf("failed loading config: %s", err)
	}
	if x != 99 {
		t.Errorf("unexpected value %d after load", x)
	}
	t.Run("struct", testConfigVarStruct)
}

func testConfigVarStruct(t *testing.T) {
	var c Config
	loader := new(ninetyNineLoader)
	c.SetLoaders(Loaders{loader})
	x := struct {
		X int
	}{X: 1}
	c.Var(&x, "", "test")
	if err := c.Load(); err != nil {
		t.Errorf("failed loading config: %s", err)
	}
	if x.X != 99 {
		t.Errorf("unexpected value %d after load", x.X)
	}
	if loader.settings[0].Path.String() != "test->X" {
		t.Errorf("unexpected name %s for test->X", loader.settings[0].Path.String())
	}
}

func testConfigScan(t *testing.T) {
	var c Config
	c.SetLoaders(Loaders{new(ninetyNineLoader)})
	x := struct {
		X int
	}{}
	if err := c.Configure(&x); err != nil {
		t.Errorf("failed loading config: %s", err)
	}
	if x.X != 99 {
		t.Errorf("field X is %d", x.X)
	}
	t.Run("nested struct", testConfigScanNestedStruct)
	t.Run("nested pointer to struct", testConfigScanNestedPointerToStruct)
	t.Run("override name", testConfigScanOverrideName)
	t.Run("override prefix", testConfigScanOverridePrefix)
	t.Run("ommitted field", testConfigScanOmittedField)
	t.Run("unexported field", testConfigScanUnexportedField)
}

func testConfigScanNestedStruct(t *testing.T) {
	var c Config
	c.SetLoaders(Loaders{new(ninetyNineLoader)})
	x := struct {
		A struct {
			B struct {
				X int
			}
		}
	}{}
	if err := c.Configure(&x); err != nil {
		t.Errorf("failed loading config: %s", err)
	}
	if x.A.B.X != 99 {
		t.Errorf("field A.B.X is %d", x.A.B.X)
	}
}

func testConfigScanNestedPointerToStruct(t *testing.T) {
	var c Config
	c.SetLoaders(Loaders{new(ninetyNineLoader)})
	x := struct {
		A *struct {
			B **struct {
				X int
			}
		}
	}{}
	if err := c.Configure(&x); err != nil {
		t.Errorf("failed loading config: %s", err)
	}
	switch {
	case x.A == nil:
		t.Error("x.A is nil")
	case x.A.B == nil:
		t.Error("x.A.B is nil")
	case *x.A.B == nil:
		t.Error("*x.A.B is nil")
	case (*x.A.B).X != 99:
		t.Errorf("field (*x.A.B).X is %d", (*x.A.B).X)
	}
}

func testConfigScanOverrideName(t *testing.T) {
	var c Config
	loader := new(ninetyNineLoader)
	c.SetLoaders(Loaders{loader})
	x := struct {
		X int `config:"a"`
		Y struct {
			X int
		} `config:"b"`
	}{}
	if err := c.Configure(&x); err != nil {
		t.Errorf("failed loading config: %s", err)
	}
	switch {
	case len(loader.settings) != 2:
		t.Errorf("len(loader.settings) = %d", len(loader.settings))
	case loader.settings[0].Path.String() != "a":
		t.Errorf("path a = %s", loader.settings[0].Path)
	case loader.settings[1].Path.String() != "b->X":
		t.Errorf("path b = %s", loader.settings[1].Path)
	}
}

func testConfigScanOverridePrefix(t *testing.T) {
	var c Config
	loader := new(ninetyNineLoader)
	c.SetLoaders(Loaders{loader})
	x := struct {
		X int
		Y struct {
			Z int
		} `prefix:"-"`
	}{}
	if err := c.Configure(&x); err != nil {
		t.Errorf("failed loading config: %s", err)
	}
	switch {
	case len(loader.settings) != 2:
		t.Errorf("len(loader.settings) = %d", len(loader.settings))
	case loader.settings[0].Path.String() != "X":
		t.Errorf("path X = %s", loader.settings[0].Path)
	case loader.settings[1].Path.String() != "Z":
		t.Errorf("path Z = %s", loader.settings[1].Path)
	}
}

func testConfigScanOmittedField(t *testing.T) {
	var c Config
	loader := new(ninetyNineLoader)
	c.SetLoaders(Loaders{loader})
	x := struct {
		X int `config:"-"`
		Y struct {
			Z int
		} `config:"-"`
	}{}
	x.X, x.Y.Z = 1, 1
	if err := c.Configure(&x); err != nil {
		t.Errorf("failed loading config: %s", err)
	}
	if x.X != 1 {
		t.Errorf("ommitted field X set to %d", x.X)
	}
	if x.Y.Z != 1 {
		t.Errorf("field Z in ommited struct Y set to %d", x.Y.Z)
	}
}

func testConfigScanUnexportedField(t *testing.T) {
	var c Config
	loader := new(ninetyNineLoader)
	c.SetLoaders(Loaders{loader})
	x := struct {
		x int
		y struct {
			Z int
		}
	}{}
	x.x, x.y.Z = 1, 1
	if err := c.Configure(&x); err != nil {
		t.Errorf("failed loading config: %s", err)
	}
	if x.x != 1 {
		t.Errorf("ommitted field X set to %d", x.x)
	}
	if x.y.Z != 1 {
		t.Errorf("field Z in ommited struct Y set to %d", x.y.Z)
	}
}

func testConfigFrom(t *testing.T) {
	var c Config
	c.SetLoaders(Loaders{new(ninetyNineLoader)})
	x, y, z := 1, 2, 3
	c.Var(&x, `from:"*"`, "x")
	c.Var(&y, `from:"flag,env,99"`, "y")
	c.Var(&z, `from:"flag"`, "z")
	if err := c.Load(); err != nil {
		t.Errorf("failed loading config: %s", err)
	}
	if x != 99 {
		t.Error(`did not load value for from:"*"`)
	}
	if y != 99 {
		t.Error(`did not load value for from:"flag,env,99"`)
	}
	if z != 3 {
		t.Errorf(`did not skip value for from:"flag"`)
	}
}
