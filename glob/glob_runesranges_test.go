package glob

import (
	"strings"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

func TestGlob_RunesRanges(t *testing.T) {
	tests := []testGlob{
		{
			glob: "[a-c]",
			want: &Glob{
				Glob: "[a-c]", Node: "[a-c]",
				MinLen: 1, MaxLen: 1,
				Items: []items.Item{items.NewRunesRanges("[a-c]")},
			},
			verify: "^[a-c]$",
			match:  []string{"a", "b", "c"},
			miss:   []string{"", "d", "ab", "ac", "a.b"},
		},
		{
			glob: "[a-c]z",
			want: &Glob{
				Glob: "[a-c]z", Node: "[a-c]z",
				MinLen: 2, MaxLen: 2, Suffix: "z",
				Items: []items.Item{items.NewRunesRanges("[a-c]")},
			},
			verify: "^[a-c]z$",
			match:  []string{"az", "bz", "cz"},
			miss:   []string{"", "d", "ab", "ac", "abz", "az.b", "a.bz"},
		},
		{
			glob: "[a-c]*",
			want: &Glob{
				Glob: "[a-c]*", Node: "[a-c]*",
				MinLen: 1, MaxLen: -1,
				Items: []items.Item{items.NewRunesRanges("[a-c]"), items.Star(0)},
			},
			verify: "^[a-c].*$",
			match:  []string{"a", "b", "c", "a.", "ab", "az", "bz", "cz", "a.b", "abz", "cab", "czb", "az.b", "a.bz"},
			miss:   []string{"", "d", "da"},
		},
		// composite
		{
			glob: "a*[b-e].b",
			want: &Glob{
				Glob: "a*[b-e].b", Node: "a*[b-e].b",
				MinLen: 4, MaxLen: -1, Prefix: "a", Suffix: ".b",
				Items: []items.Item{items.Star(0), items.NewRunesRanges("[b-e]")},
			},
			verify: `^a.*[b-e]\.b$`,
			match: []string{
				"ab.b", "ac.b", "ae.b",
				"aSTc.b",
				"acbec.b", "abbece.b",
			},
			miss: []string{"", "ab", "c", "a.b", "a.bd", "aa.b", "af.b"},
		},
	}
	for n, tt := range tests {
		runTestGlob(t, n, tt)
	}
}

func TestGlob_RunesRanges_Broken(t *testing.T) {
	tests := []testGlob{
		// broken
		// compare with graphite-clickhouse. Now It's not error, but filter
		// (Path LIKE 'z%' AND match(Path, '^z[ac$')))
		{glob: "[ac", wantErr: true},
		{glob: "a]c", wantErr: true},
		// skip empty
		{
			glob: "[]a",
			want: &Glob{
				Glob: "[]a", Node: "a",
				MinLen: 1, MaxLen: 1,
			},
			match: []string{"a"},
			miss:  []string{"", "b", "ab"},
		},
	}
	for n, tt := range tests {
		runTestGlob(t, n, tt)
	}
}

// benchmark for RunesRanges gready skip scan optimization (ASCII)
var (
	globStarRuneRangesAscii   = "find*[P-Z]*_ENDOF"
	stringStarRuneRangesAscii = strings.Repeat("findПОИСК", 20) + "LAST_ENDOF"
)

func Benchmark_Star_RunesRanges_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globStarRuneRangesAscii)
		if !g.Match(stringStarRuneRangesAscii) {
			b.Fatal(stringStarRuneRangesAscii)
		}
	}
}

func Benchmark_Star_RunesRanges_ASCII_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globStarRuneRangesAscii)
		if !w.MatchString(stringStarRuneRangesAscii) {
			b.Fatal(stringStarRuneRangesAscii)
		}
	}
}

func Benchmark_Star_RunesRanges_ASCII_Precompiled(b *testing.B) {
	g := ParseMust(globStarRuneRangesAscii)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringStarRuneRangesAscii) {
			b.Fail()
		}
	}
}

func Benchmark_Star_RunesRanges_ASCII_Regex_Precompiled(b *testing.B) {
	g := tests.BuildGlobRegexp(globStarRuneRangesAscii)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.MatchString(stringStarRuneRangesAscii) {
			b.Fatal(stringStarRuneRangesAscii)
		}
	}
}

// benchmark for RunesRanges gready skip scan optimization (Unicode)
var (
	globStarRuneRangesUnicode   = "find*[ЙЗЩА-ВЕ]*_ENDOF"
	stringStarRuneRangesUnicode = strings.Repeat("findПОИСК", 20) + "LAЕT_ENDOF"
)

func Benchmark_Star_RunesRanges_Unicode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globStarRuneRangesUnicode)
		if !g.Match(stringStarRuneRangesUnicode) {
			b.Fatal(stringStarRuneRangesUnicode)
		}
	}
}

func Benchmark_Star_RunesRanges_Unicode_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globStarRuneRangesUnicode)
		if !w.MatchString(stringStarRuneRangesUnicode) {
			b.Fatal(stringStarRuneRangesUnicode)
		}
	}
}

func Benchmark_Star_RunesRanges_Unicode_Precompiled(b *testing.B) {
	g := ParseMust(globStarRuneRangesUnicode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringStarRuneRangesUnicode) {
			b.Fail()
		}
	}
}

func Benchmark_Star_RunesRanges_Unicode_Regex_Precompiled(b *testing.B) {
	g := tests.BuildGlobRegexp(globStarRuneRangesUnicode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.MatchString(stringStarRuneRangesUnicode) {
			b.Fatal(stringStarRuneRangesUnicode)
		}
	}
}
