package gtags

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

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
										Glob: &WildcardItems{
											MinSize: 1, MaxSize: -1, P: "c", Inners: []items.InnerItem{items.ItemStar{}},
										},
									},
									Terminated: []string{"seriesByTag('name=a', 'b!=c*')"},
								},
							},
						},
					},
				},
				Queries: map[string]bool{"seriesByTag('name=a', 'b!=c*')": true},
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
				Queries: map[string]bool{"seriesByTag('name=a', 'b!=c[a]')": true},
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
		t.Run(tt.name, func(t *testing.T) {
			runTestTagsMatcher(t, tt)
		})
	}
}
