package gglob

import (
	"unicode/utf8"
)

// runesExpand expand runes like [a-z0]
func runesExpand(runes []rune) (m map[rune]struct{}, failed bool) {
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
