package items

import (
	"sort"
	"strings"
)

func interception(a []string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	default:
		for i := range a[0] {
			for n := 1; n < len(a); n++ {
				if i == len(a[n]) || a[0][i] != a[n][i] {
					return a[0][:i]
				}
			}
		}
		return a[0]
	}
}

func removeDuplicated(a []string) []string {
	n := len(a)
	i := 0
	for ; i < n-1; i++ {
		if a[i] != a[i+1] {
			break
		}
	}
	if i < n {
		a = a[i:]
		n = len(a)
	}

	// this index will move only when we modify the array in-place to include a new	non-duplicate element.
	j := 0

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
			// remove empty string  from start
			if list[0] == "" {
				list = list[1:]
			}

		}
	} else {
		failed = true
	}

	return
}

type ItemList struct {
	// nodeList
	Vals    []string // strings
	ValsMin int      // min len in vals or min rune in range
	ValsMax int      // max len in vals or max rune in range
}

// func (*ItemList) Type() NodeType {
// 	return NodeList
// }

func (item *ItemList) IsString() (string, bool) {
	return "", false
}

func (item *ItemList) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	// TODO: nodeList Skip scan
	l := len(part)
	if l < item.ValsMin {
		return
	}
	if len(nextItems) == 0 && l > item.ValsMax {
		return
	}
	// TODO: may be optimize scan of duplicate with prefix tree ?
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
