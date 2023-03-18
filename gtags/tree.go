package gtags

import (
	"github.com/msaf1980/go-matcher/glob"
	"github.com/msaf1980/go-matcher/pkg/items"
)

// GTagsTree is batch seriesByTag matcher, writted for graphite project (use on large seriesByTag set)
type GTagsTree struct {
	// seriesByTag()
	items.Terminated

	Root       *TaggedItem
	Queries    map[string]int
	QueryIndex map[int]string
}

func NewTree() *GTagsTree {
	return &GTagsTree{
		Root:       new(TaggedItem),
		Queries:    make(map[string]int),
		QueryIndex: make(map[int]string),
	}
}

func (gtree *GTagsTree) Add(queryString string, index int) (normalized string, n int, err error) {
	if queryString == "" {
		return
	}
	if index < 0 {
		err = glob.ErrIndexInvalid
		normalized = queryString
		return
	}
	var ok bool
	if n, ok = gtree.Queries[queryString]; ok {
		// aleady added
		err = glob.ErrGlobExist
		normalized = queryString
		return
	}
	if normalized, ok = gtree.QueryIndex[index]; ok {
		err = glob.ErrIndexDup
		return
	}

	var terms TaggedTermList
	if terms, err = ParseSeriesByTag(queryString); err != nil {
		return
	}
	normalized = terms.String()

	if n, ok = gtree.Queries[normalized]; ok {
		// aleady added
		err = glob.ErrGlobExist
		return
	}

	if len(terms) == 0 {
		gtree.Terminate = true
		gtree.Query = normalized
		gtree.Index = index

	} else {
		lastItem := gtree.Root.Parse(terms, normalized, index)
		lastItem.Terminate = true
		lastItem.Query = normalized
		lastItem.Index = index
	}

	gtree.Queries[queryString] = index
	if normalized != queryString {
		gtree.Queries[normalized] = index
	}
	gtree.QueryIndex[index] = normalized

	n = index

	return
}

func (gtree *GTagsTree) AddTerms(terms TaggedTermList, index int) (normalized string, n int, err error) {
	if len(terms) == 0 {
		return
	}
	if index < 0 {
		err = glob.ErrIndexInvalid
		return
	}
	normalized = terms.String()

	var ok bool
	if n, ok = gtree.Queries[normalized]; ok {
		// aleady added
		err = glob.ErrGlobExist
		return
	}
	if normalized, ok = gtree.QueryIndex[index]; ok {
		err = glob.ErrIndexDup
		return
	}

	if len(terms) == 0 {
		gtree.Terminated = items.Terminated{
			Terminate: true, Index: index, Query: normalized,
		}
	} else {
		gtree.Root.Parse(terms, normalized, index)
	}

	gtree.Queries[normalized] = index
	gtree.QueryIndex[index] = normalized

	n = index

	return
}

func (gtree *GTagsTree) MatchByTagsMap(tags map[string]string, queries *[]string, index *[]int, first items.Store) (matched int) {
	return gtree.Root.MatchByTagsMap(tags, queries, index, first)
}

func (gtree *GTagsTree) MatchByTags(tags []Tag, queries *[]string, index *[]int, first items.Store) (matched int) {
	if gtree.Terminate {
		gtree.Terminated.Append(queries, index, first)
	}
	return gtree.Root.MatchByTags(tags, queries, index, first)
}
