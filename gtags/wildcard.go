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

		matched = node.Inners[0].Match(part, "", node.Inners[1:])
	}
	return
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
			var (
				prev  items.ItemType
				prevS string
				prevC rune
			)

			for glob != "" {
				inner, glob, min, max, err = items.NextWildcardItem(glob)
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
				// try to in-palce merge
				if s, ok := inner.IsString(); ok {
					switch prev {
					case items.ItemTypeString:
						prevS += s
						inners[len(inners)-1] = items.ItemString(prevS)
					case items.ItemTypeChar:
						prevS = string(prevC) + s
						inners[len(inners)-1] = items.ItemString(prevS)
					default:
						if len(inners) == 0 {
							if node.P == "" {
								node.P = s
							} else {
								node.P += s
							}
						} else {
							prev = items.ItemTypeString
							prevS = s
							inners = append(inners, inner)
						}
					}
				} else if c, ok := inner.IsRune(); ok {
					switch prev {
					case items.ItemTypeString:
						var sb strings.Builder
						sb.Grow(len(prevS) + 1)
						sb.WriteString(prevS)
						sb.WriteRune(c)
						prev = items.ItemTypeString
						prevS = sb.String()
						inners[len(inners)-1] = items.ItemString(prevS)
					case items.ItemTypeChar:
						var sb strings.Builder
						sb.Grow(2)
						sb.WriteRune(prevC)
						sb.WriteRune(c)
						prev = items.ItemTypeString
						prevS = sb.String()
						inners[len(inners)-1] = items.ItemString(prevS)
					default:
						if len(inners) == 0 {
							if node.P == "" {
								node.P = string(c)
							} else {
								var sb strings.Builder
								sb.Grow(len(node.P) + 1)
								sb.WriteString(node.P)
								sb.WriteRune(c)
								node.P = sb.String()
							}
						} else {
							prev = items.ItemTypeChar
							prevC = c
							inners = append(inners, inner)
						}
					}
				} else {
					inners = append(inners, inner)
				}
			}
			if len(inners) > 1 {
				last := len(inners) - 1
				// var size int
				if s, ok := inners[last].IsString(); ok {
					if node.Suffix == "" {
						node.Suffix = s
					} else {
						node.Suffix = s + node.Suffix
					}
					inners = inners[:last]
				} else if c, ok := inners[last].IsRune(); ok {
					var sb strings.Builder
					sb.Grow(len(node.Suffix) + 1)
					sb.WriteRune(c)
					sb.WriteString(node.Suffix)
					node.Suffix = sb.String()
					inners = inners[:last]
				}
			}
			if len(inners) == 0 {
				if node.Suffix != "" {
					if node.P == "" {
						node.P = node.Suffix
					} else {
						node.P += node.Suffix
					}
					node.Suffix = ""
				}
			} else {
				node.Inners = inners
			}
		}
	}
	return
}
