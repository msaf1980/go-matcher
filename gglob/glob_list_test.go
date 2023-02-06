package gglob

import "testing"

func TestGlobMatcher_List(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: `{"{a,bc}"}`, globs: []string{"{a,bc}"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "{a,bc}", Terminated: "{a,bc}", MinSize: 1, MaxSize: 2,
								InnerItem: InnerItem{Typ: NodeList, Vals: []string{"a", "bc"}, ValsMin: 1, ValsMax: 2},
							},
						},
					},
				},
				Globs: map[string]bool{"{a,bc}": true},
			},
			matchGlobs: map[string][]string{"a": {"{a,bc}"}, "bc": {"{a,bc}"}},
			miss:       []string{"", "b", "ab", "ba", "abc"},
		},
		{
			name: `{"a{a,bc}{qa,q}c"}`, globs: []string{"a{a,bc}{qa,q}c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a{a,bc}{qa,q}c", Terminated: "a{a,bc}{qa,q}c", MinSize: 4, MaxSize: 6,
								InnerItem: InnerItem{Typ: NodeInners, P: "a"}, Suffix: "c",
								Inners: []*InnerItem{
									{Typ: NodeList, Vals: []string{"a", "bc"}, ValsMin: 1, ValsMax: 2},
									{Typ: NodeList, Vals: []string{"q", "qa"}, ValsMin: 1, ValsMax: 2},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"a{a,bc}{qa,q}c": true},
			},
			matchGlobs: map[string][]string{"aaqac": {"a{a,bc}{qa,q}c"}, "abcqac": {"a{a,bc}{qa,q}c"}, "aaqc": {"a{a,bc}{qa,q}c"}},
			miss:       []string{"", "b", "ab", "ba", "abc", "aabc", "aaqbc"},
		},
		{
			name: `{"a{a,bc}Z{qa,q}c"}`, globs: []string{"a{a,bc}Z{qa,q}c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a{a,bc}Z{qa,q}c", Terminated: "a{a,bc}Z{qa,q}c", MinSize: 5, MaxSize: 7,
								InnerItem: InnerItem{Typ: NodeInners, P: "a"}, Suffix: "c",
								Inners: []*InnerItem{
									{Typ: NodeList, Vals: []string{"a", "bc"}, ValsMin: 1, ValsMax: 2},
									{Typ: NodeString, P: "Z"},
									{Typ: NodeList, Vals: []string{"q", "qa"}, ValsMin: 1, ValsMax: 2},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"a{a,bc}Z{qa,q}c": true},
			},
			matchGlobs: map[string][]string{"aaZqac": {"a{a,bc}Z{qa,q}c"}, "abcZqac": {"a{a,bc}Z{qa,q}c"}, "aaZqc": {"a{a,bc}Z{qa,q}c"}},
			miss:       []string{"", "b", "ab", "ba", "abc", "aabc", "aaqbc"},
		},
		// one item optimization
		{
			name: `{"{a}"}`, globs: []string{"{a}"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "{a}", Terminated: "{a}", MinSize: 1, MaxSize: 1,
								InnerItem: InnerItem{Typ: NodeString, P: "a"},
							},
						},
					},
				},
				Globs: map[string]bool{"{a}": true},
			},
			matchGlobs: map[string][]string{"a": {"{a}"}},
			miss:       []string{"", "b", "d", "ab", "a.b"},
		},
		{
			name: `{"{a,}"}`, globs: []string{"{a,}"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "{a,}", Terminated: "{a,}", MinSize: 1, MaxSize: 1,
								InnerItem: InnerItem{Typ: NodeString, P: "a"},
							},
						},
					},
				},
				Globs: map[string]bool{"{a,}": true},
			},
			matchGlobs: map[string][]string{"a": {"{a,}"}},
			miss:       []string{"", "b", "d", "ab", "a.b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

func TestGlobMatcher_List_Broken(t *testing.T) {
	tests := []testGlobMatcher{
		// broken
		{name: `{"z{ac"}`, globs: []string{"{ac"}, wantErr: true},
		{name: `{"a}c"}`, globs: []string{"a}c"}, wantErr: true},
		// skip empty
		{
			name: `{"{}a"}`, globs: []string{"{}a"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "{}a", Terminated: "{}a", MinSize: 1, MaxSize: 1,
								InnerItem: InnerItem{Typ: NodeString, P: "a"},
							},
						},
					},
				},
				Globs: map[string]bool{"{}a": true},
			},
			matchGlobs: map[string][]string{"a": {"{}a"}},
			miss:       []string{"", "b", "ab"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
