package srcnote

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type AnnotationKind string

const (
	KindConditional = "Conditional"
	KindNote        = "Note"
)

func validKind(s string) bool {
	switch s {
	case KindConditional:
		return true
	case KindNote:
		return true
	}
	return false
}

type Annotation struct {
	Kind     string
	FileName string
	Line     int
	Text     string
}

type anList []*Annotation

func (l anList) Len() int           { return len(l) }
func (l anList) Less(i, j int) bool { return l[i].Line < l[j].Line }
func (l anList) Swap(i, j int) {
	t := l[i]
	l[i] = l[j]
	l[j] = t
}

func NewSet() *Set {
	return &Set{
		files: map[string]anList{},
	}
}

type Set struct {
	// files is a map from fileName to list of Annotations.
	// TODO: sort for binary search.
	files map[string]anList
}

// Lookup any annotations that match the range from [start, end).
func (a *Set) Lookup(fileName string, start, end int) []*Annotation {
	var matches anList

	ans, ok := a.files[fileName]
	if !ok {
		return nil
	}

	for _, an := range ans {
		if start <= an.Line && an.Line < end {
			matches = append(matches, an)
		}
	}

	sort.Sort(matches)

	return matches
}

func (a *Set) Add(an *Annotation) {
	a.files[an.FileName] = append(a.files[an.FileName], an)
}

func (a *Set) AddList(al []*Annotation) {
	for _, an := range al {
		a.Add(an)
	}
}

func Load(fileNames ...string) (*Set, error) {
	ret := NewSet()
	for _, fileName := range fileNames {
		alist, err := ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		ret.AddList(alist)
	}
	return ret, nil
}

// ReadFile reads in a file of the format:
//
//	filename:line:kind:text...
//
// Empty lines and lines beginning with "#" will be ignored.
func ReadFile(fileName string) ([]*Annotation, error) {
	var ret []*Annotation

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var curLine int
	for scanner.Scan() {
		curLine++
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}

		const partsCount = 4
		parts := strings.SplitN(line, ":", partsCount)
		if len(parts) != partsCount {
			return nil, fmt.Errorf("%s:%d: invalid format: %q", fileName, curLine, line)
		}

		srcLine, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("%s:%d: line value is not an integer: %q", fileName, curLine, line)
		}

		if !validKind(parts[2]) {
			return nil, fmt.Errorf("%s:%d: kind value is invalid: %q", fileName, curLine, line)
		}

		an := &Annotation{
			Kind:     parts[2],
			FileName: parts[0],
			Line:     srcLine,
			Text:     parts[3],
		}
		ret = append(ret, an)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ReadAnnotationFile:%w", err)
	}

	return ret, nil
}
