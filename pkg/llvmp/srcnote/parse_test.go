package srcnote

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	alist, err := ReadFile("testinput_valid.txt")
	if err != nil {
		t.Errorf("ReadFile() = %v, want nil", err)
	}

	if diff := cmp.Diff(alist, []*Annotation{
		{"file1", 100, KindConditional, nil, "some text"},
		{"file1", 101, KindConditional, map[string]string{"tag1": ""}, "Note:some text"},
		{"file1", 102, KindNote, map[string]string{"tag2": "abc", "tag3": ""}, "Note:some text:some text"},
	}); diff != "" {
		t.Errorf("Diff =\n%s", diff)
	}

	t.Logf("%+v", alist)
	_, err = ReadFile("testinput_invalid.txt")
	if err == nil {
		t.Errorf("ReadFile() = %v, want != nil", err)
	}
}
