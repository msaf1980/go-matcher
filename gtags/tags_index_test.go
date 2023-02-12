package gtags

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestTagsMatcherIndex(t *testing.T) {
	tests := []testTagsMatcherIndex{
		{
			name: `{"seriesByTag('name=a', 'b=c*')", seriesByTag('name=a', 'b=*a')"}`,
			queries: []string{
				"seriesByTag('name=a', 'b=c*')",
				"seriesByTag('name=a', 'b=*a')",
			},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermEq, Value: "c*", HasWildcard: true,
										Glob: &WildcardItems{
											MinSize: 1, MaxSize: -1, P: "c", Inners: []items.InnerItem{items.ItemStar{}},
										},
									},
									Terminated: []string{"seriesByTag('name=a', 'b=c*')"},
									TermIndex:  []int{0},
								},
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermEq, Value: "*a", HasWildcard: true,
										Glob: &WildcardItems{
											MinSize: 1, MaxSize: -1, Suffix: "a",
											Inners: []items.InnerItem{items.ItemStar{}},
										},
									},
									Terminated: []string{"seriesByTag('name=a', 'b=*a')"},
									TermIndex:  []int{1},
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a', 'b=c*')": 0,
					"seriesByTag('name=a', 'b=*a')": 1,
				},
			},
			matchPaths: map[string][]int{
				"a?a=v1&b=ca":      {0, 1},
				"a?b=c":            {0},
				"a?a=v1&b=c&e=v3":  {0},
				"a?a=v1&b=ca&e=v3": {0, 1},
				"a?b=ba":           {1},
				"a?b=dac":          {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestTagsMatcherIndex(t, tt)
		})
	}
}