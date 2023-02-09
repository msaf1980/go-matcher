package gtags

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type testTagsMatcher struct {
	name       string
	queries    []string
	wantW      *TagsMatcher
	wantErr    bool
	matchPaths map[string][]string // must match with glob
	missPaths  []string
}

var (
	cmpTransform = cmp.Transformer("Re", func(in *regexp.Regexp) string {
		if in == nil {
			return "<nil>"
		}
		return in.String()
	})
)

func runTestTagsMatcher(t *testing.T, tt testTagsMatcher) {
	w := NewTagsMatcher()
	err := w.Adds(tt.queries)
	if (err != nil) != tt.wantErr {
		t.Errorf("TagsMatcher.Add() error = %v, wantErr %v", err, tt.wantErr)
		return
	}
	if err == nil {
		if !cmp.Equal(w, tt.wantW, cmpTransform) {
			t.Errorf("TagsMatcher.Add() = %s", cmp.Diff(tt.wantW, w, cmpTransform))
		}
		verifyTagsMatcher(t, tt.matchPaths, tt.missPaths, w)
	}
	if tt.wantErr {
		assert.Equal(t, 0, len(tt.matchPaths), "can't check on error")
		assert.Equal(t, 0, len(tt.missPaths), "can't check on error")
	}
}

func verifyTagsMatcher(t *testing.T, matchTags map[string][]string, miss []string, w *TagsMatcher) {
	for path, wantTags := range matchTags {
		if queries := w.MatchByPath(path); !cmp.Equal(wantTags, queries) {
			t.Errorf("TagsMatcher.MatchByTags(%q) = %s", path, cmp.Diff(wantTags, queries))
		}
		tags, err := PathTagsMap(path)
		if err != nil {
			t.Errorf("ParsePath(%q) err = %q", path, err.Error())
		}
		if queries := w.MatchByTags(tags); !cmp.Equal(wantTags, queries) {
			t.Errorf("TagsMatcher.MatchByTags(%q) = %s", path, cmp.Diff(wantTags, queries))
		}
	}
	for _, path := range miss {
		if queries := w.MatchByPath(path); len(queries) != 0 {
			t.Errorf("TagsMatcher.MatchByPath(%q) != %q", path, queries)
		}
		tags, err := PathTagsMap(path)
		if err != nil {
			t.Errorf("ParsePath(%q) err = %q", path, err.Error())
		}
		if queries := w.MatchByTags(tags); len(queries) != 0 {
			t.Errorf("TagsMatcher.MatchByPath(%q) != %q", path, queries)
		}
	}
}

//////////////////////////////////////////////////////////////////////////////
// Benchmarks
//////////////////////////////////////////////////////////////////////////////
