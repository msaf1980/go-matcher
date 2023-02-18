package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

func TestGlobMatcher_One(t *testing.T) {
	tests := []testGlobMatcher{
		// ? match
		{
			name: `{"?"}`, globs: []string{"?"},
			wantW: &GlobMatcher{
				Root: map[int]*wildcards.NodeItem{
					1: {
						Childs: []*wildcards.NodeItem{
							{
								Node: "?", Terminated: "?", TermIndex: -1, MinSize: 1, MaxSize: 1,
								Inners: []wildcards.InnerItem{wildcards.ItemOne{}},
							},
						},
					},
				},
				Globs: map[string]int{"?": -1},
			},
			matchPaths: map[string][]string{"a": {"?"}, "c": {"?"}},
			missPaths:  []string{"", "ab", "a.b"},
		},
		{
			name: `{"a?"}`, globs: []string{"a?"},
			wantW: &GlobMatcher{
				Root: map[int]*wildcards.NodeItem{
					1: {
						Childs: []*wildcards.NodeItem{
							{
								Node: "a?", Terminated: "a?", TermIndex: -1, MinSize: 2, MaxSize: 2,
								P: "a", Inners: []wildcards.InnerItem{wildcards.ItemOne{}},
							},
						},
					},
				},
				Globs: map[string]int{"a?": -1},
			},
			matchPaths: map[string][]string{"ac": {"a?"}, "az": {"a?"}},
			missPaths:  []string{"", "a", "bc", "ace", "a.c"},
		},
		{
			name: `{"a?c"}`, globs: []string{"a?c"},
			wantW: &GlobMatcher{
				Root: map[int]*wildcards.NodeItem{
					1: {
						Childs: []*wildcards.NodeItem{
							{
								Node: "a?c", Terminated: "a?c", TermIndex: -1, P: "a",
								Inners: []wildcards.InnerItem{wildcards.ItemOne{}}, Suffix: "c",
								MinSize: 3, MaxSize: 3,
							},
						},
					},
				},
				Globs: map[string]int{"a?c": -1},
			},
			matchPaths: map[string][]string{"acc": {"a?c"}, "aec": {"a?c"}},
			missPaths:  []string{"", "ab", "ac", "ace", "a.c"},
		},
		{
			name: `{"a?[c]?d"}`, globs: []string{"a?[c]?d"},
			wantW: &GlobMatcher{
				Root: map[int]*wildcards.NodeItem{
					1: {
						Childs: []*wildcards.NodeItem{
							{
								Node: "a?[c]?d", Terminated: "a?[c]?d", TermIndex: -1, P: "a", Suffix: "d",
								Inners: []wildcards.InnerItem{
									wildcards.ItemOne{},
									wildcards.ItemRune('c'),
									wildcards.ItemOne{},
								},
								MinSize: 5, MaxSize: 5,
							},
						},
					},
				},
				Globs: map[string]int{"a?[c]?d": -1},
			},
			matchPaths: map[string][]string{"aZccd": {"a?[c]?d"}, "aZcAd": {"a?[c]?d"}},
			missPaths:  []string{"", "ab", "ac", "ace", "aZDAd", "a.c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
