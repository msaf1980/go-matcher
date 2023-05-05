package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

func BenchmarkGready_List_Tree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		_, _, err := w.Add(globsBatchList[0], 1)
		if err != nil {
			b.Fatal(err)
		}
		var store items.AllStore
		store.Init()
		_ = w.Match(pathsBatchList[0], &store)
		if len(store.S.S) != 1 {
			b.Fatal(store.S.S)
		}
	}
}

func _BenchmarkGready_List_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globsBatchList[0])
		if !w.MatchString(pathsBatchList[0]) {
			b.Fatal(pathsBatchList[0])
		}
	}
}

func BenchmarkGready_List_Tree_Precompiled(b *testing.B) {
	g := ParseMust(globsBatchList[0])

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewTree()
		w.AddGlob(g, 1)

		var store items.AllStore
		store.Init()
		store.Grow(4)
		_ = w.Match(pathsBatchList[0], &store)
		if len(store.S.S) != 1 {
			b.Fatal(store.S.S)
		}
	}
}

func BenchmarkGready_List_Tree_Precompiled2(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globsBatchList[0], 1)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var store items.AllStore
		store.Init()
		_ = w.Match(pathsBatchList[0], &store)
		if len(store.S.S) != 1 {
			b.Fatal(store.S.S)
		}
	}
}

func BenchmarkGready_List_Tree_Prealloc(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globsBatchList[0], 1)
	if err != nil {
		b.Fatal(err)
	}
	var store items.AllStore
	store.Init()
	store.Grow(4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Init()
		_ = w.Match(pathsBatchList[0], &store)
		if len(store.S.S) != 1 {
			b.Fatal(store.S.S)
		}
	}
}

func _BenchmarkGready_List_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globsBatchList[0])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !w.MatchString(pathsBatchList[0]) {
			b.Fatalf("\n%s\n%s\n%s", w.String(), globsBatchList[0], pathsBatchList[0])
		}
	}
}
