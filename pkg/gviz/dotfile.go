package gviz

import (
	"fmt"
	"strings"
)

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
