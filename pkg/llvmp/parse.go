package llvmp

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/bowei/cilium-bpf-hack/pkg/cilconst"
)

func newParseContext() *parseContext {
	return &parseContext{
		m:             newModule(),
		lines:         lineRing{limit: 30},
		all:           map[int]interface{}{},
		locations:     map[int]location{},
		files:         map[int]sourceFile{},
		lexicalBlocks: map[int]lexicalBlock{},
		subprogram:    map[int]subprogram{},
	}
}

type parseContext struct {
	m     *Module
	curFn *FnDef

	lines lineRing

	all           map[int]interface{}
	locations     map[int]location
	files         map[int]sourceFile
	lexicalBlocks map[int]lexicalBlock
	subprogram    map[int]subprogram
}

type sourceRef struct {
	file string
	line int
}

func (c *parseContext) String() string {
	return fmt.Sprintf("parseContext:%d:%q", c.lines.count, c.lines.cur())
}

func (c *parseContext) dumpStdout() {
	fmt.Println("----")
	fmt.Println(c.m.Dump())
	fmt.Println(c.curFn)
	fmt.Println(c.lines)

	fmt.Println("-- locations:")
	for k, x := range c.locations {
		fmt.Printf("%v: %+v\n", k, x)
	}
	fmt.Println("-- files:")
	for k, x := range c.files {
		fmt.Printf("%v: %+v\n", k, x)
	}
	fmt.Println("-- lexicalBlocks:")
	for k, x := range c.lexicalBlocks {
		fmt.Printf("%v: %+v\n", k, x)
	}
	fmt.Println("-- subprogram:")
	for k, x := range c.subprogram {
		fmt.Printf("%v: %+v\n", k, x)
	}
	fmt.Println("----")
}

func (c *parseContext) lookupFunc(id int) (*sourceRef, error) {
	sp, ok := c.subprogram[id]
	if !ok {
		return nil, fmt.Errorf("XXX")
	}
	f, ok := c.files[sp.file]
	if !ok {
		return nil, fmt.Errorf("XXX")
	}
	ret := &sourceRef{
		file: f.fileName,
		line: sp.line,
	}
	return ret, nil
}

func (c *parseContext) lookupLocation(id int) (*sourceRef, error) {
	l, ok := c.locations[id]
	if !ok {
		return nil, fmt.Errorf("XXX")
	}

	var fileName string
	sp, ok := c.subprogram[l.scope]
	if ok {
		f, ok := c.files[sp.file]
		if !ok {
			return nil, fmt.Errorf("XXX")
		}
		fileName = f.fileName
	} else if lb, ok := c.lexicalBlocks[l.scope]; ok {
		f, ok := c.files[lb.file]
		if !ok {
			return nil, fmt.Errorf("XXX")
		}
		fileName = f.fileName
	} else {
		return nil, fmt.Errorf("XXX")
	}

	ret := &sourceRef{
		file: fileName,
		line: l.line,
	}

	return ret, nil
}

// ParseLL output from a compilation.
func ParseLL(fileName string) (*Module, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("ParseLL:Open:%w", err)
	}
	defer f.Close()

	pc := newParseContext()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		pc.lines.push(line)

		for _, m := range []struct {
			r *regexp.Regexp
			f func(*parseContext) error
		}{
			{fnStartRe, parseFnStart},
			{fnEndRe, parseFnEnd},
			{tcInternalRe, parseTCInternal},
			{tcDyanmicRe, parseTCDynamic},
			{tcPolicyRe, parseTCPolicy},
			{tcEgressPolicyRe, parseTCEgressPolicy},
			{callRe, parseCall},
			{retRe, parseRet},
			{diLexicalBlockRe, parseDILexicalBlock},
			{diLocationRe, parseDILocation},
			{diFileRe, parseDIFile},
			{diSubprogramRe, parseDISubprogram},
		} {
			if !m.r.MatchString(line) {
				continue
			}
			err := m.f(pc)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ParseLL:scanner:%w", err)
	}

	if err := resolveSources(pc); err != nil {
		return nil, err
	}

	// pc.dumpStdout()

	return pc.m, nil
}

func resolveSources(pc *parseContext) error {
	for _, fn := range pc.m.Functions {
		sref, err := pc.lookupFunc(fn.dbgRef)
		if err != nil {
			fn.File = "not found"
		} else {
			fn.File = sref.file
			fn.Line = sref.line
		}
		for _, st := range fn.Steps {
			sref, err := pc.lookupLocation(st.dbgRef)
			if err != nil {
				st.File = "not found"
			} else {
				st.File = sref.file
				st.Line = sref.line
			}
		}
	}
	return nil
}

var (
	fnStartRe   = regexp.MustCompile("^define.*{")
	fnSectionRe = regexp.MustCompile(` +section "[0-9]+/[0-9]+" +`)
	fnStart2Re  = regexp.MustCompile(`^define (internal|dso_local) [a-zA-Z0-9_]+ @([a-zA-Z0-9_]+)\(.* !dbg ![0-9]+ {`)
	fnEndRe     = regexp.MustCompile("^}")
)

func parseFnStart(pc *parseContext) error {
	line := pc.lines.cur()
	matches := fnStart2Re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return fmt.Errorf("parseFnStart:%v", pc)
	}

	fnLinkage, fnName := matches[1], matches[2]
	curFn := pc.m.addFn(fnName)
	pc.curFn = curFn

	curFn.dbgRef = debugRef(line)
	curFn.Linkage = fnLinkage

	switch fnLinkage {
	case "internal":
		if fnSectionRe.MatchString(pc.lines.cur()) {
			curFn.Kind = FnKindTail
		} else {
			curFn.Kind = FnKindInternal
		}
	case "dso_local":
		curFn.Kind = FnKindTail
	}

	return nil
}

func parseFnEnd(pc *parseContext) error {
	if pc.curFn == nil {
		return fmt.Errorf("parseFnEnd:no fn:%v", pc)
	}
	pc.curFn = nil
	return nil
}

var tcInternalRe = regexp.MustCompile(` *(%[0-9]+ =|) *call i32 @tail_call_internal\(ptr noundef %[0-9]+, i32 noundef ([0-9]+),.*\).*`)

func parseTCInternal(pc *parseContext) error {
	line := pc.lines.cur()

	matches := tcInternalRe.FindStringSubmatch(line)
	if len(matches) == 0 {
		return fmt.Errorf("parseTCInternal:no match:%v", pc)
	}

	idx, err := strconv.Atoi(matches[2])

	if err != nil {
		return fmt.Errorf("parseTCInternal:bad idx:%v:%q:%w", pc, matches[1], err)
	}

	if pc.curFn == nil {
		return fmt.Errorf("parseFnEnd:no fn:no cur_fn:%v", pc)
	}

	step := pc.curFn.addStep()
	step.Kind = StepTailCall
	step.Function = cilconst.TailCallMap[idx]
	step.dbgRef = debugRef(line)
	step.line = line

	return nil
}

var tcDyanmicRe = regexp.MustCompile(` *(%[0-9]+ =|) *call void @tail_call_dynamic\(ptr noundef %[0-9]+,.*\).*`)

func parseTCDynamic(pc *parseContext) error {
	return nil
}

var tcPolicyRe = regexp.MustCompile(` *(%[0-9]+ =|) *call i32 @tail_call_policy\(ptr noundef %[0-9]+.*\).*`)

func parseTCPolicy(pc *parseContext) error {
	return nil
}

var tcEgressPolicyRe = regexp.MustCompile(` *(%[0-9]+ =|) *call i32 @tail_call_egress_policy\(ptr noundef %[0-9]+,.*\).*`)

func parseTCEgressPolicy(pc *parseContext) error {
	return nil
}

var (
	callRe = regexp.MustCompile(`^ +(%[0-9]+ = call|call)`)
	// TODO: looks like there is sometimes an indirect jmp that is loaded.
	//
	//   %7 = load ptr, ptr @map_lookup_elem, align 8, !dbg !11158
	//   %8 = call ptr %7(ptr noundef @test_cilium_lxc, ptr noundef %3), !dbg !11158
	callIndirectRe      = regexp.MustCompile(` *(%[0-9]+ = call.*|call) (ptr|i32|i64|void) %[0-9]+.*`)
	callSymRe           = regexp.MustCompile(` *(%[0-9]+ = call.*|call).*@([a-zA-Z0-9_]+).*`)
	callAsmSideEffectRe = regexp.MustCompile(` *(%[0-9]+ =|) *call [a-zA-Z0-9]+ asm sideeffect.*`)
)

func parseCall(pc *parseContext) error {
	line := pc.lines.cur()

	// Ignore these for now.
	switch {
	case callIndirectRe.MatchString(line):
		return nil
	case callAsmSideEffectRe.MatchString(line):
		return nil
	}

	matches := callSymRe.FindStringSubmatch(line)
	if len(matches) != 3 {
		return fmt.Errorf("parseCall:%v", pc)
	}

	// ignore matches[1]
	fnName := matches[2]

	step := pc.curFn.addStep()
	step.Kind = StepFnCall
	step.Function = fnName
	step.dbgRef = debugRef(line)
	step.line = line

	return nil
}

var retRe = regexp.MustCompile(` *ret.*!dbg !([0-9]+)`)

func parseRet(pc *parseContext) error {
	line := pc.lines.cur()

	matches := retRe.FindStringSubmatch(line)
	if len(matches) != 2 {
		return fmt.Errorf("parseRet:%v", pc)
	}

	dbg, err := strconv.Atoi(matches[1])
	if err != nil {
		return fmt.Errorf("parseRet:bad_int:%v", pc)
	}

	step := pc.curFn.addStep()
	step.Kind = StepRet
	step.dbgRef = dbg

	return nil
}

type lexicalBlock struct {
	id    int
	scope int
	file  int
}

var diLexicalBlockRe = regexp.MustCompile(` *!([0-9]+) = distinct !DILexicalBlock\(scope: !([0-9]+), file: !([0-9]+),.*\)`)

func parseDILexicalBlock(pc *parseContext) error {
	matches := diLexicalBlockRe.FindStringSubmatch(pc.lines.cur())
	if len(matches) != 4 {
		return fmt.Errorf("diLexicalBlock:no_match:%v", pc)
	}

	sid, sscope, sfile := matches[1], matches[2], matches[3]
	id, err := strconv.Atoi(sid)
	if err != nil {
		return fmt.Errorf("diLexicalBlock:bad_int:%v:%v", pc, err)
	}
	scope, err := strconv.Atoi(sscope)
	if err != nil {
		return fmt.Errorf("diLexicalBlock:bad_int:%v:%v", pc, err)
	}
	file, err := strconv.Atoi(sfile)
	if err != nil {
		return fmt.Errorf("diLexicalBlock:bad_int:%v:%v", pc, err)
	}

	lb := lexicalBlock{id: id, file: file, scope: scope}
	pc.all[id] = lb
	pc.lexicalBlocks[id] = lb

	return nil
}

type location struct {
	id    int
	line  int
	col   int
	scope int
}

var diLocationRe = regexp.MustCompile(` *!([0-9]+) = !DILocation\(line: ([0-9]+), column: ([0-9]+), scope: !([0-9]+)\)`)

func parseDILocation(pc *parseContext) error {
	matches := diLocationRe.FindStringSubmatch(pc.lines.cur())
	if len(matches) != 5 {
		return fmt.Errorf("parseDILocation:no_match:%v", pc)
	}

	sid, sline, scol, sscope := matches[1], matches[2], matches[3], matches[4]
	id, err := strconv.Atoi(sid)
	if err != nil {
		return fmt.Errorf("parseDILocation:bad_int:%v:%v", pc, err)
	}
	line, err := strconv.Atoi(sline)
	if err != nil {
		return fmt.Errorf("parseDILocation:bad_int:%v:%v", pc, err)
	}
	col, err := strconv.Atoi(scol)
	if err != nil {
		return fmt.Errorf("parseDILocation:bad_int:%v:%v", pc, err)
	}
	scope, err := strconv.Atoi(sscope)
	if err != nil {
		return fmt.Errorf("parseDILocation:bad_int:%v:%v", pc, err)
	}

	l := location{
		id:    id,
		line:  line,
		col:   col,
		scope: scope,
	}
	pc.locations[id] = l
	pc.all[id] = l

	return nil
}

type subprogram struct {
	id    int
	name  string
	file  int
	line  int
	scope int
}

var diSubprogramRe = regexp.MustCompile(`!([0-9]+) = distinct !DISubprogram\(name: "([^"]+)", scope: !([0-9]+), file: !([0-9]+), line: ([0-9]+).*\)`)

func parseDISubprogram(pc *parseContext) error {
	matches := diSubprogramRe.FindStringSubmatch(pc.lines.cur())
	if len(matches) != 6 {
		return fmt.Errorf("parseDISubprogram:no_match:%v", pc)
	}

	sid, name, sscope, sfile, sline := matches[1], matches[2], matches[3], matches[4], matches[5]
	id, err := strconv.Atoi(sid)
	if err != nil {
		return fmt.Errorf("parseDISubprogram:bad_int:%v:%v", pc, err)
	}
	scope, err := strconv.Atoi(sscope)
	if err != nil {
		return fmt.Errorf("parseDISubprogram:bad_int:%v:%v", pc, err)
	}
	file, err := strconv.Atoi(sfile)
	if err != nil {
		return fmt.Errorf("parseDISubprogram:bad_int:%v:%v", pc, err)
	}
	line, err := strconv.Atoi(sline)
	if err != nil {
		return fmt.Errorf("parseDISubprogram:bad_int:%v:%v", pc, err)
	}

	sp := subprogram{
		id:    id,
		name:  name,
		file:  file,
		line:  line,
		scope: scope,
	}
	pc.subprogram[id] = sp
	pc.all[id] = sp

	return nil
}

type sourceFile struct {
	id       int
	fileName string
}

var diFileRe = regexp.MustCompile(`!([0-9]+) = !DIFile\(filename: "([^"]+)", directory: "[^"]+", checksumkind: .*, checksum: "[^"]+"\)`)

func parseDIFile(pc *parseContext) error {
	matches := diFileRe.FindStringSubmatch(pc.lines.cur())
	if len(matches) != 3 {
		return fmt.Errorf("parseDIFile:no_match:%v", pc)
	}

	sid, fileName := matches[1], matches[2]
	id, err := strconv.Atoi(sid)
	if err != nil {
		return fmt.Errorf("parseDIFile:bad_int:%v:%v", pc, err)
	}

	sf := sourceFile{id: id, fileName: fileName}
	pc.all[id] = sf
	pc.files[id] = sf

	return nil
}

var debugRefRe = regexp.MustCompile(`!dbg !([0-9]+)`)

// debugRef returns the !dbg reference or -1 if there is no matching reference.
func debugRef(line string) int {
	matches := debugRefRe.FindStringSubmatch(line)
	if len(matches) != 2 {
		return -1
	}
	// TODO: this eats the Atoi error.
	ref, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1
	}
	return ref
}
