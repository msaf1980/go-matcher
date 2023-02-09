package gglob

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/stretchr/testify/assert"
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
	err := w.Adds(tt.globs)
	if (err != nil) != tt.wantErr {
		t.Errorf("GlobMatcher.Add() error = %v, wantErr %v", err, tt.wantErr)
		return
	}
	if err == nil {
		if !reflect.DeepEqual(w, tt.wantW) {
			t.Errorf("GlobMatcher.Add() = %s", cmp.Diff(tt.wantW, w))
		}
		verifyGlobMatcher(t, tt.matchPaths, tt.missPaths, w)
	}
	if tt.wantErr {
		assert.Equal(t, 0, len(tt.matchPaths), "can't check on error")
		assert.Equal(t, 0, len(tt.missPaths), "can't check on error")
	}
}

func verifyGlobMatcher(t *testing.T, matchGlobs map[string][]string, miss []string, w *GlobMatcher) {
	for path, wantGlobs := range matchGlobs {
		if globs := w.Match(path); !reflect.DeepEqual(wantGlobs, globs) {
			t.Errorf("GlobMatcher.Match(%q) = %s", path, cmp.Diff(wantGlobs, globs))
		}
		parts := items.PathSplit(path)
		if globs := w.MatchByParts(parts); !reflect.DeepEqual(wantGlobs, globs) {
			t.Errorf("GlobMatcher.MatchByParts(%q) = %s", path, cmp.Diff(wantGlobs, globs))
		}
	}
	for _, path := range miss {
		if globs := w.Match(path); len(globs) != 0 {
			t.Errorf("GlobMatcher.Match(%q) != %q", path, globs)
		}
		parts := items.PathSplit(path)
		if globs := w.MatchByParts(parts); len(globs) != 0 {
			t.Errorf("GlobMatcher.MatchByParts(%q) != %q", path, globs)
		}
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
