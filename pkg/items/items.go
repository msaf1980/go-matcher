package items

import (
	"io"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/msaf1980/go-matcher/pkg/utils"
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

type InnerItem interface {
	Type() NodeType
	Match(part string, nextParts string, nextItems []InnerItem) (found bool)
}

type ItemRune map[rune]struct{}

func (ItemRune) Type() NodeType {
	return NodeRune
}

func (item ItemRune) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	if c, n := utf8.DecodeRuneInString(part); c != utf8.RuneError {
		if _, ok := item[c]; ok {
			found = true
			part = part[n:]
		}
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

type ItemString string

func (ItemString) Type() NodeType {
	return NodeString
}

func (item ItemString) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	s := string(item)
	if part == s {
		// full match
		found = true
		part = ""
	} else if strings.HasPrefix(part, s) {
		// strip prefix
		found = true
		part = part[len(s):]
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

type ItemOne struct{}

func (ItemOne) Type() NodeType {
	return NodeOne
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

type ItemStar struct{}

func (ItemStar) Type() NodeType {
	return NodeStar
}

func (item ItemStar) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
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
			switch v := nextItem.(type) {
			// speedup NodeString find
			case ItemString:
				s := string(v)
				if idx := strings.Index(part, s); idx == -1 {
					// string not found, no need star scan
					break LOOP
				} else {
					nextOffset += idx
					idx += len(s)
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
				found = nextItems[0].Match(part, nextParts, nextItems[1:])
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

type ItemList struct {
	// nodeList
	Vals    []string // strings
	ValsMin int      // min len in vals or min rune in range
	ValsMax int      // max len in vals or max rune in range
}

func (*ItemList) Type() NodeType {
	return NodeList
}

func (item *ItemList) Match(part string, nextParts string, nextItems []InnerItem) (found bool) {
	// TODO: nodeList Skip scan
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
				found = nextItems[0].Match(part, nextParts, nextItems[1:])
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

// NextInnerItem extract InnerItem
func NextInnerItem(s string) (item InnerItem, next string, minLen int, maxLen int, err error) {
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
		runes, failed := RunesExpand([]rune(s))
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
			return ItemString(v), next, 1, 1, nil
		}
		return ItemRune(runes), next, 1, 1, nil
	case '{':
		if idx := strings.Index(s, "}"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		vals, failed := ListExpand(s)
		if failed {
			return nil, s, 0, 0, ErrNodeMissmatch{NodeRune, s}
		}
		if len(vals) == 0 {
			return nil, next, 0, 0, nil
		}
		if len(vals) == 1 {
			// one item optimization
			return ItemString(vals[0]), next, len(vals[0]), len(vals[0]), nil
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
		return &ItemList{Vals: vals, ValsMin: minLen, ValsMax: maxLen}, next, minLen, maxLen, nil
	case '*':
		var next string
		for i, c := range s {
			if c != '*' {
				next = s[i:]
				break
			}
		}
		return ItemStar{}, next, 0, 0, nil
	case '?':
		next := s[1:]
		return ItemOne{}, next, 1, 1, nil
	case ']', '}':
		return nil, s, 0, 0, ErrNodeUnclosed{s}
	default:
		// string segment
		end := IndexWildcard(s)
		v, next := SplitString(s, end)
		return ItemString(v), next, len(v), len(v), nil
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

	Inners []InnerItem // inner segments
	Childs []*NodeItem // next possible parts tree
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

// Merge is trying to merge inners
func (node *NodeItem) Merge(inners []InnerItem) {
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
		case ItemString:
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
		case ItemString:
			var sb strings.Builder
			// merge strings from prefix
			s := string(v)
			sb.Grow(len(node.P) + len(node.Suffix) + len(s))
			sb.WriteString(node.P)
			sb.WriteString(s)
			i := 1
			for i < len(inners) && inners[i].Type() == NodeString {
				s := string(inners[i].(ItemString))
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
		case ItemString:
			//merge to suffix
			size := len(node.Suffix) + len(v)
			i := last - 1
			for i > 0 && inners[i].Type() == NodeString {
				size += len(inners[i].(ItemString))
				i--
			}
			i++
			last = i
			var sb strings.Builder
			sb.Grow(size)
			for ; i < len(inners); i++ {
				sb.WriteString(string(inners[i].(ItemString)))
			}
			sb.WriteString(node.Suffix)
			node.Suffix = sb.String()
			inners = inners[:last]
		}

		node.Inners = inners
	}
}

func (node *NodeItem) Match(path string, matched *[]string) {
	for _, node := range node.Childs {
		part, nextParts, _ := strings.Cut(path, ".")
		// match first node
		node.MatchItems(part, nextParts, matched)
	}
}

func (node *NodeItem) MatchByParts(parts []string, matched *[]string) {
	for _, node := range node.Childs {
		// match first node
		node.MatchItemsPart(parts[0], parts[1:], matched)
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

					for part != "" {
						inner, part, min, max, err = NextInnerItem(part)
						if err != nil {
							return
						}
						if inner == nil {
							continue
						}
						lastNode.MinSize += min
						if lastNode.MaxSize != -1 {
							if inner.Type() == NodeStar {
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

func ParseItems(root map[int]*NodeItem, glob string) (lastNode *NodeItem, err error) {
	glob, partsCount := PathLevel(glob)

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
	path, partsCount := PathLevel(path)
	if node, ok := w.Root[partsCount]; ok {
		globs = make([]string, 0, utils.Min(4, len(node.Childs)))
		node.Match(path, &globs)
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
		node.Match(path, globs)
	}
}

// NextWildcardItem extract InnerItem
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
		runes, failed := RunesExpand([]rune(s))
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
			return ItemString(v), next, 1, 1, nil
		}
		return ItemRune(runes), next, 1, 1, nil
	case '{':
		if idx := strings.Index(s, "}"); idx != -1 {
			idx++
			next = s[idx:]
			s = s[:idx]
		}
		vals, failed := ListExpand(s)
		if failed {
			return nil, s, 0, 0, ErrNodeMissmatch{NodeRune, s}
		}
		if len(vals) == 0 {
			return nil, next, 0, 0, nil
		}
		if len(vals) == 1 {
			// one item optimization
			return ItemString(vals[0]), next, len(vals[0]), len(vals[0]), nil
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
		return &ItemList{Vals: vals, ValsMin: minLen, ValsMax: maxLen}, next, minLen, maxLen, nil
	case '*':
		var next string
		for i, c := range s {
			if c != '*' {
				next = s[i:]
				break
			}
		}
		return ItemStar{}, next, 0, 0, nil
	case '?':
		next := s[1:]
		return ItemOne{}, next, 1, 1, nil
	case ']', '}':
		return nil, s, 0, 0, ErrNodeUnclosed{s}
	default:
		// string segment
		end := IndexWildcard(s)
		v, next := SplitString(s, end)
		return ItemString(v), next, len(v), len(v), nil
	}
}
