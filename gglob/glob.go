package gglob

import (
	"io"
	"math"
	"strings"

	"github.com/msaf1980/go-matcher/pkg/items"
)

// nextInnerItem extract InnerItem
func nextInnerItem(s string) (item items.InnerItem, next string, minLen int, maxLen int, err error) {
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
		runes, failed := items.RunesExpand([]rune(s))
		if failed {
			return nil, s, 0, 0, items.ErrNodeMissmatch{items.NodeRune, s}
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
			return items.ItemString(v), next, 1, 1, nil
		}
		return items.ItemRune(runes), next, 1, 1, nil
	case '{':
		if idx := strings.Index(s, "}"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		vals, failed := items.ListExpand(s)
		if failed {
			return nil, s, 0, 0, items.ErrNodeMissmatch{items.NodeRune, s}
		}
		if len(vals) == 0 {
			return nil, next, 0, 0, nil
		}
		if len(vals) == 1 {
			// one item optimization
			return items.ItemString(vals[0]), next, len(vals[0]), len(vals[0]), nil
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
		return &items.ItemList{Vals: vals, ValsMin: minLen, ValsMax: maxLen}, next, minLen, maxLen, nil
	case '*':
		var next string
		for i, c := range s {
			if c != '*' {
				next = s[i:]
				break
			}
		}
		return items.ItemStar{}, next, 0, 0, nil
	case '?':
		next := s[1:]
		return items.ItemOne{}, next, 1, 1, nil
	case ']', '}':
		return nil, s, 0, 0, items.ErrNodeUnclosed{s}
	default:
		// string segment
		end := items.IndexWildcard(s)
		v, next := items.SplitString(s, end)
		return items.ItemString(v), next, len(v), len(v), nil
	}
}

// NodeItem contains pattern node item
type NodeItem struct {
	Node string // raw string (or full glob for terminated)

	Terminated string // end of chain (resulting glob)

	// size check optimization
	MinSize int
	MaxSize int // 0 or -1 for unlimited

	P      string // prefix or full string if len(inners) == 0
	Suffix string // suffix

	Inners []items.InnerItem // inner segments
	Childs []*NodeItem       // next possible parts tree
}

func (node *NodeItem) Match(part string, nextParts string, items *[]string) {
	var found bool

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
		found = (node.P == part)
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

		found = node.Inners[0].Match(part, nextParts, node.Inners[1:])
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
func (node *NodeItem) Merge(inners []items.InnerItem) {
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
			for i < len(inners) && inners[i].Type() == items.NodeString {
				s := string(inners[i].(items.ItemString))
				sb.WriteString(s)
				i++
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
			for i > 0 && inners[i].Type() == items.NodeString {
				size += len(inners[i].(items.ItemString))
				i--
			}
			i++
			last = i
			var sb strings.Builder
			sb.Grow(size)
			for ; i < len(inners); i++ {
				sb.WriteString(string(inners[i].(items.ItemString)))
			}
			sb.WriteString(node.Suffix)
			node.Suffix = sb.String()
			inners = inners[:last]
		}

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

func (node *NodeItem) Parse(glob string, partsCount int) (lastNode *NodeItem, err error) {
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
			return nil, items.ErrNodeEmpty{glob}
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
			pos := items.IndexWildcard(part)
			if pos == -1 {
				lastNode.Inners = []items.InnerItem{items.ItemString(part)}
			} else {
				if pos > 0 {
					lastNode.P = part[:pos] // prefix
					part = part[pos:]
					lastNode.MinSize = len(lastNode.P)
					lastNode.MaxSize = len(lastNode.P)
				}
				end := items.IndexLastWildcard(part)
				if end == 0 && part[0] != '?' && part[0] != '*' {
					err = items.ErrNodeUnclosed{part}
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
					lastNode.Inners = []items.InnerItem{items.ItemStar{}}
					lastNode.MaxSize = -1 // unlimited
				case "?":
					lastNode.Inners = []items.InnerItem{items.ItemOne{}}
					lastNode.MinSize++
					if lastNode.MaxSize != -1 {
						lastNode.MaxSize++
					}
				default:
					var (
						inner    items.InnerItem
						min, max int
					)
					innerCount := items.WildcardCount(part)
					inners := make([]items.InnerItem, 0, innerCount)

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
							if inner.Type() == items.NodeStar {
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
		err = items.ErrGlobNotExpanded{glob}
		return
	}
	if i != partsCount || (len(lastNode.Childs) == 0 && lastNode.Terminated == "") {
		// child  or/and terminated node
		err = items.ErrNodeNotEnd{lastNode.Node}
	}
	return
}

func ParseItems(root map[int]*NodeItem, glob string) (lastNode *NodeItem, err error) {
	glob, partsCount := items.PathLevel(glob)

	node, ok := root[partsCount]
	if !ok {
		node = &NodeItem{}
		root[partsCount] = node
	}
	_, err = node.Parse(glob, partsCount)

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
	path, partsCount := items.PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		globs = make([]string, 0, items.Min(4, len(node.Childs)))
		node.MatchRoot(path, &globs)
	}

	return globs
}

func (w *GlobMatcher) MatchP(path string, globs *[]string) {
	*globs = (*globs)[:0]
	if path == "" {
		return
	}
	path, partsCount := items.PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		node.MatchRoot(path, globs)
	}
}
