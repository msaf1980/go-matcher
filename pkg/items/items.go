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
	FindList  // list item
	FindStar  // gready item
	FindGroup // group like {a*,b?}
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
		if a[i].Equal(b[i]) {
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

// MatchItems check string against []NodeItems (parsed wildcards or simple regular expression)
func MatchItems(s string, items []Item) bool {
	matched, _ := matchItems(s, items)
	return matched
}

// matchItems check string against []NodeItems (parsed wildcards or simple regular expression)
//
// return
//
// @matched flag for string is matched
//
// @abortGready flag for not matched, but scan is aborted (for example by gready skip scan results)
func matchItems(s string, items []Item) (matched, abortGready bool) {
	var (
		pos, offset int
		flag        FindFlag
	)
	for {
		if pos == len(items) {
			matched = s == ""
			return
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

			for i := 0; i < list.Len(); i++ {
				if offset = list.MatchN(s, i); offset == -1 {
					continue
				}
				s := s[offset:]
				if len(items) == pos {
					if s == "" {
						matched = true
						return
					}
				} else if matched, _ = matchItems(s, items[pos:]); matched {
					return
				}
			}

			if list.IsOptional() {
				if len(items) == 1 {
					if s == "" {
						matched = true
						return
					}
				} else if matched, _ = matchItems(s, items[pos:]); matched {
					return
				}
			}

			return
		case FindNotSupported:
			panic("not supported in match")
		case FindForwarded:
			panic("forwarded match")
		case FindStar:
			pos := pos + 1
			if len(items) <= pos {
				matched = true
				return
			} else {
				return matchStarItems(s, items[pos:])
			}
		default:
			panic(fmt.Errorf("unsupported find flag: %d", flag))
		}
	}
}

// matchStarItems check string against []NodeItems after Star item
//
// return
//
// @matched flag for string is matched
//
// @abortGready flag for not matched, but scan is aborted (for example by gready skip scan results)
func matchStarItems(s string, items []Item) (matched, abortGready bool) {
	var (
		offset, length int
		flag           FindFlag
	)
	if len(items) == 0 {
		matched = s == ""
		abortGready = true
		return
	}
	for {
		if len(items) == 0 || items[0] == nil {
			matched = s == ""
			return
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
			if matched, abortGready = matchItems(sub, items[1:]); matched || abortGready {
				return
			}
		case FindList:
			list := items[0].(ItemList)

			if list.IsOptional() {
				s := s
				for {
					if len(items) == 1 {
						if s == "" {
							matched = true
							return
						}
					} else if matched, _ = matchStarItems(s, items[1:]); matched {
						return
					}
					// skip one symbol and retry scan
					if _, length = utf8.DecodeRuneInString(s); length < 1 {
						break
					}
					s = s[length:]
				}
			}

			// TODO: abort gready with FindFirst
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
						if s == "" {
							matched = true
							return
						}
					} else if matched, _ = matchItems(s, items[1:]); matched {
						return
					}
				}
				if offset == -1 {
					break
				}
				s = s[offset:]
				if _, length = utf8.DecodeRuneInString(s); length < 1 {
					break
				}
				s = s[length:]
			}

			return
		case FindNotSupported:
			panic("not supported in match")
		case FindForwarded:
			panic("forwarded match")
		case FindStar:
			if len(items) == 1 {
				matched = true
				return
			}
			items = items[1:]
			continue
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
