package gtags

import (
	"regexp"
	"testing"
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
			queries: []string{`seriesByTag('name=a', 'b!=~c(a|z)\.a')`},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Childs: []*taggedItemStr{
						{
							Term: "__name__=a",
							Childs: []*taggedItemStr{
								{
									Term: `b!=~c(a|z)\.a`, Terminate: true, TermIndex: 0,
									Terminated: `seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
								},
							},
						},
					},
				},
				Queries: map[string]int{
					`seriesByTag('name=a', 'b!=~c(a|z)\.a')`:    0,
					`seriesByTag('__name__=a','b!=~c(a|z)\.a')`: 0,
				},
				QueryIndex: map[int]string{0: `seriesByTag('__name__=a','b!=~c(a|z)\.a')`},
			},
			match: map[string][]string{
				"a?a=v1&b=ca.b":      {`seriesByTag('__name__=a','b!=~c(a|z)\.a')`},
				"a?b=ca.b":           {`seriesByTag('__name__=a','b!=~c(a|z)\.a')`},
				"a?a=v1&b=c.a&e=v3":  {`seriesByTag('__name__=a','b!=~c(a|z)\.a')`},
				"a?a=v1&b=ca.z&e=v3": {`seriesByTag('__name__=a','b!=~c(a|z)\.a')`},

				"a?a=v1&b=ca.a": {}, "b?a=v1": {},
			},
		},
	}
	for n, tt := range tests {
		runTestGTagsTree(t, n, tt)
	}
}
