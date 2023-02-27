package items

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRune(t *testing.T) {
	tests := []struct {
		s           string
		find        rune
		wantMatch   int
		wantFind    int
		wantFindLen int
		wantFindStr string
	}{
		{s: "", find: 'f', wantMatch: -1, wantFind: -1},
		{s: "f", find: 'f', wantMatch: 1, wantFind: 0, wantFindLen: 1},
		{s: "af", find: 'f', wantMatch: -1, wantFind: 1, wantFindLen: 1},
		{s: "ЯB界Cd", find: '界', wantMatch: -1, wantFind: 3, wantFindLen: 3, wantFindStr: "Cd"},
		{s: "ЯB界Cd", find: 'Я', wantMatch: 2, wantFind: 0, wantFindLen: 2, wantFindStr: "B界Cd"},
	}
	for _, tt := range tests {
		t.Run(tt.s+"#"+string(tt.find), func(t *testing.T) {
			item := Rune(tt.find)

			matched, flag := item.Match(tt.s)
			assert.Equal(t, FindDone, flag, "Rune(%q).Match(%q) flag", tt.find, tt.s)
			if matched != tt.wantMatch {
				t.Errorf("Rune(%q).Match(%q) = %v, want %v", tt.find, tt.s, matched, tt.wantMatch)
			}

			next, length, _ := item.Find(tt.s)
			assert.Equal(t, FindDone, flag, "Rune(%q).Find(%q) flag", tt.find, tt.s)
			if next != tt.wantFind {
				t.Errorf("Rune(%q).Find(%q) = %d, want %d", tt.find, tt.s, next, tt.wantFind)
			} else if next >= 0 {
				if length != tt.wantFindLen {
					t.Errorf("Rune(%q).Find(%q) length = %d, want %d", tt.find, tt.s, length, tt.wantFindLen)
				}
				if nextS := tt.s[next+length:]; nextS != tt.wantFindStr {
					t.Errorf("Rune(%q).Find(%q) = %q, want %q", tt.find, tt.s, nextS, tt.wantFindStr)
				}
			}
		})
	}
}

var (
	findRuneASCII   = 'e'
	findRuneUnicode = '界'
)

func BenchmarkRune_Find_Unicode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var item Item = Rune(findRuneUnicode)
		_, _, _ = item.Find(stringUnicode)
	}
}

func BenchmarkRune_Find_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var item Item = Rune(findRuneASCII)
		_, _, _ = item.Find(stringASCII)
	}
}
