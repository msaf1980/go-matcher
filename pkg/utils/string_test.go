package utils

import "testing"

func TestStringSkipRunes(t *testing.T) {
	tests := []struct {
		s        string
		runes    int
		wantNext int
		want     string
	}{
		{s: "", runes: 1, wantNext: -1},
		{s: "ac", runes: 2, wantNext: 2}, // end of
		{s: "", runes: 0, wantNext: 0},
		{s: "aBcd", runes: 2, wantNext: 2, want: "cd"},
		{s: "ЯB界Cd", runes: 3, wantNext: 6, want: "Cd"},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if next := StringSkipRunes(tt.s, tt.runes); next != tt.wantNext {
				t.Errorf("StringSkipRunes() = %v, want %v", next, tt.wantNext)
			} else if next >= 0 {
				if nextS := tt.s[next:]; nextS != tt.want {
					t.Errorf("StringSkipRunes() = %q, want %q", nextS, tt.want)
				}
			}
		})
	}
}

func Benchmark_String_SkipRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = StringSkipRunes("ЯB界Cd", 3)
	}
}

func Benchmark_String_SkipRunesEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = StringSkipRunes("ЯB界Cd", 0)
	}
}
