package gtags

import (
	"strings"

	"github.com/msaf1980/go-matcher/pkg/utils"
)

// TagsMatcher is tags matcher, writted for graphite project
type TagsMatcher struct {
	Root    *TaggedItem // by sorted first key (__name__ prefered)
	Queries map[string]int
}

func NewTagsMatcher() *TagsMatcher {
	return &TagsMatcher{
		Root:    &TaggedItem{Childs: make([]*TaggedItem, 0, 8)},
		Queries: make(map[string]int),
	}
}

func (w *TagsMatcher) Adds(queries []string) (err error) {
	var buf strings.Builder
	for _, query := range queries {
		buf.Grow(len(query))
		if _, err = w.Add(query, &buf); err != nil {
			return err
		}
	}
	return
}

func (w *TagsMatcher) Add(query string, buf *strings.Builder) (newQuery string, err error) {
	return w.AddIndexed(query, -1, buf)
}

func (w *TagsMatcher) AddIndexed(query string, termIndex int, buf *strings.Builder) (newQuery string, err error) {
	if query == "" {
		return
	}
	if _, ok := w.Queries[query]; ok {
		// aleady added
		newQuery = query
		return
	}
	var (
		terms TaggedTermList
		item  *TaggedItem
	)
	buf.Reset()
	if terms, err = ParseSeriesByTag(query); err != nil {
		return
	}
	if err = terms.Build(); err != nil {
		return
	}
	terms.Rewrite(buf)

	if item, err = w.Root.Parse(terms); err != nil {
		return
	}
	item.Terminated = append(item.Terminated, query)
	if termIndex > -1 {
		item.TermIndex = append(item.TermIndex, termIndex)
	}

	w.Queries[query] = termIndex

	newQuery = buf.String()
	if newQuery == query {
		newQuery = query
	} else {
		if _, ok := w.Queries[newQuery]; !ok {
			w.Queries[newQuery] = termIndex // write optimized glob
			item.Terminated = append(item.Terminated, newQuery)
			w.Queries[newQuery] = termIndex
		}
	}

	return
}

func (w *TagsMatcher) MatchByTagsMap(tags map[string]string) (queries []string) {
	if len(tags) == 0 {
		return
	}
	queries = make([]string, 0, utils.Min(8, len(w.Root.Childs)))
	w.Root.MatchByTagsMap(tags, &queries)

	return
}

func (w *TagsMatcher) MatchByTagsMapB(tags map[string]string, queries *[]string) {
	// *queries = (*queries)[:0]
	if len(tags) == 0 {
		return
	}
	w.Root.MatchByTagsMap(tags, queries)
}

func (w *TagsMatcher) MatchIndexedByTagsMap(tags map[string]string) (queries []int) {
	if len(tags) == 0 {
		return
	}
	queries = make([]int, 0, utils.Min(8, len(w.Root.Childs)))
	w.Root.MatchIndexedByTagsMap(tags, &queries)

	return
}

func (w *TagsMatcher) MatchIndexedByTagsMapB(tags map[string]string, queries *[]int) {
	// *queries = (*queries)[:0]
	if len(tags) == 0 {
		return
	}
	w.Root.MatchIndexedByTagsMap(tags, queries)
}

func (w *TagsMatcher) MatchFirstByTagsMap(tags map[string]string, queryIndex *int) {
	// *queries = (*queries)[:0]
	if len(tags) == 0 {
		return
	}
	w.Root.MatchFirstByTagsMap(tags, queryIndex)
}

func (w *TagsMatcher) MatchByTags(tags []Tag) (queries []string) {
	if len(tags) == 0 {
		return
	}
	queries = make([]string, 0, utils.Min(8, len(w.Root.Childs)))
	w.Root.MatchByTags(tags, &queries)

	return
}

func (w *TagsMatcher) MatchByTagsB(tags []Tag, queries *[]string) {
	// *queries = (*queries)[:0]
	if len(tags) == 0 {
		return
	}
	w.Root.MatchByTags(tags, queries)
}

func (w *TagsMatcher) MatchIndexedByTags(tags []Tag) (queries []int) {
	if len(tags) == 0 {
		return
	}
	queries = make([]int, 0, utils.Min(8, len(w.Root.Childs)))
	w.Root.MatchIndexedByTags(tags, &queries)

	return
}

func (w *TagsMatcher) MatchIndexedByTagsB(tags []Tag, queries *[]int) {
	// *queries = (*queries)[:0]
	if len(tags) == 0 {
		return
	}
	w.Root.MatchIndexedByTags(tags, queries)
}

func (w *TagsMatcher) MatchFirstByTags(tags []Tag, queryIndex *int) {
	// *queries = (*queries)[:0]
	if len(tags) == 0 {
		return
	}
	w.Root.MatchFirstByTags(tags, queryIndex)
}
