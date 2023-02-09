package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_Multi(t *testing.T) {
	tests := []testGlobMatcher{
		// composite
		{
			name: `{"a*c", "a*c*", "a*b?c", "a.b?d", "a*c.b"}`, globs: []string{"a*c", "a*c*", "a*b?c", "a.b?d", "a*c.b"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a*c", Terminated: "a*c", P: "a", Suffix: "c",
								Inners:  []items.InnerItem{items.ItemStar{}},
								MinSize: 2, MaxSize: -1,
							},
							{
								Node: "a*c*", Terminated: "a*c*", P: "a",
								MinSize: 2, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}, items.ItemString("c"), items.ItemStar{}},
							},
							{
								Node: "a*b?c", Terminated: "a*b?c", P: "a", Suffix: "c",
								MinSize: 4, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}, items.ItemString("b"), items.ItemOne{}},
							},
						},
					},
					2: {
						Childs: []*items.NodeItem{
							{
								Node: "a", P: "a",
								Childs: []*items.NodeItem{
									{
										Node: "b?d", Terminated: "a.b?d", P: "b", Suffix: "d", MinSize: 3, MaxSize: 3,
										Inners: []items.InnerItem{items.ItemOne{}},
									},
								},
							},
							{
								Node: "a*c", P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}},
								Childs: []*items.NodeItem{
									{Node: "b", Terminated: "a*c.b", P: "b", MinSize: 0, MaxSize: 0},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"a*c": true, "a*c*": true, "a*b?c": true, "a*c.b": true, "a.b?d": true},
			},
			matchPaths: map[string][]string{
				"acbec":  {"a*c", "a*c*", "a*b?c"},
				"abbece": {"a*c*"},
				"a.bfd":  {"a.b?d"},
			},
			missPaths: []string{"", "ab", "c", "a.b", "a.bd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
