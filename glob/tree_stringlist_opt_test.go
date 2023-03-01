package glob

import (
	"testing"
)

func TestGlobTree_StringList_Opt(t *testing.T) {
	tests := []testGlobTree{
		{
			globs: []string{
				"bc.*{,b,cd}*",
				"bc.*{,b,cd}",
			},
			want: &globTreeStr{
				Root: &TreeItemStr{
					Childs: []*TreeItemStr{
						{
							Node: "bc.", Childs: []*TreeItemStr{
								{
									Node: "*", Childs: []*TreeItemStr{
										{
											Node:       "{,b,cd}",
											Terminated: "bc.*{,b,cd}",
											TermIndex:  1,
											Childs: []*TreeItemStr{
												{
													Node: "*", Childs: []*TreeItemStr{},
													Terminated: "bc.*{,b,cd}*",
													TermIndex:  0,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Globs:      map[string]int{"bc.*{,b,cd}*": 0, "bc.*{,b,cd}": 1},
				GlobsIndex: map[int]string{0: "bc.*{,b,cd}*", 1: "bc.*{,b,cd}"},
			},
			match: map[string][]string{
				"bc..df": {"bc.*{,b,cd}", "bc.*{,b,cd}*"},
				"bc.":    {"bc.*{,b,cd}", "bc.*{,b,cd}*"},
			},
		},
	}
	for n, tt := range tests {
		runTestGlobTree(t, n, tt)
	}
}
