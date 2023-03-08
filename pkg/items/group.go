package items

import (
	"math/rand"
	"strings"
)

// Group is list of items
type Group struct {
	Vals    []Item
	MinSize int
	MaxSize int
}

func (item *Group) Equal(a Item) bool {
	if v, ok := a.(*Group); ok {
		if len(item.Vals) == len(v.Vals) {
			for i := 0; i < len(item.Vals); i++ {
				if !item.Vals[i].Equal(v.Vals[i]) {
					return false
				}
			}
			return true
		}

	}
	return false
}

func (item *Group) WriteRandom(buf *strings.Builder) {
	item.Vals[rand.Intn(len(item.Vals))].WriteRandom(buf)
}

func (item *Group) WriteString(buf *strings.Builder) string {
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

func (item *Group) String() string {
	var buf strings.Builder
	return item.WriteString(&buf)
}

func (item *Group) MinLen() int {
	return item.MinSize
}

func (item *Group) MaxLen() int {
	return item.MaxSize
}

func (item *Group) Find(s string) (index, length int, support FindFlag) {
	support = FindGroup
	return
}

func (item *Group) Match(s string) (offset int, support FindFlag) {
	support = FindGroup
	return
}

// FindFirst is try find first symbol (exclude empty)
func (item *Group) FindFirst(s string) (index int, supported bool) {
	return
}

func (item *Group) MatchLast(s string) (offset int, support FindFlag) {
	support = FindNotSupported
	return
}

// IsOptional check when contain empty value and can be skipped
func (item *Group) IsOptional() bool {
	return item.MinSize == 0
}

func NewGroup(vals []string) (item Item, err error) {
	items := make([]Item, 0, len(vals))
	for i := 0; i < len(vals); i++ {
		if HasWildcard(vals[i]) {
			v := vals[i]

			if v == "*" {
				return Star(0), nil
			} else if v == "?" {
				items = append(items, Any(1))
			} else {
				chain := NewChain()
				for v != "" {
					item, v, err = NextWildcardItem(v)
					if err != nil {
						return
					}
					if item == nil {
						continue
					}
					chain.Append(item)
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
		return nil, nil
	}
	minLen := items[0].MinLen()
	maxLen := items[0].MaxLen()
	for i := 1; i < len(items); i++ {
		if minLen > items[i].MinLen() {
			minLen = items[i].MinLen()
		}
		l := items[i].MaxLen()
		if maxLen != -1 && (l == -1 || maxLen < l) {
			maxLen = l
		}
	}
	return &Group{Vals: items, MinSize: minLen, MaxSize: maxLen}, nil
}
