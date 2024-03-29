package gtags

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/msaf1980/go-matcher/glob"
	"github.com/msaf1980/go-matcher/pkg/items"
)

type taggedItemsStr struct {
	Key        string
	Matched    []*taggedItemStr
	NotMatched []*taggedItemStr
}

func StringTaggedItems(items TaggedItems) taggedItemsStr {
	t := taggedItemsStr{
		Key: items.Key,
	}
	if items.Matched != nil {
		t.Matched = make([]*taggedItemStr, len(items.Matched))
		for i := 0; i < len(t.Matched); i++ {
			t.Matched[i] = StringTaggedItem(items.Matched[i])
		}
	}
	if items.NotMatched != nil {
		t.NotMatched = make([]*taggedItemStr, len(items.NotMatched))
		for i := 0; i < len(t.NotMatched); i++ {
			t.NotMatched[i] = StringTaggedItem(items.NotMatched[i])
		}
	}

	return t
}

type taggedItemStr struct {
	Term string `json:"node"`

	items.Terminated

	Items []taggedItemsStr `json:"match"` // next possible parts slice
}

func StringTaggedItem(treeItem *TaggedItem) *taggedItemStr {
	if treeItem == nil {
		return nil
	}
	var term string
	if treeItem.Term != nil {
		term = treeItem.Term.String()
	}
	treeItemStr := &taggedItemStr{
		Term:       term,
		Terminated: treeItem.Terminated,
	}

	if treeItem.Items != nil {
		for _, childs := range treeItem.Items {
			treeItemStr.Items = append(treeItemStr.Items, StringTaggedItems(childs))
		}
	}
	return treeItemStr
}

type gTagsTreeStr struct {
	Root *taggedItemStr
	items.Terminated
	Queries    map[string]int
	QueryIndex map[int]string
}

type verify struct {
	glob  string
	index int
}

func mergeVerify(globs []string, index []int) []verify {
	if len(globs) != len(index) {
		return nil
	}
	v := make([]verify, len(globs))
	for i := 0; i < len(globs); i++ {
		v[i].glob = globs[i]
		v[i].index = index[i]
	}
	return v
}

type testGTagsTree struct {
	queries []string
	skipCmp bool // don't compare glob tree, only glob maps
	want    *gTagsTreeStr
	match   map[string][]string
}

func runTestGTagsTree(t *testing.T, n int, tt testGTagsTree) {
	t.Run(fmt.Sprintf("%d#%#v", n, tt.queries), func(t *testing.T) {
		gtree := NewTree()
		for i, g := range tt.queries {
			_, _, err := gtree.Add(g, i)
			if err != nil && err != glob.ErrGlobExist {
				t.Fatalf("GlobTree.Add(%q) error = %v", g, err)
			}
		}

		var globTree *gTagsTreeStr
		if tt.skipCmp {
			globTree = &gTagsTreeStr{
				Queries:    gtree.Queries,
				QueryIndex: gtree.QueryIndex,
			}
		} else {
			globTree = &gTagsTreeStr{
				Root:       StringTaggedItem(gtree.Root),
				Terminated: gtree.Terminated,
				Queries:    gtree.Queries,
				QueryIndex: gtree.QueryIndex,
			}
		}
		if !reflect.DeepEqual(globTree, tt.want) {
			t.Fatalf("GlobTree(%#v) = %s", tt.queries, cmp.Diff(tt.want, globTree))
		}

		verifyGTagsTree(t, tt.queries, tt.match, gtree)
	})
}

func verifyGTagsTree(t *testing.T, inGlobs []string, match map[string][]string, gtree *GTagsTree) {
	for path, wantQueries := range match {
		t.Run("#path="+path, func(t *testing.T) {
			var store items.AllStore
			store.Init()
			store.Grow(1)
			tags, err := PathTags(path)
			if err != nil {
				panic(err)
			}
			matched := gtree.MatchByTags(tags, &store)

			verify := mergeVerify(store.S.S, store.Index.N)

			sort.Strings(store.S.S)
			sort.Strings(wantQueries)
			sort.Ints(store.Index.N)

			if !reflect.DeepEqual(wantQueries, store.S.S) {
				t.Fatalf("GTagsTree(%#v).MatchByTags(%q) globs = %s", inGlobs, path, cmp.Diff(wantQueries, store.S.S))
			}

			if matched != len(store.S.S) || len(store.S.S) != len(store.Index.N) {
				t.Fatalf("GTagsTree(%#v).MatchByTags(%q) = %d, want %d, index = %d", inGlobs, path, matched, len(store.S.S), len(store.Index.N))
			}

			for _, v := range verify {
				if v.glob != gtree.QueryIndex[v.index] {
					t.Errorf("GTagsTree(%#v).MatchByTags(%q) index = %d glob = %s, want %s",
						inGlobs, path, v.index, gtree.QueryIndex[v.index], v.glob)
				}
			}

			if len(store.Index.N) > 0 {
				if store.Min.Min != store.Index.N[0] {
					t.Errorf("GTagsTree(%#v).MatchByTags(%q) first index = %d, want %d",
						inGlobs, path, store.Min.Min, store.Index.N[0])
				}
			}

			tagsMap, err := PathTagsMap(path)
			if err != nil {
				panic(err)
			}
			store.Init()
			matched = gtree.MatchByTagsMap(tagsMap, &store)

			if !reflect.DeepEqual(wantQueries, store.S.S) {
				t.Fatalf("GTagsTree(%#v).MatchByTagsMap(%q) globs = %s", inGlobs, path, cmp.Diff(wantQueries, store.S.S))
			}

			if matched != len(store.S.S) || len(store.S.S) != len(store.Index.N) {
				t.Fatalf("GTagsTree(%#v).MatchByTagsMap(%q) = %d, want %d, index = %d", inGlobs, path, matched, len(store.S.S), len(store.Index.N))
			}

			for _, v := range verify {
				if v.glob != gtree.QueryIndex[v.index] {
					t.Errorf("GTagsTree(%#v).MatchByTagsMap(%q) index = %d glob = %s, want %s",
						inGlobs, path, v.index, gtree.QueryIndex[v.index], v.glob)
				}
			}

			if len(store.Index.N) > 0 {
				if store.Min.Min != store.Index.N[0] {
					t.Errorf("GTagsTree(%#v).MatchByTagsMap(%q) first index = %d, want %d",
						inGlobs, path, store.Min.Min, store.Index.N[0])
				}
			}
		})
	}
}
