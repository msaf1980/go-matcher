package items

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/msaf1980/go-matcher/pkg/utils"
)

func interceptionLeft(a []string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	default:
		for i, c := range a[0] {
			for n := 1; n < len(a); n++ {
				if i == len(a[n]) {
					return a[0][:i]
				}
				r, _ := utf8.DecodeRuneInString(a[n][i:])
				if c != r {
					return a[0][:i]
				}
			}
		}
		return a[0]
	}
}

func removeDuplicated(a []string) []string {
	// try to truncate from start
	n := len(a)
	i := 0
	for ; i < n-1; i++ {
		if a[i] != a[i+1] {
			break
		}
	}
	j := n - 1
	for ; j > i; j-- {
		if a[j-1] != a[j] {
			break
		}
	}
	a = a[i : j+1]
	n = len(a)

	// this index will move only when we modify the array in-place to include a new	non-duplicate element.
	j = 0

	for i = 0; i < n; i++ {
		//  If the current element is equal to the next element, then skip the current element because it's a duplicate.
		if i < n-1 && a[i] == a[i+1] {
			continue
		}

		a[j] = a[i]
		j++
	}

	return a[:j]
}

func ListExpand(s string) (list []string, ok bool) {
	last := len(s) - 1
	if len(s) > 1 && s[0] == '{' && s[last] == '}' {
		ok = true
		s = s[1:last]
		if s == "" {
			return
		}
		list = strings.Split(s, ",")
		if len(list) > 0 {
			sort.Strings(list)
			// cleanup duplicated
			list = removeDuplicated(list)
		}
	}

	return
}

// StringList is a alternate list: {a,b,cd}
type StringList struct {
	// nodeList
	ASCIIStarted bool // gready skip scan by first symbol
	FirstASCII   utils.ASCIISet
	Vals         []string // strings
	MinSize      int      // min len in vals or min rune in range, also flag for contain epmty value
	MaxSize      int      // max len in vals or max rune in range
}

func (item *StringList) WriteString(buf *strings.Builder) string {
	l := buf.Len()
	buf.WriteByte('{')
	for i, s := range item.Vals {
		if i > 0 || item.MinSize == 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(s)
	}
	buf.WriteByte('}')
	return buf.String()[l:]
}

func (item *StringList) String() string {
	var buf strings.Builder
	return item.WriteString(&buf)
}

func (item *StringList) MinLen() int {
	return item.MinSize
}

func (item *StringList) MaxLen() int {
	return item.MaxSize
}

func (item *StringList) Find(s string) (index, length int, support FindFlag) {
	support = FindList
	return
}

func (item *StringList) Match(s string) (offset int, support FindFlag) {
	support = FindList
	return
}

func (item *StringList) MatchLast(s string) (offset int, support FindFlag) {
	support = FindNotSupported
	return
}

// IsOptional check when contain empty value and can be skipped
func (item *StringList) IsOptional() bool {
	return item.MinSize == 0
}

// Len is a values count (exclude empty)
func (item *StringList) Len() int {
	return len(item.Vals)
}

// FindFirst is try find first symbol (exclude empty)
func (item *StringList) FindFirst(s string) (index int, supported bool) {
	if item.ASCIIStarted {
		supported = true
		index = item.FirstASCII.Index(s)
	}
	return
}

// MatchFirst is try match first symbol (exclude empty)
func (item *StringList) MatchFirst(s string) (ok, supported bool) {
	if s == "" {
		supported = true
		return
	}
	if item.ASCIIStarted {
		supported = true
		ok = item.FirstASCII.Contains(s[0])
	} else {
		ok = true
	}
	return
}

func (item *StringList) FindN(s string, n int) (index, length int) {
	v := item.Vals[n]
	if index = strings.Index(s, v); index != -1 {
		length = len(v)
	}
	return
}

func (item *StringList) MatchN(s string, n int) (offset int) {
	v := item.Vals[n]
	if strings.HasPrefix(s, v) {
		// strip prefix
		offset = len(v)
	} else {
		offset = -1
	}
	return
}

// func NewItemList return optimized version of InnerItem
func NewItemList(vals []string) (item Item) {
	// TODO: support escape symbols
	if len(vals) == 0 {
		return
	}
	if len(vals) == 1 {
		if vals[0] == "" {
			return nil
		}
		// one item optimization
		c, n := utf8.DecodeRuneInString(vals[0])
		if n == len(vals[0]) {
			if c <= 127 {
				return Byte(c)
			}
			return Rune(c)
		}
		return NewString(vals[0])
	}
	minLen := math.MaxInt
	maxLen := 0

	asciiStarted := true
	var firstASCII utils.ASCIISet
	for _, v := range vals {
		l := len(v)
		if maxLen < l {
			maxLen = l
		}
		if minLen > l {
			minLen = l
		}
		if l > 0 {
			if !firstASCII.Add(v[0]) {
				asciiStarted = false
			}
		}
	}

	if minLen == 0 {
		if vals[0] != "" {
			panic(fmt.Errorf("must be empty values in list: %v", vals))
		}
		vals = vals[1:]
	}

	if asciiStarted {
		item = &StringList{
			Vals: vals, MinSize: minLen, MaxSize: maxLen,
			ASCIIStarted: asciiStarted, FirstASCII: firstASCII,
		}
	} else {
		item = &StringList{Vals: vals, MinSize: minLen, MaxSize: maxLen}
	}

	return
}
