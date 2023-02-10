package gtags

import (
	"strings"

	"github.com/msaf1980/go-matcher/pkg/items"
)

// WildcardItems contains pattern node item
type WildcardItems struct {
	// size check optimization
	MinSize int
	MaxSize int // 0 or -1 for unlimited

	P      string // prefix or full string if len(inners) == 0
	Suffix string // suffix

	Inners []items.InnerItem // inner segments
}

func (node *WildcardItems) Match(part string) (matched bool) {
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

		matched, _ = node.Inners[0].Match(part, "", node.Inners[1:], false)
	}
	return
}

// Merge is trying to merge inners
func (node *WildcardItems) Merge(inners []items.InnerItem) {
	if len(inners) == 0 {
		if node.P != "" && node.Suffix != "" {
			node.P += node.Suffix
			node.Suffix = ""
		} else if node.Suffix != "" {
			node.P = node.Suffix
			node.Suffix = ""
		}
		return
	}
	if len(inners) == 1 {
		// merge
		switch v := inners[0].(type) {
		case items.ItemString:
			if node.P != "" || node.Suffix != "" {
				node.P = node.P + string(v) + node.Suffix
				node.Suffix = ""
			} else {
				node.P = string(v)
			}
		default:
			node.Inners = inners
		}
		return
	} else {
		switch v := inners[0].(type) {
		case items.ItemString:
			var sb strings.Builder
			// merge strings from prefix
			s := string(v)
			sb.Grow(len(node.P) + len(node.Suffix) + len(s))
			sb.WriteString(node.P)
			sb.WriteString(s)
			i := 1
			for i < len(inners) {
				if s, ok := inners[i].IsString(); ok {
					sb.WriteString(s)
					i++
				} else {
					break
				}
			}

			if i == len(inners) {
				// merge all strings from start to last string
				sb.WriteString(node.Suffix)
				node.P = sb.String()
				node.Suffix = ""
				return
			} else {
				node.P = sb.String()
				inners = inners[i:]
			}
		}

		last := len(inners) - 1
		switch v := inners[last].(type) {
		case items.ItemString:
			//merge to suffix
			size := len(node.Suffix) + len(v)
			i := last - 1
			for i > 0 {
				if s, ok := inners[i].IsString(); ok {
					size += len(s)
					i--
				} else {
					break
				}
			}
			i++
			last = i
			var sb strings.Builder
			sb.Grow(size)
			for ; i < len(inners); i++ {
				s, _ := inners[i].IsString()
				sb.WriteString(s)
			}
			sb.WriteString(node.Suffix)
			node.Suffix = sb.String()
			inners = inners[:last]
		}

		node.Inners = inners
	}
}

func (node *WildcardItems) Parse(glob string) (err error) {
	pos := items.IndexWildcard(glob)
	if pos == -1 {
		node.P = glob
	} else {
		if pos > 0 {
			node.P = glob[:pos] // prefix
			glob = glob[pos:]
			node.MinSize = len(node.P)
			node.MaxSize = len(node.P)
		}
		end := items.IndexLastWildcard(glob)
		if end == 0 && glob[0] != '?' && glob[0] != '*' {
			err = items.ErrNodeUnclosed{glob}
			return
		}
		if end < len(glob)-1 {
			end++
			node.Suffix = glob[end:]
			glob = glob[:end]
			node.MinSize += len(node.Suffix)
			node.MaxSize += len(node.Suffix)
		}

		switch glob {
		case "*":
			node.Inners = []items.InnerItem{items.ItemStar{}}
			node.MaxSize = -1 // unlimited
		case "?":
			node.Inners = []items.InnerItem{items.ItemOne{}}
			node.MinSize++
			if node.MaxSize != -1 {
				node.MaxSize++
			}
		default:
			var (
				inner    items.InnerItem
				min, max int
			)
			innerCount := items.WildcardCount(glob)
			inners := make([]items.InnerItem, 0, innerCount)

			for glob != "" {
				inner, glob, min, max, err = items.NextInnerItem(glob)
				if err != nil {
					return
				}
				if inner == nil {
					continue
				}
				node.MinSize += min
				if node.MaxSize != -1 {
					if max == -1 {
						node.MaxSize = -1
					} else {
						node.MaxSize += max
					}
				}
				inners = append(inners, inner)
			}
			node.Merge(inners)
		}
	}
	return
}
