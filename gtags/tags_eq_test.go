package gtags

import (
	"regexp"
	"testing"
)

func TestTagsMatcherEqual(t *testing.T) {
	tests := []testTagsMatcher{
		{
			name: "empty #1", queries: []string{},
			wantW: &TagsMatcher{
				Root:    &TaggedItem{Childs: []*TaggedItem{}},
				Queries: map[string]bool{},
			},
		},
		{
			name: "empty #2", queries: []string{""},
			wantW: &TagsMatcher{
				Root:    &TaggedItem{Childs: []*TaggedItem{}},
				Queries: map[string]bool{},
			},
		},
		// empty
		{
			name: `{"seriesByTag()"}`, queries: []string{"seriesByTag()"},
			wantErr: true,
		},
		// match
		{
			name: `{"seriesByTag('name=a', 'b=c')"}`, queries: []string{"seriesByTag('name=a', 'b=c')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term:       &TaggedTerm{Key: "b", Op: TaggedTermEq, Value: "c"},
									Terminated: []string{"seriesByTag('name=a', 'b=c')"},
								},
							},
						},
					},
				},
				Queries: map[string]bool{"seriesByTag('name=a', 'b=c')": true},
			},
			matchPaths: map[string][]string{
				"a?a=v1&b=c":      {"seriesByTag('name=a', 'b=c')"},
				"a?b=c":           {"seriesByTag('name=a', 'b=c')"},
				"a?a=v1&b=c&e=v3": {"seriesByTag('name=a', 'b=c')"},
			},
			missPaths: []string{"a?b=ca", "a?b=v1", "a?c=v1", "b?a=v1"},
		},
		{
			name:    `{"seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=crm')"}`,
			queries: []string{"seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=crm')"},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "cpu.load_avg"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{Key: "app", Op: TaggedTermEq, Value: "postgresql"},
									Childs: []*TaggedItem{
										{
											Term: &TaggedTerm{Key: "project", Op: TaggedTermEq, Value: "sales"},
											Childs: []*TaggedItem{
												{
													Term:       &TaggedTerm{Key: "subproject", Op: TaggedTermEq, Value: "crm"},
													Terminated: []string{"seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=crm')"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Queries: map[string]bool{
					"seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=crm')": true,
				},
			},
			matchPaths: map[string][]string{
				"cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=sales&subproject=crm": {
					"seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=crm')",
				},
			},
			missPaths: []string{
				"cpu.load_avg?app=crm&dc=dc1&host=node1-crm&project=sales&subproject=crm",
				"cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=backoffice&subproject=card",
			},
		},
		// duplicate seriesByTag item,
		{
			name: `{"seriesByTag('name=a', 'b=c', 'name=c')"}`, queries: []string{
				"seriesByTag('b=c','name=a', 'name=c')",
				"seriesByTag('name=a', 'b=c', 'name=c')",
			},
			wantW: &TagsMatcher{
				Root: &TaggedItem{
					Childs: []*TaggedItem{
						{
							Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "a"},
							Childs: []*TaggedItem{
								{
									Term: &TaggedTerm{Key: "__name__", Op: TaggedTermEq, Value: "c"},
									Childs: []*TaggedItem{
										{
											Term: &TaggedTerm{Key: "b", Op: TaggedTermEq, Value: "c"},
											Terminated: []string{
												"seriesByTag('b=c','name=a', 'name=c')",
												"seriesByTag('name=a', 'b=c', 'name=c')",
											},
										},
									},
								},
							},
						},
					},
				},
				Queries: map[string]bool{
					"seriesByTag('b=c','name=a', 'name=c')":  true,
					"seriesByTag('name=a', 'b=c', 'name=c')": true,
				},
			},
			matchPaths: map[string][]string{},
			missPaths:  []string{"a?a=v1&b=c", "c?a=v1&b=c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestTagsMatcher(t, tt)
		})
	}
}

var (
	queryEqual = "seriesByTag('name=cpu.load_avg', 'app=postgresql', 'project=sales', 'subproject=crm')"
	pathEqual  = "cpu.load_avg?app=postgresql&dc=dc1&host=node1-db&project=sales&subproject=crm"
	regexEqual = `^cpu\.load_avg\?(.*&)?app=postgresql(.*&)?project=sales(.*&)?subproject=crm(&|$)`
)

func BenchmarkEqual_ByTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTagsMatcher()
		err := w.Add(queryEqual)
		if err != nil {
			b.Fatal(err)
		}
		tags, err := PathTags(pathEqual)
		if err != nil {
			b.Fatal(err)
		}

		queries := w.MatchByTags(tags)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func BenchmarkEqual_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := regexp.MustCompile(regexEqual)
		if !w.MatchString(pathEqual) {
			b.Fatal(pathEqual)
		}
	}
}

func BenchmarkEqual_Precompiled_ByTags(b *testing.B) {
	w := NewTagsMatcher()
	err := w.Add(queryEqual)
	if err != nil {
		b.Fatal(err)
	}
	queries := make([]string, 0, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tags, err := PathTags(pathEqual)
		if err != nil {
			b.Fatal(err)
		}
		w.MatchByTagsB(tags, &queries)
		if len(queries) != 1 {
			b.Fatal(queries)
		}
	}
}

func BenchmarkEqual_Precompiled_ByTags2(b *testing.B) {
	w := NewTagsMatcher()
	err := w.Add(queryEqual)
	if err != nil {
		b.Fatal(err)
	}
	queries := make([]string, 0, 1)

	tags, err := PathTags(pathEqual)
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

func BenchmarkEqual_Precompiled_Regex(b *testing.B) {
	w := regexp.MustCompile(regexEqual)
	for i := 0; i < b.N; i++ {
		if !w.MatchString(pathEqual) {
			b.Fatal(pathEqual)
		}
	}
}
