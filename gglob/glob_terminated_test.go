package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

func TestGlobMatcherTerminated(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: `{"a", "a.bc", "a.dc", "b.bc"}`, globs: []string{"a", "a.bc", "a.dc", "b.bc"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a", Terminated: "a", TermIndex: -1,
								WildcardItems: wildcards.WildcardItems{P: "a"},
							},
						},
					},
					2: {
						Childs: []*NodeItem{
							{
								Node: "a", WildcardItems: wildcards.WildcardItems{P: "a"},
								Childs: []*NodeItem{
									{
										Node: "bc", Terminated: "a.bc", TermIndex: -1,
										WildcardItems: wildcards.WildcardItems{P: "bc"},
									},
									{
										Node: "dc", Terminated: "a.dc", TermIndex: -1,
										WildcardItems: wildcards.WildcardItems{P: "dc"},
									},
								},
							},
							{
								Node: "b", WildcardItems: wildcards.WildcardItems{P: "b"},
								Childs: []*NodeItem{
									{
										Node: "bc", Terminated: "b.bc", TermIndex: -1,
										WildcardItems: wildcards.WildcardItems{P: "bc"},
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
