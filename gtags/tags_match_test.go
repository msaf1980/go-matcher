package gtags

import (
	"regexp"
	"testing"
)

func TestTagsMatcher_Regex_Match(t *testing.T) {
	tests := []testTagsMatcher{
		{
			name: `{"seriesByTag('name=a', 'b=~c(a|z)\.a')"}`, queries: []string{`seriesByTag('name=a', 'b=~c(a|z)\.a')`},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermMatch, Value: `c(a|z)\.a`,
										Re: regexp.MustCompile(`c(a|z)\.a`),
									},
									Terminated: []string{`seriesByTag('name=a', 'b=~c(a|z)\.a')`},
								},
							},
						},
					},
				},
				Queries: map[string]bool{`seriesByTag('name=a', 'b=~c(a|z)\.a')`: true},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=ca.a":      {`seriesByTag('name=a', 'b=~c(a|z)\.a')`},
				"a?b=ca.a":           {`seriesByTag('name=a', 'b=~c(a|z)\.a')`},
				"a?a=v1&b=cz.a&e=v3": {`seriesByTag('name=a', 'b=~c(a|z)\.a')`},
				"a?a=v1&b=ca.a&e=v3": {`seriesByTag('name=a', 'b=~c(a|z)\.a')`},
			},
			missPaths: []string{"a?a=v1&b=ca.b", "a?b=da", "a?b=v1", "a?c=v1", "b?a=v1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestTagsMatcher(t, tt)
		})
	}
}

var (
	queryEqualR = "seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=~c(r|m)')"
	pathEqualR  = "cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=sales&subproject=crm"
	regexEqualR = `^cpu\.load_avg\?(.*&)?app=postgresql(.*&)?project=sales(.*&)?subproject=c(r|m)`
)

func BenchmarkEqualR_ByTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTagsMatcher()
		err := w.Add(queryEqualR)
		if err != nil {
			b.Fatal(err)
		}
		tags, err := PathTags(pathEqualR)
		if err != nil {
			b.Fatal(err)
		}

		queries := w.MatchByTags(tags)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func BenchmarkEqualR_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := regexp.MustCompile(regexEqualR)
		if !w.MatchString(pathEqualR) {
			b.Fatal(pathEqualR)
		}
	}
}

func BenchmarkEqualR_Precompiled_ByTags(b *testing.B) {
	w := NewTagsMatcher()
	err := w.Add(queryEqualR)
	if err != nil {
		b.Fatal(err)
	}
	queries := make([]string, 0, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tags, err := PathTags(pathEqualR)
		if err != nil {
			b.Fatal(err)
		}
		w.MatchByTagsB(tags, &queries)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func BenchmarkEqualR_Precompiled_ByTags2(b *testing.B) {
	w := NewTagsMatcher()
	err := w.Add(queryEqualR)
	if err != nil {
		b.Fatal(err)
	}
	queries := make([]string, 0, 1)
	tags, err := PathTags(pathEqualR)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.MatchByTagsB(tags, &queries)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func BenchmarkEqualR_Precompiled_Regex(b *testing.B) {
	w := regexp.MustCompile(regexEqualR)
	for i := 0; i < b.N; i++ {
		if !w.MatchString(pathEqualR) {
			b.Fatal(pathEqualR)
		}
	}
}
