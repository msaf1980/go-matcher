package gglob

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestNodeItem_Merge(t *testing.T) {
	tests := []struct {
		name       string
		node       *items.NodeItem
		inners     []items.InnerItem
		wantNode   *items.NodeItem
		matchGlobs map[string][]string // must match with glob
		miss       []string
	}{
		{
			name: "merge strings #all",
			node: &items.NodeItem{
				Node: "a[a-]Z[Q]", Terminated: "a[a-]Z[Q]", MinSize: 4, MaxSize: 4,
				P: "a",
			},
			inners: []items.InnerItem{
				items.ItemString("a"), items.ItemString("Z"), items.ItemString("Q"),
			},
			wantNode: &items.NodeItem{
				Node: "a[a-]Z[Q]", Terminated: "a[a-]Z[Q]", P: "aaZQ", MinSize: 4, MaxSize: 4,
			},
			matchGlobs: map[string][]string{
				"aaZQ": {"a[a-]Z[Q]"},
			},
			miss: []string{"", "ab", "aaZQa"},
		},
		{
			name: "merge strings #prefix",
			node: &items.NodeItem{
				Node: "a[a-]Z[Q]*", Terminated: "a[a-]Z[Q]*", P: "a", MinSize: 4, MaxSize: -1,
			},
			inners: []items.InnerItem{
				items.ItemString("a"), items.ItemString("Z"), items.ItemString("Q"), items.ItemStar{},
			},
			wantNode: &items.NodeItem{
				Node: "a[a-]Z[Q]*", Terminated: "a[a-]Z[Q]*", P: "aaZQ", MinSize: 4, MaxSize: -1,
				Inners: []items.InnerItem{items.ItemStar{}},
			},
			matchGlobs: map[string][]string{
				"aaZQ":  {"a[a-]Z[Q]*"},
				"aaZQa": {"a[a-]Z[Q]*"},
			},
			miss: []string{"", "ab", "aaZqa"},
		},
		{
			name: "merge strings #suffix",
			node: &items.NodeItem{
				Node: "a[a-]Z[Q]*[z-]l", Terminated: "a[a-]Z[Q]*[z-]l", MinSize: 6, MaxSize: -1,
				P: "a", Suffix: "l",
			},
			inners: []items.InnerItem{
				items.ItemString("a"), items.ItemString("Z"), items.ItemString("Q"), items.ItemStar{}, items.ItemString("z"),
			},
			wantNode: &items.NodeItem{
				Node: "a[a-]Z[Q]*[z-]l", Terminated: "a[a-]Z[Q]*[z-]l", MinSize: 6, MaxSize: -1,
				P: "aaZQ", Suffix: "zl",
				Inners: []items.InnerItem{items.ItemStar{}},
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
				t.Errorf("items.NodeItem.Merge() = %s", cmp.Diff(tt.wantNode, tt.node))
			} else {
				w := NewGlobMatcher()
				rootNode := &items.NodeItem{}
				w.Root[1] = rootNode
				rootNode.Childs = append(rootNode.Childs, tt.node)
				verifyGlobMatcher(t, tt.matchGlobs, tt.miss, w)
			}
		})
	}
}
