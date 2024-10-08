// Package rawcg generates a raw call graph starting from a given function.
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

	err := llvmp.Closure(r.m, r.params.Start, r.createNode, llvmp.ClosureOptions{})
	if err != nil {
		fmt.Printf("// ERROR: %v\n", fmt.Errorf("RawCG: %w", err))
		// TODO: return code.
	}

	r.createEdges() // TODO: error
	r.hideUnreachable()

	return gviz.DotFile(r.g), nil
}

func (r *runner) createNode(_ *llvmp.Module, fn *llvmp.FnDef) bool {
	fmt.Printf("// Function %q (%s:%d)\n", fn.Name, fn.File, fn.Line)

	fNode := r.g.NewNode(fn.Name)
	fNode.Attribs("shape", "rectangle")

	if r.params.Ignored.Match(fn.Name) {
		fNode.Hidden = true
	}

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

		r.addAnnotations(fn.File, prevLine, step.Line, fNode)
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

func (r *runner) addAnnotations(fileName string, start, end int, node *gviz.Node) {
	for _, an := range r.params.SrcAn.Lookup(fileName, start, end) {
		for k, v := range an.Tags {
			// TODO: handle overwrite or multiple
			node.Tags[k] = v
		}
		switch an.Kind {
		case srcnote.KindConditional:
			node.AddRow([]gviz.NodeCol{
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
			node.AddRow([]gviz.NodeCol{
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
}

func (r *runner) createEdges() {
	for _, d := range r.f2n {
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
				e.Attribs("color", "orange")
			case llvmp.StepRet:
				// Ret does not create a link.
			default:
				fmt.Printf("// ERROR: Edge: Step (skipped) %v\n", step)
			}
		}
	}
}

func (r *runner) hideUnreachable() {
	start := r.f2n[r.params.Start]

	visible := map[*gviz.Node]bool{}
	gviz.Traverse(
		start.node,
		func(n *gviz.Node) bool {
			if !n.Hidden {
				visible[n] = true
				return true
			}
			return false
		},
		func(e *gviz.Edge) bool { return !e.B.Hidden })
	for _, n := range r.g.Nodes {
		if !visible[n] {
			n.Hidden = true
		}
	}
}
