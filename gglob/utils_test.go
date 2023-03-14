package gglob

import "testing"

func BenchmarkPathSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = PathSplit(stringBenchASCII)
	}
}

func BenchmarkPathSplitB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parts := make([]string, 5)
		_ = PathSplitB(stringBenchASCII, &parts)
	}
}

func BenchmarkPathSplitB_Prealloc(b *testing.B) {
	parts := make([]string, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PathSplitB(stringBenchASCII, &parts)
	}
}
