package glob

import (
	"testing"
)

func TestGlobTree_StringList(t *testing.T) {
	tests := []testGlobTree{
		{
			globs: []string{
				"a.*.{cd,b}*.e", "a.*.{b,cd,b}*.e",
			},
			want: &globTreeStr{
				Root: &TreeItemStr{
					Childs: []*TreeItemStr{
						{
							Node: ".e", Reverse: true, Childs: []*TreeItemStr{
								{
									Node: "a.", Childs: []*TreeItemStr{
										{
											Node: "*", Childs: []*TreeItemStr{
												{
													Node: ".", Childs: []*TreeItemStr{
														{
															Node: "{b,cd}", Childs: []*TreeItemStr{
																{
																	Node: "*", Childs: []*TreeItemStr{},
																	Terminated: "a.*.{b,cd}*.e", TermIndex: 0, Term: true,
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
				Globs:      map[string]int{"a.*.{b,cd}*.e": 0, "a.*.{cd,b}*.e": 0},
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
				Root: &TreeItemStr{
					Childs: []*TreeItemStr{
						{
							Node: ".df", Reverse: true, Childs: []*TreeItemStr{
								{
									Node: "*", Childs: []*TreeItemStr{
										{
											Node: "{b,cd}", Childs: []*TreeItemStr{
												{
													Node: "*", Childs: []*TreeItemStr{},
													Terminated: "*{b,cd}*.df", TermIndex: 0, Term: true,
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
				Root: &TreeItemStr{
					Childs: []*TreeItemStr{
						{
							Node: "a.b", Childs: []*TreeItemStr{
								{
									Node: "*", Childs: []*TreeItemStr{
										{
											Node: ".", Childs: []*TreeItemStr{
												{
													Node: "{bc,c}", Childs: []*TreeItemStr{},
													Terminated: "a.b*.{bc,c}", TermIndex: 0, Term: true,
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
				"a.b.b.bed.cd": nil, "a.b.b.bed.bcd": nil, "a.b.d.bc": nil,
			},
		},
	}
	for n, tt := range tests {
		runTestGlobTree(t, n, tt)
	}
}
