package gtags

import (
	"regexp"
	"testing"

	"github.com/msaf1980/go-matcher/glob"
	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/utils"
)

func TestTaggedTermListEqual_Wildcard(t *testing.T) {
	tests := []testTaggedTermList{
		{
			query:     "seriesByTag('name=a', 'b=c*')",
			wantQuery: "seriesByTag('__name__=a','b=c*')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a"},
				{
					Key: "b", Op: TaggedTermEq, Value: "c*", HasWildcard: true,
					Glob: &glob.Glob{
						Glob: "c*", Node: "c*", MinLen: 1, MaxLen: -1, Prefix: "c",
						Items: []items.Item{items.Star(0)},
					},
				},
			},
			matchPaths: []string{"a?a=v1&b=ca", "a?b=c", "a?a=v1&b=c&e=v3", "a?a=v1&b=ca&e=v3"},
			missPaths:  []string{"a?b=da", "a?b=v1", "a?c=v1", "b?a=v1"},
		},
		{
			query:     "seriesByTag('name=a.b', 'b=c*.a')",
			wantQuery: "seriesByTag('__name__=a.b','b=c*.a')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a.b"},
				{
					Key: "b", Op: TaggedTermEq, Value: "c*.a", HasWildcard: true,
					Glob: &glob.Glob{
						Glob: "c*.a", Node: "c*.a",
						MinLen: 3, MaxLen: -1, Prefix: "c", Suffix: ".a",
						Items: []items.Item{items.Star(0)},
					},
				},
			},
			matchPaths: []string{
				"a.b?a=v1&b=ca.a", "a.b?b=c.a", "a.b?a=v1&b=c.a&e=v3", "a.b?a=v1&b=ca.a&e=v3", "a.b?a=v1&b=cb.a&e=v3",
			},
			missPaths: []string{"a?b=c.a", "a.b?b=da", "a.b?b=ca", "a.b?b=ca.b", "a.b?b=v1", "a.b?c=v1", "b?a=v1"},
		},
		{
			query:     "seriesByTag('name=a.b', 'b=a{a,bc}Z{qa,q}c.a')",
			wantQuery: "seriesByTag('__name__=a.b','b=a{a,bc}Z{q,qa}c.a')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a.b"},
				{
					Key: "b", Op: TaggedTermEq, Value: "a{a,bc}Z{q,qa}c.a", HasWildcard: true,
					Glob: &glob.Glob{
						Glob: "a{a,bc}Z{qa,q}c.a", Node: "a{a,bc}Z{q,qa}c.a",
						Prefix: "a", Suffix: "c.a", MinLen: 7, MaxLen: 9,
						Items: []items.Item{
							&items.StringList{
								Vals: []string{"a", "bc"}, MinSize: 1, MaxSize: 2,
								FirstASCII: utils.MakeASCIISetMust("ab"), ASCIIStarted: true,
							},
							items.Byte('Z'),
							&items.StringList{
								Vals: []string{"q", "qa"}, MinSize: 1, MaxSize: 2,
								FirstASCII: utils.MakeASCIISetMust("q"), ASCIIStarted: true,
							},
						},
					},
				},
			},
			matchPaths: []string{
				"a.b?a=v1&b=aaZqc.a", "a.b?a=v1&b=abcZqac.a",
			},
			missPaths: []string{"a?b=c.a", "a.b?b=da", "a.b?b=ca", "a.b?b=ca.b", "a.b?b=v1", "a.b?c=v1", "b?a=v1"},
		},
		// compaction
		{
			query:     "seriesByTag('name=a', 'b=c[a]')",
			wantQuery: "seriesByTag('__name__=a','b=ca')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a"},
				{Key: "b", Op: TaggedTermEq, Value: "ca"},
			},
			matchPaths: []string{"a?a=v1&b=ca", "a?b=ca", "a?a=v1&b=ca&e=v3"},
			missPaths:  []string{"a?b=c", "a?b=v1", "a?c=v1", "b?a=v1"},
		},
		{
			query:     "seriesByTag('name=a', 'b=a?*??c')",
			wantQuery: "seriesByTag('__name__=a','b=a*???c')",
			want: TaggedTermList{
				{Key: "__name__", Op: TaggedTermEq, Value: "a"},
				{
					Key: "b", Op: TaggedTermEq, Value: "a*???c", HasWildcard: true,
					Glob: &glob.Glob{
						Glob: "a?*??c", Node: "a*???c",
						Prefix: "a", Suffix: "c", MinLen: 5, MaxLen: -1,
						Items: []items.Item{items.Star(3)},
					},
				},
			},
			matchPaths: []string{"a?a=v1&b=aBCDc", "a?b=aAFCDc", "a?a=v1&b=aAFCDc&e=v3"},
			missPaths:  []string{"a?b=c", "a?b=v1", "a?b=aCDc", "a?c=v1", "b?a=v1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			runTestTaggedTermList(t, tt)
		})
	}
}

func TestGTagsTree_Equal_Wildcard(t *testing.T) {
	tests := []testGTagsTree{
		{
			queries: []string{"seriesByTag('name=a', 'b=c*')"},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Childs: []*taggedItemStr{
						{
							Term: "__name__=a",
							Childs: []*taggedItemStr{
								{
									Term: "b=c*", Terminate: true, TermIndex: 0,
									Terminated: "seriesByTag('__name__=a','b=c*')",
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a', 'b=c*')": 0, "seriesByTag('__name__=a','b=c*')": 0,
				},
				QueryIndex: map[int]string{0: "seriesByTag('__name__=a','b=c*')"},
			},
			match: map[string][]string{
				"a?a=v1&b=ca":      {"seriesByTag('__name__=a','b=c*')"},
				"a?b=c":            {"seriesByTag('__name__=a','b=c*')"},
				"a?a=v1&b=c&e=v3":  {"seriesByTag('__name__=a','b=c*')"},
				"a?a=v1&b=ca&e=v3": {"seriesByTag('__name__=a','b=c*')"},
				"a?b=da":           {}, "a?b=v1": {}, "a?c=v1": {}, "b?a=v1": {},
			},
		},
		{
			queries: []string{"seriesByTag('name=a.b', 'b=c*.a')"},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Childs: []*taggedItemStr{
						{
							Term: "__name__=a.b",
							Childs: []*taggedItemStr{
								{
									Term: "b=c*.a", Terminate: true, TermIndex: 0,
									Terminated: "seriesByTag('__name__=a.b','b=c*.a')",
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a.b', 'b=c*.a')":    0,
					"seriesByTag('__name__=a.b','b=c*.a')": 0,
				},
				QueryIndex: map[int]string{0: "seriesByTag('__name__=a.b','b=c*.a')"},
			},
			match: map[string][]string{
				"a.b?a=v1&b=ca.a":      {"seriesByTag('__name__=a.b','b=c*.a')"},
				"a.b?b=c.a":            {"seriesByTag('__name__=a.b','b=c*.a')"},
				"a.b?a=v1&b=c.a&e=v3":  {"seriesByTag('__name__=a.b','b=c*.a')"},
				"a.b?a=v1&b=ca.a&e=v3": {"seriesByTag('__name__=a.b','b=c*.a')"},
				"a.b?a=v1&b=cb.a&e=v3": {"seriesByTag('__name__=a.b','b=c*.a')"},

				"a?b=c.a": {}, "a.b?b=da": {}, "a.b?b=ca": {},
				"a.b?b=ca.b": {}, "a.b?b=v1": {},
				"a.b?c=v1": {}, "b?a=v1": {},
			},
		},
		{
			queries: []string{"seriesByTag('name=a.b', 'b=a{a,bc}Z{qa,q}c.a')"},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Childs: []*taggedItemStr{
						{
							Term: "__name__=a.b",
							Childs: []*taggedItemStr{
								{
									Term:      "b=a{a,bc}Z{q,qa}c.a",
									Terminate: true, TermIndex: 0,
									Terminated: "seriesByTag('__name__=a.b','b=a{a,bc}Z{q,qa}c.a')",
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a.b', 'b=a{a,bc}Z{qa,q}c.a')":    0,
					"seriesByTag('__name__=a.b','b=a{a,bc}Z{q,qa}c.a')": 0,
				},
				QueryIndex: map[int]string{0: "seriesByTag('__name__=a.b','b=a{a,bc}Z{q,qa}c.a')"},
			},
			match: map[string][]string{
				"a.b?a=v1&b=aaZqc.a":   {"seriesByTag('__name__=a.b','b=a{a,bc}Z{q,qa}c.a')"},
				"a.b?a=v1&b=abcZqac.a": {"seriesByTag('__name__=a.b','b=a{a,bc}Z{q,qa}c.a')"},

				"a?b=c.a": {}, "a.b?b=da": {}, "a.b?b=ca": {}, "a.b?b=ca.b": {},
				"a.b?b=v1": {}, "a.b?c=v1": {}, "b?a=v1": {},
			},
		},
		// compaction
		{
			queries: []string{"seriesByTag('name=a', 'b=c[a]')"},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Childs: []*taggedItemStr{
						{
							Term: "__name__=a",
							Childs: []*taggedItemStr{
								{
									Term: "b=ca", Terminate: true, TermIndex: 0,
									Terminated: "seriesByTag('__name__=a','b=ca')",
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a', 'b=c[a]')":  0,
					"seriesByTag('__name__=a','b=ca')": 0,
				},
				QueryIndex: map[int]string{0: "seriesByTag('__name__=a','b=ca')"},
			},
			match: map[string][]string{
				"a?a=v1&b=ca":      {"seriesByTag('__name__=a','b=ca')"},
				"a?b=ca":           {"seriesByTag('__name__=a','b=ca')"},
				"a?a=v1&b=ca&e=v3": {"seriesByTag('__name__=a','b=ca')"},

				"a?b=c": {}, "a?b=v1": {}, "a?c=v1": {}, "b?a=v1": {},
			},
		},
		{
			queries: []string{
				"seriesByTag('name=a', 'b=a?*??c')",
				"seriesByTag('name=a', 'b=a*???c')",
			},
			want: &gTagsTreeStr{
				Root: &taggedItemStr{
					Childs: []*taggedItemStr{
						{
							Term: "__name__=a",
							Childs: []*taggedItemStr{
								{
									Term: "b=a*???c", Terminate: true, TermIndex: 0,
									Terminated: "seriesByTag('__name__=a','b=a*???c')",
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a', 'b=a?*??c')":    0,
					"seriesByTag('__name__=a','b=a*???c')": 0,
				},
				QueryIndex: map[int]string{0: "seriesByTag('__name__=a','b=a*???c')"},
			},
			match: map[string][]string{
				"a?a=v1&b=aBCDc":       {"seriesByTag('__name__=a','b=a*???c')"},
				"a?b=aAFCDc":           {"seriesByTag('__name__=a','b=a*???c')"},
				"a?a=v1&b=aAFCDc&e=v3": {"seriesByTag('__name__=a','b=a*???c')"},

				"a?b=c": {}, "a?b=v1": {}, "a?b=aCDc": {}, "a?c=v1": {}, "b?a=v1": {},
			},
		},
	}
	for n, tt := range tests {
		runTestGTagsTree(t, n, tt)
	}
}

var (
	queryEqualW = "seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=c*')"
	pathEqualW  = "cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=sales&subproject=crm"
	regexEqualW = `^cpu\.load_avg\?(.*&)?app=postgresql(.*&)?project=sales(.*&)?subproject=c`
)

func BenchmarkEqualW_Terms(b *testing.B) {
	for i := 0; i < b.N; i++ {
		terms, err := ParseSeriesByTag(queryEqualW)
		if err != nil {
			b.Fatal(err)
		}

		tags, err := PathTags(pathEqualW)
		if err != nil {
			b.Fatal(err)
		}

		if !terms.MatchByTags(tags) {
			b.Fatal(pathEqualW)
		}
	}
}

func BenchmarkEqualW_Tree_ByTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		_, _, err := w.Add(queryEqualW, 0)
		if err != nil {
			b.Fatal(err)
		}
		tags, err := PathTags(pathEqualW)
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

func _BenchmarkEqualW_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := regexp.MustCompile(regexEqualW)
		if !w.MatchString(pathEqualW) {
			b.Fatal(pathEqualW)
		}
	}
}

func BenchmarkEqualW_Terms_Precompiled(b *testing.B) {
	terms, err := ParseSeriesByTag(queryEqualW)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tags, err := PathTags(pathEqualW)
		if err != nil {
			b.Fatal(err)
		}

		if !terms.MatchByTags(tags) {
			b.Fatal(pathEqualW)
		}
	}
}

func BenchmarkEqualW_Terms_Prealloc(b *testing.B) {
	terms, err := ParseSeriesByTag(queryEqual)
	if err != nil {
		b.Fatal(err)
	}

	tags, err := PathTags(pathEqualW)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !terms.MatchByTags(tags) {
			b.Fatal(pathEqualW)
		}
	}
}

func BenchmarkEqualW_Tree_ByTags_Precompiled(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(queryEqualW, 0)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tags, err := PathTags(pathEqualW)
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

func BenchmarkEqualW_Tree_ByTags_Prealloc(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(queryEqualW, 0)
	if err != nil {
		b.Fatal(err)
	}
	tags, err := PathTags(pathEqualW)
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

func _BenchmarkEqualW_Precompiled_Regex(b *testing.B) {
	w := regexp.MustCompile(regexEqualW)
	for i := 0; i < b.N; i++ {
		if !w.MatchString(pathEqualW) {
			b.Fatal(pathEqualW)
		}
	}
}
