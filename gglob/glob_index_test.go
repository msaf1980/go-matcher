package gglob

import (
	"testing"

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
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a*c", Terminated: "a*c", TermIndex: 0, P: "a", Suffix: "c",
								Inners:  []items.InnerItem{items.ItemStar{}},
								MinSize: 2, MaxSize: -1,
							},
							{
								Node: "a*c*", Terminated: "a*c*", TermIndex: 1, P: "a",
								MinSize: 2, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}, items.ItemString("c"), items.ItemStar{}},
							},
							{
								Node: "a*b?c", Terminated: "a*b?c", TermIndex: 2, P: "a", Suffix: "c",
								MinSize: 4, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}, items.ItemString("b"), items.ItemOne{}},
							},
						},
					},
					2: {
						Childs: []*items.NodeItem{
							{
								Node: "a*c", P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}},
								Childs: []*items.NodeItem{
									{Node: "b", Terminated: "a*c.b", TermIndex: 3, P: "b", MinSize: 0, MaxSize: 0},
								},
							},
							{
								Node: "a", P: "a",
								Childs: []*items.NodeItem{
									{
										Node: "b?d", Terminated: "a.b?d", TermIndex: 4, P: "b", Suffix: "d", MinSize: 3, MaxSize: 3,
										Inners: []items.InnerItem{items.ItemOne{}},
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
