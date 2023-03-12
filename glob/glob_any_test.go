package glob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

func TestGlob_Any(t *testing.T) {
	tests := []testGlob{
		{
			glob: "?",
			want: &Glob{
				Glob: "?", Node: "?", MinLen: 1, MaxLen: 4,
				Items: []items.Item{items.Any(1)},
			},
			verify: "^.$",
			match:  []string{"a", "c"},
			miss:   []string{"", "ab", "a.b"},
		},
		{
			glob: "??",
			want: &Glob{
				Glob: "??", Node: "??", MinLen: 2, MaxLen: 8,
				Items: []items.Item{items.Any(2)},
			},
			verify: "^.{2}$",
			match:  []string{"ab", "cc"},
			miss:   []string{"", "a", "a.b"},
		},
		{
			glob: "???",
			want: &Glob{
				Glob: "???", Node: "???",
				MinLen: 3, MaxLen: 12,
				Items: []items.Item{items.Any(3)},
			},
			verify: "^.{3}$",
			match:  []string{"abc", "ccc"},
			miss:   []string{"", "ab", "abcd", "a.bd"},
		},
		{
			glob: "a?",
			want: &Glob{
				Glob: "a?", Node: "a?",
				MinLen: 2, MaxLen: 5, Prefix: "a",
				Items: []items.Item{items.Any(1)},
			},
			verify: "^a.$",
			match:  []string{"ac", "az"},
			miss:   []string{"", "a", "b", "bc", "ace", "a.c"},
		},
		{
			glob: "a?c",
			want: &Glob{
				Glob: "a?c", Node: "a?c",
				MinLen: 3, MaxLen: 6, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Any(1)},
			},
			verify: "^a.c$",
			match:  []string{"acc", "aec"},
			miss:   []string{"", "ab", "ac", "ace", "a.e", "bec"},
		},
		{
			glob: "a?[c]?d",
			want: &Glob{
				Glob: "a?[c]?d", Node: "a?c?d",
				MinLen: 5, MaxLen: 11, Prefix: "a", Suffix: "d",
				Items: []items.Item{items.Any(1), items.Byte('c'), items.Any(1)},
			},
			verify: "^a.c.d$",
			match:  []string{"aZccd", "aZcAd"},
			miss:   []string{"", "ab", "ac", "ace", "aZDAd", "a.c"},
		},
		{
			glob: "a*?c?d",
			want: &Glob{
				Glob: "a*?c?d", Node: "a*?c?d",
				MinLen: 5, MaxLen: -1, Prefix: "a", Suffix: "d",
				Items: []items.Item{items.Star(1), items.Byte('c'), items.Any(1)},
			},
			verify: "^a.*.c.d$",
			match:  []string{"aZccd", "aZcAd", "aIZcAd"},
			miss:   []string{"", "ab", "ac", "ace", "aZDAd", "a.c"},
		},
	}
	for n, tt := range tests {
		runTestGlob(t, n, tt)
	}
}

var (
	globAnySkipASCII   = "DB*?web*_Status"
	stringAnySkipASCII = "DBCassandraSalesSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIP_we_Status"
)

func Benchmark_Star_Any_Miss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globAnySkipASCII)
		if g.Match(stringAnySkipASCII) {
			b.Fatal(stringAnySkipASCII)
		}
	}
}

func _Benchmark_Star_Any_Miss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globAnySkipASCII)
		if w.MatchString(stringAnySkipASCII) {
			b.Fatal(stringAnySkipASCII)
		}
	}
}

func Benchmark_Star_Any_Miss_Precompiled(b *testing.B) {
	g := ParseMust(globAnySkipASCII)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.Match(stringAnySkipASCII) {
			b.Fatal(stringAnySkipASCII)
		}
	}
}

func _Benchmark_Star_Any_Miss_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globAnySkipASCII)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(stringAnySkipASCII) {
			b.Fatal(stringAnySkipASCII)
		}
	}
}
