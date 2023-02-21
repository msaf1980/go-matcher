package wildcards

import (
	"strings"
	"unicode/utf8"
)

// ItemOne is a any symbol: ?
type ItemOne struct{}

func (item ItemOne) Type() (typ ItemType, s string, c rune) {
	return ItemTypeOther, "", utf8.RuneError
}

func (item ItemOne) Strings() []string {
	return nil
}

func (item ItemOne) WriteString(buf *strings.Builder) {
	buf.WriteRune('?')
}

func (item ItemOne) Locate(part string, nextItems []InnerItem) (offset int, support bool, skip int) {
	if len(nextItems) == 0 || part == "" {
		return -1, false, 0
	}
	_, n := utf8.DecodeRuneInString(part)
	if n == 0 {
		// failback
		return -1, false, 0
	}
	offset, support, skip = nextItems[0].Locate(part[n:], nextItems[1:])
	if support && offset != -1 {
		skip++
	}
	return
}

func (item ItemOne) Match(part string, nextItems []InnerItem) (found bool) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		found = true
		part = part[n:]

		if len(nextItems) > 0 {
			found = nextItems[0].Match(part, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}
