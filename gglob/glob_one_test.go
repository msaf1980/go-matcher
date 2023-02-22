package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/globs"
	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_One(t *testing.T) {
	tests := []testGlobMatcher{
		// ? match
		{
			name: `{"?"}`, globs: []string{"?"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "?", Terminated: []string{"?"},
								NodeItem: items.NodeItem{
									MinSize: 1, MaxSize: 1,
									Inners: []items.Item{items.ItemOne{}},
								},
							},
						},
					},
				},
				Globs: map[string]int{"?": -1},
			},
			matchPaths: map[string][]string{"a": {"?"}, "c": {"?"}},
			missPaths:  []string{"", "ab", "a.b"},
		},
		{
			name: `{"a?"}`, globs: []string{"a?"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a?", Terminated: []string{"a?"},
								NodeItem: items.NodeItem{
									P: "a", MinSize: 2, MaxSize: 2,
									Inners: []items.Item{items.ItemOne{}},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a?": -1},
			},
			matchPaths: map[string][]string{"ac": {"a?"}, "az": {"a?"}},
			missPaths:  []string{"", "a", "bc", "ace", "a.c"},
		},
		{
			name: `{"a?c"}`, globs: []string{"a?c"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a?c", Terminated: []string{"a?c"},
								NodeItem: items.NodeItem{
									P: "a", MinSize: 3, MaxSize: 3,
									Inners: []items.Item{items.ItemOne{}}, Suffix: "c",
								},
							},
						},
					},
				},
				Globs: map[string]int{"a?c": -1},
			},
			matchPaths: map[string][]string{"acc": {"a?c"}, "aec": {"a?c"}},
			missPaths:  []string{"", "ab", "ac", "ace", "a.c"},
		},
		{
			name: `{"a?[c]?d"}`, globs: []string{"a?[c]?d"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a?c?d", Terminated: []string{"a?[c]?d", "a?c?d"},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "d", MinSize: 5, MaxSize: 5,
									Inners: []items.Item{
										items.ItemOne{},
										items.ItemRune('c'),
										items.ItemOne{},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a?[c]?d": -1, "a?c?d": -1},
			},
			matchPaths: map[string][]string{"aZccd": {"a?[c]?d", "a?c?d"}, "aZcAd": {"a?[c]?d", "a?c?d"}},
			missPaths:  []string{"", "ab", "ac", "ace", "aZDAd", "a.c"},
		},
		{
			name: `{"a*?c?d"}`, globs: []string{"a*?c?d"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a*?c?d", Terminated: []string{"a*?c?d"},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "d", MinSize: 5, MaxSize: -1,
									Inners: []items.Item{
										items.ItemNStar(1),
										items.ItemRune('c'),
										items.ItemOne{},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a*?c?d": -1},
			},
			matchPaths: map[string][]string{"aZccd": {"a*?c?d"}, "aZcAd": {"a*?c?d"}, "aIZcAd": {"a*?c?d"}},
			missPaths:  []string{"", "ab", "ac", "ace", "aZDAd", "a.c"},
		},
	}
	for _, tt := range tests {
		runTestGlobMatcher(t, tt)
	}
}
