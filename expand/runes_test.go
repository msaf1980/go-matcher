package expand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunesRanges_Expand(t *testing.T) {
	tests := []struct {
		s          string
		want       []runes
		wantFailed bool
	}{
		{s: "a-c", want: []runes{{'a', 'c', 'a'}}},
		{s: "a-cf", want: []runes{{'a', 'c', 'a'}, {'f', 'f', 'f'}}},
		{s: "za-c", want: []runes{{'a', 'c', 'a'}, {'z', 'z', 'z'}}},
		{s: "Za-c", want: []runes{{'Z', 'Z', 'Z'}, {'a', 'c', 'a'}}},
		{s: "a-czA-C", want: []runes{{'A', 'C', 'A'}, {'a', 'c', 'a'}, {'z', 'z', 'z'}}},
		{s: "a-czqj-mA-C", want: []runes{{'A', 'C', 'A'}, {'a', 'c', 'a'}, {'j', 'm', 'j'}, {'q', 'q', 'q'}, {'z', 'z', 'z'}}},
		// contains unicode
		{s: "界Э-ЯW-Z", want: []runes{{'W', 'Z', 'W'}, {'Э', 'Я', 'Э'}, {'界', '界', '界'}}},
		// partially broken
		{s: "-q", want: []runes{{'q', 'q', 'q'}}},
		{s: "a-", want: []runes{{'a', 'a', 'a'}}},
		{s: "-a-c", want: []runes{{'a', 'c', 'a'}}},
		// duplicated ASCII ranges
		{s: "a-cb-ld-fa-c", want: []runes{{'a', 'l', 'a'}}},
		{s: "1-32-45-76-78-95789", want: []runes{{'1', '9', '1'}}},
		// duplicated Unicode ranges
		{s: "а-йв-у1-9b-dА-ДО-П你好世", want: []runes{
			{'1', '9', '1'}, {'b', 'd', 'b'},
			{'А', 'Д', 'А'}, {'О', 'П', 'О'}, {'а', 'у', 'а'}, {'世', '世', '世'}, {'你', '你', '你'}, {'好', '好', '好'},
		}},
		// overlapped
		{s: "z-a", wantFailed: true},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			rs, ok := runesRangeExpand(tt.s)
			if ok != !tt.wantFailed {
				t.Fatalf("RunesRangeExpand(%q) = %v, want %v", tt.s, ok, !tt.wantFailed)
			}
			if ok {
				assert.Equal(t, tt.want, rs)
			}
		})
	}
}
