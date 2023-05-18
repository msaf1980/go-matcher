package expand

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseExpr(t *testing.T) {
	tests := []struct {
		in   string
		want []Expression
	}{
		{
			in: "x{12,b}xxxxx",
			want: []Expression{
				{body: "x"},
				{typ: expList, body: "{12,b}", list: []string{"12", "b"}},
				{body: "xxxxx"},
			},
		},
		{
			in: "ab[124-5]xxxxx",
			want: []Expression{
				{body: "ab"},
				{typ: expRunes, body: "[124-5]", runes: runesRangeExpandMust("124-5")},
				{body: "xxxxx"},
			},
		},
		// complex, not expand at now
		{
			in: "{x{12,{}}xxxxx",
			want: []Expression{
				{typ: expWildcard, body: "{x{12,{}}"},
				{body: "xxxxx"},
			},
		},
		{
			in: "x{12,{}}{{,13}",
			want: []Expression{
				{body: "x"},
				{typ: expWildcard, body: "{12,{}}{{,13}"},
			},
		},
		{
			in: "{x{12,{}}{{,13}",
			want: []Expression{
				{typ: expWildcard, body: "{x{12,{}}{{,13}"},
			},
		},
		{
			in: "{{x{12,{}}{{,13}",
			want: []Expression{
				{typ: expWildcard, body: "{{x{12,{}}{{,13}"},
			},
		},
		// unclosed
		{
			in: "}some{",
			want: []Expression{
				{typ: expWildcard, body: "}some{"},
			},
		},
		{
			in: "}some",
			want: []Expression{
				{typ: expWildcard, body: "}"},
				{body: "some"},
			},
		},
		{
			in: "some{",
			want: []Expression{
				{body: "some"},
				{typ: expWildcard, body: "{"},
			},
		},
		{
			in: "]some[",
			want: []Expression{
				{typ: expWildcard, body: "]some["},
			},
		},
		{
			in: "]some",
			want: []Expression{
				{typ: expWildcard, body: "]"},
				{body: "some"},
			},
		},
		{
			in: "some[",
			want: []Expression{
				{body: "some"},
				{typ: expWildcard, body: "["},
			},
		},
		// star
		{
			in: "x[124-5]a*xxxx",
			want: []Expression{
				{body: "x"},
				{typ: expRunes, body: "[124-5]", runes: runesRangeExpandMust("124-5")},
				{body: "a"},
				{typ: expWildcard, body: "*"},
				{body: "xxxx"},
			},
		},
		{
			in: "x*[124-5]xxxxx",
			want: []Expression{
				{body: "x"},
				{typ: expWildcard, body: "*[124-5]"},
				{body: "xxxxx"},
			},
		},
		{
			in: "x?[124-5]xxxxx",
			want: []Expression{
				{body: "x"},
				{typ: expWildcard, body: "?[124-5]"},
				{body: "xxxxx"},
			},
		},
		{
			in: "x[*124-5]xxxxx",
			want: []Expression{
				{body: "x"},
				{typ: expWildcard, body: "[*124-5]"},
				{body: "xxxxx"},
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %s", n, tt.in), func(t *testing.T) {
			gotExps := ParseExpr(tt.in)
			assert.Equal(t, tt.want, gotExps.exps, "exps")
		})
	}
}

func TestExpand(t *testing.T) {
	type data struct {
		in    string
		max   int
		depth int
		out   []string
	}

	tests := []data{
		{in: "{b,c}", max: -1, out: []string{"b", "c"}},
		{in: "a{b,c}d", max: -1, out: []string{"abd", "acd"}},
		{in: "a{,b,c}d", max: -1, out: []string{"ad", "abd", "acd"}},
		{in: "a{b,,c}d", max: -1, out: []string{"abd", "ad", "acd"}},
		{in: "a{b,c,}d", max: -1, out: []string{"abd", "acd", "ad"}},
		{in: "a{b,}d", max: -1, out: []string{"abd", "ad"}},
		{in: "a{b}d", max: -1, out: []string{"abd"}},
		{in: "a{b,c,}d", max: -1, out: []string{"abd", "acd", "ad"}},
		{in: "[2-4]", max: -1, out: []string{"2", "3", "4"}},
		{in: "[24]", max: -1, out: []string{"2", "4"}},
		{in: "[243]", max: -1, out: []string{"2", "3", "4"}},
		{in: "1[b-e]2[a-c]3", max: -1, out: []string{"1b2a3", "1b2b3", "1b2c3", "1c2a3", "1c2b3", "1c2c3", "1d2a3", "1d2b3", "1d2c3", "1e2a3", "1e2b3", "1e2c3"}},
		{in: "as{12,32}[a-c]{2}", max: -1, out: []string{"as12a2", "as12b2", "as12c2", "as32a2", "as32b2", "as32c2"}},
		{in: "as{12,32}[a-c]{2}", max: 1, out: []string{"as{12,32}[a-c]2"}},
		{in: "as{12,32}[a-c]{2}", max: 2, out: []string{"as12[a-c]2", "as32[a-c]2"}},
		{in: "as{12,32}[a-c]{2}", max: 5, out: []string{"as12[a-c]2", "as32[a-c]2"}},
		{in: "as{12,32}[a-c]{2}", max: 6, out: []string{"as12a2", "as12b2", "as12c2", "as32a2", "as32b2", "as32c2"}},
		{in: "as{12,32}.[a-c].{2}", max: -1, depth: 1, out: []string{"as12.[a-c].2", "as32.[a-c].2"}},                                  // expand only first founded node
		{in: "as{12,32}{2}.[a-c]", max: -1, depth: 1, out: []string{"as122.[a-c]", "as322.[a-c]"}},                                     // expand only first founded node
		{in: "as{12,32}.2[a-c]", max: -1, depth: 1, out: []string{"as12.2[a-c]", "as32.2[a-c]"}},                                       // expand only first founded node
		{in: "as{12,32}[a-c]{2}", max: -1, depth: 2, out: []string{"as12a2", "as12b2", "as12c2", "as32a2", "as32b2", "as32c2"}},        // expand only two founded nodes
		{in: "as{12,32}.{2}[a-c]", max: -1, depth: 2, out: []string{"as12.2a", "as12.2b", "as12.2c", "as32.2a", "as32.2b", "as32.2c"}}, // expand only two founded nodes
		{in: "as{12,32}.[a-c].{2,a}", max: -1, depth: 2, out: []string{ // expand only two founded nodes
			"as12.a.{2,a}", "as12.b.{2,a}", "as12.c.{2,a}", "as32.a.{2,a}", "as32.b.{2,a}", "as32.c.{2,a}",
		}},
		{in: "as{12,32}.[a-c].{2,a}", max: -1, depth: 3, out: []string{ // expand only three founded nodes
			"as12.a.2", "as12.a.a", "as12.b.2", "as12.b.a", "as12.c.2", "as12.c.a", "as32.a.2", "as32.a.a", "as32.b.2", "as32.b.a", "as32.c.2", "as32.c.a",
		}},
		// star
		{in: "a{b,c}*d", max: -1, out: []string{"ab*d", "ac*d"}},
		{in: "a*{b,c}d", max: -1, out: []string{"a*{b,c}d"}}, // no reverse expand
		{in: "a?{b,c}d", max: -1, out: []string{"a?{b,c}d"}}, // no reverse expand
		{in: "a{*b,c}d", max: -1, out: []string{"a{*b,c}d"}},
	}

	for n, tt := range tests {
		t.Run(fmt.Sprintf("[%d] [%d:%d] %s", n, tt.max, tt.depth, tt.in), func(t *testing.T) {
			result, err := Expand(tt.in, tt.max, tt.depth)
			require.NoError(t, err, "input %q", tt.in)
			assert.Equal(t, tt.out, result, "input %q", tt.in)
		})
	}
}
