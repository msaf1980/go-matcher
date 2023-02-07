package items

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func Test_runesExpand(t *testing.T) {
	tests := []struct {
		runes      string
		wantM      map[rune]struct{}
		wantFailed bool
	}{
		{runes: "[-q]", wantM: map[rune]struct{}{'q': {}}, wantFailed: false},
		{runes: "[a-c]", wantM: map[rune]struct{}{'a': {}, 'b': {}, 'c': {}}, wantFailed: false},
		{runes: "[a-cf]", wantM: map[rune]struct{}{'a': {}, 'b': {}, 'c': {}, 'f': {}}, wantFailed: false},
		{runes: "[za-c]", wantM: map[rune]struct{}{'a': {}, 'b': {}, 'c': {}, 'z': {}}, wantFailed: false},
		{runes: "[a-czA-C]", wantM: map[rune]struct{}{'a': {}, 'b': {}, 'c': {}, 'z': {}, 'A': {}, 'B': {}, 'C': {}}, wantFailed: false},
		// partially broken
		{runes: "[-]", wantM: map[rune]struct{}{}, wantFailed: false},
		{runes: "[a-]", wantM: map[rune]struct{}{'a': {}}, wantFailed: false},
		{runes: "[-a-c]", wantM: map[rune]struct{}{'a': {}, 'b': {}, 'c': {}}, wantFailed: false},
		// broken
		{runes: "", wantFailed: true},
		{runes: "[a", wantFailed: true},
		{runes: "a]", wantFailed: true},
	}
	for _, tt := range tests {
		t.Run(tt.runes, func(t *testing.T) {
			gotM, gotSuccess := RunesExpand([]rune(tt.runes))
			if !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("runesExpand(%s) = %s", tt.runes, cmp.Diff(tt.wantM, gotM))
			}
			assert.Equal(t, tt.wantFailed, gotSuccess)
		})
	}
}
