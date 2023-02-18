package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

func TestGlobMatcher_Star(t *testing.T) {
	tests := []testGlobMatcher{
		// deduplication
		{
			name: `{"a******c"}`, globs: []string{"a******c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a******c", Terminated: "a******c", TermIndex: -1,
								P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
								Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
							},
						},
					},
				},
				Globs: map[string]int{"a******c": -1},
			},
			matchPaths: map[string][]string{"ac": {"a******c"}, "abc": {"a******c"}, "abcc": {"a******c"}},
			missPaths:  []string{"", "acb"},
		},
		// * match
		{
			name: `{"*"}`, globs: []string{"*"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "*", Terminated: "*", TermIndex: -1,
								Inners: []wildcards.InnerItem{wildcards.ItemStar{}}, MaxSize: -1,
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
								Node: "a*c", Terminated: "a*c", TermIndex: -1,
								P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
								Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
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
								Node: "a*b?c", Terminated: "a*b?c", TermIndex: -1,
								P: "a", Suffix: "c", MinSize: 4, MaxSize: -1,
								Inners: []wildcards.InnerItem{wildcards.ItemStar{}, wildcards.ItemRune('b'), wildcards.ItemOne{}},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
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
		err := w.Add(targetStarMiss)
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
	err := w.Add(targetStarMiss)
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
	err := w.Add(targetStarMiss)
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
