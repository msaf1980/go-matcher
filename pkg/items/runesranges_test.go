package items

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunesRanges(t *testing.T) {
	tests := []struct {
		s           string
		ranges      string
		wantMatch   int
		wantFind    int
		wantFindLen int
		wantFindStr string
	}{
		{s: "", ranges: "[f]", wantMatch: -1, wantFind: -1},
		{s: "f", ranges: "[f]", wantMatch: 1, wantFind: 0, wantFindLen: 1},
		{s: "af", ranges: "[f]", wantMatch: -1, wantFind: 1, wantFindLen: 1},
		{s: "ЯB界Cd", ranges: "[界]", wantMatch: -1, wantFind: 3, wantFindLen: 3, wantFindStr: "Cd"},
		{s: "ЯB界Cd", ranges: "[ЯВ]", wantMatch: 2, wantFind: 0, wantFindLen: 2, wantFindStr: "B界Cd"},
		{
			s: "ЗКйлМН", ranges: "[а-йв-у1-9b-dА-ДО-П你好世]",
			wantMatch: -1, wantFind: 4, wantFindLen: 2, wantFindStr: "лМН",
		},
	}
	for _, tt := range tests {
		t.Run(tt.s+"#"+string(tt.ranges), func(t *testing.T) {
			item := NewRunesRanges(tt.ranges)
			if item == nil {
				t.Fail()
			}

			matched, flag := item.Match(tt.s)
			assert.Equal(t, FindDone, flag, "RunesRanges(%q).Match(%q) flag", tt.ranges, tt.s)
			if matched != tt.wantMatch {
				t.Errorf("RunesRanges(%q).Match(%q) = %v, want %v", tt.ranges, tt.s, matched, tt.wantMatch)
			}

			next, length, flag := item.Find(tt.s)
			assert.Equal(t, FindDone, flag, "RunesRanges(%q).Find(%q) flag", tt.ranges, tt.s)
			if next != tt.wantFind {
				t.Errorf("RunesRanges(%q).Find(%q) = %d, want %d", tt.ranges, tt.s, next, tt.wantFind)
			} else if next >= 0 {
				if length != tt.wantFindLen {
					t.Errorf("RunesRanges(%q).Find(%q) length = %q, want %q", tt.ranges, tt.s, length, tt.wantFindLen)
				}
				if nextS := tt.s[next+length:]; nextS != tt.wantFindStr {
					t.Errorf("RunesRanges(%q).Find(%q) = %q, want %q", tt.ranges, tt.s, nextS, tt.wantFindStr)
				}
			}
		})
	}
}

var (
	findRunesRangesASCII   = "{ac,"
	findRunesRangesUnicode = "[а-йв-у1-9b-dА-ДО-П你好世]"
)

func BenchmarkRunesRanges_Find_Unicode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var item Item = NewRunesRanges(findRunesRangesUnicode)
		_, _, _ = item.Find(stringUnicode)
	}
}

func BenchmarkRunesRanges_Find_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var item Item = NewRunesRanges(findRunesRangesASCII)
		_, _, _ = item.Find(stringASCII)
	}
}
