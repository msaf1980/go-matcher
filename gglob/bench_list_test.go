package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

func BenchmarkGready_List_Tree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		_, _, err := w.Add(targetsBatchList[0], 1)
		if err != nil {
			b.Fatal(err)
		}
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathsBatchList[0], &globs, nil, &first)
		if len(globs) != 1 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_List_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(targetsBatchList[0])
		if !w.MatchString(pathsBatchList[0]) {
			b.Fatal(pathsBatchList[0])
		}
	}
}

func BenchmarkGready_List_Tree_Precompiled(b *testing.B) {
	g := ParseMust(targetsBatchList[0])

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewTree()
		w.AddGlob(g, 1)

		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathsBatchList[0], &globs, nil, &first)
		if len(globs) != 1 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_List_Tree_Precompiled2(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(targetsBatchList[0], 1)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathsBatchList[0], &globs, nil, &first)
		if len(globs) != 1 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_List_Tree_Prealloc(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(targetsBatchList[0], 1)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		first.Init()
		_ = w.Match(pathsBatchList[0], &globs, nil, &first)
		if len(globs) != 1 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_List_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(targetsBatchList[0])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathsBatchList[0]) {
			b.Fatal(pathsBatchList[0])
		}
	}
}
