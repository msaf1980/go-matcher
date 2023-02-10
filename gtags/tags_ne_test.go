package gtags

import (
	"testing"
)

func TestTaggedTermListNe(t *testing.T) {
	tests := []testTaggedTermList{
		{
			query: "seriesByTag('name=a', 'b=c', 'c!=vc')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a"},
				{Key: "b", Op: TaggedTermEq, Value: "c"},
				{Key: "c", Op: TaggedTermNe, Value: "vc"},
			},
			matchPaths: []string{"a?a=v1&b=c", "a?b=c", "a?a=v1&b=c&e=v3", "a?a=v1&b=c&c=v1&e=v3"},
			missPaths:  []string{"a?b=ca", "a?b=v1", "a?c=v1", "b?a=v1", "a?a=v1&b=c&c=vc&e=v3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			runTestTaggedTermList(t, tt)
		})
	}
}

func TestTagsMatcherNe(t *testing.T) {
	tests := []testTagsMatcher{
		{
			name: `{"seriesByTag('name=a', 'b=c', 'c!=vc')"}`, queries: []string{"seriesByTag('name=a', 'b=c', 'c!=vc')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{Key: "b", Op: TaggedTermEq, Value: "c"},
									Childs: []*TaggedItem{
										{
											Term:       &TaggedTerm{Key: "c", Op: TaggedTermNe, Value: "vc"},
											Terminated: []string{"seriesByTag('name=a', 'b=c', 'c!=vc')"},
										},
									},
								},
							},
						},
					},
				},
				Queries: map[string]bool{"seriesByTag('name=a', 'b=c', 'c!=vc')": true},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=c":           {"seriesByTag('name=a', 'b=c', 'c!=vc')"},
				"a?b=c":                {"seriesByTag('name=a', 'b=c', 'c!=vc')"},
				"a?a=v1&b=c&e=v3":      {"seriesByTag('name=a', 'b=c', 'c!=vc')"},
				"a?a=v1&b=c&c=v1&e=v3": {"seriesByTag('name=a', 'b=c', 'c!=vc')"},
			},
			missPaths: []string{"a?b=ca", "a?b=v1", "a?c=v1", "b?a=v1", "a?a=v1&b=c&c=vc&e=v3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestTagsMatcher(t, tt)
		})
	}
}
