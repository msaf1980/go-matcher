package items

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		s           string
		find        string
		wantMatch   int
		wantFind    int
		wantFindLen int
		wantFindStr string
	}{
		{s: "", find: ""},
		{s: "", find: "f", wantMatch: -1, wantFind: -1},
		{s: "f", find: "f", wantMatch: 1, wantFind: 0, wantFindLen: 1},
		{s: "af", find: "f", wantMatch: -1, wantFind: 1, wantFindLen: 1},
		{s: "ЯB界Cd", find: "界C", wantMatch: -1, wantFind: 3, wantFindLen: 4, wantFindStr: "d"},
	}
	for _, tt := range tests {
		t.Run(tt.s+"#"+tt.find, func(t *testing.T) {
			item := NewString(tt.find)

			matched, flag := item.Match(tt.s)
			assert.Equal(t, FindDone, flag, "String(%q).Match(%q) flag", tt.find, tt.s)
			if matched != tt.wantMatch {
				t.Errorf("String(%q).Match(%q) = %v, want %v", tt.find, tt.s, matched, tt.wantMatch)
			}

			pos, length, flag := item.Find(tt.s)
			assert.Equal(t, FindDone, flag, "String(%q).Find(%q) flag", tt.find, tt.s)
			if pos != tt.wantFind {
				t.Errorf("String(%q).Find(%q) pos = %d, want %d", tt.find, tt.s, pos, tt.wantFind)
			} else if pos >= 0 {
				if length != tt.wantFindLen {
					t.Errorf("String(%q).Find(%q) length = %d, want %d", tt.find, tt.s, length, tt.wantFindLen)
				}
				if nextS := tt.s[pos+length:]; nextS != tt.wantFindStr {
					t.Errorf("String(%q).Find(%q) next = %q, want %q", tt.find, tt.s, nextS, tt.wantFindStr)
				}
			}
		})
	}
}

var (
	findASCII     = "ABCDe"
	stringASCII   = strings.Repeat("IKDEl", 60) + findASCII
	findUnicode   = "ЯB界Cd"
	stringUnicode = strings.Repeat("你好世", 60) + findUnicode
)

func BenchmarkString_Find_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var item Item = NewString(findASCII)
		_, _, _ = item.Find(stringASCII)
	}
}

func BenchmarkString_Find_Unicode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var item Item = NewString(findUnicode)
		_, _, _ = item.Find(stringUnicode)
	}
}
