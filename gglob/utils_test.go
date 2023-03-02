package gglob

import "testing"

func Benchmark_PathSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = PathSplit(stringBenchASCII)
	}
}

func Benchmark_PathSplitB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parts := make([]string, 5)
		_ = PathSplitB(stringBenchASCII, &parts)
	}
}

func Benchmark_PathSplitB_Prealloc(b *testing.B) {
	parts := make([]string, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PathSplitB(stringBenchASCII, &parts)
	}
}
