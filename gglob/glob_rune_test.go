package gglob

import "testing"

func TestGlobMatcher_Rune(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: `{"[a-c]"}`, globs: []string{"[a-c]"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "[a-c]", Terminated: "[a-c]", MinSize: 1, MaxSize: 1,
								InnerItem: InnerItem{
									Typ: NodeRune, Runes: map[int32]struct{}{'a': {}, 'b': {}, 'c': {}},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"[a-c]": true},
			},
			matchGlobs: map[string][]string{"a": {"[a-c]"}, "c": {"[a-c]"}, "b": {"[a-c]"}},
			miss:       []string{"", "d", "ab", "a.b"},
		},
		{
			name: `{"[a-c]z"}`, globs: []string{"[a-c]z"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "[a-c]z", Terminated: "[a-c]z", MinSize: 2, MaxSize: 2, Suffix: "z",
								InnerItem: InnerItem{
									Typ: NodeRune, Runes: map[int32]struct{}{'a': {}, 'b': {}, 'c': {}},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"[a-c]z": true},
			},
			matchGlobs: map[string][]string{"az": {"[a-c]z"}, "cz": {"[a-c]z"}, "bz": {"[a-c]z"}},
			miss:       []string{"", "d", "ab", "dz", "a.z"},
		},
		{
			name: `{"[a-c]*"}`, globs: []string{"[a-c]*"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "[a-c]*", Terminated: "[a-c]*", MinSize: 1, MaxSize: -1,
								InnerItem: InnerItem{Typ: NodeInners},
								Inners: []*InnerItem{
									{Typ: NodeRune, Runes: map[int32]struct{}{'a': {}, 'b': {}, 'c': {}}},
									{Typ: NodeStar},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"[a-c]*": true},
			},
			matchGlobs: map[string][]string{
				"a": {"[a-c]*"}, "c": {"[a-c]*"},
				"az": {"[a-c]*"}, "cz": {"[a-c]*"}, "bz": {"[a-c]*"},
			},
			miss: []string{"", "d", "dz", "a.z"},
		},
		// one item optimization
		{
			name: `{"[a-]"}`, globs: []string{"[a-]"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "[a-]", Terminated: "[a-]", MinSize: 1, MaxSize: 1,
								InnerItem: InnerItem{Typ: NodeString, P: "a"},
							},
						},
					},
				},
				Globs: map[string]bool{"[a-]": true},
			},
			matchGlobs: map[string][]string{"a": {"[a-]"}},
			miss:       []string{"", "b", "d", "ab", "a.b"},
		},
		{
			name: `{"a[a-]Z"}`, globs: []string{"a[a-]Z"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a[a-]Z", Terminated: "a[a-]Z", MinSize: 3, MaxSize: 3,
								InnerItem: InnerItem{Typ: NodeString, P: "aaZ"},
							},
						},
					},
				},
				Globs: map[string]bool{"a[a-]Z": true},
			},
			matchGlobs: map[string][]string{"aaZ": {"a[a-]Z"}},
			miss:       []string{"", "a", "b", "d", "ab", "aaz", "aaZa", "a.b"},
		},
		{
			name: `{"a[a-]Z[Q]"}`, globs: []string{"a[a-]Z[Q]"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a[a-]Z[Q]", Terminated: "a[a-]Z[Q]", MinSize: 4, MaxSize: 4,
								InnerItem: InnerItem{Typ: NodeString, P: "aaZQ"},
							},
						},
					},
				},
				Globs: map[string]bool{"a[a-]Z[Q]": true},
			},
			matchGlobs: map[string][]string{"aaZQ": {"a[a-]Z[Q]"}},
			miss:       []string{"", "a", "Q", "aaZ", "aaZQa", "a.b"},
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
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "[]a", Terminated: "[]a", MinSize: 1, MaxSize: 1,
								InnerItem: InnerItem{Typ: NodeString, P: "a"},
							},
						},
					},
				},
				Globs: map[string]bool{"[]a": true},
			},
			matchGlobs: map[string][]string{"a": {"[]a"}},
			miss:       []string{"", "b", "ab"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
