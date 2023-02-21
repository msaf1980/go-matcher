package gtags

import (
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

var (
	cmpTransform = cmp.Transformer("Re", func(in *regexp.Regexp) string {
		if in == nil {
			return "<nil>"
		}
		return in.String()
	})
)

type testTagsMatcher struct {
	name       string
	queries    []string
	wantW      *TagsMatcher
	wantErr    bool
	matchPaths map[string][]string // must match with queries
	missPaths  []string
}

func runTestTagsMatcher(t *testing.T, tt testTagsMatcher) {
	w := NewTagsMatcher()
	var err error
	t.Run(tt.name, func(t *testing.T) {
		err = w.Adds(tt.queries)
	})
	if (err != nil) != tt.wantErr {
		t.Fatalf("TagsMatcher.Add() error = %v, wantErr %v", err, tt.wantErr)
	}
	if err == nil {
		if !cmp.Equal(w, tt.wantW, cmpTransform) {
			t.Errorf("TagsMatcher.Add() = %s", cmp.Diff(tt.wantW, w, cmpTransform))
		}
		verifyTagsMatcher(t, tt.name, tt.matchPaths, tt.missPaths, w)
	}
	if tt.wantErr {
		assert.Equal(t, 0, len(tt.matchPaths), "can't check on error")
		assert.Equal(t, 0, len(tt.missPaths), "can't check on error")
	}
}

func verifyTagsMatcher(t *testing.T, name string, matchTags map[string][]string, miss []string, w *TagsMatcher) {
	for path, wantTags := range matchTags {
		t.Run("GlobMatcher.Add("+name+") path="+path, func(t *testing.T) {
			tags, err := PathTags(path)
			if err != nil {
				t.Errorf("PathTags(%q) err = %q", path, err.Error())
			}
			if queries := w.MatchByTags(tags); !cmp.Equal(wantTags, queries) {
				t.Errorf("TagsMatcher.MatchByTags(%q) = %s", path, cmp.Diff(wantTags, queries))
			}
			tagsMap, err := PathTagsMap(path)
			if err != nil {
				t.Errorf("PathTagsMap(%q) err = %q", path, err.Error())
			}
			if queries := w.MatchByTagsMap(tagsMap); !cmp.Equal(wantTags, queries) {
				t.Errorf("TagsMatcher.MatchByTagsMap(%q) = %s", path, cmp.Diff(wantTags, queries))
			}
		})
	}
	for _, path := range miss {
		t.Run("GlobMatcher.Add("+name+") path="+path, func(t *testing.T) {
			tags, err := PathTags(path)
			if err != nil {
				t.Errorf("PathTags(%q) err = %q", path, err.Error())
			}
			if queries := w.MatchByTags(tags); len(queries) != 0 {
				t.Errorf("TagsMatcher.MatchByPath(%q) != %q", path, queries)
			}
			tagsMap, err := PathTagsMap(path)
			if err != nil {
				t.Errorf("PathTagsMap(%q) err = %q", path, err.Error())
			}
			if queries := w.MatchByTagsMap(tagsMap); len(queries) != 0 {
				t.Errorf("TagsMatcher.MatchByPathMap(%q) != %q", path, queries)
			}
		})
	}
}

type testTagsMatcherIndex struct {
	name       string
	queries    []string
	wantW      *TagsMatcher
	wantErr    bool
	matchPaths map[string][]int
}

func runTestTagsMatcherIndex(t *testing.T, tt testTagsMatcherIndex) {
	w := NewTagsMatcher()
	var (
		err error
		buf strings.Builder
	)
	t.Run(tt.name, func(t *testing.T) {
		for n, query := range tt.queries {
			_, err = w.AddIndexed(query, n, &buf)
			if err != nil {
				return
			}
		}
	})
	if err != nil {
		t.Fatalf("TagsMatcher.Add() error = %v, wantErr %v", err, tt.wantErr)
	}
	if !cmp.Equal(w, tt.wantW, cmpTransform) {
		t.Errorf("TagsMatcher.Add() = %s", cmp.Diff(tt.wantW, w, cmpTransform))
	}
	verifyTagsMatcherIndex(t, tt.name, tt.matchPaths, w)
}

func verifyTagsMatcherIndex(t *testing.T, name string, matchTags map[string][]int, w *TagsMatcher) {
	for path, wantN := range matchTags {
		t.Run("GlobMatcher.Add("+name+") path="+path, func(t *testing.T) {
			sort.Ints(wantN)
			wantFirst := -1
			if len(wantN) > 0 {
				wantFirst = wantN[0]
			}

			tags, err := PathTags(path)
			if err != nil {
				t.Errorf("PathTags(%q) err = %q", path, err.Error())
			}
			if queries := w.MatchIndexedByTags(tags); !cmp.Equal(wantN, queries) {
				t.Errorf("TagsMatcher.MatchIndexedByTags(%q) = %s", path, cmp.Diff(wantN, queries))
			}
			first := -1
			w.MatchFirstByTags(tags, &first)
			if first != wantFirst {
				t.Errorf("TagsMatcher.MatchFirstByTags(%q) = want %d, got %d", path, wantFirst, first)
			}

			tagsMap, err := PathTagsMap(path)
			if err != nil {
				t.Errorf("PathTagsMap(%q) err = %q", path, err.Error())
			}
			if queries := w.MatchIndexedByTagsMap(tagsMap); !cmp.Equal(wantN, queries) {
				t.Errorf("TagsMatcher.MatchIndexedByTagsMap(%q) = %s", path, cmp.Diff(wantN, queries))
			}
			first = -1
			w.MatchFirstByTagsMap(tagsMap, &first)
			if first != wantFirst {
				t.Errorf("TagsMatcher.MatchFirstByTagsMap(%q) = want %d, got %d", path, wantFirst, first)
			}
		})
	}
}

//////////////////////////////////////////////////////////////////////////////
// Benchmarks
//////////////////////////////////////////////////////////////////////////////
