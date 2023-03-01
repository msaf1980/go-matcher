package glob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobTree_StringList(t *testing.T) {
	tests := []testGlobTree{
		{
			globs: []string{
				"a.*.{cd,b}*.e", "a.*.{b,cd,b}*.e",
			},
			want: &globTreeStr{
				Root: &items.TreeItemStr{
					Childs: []*items.TreeItemStr{
						{
							Node: ".e", Reverse: true, Childs: []*items.TreeItemStr{
								{
									Node: "a.", Childs: []*items.TreeItemStr{
										{
											Node: "*", Childs: []*items.TreeItemStr{
												{
													Node: ".", Childs: []*items.TreeItemStr{
														{
															Node: "{b,cd}", Childs: []*items.TreeItemStr{
																{
																	Node: "*", Childs: []*items.TreeItemStr{},
																	Terminated: []string{"a.*.{b,cd}*.e"},
																	TermIndex:  []int{0},
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
					},
				},
				Globs:      map[string]int{"a.*.{b,cd}*.e": 0},
				GlobsIndex: map[int]string{0: "a.*.{b,cd}*.e"},
			},
			match: map[string][]string{
				"a.b.bce.e":  {"a.*.{b,cd}*.e"},
				"a.b.cd.e":   {"a.*.{b,cd}*.e"},
				"a.D.bce.e":  {"a.*.{b,cd}*.e"},
				"a.D.dbce.e": nil,
				"a.D.dcd.e":  nil,
			},
		},
		{
			globs: []string{
				"*{b,cd}*.df",
			},
			want: &globTreeStr{
				Root: &items.TreeItemStr{
					Childs: []*items.TreeItemStr{
						{
							Node: ".df", Reverse: true, Childs: []*items.TreeItemStr{
								{
									Node: "*", Childs: []*items.TreeItemStr{
										{
											Node: "{b,cd}", Childs: []*items.TreeItemStr{
												{
													Node: "*", Childs: []*items.TreeItemStr{},
													Terminated: []string{"*{b,cd}*.df"},
													TermIndex:  []int{0},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Globs:      map[string]int{"*{b,cd}*.df": 0},
				GlobsIndex: map[int]string{0: "*{b,cd}*.df"},
			},
			match: map[string][]string{
				"a.D.bcd.df": {"*{b,cd}*.df", "*{b,cd}*.df"},
				"a.D.b.df":   {"*{b,cd}*.df"},
			},
		},
		{
			globs: []string{"a.b*.{bc,c}"},
			want: &globTreeStr{
				Root: &items.TreeItemStr{
					Childs: []*items.TreeItemStr{
						{
							Node: "a.b", Childs: []*items.TreeItemStr{
								{
									Node: "*", Childs: []*items.TreeItemStr{
										{
											Node: ".", Childs: []*items.TreeItemStr{
												{
													Node: "{bc,c}", Childs: []*items.TreeItemStr{},
													Terminated: []string{"a.b*.{bc,c}"},
													TermIndex:  []int{0},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Globs:      map[string]int{"a.b*.{bc,c}": 0},
				GlobsIndex: map[int]string{0: "a.b*.{bc,c}"},
			},
			match: map[string][]string{
				"a.b.bed.bc": {"a.b*.{bc,c}"},
			},
		},
		{
			globs: []string{
				"bc.*{,b,cd}",
				"bc.*{,b,cd}*",
			},
			skipCmp: true,
			want: &globTreeStr{
				Globs:      map[string]int{"bc.*{,b,cd}": 0, "bc.*{,b,cd}*": 1},
				GlobsIndex: map[int]string{0: "bc.*{,b,cd}", 1: "bc.*{,b,cd}*"},
			},
			match: map[string][]string{
				"bc.cd.df": {"bc.*{,b,cd}", "bc.*{,b,cd}*", "bc.*{,b,cd}*"}, // 2 match: (bc.*) and (bc.cd*)
			},
		},
		{
			globs:   []string{"a.b.b*.{bc,c}"},
			skipCmp: true,
			want: &globTreeStr{
				Globs:      map[string]int{"a.b.b*.{bc,c}": 0},
				GlobsIndex: map[int]string{0: "a.b.b*.{bc,c}"},
			},
			match: map[string][]string{
				"a.b.b.bed.bc": {"a.b.b*.{bc,c}"},
			},
		},
	}
	for n, tt := range tests {
		runTestGlobTree(t, n, tt)
	}
}
