package gtags

import (
	"regexp"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestTagsMatcherEqual_Wildcard(t *testing.T) {
	tests := []testTagsMatcher{
		{
			name: `{"seriesByTag('name=a', 'b=c*')"}`, queries: []string{"seriesByTag('name=a', 'b=c*')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermEq, Value: "c*", HasWildcard: true,
										Glob: &WildcardItems{
											MinSize: 1, MaxSize: -1, P: "c", Inners: []items.InnerItem{items.ItemStar{}},
										},
									},
									Terminated: []string{"seriesByTag('name=a', 'b=c*')"},
								},
							},
						},
					},
				},
				Queries: map[string]bool{"seriesByTag('name=a', 'b=c*')": true},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=ca":      {"seriesByTag('name=a', 'b=c*')"},
				"a?b=c":            {"seriesByTag('name=a', 'b=c*')"},
				"a?a=v1&b=c&e=v3":  {"seriesByTag('name=a', 'b=c*')"},
				"a?a=v1&b=ca&e=v3": {"seriesByTag('name=a', 'b=c*')"},
			},
			missPaths: []string{"a?b=da", "a?b=v1", "a?c=v1", "b?a=v1"},
		},
		{
			name: `{"seriesByTag('name=a.b', 'b=c*.a')"}`, queries: []string{"seriesByTag('name=a.b', 'b=c*.a')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a.b"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermEq, Value: "c*.a", HasWildcard: true,
										Glob: &WildcardItems{
											MinSize: 3, MaxSize: -1, P: "c", Suffix: ".a",
											Inners: []items.InnerItem{items.ItemStar{}},
										},
									},
									Terminated: []string{"seriesByTag('name=a.b', 'b=c*.a')"},
								},
							},
						},
					},
				},
				Queries: map[string]bool{"seriesByTag('name=a.b', 'b=c*.a')": true},
			},
			matchPaths: map[string][]string{
				"a.b?a=v1&b=ca.a":      {"seriesByTag('name=a.b', 'b=c*.a')"},
				"a.b?b=c.a":            {"seriesByTag('name=a.b', 'b=c*.a')"},
				"a.b?a=v1&b=c.a&e=v3":  {"seriesByTag('name=a.b', 'b=c*.a')"},
				"a.b?a=v1&b=ca.a&e=v3": {"seriesByTag('name=a.b', 'b=c*.a')"},
				"a.b?a=v1&b=cb.a&e=v3": {"seriesByTag('name=a.b', 'b=c*.a')"},
			},
			missPaths: []string{"a?b=c.a", "a.b?b=da", "a.b?b=ca", "a.b?b=ca.b", "a.b?b=v1", "a.b?c=v1", "b?a=v1"},
		},
		// compaction
		{
			name: `{"seriesByTag('name=a', 'b=c[a]')"}`, queries: []string{"seriesByTag('name=a', 'b=c[a]')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term:       &TaggedTerm{Key: "b", Op: TaggedTermEq, Value: "ca"},
									Terminated: []string{"seriesByTag('name=a', 'b=c[a]')"},
								},
							},
						},
					},
				},
				Queries: map[string]bool{"seriesByTag('name=a', 'b=c[a]')": true},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=ca":      {"seriesByTag('name=a', 'b=c[a]')"},
				"a?b=ca":           {"seriesByTag('name=a', 'b=c[a]')"},
				"a?a=v1&b=ca&e=v3": {"seriesByTag('name=a', 'b=c[a]')"},
			},
			missPaths: []string{"a?b=c", "a?b=v1", "a?c=v1", "b?a=v1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestTagsMatcher(t, tt)
		})
	}
}

var (
	queryEqualW = "seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=c*')"
	pathEqualW  = "cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=sales&subproject=crm"
	regexEqualW = `^cpu\.load_avg\?(.*&)?app=postgresql(.*&)?project=sales(.*&)?subproject=c`
)

func BenchmarkEqualW_ByTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTagsMatcher()
		err := w.Add(queryEqualW)
		if err != nil {
			b.Fatal(err)
		}
		tags, err := PathTags(pathEqualW)
		if err != nil {
			b.Fatal(err)
		}

		queries := w.MatchByTags(tags)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func BenchmarkEqualW_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := regexp.MustCompile(regexEqualW)
		if !w.MatchString(pathEqualW) {
			b.Fatal(pathEqualW)
		}
	}
}

func BenchmarkEqualW_Precompiled_ByTags(b *testing.B) {
	w := NewTagsMatcher()
	err := w.Add(queryEqualW)
	if err != nil {
		b.Fatal(err)
	}
	queries := make([]string, 0, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tags, err := PathTags(pathEqualW)
		if err != nil {
			b.Fatal(err)
		}
		w.MatchByTagsB(tags, &queries)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func BenchmarkEqualW_Precompiled_ByTags2(b *testing.B) {
	w := NewTagsMatcher()
	err := w.Add(queryEqualW)
	if err != nil {
		b.Fatal(err)
	}
	queries := make([]string, 0, 1)
	tags, err := PathTags(pathEqualW)
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

func BenchmarkEqualW_Precompiled_Regex(b *testing.B) {
	w := regexp.MustCompile(regexEqualW)
	for i := 0; i < b.N; i++ {
		if !w.MatchString(pathEqualW) {
			b.Fatal(pathEqualW)
		}
	}
}
