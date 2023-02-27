package items

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	tests := []struct {
		s           string
		n           int
		wantMatch   int
		wantFind    int
		wantFindLen int
		wantFindStr string
	}{
		{s: "", n: 1, wantMatch: -1, wantFind: -1},
		{s: "f", n: 1, wantMatch: 1, wantFind: 0, wantFindLen: 1},
		{s: "f", n: 2, wantMatch: -1, wantFind: -1},
		{s: "af", n: 2, wantMatch: 2, wantFind: 0, wantFindLen: 2},
		{s: "Яz界Cd", n: 3, wantMatch: 6, wantFind: 0, wantFindLen: 6, wantFindStr: "Cd"},
	}
	for _, tt := range tests {
		t.Run(tt.s+"#"+strconv.Itoa(tt.n), func(t *testing.T) {
			item := Any(tt.n)

			offset, flag := item.Match(tt.s)
			assert.Equal(t, FindDone, flag, "Any(%d).Match(%q) flag", tt.n, tt.s)
			if offset != tt.wantMatch {
				t.Errorf("Any(%d).Match(%q) offset = %v, want %v", tt.n, tt.s, offset, tt.wantMatch)
			}

			next, length, flag := item.Find(tt.s)
			assert.Equal(t, FindForwarded, flag, "Any(%d).Find(%q) flag", tt.n, tt.s)
			if next != tt.wantFind {
				t.Errorf("Any(%d).Find(%q) index = %d, want %d", tt.n, tt.s, next, tt.wantFind)
			} else if next >= 0 {
				if length != tt.wantFindLen {
					t.Errorf("Any(%d).Find(%q) length = %d, want %d", tt.n, tt.s, length, tt.wantFindLen)
				}
				if nextS := tt.s[next+length:]; nextS != tt.wantFindStr {
					t.Errorf("Any(%d).Find(%q) next = %q, want %q", tt.n, tt.s, nextS, tt.wantFindStr)
				}
			}
		})
	}
}
