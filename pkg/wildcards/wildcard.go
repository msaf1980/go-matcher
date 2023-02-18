package wildcards

import (
	"strings"
)

// WildcardItems contains pattern node item
type WildcardItems struct {
	// size check optimization
	MinSize int
	MaxSize int // 0 or -1 for unlimited

	P      string // prefix or full string if len(inners) == 0
	Suffix string // suffix

	Inners []InnerItem // inner segments
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
	pos := IndexWildcard(glob)
	if pos == -1 {
		node.P = glob
	} else {
		if pos > 0 {
			node.P = glob[:pos] // prefix
			glob = glob[pos:]
			node.MinSize = len(node.P)
			node.MaxSize = len(node.P)
		}
		end := IndexLastWildcard(glob)
		if end == 0 && glob[0] != '?' && glob[0] != '*' {
			err = ErrNodeUnclosed{glob}
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
			node.Inners = []InnerItem{ItemStar{}}
			node.MaxSize = -1 // unlimited
		case "?":
			node.Inners = []InnerItem{ItemOne{}}
			node.MinSize++
			if node.MaxSize != -1 {
				node.MaxSize++
			}
		default:
			var (
				inner    InnerItem
				min, max int
			)
			innerCount := WildcardCount(glob)
			inners := make([]InnerItem, 0, innerCount)
			var (
				prev  ItemType
				prevS string
				prevC rune
			)

			for glob != "" {
				inner, glob, min, max, err = NextWildcardItem(glob)
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
				typ, s, c := inner.Type()
				switch typ {
				case ItemTypeString:
					switch prev {
					case ItemTypeString:
						prevS += s
						inners[len(inners)-1] = ItemString(prevS)
					case ItemTypeRune:
						var sb strings.Builder
						sb.Grow(len(s) + 1)
						sb.WriteRune(prevC)
						sb.WriteString(s)
						prevS = sb.String()
						inners[len(inners)-1] = ItemString(prevS)
					default:
						if len(inners) == 0 {
							if node.P == "" {
								node.P = s
							} else {
								node.P += s
							}
						} else {
							prev = ItemTypeString
							prevS = s
							inners = append(inners, inner)
						}
					}
				case ItemTypeRune:
					switch prev {
					case ItemTypeString:
						var sb strings.Builder
						sb.Grow(len(prevS) + 1)
						sb.WriteString(prevS)
						sb.WriteRune(c)
						prevS = sb.String()
						inners[len(inners)-1] = ItemString(prevS)
					case ItemTypeRune:
						var sb strings.Builder
						sb.Grow(2)
						sb.WriteRune(prevC)
						sb.WriteRune(c)
						prev = ItemTypeString
						prevS = sb.String()
						inners[len(inners)-1] = ItemString(prevS)
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
							prev = ItemTypeRune
							prevC = c
							inners = append(inners, inner)
						}
					}
				default:
					prev = ItemTypeOther
					inners = append(inners, inner)
				}
			}
			if len(inners) > 1 {
				last := len(inners) - 1
				typ, s, c := inners[last].Type()
				switch typ {
				case ItemTypeString:
					if node.Suffix == "" {
						node.Suffix = s
					} else {
						node.Suffix = s + node.Suffix
					}
					inners = inners[:last]
				case ItemTypeRune:
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
