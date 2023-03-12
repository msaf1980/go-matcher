package glob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

func TestGlob_Group(t *testing.T) {
	tests := []testGlob{
		{
			glob: "{b*,a?cd*,cd[a-z]}bc*c*e",
			want: &Glob{
				Glob: "{b*,a?cd*,cd[a-z]}bc*c*e", Node: "{a?cd*,b*,cd[a-z]}bc*c*e",
				Suffix: "e", MinLen: 5, MaxLen: -1,
				Items: []items.Item{
					&items.Group{
						MinSize: 1, MaxSize: -1,
						Vals: []items.Item{
							&items.Chain{
								Items: []items.Item{
									items.Byte('a'), items.Any(1), items.NewString("cd"), items.Star(0),
								},
								MinSize: 4, MaxSize: -1,
							},
							&items.Chain{
								Items: []items.Item{items.Byte('b'), items.Star(0)}, MinSize: 1, MaxSize: -1,
							},
							&items.Chain{
								Items:   []items.Item{items.NewString("cd"), items.NewRunesRanges("[a-z]")},
								MinSize: 3, MaxSize: 3,
							},
						},
					},
					items.NewString("bc"), items.Star(0), items.Byte('c'), items.Star(0),
				},
			},
			match: []string{"aZcdbcce", "aZcdQAbcZWcIe", "cdqbcZcIe", "bCDbcZIce"},
			miss:  []string{"", "aZcdbcc", "aZcdcce", "aZcdQAbcZWIe"},
		},
		{
			glob: "?*{b*,a?cd*,cd[a-z]}bc*c*e",
			want: &Glob{
				Glob: "?*{b*,a?cd*,cd[a-z]}bc*c*e", Node: "*?{a?cd*,b*,cd[a-z]}bc*c*e",
				Suffix: "e", MinLen: 6, MaxLen: -1,
				Items: []items.Item{
					items.Star(1),
					&items.Group{
						MinSize: 1, MaxSize: -1,
						Vals: []items.Item{
							&items.Chain{
								Items: []items.Item{
									items.Byte('a'), items.Any(1), items.NewString("cd"), items.Star(0),
								},
								MinSize: 4, MaxSize: -1,
							},
							&items.Chain{
								Items: []items.Item{items.Byte('b'), items.Star(0)}, MinSize: 1, MaxSize: -1,
							},
							&items.Chain{
								Items:   []items.Item{items.NewString("cd"), items.NewRunesRanges("[a-z]")},
								MinSize: 3, MaxSize: 3,
							},
						},
					},
					items.NewString("bc"), items.Star(0), items.Byte('c'), items.Star(0),
				},
			},
			match: []string{"ZaZcdbcce", "ЙaЮcdQAbcZWcIe", "ЙacdqbcZcIe", "ЙabCDbcZIce"},
			miss:  []string{"", "aЯcdbcce", "aZcdQAbcZWcIe", "aZcdbcc", "aZcdcce", "aZcdQAbcZWIe"},
		},
		{
			glob: "?*{b*,a?cd*,cd[a-z]}bc*c*e",
			want: &Glob{
				Glob: "?*{b*,a?cd*,cd[a-z]}bc*c*e", Node: "*?{a?cd*,b*,cd[a-z]}bc*c*e",
				Suffix: "e", MinLen: 6, MaxLen: -1,
				Items: []items.Item{
					items.Star(1),
					&items.Group{
						MinSize: 1, MaxSize: -1,
						Vals: []items.Item{
							&items.Chain{
								Items: []items.Item{
									items.Byte('a'), items.Any(1), items.NewString("cd"), items.Star(0),
								},
								MinSize: 4, MaxSize: -1,
							},
							&items.Chain{
								Items: []items.Item{items.Byte('b'), items.Star(0)}, MinSize: 1, MaxSize: -1,
							},
							&items.Chain{
								Items:   []items.Item{items.NewString("cd"), items.NewRunesRanges("[a-z]")},
								MinSize: 3, MaxSize: 3,
							},
						},
					},
					items.NewString("bc"), items.Star(0), items.Byte('c'), items.Star(0),
				},
			},
			match: []string{"ZaZcdbcce", "ЙaЮcdQAbcZWcIe", "ЙacdqbcZcIe", "ЙabCDbcZIce"},
			miss:  []string{"", "aЯcdbcce", "aZcdQAbcZWcIe", "aZcdbcc", "aZcdcce", "aZcdQAbcZWIe"},
		},
	}
	for n, tt := range tests {
		runTestGlob(t, n, tt)
	}
}

// becnmark for group
var (
	globGroup   = "{b*,a?cd*,cd[a-z]}bc*c*e"
	stringGroup = "cdqbcZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZcIe"
)

func Benchmark_Group(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globGroup)
		if !g.Match(stringGroup) {
			b.Fatal(stringGroup)
		}
	}
}

func _Benchmark_Group_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globGroup)
		if !w.MatchString(stringGroup) {
			b.Fatal(stringGroup)
		}
	}
}

func Benchmark_Group_Precompiled(b *testing.B) {
	g := ParseMust(globGroup)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringGroup) {
			b.Fatal(stringGroup)
		}
	}
}

// becnmark for group after star
var (
	globStarGroup   = "*{b*,a?cd*,cd[a-z]}bc*c*e"
	stringStarGroup = "ZZZZZZZZcZZZZZZZZZZZZcdqbcZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZcIe"
)

func Benchmark_Star_Group(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globStarGroup)
		if !g.Match(stringStarGroup) {
			b.Fatal(stringStarGroup)
		}
	}
}

func _Benchmark_Star_Group_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globStarGroup)
		if !w.MatchString(stringStarGroup) {
			b.Fatal(stringStarGroup)
		}
	}
}

func Benchmark_Star_Group_Precompiled(b *testing.B) {
	g := ParseMust(globStarGroup)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringStarGroup) {
			b.Fatal(stringStarGroup)
		}
	}
}
