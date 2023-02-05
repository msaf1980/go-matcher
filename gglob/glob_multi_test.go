package gglob

import "testing"

func TestGlobMatcher_Multi(t *testing.T) {
	tests := []testGlobMatcher{
		// composite
		{
			name: `{"a*c", "a*c*", "a*b?c", "a.b?d", "a*c.b"}`, globs: []string{"a*c", "a*c*", "a*b?c", "a.b?d", "a*c.b"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a*c", Terminated: "a*c", InnerItem: InnerItem{Typ: NodeStar, P: "a"}, Suffix: "c",
								MinSize: 2, MaxSize: -1,
							},
							{
								Node: "a*c*", Terminated: "a*c*", InnerItem: InnerItem{Typ: NodeInners, P: "a"},
								MinSize: 2, MaxSize: -1,
								Inners: []*InnerItem{
									{Typ: NodeStar},
									{Typ: NodeString, P: "c"},
									{Typ: NodeStar},
								},
							},
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
					2: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a", InnerItem: InnerItem{Typ: NodeString, P: "a"},
								Childs: []*NodeItem{
									{Node: "b?d", Terminated: "a.b?d", InnerItem: InnerItem{Typ: NodeOne, P: "b"}, Suffix: "d", MinSize: 3, MaxSize: 3},
								},
							},
							{
								Node: "a*c", InnerItem: InnerItem{Typ: NodeStar, P: "a"}, Suffix: "c", MinSize: 2, MaxSize: -1,
								Childs: []*NodeItem{
									{Node: "b", Terminated: "a*c.b", InnerItem: InnerItem{Typ: NodeString, P: "b"}, MinSize: 0, MaxSize: 0},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"a*c": true, "a*c*": true, "a*b?c": true, "a*c.b": true, "a.b?d": true},
			},
			matchGlobs: map[string][]string{
				"acbec":  {"a*c", "a*c*", "a*b?c"},
				"abbece": {"a*c*"},
				"a.bfd":  {"a.b?d"},
			},
			miss: []string{"", "ab", "c", "a.b", "a.bd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}
