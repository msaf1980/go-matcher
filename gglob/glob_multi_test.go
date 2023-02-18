package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

func TestGlobMatcher_Multi(t *testing.T) {
	tests := []testGlobMatcher{
		// composite
		{
			name:  `{"a*c", "a*c*", "a*b?c", "a*bd?c", "a*{Z,Q}bd?c", "a.b?d", "a*c.b", "a*[b-e].b"}`,
			globs: []string{"a*c", "a*c*", "a*b?c", "a*bd?c", "a*{Z,Q}bd?c", "a.b?d", "a*c.b", "a*[b-e].b"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a*c", Terminated: "a*c", TermIndex: -1,
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c",
									Inners:  []wildcards.InnerItem{wildcards.ItemStar{}},
									MinSize: 2, MaxSize: -1,
								},
							},
							{
								Node: "a*c*", Terminated: "a*c*", TermIndex: -1,
								WildcardItems: wildcards.WildcardItems{
									P:       "a",
									MinSize: 2, MaxSize: -1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemStar{}, wildcards.ItemRune('c'), wildcards.ItemStar{},
									},
								},
							},
							{
								Node: "a*b?c", Terminated: "a*b?c", TermIndex: -1,
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c",
									MinSize: 4, MaxSize: -1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemStar{}, wildcards.ItemRune('b'), wildcards.ItemOne{},
									},
								},
							},
							{
								Node: "a*bd?c", Terminated: "a*bd?c", TermIndex: -1,
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c",
									MinSize: 5, MaxSize: -1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemStar{}, wildcards.ItemString("bd"), wildcards.ItemOne{},
									},
								},
							},
							{
								Node: "a*{Z,Q}bd?c", Terminated: "a*{Z,Q}bd?c", TermIndex: -1,
								WildcardItems: wildcards.WildcardItems{
									MinSize: 6, MaxSize: -1, P: "a", Suffix: "c",
									Inners: []wildcards.InnerItem{
										wildcards.ItemStar{},
										&wildcards.ItemList{Vals: []string{"Q", "Z"}, ValsMin: 1, ValsMax: 1},
										wildcards.ItemString("bd"), wildcards.ItemOne{},
									},
								},
							},
						},
					},
					2: {
						Childs: []*NodeItem{
							{
								Node: "a", WildcardItems: wildcards.WildcardItems{P: "a"},
								Childs: []*NodeItem{
									{
										Node: "b?d", Terminated: "a.b?d", TermIndex: -1,
										WildcardItems: wildcards.WildcardItems{
											P: "b", Suffix: "d", MinSize: 3, MaxSize: 3,
											Inners: []wildcards.InnerItem{wildcards.ItemOne{}},
										},
									},
								},
							},
							{
								Node: "a*c", WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
									Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
								},
								Childs: []*NodeItem{
									{
										Node: "b", Terminated: "a*c.b", TermIndex: -1,
										WildcardItems: wildcards.WildcardItems{P: "b"},
									},
								},
							},
							{
								Node: "a*[b-e]",
								WildcardItems: wildcards.WildcardItems{
									MinSize: 2,
									MaxSize: -1,
									P:       "a",
									Inners: []wildcards.InnerItem{
										wildcards.ItemStar{},
										wildcards.ItemRuneRanges{{'b', 'e'}},
									},
								},
								Childs: []*NodeItem{
									{
										Node:          "b",
										Terminated:    "a*[b-e].b",
										TermIndex:     -1,
										WildcardItems: wildcards.WildcardItems{P: "b"},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{
					"a*c": -1, "a*c*": -1, "a*b?c": -1, "a*bd?c": -1, "a*{Z,Q}bd?c": -1, "a*[b-e].b": -1,
					"a*c.b": -1, "a.b?d": -1,
				},
			},
			matchPaths: map[string][]string{
				"acbec":   {"a*c", "a*c*", "a*b?c"},
				"abbece":  {"a*c*"},
				"acbdc":   {"a*c", "a*c*", "a*b?c"},
				"acZbdc":  {"a*c", "a*c*", "a*b?c"},
				"acZbdIc": {"a*c", "a*c*", "a*bd?c", "a*{Z,Q}bd?c"},
				"a.bfd":   {"a.b?d"},
				"ac.b":    {"a*c.b", "a*[b-e].b"},
				"ae.b":    {"a*[b-e].b"},
				"aSTc.b":  {"a*c.b", "a*[b-e].b"},
			},
			missPaths: []string{"", "ab", "c", "a.b", "a.bd", "aa.b", "af.b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
