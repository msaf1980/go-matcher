package gglob

import "testing"

func TestGlobMatcher_Star(t *testing.T) {
	tests := []testGlobMatcher{
		// deduplication
		{
			name: `{"a******c"}`, globs: []string{"a******c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a******c", Terminated: "a******c", InnerItem: InnerItem{Typ: NodeStar, P: "a"},
								MinSize: 2, MaxSize: -1, Suffix: "c",
							},
						},
					},
				},
				Globs: map[string]bool{"a******c": true},
			},
			matchGlobs: map[string][]string{"ac": {"a******c"}, "abc": {"a******c"}, "abcc": {"a******c"}},
			miss:       []string{"", "acb"},
		},
		// * match
		{
			name: `{"*"}`, globs: []string{"*"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{Node: "*", Terminated: "*", InnerItem: InnerItem{Typ: NodeStar}, MaxSize: -1},
						},
					},
				},
				Globs: map[string]bool{"*": true},
			},
			matchGlobs: map[string][]string{"a": {"*"}, "b": {"*"}, "ce": {"*"}},
			miss:       []string{"", "b.c"},
		},
		{
			name: `{"a*c"}`, globs: []string{"a*c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{Node: "a*c", Terminated: "a*c", InnerItem: InnerItem{Typ: NodeStar, P: "a"}, Suffix: "c", MinSize: 2, MaxSize: -1},
						},
					},
				},
				Globs: map[string]bool{"a*c": true},
			},
			matchGlobs: map[string][]string{
				"ac": {"a*c"}, "acc": {"a*c"}, "aec": {"a*c"}, "aebc": {"a*c"},
				"aecc": {"a*c"}, "aecec": {"a*c"}, "abecec": {"a*c"},
			},
			miss: []string{"", "ab", "c", "ace", "a.c"},
		},
		// composite
		{
			name: `{"a*b?c"}`, globs: []string{"a*b?c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a*b?c", Terminated: "a*b?c", InnerItem: InnerItem{Typ: NodeInners, P: "a"}, Suffix: "c",
								MinSize: 4, MaxSize: -1,
								Inners: []*InnerItem{
									{Typ: NodeStar},
									{Typ: NodeString, P: "b"},
									{Typ: NodeOne},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"a*b?c": true},
			},
			matchGlobs: map[string][]string{
				"abec":   {"a*b?c"}, // skip *
				"abbec":  {"a*b?c"}, /// shit first b
				"acbbc":  {"a*b?c"},
				"aecbec": {"a*b?c"},
			},
			miss: []string{"", "ab", "c", "ace", "a.c", "abbece"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
