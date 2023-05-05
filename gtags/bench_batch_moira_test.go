package gtags

import (
	"math/rand"
	"testing"
	"time"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func BenchmarkBatchHuge_Moira_GTag_Precompiled(b *testing.B) {
	g := taggedTermListList(queriesBatchHugeMoira)
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, b.N)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := pathsBatchHugeMoira[rand.Intn(len(pathsBatchHugeMoira))]
		tags, _ := PathTags(path)
		for k := 0; k < len(g); k++ {
			_ = g[k].MatchByTags(tags)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Moira_Tree_ByTags_Precompiled(b *testing.B) {
	w := NewTree()
	for j := 0; j < len(queriesBatchHugeMoira); j++ {
		_, _, err := w.Add(queriesBatchHugeMoira[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, b.N)
	var store items.AllStore
	store.Init()
	store.Grow(1)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Init()
		path := pathsBatchHugeMoira[rand.Intn(len(pathsBatchHugeMoira))]
		tags, _ := PathTags(path)
		_ = w.MatchByTags(tags, &store)
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Moira_Tree_ByTagsMap_Precompiled(b *testing.B) {
	w := NewTree()
	for j := 0; j < len(queriesBatchHugeMoira); j++ {
		_, _, err := w.Add(queriesBatchHugeMoira[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, b.N)
	var store items.AllStore
	store.Init()
	store.Grow(1)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Init()
		path := pathsBatchHugeMoira[rand.Intn(len(pathsBatchHugeMoira))]
		tags, _ := PathTagsMap(path)
		_ = w.MatchByTagsMap(tags, &store)
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Moira_Tree_ByTagsMap_PrecompiledB(b *testing.B) {
	w := NewTree()
	for j := 0; j < len(queriesBatchHugeMoira); j++ {
		_, _, err := w.Add(queriesBatchHugeMoira[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, b.N)
	var store items.AllStore
	store.Init()
	store.Grow(1)

	start := time.Now()
	b.ResetTimer()
	tags := make(map[string]string)
	for i := 0; i < b.N; i++ {
		store.Init()
		path := pathsBatchHugeMoira[rand.Intn(len(pathsBatchHugeMoira))]
		_ = PathTagsMapB(path, tags)
		_ = w.MatchByTagsMap(tags, &store)
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Moira_Tree_ByTags_Prealloc(b *testing.B) {
	w := NewTree()
	for j := 0; j < len(queriesBatchHugeMoira); j++ {
		_, _, err := w.Add(queriesBatchHugeMoira[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, b.N)
	tagsList := tagsList(pathsBatchHugeMoira)
	var store items.AllStore
	store.Init()
	store.Grow(1)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Init()
		tags := tagsList[rand.Intn(len(tagsList))]
		_ = w.MatchByTags(tags, &store)
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}

func BenchmarkBatchHuge_Moira_Tree_ByTagsMap_Prealloc(b *testing.B) {
	w := NewTree()
	for j := 0; j < len(queriesBatchHugeMoira); j++ {
		_, _, err := w.Add(queriesBatchHugeMoira[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	pathsBatchHugeMoira := generateTaggedMetrics(termsBatchHugeMoira, b.N)
	tagMapList := tagMapList(pathsBatchHugeMoira)
	var store items.AllStore
	store.Init()
	store.Grow(1)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Init()
		tags := tagMapList[rand.Intn(len(tagMapList))]
		_ = w.MatchByTagsMap(tags, &store)
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N)/d.Seconds(), "match/s")
}
