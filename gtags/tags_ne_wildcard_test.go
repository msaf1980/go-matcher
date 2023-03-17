package gtags

import (
	"testing"

	"github.com/msaf1980/go-matcher/glob"
	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestTaggedTermListNe_Wildcard(t *testing.T) {
	tests := []testTaggedTermList{
		{
			query:     "seriesByTag('name=a', 'b!=c*')",
			wantQuery: "seriesByTag('__name__=a','b!=c*')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a"},
				{
					Key: "b", Op: TaggedTermNe, Value: "c*", HasWildcard: true,
					Glob: &glob.Glob{
						Glob: "c*", Node: "c*",
						MinLen: 1, MaxLen: -1, Prefix: "c",
						Items: []items.Item{items.Star(0)},
					},
				},
			},
			matchPaths: []string{"a?a=v1&b=ba", "a?c=ca", "a?a=v1&b=b&e=v3", "a?a=v1&b=ba&e=v3"},
			missPaths:  []string{"a?b=c", "a?b=ca", "b?a=v1"},
		},
		// compaction
		{
			query:     "seriesByTag('name=a', 'b!=c[a]')",
			wantQuery: "seriesByTag('__name__=a','b!=ca')",
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
	tests := []testGTagsTree{
		{
			queries: []string{"seriesByTag('name=a', 'b!=c*')"},
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
											NotMatched: []*taggedItemStr{
												{
													Term: "b!=c*",
													Terminated: items.Terminated{
														Terminate: true,
														Query:     "seriesByTag('__name__=a','b!=c*')",
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
					"seriesByTag('name=a', 'b!=c*')":    0,
					"seriesByTag('__name__=a','b!=c*')": 0,
				},
				QueryIndex: map[int]string{0: "seriesByTag('__name__=a','b!=c*')"},
			},
			match: map[string][]string{
				"a?a=v1&b=ba":      {"seriesByTag('__name__=a','b!=c*')"},
				"a?c=ca":           {"seriesByTag('__name__=a','b!=c*')"},
				"a?a=v1&b=b&e=v3":  {"seriesByTag('__name__=a','b!=c*')"},
				"a?a=v1&b=ba&e=v3": {"seriesByTag('__name__=a','b!=c*')"},

				"a?b=c": {}, "a?b=ca": {}, "b?a=v1": {},
			},
		},
		// compaction
		{
			queries: []string{"seriesByTag('name=a', 'b!=c[a]')"},
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
											NotMatched: []*taggedItemStr{
												{
													Term: "b!=ca",
													Terminated: items.Terminated{
														Terminate: true,
														Query:     "seriesByTag('__name__=a','b!=ca')",
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
					"seriesByTag('name=a', 'b!=c[a]')":  0,
					"seriesByTag('__name__=a','b!=ca')": 0,
				},
				QueryIndex: map[int]string{0: "seriesByTag('__name__=a','b!=ca')"},
			},
			match: map[string][]string{
				"a?b=c":            {"seriesByTag('__name__=a','b!=ca')"},
				"a?a=v1&b=ba":      {"seriesByTag('__name__=a','b!=ca')"},
				"a?b=ba":           {"seriesByTag('__name__=a','b!=ca')"},
				"a?a=v1&b=ba&e=v3": {"seriesByTag('__name__=a','b!=ca')"},

				"a?b=ca": {}, "b?a=v1": {},
			},
		},
	}
	for n, tt := range tests {
		runTestGTagsTree(t, n, tt)
	}
}
