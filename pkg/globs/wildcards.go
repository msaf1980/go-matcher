package globs

import (
	"io"
	"strings"

	"github.com/msaf1980/go-matcher/pkg/items"
)

// NextWildcardItem extract InnerItem from glob (not regexp)
func NextWildcardItem(s string) (item items.Item, next string, minLen int, maxLen int, err error) {
	if s == "" {
		return nil, s, 0, 0, io.EOF
	}
	switch s[0] {
	case '[':
		if idx := strings.Index(s, "]"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		runes, failed := items.RunesExpand(s)
		if failed {
			return nil, s, 0, 0, items.ErrNodeMissmatch{"rune", s}
		}
		if len(runes) == 0 {
			return nil, next, 0, 0, nil
		}
		if len(runes) == 1 {
			if runes[0].First == runes[0].Last {
				// one item optimization
				return items.ItemRune(runes[0].First), next, 1, 1, nil
			}
		}
		return items.ItemRuneRanges(runes), next, 1, 1, nil
	case '{':
		if idx := strings.Index(s, "}"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		vals, failed := items.ListExpand(s)
		if failed {
			return nil, s, 0, 0, items.ErrNodeMissmatch{"list", s}
		}
		item, minLen, maxLen = items.NewItemList(vals)
		return
	case '*':
		var next string
		for i, c := range s {
			if c != '*' {
				next = s[i:]
				break
			}
		}
		return items.ItemStar{}, next, 0, -1, nil
	case '?':
		next := s[1:]
		return items.ItemOne{}, next, 1, 1, nil
	case ']', '}':
		return nil, s, 0, 0, items.ErrNodeUnclosed{s}
	default:
		// string segment
		end := items.IndexWildcard(s)
		if end == -1 {
			return items.ItemString(s), next, len(s), len(s), nil
		}
		v, next := items.SplitString(s, end)
		if len(v) == 1 {
			return items.ItemRune(v[0]), next, len(v), len(v), nil
		}
		return items.ItemString(v), next, len(v), len(v), nil
	}
}
