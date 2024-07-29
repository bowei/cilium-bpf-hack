package llvmp

// lineRing is a running window that keeps the last .limit lines of context for
// the parser. This allows us to parse multi-line constructs.
type lineRing struct {
	limit int
	lines []string
	count int
}

func (r *lineRing) push(line string) {
	r.lines = append(r.lines, line)
	if len(r.lines) > r.limit {
		r.lines = r.lines[len(r.lines)-r.limit:]
	}
	r.count++
}

func (r *lineRing) cur() string { return r.lines[len(r.lines)-1] }
