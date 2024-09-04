package rawcg

import (
	"fmt"

	"github.com/bowei/cilium-bpf-hack/pkg/gviz"
	"github.com/bowei/cilium-bpf-hack/pkg/llvmp"
	"github.com/bowei/cilium-bpf-hack/pkg/llvmp/ignore"
	"github.com/bowei/cilium-bpf-hack/pkg/llvmp/srcnote"
)

type Params struct {
	Start   string
	Ignored ignore.Set
	SrcAn   *srcnote.Set
}

func Run(m *llvmp.Module, params *Params) (string, error) {
	r := runner{
		m:      m,
		params: params,
		g:      gviz.NewGraph("cfg"),
		f2n:    map[string]rawCGData{},
	}
	return r.do()
}

var (
	condAttrib       = gviz.NewAt().Align("left").BGColor("yellow").Map()
	entryPointAttrib = gviz.NewAt().Align("left").BGColor("pink").Map()
	fnAttrib         = gviz.NewAt().Align("left").BGColor("green").Map()
	noteAttrib       = gviz.NewAt().Align("left").BGColor("lemonchiffon").Map()
	stepAttrib       = gviz.NewAt().Align("left").Map()
	tailCallAttrib   = gviz.NewAt().Align("left").BGColor("orange").Map()
)

type rawCGData struct {
	node *gviz.Node
	fn   *llvmp.FnDef
}

type runner struct {
	m      *llvmp.Module
	params *Params
	g      *gviz.Graph
	f2n    map[string]rawCGData
}

func (r *runner) do() (string, error) {
	fmt.Printf("// RawCG %s\n", r.params.Start)

	// Create the function nodes.
	err := llvmp.Closure(r.m, r.params.Start, r.createNode, llvmp.ClosureOptions{
		FollowTailCalls: true,
	})
	if err != nil {
		fmt.Printf("// ERROR: %v\n", fmt.Errorf("RawCG: %w", err))
		// TODO: return code.
	}

	r.createEdges() // TODO: error

	return gviz.DotFile(r.g), nil
}

func (r *runner) createNode(_ *llvmp.Module, fn *llvmp.FnDef) bool {
	fmt.Printf("// Function %q (%s:%d)\n", fn.Name, fn.File, fn.Line)

	if r.params.Ignored.Match(fn.Name) {
		fmt.Printf("// Node (skipped) %s\n", fn.Name)
		return true
	}

	fNode := r.g.NewNode(fn.Name)
	fNode.Attribs("shape", "rectangle")

	r.f2n[fn.Name] = rawCGData{
		node: fNode,
		fn:   fn,
	}

	if fn.Name == r.params.Start {
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
		for _, an := range r.params.SrcAn.Lookup(fn.File, prevLine, step.Line) {
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
		case llvmp.StepFnCall:
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
		case llvmp.StepTailCall:
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
		case llvmp.StepRet:
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

func (r *runner) createEdges() {
	for _, d := range r.f2n {
		if r.params.Ignored.Match(d.fn.Name) {
			fmt.Printf("// Edge: Node (skipped) %s\n", d.fn.Name)
			continue
		}
		for i, step := range d.fn.Steps {
			switch step.Kind {
			case llvmp.StepFnCall:
				switch {
				case step.Function == "":
					fmt.Printf("// ERROR: Edge: Step (skipped) fname is empty: %v\n", step)
				case step.Function == "tail_call_internal":
					// This is handled by the StepTailCall. Skip.
				case !r.params.Ignored.Match(step.Function):
					targetD, ok := r.f2n[step.Function]
					if !ok {
						continue
					}
					e := r.g.NewEdge(d.node, targetD.node)
					e.APort = fmt.Sprintf("s%d", i)
					e.BPort = "Start0"
				default:
					// ignored
				}
			case llvmp.StepTailCall:
				targetD, ok := r.f2n[step.Function]
				if !ok {
					continue
				}
				e := r.g.NewEdge(d.node, targetD.node)
				e.APort = fmt.Sprintf("s%d", i)
				e.BPort = "Start0"
			case llvmp.StepRet:
				// Ret does not create a link.
			default:
				fmt.Printf("// ERROR: Edge: Step (skipped) %v\n", step)
			}
		}
	}
}
