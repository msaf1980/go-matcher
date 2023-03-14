package gtags

import (
	"regexp"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestTaggedTermList_Equal(t *testing.T) {
	tests := []testTaggedTermList{
		{
			query:   "",
			wantErr: true,
		},
		// incomplete
		{
			query:   "seriesByTag('a')",
			wantErr: true,
		},
		{
			query:   "seriesByTag(' ')",
			wantErr: true,
		},
		// empty
		{
			query:      "seriesByTag('')",
			wantQuery:  "seriesByTag()",
			matchPaths: []string{"a?a=v1&b=c", "a?b=c", "a?a=v1&b=c&e=v3"},
		},
		{
			query:      "seriesByTag()",
			wantQuery:  "seriesByTag()",
			matchPaths: []string{"a?a=v1&b=c", "a?b=c", "a?a=v1&b=c&e=v3"},
		},
		// match
		{
			query:     "seriesByTag('name=a', 'b=c')",
			wantQuery: "seriesByTag('__name__=a','b=c')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a"},
				{Key: "b", Op: TaggedTermEq, Value: "c"},
			},
			matchPaths: []string{"a?a=v1&b=c", "a?b=c", "a?a=v1&b=c&e=v3"},
			missPaths:  []string{"a?b=ca", "a?b=v1", "a?c=v1", "b?a=v1"},
		},
		{
			query:     "seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=crm')",
			wantQuery: "seriesByTag('__name__=cpu.load_avg','app=postgresql','project=sales','subproject=crm')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "cpu.load_avg"},
				{Key: "app", Op: TaggedTermEq, Value: "postgresql"},
				{Key: "project", Op: TaggedTermEq, Value: "sales"},
				{Key: "subproject", Op: TaggedTermEq, Value: "crm"},
			},
			matchPaths: []string{
				"cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=sales&subproject=crm",
			},
			missPaths: []string{
				"cpu.load_avg?app=crm&dc=dc1&host=node1-crm&project=sales&subproject=crm",
				"cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=backoffice&subproject=card",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			runTestTaggedTermList(t, tt)
		})
	}
}

func TestGTagsTree_Equal(t *testing.T) {
	tests := []testGTagsTree{
		{
			queries: []string{},
			want: &gTagsTreeStr{
				Root:       &taggedItemStr{},
				Queries:    map[string]int{},
				QueryIndex: map[int]string{},
			},
		},
		{
			queries: []string{""},
			want: &gTagsTreeStr{
				Root:       &taggedItemStr{},
				Queries:    map[string]int{},
				QueryIndex: map[int]string{},
			},
		},
		// empty
		{
			queries: []string{"seriesByTag( )"},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Childs: []*taggedItemStr{
						{Terminate: true, Terminated: "seriesByTag()"},
					},
				},
				Queries:    map[string]int{"seriesByTag()": 0, "seriesByTag( )": 0},
				QueryIndex: map[int]string{0: "seriesByTag()"},
			},
		},
		// match
		{
			queries: []string{
				"seriesByTag('name=a', 'b=c')",
			},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Childs: []*taggedItemStr{
						{
							Term: "__name__=a",
							Childs: []*taggedItemStr{
								{
									Term: "b=c", Terminate: true, TermIndex: 0,
									Terminated: "seriesByTag('__name__=a','b=c')",
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('__name__=a','b=c')": 0,
					"seriesByTag('name=a', 'b=c')":    0,
				},
				QueryIndex: map[int]string{0: "seriesByTag('__name__=a','b=c')"},
			},
			match: map[string][]string{
				"a?a=v1&b=c":      {"seriesByTag('__name__=a','b=c')"},
				"a?b=c":           {"seriesByTag('__name__=a','b=c')"},
				"a?a=v1&b=c&e=v3": {"seriesByTag('__name__=a','b=c')"},
				"a?b=ca":          {}, "a?b=v1": {}, "a?c=v1": {}, "b?a=v1": {},
			},
		},
		{
			queries: []string{"seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=crm')"},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Childs: []*taggedItemStr{
						{
							Term: "__name__=cpu.load_avg",
							Childs: []*taggedItemStr{
								{
									Term: "app=postgresql",
									Childs: []*taggedItemStr{
										{
											Term: "project=sales",
											Childs: []*taggedItemStr{
												{
													Term:      "subproject=crm",
													Terminate: true, TermIndex: 0,
													Terminated: "seriesByTag('__name__=cpu.load_avg','app=postgresql','project=sales','subproject=crm')",
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
					"seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=crm')":  0,
					"seriesByTag('__name__=cpu.load_avg','app=postgresql','project=sales','subproject=crm')": 0,
				},
				QueryIndex: map[int]string{
					0: "seriesByTag('__name__=cpu.load_avg','app=postgresql','project=sales','subproject=crm')",
				},
			},
			match: map[string][]string{
				"cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=sales&subproject=crm": {
					"seriesByTag('__name__=cpu.load_avg','app=postgresql','project=sales','subproject=crm')",
				},
				"cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=sales&subproject=crm&z=v": {
					"seriesByTag('__name__=cpu.load_avg','app=postgresql','project=sales','subproject=crm')",
				},
				"cpu.load_avg?app=postgresql&b=v&dc=dc1&host=node1-db&project=sales&subproject=crm&z=v": {
					"seriesByTag('__name__=cpu.load_avg','app=postgresql','project=sales','subproject=crm')",
				},
				"cpu.load_avg?app=crm&dc=dc1&host=node1-crm&project=sales&subproject=crm":             {},
				"cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=backoffice&subproject=card": {},
			},
		},
		// duplicate seriesByTag item,
		{
			queries: []string{
				"seriesByTag('name=a', 'b=c')", "seriesByTag('b=c','name=a',)",
			},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Childs: []*taggedItemStr{
						{
							Term: "__name__=a",
							Childs: []*taggedItemStr{
								{
									Term: "b=c", Terminate: true, TermIndex: 0,
									Terminated: "seriesByTag('__name__=a','b=c')",
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('__name__=a','b=c')": 0,
					"seriesByTag('name=a', 'b=c')":    0,
				},
				QueryIndex: map[int]string{0: "seriesByTag('__name__=a','b=c')"},
			},
			match: map[string][]string{
				"a?a=v1&b=c":      {"seriesByTag('__name__=a','b=c')"},
				"a?b=c":           {"seriesByTag('__name__=a','b=c')"},
				"a?a=v1&b=c&e=v3": {"seriesByTag('__name__=a','b=c')"},
				"a?b=ca":          {}, "a?b=v1": {}, "a?c=v1": {}, "b?a=v1": {},
			},
		},
	}
	for n, tt := range tests {
		runTestGTagsTree(t, n, tt)
	}
}

var (
	queryEqual = "seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=crm')"
	pathEqual  = "cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=sales&subproject=crm"
	regexEqual = `^cpu\.load_avg\?(.*&)?app=postgresql(.*&)?project=sales(.*&)?subproject=crm(&|$)`
)

func BenchmarkEqual_Terms(b *testing.B) {
	for i := 0; i < b.N; i++ {
		terms, err := ParseSeriesByTag(queryEqual)
		if err != nil {
			b.Fatal(err)
		}
		tags, err := PathTags(pathEqual)
		if err != nil {
			b.Fatal(err)
		}

		if !terms.MatchByTags(tags) {
			b.Fatal(pathEqual)
		}
	}
}

func BenchmarkEqual_Tree_ByTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		_, _, err := w.Add(queryEqual, 0)
		if err != nil {
			b.Fatal(err)
		}
		tags, err := PathTags(pathEqual)
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

func _BenchmarkEqual_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := regexp.MustCompile(regexEqual)
		if !w.MatchString(pathEqual) {
			b.Fatal(pathEqual)
		}
	}
}

func BenchmarkEqual_Terms_Precompiled(b *testing.B) {
	terms, err := ParseSeriesByTag(queryEqual)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tags, err := PathTags(pathEqual)
		if err != nil {
			b.Fatal(err)
		}

		if !terms.MatchByTags(tags) {
			b.Fatal(pathEqual)
		}
	}
}

func BenchmarkEqual_Terms_Prealloc(b *testing.B) {
	terms, err := ParseSeriesByTag(queryEqual)
	if err != nil {
		b.Fatal(err)
	}

	tags, err := PathTags(pathEqual)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !terms.MatchByTags(tags) {
			b.Fatal(pathEqual)
		}
	}
}

func BenchmarkEqual_Tree_ByTags_Precompiled(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(queryEqual, 0)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tags, err := PathTags(pathEqual)
		if err != nil {
			b.Fatal(err)
		}
		queries := make([]string, 0, 1)
		index := make([]int, 0, 1)
		first := items.MinStore{-1}
		_ = w.MatchByTags(tags, &queries, &index, &first)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func BenchmarkEqual_Tree_ByTags_Prealloc(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(queryEqual, 0)
	if err != nil {
		b.Fatal(err)
	}
	tags, err := PathTags(pathEqual)
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
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func _BenchmarkEqual_Regex_Precompiled(b *testing.B) {
	w := regexp.MustCompile(regexEqual)
	for i := 0; i < b.N; i++ {
		if !w.MatchString(pathEqual) {
			b.Fatal(pathEqual)
		}
	}
}
