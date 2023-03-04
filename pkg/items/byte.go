package items

import (
	"strings"
)

// Byte is a unicode symbol
type Byte byte

func (item Byte) Equal(a Item) bool {
	if v, ok := a.(Byte); ok {
		return item == v
	}
	return false
}

func (item Byte) WriteRandom(buf *strings.Builder) {
	buf.WriteByte(byte(item))
}

func (item Byte) WriteString(buf *strings.Builder) string {
	l := buf.Len()
	buf.WriteByte(byte(item))
	return buf.String()[l:]
}

func (item Byte) String() string {
	var buf strings.Builder
	return item.WriteString(&buf)
}

func (item Byte) Append(s string) *String {
	return &String{S: string(item) + s}
}

func (item Byte) AppendByte(b byte) *String {
	return &String{S: string(item) + string(b)}
}

func (item Byte) AppendRune(r rune) *String {
	return &String{S: string(item) + string(r)}
}

func (item Byte) MinLen() int {
	return 1
}

func (item Byte) MaxLen() int {
	return 1
}

func (item Byte) Find(s string) (index, length int, support FindFlag) {
	c := byte(item)
	if index = strings.IndexByte(s, c); index != -1 {
		length = 1
	}
	return
}

func (item Byte) Match(s string) (offset int, support FindFlag) {
	if s == "" {
		offset = -1
	} else {
		c := byte(item)
		if c == s[0] {
			offset = 1
		} else {
			offset = -1
		}
	}
	return
}

func (item Byte) MatchLast(s string) (offset int, support FindFlag) {
	n := len(s) - 1
	if c := s[n]; c != byte(item) {
		offset = -1
	}
	return
}
