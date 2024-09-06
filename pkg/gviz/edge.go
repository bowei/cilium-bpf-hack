package gviz

import (
	"fmt"
	"strings"
)

type Edge struct {
	A     *Node
	APort string
	B     *Node
	BPort string
	Tags  map[string]string

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
