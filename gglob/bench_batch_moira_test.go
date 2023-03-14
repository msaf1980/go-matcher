package gglob

import (
	"math/rand"
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
		// _ = w.Match(pathsBatchHugeMoira[i], &globs, nil, &first)
		path := pathsBatchHugeMoira[rand.Intn(len(pathsBatchHugeMoira))]
		_ = w.Match(path, &globs, nil, &first)
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
		path := pathsBatchHugeMoira[rand.Intn(len(pathsBatchHugeMoira))]
		for k := 0; k < len(g); k++ {
			// _ = g[k].Match(pathsBatchHugeMoira[i])
			_ = g[k].Match(path)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Moira_GGlob_ByParts_Prealloc(b *testing.B) {
	g := parseGGlobs(globsBatchHugeMoira)
	parts := make([]string, 8)
	pathsBatchHugeMoira := generatePaths(gGlobsBatchHugeMoira, b.N)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := pathsBatchHugeMoira[rand.Intn(len(pathsBatchHugeMoira))]
		length := len(path)
		PathSplitB(path, &parts)
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
		path := pathsBatchHugeMoira[rand.Intn(len(pathsBatchHugeMoira))]
		_ = w.Match(path, &globs, nil, &first)
		// _ = w.Match(pathsBatchHugeMoira[i], &globs, nil, &first)
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
		// _ = PathSplitB(pathsBatchHugeMoira[i], &parts)
		path := pathsBatchHugeMoira[rand.Intn(len(pathsBatchHugeMoira))]
		_ = PathSplitB(path, &parts)
		_ = w.MatchByParts(parts, &globs, nil, &first)
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}
