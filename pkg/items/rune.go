package items

import (
	"strings"
	"unicode/utf8"
)

// Rune is a unicode symbol
type Rune rune

func (item Rune) Equal(a Item) bool {
	if v, ok := a.(Rune); ok {
		return item == v
	}
	return false
}

func (item Rune) WriteString(buf *strings.Builder) string {
	l := buf.Len()
	buf.WriteRune(rune(item))
	return buf.String()[l:]
}

func (item Rune) String() string {
	var buf strings.Builder
	return item.WriteString(&buf)
}

func (item Rune) Append(s string) *String {
	return &String{S: string(item) + s}
}

func (item Rune) AppendByte(b byte) *String {
	return &String{S: string(item) + string(b)}
}

func (item Rune) AppendRune(r rune) *String {
	return &String{S: string(item) + string(r)}
}

func (item Rune) MinLen() int {
	return utf8.RuneLen(rune(item))
}

func (item Rune) MaxLen() int {
	return utf8.RuneLen(rune(item))
}

func (item Rune) Find(s string) (index, length int, support FindFlag) {
	c := rune(item)
	if index = strings.IndexRune(s, c); index != -1 {
		length = utf8.RuneLen(c)
	}
	return
}

func (item Rune) Match(s string) (offset int, support FindFlag) {
	if c, n := utf8.DecodeRuneInString(s); c != utf8.RuneError {
		if c == rune(item) {
			offset = n
		} else {
			offset = -1
		}
	} else {
		offset = -1
	}
	return
}
