package gtags

import (
	"github.com/msaf1980/go-matcher/pkg/utils"
)

type TagsItem struct {
	Query string // seriesByTag
	Terms TaggedTermList
}

// TagsMatcher is tags matcher, writted for graphite project
type TagsMatcher struct {
	Root    []*TagsItem
	Queries map[string]bool
}

func NewTagsMatcher() *TagsMatcher {
	return &TagsMatcher{
		Root:    make([]*TagsItem, 0),
		Queries: make(map[string]bool),
	}
}

func (w *TagsMatcher) Adds(queries []string) (err error) {
	for _, query := range queries {
		if err = w.Add(query); err != nil {
			return err
		}
	}
	return
}

func (w *TagsMatcher) Add(query string) (err error) {
	if query == "" {
		return
	}
	if w.Queries[query] {
		// aleady added
		return
	}
	var terms TaggedTermList
	if terms, err = ParseSeriesByTag(query); err != nil {
		return err
	}
	w.Root = append(w.Root, &TagsItem{Query: query, Terms: terms})

	w.Queries[query] = true

	return
}

func (w *TagsMatcher) MatchByTags(tags map[string]string) (queries []string) {
	if len(tags) == 0 {
		return
	}
	queries = make([]string, 0, utils.Min(4, len(w.Root)))
	w.MatchByTagsB(tags, &queries)

	return
}

func (w *TagsMatcher) MatchByTagsB(tags map[string]string, queries *[]string) {
	*queries = (*queries)[:0]
	for _, terms := range w.Root {
		if terms.Terms.MatchByTags(tags) {
			*queries = append(*queries, terms.Query)
		}
	}
}

func (w *TagsMatcher) MatchByPath(path string) (queries []string) {
	if len(path) == 0 {
		return
	}
	queries = make([]string, 0, utils.Min(4, len(w.Root)))
	w.MatchByPathB(path, &queries)

	return
}

func (w *TagsMatcher) MatchByPathB(path string, queries *[]string) {
	*queries = (*queries)[:0]
	for _, terms := range w.Root {
		if terms.Terms.MatchByPath(path) {
			*queries = append(*queries, terms.Query)
		}
	}
}
