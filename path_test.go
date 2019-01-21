package config

import (
	"testing"
)

func TestPath(t *testing.T) {
	t.Run("Elements", func(t *testing.T) {
		root := NewRootPath("")
		node := root.AddNodePath(root.NewNodePath("a"))
		node = node.AddNodePath(node.NewNodePath("b"))
		p := node.AddPath(node.NewPath("c"))
		elements := p.Elements()
		switch {
		case len(elements) != 3,
			elements[0] != "a",
			elements[1] != "b",
			elements[2] != "c":
			t.Errorf("path with anonymous root returned elements %v", elements)
		}
		root = NewRootPath("a")
		node = root.AddNodePath(root.NewNodePath("b"))
		node = node.AddNodePath(node.NewNodePath("c"))
		p = node.AddPath(node.NewPath("d"))
		elements = p.Elements()
		switch {
		case len(elements) != 4,
			elements[0] != "a",
			elements[1] != "b",
			elements[2] != "c",
			elements[3] != "d":
			t.Errorf("path with named root returned elements %v", elements)
		}
	})
	t.Run("Cmp", func(t *testing.T) {
		root := NewRootPath("")
		a := root.AddNodePath(root.NewNodePath("a"))
		b := a.AddNodePath(a.NewNodePath("b"))
		c := b.AddPath(b.NewPath("c"))
		d := b.AddPath(b.NewPath("d"))
		if i := PathCmp(c, d); i != -1 {
			t.Errorf("comparing %s with %s returned %d", c, d, i)
		}
		if i := PathCmp(d, c); i != 1 {
			t.Errorf("comparing %s with %s returned %d", d, c, i)
		}
		if i := PathCmp(d, d); i != 0 {
			t.Errorf("comparing %s with %s returned %d", d, d, i)
		}

		aa := a.AddPath(a.NewPath("a"))
		if i := PathCmp(aa, d); i != -1 {
			t.Errorf("comparing %s with %s returned %d", aa, d, i)
		}
		if i := PathCmp(d, aa); i != 1 {
			t.Errorf("comparing %s with %s returned %d", d, aa, i)
		}
	})
}

func TestNodePath(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		root := NewRootPath("")
		elements := root.Elements()
		if len(elements) != 0 {
			t.Errorf("anonymous root path returned elements %v", elements)
		}
	})
	t.Run("NewPath", func(t *testing.T) {
		root := NewRootPath("root")
		p := root.NewPath("test")
		elements := p.Elements()
		if len(elements) != 2 || elements[0] != "root" || elements[1] != "test" {
			t.Errorf("path root->test returned elements %v", elements)
		}
		if found := root.FindPath("test"); found != nil {
			t.Errorf("finding unadded path returned %s", found)
		}
		if found := root.FindNodePath("test"); found != nil {
			t.Errorf("finding unadded node path returned %s", found)
		}
	})
	t.Run("AddPath", func(t *testing.T) {
		root := NewRootPath("root")
		p := root.AddPath(root.NewPath("test"))
		elements := p.Elements()
		if len(elements) != 2 || elements[0] != "root" || elements[1] != "test" {
			t.Errorf("path root->test returned elements %v", elements)
		}
		if found := root.FindPath("test"); found == nil {
			t.Errorf("finding added path returned nil")
		} else if found != p {
			t.Errorf("finding added path returned %s", found)
		}
		if found := root.FindNodePath("test"); found != nil {
			t.Errorf("finding unadded node path returned %s", found)
		}
	})
	t.Run("NewNodePath", func(t *testing.T) {
		root := NewRootPath("root")
		p := root.NewNodePath("test")
		elements := p.Elements()
		if len(elements) != 2 || elements[0] != "root" || elements[1] != "test" {
			t.Errorf("path root->test returned elements %v", elements)
		}
		if found := root.FindNodePath("test"); found != nil {
			t.Errorf("finding unadded node path returned %s", found)
		}
		if found := root.FindPath("test"); found != nil {
			t.Errorf("finding unadded path returned %s", found)
		}
	})
	t.Run("AddNodePath", func(t *testing.T) {
		root := NewRootPath("root")
		p := root.AddNodePath(root.NewNodePath("test"))
		elements := p.Elements()
		if len(elements) != 2 || elements[0] != "root" || elements[1] != "test" {
			t.Errorf("path root->test returned elements %v", elements)
		}
		if found := root.FindNodePath("test"); found == nil {
			t.Errorf("finding added node path returned nil")
		} else if found != p {
			t.Errorf("finding added node path returned %s", found)
		}
		if found := root.FindPath("test"); found != nil {
			t.Errorf("finding unadded path returned %s", found)
		}
	})
	t.Run("FindPath", func(t *testing.T) {
		root := NewRootPath("")
		node := root.AddNodePath(root.NewNodePath("a"))
		node = node.AddNodePath(node.NewNodePath("b"))
		p := node.AddPath(node.NewPath("c"))
		if found := root.FindPath("a", "b", "c"); found == nil {
			t.Errorf("finding added path returned nil")
		} else if found != p {
			t.Errorf("finding added path returned %s", found)
		}
	})
	t.Run("FindNodePath", func(t *testing.T) {
		root := NewRootPath("root")
		node := root.AddNodePath(root.NewNodePath("a"))
		node = node.AddNodePath(node.NewNodePath("b"))
		node = node.AddNodePath(node.NewNodePath("c"))
		if found := root.FindNodePath("a", "b", "c"); found == nil {
			t.Errorf("finding added node path returned nil")
		} else if found != node {
			t.Errorf("finding added node path returned %s", found)
		}
	})
}
