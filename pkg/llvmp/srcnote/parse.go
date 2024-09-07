package srcnote

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type anList []*Annotation

func (l anList) Len() int { return len(l) }
func (l anList) Less(i, j int) bool {
	if l[i].FileName == l[j].FileName {
		return l[i].Line < l[j].Line
	}
	return l[i].FileName < l[j].FileName
}

func (l anList) Swap(i, j int) {
	t := l[i]
	l[i] = l[j]
	l[j] = t
}

// Load source annotations from the fileNames.
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

type parseContext struct {
	fileName string
	line     int
}

// ReadFile parses a single annotation file.
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
		line = strings.TrimSpace(line)

		if line == "" || line[0] == '#' {
			continue
		}

		an, err := parseAnnotation(parseContext{fileName: fileName, line: curLine}, line)
		if err != nil {
			return nil, err
		}
		ret = append(ret, an)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ReadAnnotationFile:%w", err)
	}

	return ret, nil
}
