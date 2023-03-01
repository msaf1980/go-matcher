package glob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobTree_StringList_Opt(t *testing.T) {
	tests := []testGlobTree{
		{
			globs: []string{
				"bc.*{,b,cd}*",
				"bc.*{,b,cd}",
			},
			want: &globTreeStr{
				Root: &items.TreeItemStr{
					Childs: []*items.TreeItemStr{
						{
							Node: "bc.", Childs: []*items.TreeItemStr{
								{
									Node: "*", Childs: []*items.TreeItemStr{
										{
											Node:       "{,b,cd}",
											Terminated: []string{"bc.*{,b,cd}"},
											TermIndex:  []int{1},
											Childs: []*items.TreeItemStr{
												{
													Node: "*", Childs: []*items.TreeItemStr{},
													Terminated: []string{"bc.*{,b,cd}*"},
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
