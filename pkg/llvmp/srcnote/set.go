package srcnote

import "sort"

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
	sort.Sort(a.files[an.FileName])
}

func (a *Set) AddList(al []*Annotation) {
	for _, an := range al {
		a.Add(an)
	}
}
