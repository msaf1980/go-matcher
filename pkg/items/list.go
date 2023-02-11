package items

import (
	"math"
	"sort"
	"strings"
	"unicode/utf8"
)

func interceptionLeft(a []string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	default:
		for i, c := range a[0] {
			for n := 1; n < len(a); n++ {
				if i == len(a[n]) {
					return a[0][:i]
				}
				r, _ := utf8.DecodeRuneInString(a[n][i:])
				if c != r {
					return a[0][:i]
				}
			}
		}
		return a[0]
	}
}

func removeDuplicated(a []string) []string {
	// try to truncate from start
	n := len(a)
	i := 0
	for ; i < n-1; i++ {
		if a[i] != a[i+1] {
			break
		}
	}
	j := n - 1
	for ; j > i; j-- {
		if a[j-1] != a[j] {
			break
		}
	}
	a = a[i : j+1]
	n = len(a)

	// this index will move only when we modify the array in-place to include a new	non-duplicate element.
	j = 0

	for i = 0; i < n; i++ {
		//  If the current element is equal to the next element, then skip the current element because it's a duplicate.
		if i < n-1 && a[i] == a[i+1] {
			continue
		}

		a[j] = a[i]
		j++
	}

	return a[:j]
}

func ListExpand(s string) (list []string, failed bool) {
	last := len(s) - 1
	if len(s) > 1 && s[0] == '{' && s[last] == '}' {
		s = s[1:last]
		if s == "" {
			return
		}
		list = strings.Split(s, ",")
		if len(list) > 0 {
			sort.Strings(list)
			// cleanup duplicated
			list = removeDuplicated(list)
		}
	} else {
		failed = true
	}

	return
}

type ItemList struct {
	// nodeList
	FirstRunes map[rune]struct{} // for gready skip scan, if no empty string in list
	Vals       []string          // strings
	ValsMin    int               // min len in vals or min rune in range
	ValsMax    int               // max len in vals or max rune in range
}

// func (*ItemList) Type() NodeType {
// 	return NodeList
// }

func (item *ItemList) IsString() (string, bool) {
	return "", false
}

func (item *ItemList) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	l := len(part)
	if l < item.ValsMin {
		return
	}
	if len(nextItems) == 0 && l > item.ValsMax {
		return
	}
	// TODO: may be optimize scan of duplicate with prefix tree (runes ?) ?
LOOP:
	for _, s := range item.Vals {
		part := part
		if part == s {
			// full match
			found = true
			part = ""
		} else if strings.HasPrefix(part, s) {
			// strip prefix
			found = true
			part = part[len(s):]
		} else {
			// try to next
			continue
		}
		if found {
			if part != "" && len(nextItems) > 0 {
				found = nextItems[0].Match(part, nextParts, nextItems[1:])
			} else if part != "" || len(nextItems) > 0 {
				found = false
			}
			if found {
				break LOOP
			}
		}
	}
	return
}

// LocateFirst find any of first runes pos (call only if ValsMin > 0)
func (item *ItemList) LocateFirst(part string) (offset int) {
	offset = -1
	for i, c := range part {
		if _, ok := item.FirstRunes[c]; !ok {
			offset += i
			return
		}
	}
	return
}

// func NewItemList return optimized version of InnerItem
func NewItemList(vals []string) (item InnerItem, minLen, maxLen int) {
	if len(vals) == 0 {
		return
	}
	if len(vals) == 1 {
		// one item optimization
		return ItemString(vals[0]), len(vals[0]), len(vals[0])
	}
	minLen = math.MaxInt
	maxLen = 0
	for _, v := range vals {
		l := len(v)
		if maxLen < l {
			maxLen = l
		}
		if minLen > l {
			minLen = l
		}
	}
	var firstsRunes map[rune]struct{}
	if minLen > 0 {
		// if no empty string in list
		var firstsRunes = make(map[rune]struct{})
		last := utf8.RuneError
		for _, v := range vals {
			c, _ := utf8.DecodeRuneInString(v)
			if c != last {
				firstsRunes[c] = struct{}{}
				last = c
			}
		}
	}

	item = &ItemList{Vals: vals, FirstRunes: firstsRunes, ValsMin: minLen, ValsMax: maxLen}

	return
}
