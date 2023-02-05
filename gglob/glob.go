package gglob

import (
	"io"
	"strings"
	"unicode/utf8"
)

type NodeType int8

const (
	NodeEmpty NodeType = iota
	NodeRoot           // root node (initial)
	NodeString
	NodeList   // {a,bc}
	NodeRange  // [a-c]
	NodeOne    // ?
	NodeStar   // *
	NodeInners // composite type, contains prefix, suffix, subitems in []inners
)

// TODO: nodeListSkip scan

type InnerItem struct {
	// can be nodeString, nodeList, nodeRange, nodeOne, nodeMany
	Typ NodeType

	// string (nodeString)
	P string

	// nodeList or 	nodeRange
	Vals    []string // strings for nodeList
	ValsMin int      // min len in vals or min rune in range
	ValsMax int      // max len in vals or max rune in range
}

func (item *InnerItem) matchStar(part string, nextParts string, nextItems []*InnerItem) (found bool) {
	if part == "" && len(nextItems) == 0 {
		return true
	}

	for ; part != ""; part = part[1:] {
		part := part           // avoid overwrite outer loop
		nextItems := nextItems // avoid overwrite outer loop
		if len(nextItems) > 0 {
			nextItem := nextItems[0]
			switch nextItem.Typ {
			// speedup NodeString find
			case NodeString:
				if idx := strings.Index(part, nextItem.P); idx == -1 {
					// string not found, no need star scan
					break
				} else {
					idx += len(nextItem.P)
					part = part[idx:]
					nextItems = nextItems[1:]
					found = true
				}
			}
		} else {
			// all of string matched to *
			part = ""
			found = true
		}
		if found {
			if part != "" && len(nextItems) > 0 {
				found = nextItems[0].matchItem(part, nextParts, nextItems[1:])
			} else if part != "" || len(nextItems) > 0 {
				found = false
			}
			if found {
				break
			}
		}
	}
	return
}

func (item *InnerItem) matchItem(part string, nextParts string, nextItems []*InnerItem) (found bool) {
	switch item.Typ {
	case NodeStar:
		return item.matchStar(part, nextParts, nextItems)
	case NodeString:
		if part == item.P {
			// full match
			found = true
			part = ""
		} else if strings.HasPrefix(part, item.P) {
			// strip prefix
			found = true
			part = part[len(item.P):]
		}
	case NodeOne:
		if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
			if len(part) >= n {
				found = true
				part = part[n:]
			}
		}
	}
	if found {
		if part != "" && len(nextItems) > 0 {
			found = nextItems[0].matchItem(part, nextParts, nextItems[1:])
		} else if part != "" || len(nextItems) > 0 {
			found = false
		}
	}
	return
}

func nextInnerItem(s string) (*InnerItem, string, int, int, error) {
	if s == "" {
		return nil, s, 0, 0, io.EOF
	}
	switch s[0] {
	case '[':
		// TODO: implement nodeRange
		return nil, s, 0, 0, ErrNodeUnclosed{s}
	case '{':
		// TODO: implement nodeList
		return nil, s, 0, 0, ErrNodeUnclosed{s}
	case '*':
		next := nextString(s, 1)
		return &InnerItem{
			Typ: NodeStar,
		}, next, 0, 0, nil
	case '?':
		next := nextString(s, 1)
		return &InnerItem{
			Typ: NodeOne,
		}, next, 1, 1, nil
	case ']', '}':
		return nil, s, 0, 0, ErrNodeUnclosed{s}
	default:
		// string segment
		end := IndexWildcard(s)
		v, next := splitString(s, end)
		return &InnerItem{
			Typ: NodeString, P: v,
		}, next, len(v), len(v), nil
	}
}

// NodeItem contains pattern node item
type NodeItem struct {
	Node string // raw string (or full glob for terminated)

	Terminated bool // end of chain

	// size check optimization
	MinSize int
	MaxSize int // 0 or -1 for unlimited

	InnerItem // if one item, no need to use []inners

	// suffix for NodeInners
	Suffix string

	Inners []*InnerItem // inner segments
	Childs []*NodeItem  // next possible parts tree
}

func (node *NodeItem) Match(part string, nextParts string, items *[]string) {
	var found bool

	if node.Typ != NodeString {
		if len(part) < node.MinSize {
			return
		}
		if node.MaxSize > 0 {
			if len(part) > node.MaxSize {
				return
			}
		}
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
	}

	if node.Typ == NodeInners {
		if len(node.Inners) == 0 {
			// some broken, skip node
			return
		}
		found = node.Inners[0].matchItem(part, nextParts, node.Inners[1:])
	} else {
		found = node.matchItem(part, nextParts, nil)
	}

	if found {
		if node.Terminated {
			*items = append(*items, node.Node)
		} else if len(nextParts) > 0 {
			for _, child := range node.Childs {
				part, nextParts, _ = strings.Cut(nextParts, ".")
				child.Match(part, nextParts, items)
			}
		}
	}
}

// GlobMatcher is dotted-separated segment glob matcher, like a.b.[c-e]?.{f-o}*, writted for graphite project
type GlobMatcher struct {
	Root  map[int]*NodeItem
	Globs map[string]bool
}

func NewGlobMatcher() *GlobMatcher {
	return &GlobMatcher{
		Root:  make(map[int]*NodeItem),
		Globs: make(map[string]bool),
	}
}

func (w *GlobMatcher) Adds(globs []string) (err error) {
	for _, glob := range globs {
		if err = w.Add(glob); err != nil {
			return err
		}
	}
	return
}

func (w *GlobMatcher) Add(glob string) (err error) {
	if glob == "" {
		return
	}
	if w.Globs[glob] {
		// aleady added
		return
	}
	parts := pathSplit(glob)
	if hasEmptyParts(parts) {
		return ErrNodeEmpty{glob}
	}

	node, ok := w.Root[len(parts)]
	if !ok {
		node = &NodeItem{InnerItem: InnerItem{Typ: NodeRoot}}
		w.Root[len(parts)] = node
	}
	var newNode *NodeItem

	last := len(parts) - 1
	for i, part := range parts {
		found := false
		for _, child := range node.Childs {
			if part == node.Node {
				node = child
				found = true
				break
			}
		}
		if !found {
			if i == last {
				// last node, so terminate match
				newNode = &NodeItem{Node: glob, Terminated: true}
			} else {
				newNode = &NodeItem{Node: part}
			}
			pos := IndexWildcard(part)
			if pos == -1 {
				newNode.Typ = NodeString
				newNode.P = part
			} else {
				if pos > 0 {
					newNode.P = part[:pos] // prefix
					part = part[pos:]
					newNode.MinSize = len(newNode.P)
					newNode.MaxSize = len(newNode.P)
				}
				end := IndexLastWildcard(part)
				if end == 0 && part[0] != '?' && part[0] != '*' {
					return ErrNodeUnclosed{part}
				}
				if end < len(part)-1 {
					end++
					newNode.Suffix = part[end:]
					part = part[:end]
					newNode.MinSize += len(newNode.Suffix)
					newNode.MaxSize += len(newNode.Suffix)
				}

				switch part {
				case "*":
					newNode.Typ = NodeStar
					newNode.MaxSize = -1 // unlimited
				case "?":
					newNode.Typ = NodeOne
					newNode.MinSize++
					if newNode.MaxSize != -1 {
						newNode.MaxSize++
					}
				default:
					var inner *InnerItem
					innerCount := WildcardCount(part)
					inners := make([]*InnerItem, 0, innerCount)

					var (
						min, max int
					)
					for part != "" {
						inner, part, min, max, err = nextInnerItem(part)
						if err != nil {
							return
						}
						newNode.MinSize += min
						if newNode.MaxSize != -1 {
							if inner.Typ == NodeStar {
								newNode.MaxSize = -1
							} else {
								newNode.MaxSize += max
							}
						}
						inners = append(inners, inner)
					}
					if len(inners) == 0 {
						// no inners for inner node
						return ErrGlobNotExpanded{newNode.Node}
					} else if len(inners) == 1 {
						if inners[0].Typ == NodeString || inners[0].P != "" {
							// bug
							panic(inners[0])
						}
						// merge
						newNode.InnerItem = *inners[0]
					} else {
						newNode.Typ = NodeInners
						newNode.Inners = inners
					}
				}
			}
			node.Childs = append(node.Childs, newNode)
			node = newNode
		}
	}

	if newNode != nil {
		if len(newNode.Childs) > 0 || !newNode.Terminated {
			// childs for terminated node
			return ErrNodeNotEnd{newNode.Node}
		}
		w.Globs[glob] = true
	}

	return
}

func (w *GlobMatcher) Match(path string) (globs []string) {
	if path == "" {
		return nil
	}
	partsCount := pathLevel(path)
	var items []string
	if node, ok := w.Root[partsCount]; ok {
		items = make([]string, 0, min(4, len(node.Childs)))
		for _, node := range node.Childs {
			part, nextParts, _ := strings.Cut(path, ".")
			// match first node
			node.Match(part, nextParts, &items)
		}
	}

	return items
}
