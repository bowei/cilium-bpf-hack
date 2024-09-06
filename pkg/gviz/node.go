package gviz

import (
	"fmt"
	"strings"
)

type Node struct {
	Name   string
	Parent *Graph
	Tags   map[string]string

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
