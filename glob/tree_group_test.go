package glob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobTree_Group(t *testing.T) {
	tests := []testGlobTree{
		{
			globs: []string{
				"{b*,a?cd*,cd[a-z]}bc*c*e",
				"{b*,a?cd*,cd[a-z]}bc*CD*e",
			},
			want: &globTreeStr{
				Root: &TreeItemStr{
					Childs: []*TreeItemStr{
						{
							Node: "e", Reverse: true, Childs: []*TreeItemStr{
								{
									Node: "{a?cd*,b*,cd[a-z]}", Childs: []*TreeItemStr{
										{
											Node: "bc", Childs: []*TreeItemStr{
												{
													Node: "*", Childs: []*TreeItemStr{
														{
															Node: "c", Childs: []*TreeItemStr{
																{
																	Node: "*", Childs: []*TreeItemStr{},
																	Terminated: items.Terminated{
																		Terminate: true, Index: 0,
																		Query: "{a?cd*,b*,cd[a-z]}bc*c*e",
																	},
																},
															},
														},
														{
															Node: "CD", Childs: []*TreeItemStr{
																{
																	Node: "*", Childs: []*TreeItemStr{},
																	Terminated: items.Terminated{
																		Terminate: true, Index: 1,
																		Query: "{a?cd*,b*,cd[a-z]}bc*CD*e",
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
				},
				Globs: map[string]int{
					"{a?cd*,b*,cd[a-z]}bc*c*e": 0, "{b*,a?cd*,cd[a-z]}bc*c*e": 0,
					"{a?cd*,b*,cd[a-z]}bc*CD*e": 1, "{b*,a?cd*,cd[a-z]}bc*CD*e": 1,
				},
				GlobsIndex: map[int]string{0: "{a?cd*,b*,cd[a-z]}bc*c*e", 1: "{a?cd*,b*,cd[a-z]}bc*CD*e"},
			},
			match: map[string][]string{
				"aZcdbcce":      {"{a?cd*,b*,cd[a-z]}bc*c*e"},
				"aЯcdbcce":      {"{a?cd*,b*,cd[a-z]}bc*c*e"},
				"aЮcdQAbcZWcIe": {"{a?cd*,b*,cd[a-z]}bc*c*e"},
				"aZcdQAbcZWcIe": {"{a?cd*,b*,cd[a-z]}bc*c*e"},
				"":              nil, "acdqbcZcIe": nil, "abCDbcZIce": nil,
				"ЙabCDbcZIce": nil, "aZcdbcc": nil, "aZcdcce": nil, "aZcdQAbcZWIe": nil,
			},
		},
		{
			globs: []string{
				"*{b*,a?cd*,cd[a-z]}bc*c*e",
				"*{b*,a?cd*,cd[a-z]}bc*CD*e",
				"*{b*,a?cd*,cd[a-z]}bc*cd*e",
			},
			want: &globTreeStr{
				Root: &TreeItemStr{
					Childs: []*TreeItemStr{
						{
							Node: "e", Reverse: true, Childs: []*TreeItemStr{
								{
									Node: "*", Childs: []*TreeItemStr{
										{
											Node: "{a?cd*,b*,cd[a-z]}", Childs: []*TreeItemStr{
												{
													Node: "bc", Childs: []*TreeItemStr{
														{
															Node: "*", Childs: []*TreeItemStr{
																{
																	Node: "c", Childs: []*TreeItemStr{
																		{
																			Node: "*", Childs: []*TreeItemStr{},
																			Terminated: items.Terminated{
																				Terminate: true, Index: 0,
																				Query: "*{a?cd*,b*,cd[a-z]}bc*c*e",
																			},
																		},
																	},
																},
																{
																	Node: "CD", Childs: []*TreeItemStr{
																		{
																			Node: "*", Childs: []*TreeItemStr{},
																			Terminated: items.Terminated{
																				Terminate: true, Index: 1,
																				Query: "*{a?cd*,b*,cd[a-z]}bc*CD*e",
																			},
																		},
																	},
																},
																{
																	Node: "cd", Childs: []*TreeItemStr{
																		{
																			Node: "*", Childs: []*TreeItemStr{},
																			Terminated: items.Terminated{
																				Terminate: true, Index: 2,
																				Query: "*{a?cd*,b*,cd[a-z]}bc*cd*e",
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
						},
					},
				},
				Globs: map[string]int{
					"*{a?cd*,b*,cd[a-z]}bc*c*e": 0, "*{b*,a?cd*,cd[a-z]}bc*c*e": 0,
					"*{a?cd*,b*,cd[a-z]}bc*CD*e": 1, "*{b*,a?cd*,cd[a-z]}bc*CD*e": 1,
					"*{a?cd*,b*,cd[a-z]}bc*cd*e": 2, "*{b*,a?cd*,cd[a-z]}bc*cd*e": 2,
				},
				GlobsIndex: map[int]string{
					0: "*{a?cd*,b*,cd[a-z]}bc*c*e", 1: "*{a?cd*,b*,cd[a-z]}bc*CD*e",
					2: "*{a?cd*,b*,cd[a-z]}bc*cd*e",
				},
			},
			match: map[string][]string{
				"aZcdbcce":      {"*{a?cd*,b*,cd[a-z]}bc*c*e"},
				"aЯcdbcce":      {"*{a?cd*,b*,cd[a-z]}bc*c*e"},
				"aЮcdQAbcZWcIe": {"*{a?cd*,b*,cd[a-z]}bc*c*e"},
				"aZcdQAbcZWcIe": {"*{a?cd*,b*,cd[a-z]}bc*c*e"},
				"abCDbcZIce":    {"*{a?cd*,b*,cd[a-z]}bc*c*e"},
				"ЙabCDbcZIce":   {"*{a?cd*,b*,cd[a-z]}bc*c*e"},
				"acdqbcZcIe":    {"*{a?cd*,b*,cd[a-z]}bc*c*e"},
				"acdqbcZCDIe":   {"*{a?cd*,b*,cd[a-z]}bc*CD*e"},
				"acdqbcZcdIe":   {"*{a?cd*,b*,cd[a-z]}bc*c*e", "*{a?cd*,b*,cd[a-z]}bc*cd*e"},
				"":              nil,
				"aZcdbcc":       nil, "aZcdcce": nil, "aZcdQAbcZWIe": nil,
			},
		},
	}
	for n, tt := range tests {
		runTestGlobTree(t, n, tt)
	}
}
