package gtags

import (
	"github.com/msaf1980/go-matcher/pkg/items"
)

type TaggedItems struct {
	Key        string
	NotMatched []*TaggedItem
	Matched    []*TaggedItem
}

func (items *TaggedItems) Get(n int) *TaggedItem {
	return items.Matched[n]
}

type TaggedItem struct {
	Term *TaggedTerm

	// seriesByTag()
	items.Terminated

	// TODO: may be scan equal without wildcards with map ?
	Items []TaggedItems // next possible parts tree (by key)
}

func hasName(items []TaggedItems) bool {
	return items[0].Key == "__name__"
}

// findOrAppend search []*TaggedItem child by key (or append them, if not exist)
func (t *TaggedItem) findOrAppend(key string) int {
	if key == "" {
		return 1
	}
	if len(t.Items) == 0 {
		t.Items = make([]TaggedItems, 1, len(t.Items)+4)
		t.Items[0].Key = key
		return 0
	}
	named := hasName(t.Items)
	if key == "__name__" {
		if named {
			return 0
		}
		// ["a", ...]
		newItems := make([]TaggedItems, 1, len(t.Items)+4)
		newItems[0].Key = key
		t.Items = append(newItems, t.Items...)
		return 0
	}
	low := 0
	if named {
		low = 1
	}
	// check __name__ order
	high := len(t.Items) - 1
	for low <= high {
		mid := (low + high) / 2
		// TODO: use strings.Compare or faster analog ?
		if key == t.Items[mid].Key {
			return mid
		}
		if key < t.Items[mid].Key {
			if low == mid {
				// break, we are not find
				high = low
				break
			}
			high = mid - 1
		} else if low == high {
			// break, we are not find, shit insert index
			high = mid + 1
			break
		} else {
			low = mid + 1
		}
	}

	if high == len(t.Items) {
		t.Items = append(t.Items, TaggedItems{Key: key})
		return high
	}
	if high <= 0 {
		newItems := make([]TaggedItems, 1, len(t.Items)+4)
		newItems[0].Key = key
		t.Items = append(newItems, t.Items...)
		return 0
	}
	newItems := make([]TaggedItems, 0, len(t.Items)+4)
	newItems = append(newItems, t.Items[:high]...)
	newItems = append(newItems, TaggedItems{Key: key})
	t.Items = append(newItems, t.Items[high:]...)
	return high
}

// find search by key from start position (except empty or name key)
func (t *TaggedItem) find(key string, start int) int {
	if len(t.Items) == 0 || key == "" {
		return -1
	}
	named := hasName(t.Items)
	if key == "__name__" {
		if named {
			return 0
		}
		return -1
	}
	if len(t.Items) == 0 || key == "" {
		return -1
	}
	low := 0
	if named {
		low = 1
	}
	high := len(t.Items) - 1
	if start > low {
		if t.Items[start].Key == key {
			return start
		}
		if t.Items[start].Key > key {
			high = start - 1
		} else {
			low = start
		}
	}

	search := (low + high) / 2
	if search > low+4 {
		// try ty search from start + 4 instead of midle
		search = low + 4
	}
	for {
		// TODO: use strings.Compare or faster analog ?
		if key == t.Items[search].Key {
			return search
		}
		if key < t.Items[search].Key {
			high = search - 1
		} else {
			low = search + 1
		}
		if low > high {
			break
		}
		if low == high {
			if key == t.Items[low].Key {
				return low
			}
			break
		}
		search = (low + high) / 2
	}

	return -1
}

// find search by key from start position (except empty or name key)
func find(tags []Tag, key string, start int) int {
	if len(tags) == 0 || key == "" {
		return -1
	}
	named := tags[0].Key == "__name__"
	if key == "__name__" {
		if named {
			return 0
		}
		return -1
	}
	low := 0
	if named {
		low = 1
	}
	high := len(tags) - 1
	if start > low {
		if tags[start].Key == key {
			return start
		}
		if tags[start].Key > key {
			high = start - 1
		} else {
			low = start
		}
	}

	search := (low + high) / 2
	if search > low+4 {
		// try ty search from start + 4 instead of midle
		search = low + 4
	}
	for {
		// TODO: use strings.Compare or faster analog ?
		if key == tags[search].Key {
			return search
		}
		if key < tags[search].Key {
			high = search - 1
		} else {
			low = search + 1
		}
		if low > high {
			break
		}
		if low == high {
			if key == tags[low].Key {
				return low
			}
			break
		}
		search = (low + high) / 2
	}

	return -1
}

func (item *TaggedItem) Parse(terms TaggedTermList, query string, index int) (lastItem *TaggedItem) {
	var (
		isMatchedOp bool
		childs      []*TaggedItem
	)
	pos := item.findOrAppend(terms[0].Key)
	switch terms[0].Op {
	case TaggedTermEq, TaggedTermMatch:
		isMatchedOp = true
		childs = item.Items[pos].Matched
	default:
		childs = item.Items[pos].NotMatched
	}
	for _, child := range childs {
		key := ("ab" + string([]byte(child.Term.Key)))[2:]
		if terms[0].Key == key && terms[0].Op == child.Term.Op && terms[0].Value == child.Term.Value {
			lastItem = child
			break
		}
	}

	if lastItem == nil {
		// not found
		// TODO: items caching ?
		lastItem = &TaggedItem{Term: &terms[0]}
		childs = append(childs, lastItem)
		if isMatchedOp {
			item.Items[pos].Matched = childs
		} else {
			item.Items[pos].NotMatched = childs
		}
	}

	if len(terms) > 1 {
		lastItem = lastItem.Parse(terms[1:], query, index)
	}

	return
}

func (item *TaggedItem) MatchByTagsMap(tags map[string]string, queries *[]string, index *[]int, first items.Store) (matched int) {
	if len(tags) == 0 {
		return
	}

	for i := 0; i < len(item.Items); i++ {
		v, ok := tags[item.Items[i].Key]
		if ok {
			for _, child := range item.Items[i].Matched {
				if !child.Term.Match(v) {
					continue
				}
				if child.Terminate {
					child.Append(queries, index, first)
					matched++
				}
				if n := child.MatchByTagsMap(tags, queries, index, first); n > 0 {
					matched += n
				}
			}
			for _, child := range item.Items[i].NotMatched {
				if !child.Term.Match(v) {
					continue
				}
				if child.Terminate {
					child.Append(queries, index, first)
					matched++
				}
				if n := child.MatchByTagsMap(tags, queries, index, first); n > 0 {
					matched += n
				}
			}
		} else {
			// tags not exist, check not matched
			for _, child := range item.Items[i].NotMatched {
				if child.Terminate {
					child.Append(queries, index, first)
					matched++
				}
				if n := child.MatchByTagsMap(tags, queries, index, first); n > 0 {
					matched += n
				}
			}
		}

	}

	return
}

func (item *TaggedItem) MatchByTags(tags []Tag, queries *[]string, index *[]int, first items.Store) (matched int) {
	if len(tags) == 0 {
		return
	}

	matchPos := 0

	for i := 0; i < len(item.Items); i++ {
		n := find(tags, item.Items[i].Key, matchPos)
		if n == -1 {
			// tags not exist, check not matched
			for _, child := range item.Items[i].NotMatched {
				if child.Terminate {
					child.Append(queries, index, first)
					matched++
				}
				if n := child.MatchByTags(tags, queries, index, first); n > 0 {
					matched += n
				}
			}
		} else {
			matchPos = n
			for _, child := range item.Items[i].Matched {
				if !child.Term.Match(tags[n].Value) {
					continue
				}
				if child.Terminate {
					child.Append(queries, index, first)
					matched++
				}
				if n := child.MatchByTags(tags, queries, index, first); n > 0 {
					matched += n
				}
			}
			for _, child := range item.Items[i].NotMatched {
				if !child.Term.Match(tags[n].Value) {
					continue
				}
				if child.Terminate {
					child.Append(queries, index, first)
					matched++
				}
				if n := child.MatchByTags(tags, queries, index, first); n > 0 {
					matched += n
				}
			}
		}

	}

	return
}
