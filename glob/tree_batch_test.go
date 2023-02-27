package glob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobTree(t *testing.T) {
	tests := []testGlobTree{
		{
			globs: []string{
				"a.*.{cd,b}*.e", "a.*.{b,cd}*.e",
				"a.*.{cd,b}*.df",
				"a.b.b*.e",
				"a.*.{bc,d}*.e",
				"*.{bc,d}*.e",
				"*{cd,b}*.df",
				"a.b.b*.{c,bc}",
				"a.b.b*{c,bc}",
				// optional list (can be empty)
				"bc.{,cd,b}*.df",
				"bc.{,cd,b}",
				"bc.*{,cd,b}",
				"bc.*{,cd,b}*",
			},
			want: &GlobTree{
				Root: &items.TreeItem{
					Childs: []*items.TreeItem{
						{
							NodeItem: items.NodeItem{Node: "a.", Item: items.NewString("a.")},
							Childs: []*items.TreeItem{
								{
									NodeItem: items.NodeItem{Node: "*", Item: items.Star(0)},
									Childs: []*items.TreeItem{
										{
											NodeItem: items.NodeItem{Node: ".", Item: items.Byte('.')},
											Childs: []*items.TreeItem{
												{
													NodeItem: items.NodeItem{Node: "{b,cd}",
														Item: items.NewItemList([]string{"b", "cd"}),
													},
													Childs: []*items.TreeItem{
														{
															NodeItem: items.NodeItem{
																Node: "*", Item: items.Star(0),
															},
															Childs: []*items.TreeItem{
																{
																	NodeItem: items.NodeItem{
																		Node: ".e",
																		Item: items.NewString(".e"),
																	},
																	Terminated: []string{
																		"a.*.{b,cd}*.e",
																	},
																	TermIndex: []int{0},
																},
																{
																	NodeItem: items.NodeItem{
																		Node: ".df",
																		Item: items.NewString(".df"),
																	},
																	Terminated: []string{"a.*.{b,cd}*.df"},
																	TermIndex:  []int{2},
																},
															},
														},
													},
												},
												{
													NodeItem: items.NodeItem{Node: "{bc,d}",
														Item: items.NewItemList([]string{"bc", "d"}),
													},
													Childs: []*items.TreeItem{
														{
															NodeItem: items.NodeItem{
																Node: "*", Item: items.Star(0),
															},
															Childs: []*items.TreeItem{
																{
																	NodeItem: items.NodeItem{
																		Node: ".e",
																		Item: items.NewString(".e"),
																	},
																	Terminated: []string{
																		"a.*.{bc,d}*.e",
																	},
																	TermIndex: []int{4},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{
							NodeItem: items.NodeItem{Node: "a.b.b", Item: items.NewString("a.b.b")},
							Childs: []*items.TreeItem{
								{
									NodeItem: items.NodeItem{Node: "*", Item: items.Star(0)},
									Childs: []*items.TreeItem{
										{
											NodeItem:   items.NodeItem{Node: ".e", Item: items.NewString(".e")},
											Terminated: []string{"a.b.b*.e"},
											TermIndex:  []int{3},
										},
										{
											NodeItem: items.NodeItem{Node: ".", Item: items.Byte('.')},
											Childs: []*items.TreeItem{
												{
													NodeItem: items.NodeItem{
														Node: "{bc,c}",
														Item: items.NewItemList([]string{"bc", "c"}),
													},
													Terminated: []string{"a.b.b*.{bc,c}"},
													TermIndex:  []int{7},
												},
											},
										},
										{
											NodeItem: items.NodeItem{
												Node: "{bc,c}",
												Item: items.NewItemList([]string{"bc", "c"}),
											},
											Terminated: []string{"a.b.b*{bc,c}"},
											TermIndex:  []int{8},
										},
									},
								},
							},
						},
						{
							NodeItem: items.NodeItem{Node: "*", Item: items.Star(0)},
							Childs: []*items.TreeItem{
								{
									NodeItem: items.NodeItem{Node: ".", Item: items.Byte('.')},
									Childs: []*items.TreeItem{
										{
											NodeItem: items.NodeItem{
												Node: "{bc,d}",
												Item: items.NewItemList([]string{"bc", "d"}),
											},
											Childs: []*items.TreeItem{
												{
													NodeItem: items.NodeItem{Node: "*", Item: items.Star(0)},
													Childs: []*items.TreeItem{
														{
															NodeItem: items.NodeItem{
																Node: ".e", Item: items.NewString(".e"),
															},
															Terminated: []string{"*.{bc,d}*.e"},
															TermIndex:  []int{5},
														},
													},
												},
											},
										},
									},
								},
								{
									NodeItem: items.NodeItem{
										Node: "{b,cd}",
										Item: items.NewItemList([]string{"b", "cd"}),
									},
									Childs: []*items.TreeItem{
										{
											NodeItem: items.NodeItem{Node: "*", Item: items.Star(0)},
											Childs: []*items.TreeItem{
												{
													NodeItem: items.NodeItem{
														Node: ".df", Item: items.NewString(".df"),
													},
													Terminated: []string{"*{b,cd}*.df"},
													TermIndex:  []int{6},
												},
											},
										},
									},
								},
							},
						},
						{
							NodeItem: items.NodeItem{Node: "bc.", Item: items.NewString("bc.")},
							Childs: []*items.TreeItem{
								{
									NodeItem: items.NodeItem{
										Node: "{,b,cd}",
										Item: items.NewItemList([]string{"", "b", "cd"}),
									},
									Terminated: []string{"bc.{,b,cd}"},
									TermIndex:  []int{10},
									Childs: []*items.TreeItem{
										{
											NodeItem: items.NodeItem{Node: "*", Item: items.Star(0)},
											Childs: []*items.TreeItem{
												{
													NodeItem: items.NodeItem{
														Node: ".df", Item: items.NewString(".df"),
													},
													Terminated: []string{"bc.{,b,cd}*.df"},
													TermIndex:  []int{9},
												},
											},
										},
									},
								},
								{
									NodeItem: items.NodeItem{Node: "*", Item: items.Star(0)},
									Childs: []*items.TreeItem{
										{
											NodeItem: items.NodeItem{
												Node: "{,b,cd}",
												Item: items.NewItemList([]string{"", "b", "cd"}),
											},
											Terminated: []string{"bc.*{,b,cd}"},
											TermIndex:  []int{11},
											Childs: []*items.TreeItem{
												{
													NodeItem:   items.NodeItem{Node: "*", Item: items.Star(0)},
													Terminated: []string{"bc.*{,b,cd}*"},
													TermIndex:  []int{12},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Globs: map[string]int{
					"a.*.{b,cd}*.e":  0,
					"a.*.{b,cd}*.df": 2,
					"a.b.b*.e":       3,
					"a.*.{bc,d}*.e":  4,
					"*.{bc,d}*.e":    5,
					"*{b,cd}*.df":    6,
					"a.b.b*.{bc,c}":  7,
					"a.b.b*{bc,c}":   8,
					"bc.{,b,cd}*.df": 9,
					"bc.{,b,cd}":     10,
					"bc.*{,b,cd}":    11,
					"bc.*{,b,cd}*":   12,
				},
				GlobsIndex: map[int]string{
					0:  "a.*.{b,cd}*.e",
					2:  "a.*.{b,cd}*.df",
					3:  "a.b.b*.e",
					4:  "a.*.{bc,d}*.e",
					5:  "*.{bc,d}*.e",
					6:  "*{b,cd}*.df",
					7:  "a.b.b*.{bc,c}",
					8:  "a.b.b*{bc,c}",
					9:  "bc.{,b,cd}*.df",
					10: "bc.{,b,cd}",
					11: "bc.*{,b,cd}",
					12: "bc.*{,b,cd}*",
				},
			},
			match: map[string][]string{
				"a.b.bce.e":  {"a.*.{b,cd}*.e", "a.b.b*.e", "a.*.{bc,d}*.e", "*.{bc,d}*.e"},
				"a.b.cd.e":   {"a.*.{b,cd}*.e"},
				"a.D.bce.e":  {"a.*.{b,cd}*.e", "a.*.{bc,d}*.e", "*.{bc,d}*.e"},
				"a.D.bcd.df": {"a.*.{b,cd}*.df", "*{b,cd}*.df", "*{b,cd}*.df"},
				"a.b.b.bc":   {"a.b.b*.{bc,c}", "a.b.b*{bc,c}", "a.b.b*{bc,c}"},
				"a.b.bed.bc": {"a.b.b*.{bc,c}", "a.b.b*{bc,c}", "a.b.b*{bc,c}"},
				// empty list
				"bc..df": {
					"*{b,cd}*.df", "bc.{,b,cd}*.df",
					"bc.*{,b,cd}",
				},
				"bc.cd.df": {
					"*{b,cd}*.df", "*{b,cd}*.df", "bc.{,b,cd}*.df", "bc.{,b,cd}*.df",
					"bc.*{,b,cd}", "bc.*{,b,cd}*",
				},
				"bc.": {"bc.{,b,cd}", "bc.*{,b,cd}", "bc.*{,b,cd}*"},
				"bcd": {},
			},
		},
	}
	for n, tt := range tests {
		runTestGlobTree(t, n, tt)
	}
}

var (
	globsBatch = []string{
		"a.*.{cd,b}*.e", "a.*.{b,cd}*.e",
		"a.*.{cd,b}*.df",
		"a.b.b*.e",
		"a.*.{bc,d}*.e",
		"*.{bc,d}*.e",
		"*{cd,b}*.df",
		"a.b.b*.{c,bc}",
		"a.b.b*{c,bc}",
		// optional list (can be empty)
		"bc.{,cd,b}*.df",
		"bc.{,cd,b}",
		"bc.*{,cd,b}",
		"bc.*{,cd,b}*",
	}
	stringsBatch = []string{
		"a.b.bce.e",
		"a.b.cd.e",
		"a.D.bce.e",
		"a.D.bcd.df",
		"a.b.b.bc",
		"a.b.bed.bc",
		"bc..df",
		"bc.cd.df",
		"bc.",
		"bcd",
	}
	gBatch = parseGlobs(globsBatch)
)

func Benchmark_Batch_GlobTree_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtree := NewTree()
		for i, g := range globsBatch {
			_, _, err := gtree.AddGlob(g, i)

			if err != nil && err != ErrGlobExist {
				b.Fatalf("GlobTree.Add(%q) error = %v", g, err)
			}
		}
	}
}

func Benchmark_Batch_GlobTree_Add_Cached(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gtree := NewTree()
		for i, g := range gBatch {
			_, _, err := gtree.Add(g, i)

			if err != nil && err != ErrGlobExist {
				b.Fatalf("GlobTree.Add(%q) error = %v", g, err)
			}
		}
	}
}

func Benchmark_Batch_GlobTree(b *testing.B) {
	gtree := NewTree()
	for i, g := range gBatch {
		_, _, err := gtree.Add(g, i)

		if err != nil && err != ErrGlobExist {
			b.Fatalf("GlobTree.Add(%q) error = %v", g, err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range stringsBatch {
			var (
				globs []string
				index []int
			)
			first := -1
			_ = gtree.Match(s, &globs, &index, &first)
		}
	}
}

func Benchmark_Batch_GlobTree_Prealloc(b *testing.B) {
	gtree := NewTree()
	for i, g := range gBatch {
		_, _, err := gtree.Add(g, i)

		if err != nil && err != ErrGlobExist {
			b.Fatalf("GlobTree.Add(%q) error = %v", g, err)
		}
	}

	var (
		globs []string
		index []int
		first int
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range stringsBatch {
			first = -1
			globs = globs[:0]
			index = index[:0]
			_ = gtree.Match(s, &globs, &index, &first)
		}
	}
}

func Benchmark_Batch_Glob_Parse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = parseGlobs(globsBatch)
	}
}

func Benchmark_Batch_Glob_Prealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range stringsBatch {
			for _, g := range gBatch {
				_ = g.Match(s)
			}
		}
	}
}
