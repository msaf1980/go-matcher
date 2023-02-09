package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcherString(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: "empty #1", globs: []string{},
			wantW: &GlobMatcher{
				Root:  map[int]*items.NodeItem{},
				Globs: map[string]bool{},
			},
		},
		{
			name: "empty #2", globs: []string{""},
			wantW: &GlobMatcher{
				Root:  map[int]*items.NodeItem{},
				Globs: map[string]bool{},
			},
		},
		// string match
		{
			name: `{"a"}`, globs: []string{"a"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{{Node: "a", Terminated: "a", P: "a"}},
					},
				},
				Globs: map[string]bool{"a": true},
			},
			matchPaths: map[string][]string{"a": {"a"}},
			missPaths:  []string{"", "b", "ab", "ba"},
		},
		{
			name: `{"a.bc"}`, globs: []string{"a.bc"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					2: {
						Childs: []*items.NodeItem{
							{
								Node: "a", P: "a",
								Childs: []*items.NodeItem{{Node: "bc", Terminated: "a.bc", P: "bc"}},
							},
						},
					},
				},
				Globs: map[string]bool{"a.bc": true},
			},
			matchPaths: map[string][]string{"a.bc": {"a.bc"}},
			missPaths:  []string{"", "b", "ab", "bc", "abc", "b.bc", "a.bce", "a.bc.e"},
		},
		{
			name: `{"a", "a.bc", "a.dc", "b.bc"}`, globs: []string{"a", "a.bc", "a.dc", "b.bc"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{Node: "a", Terminated: "a", P: "a"},
						},
					},
					2: {
						Childs: []*items.NodeItem{
							{
								Node: "a", P: "a",
								Childs: []*items.NodeItem{
									{Node: "bc", Terminated: "a.bc", P: "bc"},
									{Node: "dc", Terminated: "a.dc", P: "dc"},
								},
							},
							{
								Node: "b", P: "b",
								Childs: []*items.NodeItem{
									{Node: "bc", Terminated: "b.bc", P: "bc"},
								},
							},
						},
					},
				},
				Globs: map[string]bool{
					"a":    true,
					"a.bc": true,
					"a.dc": true,
					"b.bc": true,
				},
			},
			matchPaths: map[string][]string{
				"a":    {"a"},
				"a.bc": {"a.bc"},
				"a.dc": {"a.dc"},
				"b.bc": {"b.bc"},
			},
			missPaths: []string{"", "b", "ab", "bc", "abc", "c.bc", "a.be", "a.bce", "a.bc.e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
