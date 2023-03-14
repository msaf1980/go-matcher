package items

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/msaf1980/go-matcher/pkg/utils"
)

// NextWildcardItem extract InnerItem from glob (not regexp)
func NextWildcardItem(s string) (item Item, next string, err error) {
	// TODO: implement escape symbols
	if s == "" {
		return nil, s, io.EOF
	}
	switch s[0] {
	case '[':
		if idx := strings.Index(s, "]"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		runes, ok := utils.RunesRangeExpand(s)
		if !ok {
			return nil, s, ErrNodeMissmatch{"rune", s}
		}
		n, c := runes.ASCII.Count()
		if n == 0 && len(runes.UnicodeRanges) == 0 {
			return nil, next, nil
		}
		if len(runes.UnicodeRanges) == 0 {
			if n == 1 {
				return Byte(c), next, nil
			}
		}
		if len(runes.UnicodeRanges) == 1 && n == 0 {
			if runes.UnicodeRanges[0].First == runes.UnicodeRanges[0].Last {
				// one item optimization
				return Rune(runes.UnicodeRanges[0].First), next, nil
			}
		}
		r := &RunesRanges{RunesRanges: runes}
		return r, next, nil
	case '{':
		if idx := strings.Index(s, "}"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		vals, ok := ListExpand(s)
		if !ok {
			return nil, s, ErrNodeMissmatch{"list", s}
		}
		item, err = NewItemList(vals)
		return
	case '*':
		var next string
		for i, c := range s {
			if c != '*' {
				next = s[i:]
				break
			}
		}
		return Star(0), next, nil
	case '?':
		next := s[1:]
		return Any(1), next, nil
	case ']', '}':
		return nil, s, ErrNodeUnclosed{s}
	default:
		// string segment
		end := IndexWildcard(s)
		if end == -1 {
			return NewString(s), next, nil
		}
		v, next := utils.SplitString(s, end)

		c, n := utf8.DecodeRuneInString(v)
		if n == len(v) {
			if c <= 127 {
				return Byte(c), next, nil
			}
			return Rune(c), next, nil
		}
		return NewString(v), next, nil
	}
}
