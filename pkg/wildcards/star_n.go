package wildcards

import (
	"strings"
	"unicode/utf8"
)

type ItemNStar int

func (item ItemNStar) Strings() []string {
	return nil
}

func (item ItemNStar) WriteString(buf *strings.Builder) {
	n := int(item)
	buf.WriteRune('*')
	for i := 0; i < n; i++ {
		buf.WriteRune('?')
	}
}

func (item ItemNStar) Type() (typ ItemType, s string, c rune) {
	return ItemTypeOther, "", utf8.RuneError
}

func (item ItemNStar) Locate(part string, nextItems []InnerItem) (offset int, support bool, _ int) {
	return -1, false, 0
}

func (item ItemNStar) Match(part string, nextItems []InnerItem) (found bool) {
	l := len(part)
	min := int(item)
	if l < min {
		return
	}
	if min > 0 {
		var ok bool
		if part, ok = stringAfter(part, int(min)); !ok {
			return
		}
	}

	return matchStar(part, nextItems)
}
