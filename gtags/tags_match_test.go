package gtags

import (
	"regexp"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestTaggedTermList_Regex_Match(t *testing.T) {
	tests := []testTaggedTermList{
		{
			query:     `seriesByTag('name=a', 'b=~c(a|z)\.a')`,
			wantQuery: `seriesByTag('__name__=a','b=~c(a|z)\.a')`,
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a"},
				{Key: "b", Op: TaggedTermMatch, Value: `c(a|z)\.a`, Re: regexp.MustCompile(`c(a|z)\.a`)},
			},
			matchPaths: []string{"a?a=v1&b=ca.a", "a?b=ca.a", "a?a=v1&b=cz.a&e=v3", "a?a=v1&b=ca.a&e=v3"},
			missPaths:  []string{"a?a=v1&b=ca.b", "a?b=da", "a?b=v1", "a?c=v1", "b?a=v1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			runTestTaggedTermList(t, tt)
		})
	}
}

func TestGTagsTree_Regex_Match(t *testing.T) {
	tests := []testGTagsTree{
		{
			queries: []string{`seriesByTag('name=a', 'b=~c(a|z)\.a')`},
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
													Term: `b=~c(a|z)\.a`,
													Terminated: items.Terminated{
														Terminate: true,
														Query:     `seriesByTag('__name__=a','b=~c(a|z)\.a')`,
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
					`seriesByTag('name=a', 'b=~c(a|z)\.a')`:    0,
					`seriesByTag('__name__=a','b=~c(a|z)\.a')`: 0,
				},
				QueryIndex: map[int]string{0: `seriesByTag('__name__=a','b=~c(a|z)\.a')`},
			},
			match: map[string][]string{
				"a?a=v1&b=ca.a":      {`seriesByTag('__name__=a','b=~c(a|z)\.a')`},
				"a?b=ca.a":           {`seriesByTag('__name__=a','b=~c(a|z)\.a')`},
				"a?a=v1&b=cz.a&e=v3": {`seriesByTag('__name__=a','b=~c(a|z)\.a')`},
				"a?a=v1&b=ca.a&e=v3": {`seriesByTag('__name__=a','b=~c(a|z)\.a')`},

				"a?a=v1&b=ca.b": {}, "a?b=da": {}, "a?b=v1": {}, "a?c=v1": {}, "b?a=v1": {},
			},
		},
	}
	for n, tt := range tests {
		runTestGTagsTree(t, n, tt)
	}
}

var (
	queryEqualR = "seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=~c(r|m)')"
	pathEqualR  = "cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=sales&subproject=crm"
	regexEqualR = `^cpu\.load_avg\?(.*&)?app=postgresql(.*&)?project=sales(.*&)?subproject=c(r|m)`
)

func BenchmarkMatch_Terms(b *testing.B) {
	for i := 0; i < b.N; i++ {
		terms, err := ParseSeriesByTag(queryEqual)
		if err != nil {
			b.Fatal(err)
		}
		tags, err := PathTags(pathEqualR)
		if err != nil {
			b.Fatal(err)
		}

		if !terms.MatchByTags(tags) {
			b.Fatal(pathEqualR)
		}
	}
}

func BenchmarkMatch_Tree_ByTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		_, _, err := w.Add(queryEqualR, 0)
		if err != nil {
			b.Fatal(err)
		}
		tags, err := PathTags(pathEqualR)
		if err != nil {
			b.Fatal(err)
		}
		var (
			queries []string
			index   []int
		)
		first := items.MinStore{-1}

		_ = w.MatchByTags(tags, &queries, &index, &first)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func _BenchmarkMatch_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := regexp.MustCompile(regexEqualR)
		if !w.MatchString(pathEqualR) {
			b.Fatal(pathEqualR)
		}
	}
}

func BenchmarkMatch_Terms_Precompiled(b *testing.B) {
	terms, err := ParseSeriesByTag(queryEqual)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tags, err := PathTags(pathEqualR)
		if err != nil {
			b.Fatal(err)
		}

		if !terms.MatchByTags(tags) {
			b.Fatal(pathEqualR)
		}
	}
}

func BenchmarkMatch_Terms_Prealloc(b *testing.B) {
	terms, err := ParseSeriesByTag(queryEqual)
	if err != nil {
		b.Fatal(err)
	}
	tags, err := PathTags(pathEqualR)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !terms.MatchByTags(tags) {
			b.Fatal(pathEqualR)
		}
	}
}

func BenchmarkMatch_Tree_ByTags_Precompiled(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(queryEqualR, 0)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tags, err := PathTags(pathEqualR)
		if err != nil {
			b.Fatal(err)
		}
		var (
			queries []string
			index   []int
		)
		first := items.MinStore{-1}

		_ = w.MatchByTags(tags, &queries, &index, &first)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func BenchmarkMatch_Tree_ByTags_Prealloc(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(queryEqualR, 0)
	if err != nil {
		b.Fatal(err)
	}
	tags, err := PathTags(pathEqualR)
	if err != nil {
		b.Fatal(err)
	}

	queries := make([]string, 0, 1)
	index := make([]int, 0, 1)
	first := items.MinStore{-1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queries = queries[:0]
		index = index[:0]
		first.Init()
		_ = w.MatchByTags(tags, &queries, &index, &first)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func _BenchmarkEqualR_Precompiled_Regex(b *testing.B) {
	w := regexp.MustCompile(regexEqualR)
	for i := 0; i < b.N; i++ {
		if !w.MatchString(pathEqualR) {
			b.Fatal(pathEqualR)
		}
	}
}
