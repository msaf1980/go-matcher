package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

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
						Childs: []*NodeItem{
							{Node: "a", Terminated: "a", Inners: []items.InnerItem{items.ItemString("a")}},
						},
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
						Childs: []*NodeItem{
							{
								Node: "a", Inners: []items.InnerItem{items.ItemString("a")},
								Childs: []*NodeItem{
									{Node: "bc", Terminated: "a.bc", Inners: []items.InnerItem{items.ItemString("bc")}},
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
						Childs: []*NodeItem{
							{Node: "a", Terminated: "a", Inners: []items.InnerItem{items.ItemString("a")}},
						},
					},
					2: {
						Childs: []*NodeItem{
							{
								Node: "a", Inners: []items.InnerItem{items.ItemString("a")},
								Childs: []*NodeItem{
									{Node: "bc", Terminated: "a.bc", Inners: []items.InnerItem{items.ItemString("bc")}},
									{Node: "dc", Terminated: "a.dc", Inners: []items.InnerItem{items.ItemString("dc")}},
								},
							},
							{
								Node: "b", Inners: []items.InnerItem{items.ItemString("b")},
								Childs: []*NodeItem{
									{Node: "bc", Terminated: "b.bc", Inners: []items.InnerItem{items.ItemString("bc")}},
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
				"a":    {"a"},
				"a.bc": {"a.bc"},
				"a.dc": {"a.dc"},
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
