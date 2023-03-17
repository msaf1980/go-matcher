package gtags

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestTaggedTermListNe(t *testing.T) {
	tests := []testTaggedTermList{
		{
			query:     "seriesByTag('name=a', 'b=c', 'c!=vc')",
			wantQuery: "seriesByTag('__name__=a','b=c','c!=vc')",
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

func TestGTagsTree_Equal_Ne(t *testing.T) {
	tests := []testGTagsTree{
		{
			queries: []string{"seriesByTag('name=a', 'b=c', 'c!=vc')"},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Items: []taggedItemsStr{
						{
							Key: "__name__",
							Matched: []*taggedItemStr{
								{
									Term: "__name__=a",
									Items: []taggedItemsStr{
										{
											Key: "b",
											Matched: []*taggedItemStr{
												{
													Term: "b=c", Items: []taggedItemsStr{
														{
															Key: "c",
															NotMatched: []*taggedItemStr{
																{
																	Term: "c!=vc",
																	Terminated: items.Terminated{
																		Terminate: true,
																		Query:     "seriesByTag('__name__=a','b=c','c!=vc')",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a', 'b=c', 'c!=vc')":   0,
					"seriesByTag('__name__=a','b=c','c!=vc')": 0,
				},
				QueryIndex: map[int]string{
					0: "seriesByTag('__name__=a','b=c','c!=vc')",
				},
			},
			match: map[string][]string{
				"a?a=v1&b=c":           {"seriesByTag('__name__=a','b=c','c!=vc')"},
				"a?b=c":                {"seriesByTag('__name__=a','b=c','c!=vc')"},
				"a?a=v1&b=c&e=v3":      {"seriesByTag('__name__=a','b=c','c!=vc')"},
				"a?a=v1&b=c&c=v1&e=v3": {"seriesByTag('__name__=a','b=c','c!=vc')"},

				"a?b=ca": {}, "a?b=v1": {}, "a?c=v1": {}, "b?a=v1": {},
				"a?a=v1&b=c&c=vc&e=v3": {},
			},
		},
	}
	for n, tt := range tests {
		runTestGTagsTree(t, n, tt)
	}
}
