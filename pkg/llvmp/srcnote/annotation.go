package srcnote

import (
	"fmt"
	"strconv"
	"strings"
)

type AnnotationKind string

const (
	KindConditional = AnnotationKind("Conditional")
	KindNote        = AnnotationKind("Note")
)

func validKind(s string) bool {
	switch AnnotationKind(s) {
	case KindConditional:
		return true
	case KindNote:
		return true
	}
	return false
}

// Annotation is a single-line entry in an annotations file. It has the
// following format with colon ':' as the separating character.
//
//	<file>:<line>:<kind>:<tags>:<text>
type Annotation struct {
	FileName string
	Line     int
	Kind     AnnotationKind
	Tags     map[string]string
	Text     string
}

func (a *Annotation) String() string {
	var tags []string
	for k, v := range a.Tags {
		if v == "" {
			tags = append(tags, k)
		} else {
			tags = append(tags, k+"="+v)
		}
	}
	return fmt.Sprintf("%s:%d:%s:%s:%s\n", a.FileName, a.Line, a.Kind, strings.Join(tags, ","), a.Text)
}

func parseAnnotation(pc parseContext, line string) (*Annotation, error) {
	const partsCount = 5
	parts := strings.SplitN(line, ":", partsCount)
	if len(parts) != partsCount {
		return nil, fmt.Errorf("%s:%d: invalid format, not enough ':' separators: %q", pc.fileName, pc.line, line)
	}

	for i, x := range parts {
		parts[i] = strings.TrimSpace(x)
	}
	fileName := parts[0]
	rawLine := parts[1]
	kind := parts[2]
	rawTags := parts[3]
	text := parts[4]

	srcLine, err := strconv.Atoi(rawLine)
	if err != nil {
		return nil, fmt.Errorf("%s:%d: line value is not an integer: %q", pc.fileName, pc.line, line)
	}
	if !validKind(kind) {
		return nil, fmt.Errorf("%s:%d: kind value is invalid: %q", pc.fileName, pc.line, line)
	}

	var tags map[string]string
	addTag := func(k, v string) {
		if tags == nil {
			tags = map[string]string{}
		}
		tags[k] = v
	}
	for _, rawTag := range strings.Split(rawTags, ",") {
		switch {
		case rawTag == "":
		case !strings.Contains(rawTag, "="):
			addTag(rawTag, "")
		default:
			tagParts := strings.Split(rawTag, "=")
			if len(tagParts) != 2 {
				return nil, fmt.Errorf("%s:%d: tag value is invalid: %q", pc.fileName, pc.line, line)
			}
			addTag(tagParts[0], tagParts[1])
		}
	}

	an := &Annotation{
		Kind:     AnnotationKind(kind),
		FileName: fileName,
		Line:     srcLine,
		Tags:     tags,
		Text:     text,
	}

	return an, nil
}
