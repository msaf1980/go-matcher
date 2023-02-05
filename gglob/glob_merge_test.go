package gglob

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNodeItem_Merge(t *testing.T) {
	tests := []struct {
		name       string
		node       *NodeItem
		inners     []*InnerItem
		wantNode   *NodeItem
		matchGlobs map[string][]string // must match with glob
		miss       []string
	}{
		{
			name: "merge strings #all",
			node: &NodeItem{
				Node: "a[a-]Z[Q]", Terminated: "a[a-]Z[Q]", MinSize: 4, MaxSize: 4,
				InnerItem: InnerItem{P: "a"},
			},
			inners: []*InnerItem{
				{Typ: NodeString, P: "a"}, {Typ: NodeString, P: "Z"}, {Typ: NodeString, P: "Q"},
			},
			wantNode: &NodeItem{
				Node: "a[a-]Z[Q]", Terminated: "a[a-]Z[Q]", MinSize: 4, MaxSize: 4,
				InnerItem: InnerItem{Typ: NodeString, P: "aaZQ"},
			},
			matchGlobs: map[string][]string{
				"aaZQ": {"a[a-]Z[Q]"},
			},
			miss: []string{"", "ab", "aaZQa"},
		},
		{
			name: "merge strings #prefix",
			node: &NodeItem{
				Node: "a[a-]Z[Q]*", Terminated: "a[a-]Z[Q]*", MinSize: 4, MaxSize: -1,
				InnerItem: InnerItem{P: "a"},
			},
			inners: []*InnerItem{
				{Typ: NodeString, P: "a"}, {Typ: NodeString, P: "Z"}, {Typ: NodeString, P: "Q"}, {Typ: NodeStar},
			},
			wantNode: &NodeItem{
				Node: "a[a-]Z[Q]*", Terminated: "a[a-]Z[Q]*", MinSize: 4, MaxSize: -1,
				InnerItem: InnerItem{Typ: NodeStar, P: "aaZQ"},
			},
			matchGlobs: map[string][]string{
				"aaZQ":  {"a[a-]Z[Q]*"},
				"aaZQa": {"a[a-]Z[Q]*"},
			},
			miss: []string{"", "ab", "aaZqa"},
		},
		{
			name: "merge strings #suffix",
			node: &NodeItem{
				Node: "a[a-]Z[Q]*[z-]l", Terminated: "a[a-]Z[Q]*[z-]l", MinSize: 6, MaxSize: -1,
				InnerItem: InnerItem{P: "a"}, Suffix: "l",
			},
			inners: []*InnerItem{
				{Typ: NodeString, P: "a"}, {Typ: NodeString, P: "Z"}, {Typ: NodeString, P: "Q"},
				{Typ: NodeStar}, {Typ: NodeString, P: "z"},
			},
			wantNode: &NodeItem{
				Node: "a[a-]Z[Q]*[z-]l", Terminated: "a[a-]Z[Q]*[z-]l", MinSize: 6, MaxSize: -1,
				InnerItem: InnerItem{Typ: NodeInners, P: "aaZQ"}, Suffix: "zl",
				Inners: []*InnerItem{
					{Typ: NodeStar},
				},
			},
			matchGlobs: map[string][]string{
				"aaZQzl":  {"a[a-]Z[Q]*[z-]l"},
				"aaZQazl": {"a[a-]Z[Q]*[z-]l"},
			},
			miss: []string{"", "ab", "aaZqa"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.node.Merge(tt.inners)
			if !reflect.DeepEqual(tt.wantNode, tt.node) {
				t.Errorf("NodeItem.Merge() = %s", cmp.Diff(tt.wantNode, tt.node))
			} else {
				w := NewGlobMatcher()
				rootNode := &NodeItem{InnerItem: InnerItem{Typ: NodeRoot}}
				w.Root[1] = rootNode
				rootNode.Childs = append(rootNode.Childs, tt.node)
				verifyGlobMatcher(t, tt.matchGlobs, tt.miss, w)
			}
		})
	}
}
