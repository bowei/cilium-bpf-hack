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
		{KindNote, "file1", 100, "some text"},
		{KindNote, "file1", 101, "some text"},
		{KindNote, "file1", 102, "some text:some text"},
	}); diff != "" {
		t.Errorf("Diff =\n%s", diff)
	}

	t.Logf("%+v", alist)
	alist, err = ReadFile("testinput_invalid.txt")
	if err == nil {
		t.Errorf("ReadFile() = %v, want != nil", err)
	}
}
