package glob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

func TestGlob_Rune(t *testing.T) {
	tests := []testGlob{
		{
			glob: "?Й?",
			want: &Glob{
				Glob: "?Й?", Node: "?Й?",
				MinLen: 4, MaxLen: 10,
				Items: []items.Item{items.Any(1), items.Rune('Й'), items.Any(1)},
			},
			verify: "^.Й.$",
			match:  []string{"aЙс", "ЙЙd", "ЯЙa"},
			miss:   []string{"", "a", "d", "aЙ", "Й", "ФiЙ", "aЙcc"},
		},
		{
			glob: "*Й?",
			want: &Glob{
				Glob: "*Й?", Node: "*Й?",
				MinLen: 3, MaxLen: -1,
				Items: []items.Item{items.Star(0), items.Rune('Й'), items.Any(1)},
			},
			verify: "^.*Й.$",
			match:  []string{"Йс", "ЙЙd", "aЙa"},
			miss:   []string{"", "a", "d", "aЙ", "Й"},
		},
	}
	for n, tt := range tests {
		runTestGlob(t, n, tt)
	}
}

// benchmark for rune gready skip scan optimization
var (
	globRune = "*Т*"
)

func BenchmarkStar_Rune_Miss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globRune)
		if !g.Match(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}

func _BenchmarkStar_Rune_Miss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globRune)
		if !w.MatchString(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}

func BenchmarkStar_Rune_Miss_Precompiled(b *testing.B) {
	g := ParseMust(globRune)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}

func _BenchmarkStar_Rune_Miss_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globRune)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !w.MatchString(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}
