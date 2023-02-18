package wildcards

import (
	"io"
	"strings"
)

type ItemType int8

const (
	ItemTypeOther ItemType = iota
	ItemTypeString
	ItemTypeRune
)

type InnerItem interface {
	Type() (typ ItemType, s string, c rune) // return type, and string or rune value (if contain)
	Strings() []string                      // return nil or string values (if contain)
	Match(part string, nextParts string, nextItems []InnerItem) (found bool)
	Locate(part string) (offset int, support bool)
}

// NextWildcardItem extract InnerItem from glob (not regexp)
func NextWildcardItem(s string) (item InnerItem, next string, minLen int, maxLen int, err error) {
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
		runes, failed := RunesExpand(s)
		if failed {
			return nil, s, 0, 0, ErrNodeMissmatch{"rune", s}
		}
		if len(runes) == 0 {
			return nil, next, 0, 0, nil
		}
		if len(runes) == 1 {
			if runes[0].First == runes[0].Last {
				// one item optimization
				return ItemRune(runes[0].First), next, 1, 1, nil
			}
		}
		return ItemRuneRanges(runes), next, 1, 1, nil
	case '{':
		if idx := strings.Index(s, "}"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		vals, failed := ListExpand(s)
		if failed {
			return nil, s, 0, 0, ErrNodeMissmatch{"list", s}
		}
		item, minLen, maxLen = NewItemList(vals)
		return
	case '*':
		var next string
		for i, c := range s {
			if c != '*' {
				next = s[i:]
				break
			}
		}
		return ItemStar{}, next, 0, -1, nil
	case '?':
		next := s[1:]
		return ItemOne{}, next, 1, 1, nil
	case ']', '}':
		return nil, s, 0, 0, ErrNodeUnclosed{s}
	default:
		// string segment
		end := IndexWildcard(s)
		v, next := SplitString(s, end)
		if len(v) == 0 {
			return nil, next, len(v), len(v), nil
		}
		if len(v) == 1 {
			return ItemRune(v[0]), next, len(v), len(v), nil
		}
		return ItemString(v), next, len(v), len(v), nil
	}
}
