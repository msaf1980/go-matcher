package expand

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpand(t *testing.T) {
	type data struct {
		in  string
		max int
		out []string
	}

	tests := []data{
		{in: "{b,c}", max: -1, out: []string{"b", "c"}},
		{in: "a{b,c}d", max: -1, out: []string{"abd", "acd"}},
		{in: "a{,b,c}d", max: -1, out: []string{"ad", "abd", "acd"}},
		{in: "a{b,,c}d", max: -1, out: []string{"abd", "ad", "acd"}},
		{in: "a{b,c,}d", max: -1, out: []string{"abd", "acd", "ad"}},
		{in: "a{b,}d", max: -1, out: []string{"abd", "ad"}},
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
		{in: "as{12,32}[a-c]{2}", max: -2, out: []string{"as12[a-c]2", "as32[a-c]2"}},                                 // expand only first founded node
		{in: "as{12,32}{2}[a-c]", max: -2, out: []string{"as122[a-c]", "as322[a-c]"}},                                 // expand only first founded node
		{in: "as{12,32}2[a-c]", max: -2, out: []string{"as122[a-c]", "as322[a-c]"}},                                   // expand only first founded node
		{in: "as{12,32}[a-c]{2}", max: -3, out: []string{"as12a2", "as12b2", "as12c2", "as32a2", "as32b2", "as32c2"}}, // expand only two founded nodes
		{in: "as{12,32}2[a-c]", max: -3, out: []string{"as122a", "as122b", "as122c", "as322a", "as322b", "as322c"}},   // expand only two founded nodes
		{in: "as{12,32}{2}[a-c]", max: -3, out: []string{"as122a", "as122b", "as122c", "as322a", "as322b", "as322c"}}, // expand only two founded nodes
		{in: "as{12,32}[a-c]{2,a}", max: -3, out: []string{ // expand only two founded nodes
			"as12a{2,a}", "as12b{2,a}", "as12c{2,a}", "as32a{2,a}", "as32b{2,a}", "as32c{2,a}",
		}},
		{in: "as{12,32}[a-c]{2,a}", max: -4, out: []string{ // expand only three founded nodes
			"as12a2", "as12aa", "as12b2", "as12ba", "as12c2", "as12ca", "as32a2", "as32aa", "as32b2", "as32ba", "as32c2", "as32ca",
		}},
	}

	for n, tt := range tests {
		t.Run(fmt.Sprintf("[%d] [%d] %s", n, tt.max, tt.in), func(t *testing.T) {
			result, err := Expand(tt.in, tt.max)
			require.NoError(t, err, "input %q", tt.in)
			assert.Equal(t, tt.out, result, "input %q", tt.in)
		})
	}
}

func TestGetPair(t *testing.T) {
	type data struct {
		in    string
		start int
		stop  int
	}
	tests := []data{
		{"x{12,b}xxxxx", 1, 6},
		{"x[124-5]xxxxx", 1, 7},
		// complex, not expand at now
		{"{x{12,{}}xxxxx", -1, -1},
		{"x{12,{}}{{,13}", -1, -1},
		{"{x{12,{}}{{,13}", -1, -1},
		{"{{x{12,{}}{{,13}", -1, -1},
		// unclosed
		{"}some{", -1, -1},
		{"}some", -1, -1},
		{"some{", -1, -1},
		{"]some[", -1, -1},
		{"]some", -1, -1},
		{"some[", -1, -1},
	}

	for _, tt := range tests {
		start, stop := getPair(tt.in)
		assert.Equal(t, tt.start, start, "start of %q", tt.in)
		assert.Equal(t, tt.stop, stop, "stop of %q", tt.in)
	}
}
