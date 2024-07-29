package llvmp

import (
	"fmt"

	"github.com/bowei/cilium-bpf-hack/pkg/gviz"
	"github.com/bowei/cilium-bpf-hack/pkg/llvmp/ignore"
	"github.com/bowei/cilium-bpf-hack/pkg/llvmp/srcnote"
)

type rawCGData struct {
	node *gviz.Node
	fn   *FnDef
}

type RawCGParams struct {
	Start   string
	Ignored ignore.Set
	SrcAn   *srcnote.Set
}

var (
	condAttrib       = gviz.NewAt().Align("left").BGColor("yellow").Map()
	entryPointAttrib = gviz.NewAt().Align("left").BGColor("pink").Map()
	fnAttrib         = gviz.NewAt().Align("left").BGColor("green").Map()
	noteAttrib       = gviz.NewAt().Align("left").BGColor("lemonchiffon").Map()
	stepAttrib       = gviz.NewAt().Align("left").Map()
	tailCallAttrib   = gviz.NewAt().Align("left").BGColor("orange").Map()
)

func RawCG(m *Module, params *RawCGParams) string {
	g := gviz.NewGraph("r")
	fmt.Printf("// RawCG %s\n", params.Start)

	f2n := map[string]rawCGData{}

	// Create the function nodes.
	err := Closure(m, params.Start, func(m *Module, fn *FnDef) bool {
		fmt.Printf("// Function %q (%s:%d)\n", fn.Name, fn.File, fn.Line)

		if params.Ignored.Match(fn.Name) {
			fmt.Printf("// Node (skipped) %s\n", fn.Name)
			return true
		}

		fNode := g.NewNode(fn.Name)
		fNode.Attribs("shape", "rectangle")

		f2n[fn.Name] = rawCGData{
			node: fNode,
			fn:   fn,
		}

		if fn.Name == params.Start {
			fNode.AddRow([]gviz.NodeCol{
				{
					Text: "-",
					Port: "E0",
				},
				{},
				{
					Text:    "ENTRYPOINT",
					Attribs: entryPointAttrib,
				},
			})
		}

		fNode.AddRow([]gviz.NodeCol{
			{
				Text: fmt.Sprintf("%d", 0),
				Port: "Start0",
			},
			{
				Text: fmt.Sprintf("%s:%d", fn.File, fn.Line),
			},
			{
				Text:    fmt.Sprintf("%s()", fn.Name),
				Port:    "start",
				Attribs: fnAttrib,
			},
		})

		prevLine := fn.Line

		for i, step := range fn.Steps {
			// This code assumes the source file does not change inside of a
			// function. If the source file changes, the annotations will not work
			// correctly.
			if step.File != fn.File {
				fmt.Printf("// ERROR: source file mismatch: %q != %q\n", step.File, fn.File)
			}
			for _, an := range params.SrcAn.Lookup(fn.File, prevLine, step.Line) {
				switch an.Kind {
				case srcnote.KindConditional:
					fNode.AddRow([]gviz.NodeCol{
						{},
						{
							Text: fmt.Sprintf("%s:%d", an.FileName, an.Line),
						},
						{
							Text:    an.Text,
							Attribs: condAttrib,
						},
					})
				case srcnote.KindNote:
					fNode.AddRow([]gviz.NodeCol{
						{},
						{
							Text: fmt.Sprintf("%s:%d", an.FileName, an.Line),
						},
						{
							Text:    an.Text,
							Attribs: noteAttrib,
						},
					})
				default:
					fmt.Printf("// ERROR: unhandled source annotation: %v\n", an.Kind)
				}
			}
			prevLine = step.Line

			switch step.Kind {
			case StepFnCall:
				switch {
				case step.Function == "llvm":
					// These are llvm synthetic steps. Ignore.
					fmt.Printf("// Node: Step Fn LLVM %v\n", step)
				case step.Function == "tail_call_internal":
					// This is handled by the StepTailCall. Skip.
				case step.Function != "":
					fNode.AddRow([]gviz.NodeCol{
						{
							Text: fmt.Sprintf("%d", i),
						},
						{
							Text: fmt.Sprintf("%s:%d", step.File, step.Line),
						},
						{
							Text:    step.Function,
							Port:    fmt.Sprintf("s%d", i),
							Attribs: stepAttrib,
						},
					})
				default:
					fmt.Printf("// ERROR: Node: Step Fn (skipped) %v\n", step)
				}
			case StepTailCall:
				fNode.AddRow([]gviz.NodeCol{
					{
						Text: fmt.Sprintf("%d", i),
					},
					{
						Text: fmt.Sprintf("%s:%d", step.File, step.Line),
					},
					{
						Text:    step.Function,
						Port:    fmt.Sprintf("s%d", i),
						Attribs: tailCallAttrib,
					},
				})
			case StepRet:
				fNode.AddRow([]gviz.NodeCol{
					{
						Text: fmt.Sprintf("%d", i),
					},
					{
						Text: fmt.Sprintf("%s:%d", step.File, step.Line),
					},
					{
						Text:    "ret",
						Attribs: stepAttrib,
					},
				})
			default:
				fmt.Printf("// ERROR: Node: Step (skipped) %v\n", step)
			}
		}
		return true
	})

	if err != nil {
		fmt.Printf("// ERROR: %v\n", fmt.Errorf("RawCG: %w", err))
		// TODO: return code.
	}

	// Create edges
	for _, d := range f2n {
		if params.Ignored.Match(d.fn.Name) {
			fmt.Printf("// Edge: Node (skipped) %s\n", d.fn.Name)
			continue
		}
		for i, step := range d.fn.Steps {
			switch step.Kind {
			case StepFnCall:
				switch {
				case step.Function == "":
					fmt.Printf("// ERROR: Edge: Step (skipped) fname is empty: %v\n", step)
				case step.Function == "tail_call_internal":
					// This is handled by the StepTailCall. Skip.
				case !params.Ignored.Match(step.Function):
					targetD, ok := f2n[step.Function]
					if !ok {
						continue
					}
					e := g.NewEdge(d.node, targetD.node)
					e.APort = fmt.Sprintf("s%d", i)
					e.BPort = "Start0"
				default:
					// ignored
				}
			case StepTailCall:
				targetD, ok := f2n[step.Function]
				if !ok {
					continue
				}
				e := g.NewEdge(d.node, targetD.node)
				e.APort = fmt.Sprintf("s%d", i)
				e.BPort = "Start0"
			case StepRet:
				// Ret does not create a link.
			default:
				fmt.Printf("// ERROR: Edge: Step (skipped) %v\n", step)
			}
		}
	}
	return gviz.DotFile(g)
}

type runner struct{}

func closureCb(m *Module, fn *FnDef) bool {
	fmt.Printf("// Function %q (%s:%d)\n", fn.Name, fn.File, fn.Line)

	if params.Ignored.Match(fn.Name) {
		fmt.Printf("// Node (skipped) %s\n", fn.Name)
		return true
	}

	fNode := g.NewNode(fn.Name)
	fNode.Attribs("shape", "rectangle")

	f2n[fn.Name] = rawCGData{
		node: fNode,
		fn:   fn,
	}

	if fn.Name == params.Start {
		fNode.AddRow([]gviz.NodeCol{
			{
				Text: "-",
				Port: "E0",
			},
			{},
			{
				Text:    "ENTRYPOINT",
				Attribs: entryPointAttrib,
			},
		})
	}

	fNode.AddRow([]gviz.NodeCol{
		{
			Text: fmt.Sprintf("%d", 0),
			Port: "Start0",
		},
		{
			Text: fmt.Sprintf("%s:%d", fn.File, fn.Line),
		},
		{
			Text:    fmt.Sprintf("%s()", fn.Name),
			Port:    "start",
			Attribs: fnAttrib,
		},
	})

	prevLine := fn.Line

	for i, step := range fn.Steps {
		// This code assumes the source file does not change inside of a
		// function. If the source file changes, the annotations will not work
		// correctly.
		if step.File != fn.File {
			fmt.Printf("// ERROR: source file mismatch: %q != %q\n", step.File, fn.File)
		}
		for _, an := range params.SrcAn.Lookup(fn.File, prevLine, step.Line) {
			switch an.Kind {
			case srcnote.KindConditional:
				fNode.AddRow([]gviz.NodeCol{
					{},
					{
						Text: fmt.Sprintf("%s:%d", an.FileName, an.Line),
					},
					{
						Text:    an.Text,
						Attribs: condAttrib,
					},
				})
			case srcnote.KindNote:
				fNode.AddRow([]gviz.NodeCol{
					{},
					{
						Text: fmt.Sprintf("%s:%d", an.FileName, an.Line),
					},
					{
						Text:    an.Text,
						Attribs: noteAttrib,
					},
				})
			default:
				fmt.Printf("// ERROR: unhandled source annotation: %v\n", an.Kind)
			}
		}
		prevLine = step.Line

		switch step.Kind {
		case StepFnCall:
			switch {
			case step.Function == "llvm":
				// These are llvm synthetic steps. Ignore.
				fmt.Printf("// Node: Step Fn LLVM %v\n", step)
			case step.Function == "tail_call_internal":
				// This is handled by the StepTailCall. Skip.
			case step.Function != "":
				fNode.AddRow([]gviz.NodeCol{
					{
						Text: fmt.Sprintf("%d", i),
					},
					{
						Text: fmt.Sprintf("%s:%d", step.File, step.Line),
					},
					{
						Text:    step.Function,
						Port:    fmt.Sprintf("s%d", i),
						Attribs: stepAttrib,
					},
				})
			default:
				fmt.Printf("// ERROR: Node: Step Fn (skipped) %v\n", step)
			}
		case StepTailCall:
			fNode.AddRow([]gviz.NodeCol{
				{
					Text: fmt.Sprintf("%d", i),
				},
				{
					Text: fmt.Sprintf("%s:%d", step.File, step.Line),
				},
				{
					Text:    step.Function,
					Port:    fmt.Sprintf("s%d", i),
					Attribs: tailCallAttrib,
				},
			})
		case StepRet:
			fNode.AddRow([]gviz.NodeCol{
				{
					Text: fmt.Sprintf("%d", i),
				},
				{
					Text: fmt.Sprintf("%s:%d", step.File, step.Line),
				},
				{
					Text:    "ret",
					Attribs: stepAttrib,
				},
			})
		default:
			fmt.Printf("// ERROR: Node: Step (skipped) %v\n", step)
		}
	}
	return true
}
