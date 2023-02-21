package gglob

import (
	"strings"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

func TestGlobMatcher_NStar(t *testing.T) {
	tests := []testGlobMatcher{
		// deduplication
		{
			name: `{"a******?c", "a?******c", "a?******?c"}`,
			globs: []string{
				"a******?c", "a?******c", "a**?****c",
				"a?******?c", "a**??c",
				"a?*??c", "a*?*?*c", "a??*?*c", "a???**c",
			},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a*?c", Terminated: []string{
									"a******?c", "a*?c", "a?******c", "a**?****c",
								},
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c", MinSize: 3, MaxSize: -1,
									Inners: []wildcards.InnerItem{wildcards.ItemNStar(1)},
								},
							},
							{
								Node: "a*??c", Terminated: []string{
									"a?******?c", "a*??c", "a**??c", "a*?*?*c",
								},
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c", MinSize: 4, MaxSize: -1,
									Inners: []wildcards.InnerItem{wildcards.ItemNStar(2)},
								},
							},
							{
								Node:       "a*???c",
								Terminated: []string{"a?*??c", "a*???c", "a??*?*c", "a???**c"},
								WildcardItems: wildcards.WildcardItems{
									MinSize: 5,
									MaxSize: -1,
									P:       "a",
									Suffix:  "c",
									Inners:  []wildcards.InnerItem{wildcards.ItemNStar(3)},
								},
							},
						},
					},
				},
				Globs: map[string]int{
					"a?******c": -1, "a**?****c": -1, "a******?c": -1, "a*?c": -1,
					"a?******?c": -1, "a**??c": -1, "a*?*?*c": -1, "a*??c": -1,
					"a?*??c": -1, "a*???c": -1, "a??*?*c": -1, "a???**c": -1,
				},
			},
			matchPaths: map[string][]string{
				"aBc": {"a******?c", "a*?c", "a?******c", "a**?****c"},
				"aBCc": {
					"a******?c", "a*?c", "a?******c", "a**?****c",
					"a?******?c", "a*??c", "a**??c", "a*?*?*c",
				},
				"aBCDc": {
					"a******?c", "a*?c", "a?******c", "a**?****c",
					"a?******?c", "a*??c", "a**??c", "a*?*?*c",
					"a?*??c", "a*???c", "a??*?*c", "a???**c",
				},
				"aBCDEc": {
					"a******?c", "a*?c", "a?******c", "a**?****c",
					"a?******?c", "a*??c", "a**??c", "a*?*?*c",
					"a?*??c", "a*???c", "a??*?*c", "a???**c",
				},
			},
			missPaths: []string{"", "ac", "acb"},
		},
		// * match
		{
			name: `{"*"}`, globs: []string{"*"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "*", Terminated: []string{"*"},
								WildcardItems: wildcards.WildcardItems{
									Inners: []wildcards.InnerItem{wildcards.ItemStar{}}, MaxSize: -1,
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
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a*c", Terminated: []string{"a*c"},
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
									Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
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
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a*b?c", Terminated: []string{"a*b?c"},
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "c", MinSize: 4, MaxSize: -1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemStar{}, wildcards.ItemRune('b'), wildcards.ItemOne{},
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
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a*?_FIND*_st", Terminated: []string{"a*?_FIND*_st"},
								WildcardItems: wildcards.WildcardItems{
									P: "a", Suffix: "_st", MinSize: 10, MaxSize: -1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemNStar(1),
										wildcards.ItemString("_FIND"), wildcards.ItemStar{},
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
				"aLBc_FIND_Star_st": {"a*?_FIND*_st"},
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
	targetNStarMiss = "sy*abcdertg*babcdertg*cabcdertg*sy*abcdertg*babcdertg*cabcdertMISSg*tem"
	pathNStarMiss   = "sysabcdertgebabcdertgicabcdertglsysabcdertgebabcdertgicabcdertgltem"
)

// becnmark for suffix optimization
func BenchmarkNStarMiss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(len(targetNStarMiss))
		_, err := w.Add(targetNStarMiss, &buf)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathNStarMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkNStarMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetNStarMiss)
		if w.MatchString(pathNStarMiss) {
			b.Fatal(pathNStarMiss)
		}
	}
}

func BenchmarkNStarMiss_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetNStarMiss))
	_, err := w.Add(targetNStarMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathNStarMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkNStarMiss_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetNStarMiss))
	_, err := w.Add(targetNStarMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathNStarMiss, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkNStarMiss_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetNStarMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathNStarMiss) {
			b.Fatal(pathNStarMiss)
		}
	}
}
