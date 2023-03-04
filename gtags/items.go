package gtags

import "github.com/msaf1980/go-matcher/pkg/items"

type TaggedItem struct {
	Term *TaggedTerm

	Terminate  bool
	Terminated string // end of chain (resulting seriesByTag)
	TermIndex  int    // resulting seriesByTag index

	Childs []*TaggedItem // next possible parts tree (by key)
}

type TaggedItemList []*TaggedItem

func (item *TaggedItem) Parse(terms TaggedTermList) (lastItem *TaggedItem, err error) {

	for _, child := range item.Childs {
		// TODO: may be normalize parts for equals like {a,z} and {z,a} ?
		if terms[0].Key == child.Term.Key && terms[0].Op == child.Term.Op && terms[0].Value == child.Term.Value {
			lastItem = child
			break
		}
	}

	if lastItem == nil {
		// not found
		// TODO: items caching
		lastItem = &TaggedItem{Term: &terms[0]}
		item.Childs = append(item.Childs, lastItem)
	}

	if len(terms) > 1 {
		if lastItem.Childs == nil {
			lastItem.Childs = make([]*TaggedItem, 0, 4)
		}
		lastItem, err = lastItem.Parse(terms[1:])
	}

	return
}

func (item *TaggedItem) MatchByTagsMap(tags map[string]string, queries *[]string, index *[]int, first items.Store) {
	for _, child := range item.Childs {
		if v, ok := tags[child.Term.Key]; !ok {
			if child.Term.Op == TaggedTermEq || child.Term.Op == TaggedTermMatch {
				// != and ~=! can be skiped and key can not exist, but other not
				continue
			}
		} else {
			if !child.Term.Match(v) {
				continue
			}
		}
		if child.Terminate {
			if queries != nil {
				*queries = append(*queries, child.Terminated)
			}
			if index != nil {
				*index = append(*index, child.TermIndex)
			}
			if first != nil {
				first.Store(child.TermIndex)
			}
		}
		if len(tags) > 0 {
			child.MatchByTagsMap(tags, queries, index, first)
		}
	}
}

func (item *TaggedItem) MatchByTags(tags []Tag, queries *[]string, index *[]int, first items.Store) {
	for _, child := range item.Childs {
		var i int
		tags := tags
		if child.Term.Key != tags[i].Key {
			// scan for tag
			for i = 1; i < len(tags); i++ {
				if tags[i].Key == child.Term.Key {
					break
				}
			}
		}
		if i == len(tags) {
			if child.Term.Op == TaggedTermEq || child.Term.Op == TaggedTermMatch {
				// != and ~=! can be skiped and key can not exist, but other not
				continue
			}
		} else {
			if i > 0 {
				tags = tags[i:]
			}
			if !child.Term.Match(tags[0].Value) {
				continue
			}
		}
		if child.Terminate {
			if queries != nil {
				*queries = append(*queries, child.Terminated)
			}
			if index != nil {
				*index = append(*index, child.TermIndex)
			}
			if first != nil {
				first.Store(child.TermIndex)
			}
		}
		if len(tags) > 0 {
			child.MatchByTags(tags, queries, index, first)
		}
	}
}
