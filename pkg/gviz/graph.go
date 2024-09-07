package gviz

import (
	"strings"
)

func NewGraph(name string) *Graph {
	g := &Graph{
		Name:   name,
		Nodes:  map[string]*Node{},
		Edges:  map[string]*Edge{},
		Graphs: map[string]*Graph{},
		Tags:   map[string]string{},
	}
	return g
}

type Graph struct {
	Prefix string

	Name   string
	Parent *Graph

	Nodes  map[string]*Node
	Edges  map[string]*Edge
	Graphs map[string]*Graph
	Tags   map[string]string

	indent int
}

func (g *Graph) NewGraph(name string) *Graph {
	sg := NewGraph(name)
	sg.Parent = g
	sg.indent = g.indent + 2
	g.Graphs[name] = sg

	return sg
}

func (g *Graph) NewNode(name string) *Node {
	n := &Node{
		Name:    name,
		Parent:  g,
		Tags:    map[string]string{},
		attribs: map[string]string{},
	}
	g.Nodes[name] = n

	return n
}

func (g *Graph) NewEdge(a, b *Node) *Edge {
	e := &Edge{
		A:      a,
		B:      b,
		Tags:   map[string]string{},
		parent: g,
	}
	g.Edges[EdgeMapKey(e)] = e

	return e
}

// FindNode corresponding to "a.b.c" path. Returns nil if the path does not
// exist in the Graph.
func (g *Graph) FindNode(path string) *Node {
	parts := strings.Split(path, ".")
	curGraph := g
	for i, p := range parts {
		if gr, ok := curGraph.Graphs[p]; ok {
			curGraph = gr
			continue
		}
		if n, ok := curGraph.Nodes[p]; ok {
			if i != len(parts)-1 {
				return nil
			}
			return n
		}
	}
	return nil
}

func (g *Graph) indentStr() string {
	var ret string
	for i := 0; i < g.indent; i++ {
		ret += " "
	}
	return ret
}
