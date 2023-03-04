package items

import (
	"math/rand"
	"strings"
)

// Chain is sequentional of items
type Chain struct {
	Items   []Item
	MinSize int
	MaxSize int
}

func NewChain() *Chain {
	return &Chain{Items: make([]Item, 0, 2)}
}

func (item *Chain) Equal(a Item) bool {
	if v, ok := a.(*Chain); ok {
		if item.MinSize != v.MinSize || item.MaxSize != v.MaxSize {
			return false
		}
		return ItemsEqual(item.Items, v.Items)
	}
	return false
}

func (c *Chain) Append(item Item) {
	c.MinSize += item.MinLen()
	if l := item.MaxLen(); l == -1 {
		c.MaxSize = -1
	} else {
		c.MaxSize += l
	}
	c.Items = AppendItem(c.Items, item)
}

func (c *Chain) WriteRandom(buf *strings.Builder) {
	for i := 0; i < len(c.Items); i++ {
		c.Items[i].WriteRandom(buf)
	}
}

func (c *Chain) WriteString(buf *strings.Builder) string {
	l := buf.Len()
	for _, v := range c.Items {
		v.WriteString(buf)
	}
	return buf.String()[l:]
}

func (c *Chain) String() string {
	var buf strings.Builder
	return c.WriteString(&buf)
}

func (c *Chain) MinLen() int {
	return c.MinSize
}

func (c *Chain) MaxLen() int {
	return c.MaxSize
}

func (c *Chain) Find(s string) (index, length int, support FindFlag) {
	support = FindNotSupported
	// TODO
	return
}

func (c *Chain) Match(s string) (offset int, support FindFlag) {
	support = FindNotSupported
	// TODO
	return
}

func (c *Chain) MatchLast(s string) (offset int, support FindFlag) {
	support = FindNotSupported
	return
}

// List is list of items
type List struct {
	Vals    []Item
	MinSize int
	MaxSize int
}

func (item *List) Equal(a Item) bool {
	if v, ok := a.(*List); ok {
		return item.Equal(v)
	}
	return false
}

func (item *List) WriteRandom(buf *strings.Builder) {
	item.Vals[rand.Intn(len(item.Vals))].WriteRandom(buf)
}

func (item *List) WriteString(buf *strings.Builder) string {
	l := buf.Len()
	buf.WriteByte('{')
	for i, v := range item.Vals {
		if i > 0 {
			buf.WriteByte(',')
		}
		v.WriteString(buf)
	}
	buf.WriteByte('}')
	return buf.String()[l:]
}

func (item *List) String() string {
	var buf strings.Builder
	return item.WriteString(&buf)
}

func (item *List) MinLen() int {
	return item.MinSize
}

func (item *List) MaxLen() int {
	return item.MaxSize
}

func (item *List) Find(s string) (index, length int, support FindFlag) {
	support = FindGroup
	return
}

func (item *List) Match(s string) (offset int, support FindFlag) {
	support = FindGroup
	return
}

func (item *List) MatchLast(s string) (offset int, support FindFlag) {
	support = FindNotSupported
	return
}

func NewList(vals []string) Item {
	items := make([]Item, 0, len(vals))
	for i := 0; i < len(vals); i++ {
		if HasStarAny(vals[i]) {
			v := vals[i]
			if v == "*" {
				return Star(0)
			} else if v == "?" {
				items = append(items, Any(1))
			} else {
				var (
					pos  int
					last rune
				)
				chain := NewChain()
				for i, c := range v {
					switch c {
					case '*':
						if last == c {
							continue
						} else if i != pos {
							chain.Append(NewString(v[pos:i]))
						}
						chain.Append(Star(0))
						pos = i
					case '?':
						if last != '?' && last != '*' && i != pos {
							chain.Append(NewString(v[pos:i]))
						}
						chain.Append(Any(1))
						pos = i
					default:
						if last == '*' || last == '?' {
							pos = i
						}
					}
					last = c
				}
				if pos < len(v)-1 {
					chain.Append(NewString(v[pos:]))
				}
				if len(chain.Items) == 1 {
					items = append(items, chain.Items[0])
				} else {
					items = append(items, chain)
				}
			}
		} else {
			items = append(items, NewString(vals[i]))
		}
	}
	if len(items) == 0 {
		return nil
	}
	minLen := items[0].MinLen()
	maxLen := items[0].MinLen()
	for i := 1; i < len(items); i++ {
		if minLen > items[i].MinLen() {
			minLen = items[i].MinLen()
		}
		l := items[i].MaxLen()
		if maxLen != -1 && (l == -1 || maxLen < l) {
			maxLen = l
		}
	}
	return &List{Vals: items, MinSize: minLen, MaxSize: maxLen}
}
