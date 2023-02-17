package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_Multi(t *testing.T) {
	tests := []testGlobMatcher{
		// composite
		{
			name:  `{"a*c", "a*c*", "a*b?c", "a*bd?c", "a*{Z,Q}bd?c", "a.b?d", "a*c.b"}`,
			globs: []string{"a*c", "a*c*", "a*b?c", "a*bd?c", "a*{Z,Q}bd?c", "a.b?d", "a*c.b"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a*c", Terminated: "a*c", TermIndex: -1, P: "a", Suffix: "c",
								Inners:  []items.InnerItem{items.ItemStar{}},
								MinSize: 2, MaxSize: -1,
							},
							{
								Node: "a*c*", Terminated: "a*c*", TermIndex: -1, P: "a",
								MinSize: 2, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}, items.ItemRune('c'), items.ItemStar{}},
							},
							{
								Node: "a*b?c", Terminated: "a*b?c", TermIndex: -1, P: "a", Suffix: "c",
								MinSize: 4, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}, items.ItemRune('b'), items.ItemOne{}},
							},
							{
								Node: "a*bd?c", Terminated: "a*bd?c", TermIndex: -1, P: "a", Suffix: "c",
								MinSize: 5, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}, items.ItemString("bd"), items.ItemOne{}},
							},
							{
								Node: "a*{Z,Q}bd?c", Terminated: "a*{Z,Q}bd?c", TermIndex: -1, P: "a", Suffix: "c",
								MinSize: 6, MaxSize: -1,
								Inners: []items.InnerItem{
									items.ItemStar{},
									&items.ItemList{
										Vals: []string{"Q", "Z"}, ValsMin: 1, ValsMax: 1,
										FirstRunes: map[int32]struct{}{'Q': {}, 'Z': {}},
									},
									items.ItemString("bd"), items.ItemOne{}},
							},
						},
					},
					2: {
						Childs: []*items.NodeItem{
							{
								Node: "a", P: "a",
								Childs: []*items.NodeItem{
									{
										Node: "b?d", Terminated: "a.b?d", TermIndex: -1, P: "b", Suffix: "d", MinSize: 3, MaxSize: 3,
										Inners: []items.InnerItem{items.ItemOne{}},
									},
								},
							},
							{
								Node: "a*c", P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}},
								Childs: []*items.NodeItem{
									{Node: "b", Terminated: "a*c.b", TermIndex: -1, P: "b"},
								},
							},
						},
					},
				},
				Globs: map[string]int{
					"a*c": -1, "a*c*": -1, "a*b?c": -1, "a*bd?c": -1, "a*{Z,Q}bd?c": -1,
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
