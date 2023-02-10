package items

import (
	"unicode/utf8"
)

// RunesExpand expand runes like [a-z0]
func RunesExpand(runes []rune) (m map[rune]struct{}, failed bool) {
	var r rune
	if len(runes) > 1 && runes[0] == '[' && runes[len(runes)-1] == ']' {
		m = make(map[rune]struct{})
		runes = runes[1 : len(runes)-1]
		if len(runes) == 0 {
			return
		}
		start := utf8.RuneError
		for i := 0; i < len(runes); i++ {
			if runes[i] == '-' {
				i++
				if i >= len(runes) {
					break
				}
				if start != utf8.RuneError {
					for r = start + 1; r <= runes[i] && r != utf8.RuneError; r++ {
						m[r] = struct{}{}
					}
					start = runes[i]
					continue
				}
			}
			r = runes[i]
			m[r] = struct{}{}
			start = r
		}
	} else {
		failed = true
	}

	return
}

type ItemRune map[rune]struct{}

// func (ItemRune) Type() NodeType {
// 	return NodeRune
// }

func (item ItemRune) IsString() (string, bool) {
	return "", false
}

func (item ItemRune) Match(part string, nextParts string, nextItems []InnerItem, _ bool) (found bool, _ int) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		if _, ok := item[c]; ok {
			found = true
			part = part[n:]
		}
	}
	if found {
		if part != "" && len(nextItems) > 0 {
			found, _ = nextItems[0].Match(part, nextParts, nextItems[1:], false)
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}
