package items

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewList(t *testing.T) {
	tests := []struct {
		vals []string
		want Item
	}{
		// star replace all
		{
			vals: []string{"a", "*"},
			want: Star(0),
		},
		// deduplicate
		{
			vals: []string{"**", "*"},
			want: Star(0),
		},
		{
			vals: []string{"*??", "*?"},
			want: &List{
				MinSize: 1, MaxSize: -1,
				Vals: []Item{
					Star(2),
					Star(1),
				},
			},
		},
		// chains
		{
			vals: []string{"a", "b*"},
			want: &List{
				MinSize: 1, MaxSize: -1,
				Vals: []Item{
					NewString("a"),
					&Chain{
						Items: []Item{NewString("b"), Star(0)}, MinSize: 1, MaxSize: -1,
					},
				},
			},
		},
		{
			vals: []string{"a?cd*", "b*", "cde"},
			want: &List{
				MinSize: 1, MaxSize: -1,
				Vals: []Item{
					&Chain{
						Items:   []Item{NewString("a"), Any(1), NewString("cd"), Star(0)},
						MinSize: 4, MaxSize: -1,
					},
					&Chain{
						Items: []Item{NewString("b"), Star(0)}, MinSize: 1, MaxSize: -1,
					},
					NewString("cde"),
				},
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("%d#%#v", n, tt.vals), func(t *testing.T) {
			if got := NewItemList(tt.vals); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewList(%#v) = %s", tt.vals, cmp.Diff(tt.want, got))
			}
		})
	}
}
