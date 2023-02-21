package gglob

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

type testGlobMatcher struct {
	name       string
	globs      []string
	wantW      *GlobMatcher
	wantErr    bool
	matchPaths map[string][]string // must match with glob
	missPaths  []string
}

func runTestGlobMatcher(t *testing.T, tt testGlobMatcher) {
	w := NewGlobMatcher()
	var err error
	t.Run(tt.name, func(t *testing.T) {
		err = w.Adds(tt.globs)
	})
	if (err != nil) != tt.wantErr {
		t.Fatalf("GlobMatcher.Add(%q) error = %v, wantErr %v", tt.globs, err, tt.wantErr)
	}
	if tt.wantErr {
		assert.Equal(t, 0, len(tt.matchPaths), "can't check on error")
		assert.Equal(t, 0, len(tt.missPaths), "can't check on error")
	}
	if err == nil {
		if !reflect.DeepEqual(w, tt.wantW) {
			t.Fatalf("GlobMatcher.Add(%q) = %s", tt.globs, cmp.Diff(tt.wantW, w))
		}
		verifyGlobMatcher(t, tt.name, tt.matchPaths, tt.missPaths, w)
	}
}

func verifyGlobMatcher(t *testing.T, name string, matchGlobs map[string][]string, miss []string, w *GlobMatcher) {
	for path, wantGlobs := range matchGlobs {
		t.Run(name+" path="+path, func(t *testing.T) {
			if globs := w.Match(path); !reflect.DeepEqual(wantGlobs, globs) {
				t.Errorf("GlobMatcher.Match(%q) = %s", path, cmp.Diff(wantGlobs, globs))
			}
			var globs []string
			w.MatchB(path, &globs)
			if !reflect.DeepEqual(wantGlobs, globs) {
				t.Errorf("GlobMatcher.MatchByParts(%q) = %s", path, cmp.Diff(wantGlobs, globs))
			}

			parts := wildcards.PathSplit(path)
			if globs := w.MatchByParts(parts); !reflect.DeepEqual(wantGlobs, globs) {
				t.Errorf("GlobMatcher.MatchByParts(%q) = %s", path, cmp.Diff(wantGlobs, globs))
			}
			globs = globs[:0]
			w.MatchByPartsB(parts, &globs)
			if !reflect.DeepEqual(wantGlobs, globs) {
				t.Errorf("GlobMatcher.MatchByParts(%q) = %s", path, cmp.Diff(wantGlobs, globs))
			}
		})
	}
	for _, path := range miss {
		t.Run("GlobMatcher.Add("+name+") path="+path, func(t *testing.T) {
			if globs := w.Match(path); len(globs) != 0 {
				t.Errorf("GlobMatcher.Match(%q) != %q", path, globs)
			}
			parts := wildcards.PathSplit(path)
			if globs := w.MatchByParts(parts); len(globs) != 0 {
				t.Errorf("GlobMatcher.MatchByParts(%q) != %q", path, globs)
			}
		})
	}
}

// Index matcher

type testGlobMatcherIndex struct {
	name       string
	globs      []string
	wantW      *GlobMatcher
	matchPaths map[string][]int
}

func runTestGlobMatcherIndex(t *testing.T, tt testGlobMatcherIndex) {
	w := NewGlobMatcher()
	var (
		err error
		buf strings.Builder
	)
	t.Run(tt.name, func(t *testing.T) {

		for n, glob := range tt.globs {
			buf.Grow(len(glob))
			_, err = w.AddIndexed(glob, n, &buf)
			if err != nil {
				err = fmt.Errorf("GlobMatcher.AddIndexed(%q) error = %v", glob, err)
				return
			}
		}
		if !reflect.DeepEqual(w, tt.wantW) {
			err = fmt.Errorf("GlobMatcher.AddIndexed(%q) = %s", tt.globs, cmp.Diff(tt.wantW, w))
		}
	})
	if err != nil {
		t.Fatal(err)
	}

	verifyGlobMatcherIndex(t, tt.name, tt.matchPaths, w)
}

func verifyGlobMatcherIndex(t *testing.T, name string, matchPaths map[string][]int, w *GlobMatcher) {
	for path, wantN := range matchPaths {
		t.Run(name+" path="+path, func(t *testing.T) {
			sort.Ints(wantN)
			wantFirst := -1
			if len(wantN) > 0 {
				wantFirst = wantN[0]
			}
			globsN := w.MatchIndexed(path)
			sort.Ints(globsN)
			if !reflect.DeepEqual(wantN, globsN) {
				t.Errorf("GlobMatcher.MatchIndexed(%q) = %s", path, cmp.Diff(wantN, globsN))
			}
			globsN = globsN[:0]
			w.MatchIndexedB(path, &globsN)
			sort.Ints(globsN)
			if !reflect.DeepEqual(wantN, globsN) {
				t.Errorf("GlobMatcher.MatchIndexed(%q) = %s", path, cmp.Diff(wantN, globsN))
			}

			first := -1
			w.MatchFirst(path, &first)
			if first != wantFirst {
				t.Errorf("GlobMatcher.MatchFirst(%q) = want %d, got %d", path, wantFirst, first)
			}

			parts := wildcards.PathSplit(path)
			globsN = w.MatchIndexedByParts(parts)
			if !reflect.DeepEqual(wantN, globsN) {
				t.Errorf("GlobMatcher.MatchIndexedByParts(%q) = %s", path, cmp.Diff(wantN, globsN))
			}
			globsN = globsN[:0]
			w.MatchIndexedByPartsB(parts, &globsN)
			sort.Ints(globsN)
			if !reflect.DeepEqual(wantN, globsN) {
				t.Errorf("GlobMatcher.MatchIndexed(%q) = %s", path, cmp.Diff(wantN, globsN))
			}

			first = -1
			w.MatchFirstByParts(parts, &first)
			if first != wantFirst {
				t.Errorf("GlobMatcher.MatchFirst(%q) = want %d, got %d", path, wantFirst, first)
			}
		})
	}
}

//////////////////////////////////////////////////////////////////////////////
// Benchmarks
//////////////////////////////////////////////////////////////////////////////

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
		var buf strings.Builder
		buf.Grow(len(targetSuffixMiss))
		_, err := w.Add(targetSuffixMiss, &buf)
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
	var buf strings.Builder
	buf.Grow(len(targetSuffixMiss))
	_, err := w.Add(targetSuffixMiss, &buf)
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
	var buf strings.Builder
	buf.Grow(len(targetSuffixMiss))
	_, err := w.Add(targetSuffixMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.MatchB(pathSuffixMiss, &globs)
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
	targetSizeCheck = "sy*abcdertg*babcdertg*cabcdertg*sy*abcdertg*babcdertg*cabcdertg*tem.sy*abcdertg*babcdertg*cabcdertg*sy*abcdertg*babcdertg*cabcdertg*tem.sy*abcdertg*babcdertg*cabcdertg*sy*abcdertg*babcdertg*cabcdertg*tem"
	pathSizeCheck   = "sysabcdertgebabcdertgicadtglsysabcdertgebabcdertgicagltem.sysabcdertgebabcdertgicadtglsysabcdertgebabcdertgicagltem.sysabcdertgebabcdertgicadtglsysabcdertgebabcdertgicagltem"
)

// skip by size
func BenchmarkSizeCheck(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetSizeCheck))
	_, err := w.Add(targetSizeCheck, &buf)
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
	var buf strings.Builder
	buf.Grow(len(targetSizeCheck))
	_, err := w.Add(targetSizeCheck, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathSizeCheck, &globs)
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
