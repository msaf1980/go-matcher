package gglob

import (
	"strings"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/globs"
	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_Star(t *testing.T) {
	tests := []testGlobMatcher{
		// deduplication
		{
			name: `{"a******c"}`, globs: []string{"a******c"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a*c", Terminated: []string{"a******c", "a*c"},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
									Inners: []items.Item{items.ItemStar{}},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a******c": -1, "a*c": -1},
			},
			matchPaths: map[string][]string{
				"ac": {"a******c", "a*c"}, "abc": {"a******c", "a*c"}, "abcc": {"a******c", "a*c"}},
			missPaths: []string{"", "acb"},
		},
		// * match
		{
			name: `{"*"}`, globs: []string{"*"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "*", Terminated: []string{"*"},
								NodeItem: items.NodeItem{
									Inners: []items.Item{items.ItemStar{}}, MaxSize: -1,
								},
							},
						},
					},
				},
				Globs: map[string]int{"*": -1},
			},
			matchPaths: map[string][]string{"a": {"*"}, "b": {"*"}, "ce": {"*"}},
			missPaths:  []string{"", "b.c"},
		},
		{
			name: `{"a*c"}`, globs: []string{"a*c"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a*c", Terminated: []string{"a*c"},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
									Inners: []items.Item{items.ItemStar{}},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a*c": -1},
			},
			matchPaths: map[string][]string{
				"ac": {"a*c"}, "acc": {"a*c"}, "aec": {"a*c"}, "aebc": {"a*c"},
				"aecc": {"a*c"}, "aecec": {"a*c"}, "abecec": {"a*c"},
			},
			missPaths: []string{"", "ab", "c", "ace", "a.c"},
		},
		// composite
		{
			name: `{"a*b?c"}`, globs: []string{"a*b?c"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a*b?c", Terminated: []string{"a*b?c"},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "c", MinSize: 4, MaxSize: -1,
									Inners: []items.Item{
										items.ItemStar{}, items.ItemRune('b'), items.ItemOne{},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a*b?c": -1},
			},
			matchPaths: map[string][]string{
				"abec":   {"a*b?c"}, // skip *
				"abbec":  {"a*b?c"}, /// shift first b
				"acbbc":  {"a*b?c"},
				"aecbec": {"a*b?c"},
			},
			missPaths: []string{"", "ab", "c", "ace", "a.c", "abbece"},
		},
		{
			name: `{"a*?_FIND*st"}`, globs: []string{"a*?_FIND*_st"},
			wantW: &GlobMatcher{
				Root: map[int]*globs.NodeItem{
					1: {
						Childs: []*globs.NodeItem{
							{
								Node: "a*?_FIND*_st", Terminated: []string{"a*?_FIND*_st"},
								NodeItem: items.NodeItem{
									P: "a", Suffix: "_st", MinSize: 10, MaxSize: -1,
									Inners: []items.Item{
										items.ItemNStar(1),
										items.ItemString("_FIND"), items.ItemStar{},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a*?_FIND*_st": -1},
			},
			matchPaths: map[string][]string{
				"ab_FIND_st":        {"a*?_FIND*_st"},
				"aLc_FIND_st":       {"a*?_FIND*_st"},
				"aLBc_FIND_st":      {"a*?_FIND*_st"},
				"aLBc_FIND_STAR_st": {"a*?_FIND*_st"},
				"aLBc_FINDB_st":     {"a*?_FIND*_st"},
			},
			missPaths: []string{"a_FIND_st", "a_FINDB_st"},
		},
	}
	for _, tt := range tests {
		runTestGlobMatcher(t, tt)
	}
}

var (
	targetStarMiss = "sy*abcdertg*babcdertg*cabcdertg*sy*abcdertg*babcdertg*cabcdertMISSg*tem"
	pathStarMiss   = "sysabcdertgebabcdertgicabcdertglsysabcdertgebabcdertgicabcdertgltem"
)

// becnmark for suffix optimization
func BenchmarkStarMiss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(len(targetStarMiss))
		_, err := w.Add(targetStarMiss, &buf)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathStarMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkStarMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetStarMiss)
		if w.MatchString(pathStarMiss) {
			b.Fatal(pathStarMiss)
		}
	}
}

func BenchmarkStarMiss_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetStarMiss))
	_, err := w.Add(targetStarMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathStarMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkStarMiss_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetStarMiss))
	_, err := w.Add(targetStarMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathStarMiss, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkStarMiss_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetStarMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathStarMiss) {
			b.Fatal(pathStarMiss)
		}
	}
}
