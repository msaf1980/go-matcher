package gglob

import "github.com/msaf1980/go-matcher/glob"

type TreeItem struct {
	Node string // raw part string

	Terminated []string // end of chain (resulting raw/normalized globs)
	TermIndex  []int    // rule num of end of chain (resulting glob), can be used in specific cases

	Item *glob.Glob

	// TODO: may be some ordered tree for complete string nodes search speedup (on large set) ?
	Childs []*TreeItem // next possible parts slice
}

// GGlobTree is batch glob matcher (dot-separated, like a.b*.c), writted for graphite project
type GGlobTree struct {
	Root       map[int]*TreeItem
	Globs      map[string]int
	GlobsIndex map[int]string
}

func NewTree() *GGlobTree {
	return &GGlobTree{
		Root:       make(map[int]*TreeItem),
		Globs:      make(map[string]int),
		GlobsIndex: make(map[int]string),
	}
}

// func (w *GlobMatcher) Adds(globs []string) (err error) {
// 	var buf strings.Builder
// 	for _, glob := range globs {
// 			buf.Grow(len(glob))
// 			if _, err = w.Add(glob, &buf); err != nil {
// 					return err
// 			}
// 	}
// 	return
// }

// func (gtree *GGlobTree) Add(s string, index int) (normalized string, err error) {
// 	var gg *GGlob
// 	if gg, err = Parse(s); err != nil {
// 		return
// 	}
// }
