package gglob

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/msaf1980/go-matcher/glob"
	"github.com/msaf1980/go-matcher/pkg/items"
)

type GTreeItemStr struct {
	Node string

	Terminated items.Terminated

	// TODO: may be some ordered tree for complete string nodes search speedup (on large set) ?
	ChildsMap map[string]*GTreeItemStr `json:"childs_map"`
	Childs    []*GTreeItemStr          `json:"childs"` // next possible parts slice
}

func StringGTreeItem(treeItem *GTreeItem) *GTreeItemStr {
	var node string
	if treeItem.Item != nil {
		node = treeItem.Item.Node
	}
	treeItemStr := &GTreeItemStr{
		Node:       node,
		Terminated: treeItem.Terminated,
	}

	if treeItem.Childs != nil {
		treeItemStr.Childs = make([]*GTreeItemStr, 0, len(treeItem.Childs))
		for _, child := range treeItem.Childs {
			treeItemStr.Childs = append(treeItemStr.Childs, StringGTreeItem(child))
		}
	}
	if treeItem.ChildsMap != nil {
		treeItemStr.ChildsMap = make(map[string]*GTreeItemStr)
		for k, child := range treeItem.ChildsMap {
			treeItemStr.ChildsMap[k] = StringGTreeItem(child)
		}
	}
	return treeItemStr
}

type globGTreeStr struct {
	Root       map[int]*GTreeItemStr
	Globs      map[string]int
	GlobsIndex map[int]string
}

type verify struct {
	glob  string
	index int
}

func mergeVerify(globs []string, index []int) []verify {
	if len(globs) != len(index) {
		return nil
	}
	v := make([]verify, len(globs))
	for i := 0; i < len(globs); i++ {
		v[i].glob = globs[i]
		v[i].index = index[i]
	}
	return v
}

type globTreeStr struct {
	Root       map[int]*GTreeItemStr
	Globs      map[string]int
	GlobsIndex map[int]string
}

type testGGlobTree struct {
	globs   []string
	skipCmp bool // don't compare glob tree, only glob maps
	want    *globTreeStr
	match   map[string][]string
}

func runTestGGlobTree(t *testing.T, n int, tt testGGlobTree) {
	t.Run(fmt.Sprintf("%d#%#v", n, tt.globs), func(t *testing.T) {
		gtree := NewTree()
		for i, g := range tt.globs {
			_, _, err := gtree.Add(g, i)

			if err != nil && err != glob.ErrGlobExist {
				t.Fatalf("GlobTree.Add(%q) error = %v", g, err)
			}
		}

		var globTree *globTreeStr
		if tt.skipCmp {
			globTree = &globTreeStr{
				Globs:      gtree.Globs,
				GlobsIndex: gtree.GlobsIndex,
			}
		} else {
			globTree = &globTreeStr{
				Root:       make(map[int]*GTreeItemStr),
				Globs:      gtree.Globs,
				GlobsIndex: gtree.GlobsIndex,
			}
			for n, t := range gtree.Root {
				globTree.Root[n] = StringGTreeItem(t)
			}
		}
		if !reflect.DeepEqual(globTree, tt.want) {
			t.Fatalf("GlobTree(%#v) = %s", tt.globs, cmp.Diff(tt.want, globTree))
		}

		verifyGGlobTree(t, tt.globs, tt.match, gtree)
	})
}

func verifyGGlobTree(t *testing.T, inGlobs []string, match map[string][]string, gtree *GGlobTree) {
	for path, wantGlobs := range match {
		t.Run("#path="+path, func(t *testing.T) {
			var store items.AllStore
			store.Init()
			store.Grow(1)
			matched := gtree.Match(path, &store)

			verify := mergeVerify(store.S.S, store.Index.N)

			sort.Strings(store.S.S)
			sort.Strings(wantGlobs)
			sort.Ints(store.Index.N)

			if !reflect.DeepEqual(wantGlobs, store.S.S) {
				t.Fatalf("GlobTree(%#v).Match(%q) globs = %s", inGlobs, path, cmp.Diff(wantGlobs, store.S.S))
			}

			if matched != len(store.S.S) || len(store.S.S) != len(store.Index.N) {
				t.Fatalf("GlobTree(%#v).Match(%q) = %d, want %d, index = %d", inGlobs, path, matched, len(store.S.S), len(store.Index.N))
			}

			for _, v := range verify {
				if v.glob != gtree.GlobsIndex[v.index] {
					t.Errorf("GlobTree(%#v).Match(%q) index = %d glob = %s, want %s",
						inGlobs, path, v.index, gtree.GlobsIndex[v.index], v.glob)
				}
			}

			if len(store.Index.N) > 0 {
				if store.Min.Min != store.Index.N[0] {
					t.Errorf("GlobTree(%#v).Match(%q) first index = %d, want %d",
						inGlobs, path, store.Min.Min, store.Index.N[0])
				}
			}

			store.Init()
			parts := PathSplit(path)
			matched = gtree.MatchByParts(parts, &store)
			if !reflect.DeepEqual(wantGlobs, store.S.S) {
				t.Fatalf("GlobTree(%#v).MatchByParts(%q) globs = %s", inGlobs, path, cmp.Diff(wantGlobs, store.S.S))
			}

			if matched != len(store.S.S) || len(store.S.S) != len(store.Index.N) {
				t.Fatalf("GlobTree(%#v).MatchByParts(%q) = %d, want %d, index = %d", inGlobs, path, matched, len(store.S.S), len(store.Index.N))
			}

		})
	}
}

func TestGGlobTree(t *testing.T) {
	tests := []testGGlobTree{
		{
			globs: []string{
				"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].DownEndpointCount",
				"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].UpStatus",
				"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.UpStatus",
			},
			want: &globTreeStr{
				Root: map[int]*GTreeItemStr{
					5: {
						ChildsMap: map[string]*GTreeItemStr{
							"DB": {
								Node: "DB", Childs: []*GTreeItemStr{
									{
										Node: "*", Childs: []*GTreeItemStr{
											{
												Node: "{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}",
												Childs: []*GTreeItemStr{
													{
														Node: "*[0-8]",
														ChildsMap: map[string]*GTreeItemStr{
															"DownEndpointCount": {
																Node: "DownEndpointCount",
																Terminated: items.Terminated{
																	Terminate: true,
																	Query:     "DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].DownEndpointCount",
																},
															},
															"UpStatus": {
																Node: "UpStatus",
																Terminated: items.Terminated{
																	Terminate: true, Index: 1,
																	Query: "DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].UpStatus",
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
					4: {
						ChildsMap: map[string]*GTreeItemStr{
							"DB": {
								Node: "DB", Childs: []*GTreeItemStr{
									{
										Node: "*", Childs: []*GTreeItemStr{
											{
												Node: "{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}",
												ChildsMap: map[string]*GTreeItemStr{
													"UpStatus": {
														Node: "UpStatus",
														Terminated: items.Terminated{
															Terminate: true, Index: 2,
															Query: "DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.UpStatus",
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
					"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].DownEndpointCount": 0,
					"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].UpStatus":          1,
					"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.UpStatus":                 2,
				},
				GlobsIndex: map[int]string{
					0: "DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].DownEndpointCount",
					1: "DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].UpStatus",
					2: "DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.UpStatus",
				},
			},
			match: map[string][]string{
				"DB.Sales.BalanceCluster.node1.DownEndpointCount": {
					"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].DownEndpointCount",
				},
				"DB.Back.WebCluster.node2.DownEndpointCount": {
					"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].DownEndpointCount",
				},
				"DB.Sales.BalanceCluster.node8.DownEndpointCount": {
					"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].DownEndpointCount",
				},
				"DB.Sales.BalanceCluster.node8.UpStatus": {
					"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.*[0-8].UpStatus",
				},
				"DB.Sales.BalanceCluster.UpStatus": {
					"DB.*.{BalanceCluster,BalanceStaging,CoreCluster,EventsCluster,SalesCluster,UpProduction,UpTesting,WebCluster}.UpStatus",
				},
				"DB.Sales.BalanceCluster..UpStatus":                 {},
				"DB.Back.WebCluster.node2.UpEndpointCount":          {},
				"DB.Back.DBCluster.node2.DownEndpointCount":         {},
				"DB.Sales.BalanceCluster.node1.DownEndpointCount.2": {},
				"DBA.Back.WebCluster.node2.DownEndpointCount":       {},
				"DB.Sales.BalanceCluster.node9.DownEndpointCount":   {},
			},
		},
	}
	for n, tt := range tests {
		runTestGGlobTree(t, n, tt)
	}
}

func parseGGlobs(globs []string) (g []*GGlob) {
	g = make([]*GGlob, len(globs))
	for i := 0; i < len(globs); i++ {
		g[i] = ParseMust(globs[i])
	}

	return
}
