package items

import (
	"strings"
)

type String struct {
	S string
}

func NewString(s string) Item {
	return &String{S: s}
}

func (item *String) WriteString(buf *strings.Builder) string {
	l := buf.Len()
	buf.WriteString(item.S)
	return buf.String()[l:]
}

func (item *String) String() string {
	var buf strings.Builder
	return item.WriteString(&buf)
}

func (item *String) Add(s string) {
	item.S += s
}

func (item *String) Prepend(s string) {
	item.S = s + item.S
}

func (item *String) AddByte(b byte) {
	item.S += string(b)
}

func (item *String) PrependByte(b byte) {
	item.S = string(b) + item.S
}

func (item *String) AddRune(r rune) {
	item.S += string(r)
}

func (item *String) PrependRune(r rune) {
	item.S = string(r) + item.S
}

func (item *String) MinLen() int {
	return len(item.S)
}

func (item *String) MaxLen() int {
	return len(item.S)
}

func (item *String) Find(s string) (index, length int, support FindFlag) {
	if index = strings.Index(s, item.S); index != -1 {
		length = len(item.S)
	}
	return
}

func (item *String) Match(s string) (offset int, support FindFlag) {
	if strings.HasPrefix(s, item.S) {
		// strip prefix
		offset = len(item.S)
	} else {
		offset = -1
	}
	return
}

func (item *String) MatchLast(s string) (offset int, support FindFlag) {
	if strings.HasSuffix(s, item.S) {
		// strip suffix
		offset = len(s) - len(item.S)
	} else {
		offset = -1
	}
	return
}
