package items

import (
	"strings"
)

type ItemType int8

const (
	ItemTypeOther ItemType = iota
	ItemTypeString
	ItemTypeRune
)

type Item interface {
	Type() (typ ItemType, s string, c rune) // return type, and string or rune value (if contain)
	Strings() []string                      // return nil or string values (if contain)
	Match(part string, nextItems []Item) (found bool)
	Locate(part string, nextItems []Item) (offset int, support bool, skipItems int)
	WriteString(buf *strings.Builder)
}

func WriteInnerItems(items []Item, buf *strings.Builder) {
	for _, item := range items {
		item.WriteString(buf)
	}
}

// NodeItem contains pattern node item
type NodeItem struct {
	// size check optimization
	MinSize int
	MaxSize int // 0 or -1 for unlimited

	P      string // prefix or full string if len(inners) == 0
	Suffix string // suffix

	Inners []Item // inner segments
}

func (node *NodeItem) Match(part string) (matched bool) {
	if len(part) < node.MinSize {
		return
	}
	if node.MaxSize > 0 {
		if len(part) > node.MaxSize {
			return
		}
	}
	if len(node.Inners) == 0 {
		if node.Suffix != "" {
			return
		}
		matched = (node.P == part)
	} else {
		if node.P != "" {
			if !strings.HasPrefix(part, node.P) {
				// prefix not match
				return
			}
			part = part[len(node.P):]
		}
		if node.Suffix != "" {
			if !strings.HasSuffix(part, node.Suffix) {
				// suffix not match
				return
			}
			part = part[:len(part)-len(node.Suffix)]
		}

		matched = node.Inners[0].Match(part, node.Inners[1:])
	}
	return
}
