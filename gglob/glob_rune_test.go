package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_Rune(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: `{"[a-c]"}`, globs: []string{"[a-c]"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "[a-c]", Terminated: "[a-c]", MinSize: 1, MaxSize: 1,
								Inners: []items.InnerItem{
									items.ItemRune(map[int32]struct{}{'a': {}, 'b': {}, 'c': {}}),
								},
							},
						},
					},
				},
				Globs: map[string]bool{"[a-c]": true},
			},
			matchPaths: map[string][]string{"a": {"[a-c]"}, "c": {"[a-c]"}, "b": {"[a-c]"}},
			missPaths:  []string{"", "d", "ab", "a.b"},
		},
		{
			name: `{"[a-c]z"}`, globs: []string{"[a-c]z"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "[a-c]z", Terminated: "[a-c]z", MinSize: 2, MaxSize: 2, Suffix: "z",
								Inners: []items.InnerItem{
									items.ItemRune(map[int32]struct{}{'a': {}, 'b': {}, 'c': {}}),
								},
							},
						},
					},
				},
				Globs: map[string]bool{"[a-c]z": true},
			},
			matchPaths: map[string][]string{"az": {"[a-c]z"}, "cz": {"[a-c]z"}, "bz": {"[a-c]z"}},
			missPaths:  []string{"", "d", "ab", "dz", "a.z"},
		},
		{
			name: `{"[a-c]*"}`, globs: []string{"[a-c]*"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "[a-c]*", Terminated: "[a-c]*", MinSize: 1, MaxSize: -1,
								Inners: []items.InnerItem{
									items.ItemRune(map[int32]struct{}{'a': {}, 'b': {}, 'c': {}}), items.ItemStar{},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"[a-c]*": true},
			},
			matchPaths: map[string][]string{
				"a": {"[a-c]*"}, "c": {"[a-c]*"},
				"az": {"[a-c]*"}, "cz": {"[a-c]*"}, "bz": {"[a-c]*"},
			},
			missPaths: []string{"", "d", "dz", "a.z"},
		},
		// one item optimization
		{
			name: `{"[a-]"}`, globs: []string{"[a-]"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "[a-]", Terminated: "[a-]", P: "a", MinSize: 1, MaxSize: 1,
							},
						},
					},
				},
				Globs: map[string]bool{"[a-]": true},
			},
			matchPaths: map[string][]string{"a": {"[a-]"}},
			missPaths:  []string{"", "b", "d", "ab", "a.b"},
		},
		{
			name: `{"a[a-]Z"}`, globs: []string{"a[a-]Z"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a[a-]Z", Terminated: "a[a-]Z", P: "aaZ", MinSize: 3, MaxSize: 3,
							},
						},
					},
				},
				Globs: map[string]bool{"a[a-]Z": true},
			},
			matchPaths: map[string][]string{"aaZ": {"a[a-]Z"}},
			missPaths:  []string{"", "a", "b", "d", "ab", "aaz", "aaZa", "a.b"},
		},
		{
			name: `{"a[a-]Z[Q]"}`, globs: []string{"a[a-]Z[Q]"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a[a-]Z[Q]", Terminated: "a[a-]Z[Q]", P: "aaZQ", MinSize: 4, MaxSize: 4,
							},
						},
					},
				},
				Globs: map[string]bool{"a[a-]Z[Q]": true},
			},
			matchPaths: map[string][]string{"aaZQ": {"a[a-]Z[Q]"}},
			missPaths:  []string{"", "a", "Q", "aaZ", "aaZQa", "a.b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

func TestGlobMatcher_Rune_Broken(t *testing.T) {
	tests := []testGlobMatcher{
		// broken
		// compare with graphite-clickhouse. Now It's not error, but filter
		// (Path LIKE 'z%' AND match(Path, '^z[ac$')))
		{name: `{"z[ac"}`, globs: []string{"[ac"}, wantErr: true},
		{name: `{"a]c"}`, globs: []string{"a]c"}, wantErr: true},
		// skip empty
		{
			name: `{"[]a"}`, globs: []string{"[]a"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{Node: "[]a", Terminated: "[]a", P: "a", MinSize: 1, MaxSize: 1},
						},
					},
				},
				Globs: map[string]bool{"[]a": true},
			},
			matchPaths: map[string][]string{"a": {"[]a"}},
			missPaths:  []string{"", "b", "ab"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
