package glob

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/msaf1980/go-matcher/pkg/items"
)

type TreeItemStr struct {
	Node string

	Reverse bool

	Terminated items.Terminated

	// TODO: may be some ordered tree for complete string nodes search speedup (on large set) ?
	Childs []*TreeItemStr `json:"childs"` // next possible parts slice
}

func StringTreeItem(treeItem *items.TreeItem) *TreeItemStr {
	if treeItem == nil {
		return nil
	}
	var node string
	if treeItem.Item != nil {
		node = treeItem.Item.String()
	}
	treeItemStr := &TreeItemStr{
		Node:       node,
		Reverse:    treeItem.Reverse,
		Childs:     make([]*TreeItemStr, 0, len(treeItem.Childs)),
		Terminated: treeItem.Terminated,
	}

	for _, child := range treeItem.Childs {
		treeItemStr.Childs = append(treeItemStr.Childs, StringTreeItem(child))
	}

	return treeItemStr
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

type globTreeStr struct {
	Root       *TreeItemStr
	Globs      map[string]int
	GlobsIndex map[int]string
}

type testGlobTree struct {
	globs   []string
	skipCmp bool // don't compare glob tree, only glob maps
	want    *globTreeStr
	match   map[string][]string
}

func runTestGlobTree(t *testing.T, n int, tt testGlobTree) {
	t.Run(fmt.Sprintf("%d#%#v", n, tt.globs), func(t *testing.T) {
		gtree := NewTree()
		for i, g := range tt.globs {
			_, _, err := gtree.Add(g, i)

			if err != nil && err != ErrGlobExist {
				t.Fatalf("GlobTree.Add(%q) error = %v", g, err)
			}
		}

		var globTree *globTreeStr
		if tt.skipCmp {
			globTree = &globTreeStr{
				Globs:      gtree.Globs,
				GlobsIndex: gtree.GlobsIndex,
			}
		} else {
			globTree = &globTreeStr{
				Root:       StringTreeItem(gtree.Root),
				Globs:      gtree.Globs,
				GlobsIndex: gtree.GlobsIndex,
			}
		}
		if !reflect.DeepEqual(globTree, tt.want) {
			t.Fatalf("GlobTree(%#v) = %s", tt.globs, cmp.Diff(tt.want, globTree))
		}

		verifyGlobTree(t, tt.globs, tt.match, gtree)
	})
}

func verifyGlobTree(t *testing.T, inGlobs []string, match map[string][]string, gtree *GlobTree) {
	for path, wantGlobs := range match {
		t.Run("#path="+path, func(t *testing.T) {
			var store items.AllStore
			store.Init()
			store.Grow(1)
			matched := gtree.Match(path, &store)

			verify := mergeVerify(store.S.S, store.Index.N)

			sort.Strings(store.S.S)
			sort.Strings(wantGlobs)
			sort.Ints(store.Index.N)

			if !reflect.DeepEqual(wantGlobs, store.S.S) {
				t.Fatalf("GlobTree(%#v).Match(%q) globs = %s", inGlobs, path, cmp.Diff(wantGlobs, store.S.S))
			}

			if matched != len(store.S.S) || len(store.S.S) != len(store.Index.N) {
				t.Fatalf("GlobTree(%#v).Match(%q) = %d, want %d, index = %d", inGlobs, path, matched, len(store.S.S), len(store.Index.N))
			}

			for _, v := range verify {
				if v.glob != gtree.GlobsIndex[v.index] {
					t.Errorf("GlobTree(%#v).Match(%q) index = %d glob = %s, want %s",
						inGlobs, path, v.index, gtree.GlobsIndex[v.index], v.glob)
				}
			}

			if len(store.Index.N) > 0 {
				if store.Min.Min != store.Index.N[0] {
					t.Errorf("GlobTree(%#v).Match(%q) first index = %d, want %d",
						inGlobs, path, store.Min.Min, store.Index.N[0])
				}
			}
		})
	}
}

func parseGlobs(globs []string) (g []*Glob) {
	g = make([]*Glob, len(globs))
	for i := 0; i < len(globs); i++ {
		g[i] = ParseMust(globs[i])
	}

	return
}
