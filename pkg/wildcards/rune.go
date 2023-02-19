package wildcards

import (
	"strings"
	"unicode/utf8"
)

type ItemRune rune

func (item ItemRune) Type() (typ ItemType, s string, c rune) {
	return ItemTypeRune, "", rune(item)
}

func (item ItemRune) Strings() []string {
	return nil
}

func (item ItemRune) Locate(part string, nextItems []InnerItem) (offset int, support bool, _ int) {
	support = true
	c := rune(item)
	if offset = strings.IndexRune(part, c); offset != -1 {
		offset += 1
	}
	return
}

func (item ItemRune) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		if c == rune(item) {
			found = true
			part = part[n:]
		}
	}
	if found {
		if len(nextItems) > 0 {
			found = nextItems[0].Match(part, nextParts, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}
