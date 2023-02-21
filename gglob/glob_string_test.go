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
				Root:  map[int]*NodeItem{},
				Globs: map[string]int{},
			},
		},
		{
			name: "empty #2", globs: []string{""},
			wantW: &GlobMatcher{
				Root:  map[int]*NodeItem{},
				Globs: map[string]int{},
			},
		},
		// string match
		{
			name: `{"a"}`, globs: []string{"a"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a", Terminated: []string{"a"},
								WildcardItems: wildcards.WildcardItems{P: "a"},
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
				Root: map[int]*NodeItem{
					2: {
						Childs: []*NodeItem{
							{
								Node:          "a",
								WildcardItems: wildcards.WildcardItems{P: "a"},
								Childs: []*NodeItem{
									{
										Node: "bc", Terminated: []string{"a.bc"},
										WildcardItems: wildcards.WildcardItems{P: "bc"},
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
