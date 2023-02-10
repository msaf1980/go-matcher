package gtags

import (
	"github.com/msaf1980/go-matcher/pkg/utils"
)

// TagsMatcher is tags matcher, writted for graphite project
type TagsMatcher struct {
	Root    *TaggedItem // by sorted first key (__name__ prefered)
	Queries map[string]bool
}

func NewTagsMatcher() *TagsMatcher {
	return &TagsMatcher{
		Root:    &TaggedItem{Childs: make([]*TaggedItem, 0, 8)},
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
	var (
		terms TaggedTermList
		item  *TaggedItem
	)
	if terms, err = ParseSeriesByTag(query); err != nil {
		return err
	}
	if item, err = w.Root.Parse(terms); err != nil {
		return err
	}

	item.Terminated = append(item.Terminated, query)

	w.Queries[query] = true

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
