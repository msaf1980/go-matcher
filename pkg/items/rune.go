package items

import "unicode/utf8"

type ItemRune rune

func (item ItemRune) IsString() (string, bool) {
	return "", false
}

func (item ItemRune) IsRune() (rune, bool) {
	return rune(item), true
}

func (ItemRune) CanString() bool {
	return true
}

func (item ItemRune) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		if c == rune(item) {
			found = true
			part = part[n:]
		}
	}
	if found {
		if part != "" && len(nextItems) > 0 {
			found = nextItems[0].Match(part, nextParts, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}
