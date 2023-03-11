package gglob

import (
	"testing"
	"time"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func BenchmarkBatchHuge_Moira_Tree_Precompiled2(b *testing.B) {
	g := parseGGlobs(globsBatchHugeMoira)
	w := NewTree()
	for j := 0; j < len(g); j++ {
		_, _, err := w.AddGlob(g[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	pathsBatchHugeMoira := generatePaths(gGlobsBatchHugeMoira, b.N)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathsBatchHugeMoira[i], &globs, nil, &first)
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Moira_GGlob_Precompiled(b *testing.B) {
	g := parseGGlobs(globsBatchHugeMoira)
	pathsBatchHugeMoira := generatePaths(gGlobsBatchHugeMoira, b.N)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for k := 0; k < len(g); k++ {
			_ = g[k].Match(pathsBatchHugeMoira[i])
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Moira_GGlob_Prealloc_ByParts(b *testing.B) {
	g := parseGGlobs(globsBatchHugeMoira)
	parts := make([]string, 8)
	pathsBatchHugeMoira := generatePaths(gGlobsBatchHugeMoira, b.N)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PathSplitB(pathsBatchHugeMoira[i], &parts)
		length := len(pathsBatchHugeMoira[i])
		for k := 0; k < len(g); k++ {
			_ = g[k].MatchByParts(parts, length)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Moira_Tree_Prealloc(b *testing.B) {
	g := parseGGlobs(globsBatchHugeMoira)
	w := NewTree()
	for j := 0; j < len(g); j++ {
		_, _, err := w.AddGlob(g[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	pathsBatchHugeMoira := generatePaths(gGlobsBatchHugeMoira, b.N)
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		first.Init()
		_ = w.Match(pathsBatchHugeMoira[i], &globs, nil, &first)
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Moira_Tree_Prealloc_ByParts(b *testing.B) {
	g := parseGGlobs(globsBatchHugeMoira)
	w := NewTree()
	for j := 0; j < len(g); j++ {
		_, _, err := w.AddGlob(g[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	pathsBatchHugeMoira := generatePaths(gGlobsBatchHugeMoira, b.N)
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	parts := make([]string, 10)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		first.Init()
		_ = PathSplitB(pathsBatchHugeMoira[i], &parts)
		_ = w.MatchByParts(parts, &globs, nil, &first)
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}
