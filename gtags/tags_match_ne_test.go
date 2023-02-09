package gtags

import (
	"regexp"
	"testing"
)

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
									Terminated: []string{`seriesByTag('name=a', 'b!=~c(a|z)\.a')`},
								},
							},
						},
					},
				},
				Queries: map[string]bool{`seriesByTag('name=a', 'b!=~c(a|z)\.a')`: true},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=ca.b":      {`seriesByTag('name=a', 'b!=~c(a|z)\.a')`},
				"a?b=ca.b":           {`seriesByTag('name=a', 'b!=~c(a|z)\.a')`},
				"a?a=v1&b=c.a&e=v3":  {`seriesByTag('name=a', 'b!=~c(a|z)\.a')`},
				"a?a=v1&b=ca.z&e=v3": {`seriesByTag('name=a', 'b!=~c(a|z)\.a')`},
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
