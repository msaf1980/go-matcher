package wildcards

import (
	"strings"
	"unicode/utf8"
)

// ItemOne is a n any symbols: ??? for 3 symbols
type ItemMany int

func (item ItemMany) Type() (typ ItemType, s string, c rune) {
	return ItemTypeOther, "", utf8.RuneError
}

func (item ItemMany) Strings() []string {
	return nil
}

func (item ItemMany) WriteString(buf *strings.Builder) {
	n := int(item)
	for i := 0; i < n; i++ {
		buf.WriteRune('?')
	}
}

func stringAfter(s string, length int) (string, bool) {
	for i := 0; i < length; i++ {
		_, n := utf8.DecodeRuneInString(s)
		if n == 0 {
			// failback
			return "", false
		}
		s = s[n:]
	}
	return s, true
}

func (item ItemMany) Locate(part string, nextItems []InnerItem) (offset int, support bool, skip int) {
	if len(nextItems) == 0 {
		return -1, false, 0
	}

	l := len(part)
	min := int(item)
	if l < min {
		return
	}
	var ok bool
	if part, ok = stringAfter(part, int(min)); !ok {
		return
	}

	offset, support, skip = nextItems[0].Locate(part, nextItems[1:])
	if support && offset != -1 {
		skip++
	}
	return
}

func (item ItemMany) Match(part string, nextItems []InnerItem) (found bool) {
	var ok bool
	if part, ok = stringAfter(part, int(item)); !ok {
		return false
	}

	if len(nextItems) > 0 {
		found = nextItems[0].Match(part, nextItems[1:])
	} else if part == "" && len(nextItems) == 0 {
		found = true
	}

	return
}
