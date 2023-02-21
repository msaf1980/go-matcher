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

func TestTagsMatcher_Regex_Match_Ne(t *testing.T) {
	tests := []testTagsMatcher{
		{
			name: `{"seriesByTag('name=a', 'b!=~c(a|z)\.a')"}`, queries: []string{`seriesByTag('name=a', 'b!=~c(a|z)\.a')`},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermNotMatch, Value: `c(a|z)\.a`,
										Re: regexp.MustCompile(`c(a|z)\.a`),
									},
									Terminated: []string{
										`seriesByTag('name=a', 'b!=~c(a|z)\.a')`,
										`seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
									},
								},
							},
						},
					},
				},
				Queries: map[string]int{
					`seriesByTag('name=a', 'b!=~c(a|z)\.a')`: -1, `seriesByTag('__name__=a','b!=~c(a|z)\.a')`: -1,
				},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=ca.b": {
					`seriesByTag('name=a', 'b!=~c(a|z)\.a')`, `seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
				},
				"a?b=ca.b": {
					`seriesByTag('name=a', 'b!=~c(a|z)\.a')`, `seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
				},
				"a?a=v1&b=c.a&e=v3": {
					`seriesByTag('name=a', 'b!=~c(a|z)\.a')`, `seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
				},
				"a?a=v1&b=ca.z&e=v3": {
					`seriesByTag('name=a', 'b!=~c(a|z)\.a')`, `seriesByTag('__name__=a','b!=~c(a|z)\.a')`,
				},
			},
			missPaths: []string{"a?a=v1&b=ca.a", "b?a=v1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestTagsMatcher(t, tt)
		})
	}
}
