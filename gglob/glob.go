package gglob

import (
	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/utils"
)

func ParseItems(root map[int]*items.NodeItem, glob string) (lastNode *items.NodeItem, err error) {
	glob, partsCount := items.PathLevel(glob)

	node, ok := root[partsCount]
	if !ok {
		node = &items.NodeItem{}
		root[partsCount] = node
	}
	_, err = node.Parse(glob, partsCount)

	return
}

// GlobMatcher is dotted-separated segment glob matcher, like a.b.[c-e]?.{f-o}*, writted for graphite project
type GlobMatcher struct {
	Root  map[int]*items.NodeItem
	Globs map[string]bool
}

func NewGlobMatcher() *GlobMatcher {
	return &GlobMatcher{
		Root:  make(map[int]*items.NodeItem),
		Globs: make(map[string]bool),
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
	if w.Globs[glob] {
		// aleady added
		return
	}
	if _, err = ParseItems(w.Root, glob); err != nil {
		return err
	}

	w.Globs[glob] = true

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
	if path == "" {
		return
	}
	*globs = (*globs)[:0]
	path, partsCount := items.PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		node.Match(path, globs)
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
	if len(parts) == 0 {
		return
	}
	*globs = (*globs)[:0]
	if node, ok := w.Root[len(parts)]; ok {
		node.MatchByParts(parts, globs)
	}
}
