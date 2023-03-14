package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestASCIISet_Contains(t *testing.T) {
	tests := []struct {
		set     string
		wantErr bool
		wantStr string
		wantAll string
		in      []byte
		notIn   []byte
	}{
		{
			set:     "aBzй",
			wantErr: true,
		},
		{
			set:     "aBz",
			wantStr: "Baz", wantAll: "Baz",
			in:    []byte("azB"),
			notIn: []byte("bZA"),
		},
		// range of symbols
		{
			set:     "acBbz",
			wantStr: "Ba-cz", wantAll: "Babcz",
			in:    []byte("abczB"),
			notIn: []byte("CZA"),
		},
		// duplicate symbol
		{
			set:     "aBzz",
			wantStr: "Baz", wantAll: "Baz",
			in:    []byte("azB"),
			notIn: []byte("bZA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.set, func(t *testing.T) {
			as, ok := MakeASCIISet(tt.set)
			if !ok != tt.wantErr {
				t.Errorf("MakeASCIISet(%q) = %v, want %v", tt.set, ok, tt.wantErr)
			} else if ok {
				if s := as.String(); s != tt.wantStr {
					t.Errorf("ASCIISet.String() = %q, want %q", s, tt.wantStr)
				}
				if s := as.StringAll(); s != tt.wantAll {
					t.Errorf("ASCIISet.StringAll() = %q, want %q", s, tt.wantAll)
				}
				for _, c := range tt.in {
					if !as.Contains(c) {
						t.Errorf("ASCIISet.Contains(%s) = false, want true", string(c))
					}
				}
				for _, c := range tt.notIn {
					if as.Contains(c) {
						t.Errorf("ASCIISet.Contains(%s) = true, want false", string(c))
					}
				}
			}
		})
	}
}

func TestASCIISet_Index(t *testing.T) {
	tests := []struct {
		set  string
		s    string
		want int
	}{
		{
			set:  "aBz",
			s:    "ABz",
			want: 1,
		},
		{
			set:  "aBz",
			s:    "aBz",
			want: 0,
		},
		{
			set:  "aBz",
			s:    "Abz",
			want: 2,
		},
		{
			set:  "aBz",
			s:    "AbZ",
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.set+"#"+tt.s, func(t *testing.T) {
			as, _ := MakeASCIISet(tt.set)

			index := as.Index(tt.s)
			assert.Equal(t, tt.want, index, "ASCIISet.Index(%s)", tt.s)

			index = as.IndexByte([]byte(tt.s))
			assert.Equal(t, tt.want, index, "ASCIISet.IndexByte(%s)", tt.s)
		})
	}
}

var (
	asciiSet      = "aBz"
	unicodeString = strings.Repeat("Abcd", 20) + "你好世, ЫВАЙz"
)

func BenchmarkContains_ASCIISet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		as, _ := MakeASCIISet(asciiSet)
		if !as.Contains('z') {
			b.Fatalf("ASCIISet(%q).Contains(%q)", asciiSet, 'z')
		}
	}
}

func BenchmarkContains_ASCIISet_Prealloc(b *testing.B) {
	as, _ := MakeASCIISet(asciiSet)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !as.Contains('z') {
			b.Fatalf("ASCIISet(%q).Contains(%q)", asciiSet, 'z')
		}
	}
}

func BenchmarkIndexASCII_ASCIISet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		as, _ := MakeASCIISet(asciiSet)
		if index := as.Index(unicodeString); index != len(unicodeString)-1 {
			b.Fatalf("ASCIISet(%q).Index(%q) = %d, want %d",
				asciiSet, unicodeString, index, len(unicodeString)-1,
			)
		}
	}
}

func BenchmarkIndexASCII_ASCIISet_Prealloc(b *testing.B) {
	as, _ := MakeASCIISet(asciiSet)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if index := as.Index(unicodeString); index != len(unicodeString)-1 {
			b.Fatalf("ASCIISet(%q).Index(%q) = %d, want %d",
				asciiSet, unicodeString, index, len(unicodeString)-1,
			)
		}
	}
}

func BenchmarkIndexASCII_StringsAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if index := strings.IndexAny(unicodeString, asciiSet); index != len(unicodeString)-1 {
			b.Fatalf("strings.IndexAny(%q, %q) = %d, want %d",
				asciiSet, unicodeString, index, len(unicodeString)-1,
			)
		}
	}
}
