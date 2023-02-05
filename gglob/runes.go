package gglob

import (
	"unicode/utf8"
)

// runesExpand expand runes like [a-z0]
func runesExpand(runes []rune) (m map[rune]struct{}) {
	var r rune
	m = make(map[rune]struct{})
	if len(runes) >= 3 && runes[0] == '[' && runes[len(runes)-1] == ']' {
		runes = runes[1 : len(runes)-1]
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
	}

	return
}
