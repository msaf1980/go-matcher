package expand

import (
	"fmt"
	"testing"
)

func BenchmarkExpand(b *testing.B) {
	var input = []struct {
		in    string
		depth int
		max   int
	}{
		{in: "1[b-e]2[a-c]3", max: -1},
		{in: "1[b-e]2[a-c]3", max: -2},
		{in: "232{ad,fdff,wwwww,asdasd}[z-A]", max: -1},
		{in: "mon-mon{i,y,ie}[13-0529]", max: -1},
		{in: "metric.{us,ru,en,de,dk,gb,in}server[1-4].cpu.[0-3].{idle,sys,user}", max: 7},
		{in: "metric.{us,ru,en,de,dk,gb,in}server[1-4].cpu.[0-3].{idle,sys,user}", max: 28},
		{in: "metric.{us,ru,en,de,dk,gb,in}server[1-4].cpu.[0-3].{idle,sys,user}", max: -1},
		{in: "1[b-e]2[a-c]3*", max: -1},
		{in: "1*[b-e]2[a-c]3", max: -1},
	}
	for _, bench := range input {
		b.Run(fmt.Sprintf("%s [%d]", bench.in, bench.max), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, _ = Expand(bench.in, bench.max, bench.depth)
			}
		})
		b.Run(fmt.Sprintf("%s [%d] try", bench.in, bench.max), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, _ = ExpandTry(bench.in, bench.max, bench.depth)
			}
		})
	}
	b.Log("\n")
}
