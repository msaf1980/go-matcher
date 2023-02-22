package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/globs"
	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_Many(t *testing.T) {
	tests := []testGlobMatcher{
		// ? match
		{
			name: `{"??"}`, globs: []string{"??"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "??", Terminated: []string{"??"},
								NodeItem: items.NodeItem{
									MinSize: 2, MaxSize: 2,
									Inners: []items.Item{items.ItemMany(2)},
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
