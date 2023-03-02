package gglob

import (
	"strings"

	"github.com/msaf1980/go-matcher/glob"
	"github.com/msaf1980/go-matcher/pkg/items"
)

type GTreeItem struct {
	Item *glob.Glob

	Terminated string // end of chain (resulting raw/normalized globs)
	TermIndex  int    // rule num of end of chain (resulting glob), can be used in specific cases

	// TODO: may be some ordered tree for complete string nodes search speedup (on large set) ?
	Childs []*GTreeItem // next possible parts slice
}

func (item *GTreeItem) MatchItems(path string, globs *[]string, index *[]int, first items.Store) (matched int) {
	var part string
	part, path, _ = strings.Cut(path, ".")
	if part == "" {
		return
	}
	for i := 0; i < len(item.Childs); i++ {
		if item.Childs[i].Item.Match(part) {
			if path == "" {
				if item.Childs[i].Terminated != "" {
					if globs != nil {
						*globs = append(*globs, item.Childs[i].Terminated)
					}
					if index != nil {
						*index = append(*index, item.Childs[i].TermIndex)
					}
					if first != nil {
						first.Store(item.Childs[i].TermIndex)
					}
					matched++
				}
			} else {
				if n := item.Childs[i].MatchItems(path, globs, index, first); n > 0 {
					matched += n
				}
			}
		}
	}

	return
}

func LocateChildGTreeItem(childs []*GTreeItem, node string) *GTreeItem {
	for _, child := range childs {
		if child.Item != nil && child.Item.Node == node {
			return child
		}
	}
	return nil
}

func addGGlob(treeMap map[int]*GTreeItem, gg *GGlob, index int) *GTreeItem {
	treeItem, _ := treeMap[len(gg.Parts)]
	if treeItem == nil {
		treeItem = &GTreeItem{Childs: make([]*GTreeItem, 0, 8)}
		treeMap[len(gg.Parts)] = treeItem
	}
	// treeItem := rootTree

	for i := 0; i < len(gg.Parts); i++ {
		newItem := LocateChildGTreeItem(treeItem.Childs, gg.Parts[i].Node)
		if newItem == nil {
			newItem = &GTreeItem{Item: gg.Parts[i]}
			treeItem.Childs = append(treeItem.Childs, newItem)
		}
		treeItem = newItem
	}

	treeItem.Terminated = gg.Node
	treeItem.TermIndex = index

	return treeItem
}

// GGlobTree is batch glob matcher (dot-separated, like a.b*.c), writted for graphite project
type GGlobTree struct {
	Root       map[int]*GTreeItem
	Globs      map[string]int
	GlobsIndex map[int]string
}

func NewTree() *GGlobTree {
	return &GGlobTree{
		Root:       make(map[int]*GTreeItem),
		Globs:      make(map[string]int),
		GlobsIndex: make(map[int]string),
	}
}

func (gtree *GGlobTree) Add(globString string, index int) (normalized string, n int, err error) {
	if globString == "" {
		return
	}
	if index < 0 {
		err = glob.ErrIndexInvalid
		normalized = globString
		return
	}
	var ok bool
	if n, ok = gtree.Globs[globString]; ok {
		// aleady added
		err = glob.ErrGlobExist
		normalized = globString
		return
	}

	var g *GGlob
	if g, err = Parse(globString); err != nil {
		return
	}

	if n, ok = gtree.Globs[g.Node]; ok {
		// aleady added
		err = glob.ErrGlobExist
		normalized = g.Node
		return
	}

	if normalized, ok = gtree.GlobsIndex[index]; ok {
		err = glob.ErrIndexDup
		return
	}

	normalized = g.Node

	addGGlob(gtree.Root, g, index)

	gtree.Globs[globString] = index
	if normalized != globString {
		gtree.Globs[normalized] = index
	}
	gtree.Globs[normalized] = index
	gtree.GlobsIndex[index] = normalized

	n = index

	return
}

func (gtree *GGlobTree) AddGlob(g *GGlob, index int) (normalized string, n int, err error) {
	if index < 0 {
		err = glob.ErrIndexInvalid
		normalized = g.Node
		return
	}
	normalized = g.Node

	var ok bool
	if n, ok = gtree.Globs[g.Glob]; ok {
		// aleady added
		err = glob.ErrGlobExist
		return
	}

	if n, ok = gtree.Globs[g.Node]; ok {
		// aleady added
		err = glob.ErrGlobExist
		return
	}

	if normalized, ok = gtree.GlobsIndex[index]; ok {
		err = glob.ErrIndexDup
		return
	}

	addGGlob(gtree.Root, g, index)

	gtree.Globs[g.Node] = index
	if normalized != g.Node {
		gtree.Globs[normalized] = index
	}
	gtree.GlobsIndex[index] = normalized

	n = index

	return
}

func (gtree *GGlobTree) Match(path string, globs *[]string, index *[]int, first items.Store) (matched int) {
	if path == "" {
		return
	}
	path, partsCount := PathLevel(path)
	if rootItem, ok := gtree.Root[partsCount]; ok {
		if n := rootItem.MatchItems(path, globs, index, first); n > 0 {
			matched += n
		}
	}

	return
}