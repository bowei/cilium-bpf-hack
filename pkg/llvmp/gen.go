package llvmp

import (
	"fmt"
	"strings"
)

func Graphviz(m *Module) string {
	var b strings.Builder

	b.WriteString("digraph {\n")

	var idx int

	for _, fn := range m.Functions {
		prev := fn.Name
		for _, st := range fn.Steps {
			if st.Function == "llvm" {
				// Ignore LLVM intrinsics.
				continue
			}
			b.WriteString(fmt.Sprintf("  %s -> %s_%d\n", prev, st.Function, idx))
			prev = fmt.Sprintf("%s_%d", st.Function, idx)
			idx++
		}
		b.WriteString(fmt.Sprintf("  %s -> end_%d\n", prev, idx))
	}

	b.WriteString("}\n")

	return b.String()
}
