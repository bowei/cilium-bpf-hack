package gviz

import (
	"fmt"
	"strings"
)

func NewGraph(name string) *Graph {
	g := &Graph{
		Name:   name,
		Nodes:  map[string]*Node{},
		Edges:  map[string]*Edge{},
		Graphs: map[string]*Graph{},
	}
	return g
}

type Graph struct {
	Prefix string

	Name string

	Nodes  map[string]*Node
	Edges  map[string]*Edge
	Graphs map[string]*Graph

	Parent *Graph

	indent int
}

func (g *Graph) NewNode(name string) *Node {
	n := &Node{
		Name:    name,
		Parent:  g,
		attribs: map[string]string{},
	}
	g.Nodes[name] = n
	return n
}

func (g *Graph) NewEdge(a, b *Node) *Edge {
	e := &Edge{A: a, B: b, parent: g}
	g.Edges[EdgeMapKey(e)] = e
	return e
}

func (g *Graph) NewGraph(name string) *Graph {
	sg := NewGraph(name)
	sg.Parent = g
	sg.indent = g.indent + 2
	g.Graphs[name] = sg
	return sg
}

// Find the node corresponding to "a.b.c" path. Returns nil if the path does not
// exist in the Graph.
func (g *Graph) Find(path string) *Node {
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

type Node struct {
	Name   string
	Parent *Graph

	rows    [][]NodeCol
	attribs map[string]string
}

type NodeCol struct {
	Port    string
	Text    string
	Attribs map[string]string
}

func (n *Node) FullName() string {
	name := "z_" + n.Name
	for g := n.Parent; g != nil; g = g.Parent {
		name = g.Name + "_" + name
	}
	return name
}

func (n *Node) AddRow(nc []NodeCol) {
	n.rows = append(n.rows, nc)
}

func (n *Node) Attribs(attribPairs ...string) {
	if len(attribPairs)%2 != 0 {
		panic("XXX")
	}
	for i := 0; i < len(attribPairs); i += 2 {
		key := attribPairs[i]
		val := attribPairs[i+1]
		n.attribs[key] = val
	}
}

func (n *Node) render(indent int) string {
	nodeAttribsStr := func() string {
		var attribs []string

		attribs = append(attribs, fmt.Sprintf(`label="%s"`, n.Name))
		for k, v := range n.attribs {
			attribs = append(attribs, fmt.Sprintf(`%s="%s"`, k, v))
		}
		return strings.Join(attribs, " ")
	}

	if len(n.rows) == 0 {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("// Node %s\n", n.FullName()))
		b.WriteString(fmt.Sprintf("%s%s [%s];", nspace(indent), n.FullName(), nodeAttribsStr()))
		return b.String()
	}

	renderCol := func(nc []NodeCol) string {
		var b strings.Builder
		for _, c := range nc {
			var attribs []string
			if c.Port != "" {
				attribs = append(attribs, fmt.Sprintf(`port="%s"`, c.Port))
			}
			for k, v := range c.Attribs {
				attribs = append(attribs, fmt.Sprintf(`%s="%s"`, k, v))
			}
			b.WriteString(fmt.Sprintf("<td %s>%s</td>", strings.Join(attribs, " "), c.Text))
		}
		return b.String()
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s%s [%s ", nspace(indent), n.FullName(), nodeAttribsStr()))
	b.WriteString(`label=<<table border="0">`)
	for _, r := range n.rows {
		b.WriteString("<tr>")
		b.WriteString(renderCol(r))
		b.WriteString("</tr>")
	}
	b.WriteString("</table>>")
	b.WriteString("];")

	return b.String()
}

type Edge struct {
	A     *Node
	APort string
	B     *Node
	BPort string

	parent *Graph
}

func (e *Edge) render(indent int) string {
	aRef := e.A.FullName()
	bRef := e.B.FullName()
	if e.APort != "" {
		aRef += ":" + e.APort
	}
	if e.BPort != "" {
		bRef += ":" + e.BPort
	}
	return fmt.Sprintf("%s%s -> %s;", nspace(indent), aRef, bRef)
}

func EdgeMapKey(e *Edge) string {
	var parts []string

	if e.A.Name < e.B.Name {
		parts = []string{
			e.A.Name, e.APort,
			e.B.Name, e.BPort,
		}
	} else {
		parts = []string{
			e.B.Name, e.BPort,
			e.A.Name, e.APort,
		}
	}
	return strings.Join(parts, ":")
}

func DotFile(g *Graph) string {
	var b strings.Builder

	b.WriteString("digraph {\n")
	b.WriteString("rankdir=\"LR\"\n")
	b.WriteString(dotFileGraph(g))

	for _, sg := range g.Graphs {
		b.WriteString(fmt.Sprintf("%ssubgraph cluster_%s {\n", sg.indentStr(), sg.Name))
		b.WriteString(dotFileGraph(sg))
		b.WriteString("}\n")
	}

	b.WriteString("}\n")

	return b.String()
}

func dotFileGraph(g *Graph) string {
	var b strings.Builder

	for _, n := range g.Nodes {
		b.WriteString(n.render(g.indent))
		b.WriteString("\n")
	}
	for _, e := range g.Edges {
		b.WriteString(e.render(g.indent))
		b.WriteString("\n")
	}

	return b.String()
}

func nspace(n int) string {
	var ret string
	for i := 0; i < n; i++ {
		ret += " "
	}
	return ret
}
