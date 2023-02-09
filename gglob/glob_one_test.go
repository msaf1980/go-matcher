package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_One(t *testing.T) {
	tests := []testGlobMatcher{
		// ? match
		{
			name: `{"?"}`, globs: []string{"?"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{Node: "?", Terminated: "?", Inners: []items.InnerItem{items.ItemOne{}}, MinSize: 1, MaxSize: 1},
						},
					},
				},
				Globs: map[string]bool{"?": true},
			},
			matchPaths: map[string][]string{"a": {"?"}, "c": {"?"}},
			missPaths:  []string{"", "ab", "a.b"},
		},
		{
			name: `{"a?"}`, globs: []string{"a?"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a?", Terminated: "a?", P: "a", Inners: []items.InnerItem{items.ItemOne{}},
								MinSize: 2, MaxSize: 2,
							},
						},
					},
				},
				Globs: map[string]bool{"a?": true},
			},
			matchPaths: map[string][]string{"ac": {"a?"}, "az": {"a?"}},
			missPaths:  []string{"", "a", "bc", "ace", "a.c"},
		},
		{
			name: `{"a?c"}`, globs: []string{"a?c"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a?c", Terminated: "a?c", P: "a", Inners: []items.InnerItem{items.ItemOne{}}, Suffix: "c",
								MinSize: 3, MaxSize: 3,
							},
						},
					},
				},
				Globs: map[string]bool{"a?c": true},
			},
			matchPaths: map[string][]string{"acc": {"a?c"}, "aec": {"a?c"}},
			missPaths:  []string{"", "ab", "ac", "ace", "a.c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

var (
	targetRuneStarMiss = "sy*abcdertg*[A-Z]*cabcdertg*[I-Q]*abcdertg*[A-Z]*babcdertg*cabcdertMISSg*tem"
	pathRuneStarMiss   = "sysabcdertgebaZbcdecabcdertglsIysabcdertgZebabcdertgicabcdertgltem"
)

// becnmark for suffix optimization
func BenchmarkRuneStarMiss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		err := w.Add(targetRuneStarMiss)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathRuneStarMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkRuneStarMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetRuneStarMiss)
		if w.MatchString(pathRuneStarMiss) {
			b.Fatal(pathRuneStarMiss)
		}
	}
}

func BenchmarkRuneStarMiss_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetRuneStarMiss)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathRuneStarMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkRuneStarMiss_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetRuneStarMiss)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.MatchB(pathRuneStarMiss, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkRuneStarMiss_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetRuneStarMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathRuneStarMiss) {
			b.Fatal(pathRuneStarMiss)
		}
	}
}
