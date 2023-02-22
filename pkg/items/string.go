package items

import (
	"strings"
	"unicode/utf8"
)

type ItemString string

func (item ItemString) Type() (typ ItemType, s string, c rune) {
	return ItemTypeString, string(item), utf8.RuneError
}

func (item ItemString) Strings() []string {
	return nil
}

func (item ItemString) WriteString(buf *strings.Builder) {
	buf.WriteString(string(item))
}

func (item ItemString) Locate(part string, nextItems []Item) (offset int, support bool, _ int) {
	s := string(item)
	support = true
	if offset = strings.Index(part, s); offset != -1 {
		offset += len(s)
	}
	return
}

func (item ItemString) Match(part string, nextItems []Item) (found bool) {
	s := string(item)
	if strings.HasPrefix(part, s) {
		// strip prefix
		found = true
		part = part[len(s):]

		if len(nextItems) > 0 {
			found = nextItems[0].Match(part, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}
