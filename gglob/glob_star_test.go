package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_Star(t *testing.T) {
	tests := []testGlobMatcher{
		// deduplication
		{
			name: `{"a******c"}`, globs: []string{"a******c"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a******c", Terminated: "a******c", P: "a", Suffix: "c",
								MinSize: 2, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}},
							},
						},
					},
				},
				Globs: map[string]bool{"a******c": true},
			},
			matchPaths: map[string][]string{"ac": {"a******c"}, "abc": {"a******c"}, "abcc": {"a******c"}},
			missPaths:  []string{"", "acb"},
		},
		// * match
		{
			name: `{"*"}`, globs: []string{"*"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{Node: "*", Terminated: "*", Inners: []items.InnerItem{items.ItemStar{}}, MaxSize: -1},
						},
					},
				},
				Globs: map[string]bool{"*": true},
			},
			matchPaths: map[string][]string{"a": {"*"}, "b": {"*"}, "ce": {"*"}},
			missPaths:  []string{"", "b.c"},
		},
		{
			name: `{"a*c"}`, globs: []string{"a*c"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a*c", Terminated: "a*c", P: "a", Suffix: "c", MinSize: 2, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}},
							},
						},
					},
				},
				Globs: map[string]bool{"a*c": true},
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
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a*b?c", Terminated: "a*b?c", P: "a", Suffix: "c",
								MinSize: 4, MaxSize: -1,
								Inners: []items.InnerItem{items.ItemStar{}, items.ItemString("b"), items.ItemOne{}},
							},
						},
					},
				},
				Globs: map[string]bool{"a*b?c": true},
			},
			matchPaths: map[string][]string{
				"abec":   {"a*b?c"}, // skip *
				"abbec":  {"a*b?c"}, /// shit first b
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
