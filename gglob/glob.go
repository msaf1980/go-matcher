package gglob

import (
	"fmt"
	"io"
	"math"
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
				// TODO: may be other optimization: may be for list
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

func (item *InnerItem) matchList(part string, nextParts string, nextItems []*InnerItem) (found bool) {
	l := len(part)
	if l < item.ValsMin {
		return
	}
	if len(nextItems) == 0 && l > item.ValsMax {
		return
	}
	// TODO: may be optimize scan of duplicate with prefix tree ?
LOOP:
	for _, s := range item.Vals {
		part := part
		if part == s {
			// full match
			found = true
			part = ""
		} else if strings.HasPrefix(part, s) {
			// strip prefix
			found = true
			part = part[len(s):]
		} else {
			// try to next
			continue
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
	case NodeList:
		return item.matchList(part, nextParts, nextItems)
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
		runes, failed := runesExpand([]rune(s))
		if failed {
			return nil, s, 0, 0, ErrNodeMissmatch{NodeRune, s}
		}
		if len(runes) == 0 {
			return nil, next, 0, 0, nil
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
		vals, failed := listExpand(s)
		if failed {
			return nil, s, 0, 0, ErrNodeMissmatch{NodeRune, s}
		}
		if len(vals) == 0 {
			return nil, next, 0, 0, nil
		}
		if len(vals) == 1 {
			// one item optimization
			return &InnerItem{
				Typ: NodeString,
				P:   vals[0],
			}, next, len(vals[0]), len(vals[0]), nil
		}
		minLen := math.MaxInt
		maxLen := 0
		for _, v := range vals {
			l := len(v)
			if maxLen < l {
				maxLen = l
			}
			if minLen > l {
				minLen = l
			}
		}
		return &InnerItem{
			Typ:  NodeList,
			Vals: vals, ValsMin: minLen, ValsMax: maxLen,
		}, next, minLen, maxLen, nil
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
	if len(inners) == 0 {
		if node.Typ != NodeString {
			// no inners for inner node, may be for like []
			node.Typ = NodeString
			if node.P != "" && node.Suffix != "" {
				node.P += node.Suffix
				node.Suffix = ""
			} else if node.Suffix != "" {
				node.P = node.Suffix
				node.Suffix = ""
			}
		}
	} else if len(inners) == 1 {
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

func (node *NodeItem) MatchRoot(path string, globs *[]string) {
	for _, node := range node.Childs {
		part, nextParts, _ := strings.Cut(path, ".")
		// match first node
		node.Match(part, nextParts, globs)
	}
}

func ParseItems(root map[int]*NodeItem, glob string) (lastNode *NodeItem, err error) {
	glob, partsCount := PathLevel(glob)

	node, ok := root[partsCount]
	if !ok {
		node = &NodeItem{InnerItem: InnerItem{Typ: NodeRoot}}
		root[partsCount] = node
	}
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
			// TODO: may be normalize parts for equals like {a,z} and {z,a} ?
			if part == child.Node {
				node = child
				found = true
				break
			}
		}
		if !found {
			if i == last {
				// last node, so terminate match
				lastNode = &NodeItem{Node: part, Terminated: glob}
			} else {
				lastNode = &NodeItem{Node: part}
			}
			pos := IndexWildcard(part)
			if pos == -1 {
				lastNode.Typ = NodeString
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
					lastNode.Typ = NodeStar
					lastNode.MaxSize = -1 // unlimited
				case "?":
					lastNode.Typ = NodeOne
					lastNode.MinSize++
					if lastNode.MaxSize != -1 {
						lastNode.MaxSize++
					}
				default:
					var (
						inner    *InnerItem
						min, max int
					)
					innerCount := WildcardCount(part)
					inners := make([]*InnerItem, 0, innerCount)

					for part != "" {
						inner, part, min, max, err = nextInnerItem(part)
						if err != nil {
							return
						}
						if inner == nil {
							continue
						}
						lastNode.MinSize += min
						if lastNode.MaxSize != -1 {
							if inner.Typ == NodeStar {
								lastNode.MaxSize = -1
							} else {
								lastNode.MaxSize += max
							}
						}
						inners = append(inners, inner)
					}
					lastNode.Merge(inners)
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
	if _, err = ParseItems(w.Root, glob); err != nil {
		return err
	}

	w.Globs[glob] = true

	return
}

func (w *GlobMatcher) Match(path string) (globs []string) {
	if path == "" {
		return nil
	}
	path, partsCount := PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		globs = make([]string, 0, min(4, len(node.Childs)))
		node.MatchRoot(path, &globs)
	}

	return globs
}

func (w *GlobMatcher) MatchP(path string, globs *[]string) {
	*globs = (*globs)[:0]
	if path == "" {
		return
	}
	path, partsCount := PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		node.MatchRoot(path, globs)
	}
}
