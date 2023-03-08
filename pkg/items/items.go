package items

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type FindFlag int8

const (
	FindNotSupported FindFlag = -1 // can't support greany scan, degrage per symbol scan
)
const (
	FindDone FindFlag = iota // find supported
	// find supported and also can be forwarded to next items (like ?, but for wildcards it's merged)
	FindForwarded
	FindList // list item
	FindStar // gready item
	FindGroup
	FindChain
)

type Item interface {
	Equal(a Item) bool

	// WriteString write string representation
	WriteString(buf *strings.Builder) string

	// MinLen minimum length (bytes)
	MinLen() int
	// MaxLen maximum length (bytes), -1 - if not
	MaxLen() int

	// String is string representation
	String() string

	// WriteRandom is generate random matched string for test
	WriteRandom(buf *strings.Builder)

	// Find is try to locate item and return
	//
	// return
	//
	// @index start index of item or -1 (if not found)
	//
	// @length length of matched
	//
	// @support NotSupported, FindDine, FindForwarded, FindStar
	Find(s string) (index, length int, support FindFlag)

	// Match is try to locate item on start and return
	//
	// return
	//
	// @offset index after item or -1 (if not found)
	//
	// @support FindDone, FindNotSupported, FindList FindStar
	Match(s string) (offset int, support FindFlag)

	// Match is try to locate item on end and return
	//
	// return
	//
	// @offset index after item or -1 (if not found)
	//
	// @support FindDone, FindNotSupported
	MatchLast(s string) (offset int, support FindFlag)
}

func ItemsEqual(a, b []Item) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

type ItemList interface {
	Item

	// IsOptional check when contain empty value and can be skipped
	IsOptional() bool // TODO: refactor with group (minN = 0, max = 1), it's more compatible with regexp

	// Map return map of values
	Map() map[string]struct{}

	// Len is a values count (exclude empty)
	Len() int

	// FindFirst is try find first symbol (exclude empty)
	FindFirst(s string) (index int, supported bool)

	// MatchFirst is try match first ASCII symbol (exclude empty)
	MatchFirst(s string) (ok, supported bool)

	// FindN is try find N value
	//
	// return
	//
	// @offset index of item or -1 (if not found)
	// @length length of founded segment
	FindN(s string, n int) (index, length int)

	// MatchN is try match N value
	//
	// return
	//
	// @offset index after item or -1 (if not found)
	MatchN(s string, n int) (offset int)
}

// NodeItem contains pattern node item (with childs and fixed depth)
type NodeItem struct {
	Node string // raw string (or full glob for terminated)

	Item
}

func NewNodeItem(node string, item Item) NodeItem {
	return NodeItem{
		Node: node,
		Item: item,
	}
}

func (node *NodeItem) WriteString(buf *strings.Builder) {
	buf.WriteString(node.Node)
}

type nextItemsTree struct {
	items []Item
	next  *nextItemsTree
}

func (t *nextItemsTree) IsEmpty() bool {
	return len(t.items) == 0 && (t.next == nil || t.next.IsEmpty())
}

func (t *nextItemsTree) Next() ([]Item, *nextItemsTree) {
	if len(t.items) > 0 {
		return t.items, t.next
	}
	if t.next == nil {
		return nil, nil
	}
	return t.next.items, t.next.next
}

// MatchItems check string against []NodeItems (parsed wildcards or simple regular expression)
func MatchItems(s string, items []Item) bool {
	matched, _ := matchItems(s, items, nil)
	return matched
}

// matchItems check string against []Item (parsed wildcards or simple regular expression)
//
// return
//
// @matched flag for string is matched
//
// @abortGready flag for not matched, but scan is aborted (for example by gready skip scan results)
func matchItems(s string, items []Item, nextItems *nextItemsTree) (matched, abortGready bool) {
	var (
		pos, offset int
		flag        FindFlag
	)
	for {
		if pos == len(items) {
			if nextItems == nil || nextItems.IsEmpty() {
				matched = s == ""
				return
			} else {
				items, nextItems = nextItems.Next()
				matched, _ = matchItems(s, items, nextItems)
				return
			}
		}
		if len(s) < items[pos].MinLen() {
			return
		}
		if offset, flag = items[pos].Match(s); offset == -1 {
			return
		}
		s = s[offset:]
		switch flag {
		case FindDone:
			pos++
		case FindList:
			list := items[pos].(ItemList)
			pos := pos + 1

			if list.IsOptional() {
				s := s
				if len(items) == pos {
					if nextItems == nil || nextItems.IsEmpty() {
						if s == "" {
							matched = true
							return
						}
					} else {
						items, nextItems := nextItems.Next()
						if matched, _ = matchItems(s, items, nextItems); matched {
							return
						}
					}
				} else if matched, _ = matchItems(s, items[pos:], nextItems); matched {
					return
				}
			}

			for i := 0; i < list.Len(); i++ {
				if offset = list.MatchN(s, i); offset == -1 {
					continue
				}
				s := s[offset:]
				if len(items) == pos {
					if nextItems == nil || nextItems.IsEmpty() {
						if s == "" {
							matched = true
							return
						}
					} else {
						items, nextItems := nextItems.Next()
						if matched, _ = matchItems(s, items, nextItems); matched {
							return
						}
					}
				} else if matched, _ = matchItems(s, items[pos:], nextItems); matched {
					return
				}
			}

			return
		case FindGroup:
			group := items[pos].(*Group)
			pos := pos + 1

			if group.IsOptional() {
				s := s
				if len(items) == pos {
					if nextItems == nil || nextItems.IsEmpty() {
						if s == "" {
							matched = true
							return
						}
					} else {
						items, nextItems := nextItems.Next()
						if matched, _ = matchItems(s, items, nextItems); matched {
							return
						}
					}
				} else if matched, _ = matchItems(s, items[pos:], nextItems); matched {
					return
				}
			}

			for i := 0; i < len(group.Vals); i++ {
				next := nextItemsTree{
					items: items[pos:],
					next:  nextItems,
				}
				if v, ok := group.Vals[i].(*Chain); ok {
					if matched, _ = matchItems(s, v.Items, &next); matched {
						return
					}
				} else {
					if matched, _ = matchItems(s, []Item{group.Vals[i]}, &next); matched {
						return
					}
				}
			}
			return
		case FindNotSupported:
			panic("not supported in match")
		case FindForwarded:
			panic("forwarded match")
		case FindStar:
			pos := pos + 1
			if len(items) == pos {
				if nextItems == nil || nextItems.IsEmpty() {
					matched = true
					return
				} else {
					items, nextItems = nextItems.Next()
					return matchStarItems(s, items, nextItems)
				}
			} else {
				return matchStarItems(s, items[pos:], nextItems)
			}
		default:
			panic(fmt.Errorf("unsupported find flag: %d", flag))
		}
	}
}

// matchStarItems check string against []Item after Star item
//
// return
//
// @matched flag for string is matched
//
// @abortGready flag for not matched, but scan is aborted (for example by gready skip scan results)
func matchStarItems(s string, items []Item, nextItems *nextItemsTree) (matched, abortGready bool) {
	var (
		offset, length int
		flag           FindFlag
	)
	if len(items) == 0 {
		if nextItems == nil || nextItems.IsEmpty() {
			matched = s == ""
			abortGready = true
			return
		} else {
			items, nextItems = nextItems.Next()
			return matchStarItems(s, items, nextItems)
		}
	}
	for {
		if len(items) == 0 || items[0] == nil {
			if nextItems == nil || nextItems.IsEmpty() {
				matched = s == ""
				return
			} else {
				items, nextItems = nextItems.Next()
				matched, _ = matchItems(s, items, nextItems)
				return
			}
		}
		if len(s) < items[0].MinLen() {
			abortGready = true
			return
		}

		if offset, length, flag = items[0].Find(s); offset == -1 {
			abortGready = true
			return
		}
		s = s[offset:]

		switch flag {
		case FindDone:
			sub := s[length:]
			if matched, abortGready = matchItems(sub, items[1:], nextItems); matched || abortGready {
				return
			}
		case FindList:
			list := items[0].(ItemList)

			if list.IsOptional() {
				s := s
				for {
					if len(items) == 1 {
						if nextItems == nil || nextItems.IsEmpty() {
							if s == "" {
								matched = true
								break
							}
						} else {
							items, nextItems := nextItems.Next()
							if matched, _ = matchItems(s, items, nextItems); matched {
								return
							}
						}
					} else if matched, _ = matchStarItems(s, items[1:], nextItems); matched {
						return
					}
					// skip one symbol and retry scan
					if _, length = utf8.DecodeRuneInString(s); length < 1 {
						break
					}
					s = s[length:]
				}
			}

			for {
				var offsetN int
				offset := -1
				for i := 0; i < list.Len(); i++ {
					if offsetN, length = list.FindN(s, i); offsetN == -1 {
						// TODO: skip on next scan ?
						continue
					}
					if offset == -1 || offset < offsetN {
						offset = offsetN
					}
					s := s[offsetN+length:]
					if len(items) == 1 {
						if nextItems == nil || nextItems.IsEmpty() {
							if s == "" {
								matched = true
								return
							}
						} else {
							items, nextItems := nextItems.Next()
							if matched, _ = matchItems(s, items, nextItems); matched {
								return
							}
						}
					} else if matched, _ = matchItems(s, items[1:], nextItems); matched {
						return
					}
				}
				if offset == -1 {
					break
				}
				s = s[offset:]
				// shift to one rune from minimal offset
				if _, length = utf8.DecodeRuneInString(s); length < 1 {
					break
				}
				s = s[length:]
			}

			return
		case FindGroup:
			group := items[0].(*Group)

			if group.IsOptional() {
				s := s
				if len(items) == 1 {
					if nextItems == nil || nextItems.IsEmpty() {
						if s == "" {
							matched = true
							return
						}
					} else {
						items, nextItems := nextItems.Next()
						if matched, _ = matchStarItems(s, items, nextItems); matched {
							return
						}
					}
				} else if matched, _ = matchStarItems(s, items[1:], nextItems); matched {
					return
				}
			}

			for i := 0; i < len(group.Vals); i++ {
				next := nextItemsTree{
					items: items[1:],
					next:  nextItems,
				}
				if v, ok := group.Vals[i].(*Chain); ok {
					if matched, _ = matchStarItems(s, v.Items, &next); matched {
						return
					}
				} else {
					if matched, _ = matchStarItems(s, []Item{group.Vals[i]}, &next); matched {
						return
					}
				}
			}

			return
		case FindNotSupported:
			panic("not supported in match")
		case FindForwarded:
			panic("forwarded match")
		case FindStar:
			if len(items) == 1 {
				if nextItems == nil || nextItems.IsEmpty() {
					matched = true
					return
				} else {
					items, nextItems := nextItems.Next()
					if matched, _ = matchStarItems(s, items, nextItems); matched {
						return
					}
				}
			}
			items = items[1:]
			// continue
		default:
			panic(fmt.Errorf("unsupported find flag: %d", flag))
		}

		if _, length = utf8.DecodeRuneInString(s); length < 1 {
			break
		}
		s = s[length:]
	}
	return
}

func AppendItem(items []Item, item Item) []Item {
	if len(items) == 0 {
		return append(items, item)
	}
	last := len(items) - 1
	switch v := item.(type) {
	case *String:
		switch vv := items[last].(type) {
		case *String:
			vv.Add(v.S)
		case Rune:
			v.PrependRune(rune(vv))
			items[last] = v
		case Byte:
			v.PrependByte(byte(vv))
			items[last] = v
		default:
			items = append(items, item)
		}
	case Rune:
		c := rune(v)
		switch vv := items[last].(type) {
		case *String:
			vv.AddRune(c)
		case Rune:
			items[last] = vv.AppendRune(c)
		case Byte:
			items[last] = vv.AppendRune(c)
		default:
			items = append(items, item)
		}
	case Byte:
		c := byte(v)
		switch vv := items[last].(type) {
		case *String:
			vv.AddByte(c)
		case Rune:
			items[last] = vv.AppendByte(c)
		case Byte:
			items[last] = vv.AppendByte(c)
		default:
			items = append(items, item)
		}

	case Any:
		switch vv := items[last].(type) {
		case Any:
			vv += v
			items[last] = vv
		case Star:
			vv += Star(v)
			items[last] = vv
		default:
			items = append(items, item)
		}
	case Star:
		switch vv := items[last].(type) {
		case Any:
			v += Star(vv)
			items[last] = v
		case Star:
			vv += v
			items[last] = vv
		default:
			items = append(items, item)
		}
	default:
		items = append(items, item)
	}
	return items
}
