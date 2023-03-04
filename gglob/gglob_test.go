package gglob

import (
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/msaf1980/go-matcher/pkg/tests"
	"github.com/msaf1980/go-matcher/pkg/utils"
	"github.com/stretchr/testify/assert"
)

type tGGlob struct {
	Glob string // raw glob
	Node string // optimized glob or value string if len(Inners) == 0

	MinLen int // min bytes len
	MaxLen int // -1 for unlimited

	Parts []string
}

func newTGGlob(g *GGlob) *tGGlob {
	t := &tGGlob{
		Glob:   g.Glob,
		Node:   g.Node,
		MinLen: g.MinLen,
		MaxLen: g.MaxLen,
		Parts:  make([]string, len(g.Parts)),
	}
	for i := 0; i < len(g.Parts); i++ {
		t.Parts[i] = g.Parts[i].String()
	}
	return t
}

type testGGlob struct {
	glob    string
	verify  string
	want    *tGGlob
	wantErr bool
	match   []string // must match with glob
	miss    []string
}

func runTestGGlob(t *testing.T, n int, tt testGGlob) {
	t.Run(strconv.Itoa(n)+"#"+tt.glob, func(t *testing.T) {
		g, err := Parse(tt.glob)

		if (err != nil) != tt.wantErr {
			t.Fatalf("Parse(%q) error = %v, wantErr %v", tt.glob, err, tt.wantErr)
		}
		if tt.wantErr {
			assert.Equal(t, 0, len(tt.match), "can't check on error")
			assert.Equal(t, 0, len(tt.miss), "can't check on error")
		} else {
			tg := newTGGlob(g)
			if !reflect.DeepEqual(tg, tt.want) {
				t.Fatalf("GGlob(%q) = %s", tt.glob, cmp.Diff(tt.want, tg))
			}
			verifyGGlob(t, tt.match, tt.miss, g, tt.verify)
		}
	})
}

func verifyGGlob(t *testing.T, match []string, miss []string, g *GGlob, verifyRegexp string) {
	var re *regexp.Regexp
	if verifyRegexp != "" {
		re = regexp.MustCompile(verifyRegexp)
	}
	for n, path := range match {
		t.Run(strconv.Itoa(n)+"#path="+path, func(t *testing.T) {
			matched := g.Match(path)
			if re == nil {
				if !matched {
					t.Errorf("GGlob(%q).Match(%q) = %v, want true", g.Node, path, matched)
				}
			} else {
				reMatched := re.MatchString(path)
				if !matched || matched != reMatched {
					t.Errorf("GGlob(%q).Match(%q) = %v, want true, re(%q) = %v", g.Node, path, matched, verifyRegexp, reMatched)
				}
			}

			parts := PathSplit(path)
			matched = g.MatchByParts(parts, len(path))
			if !matched {
				t.Errorf("GGlob(%q).MatchByParts(%q) = %v, want true", g.Node, path, matched)
			}
		})
	}
	for n, path := range miss {
		t.Run(strconv.Itoa(n)+"#MISS#path="+path, func(t *testing.T) {
			matched := g.Match(path)
			if re == nil {
				if matched {
					t.Errorf("GGlob(%q).Match(%q) = %v, want false", g.Node, path, matched)
				}
			} else {
				reMatched := re.MatchString(path)
				if matched || matched != reMatched {
					t.Errorf("GGlob(%q).Match(%q) = %v, want false, re(%q) = %v", g.Node, path, matched, verifyRegexp, reMatched)
				}
			}
		})
	}
}

func TestGGlob(t *testing.T) {
	tests := []testGGlob{
		{
			glob: "DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount",
			want: &tGGlob{
				Glob:   "DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount",
				Node:   "DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount",
				MinLen: 28,
				MaxLen: -1,
				Parts: []string{
					"DB",
					"*",
					"{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}",
					"*",
					"DownEndpointCount",
				},
			},
			match: []string{
				"DB.Sales.BalanceCluster.node1.DownEndpointCount",
				"DB.Back.WebCluster.node2.DownEndpointCount",
			},
			miss: []string{
				"B.Sales.BalanceCluster.node1.DownEndpointCount",
				"DBA.Back.WebCluster.node2.DownEndpointCount",
				"DB.Sales.BalanceCluster2.node1.DownEndpointCount",
				"DB.Back.WebCluster.node2.DownEndpointCount2",
				"DB.DC1.Sales.BalanceCluster.node1.DownEndpointCount",
			},
		},
		{
			glob:    "DB.*..{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount",
			wantErr: true,
		},
	}
	for n, tt := range tests {
		runTestGGlob(t, n, tt)
	}
}

func generatePaths(globs []*GGlob, count int) []string {
	result := make([]string, 0, count)
	i := 0
	var buf strings.Builder
	for {
		if i >= count {
			break
		}
		n := rand.Intn(len(globs))
		if len(globs[n].Parts) == 0 {
			result = append(result, globs[n].Node)
		} else {
			buf.Reset()
			for j := 0; j < len(globs[n].Parts); j++ {
				if j > 0 {
					buf.WriteByte('.')
				}
				globs[n].Parts[j].WriteRandom(&buf)
			}
			result = append(result, utils.CloneString(buf.String()))
		}
		i++
	}

	return result
}

var (
	globBenchASCII   = "DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount"
	stringBenchASCII = "DB.Sales.BalanceCluster.node1.DownEndpointCount"
)

func Benchmark_GGlob_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globBenchASCII)
		if !g.Match(stringBenchASCII) {
			b.Fatal(stringBenchASCII)
		}
	}
}

func Benchmark_Regex_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globBenchASCII)
		if !w.MatchString(stringBenchASCII) {
			b.Fatal(stringBenchASCII)
		}
	}
}

func Benchmark_GGlob_ASCII_Precompiled(b *testing.B) {
	g := ParseMust(globBenchASCII)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringBenchASCII) {
			b.Fatal(stringBenchASCII)
		}
	}
}

func Benchmark_Regex_ASCII_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globBenchASCII)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !w.MatchString(stringBenchASCII) {
			b.Fatal(stringBenchASCII)
		}
	}
}
