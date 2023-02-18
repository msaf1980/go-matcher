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

	wildcards.WildcardItems

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
func (node *NodeItem) ParseNode(glob string, partsCount int, termIdx int) (lastNode *NodeItem, err error) {
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
			if err = lastNode.Parse(part); err != nil {
				return
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
