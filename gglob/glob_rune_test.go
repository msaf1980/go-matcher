package gglob

import (
	"strings"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

func TestGlobMatcher_Rune(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: `{"[a-c]"}`, globs: []string{"[a-c]"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "[a-c]", Terminated: []string{"[a-c]"},
								WildcardItems: wildcards.WildcardItems{
									MinSize: 1, MaxSize: 1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemRuneRanges{{'a', 'c'}},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{"[a-c]": -1},
			},
			matchPaths: map[string][]string{"a": {"[a-c]"}, "c": {"[a-c]"}, "b": {"[a-c]"}},
			missPaths:  []string{"", "d", "ab", "a.b"},
		},
		{
			name: `{"[a-c]z"}`, globs: []string{"[a-c]z"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "[a-c]z", Terminated: []string{"[a-c]z"},
								WildcardItems: wildcards.WildcardItems{
									MinSize: 2, MaxSize: 2, Suffix: "z",
									Inners: []wildcards.InnerItem{
										wildcards.ItemRuneRanges{{'a', 'c'}},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{"[a-c]z": -1},
			},
			matchPaths: map[string][]string{"az": {"[a-c]z"}, "cz": {"[a-c]z"}, "bz": {"[a-c]z"}},
			missPaths:  []string{"", "d", "ab", "dz", "a.z"},
		},
		{
			name: `{"[a-c]*"}`, globs: []string{"[a-c]*"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "[a-c]*", Terminated: []string{"[a-c]*"},
								WildcardItems: wildcards.WildcardItems{
									MinSize: 1, MaxSize: -1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemRuneRanges{{'a', 'c'}}, wildcards.ItemStar{},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{"[a-c]*": -1},
			},
			matchPaths: map[string][]string{
				"a": {"[a-c]*"}, "c": {"[a-c]*"},
				"az": {"[a-c]*"}, "cz": {"[a-c]*"}, "bz": {"[a-c]*"},
			},
			missPaths: []string{"", "d", "dz", "a.z"},
		},
		// one item optimization
		{
			name: `{"[a-]"}`, globs: []string{"[a-]"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a", Terminated: []string{"[a-]", "a"},
								WildcardItems: wildcards.WildcardItems{P: "a", MinSize: 1, MaxSize: 1},
							},
						},
					},
				},
				Globs: map[string]int{"[a-]": -1, "a": -1},
			},
			matchPaths: map[string][]string{"a": {"[a-]", "a"}},
			missPaths:  []string{"", "b", "d", "ab", "a.b"},
		},
		{
			name: `{"a[a-]Z"}`, globs: []string{"a[a-]Z"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "aaZ", Terminated: []string{"a[a-]Z", "aaZ"},
								WildcardItems: wildcards.WildcardItems{P: "aaZ", MinSize: 3, MaxSize: 3},
							},
						},
					},
				},
				Globs: map[string]int{"a[a-]Z": -1, "aaZ": -1},
			},
			matchPaths: map[string][]string{"aaZ": {"a[a-]Z", "aaZ"}},
			missPaths:  []string{"", "a", "b", "d", "ab", "aaz", "aaZa", "a.b"},
		},
		{
			name: `{"a[a-]Z[Q]"}`, globs: []string{"a[a-]Z[Q]"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "aaZQ", Terminated: []string{"a[a-]Z[Q]", "aaZQ"},
								WildcardItems: wildcards.WildcardItems{P: "aaZQ", MinSize: 4, MaxSize: 4},
							},
						},
					},
				},
				Globs: map[string]int{"a[a-]Z[Q]": -1, "aaZQ": -1},
			},
			matchPaths: map[string][]string{"aaZQ": {"a[a-]Z[Q]", "aaZQ"}},
			missPaths:  []string{"", "a", "Q", "aaZ", "aaZQa", "a.b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

func TestGlobMatcher_Rune_Broken(t *testing.T) {
	tests := []testGlobMatcher{
		// broken
		// compare with graphite-clickhouse. Now It's not error, but filter
		// (Path LIKE 'z%' AND match(Path, '^z[ac$')))
		{name: `{"z[ac"}`, globs: []string{"[ac"}, wantErr: true},
		{name: `{"a]c"}`, globs: []string{"a]c"}, wantErr: true},
		// skip empty
		{
			name: `{"[]a"}`, globs: []string{"[]a"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "a", Terminated: []string{"[]a", "a"},
								WildcardItems: wildcards.WildcardItems{P: "a", MinSize: 1, MaxSize: 1},
							},
						},
					},
				},
				Globs: map[string]int{"[]a": -1, "a": -1},
			},
			matchPaths: map[string][]string{"a": {"[]a", "a"}},
			missPaths:  []string{"", "b", "ab"},
		},
	}
	for _, tt := range tests {
		runTestGlobMatcher(t, tt)
	}
}

var (
	targetRune = "{a-bd-ef-kq-zA-QZ}"
	pathRune   = "Z"
)

// becnmark for suffix optimization
func BenchmarkRune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(len(targetRune))
		_, err := w.Add(targetRune, &buf)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathRune)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkRune_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetRune)
		if w.MatchString(pathRune) {
			b.Fatal(pathRune)
		}
	}
}

func BenchmarkRune_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetRune))
	_, err := w.Add(targetRune, &buf)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathRune)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkRune_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetRune))
	_, err := w.Add(targetRune, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathRune, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkRune_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetRune)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathRune) {
			b.Fatal(pathRune)
		}
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
		var buf strings.Builder
		buf.Grow(len(targetRuneStarMiss))
		_, err := w.Add(targetRuneStarMiss, &buf)
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
	var buf strings.Builder
	buf.Grow(len(targetRuneStarMiss))
	_, err := w.Add(targetRuneStarMiss, &buf)
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
	var buf strings.Builder
	buf.Grow(len(targetRuneStarMiss))
	_, err := w.Add(targetRuneStarMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
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
