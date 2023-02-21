package gglob

import (
	"sort"
	"strings"

	"github.com/msaf1980/go-matcher/pkg/utils"
	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

// NodeItem contains pattern node item (with childs and fixed depth)
type NodeItem struct {
	Node string // raw string (or full glob for terminated)

	Terminated []string // end of chain (resulting glob)
	TermIndex  []int    // rule num of end of chain (resulting glob), can be used in specific cases

	wildcards.WildcardItems

	// TODO: may be some ordered tree for complete string nodes search speedup (on large set) ?
	Childs []*NodeItem // next possible parts slice
}

func (node *NodeItem) WriteString(buf *strings.Builder) {
	buf.WriteString(node.P)
	wildcards.WriteInnerItems(node.Inners, buf)
	buf.WriteString(node.Suffix)
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
			// bug, can not there, must be megred with prefix
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

		matched = node.Inners[0].Match(part, node.Inners[1:])
	}
	return
}

func (node *NodeItem) MatchItems(part string, nextParts string, matched *[]string) {
	if node.MatchNode(part) {
		if len(node.Terminated) > 0 {
			*matched = append(*matched, node.Terminated...)
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
		if len(node.TermIndex) > 0 {
			*matched = append(*matched, node.TermIndex...)
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
		if len(node.Terminated) > 0 {
			if *minMatched == -1 || *minMatched > node.TermIndex[0] {
				*minMatched = node.TermIndex[0]
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
		if len(node.TermIndex) > 0 {
			*matched = append(*matched, node.TermIndex...)
		} else if len(parts) > 0 {
			for _, child := range node.Childs {
				child.MatchIndexedItemsPart(parts[0], parts[1:], matched)
			}
		}
	}
}

func (node *NodeItem) MatchFirstItemsPart(part string, parts []string, minMatched *int) {
	if node.MatchNode(part) {
		if len(node.Terminated) > 0 {
			if *minMatched == -1 || *minMatched > node.TermIndex[0] {
				*minMatched = node.TermIndex[0]
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
		if len(node.Terminated) > 0 {
			*matched = append(*matched, node.Terminated...)
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
func (node *NodeItem) ParseNode(glob string, termIdx int, buf *strings.Builder) (newGlob string, lastNode *NodeItem, err error) {
	var (
		i    int
		part string
	)

	nextParts := glob
	for nextParts != "" {
		found := false
		part, nextParts, _ = strings.Cut(nextParts, ".")
		if part == "" {
			return glob, nil, wildcards.ErrNodeEmpty{glob}
		}

		lastNode = &NodeItem{Node: part}
		if err = lastNode.Parse(part); err != nil {
			return
		}

		// write glob segments
		if buf.Len() > 0 {
			buf.WriteRune('.')
		}
		start := buf.Len()
		lastNode.WriteString(buf)
		if nextParts == "" {
			// end of glob
			// TODO : get Terminated from buffer
			lastNode.Terminated = append(lastNode.Terminated, glob)
			if termIdx > -1 {
				lastNode.TermIndex = append(lastNode.TermIndex, termIdx)
				sort.Ints(lastNode.TermIndex)
			}
			nodeName := buf.String()
			if nodeName == glob {
				newGlob = glob
			} else {
				newGlob = string(utils.CloneString(nodeName))
				lastNode.Terminated = append(lastNode.Terminated, newGlob)
			}
		}
		nodeName := buf.String()[start:]
		if nodeName != lastNode.Node {
			lastNode.Node = string(utils.CloneString(nodeName))
		}

		for _, child := range node.Childs {
			if lastNode.Node == child.Node {
				node = child
				found = true
				if nextParts == "" {
					node.Terminated = append(node.Terminated, glob)
					if termIdx > -1 {
						node.TermIndex = append(node.TermIndex, termIdx)
						sort.Ints(node.TermIndex)
					}
				}
				break
			}
		}
		if !found {
			node.Childs = append(node.Childs, lastNode)
			node = lastNode
		}
		i++
	}
	if lastNode == nil {
		err = wildcards.ErrGlobNotExpanded{glob}
	}
	return
}
