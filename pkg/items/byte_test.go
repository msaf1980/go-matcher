package items

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByte(t *testing.T) {
	tests := []struct {
		s           string
		find        byte
		wantMatch   int
		wantFind    int
		wantFindLen int
		wantFindStr string
	}{
		{s: "", find: 'f', wantMatch: -1, wantFind: -1},
		{s: "f", find: 'f', wantMatch: 1, wantFind: 0, wantFindLen: 1},
		{s: "af", find: 'f', wantMatch: -1, wantFind: 1, wantFindLen: 1},
		{s: "ЯB界Cd", find: 'C', wantMatch: -1, wantFind: 6, wantFindLen: 1, wantFindStr: "d"},
	}
	for _, tt := range tests {
		t.Run(tt.s+"#"+string(tt.find), func(t *testing.T) {
			item := Byte(tt.find)

			offset, flag := item.Match(tt.s)
			assert.Equal(t, FindDone, flag, "Byte(%q).Match(%q) flag", tt.find, tt.s)
			if offset != tt.wantMatch {
				t.Errorf("Byte(%q).Match(%q) offset = %v, want %v", tt.find, tt.s, offset, tt.wantMatch)
			}

			next, length, flag := item.Find(tt.s)
			assert.Equal(t, FindDone, flag, "Byte(%q).Find(%q) flag", tt.find, tt.s)
			if next != tt.wantFind {
				t.Errorf("Byte(%q).Find(%q) index = %d, want %d", tt.find, tt.s, next, tt.wantFind)
			} else if next >= 0 {
				if length != tt.wantFindLen {
					t.Errorf("Byte(%q).Find(%q) length = %d, want %d", tt.find, tt.s, length, tt.wantFindLen)
				}
				if nextS := tt.s[next+length:]; nextS != tt.wantFindStr {
					t.Errorf("Byte(%q).Find(%q) next = %q, want %q", tt.find, tt.s, nextS, tt.wantFindStr)
				}
			}
		})
	}
}

var (
	findByte = 'e'
)

func BenchmarkByte_Find_Unicode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var item Item = Byte(findByte)
		_, _, _ = item.Find(stringUnicode)
	}
}

func BenchmarkByte_Find_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var item Item = Byte(findByte)
		_, _, _ = item.Find(stringASCII)
	}
}
