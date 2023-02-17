package items

import (
	"io"
	"strings"
	"unicode/utf8"
)

type InnerItem interface {
	CanString() bool
	IsRune() (c rune, ok bool)
	IsString() (s string, ok bool)
	Match(part string, nextParts string, nextItems []InnerItem) (found bool)
}

type ItemOne struct{}

func (item ItemOne) IsRune() (rune, bool) {
	return utf8.RuneError, false
}

func (item ItemOne) IsString() (string, bool) {
	return "", false
}

func (ItemOne) CanString() bool {
	return false
}

func (item ItemOne) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		found = true
		part = part[n:]
	}
	if found {
		if part != "" && len(nextItems) > 0 {
			found = nextItems[0].Match(part, nextParts, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}

// NodeItem contains pattern node item
type NodeItem struct {
	Node string // raw string (or full glob for terminated)

	Terminated string // end of chain (resulting glob)
	TermIndex  int    // rule num of end of chain (resulting glob), can be used in specific cases

	// size check optimization
	MinSize int
	MaxSize int // 0 or -1 for unlimited

	P      string // prefix or full string if len(inners) == 0
	Suffix string // suffix

	Inners []InnerItem // inner segments
	// TODO: may be some ordered tree for complete string nodes search speedup (on large set) ?
	Childs []*NodeItem // next possible parts slice
}

func (node *NodeItem) MatchNode(part string) (matched bool) {
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

func (node *NodeItem) MatchItems(part string, nextParts string, matched *[]string) {
	if node.MatchNode(part) {
		if node.Terminated != "" {
			*matched = append(*matched, node.Terminated)
		} else if nextParts != "" {
			part, nextParts, _ = strings.Cut(nextParts, ".")
			for _, child := range node.Childs {
				child.MatchItems(part, nextParts, matched)
			}
		}
	}
}

func (node *NodeItem) MatchIndexedItems(part string, nextParts string, matched *[]int) {
	if node.MatchNode(part) {
		if node.Terminated != "" {
			*matched = append(*matched, node.TermIndex)
		} else if nextParts != "" {
			part, nextParts, _ = strings.Cut(nextParts, ".")
			for _, child := range node.Childs {
				child.MatchIndexedItems(part, nextParts, matched)
			}
		}
	}
}

func (node *NodeItem) MatchFirstItems(part string, nextParts string, minMatched *int) {
	if node.MatchNode(part) {
		if node.Terminated != "" {
			if *minMatched == -1 || *minMatched > node.TermIndex {
				*minMatched = node.TermIndex
			}
		} else if nextParts != "" {
			part, nextParts, _ = strings.Cut(nextParts, ".")
			for _, child := range node.Childs {
				child.MatchFirstItems(part, nextParts, minMatched)
			}
		}
	}
}

func (node *NodeItem) MatchIndexedItemsPart(part string, parts []string, matched *[]int) {
	if node.MatchNode(part) {
		if node.Terminated != "" {
			*matched = append(*matched, node.TermIndex)
		} else if len(parts) > 0 {
			for _, child := range node.Childs {
				child.MatchIndexedItemsPart(parts[0], parts[1:], matched)
			}
		}
	}
}

func (node *NodeItem) MatchFirstItemsPart(part string, parts []string, minMatched *int) {
	if node.MatchNode(part) {
		if node.Terminated != "" {
			if *minMatched == -1 || *minMatched > node.TermIndex {
				*minMatched = node.TermIndex
			}
		} else if len(parts) > 0 {
			for _, child := range node.Childs {
				child.MatchFirstItemsPart(parts[0], parts[1:], minMatched)
			}
		}
	}
}

func (node *NodeItem) MatchItemsPart(part string, parts []string, matched *[]string) {
	if node.MatchNode(part) {
		if node.Terminated != "" {
			*matched = append(*matched, node.Terminated)
		} else if len(parts) > 0 {
			for _, child := range node.Childs {
				child.MatchItemsPart(parts[0], parts[1:], matched)
			}
		}
	}
}

// Match match root node item (recursieve with childs) for graphite path (dot-delimited, like a.b.c)
func (node *NodeItem) Match(path string, matched *[]string) {
	for _, node := range node.Childs {
		part, nextParts, _ := strings.Cut(path, ".")
		// match first node
		node.MatchItems(part, nextParts, matched)
	}
}

func (node *NodeItem) MatchIndexed(path string, matched *[]int) {
	for _, node := range node.Childs {
		part, nextParts, _ := strings.Cut(path, ".")
		// match first node
		node.MatchIndexedItems(part, nextParts, matched)
	}
}

func (node *NodeItem) MatchFirst(path string, minMatched *int) {
	for _, node := range node.Childs {
		part, nextParts, _ := strings.Cut(path, ".")
		// match first node
		node.MatchFirstItems(part, nextParts, minMatched)
	}
}

// Match match root node item (recursieve with childs) for splitted path parts
func (node *NodeItem) MatchByParts(parts []string, matched *[]string) {
	for _, node := range node.Childs {
		// match first node
		node.MatchItemsPart(parts[0], parts[1:], matched)
	}
}

func (node *NodeItem) MatchIndexedByParts(parts []string, matched *[]int) {
	for _, node := range node.Childs {
		// match first node
		node.MatchIndexedItemsPart(parts[0], parts[1:], matched)
	}
}

func (node *NodeItem) MatchFirstByParts(parts []string, minMatched *int) {
	for _, node := range node.Childs {
		// match first node
		node.MatchFirstItemsPart(parts[0], parts[1:], minMatched)
	}
}

type ItemType int8

const (
	ItemTypeOther ItemType = iota
	ItemTypeString
	ItemTypeChar
)

// Parse add glob for graphite path (dot-delimited, like a.b*.[a-c].{wait,idle}
func (node *NodeItem) Parse(glob string, partsCount int, termIdx int) (lastNode *NodeItem, err error) {
	var (
		i    int
		part string
	)

	last := partsCount - 1
	nextParts := glob
	for nextParts != "" {
		found := false
		part, nextParts, _ = strings.Cut(nextParts, ".")
		if part == "" {
			return nil, ErrNodeEmpty{glob}
		}
		for _, child := range node.Childs {
			if part == child.Node {
				node = child
				found = true
				break
			}
		}
		if !found {
			if i == last {
				// last node, so terminate match
				lastNode = &NodeItem{Node: part, Terminated: glob, TermIndex: termIdx}
			} else {
				lastNode = &NodeItem{Node: part}
			}
			pos := IndexWildcard(part)
			if pos == -1 {
				lastNode.P = part
			} else {
				if pos > 0 {
					lastNode.P = part[:pos] // prefix
					part = part[pos:]
					lastNode.MinSize = len(lastNode.P)
					lastNode.MaxSize = len(lastNode.P)
				}
				end := IndexLastWildcard(part)
				if end == 0 && part[0] != '?' && part[0] != '*' {
					err = ErrNodeUnclosed{part}
					return
				}
				if end < len(part)-1 {
					end++
					lastNode.Suffix = part[end:]
					part = part[:end]
					lastNode.MinSize += len(lastNode.Suffix)
					lastNode.MaxSize += len(lastNode.Suffix)
				}

				switch part {
				case "*":
					lastNode.Inners = []InnerItem{ItemStar{}}
					lastNode.MaxSize = -1 // unlimited
				case "?":
					lastNode.Inners = []InnerItem{ItemOne{}}
					lastNode.MinSize++
					if lastNode.MaxSize != -1 {
						lastNode.MaxSize++
					}
				default:
					var (
						inner    InnerItem
						min, max int
					)
					innerCount := WildcardCount(part)
					inners := make([]InnerItem, 0, innerCount)
					var (
						prev  ItemType
						prevS string
						prevC rune
					)

					for part != "" {
						inner, part, min, max, err = NextWildcardItem(part)
						if err != nil {
							return
						}
						if inner == nil {
							continue
						}
						lastNode.MinSize += min
						if lastNode.MaxSize != -1 {
							if max == -1 {
								lastNode.MaxSize = -1
							} else {
								lastNode.MaxSize += max
							}
						}
						// try to in-palce merge
						if s, ok := inner.IsString(); ok {
							switch prev {
							case ItemTypeString:
								prevS += s
								inners[len(inners)-1] = ItemString(prevS)
							case ItemTypeChar:
								prevS = string(prevC) + s
								inners[len(inners)-1] = ItemString(prevS)
							default:
								if len(inners) == 0 {
									if lastNode.P == "" {
										lastNode.P = s
									} else {
										lastNode.P += s
									}
								} else {
									prev = ItemTypeString
									prevS = s
									inners = append(inners, inner)
								}
							}
						} else if c, ok := inner.IsRune(); ok {
							switch prev {
							case ItemTypeString:
								var sb strings.Builder
								sb.Grow(len(prevS) + 1)
								sb.WriteString(prevS)
								sb.WriteRune(c)
								prev = ItemTypeString
								prevS = sb.String()
								inners[len(inners)-1] = ItemString(prevS)
							case ItemTypeChar:
								var sb strings.Builder
								sb.Grow(2)
								sb.WriteRune(prevC)
								sb.WriteRune(c)
								prev = ItemTypeString
								prevS = sb.String()
								inners[len(inners)-1] = ItemString(prevS)
							default:
								if len(inners) == 0 {
									if lastNode.P == "" {
										lastNode.P = string(c)
									} else {
										var sb strings.Builder
										sb.Grow(len(lastNode.P) + 1)
										sb.WriteString(lastNode.P)
										sb.WriteRune(c)
										lastNode.P = sb.String()
									}
								} else {
									prev = ItemTypeChar
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
							if lastNode.Suffix == "" {
								lastNode.Suffix = s
							} else {
								lastNode.Suffix = s + lastNode.Suffix
							}
							inners = inners[:last]
						} else if c, ok := inners[last].IsRune(); ok {
							var sb strings.Builder
							sb.Grow(len(lastNode.Suffix) + 1)
							sb.WriteRune(c)
							sb.WriteString(lastNode.Suffix)
							lastNode.Suffix = sb.String()
							inners = inners[:last]
						}
					}
					if len(inners) == 0 {
						if lastNode.Suffix != "" {
							if lastNode.P == "" {
								lastNode.P = lastNode.Suffix
							} else {
								lastNode.P += lastNode.Suffix
							}
							lastNode.Suffix = ""
						}
					} else {
						lastNode.Inners = inners
					}
				}
			}
			node.Childs = append(node.Childs, lastNode)
			node = lastNode
		}
		i++
	}
	if lastNode == nil {
		err = ErrGlobNotExpanded{glob}
		return
	}
	if i != partsCount || (len(lastNode.Childs) == 0 && lastNode.Terminated == "") {
		// child  or/and terminated node
		err = ErrNodeNotEnd{lastNode.Node}
	}
	return
}

// NextWildcardItem extract InnerItem from glob (not regexp)
func NextWildcardItem(s string) (item InnerItem, next string, minLen int, maxLen int, err error) {
	if s == "" {
		return nil, s, 0, 0, io.EOF
	}
	switch s[0] {
	case '[':
		if idx := strings.Index(s, "]"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		runes, failed := RunesExpand(s)
		if failed {
			return nil, s, 0, 0, ErrNodeMissmatch{"rune", s}
		}
		if len(runes) == 0 {
			return nil, next, 0, 0, nil
		}
		if len(runes) == 1 {
			if runes[0].First == runes[0].Last {
				// one item optimization
				return ItemRune(runes[0].First), next, 1, 1, nil
			}
		}
		return ItemRuneRanges(runes), next, 1, 1, nil
	case '{':
		if idx := strings.Index(s, "}"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		vals, failed := ListExpand(s)
		if failed {
			return nil, s, 0, 0, ErrNodeMissmatch{"list", s}
		}
		item, minLen, maxLen = NewItemList(vals)
		return
	case '*':
		var next string
		for i, c := range s {
			if c != '*' {
				next = s[i:]
				break
			}
		}
		return ItemStar{}, next, 0, -1, nil
	case '?':
		next := s[1:]
		return ItemOne{}, next, 1, 1, nil
	case ']', '}':
		return nil, s, 0, 0, ErrNodeUnclosed{s}
	default:
		// string segment
		end := IndexWildcard(s)
		v, next := SplitString(s, end)
		if len(v) == 0 {
			return nil, next, len(v), len(v), nil
		}
		if len(v) == 1 {
			return ItemRune(v[0]), next, len(v), len(v), nil
		}
		return ItemString(v), next, len(v), len(v), nil
	}
}
