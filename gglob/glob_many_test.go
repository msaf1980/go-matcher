package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

func TestGlobMatcher_Many(t *testing.T) {
	tests := []testGlobMatcher{
		// ? match
		{
			name: `{"??"}`, globs: []string{"??"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "??", Terminated: []string{"??"},
								WildcardItems: wildcards.WildcardItems{
									MinSize: 2, MaxSize: 2,
									Inners: []wildcards.InnerItem{wildcards.ItemMany(2)},
								},
							},
						},
					},
				},
				Globs: map[string]int{"??": -1},
			},
			matchPaths: map[string][]string{"ab": {"??"}, "ac": {"??"}},
			missPaths:  []string{"", "a", "abc", "a.b"},
		},
	}
	for _, tt := range tests {
		runTestGlobMatcher(t, tt)
	}
}
