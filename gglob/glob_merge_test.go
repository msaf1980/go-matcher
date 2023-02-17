package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestNodeItem_Merge(t *testing.T) {
	tests := []testGlobMatcher{
		// {
		// 	name: "merge strings #all", globs: []string{"a[a-]Z[Q]"},
		// 	wantW: &GlobMatcher{
		// 		Root: map[int]*items.NodeItem{
		// 			1: {
		// 				Childs: []*items.NodeItem{
		// 					{
		// 						Node: "a[a-]Z[Q]", Terminated: "a[a-]Z[Q]", TermIndex: -1,
		// 						P: "aaZQ", MinSize: 4, MaxSize: 4,
		// 					},
		// 				},
		// 			},
		// 		},
		// 		Globs: map[string]int{"a[a-]Z[Q]": -1},
		// 	},
		// 	matchPaths: map[string][]string{"aaZQ": {"a[a-]Z[Q]"}},
		// 	missPaths:  []string{"", "ab", "aaZQa"},
		// },
		// {
		// 	name: "merge strings #prefix", globs: []string{"a[a-]Z[Q]*"},
		// 	wantW: &GlobMatcher{
		// 		Root: map[int]*items.NodeItem{
		// 			1: {
		// 				Childs: []*items.NodeItem{
		// 					{
		// 						Node: "a[a-]Z[Q]*", Terminated: "a[a-]Z[Q]*", TermIndex: -1,
		// 						P: "aaZQ", MinSize: 4, MaxSize: -1,
		// 						Inners: []items.InnerItem{items.ItemStar{}},
		// 					},
		// 				},
		// 			},
		// 		},
		// 		Globs: map[string]int{"a[a-]Z[Q]*": -1},
		// 	},
		// 	matchPaths: map[string][]string{
		// 		"aaZQ":  {"a[a-]Z[Q]*"},
		// 		"aaZQa": {"a[a-]Z[Q]*"},
		// 	},
		// 	missPaths: []string{"", "ab", "aaZqa"},
		// },
		{
			name: "merge strings #suffix", globs: []string{"a[a-]Z[Q]*[z-][a]l"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a[a-]Z[Q]*[z-][a]l", Terminated: "a[a-]Z[Q]*[z-][a]l", TermIndex: -1,
								P: "aaZQ", Suffix: "zal", MinSize: 7, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}},
							},
						},
					},
				},
				Globs: map[string]int{"a[a-]Z[Q]*[z-][a]l": -1},
			},
			matchPaths: map[string][]string{
				"aaZQzal":  {"a[a-]Z[Q]*[z-][a]l"},
				"aaZQazal": {"a[a-]Z[Q]*[z-][a]l"},
			},
			missPaths: []string{"", "ab", "aaZqa"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

func BenchmarkMergeAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		err := w.Add("a[a-]Z[Q]")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMergePrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		err := w.Add("a[a-]Z[Q]*")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMergeSuffix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		err := w.Add("a[a-]Z[Q]*[z-]l")
		if err != nil {
			b.Fatal(err)
		}
	}
}
