package glob

import (
	"strings"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

func TestGlob_String(t *testing.T) {
	tests := []testGlob{
		{want: &Glob{}}, // empty
		// string match
		{
			glob:  "a",
			want:  &Glob{Glob: "a", Node: "a", MinLen: 1, MaxLen: 1},
			match: []string{"a"},
			miss:  []string{"", "b", "ab", "ba", "a.", "a.b"},
		},
		{
			glob:   "a.bc",
			want:   &Glob{Glob: "a.bc", Node: "a.bc", MinLen: 4, MaxLen: 4},
			verify: `^a\.bc$`,
			match:  []string{"a.bc"},
			miss:   []string{"", "b", "ab", "bc", "abc", "a.b", "a.bc.", "b.bc", "a.bce", "a.bc.e"},
		},
		// one item optimization
		{
			glob: "[a-]",
			want: &Glob{
				Glob: "[a-]", Node: "a",
				MinLen: 1, MaxLen: 1,
			},
			match: []string{"a"},
			miss:  []string{"", "b", "d", "ab", "a.b"},
		},
		{
			glob:  "[Й-]",
			want:  &Glob{Glob: "[Й-]", Node: "Й", MinLen: 2, MaxLen: 2},
			match: []string{"Й"},
			miss:  []string{"", "ф", "d", "ab", "a.b", "Ц", "ЙЦ"},
		},
		// one item optimization
		{
			glob:  "{a}",
			want:  &Glob{Glob: "{a}", Node: "a", MinLen: 1, MaxLen: 1},
			match: []string{"a"},
			miss:  []string{"", "b", "d", "ab", "a.b"},
		},
		{
			glob:  "{Й}",
			want:  &Glob{Glob: "{Й}", Node: "Й", MinLen: 2, MaxLen: 2},
			match: []string{"Й"},
			miss:  []string{"", "ф", "d", "ab", "a.b", "Ц", "ЙЦ"},
		},
		{
			glob:  "{str}",
			want:  &Glob{Glob: "{str}", Node: "str", MinLen: 3, MaxLen: 3},
			match: []string{"str"},
			miss:  []string{"", "st", "strt"},
		},
		// merge strings
		{
			glob: "a[a-]Z",
			want: &Glob{
				Glob: "a[a-]Z", Node: "aaZ",
				MinLen: 3, MaxLen: 3,
			},
			verify: `^aaZ$`,
			match:  []string{"aaZ"},
			miss:   []string{"", "a", "b", "d", "aa", "ab", "aaz", "aaZ.", "aaZa", "a.b"},
		},
		{
			glob: "a[a-]Z[Q]",
			want: &Glob{
				Glob: "a[a-]Z[Q]", Node: "aaZQ", MinLen: 4, MaxLen: 4,
			},
			match: []string{"aaZQ"},
			miss:  []string{"", "ab", "aaZQa"},
		},
		// merge strings to prefix
		{
			glob: "a[a-]Z[Q]*",
			want: &Glob{
				Glob: "a[a-]Z[Q]*", Node: "aaZQ*", MinLen: 4, MaxLen: -1,
				Prefix: "aaZQ", Items: []items.NodeItem{{Node: "*", Item: items.Star(0)}},
			},
			verify: `^aaZQ.*$`,
			match:  []string{"aaZQ", "aaZQa"},
			miss:   []string{"", "ab", "aaZqa"},
		},
		// merge string to prefix and suffix
		{
			glob: "a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l",
			want: &Glob{
				Glob:   "a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l",
				Node:   "aaZQstLT*INN*zaSTltl",
				MinLen: 18, MaxLen: -1, Prefix: "aaZQstLT", Suffix: "zaSTltl",
				Items: []items.NodeItem{
					{Node: "*", Item: items.Star(0)},
					{Node: "INN", Item: items.NewString("INN")},
					{Node: "*", Item: items.Star(0)},
				},
			},
			verify: `^aaZQstLT.*INN.*zaSTltl$`,
			match:  []string{"aaZQstLTINNzaSTltl", "aaZQstLTaINNzaSTltl"},
			miss:   []string{"", "ab", "aaZqa"},
		},
		{
			glob: "a[a-]Z[Q]st{LT}*Ц{NN}*[z-][a]ST{lt}l",
			want: &Glob{
				Glob:   "a[a-]Z[Q]st{LT}*Ц{NN}*[z-][a]ST{lt}l",
				Node:   "aaZQstLT*ЦNN*zaSTltl",
				MinLen: 19, MaxLen: -1, Prefix: "aaZQstLT", Suffix: "zaSTltl",
				Items: []items.NodeItem{
					{Node: "*", Item: items.Star(0)},
					{Node: "ЦNN", Item: items.NewString("ЦNN")},
					{Node: "*", Item: items.Star(0)},
				},
			},
			verify: `^aaZQstLT.*ЦNN.*zaSTltl$`,
			match:  []string{"aaZQstLTЦNNzaSTltl", "aaZQstLTaЦNNzaSTltl"},
			miss:   []string{"", "ab", "aaZqa", "aaZQstLTaЦaNNzaSTltl", "aaZQstLTaaNNzaSTltl"},
		},
		{
			glob: "a[Л]*И{Ф}*b{В}*В{NN}*nn[в]*В{N}*{Ц}l",
			want: &Glob{
				Glob:   "a[Л]*И{Ф}*b{В}*В{NN}*nn[в]*В{N}*{Ц}l",
				Node:   "aЛ*ИФ*bВ*ВNN*nnв*ВN*Цl",
				MinLen: 24, MaxLen: -1, Prefix: "aЛ", Suffix: "Цl",
				Items: []items.NodeItem{
					{Node: "*", Item: items.Star(0)},
					{Node: "ИФ", Item: items.NewString("ИФ")},
					{Node: "*", Item: items.Star(0)},
					{Node: "bВ", Item: items.NewString("bВ")},
					{Node: "*", Item: items.Star(0)},
					{Node: "ВNN", Item: items.NewString("ВNN")},
					{Node: "*", Item: items.Star(0)},
					{Node: "nnв", Item: items.NewString("nnв")},
					{Node: "*", Item: items.Star(0)},
					{Node: "ВN", Item: items.NewString("ВN")},
					{Node: "*", Item: items.Star(0)},
				},
			},
			verify: `^aЛ.*ИФ.*bВ.*ВNN.*nnв.*ВN.*Цl$`,
			match:  []string{"aЛИФbВВNNnnвВNЦl", "aЛЦИФrbВ_ВNN_nnвВN_Цl"},
			miss:   []string{"", "ab", "aaZqa", "aИФbВВNNnnвВNЦl"},
		},
		{
			glob: "a[a-]Z[Q]st{LT}*{NN}I*[z-][a]ST{lt}",
			want: &Glob{
				Glob:   "a[a-]Z[Q]st{LT}*{NN}I*[z-][a]ST{lt}",
				Node:   "aaZQstLT*NNI*zaSTlt",
				MinLen: 17, MaxLen: -1, Prefix: "aaZQstLT", Suffix: "zaSTlt",
				Items: []items.NodeItem{
					{Node: "*", Item: items.Star(0)},
					{Node: "NNI", Item: items.NewString("NNI")},
					{Node: "*", Item: items.Star(0)},
				},
			},
			verify: `^aaZQstLT.*NNI.*zaSTlt$`,
			match:  []string{"aaZQstLTNNIzaSTlt", "aaZQstLTaNNIzaSTlt"},
			miss:   []string{"", "ab", "aaZqa", "aaZQstLTaNNIzaSTltl"},
		},
		{
			glob: "a[a-]Z[Q]*[a]c",
			want: &Glob{
				Glob:   "a[a-]Z[Q]*[a]c",
				Node:   "aaZQ*ac",
				MinLen: 6, MaxLen: -1, Prefix: "aaZQ", Suffix: "ac",
				Items: []items.NodeItem{
					{Node: "*", Item: items.Star(0)},
				},
			},
			verify: `^aaZQ.*ac$`,
			match:  []string{"aaZQac", "aaZQaac", "aaZQbac"},
			miss:   []string{"", "ab", "aaZqa"},
		},
	}
	for n, tt := range tests {
		runTestGlob(t, n, tt)
	}
}

// benchmark for suffix miss
var (
	globSuffixMissAscii   = strings.Repeat("find*", 20) + "LAST_NOT*ENDOFT"
	stringSuffixMissAscii = strings.Repeat("findПОИСК", 20) + "LAST_ENDOF"
)

func Benchmark_Suffix_Miss_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globSuffixMissAscii)
		if g.Match(stringSuffixMissAscii) {
			b.Fatal(stringSuffixMissAscii)
		}
	}
}

func Benchmark_Suffix_Miss_ASCII_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globSuffixMissAscii)
		if w.MatchString(stringSuffixMissAscii) {
			b.Fatal(stringSuffixMissAscii)
		}
	}
}

func Benchmark_Suffix_Miss_ASCII_Precompiled(b *testing.B) {
	g := ParseMust(globSuffixMissAscii)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.Match(stringSuffixMissAscii) {
			b.Fatal(stringSuffixMissAscii)
		}
	}
}

func Benchmark_Suffix_Miss_ASCII_Regex_Precompiled(b *testing.B) {
	g := tests.BuildGlobRegexp(globSuffixMissAscii)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.MatchString(stringSuffixMissAscii) {
			b.Fatal(stringSuffixMissAscii)
		}
	}
}

// benchmark forstring gready skip scan optimization (ASCII)
var (
	globStarStringMissAscii = strings.Repeat("find*", 20) + "LAST_NOT*ENDOF"
	stringStarEndAscii      = strings.Repeat("findПОИСК", 20) + "LAST_ENDOF"
)

func Benchmark_Star_String_Miss_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globStarStringMissAscii)
		if g.Match(stringStarEndAscii) {
			b.Fatal(stringStarEndAscii)
		}
	}
}

func Benchmark_Star_String_Miss_ASCII_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globStarStringMissAscii)
		if w.MatchString(stringStarEndAscii) {
			b.Fatal(stringStarEndAscii)
		}
	}
}

func Benchmark_Star_String_Miss_ASCII_Precompiled(b *testing.B) {
	g := ParseMust(globStarStringMissAscii)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.Match(stringStarEndAscii) {
			b.Fatal(stringStarEndAscii)
		}
	}
}

func Benchmark_Star_String_Miss_ASCII_Regex_Precompiled(b *testing.B) {
	g := tests.BuildGlobRegexp(globStarStringMissAscii)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.MatchString(stringStarEndAscii) {
			b.Fatal(stringStarEndAscii)
		}
	}
}

// benchmark for string gready skip scan optimization (Unicode)
var (
	globStarStringMissUnicode = strings.Repeat("ПОИСК*", 20) + "ЛАСТ_НЕТ*КОНЕЦ"
	stringStarEndUnicode      = strings.Repeat("ПОИСКfind", 20) + "ЛАСТ_КОНЕЦ"
)

func Benchmark_Star_String_Miss_Unicode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globStarStringMissUnicode)
		if g.Match(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}

func Benchmark_Star_String_Miss_Unicode_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globStarStringMissUnicode)
		if w.MatchString(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}

func Benchmark_Star_String_Miss_Unicode_Precompiled(b *testing.B) {
	g := ParseMust(globStarStringMissUnicode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.Match(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}

func Benchmark_Star_String_Miss_Unicode_Regex_Precompiled(b *testing.B) {
	g := tests.BuildGlobRegexp(globStarStringMissUnicode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.MatchString(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}

// benchmark for string gready skip scan optimization (Unicode)
var (
	globStarSuffixMissUnicode = strings.Repeat("ПОИСК*", 20) + "ЛАСТ_НЕТ*ЙОНЕЦ"
)

func Benchmark_Star_Suffix_Miss_Unicode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globStarSuffixMissUnicode)
		if g.Match(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}

func Benchmark_Star_Suffix_Miss_Unicode_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globStarSuffixMissUnicode)
		if w.MatchString(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}

func Benchmark_Star_Suffix_Miss_Unicode_Precompiled(b *testing.B) {
	g := ParseMust(globStarSuffixMissUnicode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.Match(stringStarEndUnicode) {
			b.Fail()
		}
	}
}

func Benchmark_Star_Suffix_Miss_Unicode_Regex_Precompiled(b *testing.B) {
	g := tests.BuildGlobRegexp(globStarSuffixMissUnicode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.MatchString(stringStarEndUnicode) {
			b.Fatal(stringStarEndUnicode)
		}
	}
}
