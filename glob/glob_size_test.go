package glob

import (
	"strings"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/tests"
)

var (
	globMaxSizeCheck   = strings.Repeat("find*", 90) + "LAST_Й*ENDOF"
	stringMaxSizeCheck = strings.Repeat("findПОИСК", 10) + "LAST_Й_ENDOF"
)

func Benchmark_Size_Max(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globMaxSizeCheck)
		if g.Match(stringMaxSizeCheck) {
			b.Fatal(stringMaxSizeCheck)
		}
	}
}

func _Benchmark_Size_Max_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globMaxSizeCheck)
		if w.MatchString(stringMaxSizeCheck) {
			b.Fatal(stringMaxSizeCheck)
		}
	}
}

func Benchmark_Size_Max_Precompiled(b *testing.B) {
	g := ParseMust(globMaxSizeCheck)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.Match(stringMaxSizeCheck) {
			b.Fatal(stringMaxSizeCheck)
		}
	}
}

func _Benchmark_Size_Max_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globMaxSizeCheck)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(stringMaxSizeCheck) {
			b.Fatal(stringMaxSizeCheck)
		}
	}
}
