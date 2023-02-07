package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_Star(t *testing.T) {
	tests := []testGlobMatcher{
		// deduplication
		{
			name: `{"a******c"}`, globs: []string{"a******c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a******c", Terminated: "a******c", P: "a", Suffix: "c",
								MinSize: 2, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}},
							},
						},
					},
				},
				Globs: map[string]bool{"a******c": true},
			},
			matchGlobs: map[string][]string{"ac": {"a******c"}, "abc": {"a******c"}, "abcc": {"a******c"}},
			miss:       []string{"", "acb"},
		},
		// * match
		{
			name: `{"*"}`, globs: []string{"*"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{Node: "*", Terminated: "*", Inners: []items.InnerItem{items.ItemStar{}}, MaxSize: -1},
						},
					},
				},
				Globs: map[string]bool{"*": true},
			},
			matchGlobs: map[string][]string{"a": {"*"}, "b": {"*"}, "ce": {"*"}},
			miss:       []string{"", "b.c"},
		},
		{
			name: `{"a*c"}`, globs: []string{"a*c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a*c", Terminated: "a*c", P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}},
							},
						},
					},
				},
				Globs: map[string]bool{"a*c": true},
			},
			matchGlobs: map[string][]string{
				"ac": {"a*c"}, "acc": {"a*c"}, "aec": {"a*c"}, "aebc": {"a*c"},
				"aecc": {"a*c"}, "aecec": {"a*c"}, "abecec": {"a*c"},
			},
			miss: []string{"", "ab", "c", "ace", "a.c"},
		},
		// composite
		{
			name: `{"a*b?c"}`, globs: []string{"a*b?c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a*b?c", Terminated: "a*b?c", P: "a", Suffix: "c",
								MinSize: 4, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}, items.ItemString("b"), items.ItemOne{}},
							},
						},
					},
				},
				Globs: map[string]bool{"a*b?c": true},
			},
			matchGlobs: map[string][]string{
				"abec":   {"a*b?c"}, // skip *
				"abbec":  {"a*b?c"}, /// shit first b
				"acbbc":  {"a*b?c"},
				"aecbec": {"a*b?c"},
			},
			miss: []string{"", "ab", "c", "ace", "a.c", "abbece"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
