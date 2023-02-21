package gglob

import (
	"strings"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/wildcards"
)

func TestNodeItem_Merge(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: "merge strings #all", globs: []string{"a[a-]Z[Q]"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "aaZQ", Terminated: []string{"a[a-]Z[Q]", "aaZQ"},
								WildcardItems: wildcards.WildcardItems{P: "aaZQ", MinSize: 4, MaxSize: 4},
							},
						},
					},
				},
				Globs: map[string]int{"a[a-]Z[Q]": -1, "aaZQ": -1},
			},
			matchPaths: map[string][]string{"aaZQ": {"a[a-]Z[Q]", "aaZQ"}},
			missPaths:  []string{"", "ab", "aaZQa"},
		},
		{
			name: "merge strings #prefix", globs: []string{"a[a-]Z[Q]*"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "aaZQ*", Terminated: []string{"a[a-]Z[Q]*", "aaZQ*"},
								WildcardItems: wildcards.WildcardItems{
									P: "aaZQ", MinSize: 4, MaxSize: -1,
									Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a[a-]Z[Q]*": -1, "aaZQ*": -1},
			},
			matchPaths: map[string][]string{
				"aaZQ":  {"a[a-]Z[Q]*", "aaZQ*"},
				"aaZQa": {"a[a-]Z[Q]*", "aaZQ*"},
			},
			missPaths: []string{"", "ab", "aaZqa"},
		},
		{
			name: "merge strings #suffix", globs: []string{"a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node:       "aaZQstLT*INN*zaSTltl",
								Terminated: []string{"a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l", "aaZQstLT*INN*zaSTltl"},
								WildcardItems: wildcards.WildcardItems{
									P: "aaZQstLT", Suffix: "zaSTltl", MinSize: 18, MaxSize: -1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemStar{}, wildcards.ItemString("INN"),
										wildcards.ItemStar{},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l": -1, "aaZQstLT*INN*zaSTltl": -1},
			},
			matchPaths: map[string][]string{
				"aaZQstLTINNzaSTltl":  {"a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l", "aaZQstLT*INN*zaSTltl"},
				"aaZQstLTaINNzaSTltl": {"a[a-]Z[Q]st{LT}*I{NN}*[z-][a]ST{lt}l", "aaZQstLT*INN*zaSTltl"},
			},
			missPaths: []string{"", "ab", "aaZqa"},
		},
		{
			name: "merge strings #suffix 2", globs: []string{"a[a-]Z[Q]st{LT}*{NN}I*[z-][a]ST{lt}l"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node:       "aaZQstLT*NNI*zaSTltl",
								Terminated: []string{"a[a-]Z[Q]st{LT}*{NN}I*[z-][a]ST{lt}l", "aaZQstLT*NNI*zaSTltl"},
								WildcardItems: wildcards.WildcardItems{
									P: "aaZQstLT", Suffix: "zaSTltl", MinSize: 18, MaxSize: -1,
									Inners: []wildcards.InnerItem{
										wildcards.ItemStar{}, wildcards.ItemString("NNI"),
										wildcards.ItemStar{},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a[a-]Z[Q]st{LT}*{NN}I*[z-][a]ST{lt}l": -1, "aaZQstLT*NNI*zaSTltl": -1},
			},
			matchPaths: map[string][]string{
				"aaZQstLTNNIzaSTltl":  {"a[a-]Z[Q]st{LT}*{NN}I*[z-][a]ST{lt}l", "aaZQstLT*NNI*zaSTltl"},
				"aaZQstLTaNNIzaSTltl": {"a[a-]Z[Q]st{LT}*{NN}I*[z-][a]ST{lt}l", "aaZQstLT*NNI*zaSTltl"},
			},
			missPaths: []string{"", "ab", "aaZqa"},
		},
		{
			name: "merge strings #suffix 3", globs: []string{"a[a-]Z[Q]*[a]c"},
			wantW: &GlobMatcher{
				Root: map[int]*NodeItem{
					1: {
						Childs: []*NodeItem{
							{
								Node: "aaZQ*ac", Terminated: []string{"a[a-]Z[Q]*[a]c", "aaZQ*ac"},
								WildcardItems: wildcards.WildcardItems{
									P: "aaZQ", Suffix: "ac", MinSize: 6, MaxSize: -1,
									Inners: []wildcards.InnerItem{wildcards.ItemStar{}},
								},
							},
						},
					},
				},
				Globs: map[string]int{"a[a-]Z[Q]*[a]c": -1, "aaZQ*ac": -1},
			},
			matchPaths: map[string][]string{
				"aaZQac":  {"a[a-]Z[Q]*[a]c", "aaZQ*ac"},
				"aaZQaac": {"a[a-]Z[Q]*[a]c", "aaZQ*ac"},
				"aaZQbac": {"a[a-]Z[Q]*[a]c", "aaZQ*ac"},
			},
			missPaths: []string{"", "ab", "aaZqa"},
		},
	}
	for _, tt := range tests {
		runTestGlobMatcher(t, tt)
	}
}

func BenchmarkMergeAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(16)
		_, err := w.Add("a[a-]Z[Q]", &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMergeAll_Prealloc(b *testing.B) {
	var buf strings.Builder
	buf.Grow(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		_, err := w.Add("a[a-]Z[Q]", &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMergePrefix(b *testing.B) {
	var buf strings.Builder
	buf.Grow(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		_, err := w.Add("a[a-]Z[Q]*", &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMergePrefix_Prealloc(b *testing.B) {
	var buf strings.Builder
	buf.Grow(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		_, err := w.Add("a[a-]Z[Q]*", &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMergeSuffix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(16)
		_, err := w.Add("a[a-]Z[Q]*[z-]l", &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMergeSuffix_Prealloc(b *testing.B) {
	var buf strings.Builder
	buf.Grow(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		_, err := w.Add("a[a-]Z[Q]*[z-]l", &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}
