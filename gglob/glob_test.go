package gglob

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type testGlobMatcher struct {
	name       string
	globs      []string
	wantW      *GlobMatcher
	wantErr    bool
	matchGlobs map[string][]string // must match with glob
	miss       []string
}

func runTestGlobMatcher(t *testing.T, tt testGlobMatcher) {
	w := NewGlobMatcher()
	err := w.Adds(tt.globs)
	if (err != nil) != tt.wantErr {
		t.Errorf("GlobMatcher.Add() error = %v, wantErr %v", err, tt.wantErr)
		return
	}
	if !reflect.DeepEqual(w, tt.wantW) {
		t.Errorf("GlobMatcher.Add() = %s", cmp.Diff(tt.wantW, w))
	}
	if err == nil {
		for path, wantGlobs := range tt.matchGlobs {
			if globs := w.Match(path); !reflect.DeepEqual(wantGlobs, globs) {
				t.Errorf("GlobMatcher.Match(%q) = %s", path, cmp.Diff(wantGlobs, globs))
			}
		}
		for _, path := range tt.miss {
			if globs := w.Match(path); len(globs) != 0 {
				t.Errorf("GlobMatcher.Match(%q) != %q", path, globs)
			}
		}
	} else {
		assert.Equal(t, 0, len(tt.matchGlobs), "can't check on error")
		assert.Equal(t, 0, len(tt.miss), "can't check on error")
	}
}

func TestGlobMatcherString(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: "empty #1", globs: []string{},
			wantW: &GlobMatcher{
				Root:  map[int]*NodeItem{},
				Globs: map[string]bool{},
			},
		},
		{
			name: "empty #2", globs: []string{""},
			wantW: &GlobMatcher{
				Root:  map[int]*NodeItem{},
				Globs: map[string]bool{},
			},
		},
		// string match
		{
			name: `{"a"}`, globs: []string{"a"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs:    []*NodeItem{{Node: "a", Terminated: "a", InnerItem: InnerItem{Typ: NodeString, P: "a"}}},
					},
				},
				Globs: map[string]bool{"a": true},
			},
			matchGlobs: map[string][]string{"a": {"a"}},
			miss:       []string{"", "b", "ab", "ba"},
		},
		{
			name: `{"a.bc"}`, globs: []string{"a.bc"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					2: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a", InnerItem: InnerItem{Typ: NodeString, P: "a"},
								Childs: []*NodeItem{
									{Node: "bc", Terminated: "a.bc", InnerItem: InnerItem{Typ: NodeString, P: "bc"}},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"a.bc": true},
			},
			matchGlobs: map[string][]string{"a.bc": {"a.bc"}},
			miss:       []string{"", "b", "ab", "bc", "abc", "b.bc", "a.bce", "a.bc.e"},
		},
		{
			name: `{"a", "a.bc", "a.dc", "b.bc"}`, globs: []string{"a", "a.bc", "a.dc", "b.bc"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs:    []*NodeItem{{Node: "a", Terminated: "a", InnerItem: InnerItem{Typ: NodeString, P: "a"}}},
					},
					2: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a", InnerItem: InnerItem{Typ: NodeString, P: "a"},
								Childs: []*NodeItem{
									{Node: "bc", Terminated: "a.bc", InnerItem: InnerItem{Typ: NodeString, P: "bc"}},
									{Node: "dc", Terminated: "a.dc", InnerItem: InnerItem{Typ: NodeString, P: "dc"}},
								},
							},
							{
								Node: "b", InnerItem: InnerItem{Typ: NodeString, P: "b"},
								Childs: []*NodeItem{
									{Node: "bc", Terminated: "b.bc", InnerItem: InnerItem{Typ: NodeString, P: "bc"}},
								},
							},
						},
					},
				},
				Globs: map[string]bool{
					"a":    true,
					"a.bc": true,
					"a.dc": true,
					"b.bc": true,
				},
			},
			matchGlobs: map[string][]string{
				// "a":    {"a"},
				// "a.bc": {"a.bc"},
				"a.dc": {"a.dc"},
				// "b.bc": {"b.bc"},
			},
			miss: []string{"", "b", "ab", "bc", "abc", "c.bc", "a.be", "a.bce", "a.bc.e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

func TestGlobMatcher_One(t *testing.T) {
	tests := []testGlobMatcher{
		// ? match
		{
			name: `{"?"}`, globs: []string{"?"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{Node: "?", Terminated: "?", InnerItem: InnerItem{Typ: NodeOne}, MinSize: 1, MaxSize: 1},
						},
					},
				},
				Globs: map[string]bool{"?": true},
			},
			matchGlobs: map[string][]string{"a": {"?"}, "c": {"?"}},
			miss:       []string{"", "ab", "a.b"},
		},
		{
			name: `{"a?"}`, globs: []string{"a?"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{Node: "a?", Terminated: "a?", InnerItem: InnerItem{Typ: NodeOne, P: "a"}, MinSize: 2, MaxSize: 2},
						},
					},
				},
				Globs: map[string]bool{"a?": true},
			},
			matchGlobs: map[string][]string{"ac": {"a?"}, "az": {"a?"}},
			miss:       []string{"", "a", "bc", "ace", "a.c"},
		},
		{
			name: `{"a?c"}`, globs: []string{"a?c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a?c", Terminated: "a?c", InnerItem: InnerItem{Typ: NodeOne, P: "a"}, Suffix: "c",
								MinSize: 3, MaxSize: 3,
							},
						},
					},
				},
				Globs: map[string]bool{"a?c": true},
			},
			matchGlobs: map[string][]string{"acc": {"a?c"}, "aec": {"a?c"}},
			miss:       []string{"", "ab", "ac", "ace", "a.c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

func TestGlobMatcher_Star(t *testing.T) {
	tests := []testGlobMatcher{
		// * match
		{
			name: `{"*"}`, globs: []string{"*"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{Node: "*", Terminated: "*", InnerItem: InnerItem{Typ: NodeStar}, MaxSize: -1},
						},
					},
				},
				Globs: map[string]bool{"*": true},
			},
			matchGlobs: map[string][]string{"a": {"*"}, "b": {"*"}, "ce": {"*"}},
			miss:       []string{"", "b.c"},
		},
		{
			name: `{"a*c"}`, globs: []string{"a*c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{Node: "a*c", Terminated: "a*c", InnerItem: InnerItem{Typ: NodeStar, P: "a"}, Suffix: "c", MinSize: 2, MaxSize: -1},
						},
					},
				},
				Globs: map[string]bool{"a*c": true},
			},
			matchGlobs: map[string][]string{
				"ac": {"a*c"}, "acc": {"a*c"}, "aec": {"a*c"}, "aebc": {"a*c"},
				"aecc": {"a*c"}, "aecec": {"a*c"}, "abecec": {"a*c"},
			},
			miss: []string{"", "ab", "c", "ace", "a.c"},
		},
		// composite
		{
			name: `{"a*b?c"}`, globs: []string{"a*b?c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a*b?c", Terminated: "a*b?c", InnerItem: InnerItem{Typ: NodeInners, P: "a"}, Suffix: "c",
								MinSize: 4, MaxSize: -1,
								Inners: []*InnerItem{
									{Typ: NodeStar},
									{Typ: NodeString, P: "b"},
									{Typ: NodeOne},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"a*b?c": true},
			},
			matchGlobs: map[string][]string{
				"abec":   {"a*b?c"}, // skip *
				"abbec":  {"a*b?c"}, /// shit first b
				"acbbc":  {"a*b?c"},
				"aecbec": {"a*b?c"},
			},
			miss: []string{"", "ab", "c", "ace", "a.c", "abbece"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

func TestGlobMatcher_Rune(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: `{"[a-c]"}`, globs: []string{"[a-c]"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "[a-c]", Terminated: "[a-c]", MinSize: 1, MaxSize: 1,
								InnerItem: InnerItem{
									Typ: NodeRune, Runes: map[int32]struct{}{'a': {}, 'b': {}, 'c': {}},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"[a-c]": true},
			},
			matchGlobs: map[string][]string{"a": {"[a-c]"}, "c": {"[a-c]"}, "b": {"[a-c]"}},
			miss:       []string{"", "d", "ab", "a.b"},
		},
		{
			name: `{"[a-c]z"}`, globs: []string{"[a-c]z"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "[a-c]z", Terminated: "[a-c]z", MinSize: 2, MaxSize: 2, Suffix: "z",
								InnerItem: InnerItem{
									Typ: NodeRune, Runes: map[int32]struct{}{'a': {}, 'b': {}, 'c': {}},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"[a-c]z": true},
			},
			matchGlobs: map[string][]string{"az": {"[a-c]z"}, "cz": {"[a-c]z"}, "bz": {"[a-c]z"}},
			miss:       []string{"", "d", "ab", "dz", "a.z"},
		},
		{
			name: `{"[a-c]*"}`, globs: []string{"[a-c]*"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "[a-c]*", Terminated: "[a-c]*", MinSize: 1, MaxSize: -1,
								InnerItem: InnerItem{Typ: NodeInners},
								Inners: []*InnerItem{
									{Typ: NodeRune, Runes: map[int32]struct{}{'a': {}, 'b': {}, 'c': {}}},
									{Typ: NodeStar},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"[a-c]*": true},
			},
			matchGlobs: map[string][]string{
				"a": {"[a-c]*"}, "c": {"[a-c]*"},
				"az": {"[a-c]*"}, "cz": {"[a-c]*"}, "bz": {"[a-c]*"},
			},
			miss: []string{"", "d", "dz", "a.z"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

func TestGlobMatcher_Multi(t *testing.T) {
	tests := []testGlobMatcher{
		// composite
		{
			name: `{"a*c", "a*c*", "a*b?c", "a.b?d", "a*c.b"}`, globs: []string{"a*c", "a*c*", "a*b?c", "a.b?d", "a*c.b"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a*c", Terminated: "a*c", InnerItem: InnerItem{Typ: NodeStar, P: "a"}, Suffix: "c",
								MinSize: 2, MaxSize: -1,
							},
							{
								Node: "a*c*", Terminated: "a*c*", InnerItem: InnerItem{Typ: NodeInners, P: "a"},
								MinSize: 2, MaxSize: -1,
								Inners: []*InnerItem{
									{Typ: NodeStar},
									{Typ: NodeString, P: "c"},
									{Typ: NodeStar},
								},
							},
							{
								Node: "a*b?c", Terminated: "a*b?c", InnerItem: InnerItem{Typ: NodeInners, P: "a"}, Suffix: "c",
								MinSize: 4, MaxSize: -1,
								Inners: []*InnerItem{
									{Typ: NodeStar},
									{Typ: NodeString, P: "b"},
									{Typ: NodeOne},
								},
							},
						},
					},
					2: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a", InnerItem: InnerItem{Typ: NodeString, P: "a"},
								Childs: []*NodeItem{
									{Node: "b?d", Terminated: "a.b?d", InnerItem: InnerItem{Typ: NodeOne, P: "b"}, Suffix: "d", MinSize: 3, MaxSize: 3},
								},
							},
							{
								Node: "a*c", InnerItem: InnerItem{Typ: NodeStar, P: "a"}, Suffix: "c", MinSize: 2, MaxSize: -1,
								Childs: []*NodeItem{
									{Node: "b", Terminated: "a*c.b", InnerItem: InnerItem{Typ: NodeString, P: "b"}, MinSize: 0, MaxSize: 0},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"a*c": true, "a*c*": true, "a*b?c": true, "a*c.b": true, "a.b?d": true},
			},
			matchGlobs: map[string][]string{
				"acbec":  {"a*c", "a*c*", "a*b?c"},
				"abbece": {"a*c*"},
				"a.bfd":  {"a.b?d"},
			},
			miss: []string{"", "ab", "c", "a.b", "a.bd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

func buildGlobRegexp(g string) *regexp.Regexp {
	s := g
	s = strings.ReplaceAll(s, ".", `\.`)
	s = strings.ReplaceAll(s, "$", `\$`)
	s = strings.ReplaceAll(s, "{", "(")
	s = strings.ReplaceAll(s, "}", ")")
	s = strings.ReplaceAll(s, "?", `\?`)
	s = strings.ReplaceAll(s, ",", "|")
	s = strings.ReplaceAll(s, "*", ".*")
	return regexp.MustCompile("^" + s + "$")
}

var (
	targetSuffixMiss = "sy?abcdertg?babcdertg?cabcdertg?sy?abcdertg?babcdertg?cabcdertg?tem"
	pathSuffixMiss   = "sysabcdertgebabcdertgicabcdertglsysabcdertgebabcdertgicabcdertgltems"
)

// becnmark for suffix optimization
func BenchmarkSuffixMiss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		err := w.Add(targetSuffixMiss)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathSuffixMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkSuffixMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetSuffixMiss)
		if w.MatchString(pathSuffixMiss) {
			b.Fatal(pathSuffixMiss)
		}
	}
}

func BenchmarkSuffixMiss_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetSuffixMiss)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathSuffixMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkSuffixMiss_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetSuffixMiss)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.MatchP(pathSuffixMiss, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkSuffixMiss_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetSuffixMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathSuffixMiss) {
			b.Fatal(pathSuffixMiss)
		}
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
		w.MatchP(pathStarMiss, &globs)
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

var (
	targetSizeCheck = "sy*abcdertg*babcdertg*cabcdertg*sy*abcdertg*babcdertg*cabcdertg*tem.sy*abcdertg*babcdertg*cabcdertg*sy*abcdertg*babcdertg*cabcdertg*tem.sy*abcdertg*babcdertg*cabcdertg*sy*abcdertg*babcdertg*cabcdertg*tem"
	pathSizeCheck   = "sysabcdertgebabcdertgicadtglsysabcdertgebabcdertgicagltem.sysabcdertgebabcdertgicadtglsysabcdertgebabcdertgicagltem.sysabcdertgebabcdertgicadtglsysabcdertgebabcdertgicagltem"
)

// skip by size
func BenchmarkSizeCheck(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetSizeCheck)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathSizeCheck)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkSizeCheck_P(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetSizeCheck)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.MatchP(pathSizeCheck, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkSizeCheck_Regex(b *testing.B) {
	w := buildGlobRegexp(targetSizeCheck)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathSizeCheck) {
			b.Fatal(pathSizeCheck)
		}
	}
}
