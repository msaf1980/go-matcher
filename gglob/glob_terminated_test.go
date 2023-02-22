package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/globs"
	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcherTerminated(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: `{"a", "a.bc", "a.dc", "b.bc"}`, globs: []string{"a", "a.bc", "a.dc", "b.bc"},
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
					2: {
						Childs: []*globs.NodeItem{
							{
								Node: "a", NodeItem: items.NodeItem{P: "a"},
								Childs: []*globs.NodeItem{
									{
										Node: "bc", Terminated: []string{"a.bc"},
										NodeItem: items.NodeItem{P: "bc"},
									},
									{
										Node: "dc", Terminated: []string{"a.dc"},
										NodeItem: items.NodeItem{P: "dc"},
									},
								},
							},
							{
								Node: "b", NodeItem: items.NodeItem{P: "b"},
								Childs: []*globs.NodeItem{
									{
										Node: "bc", Terminated: []string{"b.bc"},
										NodeItem: items.NodeItem{P: "bc"},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{
					"a":    -1,
					"a.bc": -1,
					"a.dc": -1,
					"b.bc": -1,
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
		runTestGlobMatcher(t, tt)
	}
}
