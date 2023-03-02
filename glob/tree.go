package glob

import (
	"errors"

	"github.com/msaf1980/go-matcher/pkg/items"
)

var (
	ErrIndexInvalid = errors.New("index can't be negative")
	ErrIndexDup     = errors.New("duplicate index")
	ErrGlobExist    = errors.New("glob already exist")
)

func addGlob(rootTree *items.TreeItem, gg *Glob, index int) *items.TreeItem {
	treeItem := rootTree

	if gg.Suffix != "" {
		newItem := items.LocateChildTreeItem(treeItem.Childs, gg.Suffix, true)
		if newItem == nil {
			node := items.NewNodeItem(gg.Suffix, items.NewString(gg.Suffix))
			newItem = &items.TreeItem{NodeItem: node, Reverse: true}
			treeItem.Childs = append(treeItem.Childs, newItem)
		}
		treeItem = newItem
	}

	if gg.Prefix != "" {
		newItem := items.LocateChildTreeItem(treeItem.Childs, gg.Prefix, false)
		if newItem == nil {
			node := items.NewNodeItem(gg.Prefix, items.NewString(gg.Prefix))
			newItem = &items.TreeItem{NodeItem: node}
			treeItem.Childs = append(treeItem.Childs, newItem)
		}
		treeItem = newItem
	}

	for i := 0; i < len(gg.Items); i++ {
		newItem := items.LocateChildTreeItem(treeItem.Childs, gg.Items[i].Node, false)
		if newItem == nil {
			newItem = &items.TreeItem{NodeItem: gg.Items[i]}
			treeItem.Childs = append(treeItem.Childs, newItem)
		}
		treeItem = newItem
	}

	treeItem.Terminated = gg.Node
	treeItem.TermIndex = index

	return treeItem
}

// GlobTree is batch glob matcher
//
// TODO: need to be optimize, bad perfomance
type GlobTree struct {
	Root       *items.TreeItem
	Globs      map[string]int
	GlobsIndex map[int]string
}

func NewTree() *GlobTree {
	return &GlobTree{
		Root:       new(items.TreeItem),
		Globs:      make(map[string]int),
		GlobsIndex: make(map[int]string),
	}
}

func (gtree *GlobTree) Add(glob string, index int) (normalized string, n int, err error) {
	if glob == "" {
		return
	}
	if index < 0 {
		err = ErrIndexInvalid
		normalized = glob
		return
	}
	var ok bool
	if n, ok = gtree.Globs[glob]; ok {
		// aleady added
		err = ErrGlobExist
		normalized = glob
		return
	}

	var g *Glob
	if g, err = Parse(glob); err != nil {
		return
	}

	if n, ok = gtree.Globs[g.Node]; ok {
		// aleady added
		err = ErrGlobExist
		normalized = g.Node
		return
	}

	if normalized, ok = gtree.GlobsIndex[index]; ok {
		err = ErrIndexDup
		return
	}

	normalized = g.Node

	addGlob(gtree.Root, g, index)

	gtree.Globs[g.Glob] = index
	if normalized != g.Glob {
		gtree.Globs[normalized] = index
	}
	gtree.GlobsIndex[index] = normalized

	n = index

	return
}

func (gtree *GlobTree) AddGlob(g *Glob, index int) (normalized string, n int, err error) {
	if index < 0 {
		err = ErrIndexInvalid
		normalized = g.Node
		return
	}
	normalized = g.Node

	var ok bool
	if n, ok = gtree.Globs[g.Glob]; ok {
		// aleady added
		err = ErrGlobExist
		return
	}

	if n, ok = gtree.Globs[g.Node]; ok {
		// aleady added
		err = ErrGlobExist
		return
	}

	if normalized, ok = gtree.GlobsIndex[index]; ok {
		err = ErrIndexDup
		return
	}

	addGlob(gtree.Root, g, index)

	gtree.Globs[g.Glob] = index
	if normalized != g.Glob {
		gtree.Globs[normalized] = index
	}
	gtree.GlobsIndex[index] = normalized

	n = index

	return
}

func (gtree *GlobTree) Match(s string, globs *[]string, index *[]int, first items.Store) (matched int) {
	return gtree.Root.Match(s, globs, index, first)
}
