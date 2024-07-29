package gviz

import (
	"fmt"
	"testing"
)

func TestGraph(t *testing.T) {
	g := NewGraph("g1")

	n1 := g.NewNode("n1")
	n2 := g.NewNode("n2")

	for i := 0; i < 50; i++ {
		n2.AddRow([]NodeCol{
			{Port: "f0", Text: fmt.Sprintf("hello%d", i)},
			{Port: "f0", Text: fmt.Sprintf("foo%d", i)},
		})
	}
	e1 := g.NewEdge(n1, n2)

	t.Error(DotFile(g))
	fmt.Println(e1)
}
