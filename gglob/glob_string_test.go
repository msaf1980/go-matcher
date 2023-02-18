package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

func TestGlobMatcherString(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: "empty #1", globs: []string{},
			wantW: &GlobMatcher{
				Root:  map[int]*wildcards.NodeItem{},
				Globs: map[string]int{},
			},
		},
		{
			name: "empty #2", globs: []string{""},
			wantW: &GlobMatcher{
				Root:  map[int]*wildcards.NodeItem{},
				Globs: map[string]int{},
			},
		},
		// string match
		{
			name: `{"a"}`, globs: []string{"a"},
			wantW: &GlobMatcher{
				Root: map[int]*wildcards.NodeItem{
					1: {
						Childs: []*wildcards.NodeItem{{Node: "a", Terminated: "a", TermIndex: -1, P: "a"}},
					},
				},
				Globs: map[string]int{"a": -1},
			},
			matchPaths: map[string][]string{"a": {"a"}},
			missPaths:  []string{"", "b", "ab", "ba"},
		},
		{
			name: `{"a.bc"}`, globs: []string{"a.bc"},
			wantW: &GlobMatcher{
				Root: map[int]*wildcards.NodeItem{
					2: {
						Childs: []*wildcards.NodeItem{
							{
								Node: "a", P: "a",
								Childs: []*wildcards.NodeItem{
									{Node: "bc", Terminated: "a.bc", TermIndex: -1, P: "bc"},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a.bc": -1},
			},
			matchPaths: map[string][]string{
				"a.bc": {"a.bc"},
				// last dot
				"a.bc.": {"a.bc"},
			},
			missPaths: []string{"", "b", "ab", "bc", "abc", "b.bc", "a.bce", "a.bc.e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
