package ignore

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExpand(t *testing.T) {
	// No t.parallel: modifies global vars.
	oldIgnoredFn := ignoredFn
	defer func() { ignoredFn = oldIgnoredFn }()

	for _, tc := range []struct {
		name      string
		ignoredFn map[string][]string
		in        []string
		want      []string
		wantErr   bool
	}{
		{
			name: "no_expansion",
			in:   []string{"a", "b"},
			want: []string{"a", "b"},
		},
		{
			name: "single_expansion",
			ignoredFn: map[string][]string{
				"@b": {"b1", "b2"},
			},
			in:   []string{"a", "@b", "c"},
			want: []string{"a", "b1", "b2", "c"},
		},
		{
			name: "multi_expansion",
			ignoredFn: map[string][]string{
				"@b":  {"@b1", "b2"},
				"@b1": {"b11", "b12"},
			},
			in:   []string{"a", "@b", "c"},
			want: []string{"a", "b11", "b12", "b2", "c"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ignoredFn = tc.ignoredFn
			got, err := expand(tc.in)
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Fatalf("err = %v, wantErr = %t", err, tc.wantErr)
			}
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("Diff (-got,+want) =\n%s", diff)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	s := Set{
		"fn1": true,
	}
	for _, tc := range []struct {
		name string
		in   string
		s    Set
		want bool
	}{
		{
			name: "no_match",
			in:   "fn2",
			s:    s,
		},

		{
			name: "no_match",
			in:   "fn1",
			s:    s,
			want: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.s.Match(tc.in); got != tc.want {
				t.Errorf("got = %t, want = %t", got, tc.want)
			}
		})
	}
}
