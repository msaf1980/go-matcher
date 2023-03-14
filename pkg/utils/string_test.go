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

func TestStringSkipRunesLast(t *testing.T) {
	tests := []struct {
		s       string
		runes   int
		wantEnd int
		want    string
	}{
		{s: "", runes: 1, wantEnd: -1},
		{s: "ac", runes: 3, wantEnd: -1}, // end of
		{s: "ac", runes: 2, wantEnd: 0},  // end of
		{s: "", runes: 0, wantEnd: 0},
		{s: "aBcd", runes: 2, wantEnd: 2, want: "aB"},
		{s: "ЯB界Cd", runes: 2, wantEnd: 6, want: "ЯB界"},
		{s: "ЯB界Cd", runes: 3, wantEnd: 3, want: "ЯB"},
		{s: "ЯB界Cd", runes: 4, wantEnd: 2, want: "Я"},
		{s: "ЯB界Cd", runes: 5, wantEnd: 0},
		{s: "ЯB界Cd", runes: 6, wantEnd: -1},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			next := StringSkipRunesLast(tt.s, tt.runes)
			if next != tt.wantEnd {
				t.Errorf("StringSkipRunes() = %v, want %v", next, tt.wantEnd)
			}
			if next >= 0 {
				if nextS := tt.s[:next]; nextS != tt.want {
					t.Errorf("StringSkipRunes() string = %q, want %q", nextS, tt.want)
				}
			}
		})
	}
}

func BenchmarkString_SkipRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = StringSkipRunes("ЯB界Cd", 3)
	}
}

func BenchmarkString_SkipRunesEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = StringSkipRunes("ЯB界Cd", 0)
	}
}
