package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/globs"
	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_Multi(t *testing.T) {
	tests := []testGlobMatcher{
		// composite
		{
			name:  `{"a*c", "a*c*", "a*b?c", "a*bd?c", "a*{Z,Q}bd?c", "a.b?d", "a*c.b", "a*[b-e].b"}`,
			globs: []string{"a*c", "a*c*", "a*b?c", "a*bd?c", "a*{Z,Q}bd?c", "a.b?d", "a*c.b", "a*[b-e].b"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a*c", Terminated: []string{"a*c"},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
									Inners: []items.Item{items.ItemStar{}},
								},
							},
							{
								Node: "a*c*", Terminated: []string{"a*c*"},
								NodeItem: items.NodeItem{
									P: "a", MinSize: 2, MaxSize: -1,
									Inners: []items.Item{
										items.ItemStar{}, items.ItemRune('c'), items.ItemStar{},
									},
								},
							},
							{
								Node: "a*b?c", Terminated: []string{"a*b?c"},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "c", MinSize: 4, MaxSize: -1,
									Inners: []items.Item{
										items.ItemStar{}, items.ItemRune('b'), items.ItemOne{},
									},
								},
							},
							{
								Node: "a*bd?c", Terminated: []string{"a*bd?c"},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "c", MinSize: 5, MaxSize: -1,
									Inners: []items.Item{
										items.ItemStar{}, items.ItemString("bd"), items.ItemOne{},
									},
								},
							},
							{
								Node: "a*{Q,Z}bd?c", Terminated: []string{"a*{Z,Q}bd?c", "a*{Q,Z}bd?c"},
								NodeItem: items.NodeItem{
									MinSize: 6, MaxSize: -1, P: "a", Suffix: "c",
									Inners: []items.Item{
										items.ItemStar{},
										&items.ItemList{Vals: []string{"Q", "Z"}, ValsMin: 1, ValsMax: 1},
										items.ItemString("bd"), items.ItemOne{},
									},
								},
							},
						},
					},
					2: {
						Childs: []*globs.NodeItem{
							{
								Node: "a", NodeItem: items.NodeItem{P: "a"},
								Childs: []*globs.NodeItem{
									{
										Node: "b?d", Terminated: []string{"a.b?d"},
										NodeItem: items.NodeItem{
											P: "b", Suffix: "d", MinSize: 3, MaxSize: 3,
											Inners: []items.Item{items.ItemOne{}},
										},
									},
								},
							},
							{
								Node: "a*c", NodeItem: items.NodeItem{
									P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
									Inners: []items.Item{items.ItemStar{}},
								},
								Childs: []*globs.NodeItem{
									{
										Node: "b", Terminated: []string{"a*c.b"},
										NodeItem: items.NodeItem{P: "b"},
									},
								},
							},
							{
								Node: "a*[b-e]",
								NodeItem: items.NodeItem{
									P: "a", MinSize: 2, MaxSize: -1,
									Inners: []items.Item{
										items.ItemStar{}, items.ItemRuneRanges{{'b', 'e'}},
									},
								},
								Childs: []*globs.NodeItem{
									{
										Node:       "b",
										Terminated: []string{"a*[b-e].b"},
										NodeItem:   items.NodeItem{P: "b"},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{
					"a*c": -1, "a*c*": -1, "a*b?c": -1, "a*bd?c": -1,
					"a*{Z,Q}bd?c": -1, "a*{Q,Z}bd?c": -1,
					"a*[b-e].b": -1,
					"a*c.b":     -1, "a.b?d": -1,
				},
			},
			matchPaths: map[string][]string{
				"acbec":   {"a*c", "a*c*", "a*b?c"},
				"abbece":  {"a*c*"},
				"acbdc":   {"a*c", "a*c*", "a*b?c"},
				"acZbdc":  {"a*c", "a*c*", "a*b?c"},
				"acZbdIc": {"a*c", "a*c*", "a*bd?c", "a*{Z,Q}bd?c", "a*{Q,Z}bd?c"},
				"a.bfd":   {"a.b?d"},
				"ac.b":    {"a*c.b", "a*[b-e].b"},
				"ae.b":    {"a*[b-e].b"},
				"aSTc.b":  {"a*c.b", "a*[b-e].b"},
			},
			missPaths: []string{"", "ab", "c", "a.b", "a.bd", "aa.b", "af.b"},
		},
	}
	for _, tt := range tests {
		runTestGlobMatcher(t, tt)
	}
}
