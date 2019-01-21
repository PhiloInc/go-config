package config

import (
	"fmt"
	"strings"
)

/*
Path provides a leaf node in a hierarchy of configuration parameters.

When scanning nested structs, for example, there will be a hierarchy such as
 TypeNameA +> TypeNameB -> TypeNameC
           |
           +-> TypeNameD -> TypeNameE
A Path can refer to the edge nodes, TypeNameC or TypeNameE in this case. Each
will have an association with one parent node.
*/
type Path struct {
	name   string
	parent *NodePath
}

// Name returns the name of node p.
func (p *Path) Name() string {
	if p == nil {
		return ""
	}
	return p.name
}

/*
Elements returns a list of names from the root node to p.

If the root *NodePath is unnamed, it will not be included in the list.
Otherwise, the first item will be the name of the root node, and the last will
be the name of p. E.g., for
	NodePath()->NodePath(usr)->NodePath(bin)->Path(python)
Elements returns
	[]string{"usr", "bin", "python"}
*/
func (p *Path) Elements() []string {
	if p == nil {
		return nil
	}
	path, i := p, 0
	for ; path.parent != nil; path = &path.parent.Path {
		i++
	}
	// if root path has a name, include it.
	if path.name != "" {
		i++
	}
	elems := make([]string, i)
	for path, i = p, i-1; path.parent != nil; path, i = &path.parent.Path, i-1 {
		elems[i] = path.name
	}
	if i == 0 {
		elems[0] = path.name
	}
	return elems
}

// String returns a string for p by joining the result of p.Elements() with "->".
func (p *Path) String() string {
	return strings.Join(p.Elements(), "->")
}

/*
PathCmp compares two Paths to determine ordering.

It returns -1 if a should order before b; 1 if a should order after b; or 0 if
a and b are equivalent.
*/
func PathCmp(a, b *Path) int {
	ei, ej := a.Elements(), b.Elements()
	for len(ei) != 0 && len(ej) != 0 {
		if ei[0] < ej[0] {
			return -1
		}
		if ei[0] > ej[0] {
			return 1
		}
		ei, ej = ei[1:], ej[1:]
	}
	if len(ei) < len(ej) {
		return -1
	}
	if len(ei) > len(ej) {
		return 1
	}
	return 0
}

type child struct {
	*NodePath
	*Path
}

/*
NodePath represents an intermediate node in a configuration path.

It is associated with one parent NodePath, or none if it is a root. It is
also associated with 0 or more child NodePath or Path elements.
*/
type NodePath struct {
	Path
	children map[string]child
}

/*
NewRootPath returns a *NodePath that represents the root of a path hierarchy.

If name is the empty string, the root is unnamed. This is equivalent to the
zero value, which is also ready to use.
*/
func NewRootPath(name string) *NodePath {
	return &NodePath{Path: Path{name: name}}
}

/*
NewPath returns a new Path but without adding it as a child of np.

The returned *Path will have np as its parent, but will not be locatable using
np.FindPath.
*/
func (np *NodePath) NewPath(name string) *Path {
	return &Path{name: name, parent: np}
}

/*
NewNodePath returns a new NodePath but without adding it as a child of np.

The returned *NodePath will have np as its parent, but will not be locatable
using np.FindNodePath.
*/
func (np *NodePath) NewNodePath(name string) *NodePath {
	return &NodePath{Path: Path{name: name, parent: np}}
}

/*
AddPath registers a child Path to np.

This makes the child locatable using np.FindPath.
*/
func (np *NodePath) AddPath(c *Path) *Path {
	c.parent = np
	if np.children == nil {
		np.children = make(map[string]child, 1)
	}
	item := np.children[c.name]
	if item.Path != nil {
		panic(fmt.Sprintf("path %s already exists", c))
	} else {
		item.Path = c
	}
	np.children[c.name] = item
	return c
}

/*
AddNodePath registers a child NodePath to np.

This makes the child locatable using np.FindNodePath.
*/
func (np *NodePath) AddNodePath(c *NodePath) *NodePath {
	c.parent = np
	if np.children == nil {
		np.children = make(map[string]child, 1)
	}
	item := np.children[c.name]
	if item.NodePath != nil {
		panic(fmt.Sprintf("node path %s already exists", c))
	} else {
		item.NodePath = c
	}
	np.children[c.name] = item
	return c
}

func (np *NodePath) find(elements ...string) child {
	var item child
	for i := range elements {
		if np == nil { // shouldn't happen
			return child{}
		}
		var ok bool
		if item, ok = np.children[elements[i]]; !ok {
			return child{}
		}
		np = item.NodePath
	}
	return item
}

/*
FindPath finds the given path.

The name of np is not considered, so the first parameter will be a child of
np. For example, in the hierarchy
	NodePath()->NodePath(usr)->NodePath(bin)->Path(python)
Calling FindPath("usr", "bin", "python") on the root node will return
Path(python).

If the path is not found within np, or is not a leaf node (*Path), FindPath
returns nil.
*/
func (np *NodePath) FindPath(elements ...string) *Path {
	return np.find(elements...).Path
}

/*
FindNodePath finds the given path.

The name of np is not considered, so the first parameter will be a child of
np. For example, in the hierarchy
	NodePath()->NodePath(usr)->NodePath(bin)->Path(python)
Calling FindPath("usr", "bin") on the root node will return NodePath(bin).

If the path is not found within np, or is not an intermediate node (*NodePath),
FindNodePath returns nil.
*/
func (np *NodePath) FindNodePath(elements ...string) *NodePath {
	return np.find(elements...).NodePath
}
