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

	attribs map[string]string
	parent  *Graph
}

func (e *Edge) Attribs(attribPairs ...string) {
	if e.attribs == nil {
		e.attribs = map[string]string{}
	}
	if len(attribPairs)%2 != 0 {
		panic(fmt.Sprintf("non-paired attribPairs (must be even length): %v", attribPairs))
	}
	for i := 0; i < len(attribPairs); i += 2 {
		key := attribPairs[i]
		val := attribPairs[i+1]
		e.attribs[key] = val
	}
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

	ret := fmt.Sprintf("%s%s -> %s", nspace(indent), aRef, bRef)
	if len(e.attribs) > 0 {
		var attribs []string
		for k, v := range e.attribs {
			attribs = append(attribs, fmt.Sprintf(`%s="%s"`, k, v))
		}
		ret += "[" + strings.Join(attribs, ",") + "]"
	}
	return ret
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
