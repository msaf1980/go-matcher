package gtags

import (
	"testing"
	"time"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

var (
	queriesBatchHugeMoira = tests.LoadPatterns("tagged_patterns.txt")
	termsBatchHugeMoira   = taggedTermListList(queriesBatchHugeMoira)
)

func BenchmarkBatchHuge_Tree_ByTags(b *testing.B) {
	start := time.Now()
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, len(termsBatchHugeMoira))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewTree()
		for j := 0; j < len(queriesBatchHugeMoira); j++ {
			_, _, err := w.Add(queriesBatchHugeMoira[j], j)
			if err != nil {
				b.Fatal(err)
			}
		}
		first := items.MinStore{-1}
		for j := 0; j < len(pathsBatchHugeMoira); j++ {
			first.Init()
			tags, _ := PathTags(pathsBatchHugeMoira[j])
			_ = w.MatchByTags(tags, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_GTag_ByPaths(b *testing.B) {
	start := time.Now()
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, len(termsBatchHugeMoira))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g := taggedTermListList(queriesBatchHugeMoira)
		for j := 0; j < len(pathsBatchHugeMoira); j++ {
			tags, _ := PathTags(pathsBatchHugeMoira[j])
			for k := 0; k < len(g); k++ {
				_ = g[k].MatchByTags(tags)
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Tree_Precompiled(b *testing.B) {
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, len(termsBatchHugeMoira))

	w := NewTree()
	for j := 0; j < len(queriesBatchHugeMoira); j++ {
		_, _, err := w.Add(queriesBatchHugeMoira[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var store items.AllStore
		store.Init()
		for j := 0; j < len(pathsBatchHugeMoira); j++ {
			store.Init()
			tags, _ := PathTags(pathsBatchHugeMoira[j])
			_ = w.MatchByTags(tags, &store)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Tree_ByMap_Precompiled(b *testing.B) {
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, len(termsBatchHugeMoira))

	w := NewTree()
	for j := 0; j < len(queriesBatchHugeMoira); j++ {
		_, _, err := w.Add(queriesBatchHugeMoira[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var store items.AllStore
		store.Init()
		for j := 0; j < len(pathsBatchHugeMoira); j++ {
			store.Init()
			tags, _ := PathTagsMap(pathsBatchHugeMoira[j])
			_ = w.MatchByTagsMap(tags, &store)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Tree_ByMap_PrecompiledB(b *testing.B) {
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, len(termsBatchHugeMoira))

	w := NewTree()
	for j := 0; j < len(queriesBatchHugeMoira); j++ {
		_, _, err := w.Add(queriesBatchHugeMoira[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tags := make(map[string]string)
		store := items.MinStore{-1}
		for j := 0; j < len(pathsBatchHugeMoira); j++ {
			store.Init()
			_ = PathTagsMapB(pathsBatchHugeMoira[j], tags)
			_ = w.MatchByTagsMap(tags, &store)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Tree_Prealloc(b *testing.B) {
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, len(termsBatchHugeMoira))

	w := NewTree()
	for j := 0; j < len(queriesBatchHugeMoira); j++ {
		_, _, err := w.Add(queriesBatchHugeMoira[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}

	tagsList := tagsList(pathsBatchHugeMoira)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var store items.AllStore
		store.Init()
		store.Grow(1)
		for j := 0; j < len(tagsList); j++ {
			store.Init()
			_ = w.MatchByTags(tagsList[j], &store)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Tree_ByMap_Prealloc(b *testing.B) {
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, len(termsBatchHugeMoira))

	w := NewTree()
	for j := 0; j < len(queriesBatchHugeMoira); j++ {
		_, _, err := w.Add(queriesBatchHugeMoira[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}

	tagsList := tagsList(pathsBatchHugeMoira)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var store items.AllStore
		store.Init()
		store.Grow(1)
		for j := 0; j < len(tagsList); j++ {
			store.Init()
			_ = w.MatchByTags(tagsList[j], &store)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_GTag_Precompiled(b *testing.B) {
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, len(termsBatchHugeMoira))
	g := taggedTermListList(queriesBatchHugeMoira)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(pathsBatchHugeMoira); j++ {
			tags, _ := PathTags(pathsBatchHugeMoira[j])
			for k := 0; k < len(g); k++ {
				_ = g[k].MatchByTags(tags)
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_GTag_Prealloc(b *testing.B) {
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, len(termsBatchHugeMoira))
	g := taggedTermListList(queriesBatchHugeMoira)
	tagsList := tagsList(pathsBatchHugeMoira)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(tagsList); j++ {
			for k := 0; k < len(g); k++ {
				_ = g[k].MatchByTags(tagsList[j])
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
}
