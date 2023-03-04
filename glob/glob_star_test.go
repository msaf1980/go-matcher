package glob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlob_Star(t *testing.T) {
	tests := []testGlob{
		// any
		{
			glob: "*",
			want: &Glob{
				Glob: "*", Node: "*",
				MaxLen: -1,
				Items:  []items.Item{items.Star(0)},
			},
			verify: "^.*$",
			match:  []string{"", "ac", "abc", "abcc"},
		},
		{
			glob: "**",
			want: &Glob{
				Glob: "**", Node: "*",
				MaxLen: -1,
				Items:  []items.Item{items.Star(0)},
			},
			verify: "^.*$",
			match:  []string{"", "ac", "abc", "abcc"},
		},
		// deduplication
		{
			glob: "a******c",
			want: &Glob{
				Glob: "a******c", Node: "a*c",
				MinLen: 2, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(0)},
			},
			verify: "^a.*c$",
			match:  []string{"ac", "abc", "abcc"},
			miss:   []string{"", "acb", "ac.", "ac.b"},
		},
		// * match
		{
			glob: "*",
			want: &Glob{
				Glob: "*", Node: "*",
				MinLen: 0, MaxLen: -1,
				Items: []items.Item{items.Star(0)},
			},
			verify: "^.*$",
			match:  []string{"", "a", "ac", "abc", "abcc", "ac.", "ac.b"},
		},
		{
			glob: "a*c",
			want: &Glob{
				Glob: "a*c", Node: "a*c",
				MinLen: 2, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(0)},
			},
			verify: "^a.*c$",
			match:  []string{"ac", "acc", "aec", "aebc", "aecc", "aecec", "abecec", "a.c"},
			miss:   []string{"", "ab", "c", "ace"},
		},
		{
			glob: "a*",
			want: &Glob{
				Glob: "a*", Node: "a*",
				MinLen: 1, MaxLen: -1, Prefix: "a",
				Items: []items.Item{items.Star(0)},
			},
			verify: "^a.*$",
			match:  []string{"ac", "ab", "acc", "ace", "aec", "aebc", "aecc", "aecec", "abecec", "a.c"},
			miss:   []string{"", "c"},
		},
		// composite
		{
			glob: "a*b?c",
			want: &Glob{
				Glob: "a*b?c", Node: "a*b?c",
				MinLen: 4, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(0), items.Byte('b'), items.Any(1)},
			},
			verify: "^a.*b.c$",
			match:  []string{"abec", "abbec", "acbbc", "aecbec"},
			miss:   []string{"", "ab", "c", "ace", "a.c", "abbece"},
		},
		{
			glob: "a*?_FIND*_st",
			want: &Glob{
				Glob: "a*?_FIND*_st", Node: "a*?_FIND*_st",
				MinLen: 10, MaxLen: -1, Prefix: "a", Suffix: "_st",
				Items: []items.Item{items.Star(1), items.NewString("_FIND"), items.Star(0)},
			},
			verify: "^a.*._FIND.*_st$",
			match: []string{
				"ab_FIND_st", "aLc_FIND_st", "aLBc_FIND_st", "aLBc_FIND_STAR_st", "aLBc_FINDB_st",
			},
			miss: []string{"a_FIND_st", "a_FINDB_st"},
		},
		// optimization
		{
			glob: "a******?c",
			want: &Glob{
				Glob: "a******?c", Node: "a*?c",
				MinLen: 3, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(1)},
			},
			verify: "^a.*.c$",
			match:  []string{"aBc", "aBCc", "aBCDc", "aBCDEc"},
			miss:   []string{"", "ac", "acb"},
		},
		{
			glob: "a?******c",
			want: &Glob{
				Glob: "a?******c", Node: "a*?c",
				MinLen: 3, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(1)},
			},
			match: []string{"aBc", "aBCc", "aBCDc", "aBCDEc"},
			miss:  []string{"", "ac", "acb"},
		},
		{
			glob: "a**?****c",
			want: &Glob{
				Glob: "a**?****c", Node: "a*?c",
				MinLen: 3, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(1)},
			},
			match: []string{"aBc", "aBCc", "aBCDc", "aBCDEc"},
			miss:  []string{"", "ac", "acb"},
		},
		{
			glob: "a?******?c",
			want: &Glob{
				Glob: "a?******?c", Node: "a*??c",
				MinLen: 4, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(2)},
			},
			verify: "^a.*.{2}c$",
			match:  []string{"aBCc", "aBCDc", "aBCDEc", "a.Cc", "aB.Dc", "aBC..c", "aBCc.aBCc"},
			miss:   []string{"abc", "aBc", "b", "cBc", "ZBCc", "aBCc.", "aBCc.a"},
		},
		{
			glob: "a**??c",
			want: &Glob{
				Glob: "a**??c", Node: "a*??c",
				MinLen: 4, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(2)},
			},
			match: []string{"aBCc", "aBCDc", "aBCDEc", "a.Cc", "aB.Dc", "aBC..c", "aBCc.aBCc"},
			miss:  []string{"abc", "aBc", "b", "cBc", "ZBCc", "aBCc.", "aBCc.a"},
		},
		{
			glob: "a*?*?*c",
			want: &Glob{
				Glob: "a*?*?*c", Node: "a*??c",
				MinLen: 4, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(2)},
			},
			match: []string{"aBCc", "aBCDc", "aBCDEc", "a.Cc", "aB.Dc", "aBC..c", "aBCc.aBCc"},
			miss:  []string{"abc", "aBc", "b", "cBc", "ZBCc", "aBCc.", "aBCc.a"},
		},
		{
			glob: "a?*??c",
			want: &Glob{
				Glob: "a?*??c", Node: "a*???c",
				MinLen: 5, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(3)},
			},
			verify: "^a.*.{3}c$",
			match:  []string{"aBCDc", "aBCDEc", "a..Cc", "aB.Dc", "aBC..c", "aBCc.aBCc"},
			miss:   []string{"abc", "aBc", "aBCc", "b", "cBc", "ZBCc", "a.Cc", "aBCc.", "aBCc.a"},
		},
		{
			glob: "a??*?*c",
			want: &Glob{
				Glob: "a??*?*c", Node: "a*???c",
				MinLen: 5, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(3)},
			},
			match: []string{"aBCDc", "aBCDEc", "a..Cc", "aB.Dc", "aBC..c", "aBCc.aBCc"},
			miss:  []string{"abc", "aBc", "aBCc", "b", "cBc", "ZBCc", "a.Cc", "aBCc.", "aBCc.a"},
		},
		{
			glob: "a???**c",
			want: &Glob{
				Glob: "a???**c", Node: "a*???c",
				MinLen: 5, MaxLen: -1, Prefix: "a", Suffix: "c",
				Items: []items.Item{items.Star(3)},
			},
			match: []string{"aBCDc", "aBCDEc", "a..Cc", "aB.Dc", "aBC..c", "aBCc.aBCc"},
			miss:  []string{"abc", "aBc", "aBCc", "b", "cBc", "ZBCc", "a.Cc", "aBCc.", "aBCc.a"},
		},
	}
	for n, tt := range tests {
		runTestGlob(t, n, tt)
	}
}
