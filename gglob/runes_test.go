package gglob

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_runesExpand(t *testing.T) {
	tests := []struct {
		runes string
		wantM map[rune]struct{}
	}{
		{runes: "", wantM: map[rune]struct{}{}},
		{runes: "[-q]", wantM: map[rune]struct{}{'q': {}}},
		{runes: "[a-c]", wantM: map[rune]struct{}{'a': {}, 'b': {}, 'c': {}}},
		{runes: "[a-cf]", wantM: map[rune]struct{}{'a': {}, 'b': {}, 'c': {}, 'f': {}}},
		{runes: "[za-c]", wantM: map[rune]struct{}{'a': {}, 'b': {}, 'c': {}, 'z': {}}},
		{runes: "[a-czA-C]", wantM: map[rune]struct{}{'a': {}, 'b': {}, 'c': {}, 'z': {}, 'A': {}, 'B': {}, 'C': {}}},
		// partially broken
		{runes: "[-]", wantM: map[rune]struct{}{}},
		{runes: "[a-]", wantM: map[rune]struct{}{'a': {}}},
		{runes: "[-a-c]", wantM: map[rune]struct{}{'a': {}, 'b': {}, 'c': {}}},
	}
	for _, tt := range tests {
		t.Run(tt.runes, func(t *testing.T) {
			if gotM := runesExpand([]rune(tt.runes)); !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("runesExpand(%s) = %s", tt.runes, cmp.Diff(tt.wantM, gotM))
			}
		})
	}
}
