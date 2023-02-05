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
						Childs:    []*NodeItem{{Node: "a", Terminated: true, InnerItem: InnerItem{Typ: NodeString, P: "a"}}},
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
									{Node: "a.bc", Terminated: true, InnerItem: InnerItem{Typ: NodeString, P: "bc"}},
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
			name: `{"a", "a.bc", "b.bc"}`, globs: []string{"a", "a.bc", "b.bc"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs:    []*NodeItem{{Node: "a", Terminated: true, InnerItem: InnerItem{Typ: NodeString, P: "a"}}},
					},
					2: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{
								Node: "a", InnerItem: InnerItem{Typ: NodeString, P: "a"},
								Childs: []*NodeItem{
									{Node: "a.bc", Terminated: true, InnerItem: InnerItem{Typ: NodeString, P: "bc"}},
								},
							},
							{
								Node: "b", InnerItem: InnerItem{Typ: NodeString, P: "b"},
								Childs: []*NodeItem{
									{Node: "b.bc", Terminated: true, InnerItem: InnerItem{Typ: NodeString, P: "bc"}},
								},
							},
						},
					},
				},
				Globs: map[string]bool{
					"a":    true,
					"a.bc": true,
					"b.bc": true,
				},
			},
			matchGlobs: map[string][]string{
				"a":    {"a"},
				"a.bc": {"a.bc"},
				"b.bc": {"b.bc"},
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
							{Node: "?", Terminated: true, InnerItem: InnerItem{Typ: NodeOne}, MinSize: 1, MaxSize: 1},
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
							{Node: "a?", Terminated: true, InnerItem: InnerItem{Typ: NodeOne, P: "a"}, MinSize: 2, MaxSize: 2},
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
							{Node: "a?c", Terminated: true, InnerItem: InnerItem{Typ: NodeOne, P: "a"}, Suffix: "c", MinSize: 3, MaxSize: 3},
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
							{Node: "*", Terminated: true, InnerItem: InnerItem{Typ: NodeStar}, MaxSize: -1},
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
							{Node: "a*c", Terminated: true, InnerItem: InnerItem{Typ: NodeStar, P: "a"}, Suffix: "c", MinSize: 2, MaxSize: -1},
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
								Node: "a*b?c", Terminated: true, InnerItem: InnerItem{Typ: NodeInners, P: "a"}, Suffix: "c",
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

func TestGlobMatcher_Multi(t *testing.T) {
	tests := []testGlobMatcher{
		// composite
		{
			name: `{"a*c", "a*b?c", "a.b?d"}`, globs: []string{"a*c", "a*b?c", "a.b?d"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						InnerItem: InnerItem{Typ: NodeRoot},
						Childs: []*NodeItem{
							{Node: "a*c", Terminated: true, InnerItem: InnerItem{Typ: NodeStar, P: "a"}, Suffix: "c", MinSize: 2, MaxSize: -1},
							{
								Node: "a*b?c", Terminated: true, InnerItem: InnerItem{Typ: NodeInners, P: "a"}, Suffix: "c",
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
									{Node: "a.b?d", Terminated: true, InnerItem: InnerItem{Typ: NodeOne, P: "b"}, Suffix: "d", MinSize: 3, MaxSize: 3}},
							},
						},
					},
				},
				Globs: map[string]bool{"a*c": true, "a*b?c": true, "a.b?d": true},
			},
			matchGlobs: map[string][]string{
				"acbec": {"a*c", "a*b?c"},
				"a.bfd": {"a.b?d"},
			},
			miss: []string{"", "ab", "c", "ace", "abbece", "a.b", "a.bd"},
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

func BenchmarkSuffixMiss_R(b *testing.B) {
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

func BenchmarkSuffixMiss_Precompiled_R(b *testing.B) {
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

func BenchmarkStarMiss_R(b *testing.B) {
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

func BenchmarkStarMiss_Precompiled_R(b *testing.B) {
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

func BenchmarkSizeCheck_R(b *testing.B) {
	w := buildGlobRegexp(targetSizeCheck)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathSizeCheck) {
			b.Fatal(pathSizeCheck)
		}
	}
}
