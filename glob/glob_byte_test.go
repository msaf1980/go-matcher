package glob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

func TestGlob_Byte(t *testing.T) {
	tests := []testGlob{
		{
			glob: "*a?",
			want: &Glob{
				Glob: "*a?", Node: "*a?",
				MinLen: 2, MaxLen: -1,
				Items: []items.Item{items.Star(0), items.Byte('a'), items.Any(1)},
			},
			match: []string{"aa", "ab", "ac", "bad"},
			miss:  []string{"", "a", "d"},
		},
	}
	for n, tt := range tests {
		runTestGlob(t, n, tt)
	}
}

// benchmark for byte gready skip scan optimization
var (
	globByte = "*T*"
)

func Benchmark_Star_Byte_Miss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globByte)
		if !g.Match(stringStarEndAscii) {
			b.Fatal(stringStarEndAscii)
		}
	}
}

func _Benchmark_Star_Byte_Miss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globByte)
		if !w.MatchString(stringStarEndAscii) {
			b.Fatal(stringStarEndAscii)
		}
	}
}

func Benchmark_Star_Byte_Miss_Precompiled(b *testing.B) {
	g := ParseMust(globByte)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringStarEndAscii) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}

func _Benchmark_Star_Byte_Miss_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globByte)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !w.MatchString(stringStarEndAscii) {
			b.Fatal(stringStarEndAscii)
		}
	}
}
