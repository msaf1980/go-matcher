package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
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
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a*c", Terminated: "a*c", TermIndex: 0,
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c",
									Inners:  []wildcards.InnerItem{wildcards.ItemStar{}},
									MinSize: 2, MaxSize: -1,
								},
							},
							{
								Node: "a*c*", Terminated: "a*c*", TermIndex: 1,
								WildcardItems: wildcards.WildcardItems{
									P: "a", MinSize: 2, MaxSize: -1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemStar{}, wildcards.ItemRune('c'), wildcards.ItemStar{},
									},
								},
							},
							{
								Node: "a*b?c", Terminated: "a*b?c", TermIndex: 2,
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c", MinSize: 4, MaxSize: -1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemStar{}, wildcards.ItemRune('b'), wildcards.ItemOne{},
									},
								},
							},
						},
					},
					2: {
						Childs: []*NodeItem{
							{
								Node: "a*c",
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
									Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
								},
								Childs: []*NodeItem{
									{
										Node: "b", Terminated: "a*c.b", TermIndex: 3,
										WildcardItems: wildcards.WildcardItems{P: "b", MinSize: 0, MaxSize: 0},
									},
								},
							},
							{
								Node: "a", WildcardItems: wildcards.WildcardItems{P: "a"},
								Childs: []*NodeItem{
									{
										Node: "b?d", Terminated: "a.b?d", TermIndex: 4,
										WildcardItems: wildcards.WildcardItems{
											P: "b", Suffix: "d", MinSize: 3, MaxSize: 3,
											Inners: []wildcards.InnerItem{wildcards.ItemOne{}},
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
