package items

import (
	"strings"

	"github.com/msaf1980/go-matcher/pkg/utils"
)

// ItemOne is a n any symbols: wildcard ??? for 3 symbols
type Any int

func (item Any) Equal(a Item) bool {
	if v, ok := a.(Any); ok {
		return item == v
	}
	return false
}

func (item Any) WriteString(buf *strings.Builder) string {
	n := int(item)
	l := buf.Len()
	buf.Grow(l + n)
	for i := 0; i < n; i++ {
		buf.WriteByte('?')
	}
	return buf.String()[l:]
}

func (item Any) String() string {
	var buf strings.Builder
	return item.WriteString(&buf)
}

func (item Any) MinLen() int {
	return int(item)
}

func (item Any) MaxLen() int {
	return int(item) * 4
}

func (item Any) Find(s string) (index, length int, support FindFlag) {
	support = FindForwarded
	n := int(item)
	if length = utils.StringSkipRunes(s, n); length == -1 {
		index = -1
	}
	return
}

func (item Any) Match(s string) (offset int, support FindFlag) {
	offset = utils.StringSkipRunes(s, int(item))
	return
}
