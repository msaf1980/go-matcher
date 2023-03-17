package gtags

import (
	"regexp"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestTaggedTermList_Regex_Match_Ne(t *testing.T) {
	tests := []testTaggedTermList{
		{
			query:     `seriesByTag('name=a', 'b!=~c(a|z)\.a')`,
			wantQuery: `seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a"},
				{Key: "b", Op: TaggedTermNotMatch, Value: `c(a|z)\.a`, Re: regexp.MustCompile(`c(a|z)\.a`)},
			},
			matchPaths: []string{"a?a=v1&b=ca.b", "a?b=ca.b", "a?a=v1&b=c.a&e=v3", "a?a=v1&b=ca.z&e=v3"},
			missPaths:  []string{"a?a=v1&b=ca.a", "b?a=v1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			runTestTaggedTermList(t, tt)
		})
	}
}

func TestGTagsTree_Regex_Match_Ne(t *testing.T) {
	tests := []testGTagsTree{
		{
			queries: []string{
				`seriesByTag('name=a', 'b!=~c(a|z)\.a')`,
				`seriesByTag('__name__=a','b!=~c(a|z)\.a','d=e')`,
			},
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
													Term: `b!=~c(a|z)\.a`,
													Terminated: items.Terminated{
														Terminate: true,
														Query:     `seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
													},
													Items: []taggedItemsStr{
														{
															Key: "d",
															Matched: []*taggedItemStr{
																{
																	Term: "d=e",
																	Terminated: items.Terminated{
																		Terminate: true, Index: 1,
																		Query: `seriesByTag('__name__=a','b!=~c(a|z)\.a','d=e')`,
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
					`seriesByTag('name=a', 'b!=~c(a|z)\.a')`:          0,
					`seriesByTag('__name__=a','b!=~c(a|z)\.a')`:       0,
					`seriesByTag('__name__=a','b!=~c(a|z)\.a','d=e')`: 1,
				},
				QueryIndex: map[int]string{
					0: `seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
					1: `seriesByTag('__name__=a','b!=~c(a|z)\.a','d=e')`,
				},
			},
			match: map[string][]string{
				"a?a=v1&b=ca.b":      {`seriesByTag('__name__=a','b!=~c(a|z)\.a')`},
				"a?b=ca.b":           {`seriesByTag('__name__=a','b!=~c(a|z)\.a')`},
				"a?a=v1&b=c.a&e=v3":  {`seriesByTag('__name__=a','b!=~c(a|z)\.a')`},
				"a?a=v1&b=ca.z&e=v3": {`seriesByTag('__name__=a','b!=~c(a|z)\.a')`},
				"a?a=v1&b=ca.z&d=e&e=v3": {
					`seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
					`seriesByTag('__name__=a','b!=~c(a|z)\.a','d=e')`,
				},
				// tag b not exist
				"a?a=v1&d=ca.b": {`seriesByTag('__name__=a','b!=~c(a|z)\.a')`},
				"a?a=v1&d=e": {
					`seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
					`seriesByTag('__name__=a','b!=~c(a|z)\.a','d=e')`,
				},
				"a?a=v1&d=e&f=g": {
					`seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
					`seriesByTag('__name__=a','b!=~c(a|z)\.a','d=e')`,
				},
				"a?d=e": {
					`seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
					`seriesByTag('__name__=a','b!=~c(a|z)\.a','d=e')`,
				},
				"a?d=e&f=g": {
					`seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
					`seriesByTag('__name__=a','b!=~c(a|z)\.a','d=e')`,
				},

				"a?a=v1&b=ca.a": {}, "a?a=v1&b=ca.a&d=e": {},
				"b?a=v1": {},
			},
		},
	}
	for n, tt := range tests {
		runTestGTagsTree(t, n, tt)
	}
}
