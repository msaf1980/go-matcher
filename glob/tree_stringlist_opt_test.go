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
				Root: &TreeItemStr{
					Childs: []*TreeItemStr{
						{
							Node: "bc.", Childs: []*TreeItemStr{
								{
									Node: "*", Childs: []*TreeItemStr{
										{
											Node: "{,b,cd}",
											Terminated: items.Terminated{
												Terminate: true, Index: 1,
												Query: "bc.*{,b,cd}",
											},
											Childs: []*TreeItemStr{
												{
													Node: "*", Childs: []*TreeItemStr{},
													Terminated: items.Terminated{
														Terminate: true, Index: 0,
														Query: "bc.*{,b,cd}*",
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
