package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/globs"
	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcherString(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: "empty #1", globs: []string{},
			wantW: &GlobMatcher{
				Root:  map[int]*globs.NodeItem{},
				Globs: map[string]int{},
			},
		},
		{
			name: "empty #2", globs: []string{""},
			wantW: &GlobMatcher{
				Root:  map[int]*globs.NodeItem{},
				Globs: map[string]int{},
			},
		},
		// string match
		{
			name: `{"a"}`, globs: []string{"a"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a", Terminated: []string{"a"},
								NodeItem: items.NodeItem{P: "a"},
							},
						},
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
				Root: map[int]*globs.NodeItem{
					2: {
						Childs: []*globs.NodeItem{
							{
								Node:     "a",
								NodeItem: items.NodeItem{P: "a"},
								Childs: []*globs.NodeItem{
									{
										Node: "bc", Terminated: []string{"a.bc"},
										NodeItem: items.NodeItem{P: "bc"},
									},
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
		runTestGlobMatcher(t, tt)
	}
}
