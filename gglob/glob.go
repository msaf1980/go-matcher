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
	NodeMany   // *
	NodeInners // composite type, contains prefix, suffix, subitems in []inners
)

type InnerItem struct {
	// can be nodeString, nodeList, nodeRange, nodeOne, nodeMany
	Typ NodeType

	// string (nodeString)
	S string

	// nodeList or 	nodeRange
	Vals    []string // strings for nodeList
	ValsMin int      // min len in vals or min rune in range
	ValsMax int      // max len in vals or max rune in range
}

func (item *InnerItem) Match(s string, nextItems []*InnerItem) (found bool, matched string, next string) {
	switch item.Typ {
	case NodeString:
		if s == item.S {
			// full match
			found, matched = true, s
		} else if strings.HasPrefix(s, item.S) {
			// strip prefix
			found, matched, next = true, item.S, s[len(item.S):]
		} else {
			next = s
		}
		// TODO: other types
	case NodeMany:
		if len(nextItems) > 1 {
			// TODO: * in multipart
		} else {
			// all of string matched to *
			found, matched = true, s
		}
	case NodeOne:
		if len(nextItems) > 1 {
			if c, n := utf8.DecodeRuneInString(s); c != utf8.RuneError {
				if len(s) > n {
					found, matched, next = true, s[:n], s[n:]
				}
			}
		} else if len(s) == 1 {
			// string matched to ? must have one element
			found, matched = true, s
		}
	}
	return
}

func nextInnerItem(s string) (*InnerItem, string, error) {
	if s == "" {
		return nil, s, io.EOF
	}
	switch s[0] {
	case '[':
		// TODO: implement nodeRange
		return nil, s, ErrNodeUnclosed{s}
	case '{':
		return nil, s, ErrNodeUnclosed{s}
	case '*':
		next := nextString(s, 1)
		return &InnerItem{
			Typ: NodeMany,
		}, next, nil
	case '?':
		next := nextString(s, 1)
		return &InnerItem{
			Typ: NodeOne,
		}, next, nil
	case ']', '}':
		return nil, s, ErrNodeUnclosed{s}
	default:
		// string segment
		end := IndexWildcard(s)
		v, next := splitString(s, end)
		return &InnerItem{
			Typ: NodeString, S: v,
		}, next, nil
	}
}

// NodeItem contains pattern node item
type NodeItem struct {
	Node string // raw string (or full glob for terminated)

	Terminated bool // end of chain

	InnerItem // if one item, no need to use []inners

	Inners []*InnerItem // inner segments
	Childs []*NodeItem  // next possible parts tree
}

func (node *NodeItem) AddMatched(parts []string, items *[]string) {
	if node.Typ == NodeInners {
		if len(node.Inners) == 0 {
			// some broken, skip node
			return
		}
		part := parts[0]
		var found bool
		for i, inner := range node.Inners {
			if found, _, part = inner.Match(part, node.Inners[i:]); found {

			} else {
				// item not found
				return
			}
		}
		if part != "" {
			// partial match
			return
		}
	} else if node.Typ == NodeString {
		if node.S != parts[0] {
			return
		}
	} else if found, _, part := node.Match(parts[0], nil); !found || part != "" {
		// not matched or partial match
		return
	}

	if node.Terminated {
		*items = append(*items, node.Node)
	} else {
		if len(parts) == 1 {
			return
		}
		parts = parts[1:]
		for _, child := range node.Childs {
			child.AddMatched(parts, items)
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
				newNode.S = part
			} else {
				switch part {
				case "*":
					newNode.Typ = NodeMany
				case "?":
					newNode.Typ = NodeOne
				default:
					var inner *InnerItem
					innerCount := WildcardCount(part) + 1
					newNode.Inners = make([]*InnerItem, 0, innerCount)
					for part != "" {
						inner, part, err = nextInnerItem(part)
						if err != nil {
							return
						}
						newNode.Inners = append(newNode.Inners, inner)
					}
					newNode.Typ = NodeInners
					if len(newNode.Inners) == 0 {
						// no inners for inner node
						return ErrGlobNotExpanded{newNode.Node}
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
	parts := pathSplit(path)
	if hasEmptyParts(parts) {
		return nil
	}
	var items []string
	if node, ok := w.Root[len(parts)]; ok {
		items = make([]string, 0, min(4, len(node.Childs)))
		for _, node := range node.Childs {
			// match first node
			node.AddMatched(parts, &items)
		}
	}

	return items
}
