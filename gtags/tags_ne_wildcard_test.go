package gtags

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

func TestTaggedTermListNe_Wildcard(t *testing.T) {
	tests := []testTaggedTermList{
		{
			query: "seriesByTag('name=a', 'b!=c*')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a"},
				{
					Key: "b", Op: TaggedTermNe, Value: "c*", HasWildcard: true,
					Glob: &wildcards.WildcardItems{MinSize: 1, MaxSize: -1, P: "c", Inners: []wildcards.InnerItem{wildcards.ItemStar{}}},
				},
			},
			matchPaths: []string{"a?a=v1&b=ba", "a?c=ca", "a?a=v1&b=b&e=v3", "a?a=v1&b=ba&e=v3"},
			missPaths:  []string{"a?b=c", "a?b=ca", "b?a=v1"},
		},
		// compaction
		{
			query: "seriesByTag('name=a', 'b!=c[a]')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a"},
				{Key: "b", Op: TaggedTermNe, Value: "ca"},
			},
			matchPaths: []string{"a?b=c", "a?a=v1&b=ba", "a?b=ba", "a?a=v1&b=ba&e=v3"},
			missPaths:  []string{"a?b=ca", "b?a=v1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			runTestTaggedTermList(t, tt)
		})
	}
}

func TestTagsMatcherNe_Wildcard(t *testing.T) {
	tests := []testTagsMatcher{
		{
			name: `{"seriesByTag('name=a', 'b!=c*')"}`, queries: []string{"seriesByTag('name=a', 'b!=c*')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermNe, Value: "c*", HasWildcard: true,
										Glob: &wildcards.WildcardItems{
											MinSize: 1, MaxSize: -1, P: "c", Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
										},
									},
									Terminated: []string{"seriesByTag('name=a', 'b!=c*')"},
								},
							},
						},
					},
				},
				Queries: map[string]int{"seriesByTag('name=a', 'b!=c*')": -1},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=ba":      {"seriesByTag('name=a', 'b!=c*')"},
				"a?c=ca":           {"seriesByTag('name=a', 'b!=c*')"},
				"a?a=v1&b=b&e=v3":  {"seriesByTag('name=a', 'b!=c*')"},
				"a?a=v1&b=ba&e=v3": {"seriesByTag('name=a', 'b!=c*')"},
			},
			missPaths: []string{"a?b=c", "a?b=ca", "b?a=v1"},
		},
		// compaction
		{
			name: `{"seriesByTag('name=a', 'b!=c[a]')"}`, queries: []string{"seriesByTag('name=a', 'b!=c[a]')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermNe, Value: "ca",
									},
									Terminated: []string{"seriesByTag('name=a', 'b!=c[a]')"},
								},
							},
						},
					},
				},
				Queries: map[string]int{"seriesByTag('name=a', 'b!=c[a]')": -1},
			},
			matchPaths: map[string][]string{
				"a?b=c":            {"seriesByTag('name=a', 'b!=c[a]')"},
				"a?a=v1&b=ba":      {"seriesByTag('name=a', 'b!=c[a]')"},
				"a?b=ba":           {"seriesByTag('name=a', 'b!=c[a]')"},
				"a?a=v1&b=ba&e=v3": {"seriesByTag('name=a', 'b!=c[a]')"},
			},
			missPaths: []string{"a?b=ca", "b?a=v1"},
		},
	}
	for _, tt := range tests {
		runTestTagsMatcher(t, tt)
	}
}
