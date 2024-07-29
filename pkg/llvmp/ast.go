package llvmp

import (
	"fmt"
	"strings"
)

func newModule() *Module {
	return &Module{
		Functions: map[string]*FnDef{},
	}
}

type Module struct {
	Functions map[string]*FnDef
}

func (m *Module) addFn(name string) *FnDef {
	fn := &FnDef{
		Name: name,
	}
	m.Functions[name] = fn
	return fn
}

func (m *Module) Dump() string {
	var b strings.Builder

	for fn, f := range m.Functions {
		b.WriteString(fmt.Sprintf("fn:%s %+v\n", fn, f))
		for _, step := range f.Steps {
			b.WriteString(fmt.Sprintf("  %+v\n", step))
		}
	}

	return b.String()
}

type FnKind string

const (
	FnKindTail     = FnKind("FnKindTail")
	FnKindInternal = FnKind("FnKindInternal")
)

type FnDef struct {
	Name    string
	Linkage string
	Kind    FnKind

	File string
	Line int

	Steps []*Step

	dbgRef int
}

func (d *FnDef) addStep() *Step {
	step := &Step{}
	d.Steps = append(d.Steps, step)
	return step
}

type StepKind string

const (
	StepFnCall   = StepKind("StepFnCall")
	StepTailCall = StepKind("StepTailCall")
	StepRet      = StepKind("StepRet")
)

type Step struct {
	Kind StepKind
	File string
	Line int

	Function string

	dbgRef int
	line   string
}
