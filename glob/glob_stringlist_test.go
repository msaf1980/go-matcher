package glob

import (
	"strings"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
	"github.com/msaf1980/go-matcher/pkg/utils"
)

func TestGlob_StringList(t *testing.T) {
	tests := []testGlob{
		{
			// TODO: convert to runesranges - but may be empty element
			glob: "{Z,Q}",
			want: &Glob{
				Glob: "{Z,Q}", Node: "{Q,Z}", MinLen: 1, MaxLen: 1,
				Items: []items.Item{
					&items.StringList{
						Vals: []string{"Q", "Z"}, MinSize: 1, MaxSize: 1,
						FirstASCII: utils.MakeASCIISetMust("ZQ"), ASCIIStarted: true,
					},
				},
			},
			verify: `^(Z|Q)$`,
			match:  []string{"Z", "Q"},
			miss:   []string{"", "z", "q", "b", "ZQ", "QZ", "Zq", "ab", "ba", "abc"},
		},
		{
			// TODO: convert to runesranges - but may be empty element
			glob: "{Z,Q,}",
			want: &Glob{
				Glob: "{Z,Q,}", Node: "{,Q,Z}", MinLen: 0, MaxLen: 1,
				Items: []items.Item{
					&items.StringList{
						Vals: []string{"Q", "Z"}, MinSize: 0, MaxSize: 1,
						FirstASCII: utils.MakeASCIISetMust("QZ"), ASCIIStarted: true,
					},
				},
			},
			verify: `^(Z|Q)?$`,
			match:  []string{"", "Z", "Q"},
			miss:   []string{"z", "q", "b", "ZQ", "QZ", "Zq", "ab", "ba", "abc"},
		},
		{
			glob: "{a,bc}",
			want: &Glob{
				Glob: "{a,bc}", Node: "{a,bc}", MinLen: 1, MaxLen: 2,
				Items: []items.Item{
					&items.StringList{
						Vals: []string{"a", "bc"}, MinSize: 1, MaxSize: 2,
						FirstASCII: utils.MakeASCIISetMust("ab"), ASCIIStarted: true,
					},
				},
			},
			verify: `^(a|bc)$`,
			match:  []string{"a", "bc"},
			miss:   []string{"", "b", "ab", "ba", "abc", "bca"},
		},
		{
			glob: "*{a,bc}",
			want: &Glob{
				Glob: "*{a,bc}", Node: "*{a,bc}", MinLen: 1, MaxLen: -1,
				Items: []items.Item{
					items.Star(0),
					&items.StringList{
						Vals: []string{"a", "bc"}, MinSize: 1, MaxSize: 2,
						FirstASCII: utils.MakeASCIISetMust("ab"), ASCIIStarted: true,
					},
				},
			},
			verify: `^.*(a|bc)$`,
			match:  []string{"a", "aa", "ba", "bc", "Ba", "abc", "bca", "Bbc", "aBa"},
			miss:   []string{"", "b", "ab"},
		},
		{
			glob: "a{a,bc}{qa,q}c",
			want: &Glob{
				Glob: "a{a,bc}{qa,q}c", Node: "a{a,bc}{q,qa}c",
				MinLen: 4, MaxLen: 6, Prefix: "a", Suffix: "c",
				Items: []items.Item{
					&items.StringList{
						Vals: []string{"a", "bc"}, MinSize: 1, MaxSize: 2,
						FirstASCII: utils.MakeASCIISetMust("ab"), ASCIIStarted: true,
					},
					&items.StringList{
						Vals: []string{"q", "qa"}, MinSize: 1, MaxSize: 2,
						FirstASCII: utils.MakeASCIISetMust("q"), ASCIIStarted: true,
					},
				},
			},
			verify: `^a(a|bc)(q|qa)c$`,
			match:  []string{"aaqac", "abcqac", "aaqc"},
			miss:   []string{"", "b", "ab", "ba", "abc", "aabc", "aaqbc"},
		},
		{
			glob: "a{a,bc}Z{qa,q}c",
			want: &Glob{
				Glob: "a{a,bc}Z{qa,q}c", Node: "a{a,bc}Z{q,qa}c",
				MinLen: 5, MaxLen: 7, Prefix: "a", Suffix: "c",
				Items: []items.Item{
					&items.StringList{
						Vals: []string{"a", "bc"}, MinSize: 1, MaxSize: 2,
						FirstASCII: utils.MakeASCIISetMust("ab"), ASCIIStarted: true,
					},
					items.Byte('Z'),
					&items.StringList{
						Vals: []string{"q", "qa"}, MinSize: 1, MaxSize: 2,
						FirstASCII: utils.MakeASCIISetMust("q"), ASCIIStarted: true,
					},
				},
			},
			verify: `^a(a|bc)Z(q|qa)c$`,
			match:  []string{"aaZqac", "abcZqac", "aaZqc"},
			miss:   []string{"", "b", "ab", "ba", "abc", "aabc", "aaqbc"},
		},
		{
			glob: "a{,a,bc}Z{qa,q}c",
			want: &Glob{
				Glob: "a{,a,bc}Z{qa,q}c", Node: "a{,a,bc}Z{q,qa}c",
				MinLen: 4, MaxLen: 7, Prefix: "a", Suffix: "c",
				Items: []items.Item{
					&items.StringList{
						Vals: []string{"a", "bc"}, MinSize: 0, MaxSize: 2,
						FirstASCII: utils.MakeASCIISetMust("ab"), ASCIIStarted: true,
					},
					items.Byte('Z'),
					&items.StringList{
						Vals: []string{"q", "qa"}, MinSize: 1, MaxSize: 2,
						FirstASCII: utils.MakeASCIISetMust("q"), ASCIIStarted: true,
					},
				},
			},
			verify: `^a(a|bc)?Z(q|qa)c$`,
			match: []string{
				"aaZqac", "abcZqac", "aaZqc",
				// empty first list segment
				"aZqc", "aZqac",
			},
			miss: []string{"", "b", "ab", "ba", "abc", "aabc", "aaqbc"},
		},
		{
			glob: "*.{bc,d}*.e",
			want: &Glob{
				Glob: "*.{bc,d}*.e", Node: "*.{bc,d}*.e", MinLen: 4, MaxLen: -1, Suffix: ".e",
				Items: []items.Item{
					items.Star(0), items.Byte('.'),
					&items.StringList{
						Vals: []string{"bc", "d"}, MinSize: 1, MaxSize: 2,
						FirstASCII: utils.MakeASCIISetMust("bd"), ASCIIStarted: true,
					},
					items.Star(0),
				},
			},
			match: []string{"a.b.bce.e"},
		},
		{
			// TODO: convert to RunesRanges with empty
			glob: "b{a,}",
			want: &Glob{
				Glob: "b{a,}", Node: "b{,a}", MinLen: 1, MaxLen: 2, Prefix: "b",
				Items: []items.Item{
					&items.StringList{
						Vals: []string{"a"}, MinSize: 0, MaxSize: 1,
						FirstASCII: utils.MakeASCIISetMust("a"), ASCIIStarted: true,
					},
				},
			},
			verify: `^ba?$`,
			match:  []string{"b", "ba"},
			miss:   []string{"", "bb", "bd", "ab", "bab", "ba.b"},
		},
		{
			// TODO: convert to RunesRanges with empty
			glob: "*{a,}",
			want: &Glob{
				Glob: "*{a,}", Node: "*{,a}", MinLen: 0, MaxLen: -1,
				Items: []items.Item{
					items.Star(0),
					&items.StringList{
						Vals: []string{"a"}, MinSize: 0, MaxSize: 1,
						FirstASCII: utils.MakeASCIISetMust("a"), ASCIIStarted: true,
					},
				},
			},
			verify: `^.*a?$`,
			match:  []string{"", "a", "b", "ba", "bb", "bd", "ab", "bab", "ba.b"},
		},
		{
			glob: "*{a,}{b,cd}",
			want: &Glob{
				Glob: "*{a,}{b,cd}", Node: "*{,a}{b,cd}", MinLen: 1, MaxLen: -1,
				Items: []items.Item{
					items.Star(0),
					&items.StringList{
						Vals: []string{"a"}, MinSize: 0, MaxSize: 1,
						FirstASCII: utils.MakeASCIISetMust("a"), ASCIIStarted: true,
					},
					&items.StringList{
						Vals: []string{"b", "cd"}, MinSize: 1, MaxSize: 2,
						FirstASCII: utils.MakeASCIISetMust("bc"), ASCIIStarted: true,
					},
				},
			},
			verify: `^.*a?(b|cd)$`,
			match: []string{"abcd",
				"b", "cd", "ab", "acd", "bb", "bcd", "bab", "abcd", "bacd", "bbb", "ba.b", "ba.cd",
			},
			miss: []string{"", "bbd", "abc"},
		},
		{
			// TODO: convert to RunesRanges with empty
			glob: "b*{a,}",
			want: &Glob{
				Glob: "b*{a,}", Node: "b*{,a}", MinLen: 1, MaxLen: -1, Prefix: "b",
				Items: []items.Item{
					items.Star(0),
					&items.StringList{
						Vals: []string{"a"}, MinSize: 0, MaxSize: 1,
						FirstASCII: utils.MakeASCIISetMust("a"), ASCIIStarted: true,
					},
				},
			},
			verify: `^b.*a?$`,
			match:  []string{"b", "ba", "bb", "bba", "bbb", "bbd", "bab", "bbab", "bba.b"},
			miss:   []string{"", "a", "ab", "aba"},
		},
		{
			// TODO: convert to RunesRanges with empty
			glob: "b*{a,}*",
			want: &Glob{
				Glob: "b*{a,}*", Node: "b*{,a}*", Prefix: "b", MinLen: 1, MaxLen: -1,
				Items: []items.Item{
					items.Star(0),
					&items.StringList{
						Vals: []string{"a"}, MinSize: 0, MaxSize: 1,
						FirstASCII: utils.MakeASCIISetMust("a"), ASCIIStarted: true,
					},
					items.Star(0),
				},
			},
			match: []string{"bc..df"},
		},
		// glob star
		{
			glob: "DB.*.Cassandra.{BalanceCluster,BalanceStaging,Billing,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount",
			want: &Glob{
				Glob:   "DB.*.Cassandra.{BalanceCluster,BalanceStaging,Billing,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount",
				Node:   "DB.*.Cassandra.{BalanceCluster,BalanceStaging,Billing,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount",
				MinLen: 40, MaxLen: -1, Prefix: "DB.", Suffix: ".DownEndpointCount",
				Items: []items.Item{
					items.Star(0), items.NewString(".Cassandra."),
					&items.StringList{
						Vals: []string{
							"BalanceCluster", "BalanceStaging", "Billing",
							"UpProduction", "UpTesting", "WebCluster",
						},
						MinSize: 7, MaxSize: 14,
						FirstASCII: utils.MakeASCIISetMust("BUW"), ASCIIStarted: true,
					},
					items.Byte('.'), items.Star(0),
				},
			},
			verify: `^DB\..*\.Cassandra\.(BalanceCluster|BalanceStaging|Billing|UpProduction|UpTesting|WebCluster)\..*\.DownEndpointCount$`,
			match: []string{
				"DB.Sales.Cassandra.UpProduction.NODE1.DownEndpointCount",
				"DB.Sales.Cassandra.UpProduction.node1.NODE2.DownEndpointCount",
				"DB..Cassandra.UpProduction..DownEndpointCount",
				"DB..Cassandra.UpProduction.node1.DownEndpointCount",
			},
			miss: []string{
				"DB.Cassandra..UpProduction.DownEndpointCount",
				"DB.Cassandra.UpProduction.DownEndpointCount",
				"DB.Cassandra.UpProduction.node1.DownEndpointCount",
				"DB.Sales.Cassandra.SalesCluster.node1.DownEndpointCount",
			},
		},
		{
			glob: "DB.*?.Cassandra.{BalanceCluster,BalanceStaging,Billing,UpProduction,UpTesting,WebCluster}.*?.DownEndpointCount",
			want: &Glob{
				Glob:   "DB.*?.Cassandra.{BalanceCluster,BalanceStaging,Billing,UpProduction,UpTesting,WebCluster}.*?.DownEndpointCount",
				Node:   "DB.*?.Cassandra.{BalanceCluster,BalanceStaging,Billing,UpProduction,UpTesting,WebCluster}.*?.DownEndpointCount",
				MinLen: 42, MaxLen: -1, Prefix: "DB.", Suffix: ".DownEndpointCount",
				Items: []items.Item{
					items.Star(1), items.NewString(".Cassandra."),
					&items.StringList{
						Vals: []string{
							"BalanceCluster", "BalanceStaging", "Billing",
							"UpProduction", "UpTesting", "WebCluster",
						},
						MinSize: 7, MaxSize: 14,
						FirstASCII: utils.MakeASCIISetMust("BUW"), ASCIIStarted: true,
					},
					items.Byte('.'), items.Star(1),
				},
			},
			verify: `^DB\..+\.Cassandra\.(BalanceCluster|BalanceStaging|Billing|UpProduction|UpTesting|WebCluster)\..+\.DownEndpointCount$`,
			match: []string{
				"DB.a.Cassandra.UpProduction.node1.DownEndpointCount",
				"DB.a.Cassandra.UpProduction.b.DownEndpointCount",
				"DB.Sales.Cassandra.UpProduction.node1.DownEndpointCount",
				"DB.Sales.Cassandra.UpProduction.node1.NODE2.DownEndpointCount",
			},
			miss: []string{
				"Web.App.UpProduction.DownEndpointCount",
				"Web.a.Cassandra.UpProduction.node1.DownEndpointCount",
				"DB.a.Cassandra.UpProduction.node1.DownEndpointCount2",
				"DB.Cassandra.UpProduction.DownEndpointCount",
				"DB.Cassandra.UpProduction.node1.DownEndpointCount",
				"DB..Cassandra.UpProduction.node1.DownEndpointCount",
				"DB..Cassandra.UpProduction..DownEndpointCount",
				"DB.Sales.Cassandra.SalesCluster.node1.DownEndpointCount",
			},
		},
		{
			glob: "DB.*{BalanceCluster,BalanceStaging,Billing,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount",
			want: &Glob{
				Glob:   "DB.*{BalanceCluster,BalanceStaging,Billing,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount",
				Node:   "DB.*{BalanceCluster,BalanceStaging,Billing,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount",
				MinLen: 29, MaxLen: -1, Prefix: "DB.", Suffix: ".DownEndpointCount",
				Items: []items.Item{
					items.Star(0),
					&items.StringList{
						Vals: []string{
							"BalanceCluster", "BalanceStaging", "Billing",
							"UpProduction", "UpTesting", "WebCluster",
						},
						MinSize: 7, MaxSize: 14,
						FirstASCII: utils.MakeASCIISetMust("BUW"), ASCIIStarted: true,
					},
					items.Byte('.'), items.Star(0),
				},
			},
			verify: `^DB\..+(BalanceCluster|BalanceStaging|Billing|UpProduction|UpTesting|WebCluster)\..*\.DownEndpointCount$`,
			match: []string{
				"DB.Cassandra.UpProduction.node1.DownEndpointCount",
				"DB..Cassandra.UpProduction.node1.DownEndpointCount",
				"DB..Cassandra.UpProduction..DownEndpointCount",
				"DB.a.Cassandra.UpProduction.node1.DownEndpointCount",
				"DB.a.Cassandra.UpProduction.b.DownEndpointCount",
				"DB.Sales.Cassandra.UpProduction.node1.DownEndpointCount",
				"DB.Sales.Cassandra.UpProduction.node1.NODE2.DownEndpointCount",
			},
			miss: []string{
				"Web.App.UpProduction.DownEndpointCount",
				"Web.a.Cassandra.UpProduction.node1.DownEndpointCount",
				"DB.a.Cassandra.UpProduction.node1.DownEndpointCount2",
				"DB.Cassandra.UpProduction.DownEndpointCount",
				"DB.Sales.Cassandra.SalesCluster.node1.DownEndpointCount",
			},
		},
	}
	for n, tt := range tests {
		runTestGlob(t, n, tt)
	}
}

func TestGlob_List_Broken(t *testing.T) {
	tests := []testGlob{
		// broken
		{glob: "{ac", wantErr: true},
		{glob: "a}c", wantErr: true},
		// skip empty
		{
			glob:  "{}a",
			want:  &Glob{Glob: "{}a", Node: "a", MinLen: 1, MaxLen: 1},
			match: []string{"a"},
			miss:  []string{"", "b", "ab"},
		},
	}
	for n, tt := range tests {
		runTestGlob(t, n, tt)
	}
}

// becnmark for list gready skip scan optimization
var (
	globStarListASCII = "DB.*{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount"
	stringListASCII   = strings.Repeat("DB.Sales.Cassandra.", 20) +
		"SalesCluster.node1.DownEndpointCount"
)

func Benchmark_Star_StringList_ASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globStarListASCII)
		if !g.Match(stringListASCII) {
			b.Fatal(stringListASCII)
		}
	}
}

func Benchmark_Star_StringList_ASCII_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globStarListASCII)
		if !w.MatchString(stringListASCII) {
			b.Fatal(stringListASCII)
		}
	}
}

func Benchmark_Star_StringList_ASCII_Precompiled(b *testing.B) {
	g := ParseMust(globStarListASCII)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringListASCII) {
			b.Fatal(stringListASCII)
		}
	}
}

func Benchmark_Star_StringList_Regex_Precompiled(b *testing.B) {
	g := tests.BuildGlobRegexp(globStarListASCII)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.MatchString(stringListASCII) {
			b.Fatal(stringListASCII)
		}
	}
}

// becnmark for list gready skip scan optimization
var (
	globStringList03 = "{BalanceCluster,BillingTesting,StoreCluster}"
	globStringList05 = "{BalanceCluster,BillingTesting,BlobStorageCluster,BusinessStatCluster,StoreCluster}"
	globStringList10 = "{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,StoreCluster}"
	globStringList50 = "{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster,QCluster,QTest,QCluster,ZProduction,ZTesting,ZCluster}"
	stringList       = "StoreCluster"
)

func Benchmark_StringList03_Precompiled(b *testing.B) {
	g := ParseMust(globStringList03)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringList) {
			b.Fatal(stringList)
		}
	}
}

func Benchmark_StringList05_Precompiled(b *testing.B) {
	g := ParseMust(globStringList05)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringList) {
			b.Fatal(stringList)
		}
	}
}

func Benchmark_StringList10_Precompiled(b *testing.B) {
	g := ParseMust(globStringList10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringList) {
			b.Fatal(stringList)
		}
	}
}

func Benchmark_StringList50_Precompiled(b *testing.B) {
	g := ParseMust(globStringList50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !g.Match(stringList) {
			b.Fatal(stringList)
		}
	}
}
