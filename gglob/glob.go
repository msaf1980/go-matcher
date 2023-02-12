package gglob

import (
	"sort"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/utils"
)

func ParseItems(root map[int]*items.NodeItem, glob string, termIdx int) (lastNode *items.NodeItem, err error) {
	glob, partsCount := items.PathLevel(glob)

	node, ok := root[partsCount]
	if !ok {
		node = &items.NodeItem{}
		root[partsCount] = node
	}
	_, err = node.Parse(glob, partsCount, termIdx)

	return
}

// GlobMatcher is dotted-separated segment glob matcher, like a.b.[c-e]?.{f-o}*, writted for graphite project
type GlobMatcher struct {
	Root  map[int]*items.NodeItem
	Globs map[string]int
}

func NewGlobMatcher() *GlobMatcher {
	return &GlobMatcher{
		Root:  make(map[int]*items.NodeItem),
		Globs: make(map[string]int),
	}
}

func (w *GlobMatcher) Adds(globs []string) (err error) {
	for _, glob := range globs {
		if err = w.Add(glob); err != nil {
			return err
		}
	}
	return
}

func (w *GlobMatcher) Add(glob string) (err error) {
	if glob == "" {
		return
	}
	if _, ok := w.Globs[glob]; ok {
		// aleady added
		return
	}
	if _, err = ParseItems(w.Root, glob, -1); err != nil {
		return err
	}

	w.Globs[glob] = -1

	return
}

func (w *GlobMatcher) AddIndexed(glob string, termIdx int) (err error) {
	if glob == "" {
		return
	}
	if _, ok := w.Globs[glob]; ok {
		// aleady added
		return
	}
	if _, err = ParseItems(w.Root, glob, termIdx); err != nil {
		return err
	}

	w.Globs[glob] = termIdx

	return
}

func (w *GlobMatcher) Match(path string) (globs []string) {
	if path == "" {
		return nil
	}
	path, partsCount := items.PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		globs = make([]string, 0, utils.Min(4, len(node.Childs)))
		node.Match(path, &globs)
	}

	return
}

func (w *GlobMatcher) MatchB(path string, globs *[]string) {
	// *globs = (*globs)[:0]
	if path == "" {
		return
	}
	path, partsCount := items.PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		node.Match(path, globs)
	}
}

func (w *GlobMatcher) MatchIndexed(path string) (globs []int) {
	if path == "" {
		return nil
	}
	path, partsCount := items.PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		globs = make([]int, 0, utils.Min(4, len(node.Childs)))
		node.MatchIndexed(path, &globs)
	}

	return
}

func (w *GlobMatcher) MatchIndexedB(path string, globs *[]int) {
	// *globs = (*globs)[:0]
	if path == "" {
		return
	}
	path, partsCount := items.PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		node.MatchIndexed(path, globs)
		sort.Ints(*globs)
	}
}

func (w *GlobMatcher) MatchFirst(path string, globIndex *int) {
	// *globs = (*globs)[:0]
	if path == "" {
		return
	}
	path, partsCount := items.PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		node.MatchFirst(path, globIndex)
	}
}

func (w *GlobMatcher) MatchByParts(parts []string) (globs []string) {
	if len(parts) == 0 {
		return nil
	}
	if node, ok := w.Root[len(parts)]; ok {
		globs = make([]string, 0, utils.Min(4, len(node.Childs)))
		node.MatchByParts(parts, &globs)
	}

	return
}

func (w *GlobMatcher) MatchByPartsB(parts []string, globs *[]string) {
	// *globs = (*globs)[:0]
	if len(parts) == 0 {
		return
	}
	if node, ok := w.Root[len(parts)]; ok {
		node.MatchByParts(parts, globs)
	}
}

func (w *GlobMatcher) MatchIndexedByParts(parts []string) (globsIndexed []int) {
	if len(parts) == 0 {
		return nil
	}
	if node, ok := w.Root[len(parts)]; ok {
		globsIndexed = make([]int, 0, utils.Min(4, len(node.Childs)))
		node.MatchIndexedByParts(parts, &globsIndexed)
		sort.Ints(globsIndexed)
	}

	return
}

func (w *GlobMatcher) MatchIndexedByPartsB(parts []string, globs *[]int) {
	// *globs = (*globs)[:0]
	if len(parts) == 0 {
		return
	}
	if node, ok := w.Root[len(parts)]; ok {
		node.MatchIndexedByParts(parts, globs)
		sort.Ints(*globs)
	}
}

func (w *GlobMatcher) MatchFirstByParts(parts []string, globIndex *int) {
	// *globs = (*globs)[:0]
	if len(parts) == 0 {
		return
	}
	if node, ok := w.Root[len(parts)]; ok {
		node.MatchFirstByParts(parts, globIndex)
	}
}
