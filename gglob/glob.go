package gglob

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"
)

type NodeType int8

const (
	NodeEmpty NodeType = iota
	NodeRoot           // root node (initial)
	NodeString
	NodeList   // {a,bc}
	NodeRune   // [a-c]
	NodeOne    // ?
	NodeStar   // *
	NodeInners // composite type, contains prefix, suffix, subitems in []inners
)

var (
	nodeTypeStrings = []string{"", "root", "string", "list", "rune", "?", "*", "inners"}
)

func (n NodeType) String() string {
	return nodeTypeStrings[n]
}

// TODO: nodeListSkip scan

type InnerItem struct {
	// can be nodeString, nodeList, nodeRange, nodeOne, nodeMany
	Typ NodeType

	// string (nodeString)
	P string

	// nodeList
	Vals    []string // strings
	ValsMin int      // min len in vals or min rune in range
	ValsMax int      // max len in vals or max rune in range

	// 	nodeRange
	Runes map[rune]struct{}
}

func (item *InnerItem) matchStar(part string, nextParts string, nextItems []*InnerItem) (found bool) {
	if part == "" && len(nextItems) == 0 {
		return true
	}

	nextOffset := 1 // string skip optimization
LOOP:
	for ; part != ""; part = part[nextOffset:] {
		part := part           // avoid overwrite outer loop
		nextItems := nextItems // avoid overwrite outer loop
		nextOffset = 1
		if len(nextItems) > 0 {
			nextItem := nextItems[0]
			switch nextItem.Typ {
			// speedup NodeString find
			case NodeString:
				if idx := strings.Index(part, nextItem.P); idx == -1 {
					// string not found, no need star scan
					break LOOP
				} else {
					nextOffset += idx
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
				break LOOP
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
			found = true
			part = part[n:]
		}
	case NodeRune:
		if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
			if _, ok := item.Runes[c]; ok {
				found = true
				part = part[n:]
			}
		}
	}
	if found {
		if part != "" && len(nextItems) > 0 {
			found = nextItems[0].matchItem(part, nextParts, nextItems[1:])
		} else if part != "" && len(nextItems) == 0 {
			found = false
		}
	}
	return
}

// nextInnerItem extract InnerItem
func nextInnerItem(s string) (item *InnerItem, next string, minLen int, maxLen int, err error) {
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
		runes := runesExpand([]rune(s))
		if len(runes) == 0 {
			return nil, s, 0, 0, ErrNodeMissmatch{NodeRune, s}
		}
		if len(runes) == 1 {
			var v string
			for k := range runes {
				v = string(k)
			}
			// one item optimization
			return &InnerItem{
				Typ: NodeString,
				P:   v,
			}, next, 1, 1, nil
		}
		return &InnerItem{
			Typ:   NodeRune,
			Runes: runes,
		}, next, 1, 1, nil
	case '{':
		if idx := strings.Index(s, "}"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		v := listExpand(s)
		if len(v) == 0 {
			return nil, s, 0, 0, ErrNodeMissmatch{NodeRune, s}
		}
		if len(v) == 1 {
			// one item optimization
			return &InnerItem{
				Typ: NodeString,
				P:   v[0],
			}, next, len(v[0]), len(v[0]), nil
		}
		sort.Strings(v)
		return &InnerItem{
			Typ:  NodeList,
			Vals: v,
		}, next, 1, 1, nil
	case '*':
		var next string
		for i, c := range s {
			if c != '*' {
				next = s[i:]
				break
			}
		}
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

	Terminated string // end of chain (resulting glob)

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
		if node.Terminated != "" {
			*items = append(*items, node.Terminated)
		} else if len(nextParts) > 0 {
			part, nextParts, _ = strings.Cut(nextParts, ".")
			for _, child := range node.Childs {
				child.Match(part, nextParts, items)
			}
		}
	}
}

// Merge is trying to merge inners
func (node *NodeItem) Merge(inners []*InnerItem) {
	if len(inners) == 1 {
		if node.Typ != NodeEmpty {
			// bug
			panic(fmt.Sprintf("%#v on %#v", inners[0], node))
		}
		// merge
		if inners[0].Typ == NodeString {
			if node.P != "" || node.Suffix != "" {
				var sb strings.Builder
				sb.Grow(len(node.P) + len(node.Suffix) + len(inners[0].P))
				sb.WriteString(node.P)
				sb.WriteString(inners[0].P)
				sb.WriteString(node.Suffix)
				node.Suffix = ""
				inners[0].P = sb.String()
			}
		} else {
			inners[0].P = node.P
		}

		node.InnerItem = *inners[0]

		return
	} else {
		var sb strings.Builder
		if inners[0].Typ == NodeString {
			// merge strings from prefix
			sb.Grow(len(node.P) + len(node.Suffix) + len(inners[0].P))
			sb.WriteString(node.P)
			sb.WriteString(inners[0].P)
			i := 1
			for i < len(inners) && inners[i].Typ == NodeString {
				sb.WriteString(inners[i].P)
				i++
			}
			inners = inners[i:]
			if len(inners) == 0 {
				//merge to string
				sb.WriteString(node.Suffix)
				node.P = sb.String()
				node.Suffix = ""
				node.Typ = NodeString

				return
			} else if len(inners) == 1 {
				//merge to prefix
				inners[0].P = sb.String()
				node.InnerItem = *inners[0]

				return
			}
			node.P = sb.String()
		}

		last := len(inners) - 1
		if inners[last].Typ == NodeString {
			//merge to suffix
			size := len(node.Suffix)
			i := last - 1
			for i > 1 && inners[i].Typ == NodeString {
				size += len(inners[i].P)
			}
			i++
			last = i
			sb.Reset()
			sb.Grow(size)
			for ; i < len(inners); i++ {
				sb.WriteString(inners[i].P)
			}
			sb.WriteString(node.Suffix)
			node.Suffix = sb.String()
			inners = inners[:last]
		}

		node.Typ = NodeInners
		node.Inners = inners
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
			if part == child.Node {
				node = child
				found = true
				break
			}
		}
		if !found {
			if i == last {
				// last node, so terminate match
				newNode = &NodeItem{Node: part, Terminated: glob}
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
					}
					newNode.Merge(inners)
				}
			}
			node.Childs = append(node.Childs, newNode)
			node = newNode
		}
	}

	if newNode != nil {
		if len(newNode.Childs) == 0 && newNode.Terminated == "" {
			// child  or/and terminated node
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
	if node, ok := w.Root[partsCount]; ok {
		globs = make([]string, 0, min(4, len(node.Childs)))
		for _, node := range node.Childs {
			part, nextParts, _ := strings.Cut(path, ".")
			// match first node
			node.Match(part, nextParts, &globs)
		}
	}

	return globs
}

func (w *GlobMatcher) MatchP(path string, globs *[]string) {
	*globs = (*globs)[:0]
	if path == "" {
		return
	}
	partsCount := pathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		for _, node := range node.Childs {
			part, nextParts, _ := strings.Cut(path, ".")
			// match first node
			node.Match(part, nextParts, globs)
		}
	}
}
