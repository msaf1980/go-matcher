package wildcards

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func Test_runesExpand(t *testing.T) {
	tests := []struct {
		s          string
		want       ItemRuneRanges
		wantFailed bool
	}{
		{s: "[-q]", want: ItemRuneRanges{{'q', 'q'}}, wantFailed: false},
		{s: "[a-c]", want: ItemRuneRanges{{'a', 'c'}}, wantFailed: false},
		{s: "[a-cf]", want: ItemRuneRanges{{'a', 'c'}, {'f', 'f'}}, wantFailed: false},
		{s: "[za-c]", want: ItemRuneRanges{{'a', 'c'}, {'z', 'z'}}, wantFailed: false},
		{s: "[a-czA-C]", want: ItemRuneRanges{{'A', 'C'}, {'a', 'c'}, {'z', 'z'}}, wantFailed: false},
		// partially broken
		{s: "[-]", want: ItemRuneRanges{}, wantFailed: false},
		{s: "[a-]", want: ItemRuneRanges{{'a', 'a'}}, wantFailed: false},
		{s: "[-a-c]", want: ItemRuneRanges{{'a', 'c'}}, wantFailed: false},
		// duplicated ranges
		{s: "[a-cb-ld-fa-c]", want: ItemRuneRanges{{'a', 'l'}}, wantFailed: false},
		{s: "[1-32-45-76-78-95789]", want: ItemRuneRanges{{'1', '9'}}, wantFailed: false},
		// broken
		{s: "", wantFailed: true},
		{s: "[a", wantFailed: true},
		{s: "a]", wantFailed: true},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got, gotSuccess := RunesExpand(tt.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runesExpand(%s) = %s", tt.s, cmp.Diff(tt.want, got))
			}
			assert.Equal(t, tt.wantFailed, gotSuccess)
		})
	}
}
