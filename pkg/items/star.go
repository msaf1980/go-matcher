package items

import (
	"strings"

	"github.com/msaf1980/go-matcher/pkg/utils"
)

// Star is a any runes count (N or greater): *? is minimum 1 rune symbol
type Star int

func (item Star) WriteString(buf *strings.Builder) string {
	n := int(item)
	l := buf.Len()
	buf.Grow(l + n + 1)
	buf.WriteByte('*')
	for i := 0; i < n; i++ {
		buf.WriteByte('?')
	}
	return buf.String()[l:]
}

func (item Star) String() string {
	var buf strings.Builder
	return item.WriteString(&buf)
}

func (item Star) MinLen() int {
	return int(item)
}

func (item Star) MaxLen() int {
	return -1
}

func (item Star) Find(s string) (index, length int, support FindFlag) {
	support = FindStar
	length = utils.StringSkipRunes(s, int(item))
	return
}

func (item Star) Match(s string) (offset int, support FindFlag) {
	support = FindStar
	offset = utils.StringSkipRunes(s, int(item))
	return
}

func (item Star) MatchLast(s string) (offset int, support FindFlag) {
	support = FindNotSupported
	return
}
