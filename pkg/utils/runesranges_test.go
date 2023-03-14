package utils

import (
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
)

var (
	cmpTransformRunesRanges = cmp.Transformer("ASCII", func(in ASCIISet) string {
		return in.String()
	})
)

func TestRunesRanges_Expand(t *testing.T) {
	tests := []struct {
		s          string
		want       RunesRanges
		wantStr    string
		wantFailed bool
		in         []rune
		notIn      []rune
	}{
		{
			s:       "[-q]",
			want:    RunesRanges{ASCII: MakeASCIISetMust("q"), MinSize: 1, MaxSize: 1},
			wantStr: "[q]",
			in:      []rune{'q'},
			notIn:   []rune{'Q', 'b', 'Я'},
		},
		{
			s:       "[a-c]",
			want:    RunesRanges{ASCII: MakeASCIISetMust("abc"), MinSize: 1, MaxSize: 1},
			wantStr: "[a-c]",
			in:      []rune{'a', 'b', 'c'},
			notIn:   []rune{'A', 'e', 'Я'},
		},
		{
			s:       "[a-cf]",
			want:    RunesRanges{ASCII: MakeASCIISetMust("abcf"), MinSize: 1, MaxSize: 1},
			wantStr: "[a-cf]",
			in:      []rune{'a', 'b', 'c', 'f'},
			notIn:   []rune{'A', 'd', 'e', 'j', 'Я'},
		},
		{
			s:       "[za-c]",
			wantStr: "[a-cz]",
			want:    RunesRanges{ASCII: MakeASCIISetMust("abcz"), MinSize: 1, MaxSize: 1},
			in:      []rune{'a', 'b', 'c', 'z'},
			notIn:   []rune{'A', 'e', 'j', 'Я'},
		},
		{
			s:       "[a-czA-C]",
			wantStr: "[A-Ca-cz]",
			want:    RunesRanges{ASCII: MakeASCIISetMust("abcABCz"), MinSize: 1, MaxSize: 1},
			in:      []rune{'a', 'b', 'c', 'z', 'A', 'B', 'C'},
			notIn:   []rune{'E', 'e', 'j', 'Я'},
		},
		{
			s:       "[a-czqj-mA-C]",
			wantStr: "[A-Ca-cj-mqz]",
			want:    RunesRanges{ASCII: MakeASCIISetMust("abcABCqzjklm"), MinSize: 1, MaxSize: 1},
			in:      []rune{'a', 'b', 'c', 'q', 'j', 'l', 'm', 'z', 'A', 'B', 'C'},
			notIn:   []rune{'E', 'e', 'i', 'n', 'Я'},
		},
		// contains unicode
		{
			s: "[界Э-ЯW-Z]",
			want: RunesRanges{
				ASCII:         MakeASCIISetMust("WXYZ"),
				UnicodeRanges: []RuneRange{{1069, 1071}, {30028, 30028}},
				MinSize:       1, MaxSize: 3,
			},
			wantStr: "[W-ZЭ-Я界]",
			in:      []rune{'W', 'Y', 'Z', 'Э', 'Ю', 'Я', '界'},
			notIn:   []rune{'E', 'e', 'i', 'n', 'П', '好'},
		},
		// partially broken
		{
			s: "[-]", want: RunesRanges{}, wantStr: "[]", notIn: []rune{'E', 'e', 'i', 'n', 'П', '好'},
		},
		{
			s:       "[a-]",
			want:    RunesRanges{ASCII: MakeASCIISetMust("a"), MinSize: 1, MaxSize: 1},
			wantStr: "[a]",
		},
		{
			s:       "[-a-c]",
			want:    RunesRanges{ASCII: MakeASCIISetMust("abc"), MinSize: 1, MaxSize: 1},
			wantStr: "[a-c]",
		},
		// duplicated ASCII ranges
		{
			s:       "[a-cb-ld-fa-c]",
			want:    RunesRanges{ASCII: MakeASCIISetMust("abcdefghijkl"), MinSize: 1, MaxSize: 1},
			wantStr: "[a-l]",
		},
		{
			s:       "[1-32-45-76-78-95789]",
			want:    RunesRanges{ASCII: MakeASCIISetMust("123456789"), MinSize: 1, MaxSize: 1},
			wantStr: "[1-9]",
			in:      []rune{'1', '2', '3', '4', '5', '6', '9'},
			notIn:   []rune{'E', 'e', 'i', 'n', 'П'},
		},
		// duplicated Unicode ranges
		{
			s: "[а-йв-у1-9b-dА-ДО-П你好世]",
			want: RunesRanges{
				ASCII: MakeASCIISetMust("123456789bcd"),
				UnicodeRanges: []RuneRange{
					{'А', 'Д'}, {'О', 'П'}, {'а', 'у'}, {'世', '世'}, {'你', '你'}, {'好', '好'},
				},
				MinSize: 1, MaxSize: 3,
			},
			wantStr: "[1-9b-dА-ДО-Па-у世你好]",
			in:      []rune{'1', '2', '3', '4', '5', '6', '9', 'А', 'Б', 'П', 'с', '你'},
			notIn:   []rune{'E', 'e', 'i', 'n', 'Ы'},
		},
		// broken
		{s: "", wantFailed: true},
		{s: "[a", wantFailed: true},
		{s: "a]", wantFailed: true},
		// overlapped
		{s: "[z-a]", want: RunesRanges{}, wantStr: "[]", notIn: []rune{'E', 'e', 'i', 'n', 'Ы'}},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			rs, ok := RunesRangeExpand(tt.s)
			if ok != !tt.wantFailed {
				t.Fatalf("RunesRangeExpand(%q) = %v, want %v", tt.s, ok, !tt.wantFailed)
			}
			if ok {
				if !reflect.DeepEqual(rs, tt.want) {
					t.Errorf("RunesRangeExpand(%q) = %s", tt.s, cmp.Diff(tt.want, rs, cmpTransformRunesRanges))
				}
				if s := rs.String(); s != tt.wantStr {
					t.Errorf("RunesRangeExpand(%q) = %q, want %q", tt.s, s, tt.wantStr)
				}

				for _, c := range tt.in {
					if !rs.Contains(c) {
						t.Errorf("RunesRangeExpand(%q).Contains(%q) = false, want true", tt.s, string(c))
					}
				}
				for _, r := range tt.notIn {
					if rs.Contains(r) {
						t.Errorf("RunesRangeExpand(%q).Contains(%q) = true, want false", tt.s, string(r))
					}
				}
			}
		})
	}
}

func TestRunesRanges_Contains(t *testing.T) {
	tests := []struct {
		ranges   string
		contains []rune
		not      []rune
	}{
		{
			ranges:   "aBz你й",
			contains: []rune("azBй你"),
			not:      []rune("bZAЫ"),
		},
		{
			ranges:   "aBzz", // duplicate symbol
			contains: []rune("azB"),
			not:      []rune("bZA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.ranges, func(t *testing.T) {
			var rs RunesRanges
			rs.Adds(tt.ranges)
			rs.Merge()

			for _, c := range tt.contains {
				if !rs.Contains(c) {
					t.Errorf("RuneSet.Contains(%s) = false, want true", string(c))
				}
			}
			for _, c := range tt.not {
				if rs.Contains(c) {
					t.Errorf("RuneSet.Contains(%s) = true, want false", string(c))
				}
			}
		})
	}
}

func TestRunesRanges_Index(t *testing.T) {
	tests := []struct {
		s      string
		ranges string

		wantIndex  int
		wantIndexC rune
		wantIndexN int

		wantStartC rune
		wantStartN int
	}{
		{
			s: "aBz你й", ranges: "[你й]",
			wantIndex: 3, wantIndexC: '你', wantIndexN: 3,
			wantStartC: utf8.RuneError, wantStartN: -1,
		},
		{
			s: "aBz你й", ranges: "[a你й]",
			wantIndex: 0, wantIndexC: 'a', wantIndexN: 1,
			wantStartC: 'a', wantStartN: 1,
		},
		{
			s: "你йaBz", ranges: "[a你й]",
			wantIndex: 0, wantIndexC: '你', wantIndexN: 3,
			wantStartC: '你', wantStartN: 3,
		},
		{
			s: "й你aBz", ranges: "[a你й]",
			wantIndex: 0, wantIndexC: 'й', wantIndexN: 2,
			wantStartC: 'й', wantStartN: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.s+"#"+tt.ranges, func(t *testing.T) {
			rs, ok := RunesRangeExpand(tt.ranges)
			if !ok {
				t.Fail()
			}

			c, n := rs.StartsWith(tt.s)
			if c != tt.wantStartC {
				t.Errorf("RunesRanges.StartsWith() got C = %v, want %v", c, tt.wantStartC)
			}
			if n != tt.wantStartN {
				t.Errorf("RunesRanges.StartsWith() got N = %v, want %v", n, tt.wantStartN)
			}

			pos, c, n := rs.Index(tt.s)
			if pos != tt.wantIndex {
				t.Errorf("RunesRanges.Index() got Pos = %v, want %v", pos, tt.wantIndex)
			}
			if c != tt.wantIndexC {
				t.Errorf("RunesRanges.Index() got C = %v, want %v", c, tt.wantIndexC)
			}
			if n != tt.wantIndexN {
				t.Errorf("RunesRanges.Index() got Index = %v, want %v", n, tt.wantIndexN)
			}
		})
	}
}

func BenchmarkIndexUnicode_ASCII_RuneSet_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var rs RunesRanges
		rs.Adds(asciiSet)
		if index, _, _ := rs.Index(unicodeString); index != len(unicodeString)-1 {
			b.Fatalf("RuneSet(%q).Index(%q) = %d, want %d",
				asciiSet, unicodeString, index, len(unicodeString)-1,
			)
		}
	}
}

func BenchmarkIndexUnicode_ASCII_RuneSet_Prealloc(b *testing.B) {
	var rs RunesRanges
	rs.Adds(asciiSet)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if index, _, _ := rs.Index(unicodeString); index != len(unicodeString)-1 {
			b.Fatalf("RuneSet(%q).Index(%q) = %d, want %d",
				asciiSet, unicodeString, index, len(unicodeString)-1,
			)
		}
	}
}

func BenchmarkIndexUnicode_ASCII_StringsAny(b *testing.B) {
	want := len(unicodeString) - 1
	for i := 0; i < b.N; i++ {
		if index := strings.IndexAny(unicodeString, asciiSet); index != want {
			b.Fatalf("strings.IndexAny(%q, %q) = %d, want %d",
				unicodeString, runeSet, index, want,
			)
		}
	}
}

var (
	runeSet = "界Яz"
)

func BenchmarkIndexUnicode_Unicode_RuneSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var rs RunesRanges
		rs.Adds(runeSet)
		if index, _, _ := rs.Index(unicodeString); index != len(unicodeString)-1 {
			b.Fatalf("RuneSet(%q).Index(%q) = %d, want %d",
				runeSet, unicodeString, index, len(unicodeString)-1,
			)
		}
	}
}

func BenchmarkIndexUnicode_Unicode_RuneSet_Prealloc(b *testing.B) {
	var rs RunesRanges
	rs.Adds(runeSet)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if index, _, _ := rs.Index(unicodeString); index != len(unicodeString)-1 {
			b.Fatalf("RuneSet(%q).Index(%q) = %d, want %d",
				asciiSet, unicodeString, index, len(unicodeString)-1,
			)
		}
	}
}

func BenchmarkIndexUnicode_Unicode_StringsAny(b *testing.B) {
	want := len(unicodeString) - 1
	for i := 0; i < b.N; i++ {
		if index := strings.IndexAny(unicodeString, runeSet); index != want {
			b.Fatalf("strings.IndexAny(%q, %q) = %d, want %d",
				unicodeString, runeSet, index, want,
			)
		}
	}
}

var (
	// largeRunesRange = "[界Яq-zefgikmoqsuQ-ZEFGIKMOQSU]"
	largeRuneSet = "界ЯqrstuvwxyzefgikmoqsuQRSTUVWXYZEFGIKMOQSU"
)

func BenchmarkIndexUnicode_Large_RuneSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var rs RunesRanges
		rs.Adds(largeRuneSet)
		if index, _, _ := rs.Index(unicodeString); index != len(unicodeString)-1 {
			b.Fatalf("RuneSet(%q).Index(%q) = %d, want %d",
				runeSet, unicodeString, index, len(unicodeString)-1,
			)
		}
	}
}

func BenchmarkIndexUnicode_Large_RuneSet_Prealloc(b *testing.B) {
	var rs RunesRanges
	rs.Adds(largeRuneSet)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if index, _, _ := rs.Index(unicodeString); index != len(unicodeString)-1 {
			b.Fatalf("RuneSet(%q).Index(%q) = %d, want %d",
				asciiSet, unicodeString, index, len(unicodeString)-1,
			)
		}
	}
}

func BenchmarkIndexUnicode_StringsAny_Large(b *testing.B) {
	want := len(unicodeString) - 1
	for i := 0; i < b.N; i++ {
		if index := strings.IndexAny(unicodeString, largeRuneSet); index != want {
			b.Fatalf("strings.IndexAny(%q, %q) = %d, want %d",
				unicodeString, runeSet, index, want,
			)
		}
	}
}
