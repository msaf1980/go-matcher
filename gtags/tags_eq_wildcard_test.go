package gtags

import (
	"regexp"
	"strings"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
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
					Glob: &wildcards.WildcardItems{
						MinSize: 1, MaxSize: -1, P: "c",
						Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
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
					Glob: &wildcards.WildcardItems{
						MinSize: 3, MaxSize: -1, P: "c", Suffix: ".a",
						Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
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
					Glob: &wildcards.WildcardItems{
						P: "a", Suffix: "c.a", MinSize: 7, MaxSize: 9,
						Inners: []wildcards.InnerItem{
							&wildcards.ItemList{Vals: []string{"a", "bc"}, ValsMin: 1, ValsMax: 2},
							wildcards.ItemRune('Z'),
							&wildcards.ItemList{Vals: []string{"q", "qa"}, ValsMin: 1, ValsMax: 2},
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
					Glob: &wildcards.WildcardItems{
						P: "a", Suffix: "c", MinSize: 5, MaxSize: -1,
						Inners: []wildcards.InnerItem{wildcards.ItemNStar(3)},
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
										Glob: &wildcards.WildcardItems{
											MinSize: 1, MaxSize: -1, P: "c",
											Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
										},
									},
									Terminated: []string{
										"seriesByTag('name=a', 'b=c*')", "seriesByTag('__name__=a','b=c*')",
									},
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a', 'b=c*')": -1, "seriesByTag('__name__=a','b=c*')": -1,
				},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=ca":      {"seriesByTag('name=a', 'b=c*')", "seriesByTag('__name__=a','b=c*')"},
				"a?b=c":            {"seriesByTag('name=a', 'b=c*')", "seriesByTag('__name__=a','b=c*')"},
				"a?a=v1&b=c&e=v3":  {"seriesByTag('name=a', 'b=c*')", "seriesByTag('__name__=a','b=c*')"},
				"a?a=v1&b=ca&e=v3": {"seriesByTag('name=a', 'b=c*')", "seriesByTag('__name__=a','b=c*')"},
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
										Glob: &wildcards.WildcardItems{
											MinSize: 3, MaxSize: -1, P: "c", Suffix: ".a",
											Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
										},
									},
									Terminated: []string{
										"seriesByTag('name=a.b', 'b=c*.a')",
										"seriesByTag('__name__=a.b','b=c*.a')",
									},
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a.b', 'b=c*.a')": -1, "seriesByTag('__name__=a.b','b=c*.a')": -1},
			},
			matchPaths: map[string][]string{
				"a.b?a=v1&b=ca.a":      {"seriesByTag('name=a.b', 'b=c*.a')", "seriesByTag('__name__=a.b','b=c*.a')"},
				"a.b?b=c.a":            {"seriesByTag('name=a.b', 'b=c*.a')", "seriesByTag('__name__=a.b','b=c*.a')"},
				"a.b?a=v1&b=c.a&e=v3":  {"seriesByTag('name=a.b', 'b=c*.a')", "seriesByTag('__name__=a.b','b=c*.a')"},
				"a.b?a=v1&b=ca.a&e=v3": {"seriesByTag('name=a.b', 'b=c*.a')", "seriesByTag('__name__=a.b','b=c*.a')"},
				"a.b?a=v1&b=cb.a&e=v3": {"seriesByTag('name=a.b', 'b=c*.a')", "seriesByTag('__name__=a.b','b=c*.a')"},
			},
			missPaths: []string{"a?b=c.a", "a.b?b=da", "a.b?b=ca", "a.b?b=ca.b", "a.b?b=v1", "a.b?c=v1", "b?a=v1"},
		},
		{
			name:    `{"seriesByTag('name=a.b', 'b=a{a,bc}Z{qa,q}c.a')"}`,
			queries: []string{"seriesByTag('name=a.b', 'b=a{a,bc}Z{qa,q}c.a')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a.b"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermEq, Value: "a{a,bc}Z{q,qa}c.a",
										HasWildcard: true,
										Glob: &wildcards.WildcardItems{
											MinSize: 7, MaxSize: 9, P: "a", Suffix: "c.a",
											Inners: []wildcards.InnerItem{
												&wildcards.ItemList{
													Vals: []string{"a", "bc"}, ValsMin: 1, ValsMax: 2,
												},
												wildcards.ItemRune('Z'),
												&wildcards.ItemList{
													Vals: []string{"q", "qa"}, ValsMin: 1, ValsMax: 2,
												},
											},
										},
									},
									Terminated: []string{
										"seriesByTag('name=a.b', 'b=a{a,bc}Z{qa,q}c.a')",
										"seriesByTag('__name__=a.b','b=a{a,bc}Z{q,qa}c.a')",
									},
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a.b', 'b=a{a,bc}Z{qa,q}c.a')":    -1,
					"seriesByTag('__name__=a.b','b=a{a,bc}Z{q,qa}c.a')": -1,
				},
			},
			matchPaths: map[string][]string{
				"a.b?a=v1&b=aaZqc.a":   {"seriesByTag('name=a.b', 'b=a{a,bc}Z{qa,q}c.a')", "seriesByTag('__name__=a.b','b=a{a,bc}Z{q,qa}c.a')"},
				"a.b?a=v1&b=abcZqac.a": {"seriesByTag('name=a.b', 'b=a{a,bc}Z{qa,q}c.a')", "seriesByTag('__name__=a.b','b=a{a,bc}Z{q,qa}c.a')"},
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
									Term: &TaggedTerm{Key: "b", Op: TaggedTermEq, Value: "ca"},
									Terminated: []string{
										"seriesByTag('name=a', 'b=c[a]')",
										"seriesByTag('__name__=a','b=ca')"},
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a', 'b=c[a]')": -1, "seriesByTag('__name__=a','b=ca')": -1,
				},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=ca":      {"seriesByTag('name=a', 'b=c[a]')", "seriesByTag('__name__=a','b=ca')"},
				"a?b=ca":           {"seriesByTag('name=a', 'b=c[a]')", "seriesByTag('__name__=a','b=ca')"},
				"a?a=v1&b=ca&e=v3": {"seriesByTag('name=a', 'b=c[a]')", "seriesByTag('__name__=a','b=ca')"},
			},
			missPaths: []string{"a?b=c", "a?b=v1", "a?c=v1", "b?a=v1"},
		},
		{
			name:    `{"seriesByTag('name=a', 'b=a?*??c')"}`,
			queries: []string{"seriesByTag('name=a', 'b=a?*??c')", "seriesByTag('name=a', 'b=a*???c')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermEq, Value: "a*???c", HasWildcard: true,
										Glob: &wildcards.WildcardItems{
											P: "a", Suffix: "c", MinSize: 5, MaxSize: -1,
											Inners: []wildcards.InnerItem{wildcards.ItemNStar(3)},
										},
									},
									Terminated: []string{
										"seriesByTag('name=a', 'b=a?*??c')",
										"seriesByTag('__name__=a','b=a*???c')",
										"seriesByTag('name=a', 'b=a*???c')",
									},
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a', 'b=a?*??c')":    -1,
					"seriesByTag('__name__=a','b=a*???c')": -1,
					"seriesByTag('name=a', 'b=a*???c')":    -1,
				},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=aBCDc": {
					"seriesByTag('name=a', 'b=a?*??c')", "seriesByTag('__name__=a','b=a*???c')",
					"seriesByTag('name=a', 'b=a*???c')",
				},
				"a?b=aAFCDc": {
					"seriesByTag('name=a', 'b=a?*??c')", "seriesByTag('__name__=a','b=a*???c')",
					"seriesByTag('name=a', 'b=a*???c')",
				},
				"a?a=v1&b=aAFCDc&e=v3": {
					"seriesByTag('name=a', 'b=a?*??c')", "seriesByTag('__name__=a','b=a*???c')",
					"seriesByTag('name=a', 'b=a*???c')",
				},
				"a?b=c": {}, "a?b=v1": {}, "a?b=aCDc": {}, "a?c=v1": {}, "b?a=v1": {},
			},
		},
		{
			name:    `{"seriesByTag('name=a', 'b=c[a][Z-]*')"}`,
			queries: []string{"seriesByTag('name=a', 'b=c[a][Z-]*')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermEq, Value: "caZ*", HasWildcard: true,
										Glob: &wildcards.WildcardItems{
											MinSize: 3, MaxSize: -1, P: "caZ",
											Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
										},
									},
									Terminated: []string{
										"seriesByTag('name=a', 'b=c[a][Z-]*')", "seriesByTag('__name__=a','b=caZ*')",
									},
								},
							},
						},
					},
				},
				Queries: map[string]int{"seriesByTag('name=a', 'b=c[a][Z-]*')": -1, "seriesByTag('__name__=a','b=caZ*')": -1},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=caZ":       {"seriesByTag('name=a', 'b=c[a][Z-]*')", "seriesByTag('__name__=a','b=caZ*')"},
				"a?b=caZb":           {"seriesByTag('name=a', 'b=c[a][Z-]*')", "seriesByTag('__name__=a','b=caZ*')"},
				"a?a=v1&b=caZ&e=v3":  {"seriesByTag('name=a', 'b=c[a][Z-]*')", "seriesByTag('__name__=a','b=caZ*')"},
				"a?a=v1&b=caZd&e=v3": {"seriesByTag('name=a', 'b=c[a][Z-]*')", "seriesByTag('__name__=a','b=caZ*')"},
			},
			missPaths: []string{"a?b=c", "a?b=ca", "a?b=caz", "a?b=v1", "a?c=v1", "b?a=v1"},
		},
		{
			name:    `{"seriesByTag('name=a', 'b=a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l')"}`,
			queries: []string{"seriesByTag('name=a', 'b=a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{
										Key: "b", Op: TaggedTermEq, Value: "aaZQstLT*INN*zaSTltl",
										HasWildcard: true, Glob: &wildcards.WildcardItems{
											P: "aaZQstLT", Suffix: "zaSTltl", MinSize: 18, MaxSize: -1,
											Inners: []wildcards.InnerItem{
												wildcards.ItemStar{}, wildcards.ItemString("INN"), wildcards.ItemStar{},
											},
										},
									},
									Terminated: []string{
										"seriesByTag('name=a', 'b=a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l')",
										"seriesByTag('__name__=a','b=aaZQstLT*INN*zaSTltl')",
									},
								},
							},
						},
					},
				},
				Queries: map[string]int{
					"seriesByTag('name=a', 'b=a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l')": -1,
					"seriesByTag('__name__=a','b=aaZQstLT*INN*zaSTltl')":              -1,
				},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=aaZQstLTINNzaSTltl": {
					"seriesByTag('name=a', 'b=a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l')",
					"seriesByTag('__name__=a','b=aaZQstLT*INN*zaSTltl')",
				},
				"a?b=aaZQstLTINN_zaSTltl": {
					"seriesByTag('name=a', 'b=a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l')",
					"seriesByTag('__name__=a','b=aaZQstLT*INN*zaSTltl')",
				},
				"a?a=v1&b=aaZQstLT_INN_SKIP_zaSTltl&e=v3": {
					"seriesByTag('name=a', 'b=a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l')",
					"seriesByTag('__name__=a','b=aaZQstLT*INN*zaSTltl')",
				},
				"a?a=v1&b=aaZQstLT_SKIP_INN___zaSTltl&e=v3": {
					"seriesByTag('name=a', 'b=a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l')",
					"seriesByTag('__name__=a','b=aaZQstLT*INN*zaSTltl')",
				},
			},
			missPaths: []string{
				"a?b=c", "a?b=ca", "a?b=caz", "a?b=cazQl", "a?b=v1",
				"a?b=aaZQstLT_IN_zaSTltl",
				"a?b=aaZQstLT_INN_zSTltl",
				"a?c=v1", "b?a=v1",
			},
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

func BenchmarkEqualW_Terms(b *testing.B) {
	for i := 0; i < b.N; i++ {
		terms, err := ParseSeriesByTag(queryEqualW)
		if err != nil {
			b.Fatal(err)
		}
		if err = terms.Build(); err != nil {
			b.Fatal(err)
		}
		var buf strings.Builder
		buf.Grow(len(queryEqualW))
		terms.Rewrite(&buf)

		tags, err := PathTags(pathEqualW)
		if err != nil {
			b.Fatal(err)
		}

		if !terms.MatchByTags(tags) {
			b.Fatal(pathEqualW)
		}
	}
}

func BenchmarkEqualW_ByTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTagsMatcher()
		var buf strings.Builder
		buf.Grow(len(queryEqualW))
		_, err := w.Add(queryEqualW, &buf)
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

func BenchmarkEqualW_Precompiled_Terms(b *testing.B) {
	terms, err := ParseSeriesByTag(queryEqualW)
	if err != nil {
		b.Fatal(err)
	}
	if err = terms.Build(); err != nil {
		b.Fatal(err)
	}
	var buf strings.Builder
	buf.Grow(len(queryEqualW))
	terms.Rewrite(&buf)

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

func BenchmarkEqualW_Precompiled_Terms2(b *testing.B) {
	terms, err := ParseSeriesByTag(queryEqual)
	if err != nil {
		b.Fatal(err)
	}
	if err = terms.Build(); err != nil {
		b.Fatal(err)
	}
	var buf strings.Builder
	buf.Grow(len(queryEqualW))
	terms.Rewrite(&buf)

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

func BenchmarkEqualW_Precompiled_ByTags(b *testing.B) {
	w := NewTagsMatcher()
	var buf strings.Builder
	buf.Grow(len(queryEqualW))
	_, err := w.Add(queryEqualW, &buf)
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
		queries = queries[:0]
		w.MatchByTagsB(tags, &queries)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func BenchmarkEqualW_Precompiled_ByTags2(b *testing.B) {
	w := NewTagsMatcher()
	var buf strings.Builder
	buf.Grow(len(queryEqualW))
	_, err := w.Add(queryEqualW, &buf)
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
		queries = queries[:0]
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
