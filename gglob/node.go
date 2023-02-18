package gglob

import (
	"strings"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

// NodeItem contains pattern node item (with childs and fixed depth)
type NodeItem struct {
	Node string // raw string (or full glob for terminated)

	Terminated string // end of chain (resulting glob)
	TermIndex  int    // rule num of end of chain (resulting glob), can be used in specific cases

	// size check optimization
	MinSize int
	MaxSize int // 0 or -1 for unlimited

	P      string // prefix or full string if len(inners) == 0
	Suffix string // suffix

	Inners []wildcards.InnerItem // inner segments
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
			return nil, wildcards.ErrNodeEmpty{glob}
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
			pos := wildcards.IndexWildcard(part)
			if pos == -1 {
				lastNode.P = part
			} else {
				if pos > 0 {
					lastNode.P = part[:pos] // prefix
					part = part[pos:]
					lastNode.MinSize = len(lastNode.P)
					lastNode.MaxSize = len(lastNode.P)
				}
				end := wildcards.IndexLastWildcard(part)
				if end == 0 && part[0] != '?' && part[0] != '*' {
					err = wildcards.ErrNodeUnclosed{part}
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
					lastNode.Inners = []wildcards.InnerItem{wildcards.ItemStar{}}
					lastNode.MaxSize = -1 // unlimited
				case "?":
					lastNode.Inners = []wildcards.InnerItem{wildcards.ItemOne{}}
					lastNode.MinSize++
					if lastNode.MaxSize != -1 {
						lastNode.MaxSize++
					}
				default:
					var (
						inner    wildcards.InnerItem
						min, max int
					)
					innerCount := wildcards.WildcardCount(part)
					inners := make([]wildcards.InnerItem, 0, innerCount)
					var (
						prev  wildcards.ItemType
						prevS string
						prevC rune
					)

					for part != "" {
						inner, part, min, max, err = wildcards.NextWildcardItem(part)
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
						typ, s, c := inner.Type()
						switch typ {
						case wildcards.ItemTypeString:
							switch prev {
							case wildcards.ItemTypeString:
								prevS += s
								inners[len(inners)-1] = wildcards.ItemString(prevS)
							case wildcards.ItemTypeRune:
								var sb strings.Builder
								prev = wildcards.ItemTypeString
								sb.Grow(len(s) + 1)
								sb.WriteRune(prevC)
								sb.WriteString(s)
								prevS = sb.String()
								inners[len(inners)-1] = wildcards.ItemString(prevS)
							default:
								if len(inners) == 0 {
									if lastNode.P == "" {
										lastNode.P = s
									} else {
										lastNode.P += s
									}
								} else {
									prev = wildcards.ItemTypeString
									prevS = s
									inners = append(inners, inner)
								}
							}
						case wildcards.ItemTypeRune:
							switch prev {
							case wildcards.ItemTypeString:
								var sb strings.Builder
								sb.Grow(len(prevS) + 1)
								sb.WriteString(prevS)
								sb.WriteRune(c)
								prevS = sb.String()
								inners[len(inners)-1] = wildcards.ItemString(prevS)
							case wildcards.ItemTypeRune:
								var sb strings.Builder
								sb.Grow(2)
								sb.WriteRune(prevC)
								sb.WriteRune(c)
								prev = wildcards.ItemTypeString
								prevS = sb.String()
								inners[len(inners)-1] = wildcards.ItemString(prevS)
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
									prev = wildcards.ItemTypeRune
									prevC = c
									inners = append(inners, inner)
								}
							}
						default:
							prev = wildcards.ItemTypeOther
							inners = append(inners, inner)
						}
					}
					if len(inners) > 1 {
						last := len(inners) - 1
						typ, s, c := inners[last].Type()
						switch typ {
						case wildcards.ItemTypeString:
							if lastNode.Suffix == "" {
								lastNode.Suffix = s
							} else {
								lastNode.Suffix = s + lastNode.Suffix
							}
							inners = inners[:last]
						case wildcards.ItemTypeRune:
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
		err = wildcards.ErrGlobNotExpanded{glob}
		return
	}
	if i != partsCount || (len(lastNode.Childs) == 0 && lastNode.Terminated == "") {
		// child  or/and terminated node
		err = wildcards.ErrNodeNotEnd{lastNode.Node}
	}
	return
}
