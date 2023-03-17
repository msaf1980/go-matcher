package gtags

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/msaf1980/go-matcher/pkg/items"
)

type testFindItems struct {
	key string
	// want
	item TaggedItems
}

func TestTaggedItemFindOrAppend(t *testing.T) {
	tests := []struct {
		item     *TaggedItem
		want     []testFindItems
		wantItem *TaggedItem
	}{
		{
			item: &TaggedItem{},
			want: []testFindItems{{key: "a", item: TaggedItems{Key: "a"}}},
			wantItem: &TaggedItem{
				Items: []TaggedItems{
					{
						Key:     "a",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<a>"}}},
					},
				},
			},
		},
		{
			item: &TaggedItem{
				Items: []TaggedItems{
					{
						Key:     "b",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "b=c"}}},
					},
				},
			},
			want: []testFindItems{
				{
					key: "b", item: TaggedItems{
						Key:     "b",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "b=c"}}},
					},
				},
			},
			wantItem: &TaggedItem{
				Items: []TaggedItems{
					{
						Key:     "b",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "b=c"}}},
					},
				},
			},
		},
		// all is added
		{
			item: &TaggedItem{},
			want: []testFindItems{
				{key: "aZ", item: TaggedItems{Key: "aZ"}},
				{key: "a", item: TaggedItems{Key: "a"}}, // new
				{
					key: "a", item: TaggedItems{
						Key:     "a",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<a>"}}}, // exist
					},
				},
				{key: "ba", item: TaggedItems{Key: "ba"}},
				{key: "c", item: TaggedItems{Key: "c"}},
				{key: "b", item: TaggedItems{Key: "b"}},
				{key: "z", item: TaggedItems{Key: "z"}},
				{key: "__name__", item: TaggedItems{Key: "__name__"}},
				{key: "__b__", item: TaggedItems{Key: "__b__"}},
				{key: "__a__", item: TaggedItems{Key: "__a__"}},
			},
			wantItem: &TaggedItem{
				Items: []TaggedItems{
					{
						Key:     "__name__",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<__name__>"}}},
					},
					{
						Key:     "__a__",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<__a__>"}}},
					},
					{
						Key:     "__b__",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<__b__>"}}},
					},
					{
						Key:     "a",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<a>"}}},
					},
					{
						Key:     "aZ",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<aZ>"}}},
					},
					{
						Key:     "b",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<b>"}}},
					},
					{
						Key:     "ba",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<ba>"}}},
					},
					{
						Key:     "c",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<c>"}}},
					},
					{
						Key:     "z",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "<z>"}}},
					},
				},
			},
		},
	}
	for n, tt := range tests {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			for _, want := range tt.want {
				pos := tt.item.findOrAppend(want.key)
				got := tt.item.Items[pos]
				if reflect.DeepEqual(want.item, got) {
					if len(tt.item.Items[pos].Matched) == 0 {
						// add fake item for test
						tt.item.Items[pos].Matched = append(tt.item.Items[pos].Matched, &TaggedItem{
							Terminated: items.Terminated{Query: "<" + want.key + ">"}},
						)
					}
				} else {
					t.Errorf("findOrAppendTaggedItems(%q) = %s", want.key, cmp.Diff(want.item, got))
				}
			}
			if !reflect.DeepEqual(tt.wantItem, tt.item) {
				t.Errorf("findOrAppendTaggedItems() = %s", cmp.Diff(tt.wantItem, tt.item))
			}
		})
	}
}

func TestFindTaggedItem(t *testing.T) {
	tests := []struct {
		item  *TaggedItem
		start int
		want  map[string]int
	}{
		{
			item: &TaggedItem{},
			want: map[string]int{"": -1, "a": -1},
		},
		{
			item: &TaggedItem{
				Items: []TaggedItems{
					{
						Key:     "b",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "b=c"}}},
					},
				},
			},
			want: map[string]int{"": -1, "a": -1, "b": 0},
		},
		{
			item: &TaggedItem{
				Items: []TaggedItems{
					{
						Key:     "a",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "a=c"}}},
					},
					{
						Key:     "b",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "b=c"}}},
					},
				},
			},
			want: map[string]int{"": -1, "a": 0, "b": 1, "c": -1},
		},
		{
			item: &TaggedItem{
				Items: []TaggedItems{
					{
						Key:     "a",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "a=c"}}},
					},
					{
						Key:     "b",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "b=c"}}},
					},
					{
						Key:     "ba",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "ba=c"}}},
					},
				},
			},
			want: map[string]int{"": -1, "a": 0, "b": 1, "ba": 2, "c": -1},
		},
		{
			item: &TaggedItem{
				Items: []TaggedItems{
					{
						Key:     "__name__",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "__name__=c"}}},
					},
					{
						Key:     "b",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "b=c"}}},
					},
					{
						Key:     "ba",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "ba=c"}}},
					},
				},
			},
			want: map[string]int{"": -1, "__name__": 0, "a": -1, "b": 1, "ba": 2, "c": -1},
		},
		{
			item: &TaggedItem{
				Items: []TaggedItems{
					{
						Key:     "__name__",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "__name__=c"}}},
					},
					{
						Key:     "__a__",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "__a__=c"}}},
					},
					{
						Key:     "a",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "a=c"}}},
					},
					{
						Key:     "aZ",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "aZ=c"}}},
					},
					{
						Key:     "b",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "b=c"}}},
					},
					{
						Key:     "ba",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "ba=c"}}},
					},
					{
						Key:     "c",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "c=c"}}},
					},
					{
						Key:     "z",
						Matched: []*TaggedItem{{Terminated: items.Terminated{Query: "z=c"}}},
					},
				},
			},
			want: map[string]int{
				"": -1, "__name__": 0, "__a__": 1, "a": 2, "b": 4, "ba": 5,
				"c": 6, "z": 7, "za": -1, "Z": -1,
			},
		},
	}
	for n, tt := range tests {
		for key, wantPos := range tt.want {
			t.Run(strconv.Itoa(n)+"#"+key, func(t *testing.T) {
				if pos := tt.item.find(key, tt.start); wantPos == pos {
					if pos != -1 {
						if tt.item.Items[pos].Key != key {
							t.Errorf("TaggedItems.find(%q, %d) pos mismatch = %q, want %q",
								key, tt.start, tt.item.Items[pos].Key, key,
							)
						}
					}
				} else {
					t.Errorf("TaggedItems.find(%q, %d) pos = %d, want %d",
						key, tt.start, pos, wantPos,
					)
				}
			})
		}
	}
}
