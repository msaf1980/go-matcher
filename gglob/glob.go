package gglob

import (
	"sort"
	"strings"

	"github.com/msaf1980/go-matcher/pkg/utils"
	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

// GlobMatcher is dotted-separated segment glob matcher, like a.b.[c-e]?.{f-o}*, writted for graphite project
type GlobMatcher struct {
	Root  map[int]*NodeItem
	Globs map[string]int
}

func NewGlobMatcher() *GlobMatcher {
	return &GlobMatcher{
		Root:  make(map[int]*NodeItem),
		Globs: make(map[string]int),
	}
}

func (w *GlobMatcher) Adds(globs []string) (err error) {
	var buf strings.Builder
	for _, glob := range globs {
		buf.Grow(len(glob))
		if _, err = w.Add(glob, &buf); err != nil {
			return err
		}
	}
	return
}

func (w *GlobMatcher) Add(glob string, buf *strings.Builder) (newGlob string, err error) {
	return w.AddIndexed(glob, -1, buf)
}

func (w *GlobMatcher) AddIndexed(glob string, termIdx int, buf *strings.Builder) (newGlob string, err error) {
	if glob == "" {
		return
	}
	if _, ok := w.Globs[glob]; ok {
		// aleady added
		newGlob = glob
		return
	}
	if newGlob, _, err = w.parseItems(glob, termIdx, buf); err != nil {
		return
	}

	return
}

func (w *GlobMatcher) parseItems(glob string, termIdx int, buf *strings.Builder) (newGlob string, lastNode *NodeItem, err error) {
	glob, partsCount := wildcards.PathLevel(glob)

	node, ok := w.Root[partsCount]
	if !ok {
		node = &NodeItem{}
		w.Root[partsCount] = node
	}
	buf.Reset()

	newGlob, lastNode, err = node.ParseNode(glob, termIdx, buf)

	w.Globs[glob] = termIdx
	if glob != newGlob {
		if _, ok := w.Globs[newGlob]; !ok {
			w.Globs[newGlob] = termIdx // write optimized glob
		}
	}

	return
}

func (w *GlobMatcher) Match(path string) (globs []string) {
	if path == "" {
		return nil
	}
	path, partsCount := wildcards.PathLevel(path)
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
	path, partsCount := wildcards.PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		node.Match(path, globs)
	}
}

func (w *GlobMatcher) MatchIndexed(path string) (globs []int) {
	if path == "" {
		return nil
	}
	path, partsCount := wildcards.PathLevel(path)
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
	path, partsCount := wildcards.PathLevel(path)
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
	path, partsCount := wildcards.PathLevel(path)
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
