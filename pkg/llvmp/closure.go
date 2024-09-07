package llvmp

import "fmt"

func newClosureQueue() *closureQueue {
	return &closureQueue{
		done: map[string]bool{},
	}
}

type closureQueue struct {
	q    []*FnDef
	done map[string]bool
}

func (q *closureQueue) maybePush(fn *FnDef) {
	if _, ok := q.done[fn.Name]; ok {
		return
	}
	q.done[fn.Name] = true
	q.q = append(q.q, fn)
}

func (q *closureQueue) pop() *FnDef {
	if q.empty() {
		return nil
	}
	ret := q.q[0]
	q.q = q.q[1:]
	return ret
}

func (q *closureQueue) empty() bool { return len(q.q) == 0 }

type ClosureOptions struct {
	IgnoreEdge func(*Module, *FnDef, *Step) bool
}

func Closure(m *Module, startFn string, fnCallback func(m *Module, fn *FnDef) bool, options ClosureOptions) error {
	fn, ok := m.Functions[startFn]
	if !ok {
		return fmt.Errorf("startFn not found: %q", startFn)
	}

	q := newClosureQueue()
	q.maybePush(fn)
	if !fnCallback(m, fn) {
		return nil
	}

	for !q.empty() {
		next := q.pop()
		if !fnCallback(m, next) {
			return nil
		}
		for _, step := range next.Steps {
			if options.IgnoreEdge != nil && options.IgnoreEdge(m, next, step) {
				continue
			}
			fnName := step.Function
			fn, ok := m.Functions[fnName]
			if !ok {
				continue
			}
			q.maybePush(fn)
		}
	}
	return nil
}
