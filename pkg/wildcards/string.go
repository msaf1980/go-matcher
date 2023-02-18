package wildcards

import (
	"strings"
	"unicode/utf8"
)

type ItemString string

func (ItemString) CanEmpty() bool {
	return false
}

func (item ItemString) IsRune() (rune, bool) {
	return utf8.RuneError, false
}

func (item ItemString) IsString() (string, bool) {
	return string(item), true
}

func (ItemString) CanString() bool {
	return true
}

func (item ItemString) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	s := string(item)
	if strings.HasPrefix(part, s) {
		// strip prefix
		found = true
		part = part[len(s):]

		if len(nextItems) > 0 {
			found = nextItems[0].Match(part, nextParts, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}

func (item ItemString) Locate(part string) (offset int, found bool) {
	s := string(item)
	if offset = strings.Index(part, s); offset != -1 {
		offset += len(s)
		found = true
	}
	return
}
