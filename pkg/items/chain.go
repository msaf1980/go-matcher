package items

import "strings"

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
	if index, length, support = c.Items[0].Find(s); index == -1 {
		return
	}
	support = FindChain

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
