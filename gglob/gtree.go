package gglob

import (
	"strings"

	"github.com/msaf1980/go-matcher/glob"
	"github.com/msaf1980/go-matcher/pkg/items"
)

type GTreeItem struct {
	Item *glob.Glob

	items.Terminated

	// TODO: may be some ordered tree for complete string nodes search speedup (on large set) ?
	ChildsMap map[string]*GTreeItem // full match
	Childs    []*GTreeItem          // next possible parts slice
}

func (item *GTreeItem) MatchItems(path string, globs *[]string, index *[]int, first items.Store) (matched int) {
	var part string
	part, path, _ = strings.Cut(path, ".")
	if part == "" {
		return
	}
	if len(item.ChildsMap) > 0 {
		if child, ok := item.ChildsMap[part]; ok {
			if path == "" {
				if child.Terminate {
					child.Append(globs, index, first)
					matched++
				}
			} else {
				if n := child.MatchItems(path, globs, index, first); n > 0 {
					matched += n
				}
			}
		}
	}
	for i := 0; i < len(item.Childs); i++ {
		if item.Childs[i].Item.Match(part) {
			if path == "" {
				if item.Childs[i].Terminate {
					item.Childs[i].Append(globs, index, first)
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

func (item *GTreeItem) MatchItemsByParts(parts []string, globs *[]string, index *[]int, first items.Store) (matched int) {
	if len(item.ChildsMap) > 0 {
		if child, ok := item.ChildsMap[parts[0]]; ok {
			if len(parts) == 1 {
				if child.Terminate {
					child.Append(globs, index, first)
					matched++
				}
			} else {
				if n := child.MatchItemsByParts(parts[1:], globs, index, first); n > 0 {
					matched += n
				}
			}
		}
	}
	for i := 0; i < len(item.Childs); i++ {
		if item.Childs[i].Item.Match(parts[0]) {
			if len(parts) == 1 {
				if item.Childs[i].Terminate {
					item.Childs[i].Append(globs, index, first)
					matched++
				}
			} else {
				if n := item.Childs[i].MatchItemsByParts(parts[1:], globs, index, first); n > 0 {
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
		treeItem = &GTreeItem{}
		treeMap[len(gg.Parts)] = treeItem
	}
	// treeItem := rootTree

	for i := 0; i < len(gg.Parts); i++ {
		if len(gg.Parts[i].Items) == 0 {
			// string
			if treeItem.ChildsMap == nil {
				treeItem.ChildsMap = make(map[string]*GTreeItem)
			}
			newItem, ok := treeItem.ChildsMap[gg.Parts[i].Node]
			if !ok {
				newItem = &GTreeItem{Item: gg.Parts[i]}
				treeItem.ChildsMap[gg.Parts[i].Node] = newItem
			}
			treeItem = newItem
		} else {
			newItem := LocateChildGTreeItem(treeItem.Childs, gg.Parts[i].Node)
			if newItem == nil {
				if treeItem.Childs == nil {
					treeItem.Childs = make([]*GTreeItem, 0, 2)
				}
				newItem = &GTreeItem{Item: gg.Parts[i]}
				treeItem.Childs = append(treeItem.Childs, newItem)
			}
			treeItem = newItem
		}
	}

	treeItem.Terminate = true
	treeItem.Query = gg.Node
	treeItem.Index = index

	return treeItem
}

// GGlobTree is batch glob matcher (dot-separated, like a.b*.c), writted for graphite project (use on large globs set)
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
	if normalized, ok = gtree.GlobsIndex[index]; ok {
		err = glob.ErrIndexDup
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

func (gtree *GGlobTree) MatchByParts(parts []string, globs *[]string, index *[]int, first items.Store) (matched int) {
	if len(parts) == 0 {
		return
	}
	if rootItem, ok := gtree.Root[len(parts)]; ok {
		if n := rootItem.MatchItemsByParts(parts, globs, index, first); n > 0 {
			matched += n
		}
	}

	return
}
