package expand

import (
	"fmt"
	"testing"
)

func BenchmarkExpand(b *testing.B) {
	var input = []struct {
		in  string
		max int
	}{
		{in: "1[b-e]2[a-c]3"},
		{in: "232{ad,fdff,wwwww,asdasd}[z-A]"},
		{in: "mon-mon{i,y,ie}[13-0529]"},
		{in: "metric.{us,ru,en,de,dk,gb,in}server[1-4].cpu.[0-3].{idle,sys,user}", max: 7},
		{in: "metric.{us,ru,en,de,dk,gb,in}server[1-4].cpu.[0-3].{idle,sys,user}", max: 28},
		{in: "metric.{us,ru,en,de,dk,gb,in}server[1-4].cpu.[0-3].{idle,sys,user}"},
	}
	for _, bench := range input {
		b.Run(fmt.Sprintf("%s [%d]", bench.in, bench.max), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, _ = Expand(bench.in, bench.max)
			}
		})
	}
	b.Log("\n")
}
