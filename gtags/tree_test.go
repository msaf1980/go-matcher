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

type taggedItemStr struct {
	Term string `json:"node"`

	Terminate  bool   `json:"terminate"`
	Terminated string `json:"terminated"` // end of chain (resulting raw/normalized globs)
	TermIndex  int    `json:"term_index"` // rule num of end of chain (resulting glob), can be used in specific cases

	Childs []*taggedItemStr `json:"childs"` // next possible parts slice
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
		Terminate:  treeItem.Terminate,
		Terminated: treeItem.Terminated,
		TermIndex:  treeItem.TermIndex,
	}

	if treeItem.Childs != nil {
		treeItemStr.Childs = make([]*taggedItemStr, 0, len(treeItem.Childs))
		for _, child := range treeItem.Childs {
			treeItemStr.Childs = append(treeItemStr.Childs, StringTaggedItem(child))
		}
	}
	return treeItemStr
}

type gTagsTreeStr struct {
	Root       *taggedItemStr
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
	for path, wantGlobs := range match {
		t.Run("#path="+path, func(t *testing.T) {
			queries := make([]string, 0, 4)
			index := make([]int, 0, 4)
			first := items.MinStore{-1}
			tags, err := PathTags(path)
			if err != nil {
				panic(err)
			}
			matched := gtree.MatchByTags(tags, &queries, &index, &first)

			verify := mergeVerify(queries, index)

			sort.Strings(queries)
			sort.Strings(wantGlobs)
			sort.Ints(index)

			if !reflect.DeepEqual(wantGlobs, queries) {
				t.Fatalf("GTagsTree(%#v).MatchByTags(%q) globs = %s", inGlobs, path, cmp.Diff(wantGlobs, queries))
			}

			if matched != len(queries) || len(queries) != len(index) {
				t.Fatalf("GTagsTree(%#v).MatchByTags(%q) = %d, want %d, index = %d", inGlobs, path, matched, len(queries), len(index))
			}

			for _, v := range verify {
				if v.glob != gtree.QueryIndex[v.index] {
					t.Errorf("GTagsTree(%#v).MatchByTags(%q) index = %d glob = %s, want %s",
						inGlobs, path, v.index, gtree.QueryIndex[v.index], v.glob)
				}
			}

			if len(index) > 0 {
				if first.N != index[0] {
					t.Errorf("GTagsTree(%#v).MatchByTags(%q) first index = %d, want %d",
						inGlobs, path, first, index[0])
				}
			}

			tagsMap, err := PathTagsMap(path)
			if err != nil {
				panic(err)
			}
			first.Init()
			queries = queries[:0]
			index = index[:0]
			matched = gtree.MatchByTagsMap(tagsMap, &queries, &index, &first)

			if !reflect.DeepEqual(wantGlobs, queries) {
				t.Fatalf("GTagsTree(%#v).MatchByTagsMap(%q) globs = %s", inGlobs, path, cmp.Diff(wantGlobs, queries))
			}

			if matched != len(queries) || len(queries) != len(index) {
				t.Fatalf("GTagsTree(%#v).MatchByTagsMap(%q) = %d, want %d, index = %d", inGlobs, path, matched, len(queries), len(index))
			}

			for _, v := range verify {
				if v.glob != gtree.QueryIndex[v.index] {
					t.Errorf("GTagsTree(%#v).MatchByTagsMap(%q) index = %d glob = %s, want %s",
						inGlobs, path, v.index, gtree.QueryIndex[v.index], v.glob)
				}
			}

			if len(index) > 0 {
				if first.N != index[0] {
					t.Errorf("GTagsTree(%#v).MatchByTagsMap(%q) first index = %d, want %d",
						inGlobs, path, first, index[0])
				}
			}

		})
	}
}
