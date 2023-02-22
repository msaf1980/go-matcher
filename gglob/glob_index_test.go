package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/globs"
	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_Index(t *testing.T) {
	tests := []testGlobMatcherIndex{
		// composite
		{
			name: `{"a*c", "a*c*", "a*b?c", "a*c.b", "a.b?d"}`,
			globs: []string{
				"a*c", "a*c*", "a*b?c", "a*c.b", "a.b?d",
			},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a*c", Terminated: []string{"a*c"}, TermIndex: []int{0},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "c",
									Inners:  []items.Item{items.ItemStar{}},
									MinSize: 2, MaxSize: -1,
								},
							},
							{
								Node: "a*c*", Terminated: []string{"a*c*"}, TermIndex: []int{1},
								NodeItem: items.NodeItem{
									P: "a", MinSize: 2, MaxSize: -1,
									Inners: []items.Item{
										items.ItemStar{}, items.ItemRune('c'), items.ItemStar{},
									},
								},
							},
							{
								Node: "a*b?c", Terminated: []string{"a*b?c"}, TermIndex: []int{2},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "c", MinSize: 4, MaxSize: -1,
									Inners: []items.Item{
										items.ItemStar{}, items.ItemRune('b'), items.ItemOne{},
									},
								},
							},
						},
					},
					2: {
						Childs: []*globs.NodeItem{
							{
								Node: "a*c",
								NodeItem: items.NodeItem{
									P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
									Inners: []items.Item{items.ItemStar{}},
								},
								Childs: []*globs.NodeItem{
									{
										Node: "b", Terminated: []string{"a*c.b"}, TermIndex: []int{3},
										NodeItem: items.NodeItem{P: "b", MinSize: 0, MaxSize: 0},
									},
								},
							},
							{
								Node: "a", NodeItem: items.NodeItem{P: "a"},
								Childs: []*globs.NodeItem{
									{
										Node: "b?d", Terminated: []string{"a.b?d"}, TermIndex: []int{4},
										NodeItem: items.NodeItem{
											P: "b", Suffix: "d", MinSize: 3, MaxSize: 3,
											Inners: []items.Item{items.ItemOne{}},
										},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a*c": 0, "a*c*": 1, "a*b?c": 2, "a*c.b": 3, "a.b?d": 4},
			},
			matchPaths: map[string][]int{
				"acbec":  {0, 1, 2},
				"abbece": {1},
				"a.bfd":  {4},
				// not matched
				"": nil, "ab": {}, "c": {}, "a.b": {}, "a.bd": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcherIndex(t, tt)
		})
	}
}
