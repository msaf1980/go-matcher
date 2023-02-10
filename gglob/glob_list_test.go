package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

func TestGlobMatcher_List(t *testing.T) {
	tests := []testGlobMatcher{
		{
			name: `{"{a,bc}"}`, globs: []string{"{a,bc}"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "{a,bc}", Terminated: "{a,bc}", MinSize: 1, MaxSize: 2,
								Inners: []items.InnerItem{
									&items.ItemList{Vals: []string{"a", "bc"}, ValsMin: 1, ValsMax: 2},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"{a,bc}": true},
			},
			matchPaths: map[string][]string{"a": {"{a,bc}"}, "bc": {"{a,bc}"}},
			missPaths:  []string{"", "b", "ab", "ba", "abc"},
		},
		{
			name: `{"a{a,bc}{qa,q}c"}`, globs: []string{"a{a,bc}{qa,q}c"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a{a,bc}{qa,q}c", Terminated: "a{a,bc}{qa,q}c", MinSize: 4, MaxSize: 6,
								P: "a", Suffix: "c",
								Inners: []items.InnerItem{
									&items.ItemList{Vals: []string{"a", "bc"}, ValsMin: 1, ValsMax: 2},
									&items.ItemList{Vals: []string{"q", "qa"}, ValsMin: 1, ValsMax: 2},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"a{a,bc}{qa,q}c": true},
			},
			matchPaths: map[string][]string{"aaqac": {"a{a,bc}{qa,q}c"}, "abcqac": {"a{a,bc}{qa,q}c"}, "aaqc": {"a{a,bc}{qa,q}c"}},
			missPaths:  []string{"", "b", "ab", "ba", "abc", "aabc", "aaqbc"},
		},
		{
			name: `{"a{a,bc}Z{qa,q}c"}`, globs: []string{"a{a,bc}Z{qa,q}c"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{
								Node: "a{a,bc}Z{qa,q}c", Terminated: "a{a,bc}Z{qa,q}c", MinSize: 5, MaxSize: 7,
								P: "a", Suffix: "c",
								Inners: []items.InnerItem{
									&items.ItemList{Vals: []string{"a", "bc"}, ValsMin: 1, ValsMax: 2},
									items.ItemString("Z"),
									&items.ItemList{Vals: []string{"q", "qa"}, ValsMin: 1, ValsMax: 2},
								},
							},
						},
					},
				},
				Globs: map[string]bool{"a{a,bc}Z{qa,q}c": true},
			},
			matchPaths: map[string][]string{"aaZqac": {"a{a,bc}Z{qa,q}c"}, "abcZqac": {"a{a,bc}Z{qa,q}c"}, "aaZqc": {"a{a,bc}Z{qa,q}c"}},
			missPaths:  []string{"", "b", "ab", "ba", "abc", "aabc", "aaqbc"},
		},
		// one item optimization
		{
			name: `{"{a}"}`, globs: []string{"{a}"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{Node: "{a}", Terminated: "{a}", MinSize: 1, MaxSize: 1, P: "a"},
						},
					},
				},
				Globs: map[string]bool{"{a}": true},
			},
			matchPaths: map[string][]string{"a": {"{a}"}},
			missPaths:  []string{"", "b", "d", "ab", "a.b"},
		},
		{
			name: `{"{a,}"}`, globs: []string{"{a,}"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{Node: "{a,}", Terminated: "{a,}", MinSize: 1, MaxSize: 1, P: "a"},
						},
					},
				},
				Globs: map[string]bool{"{a,}": true},
			},
			matchPaths: map[string][]string{"a": {"{a,}"}},
			missPaths:  []string{"", "b", "d", "ab", "a.b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

func TestGlobMatcher_List_Broken(t *testing.T) {
	tests := []testGlobMatcher{
		// broken
		{name: `{"z{ac"}`, globs: []string{"{ac"}, wantErr: true},
		{name: `{"a}c"}`, globs: []string{"a}c"}, wantErr: true},
		// skip empty
		{
			name: `{"{}a"}`, globs: []string{"{}a"},
			wantW: &GlobMatcher{
				Root: map[int]*items.NodeItem{
					1: {
						Childs: []*items.NodeItem{
							{Node: "{}a", Terminated: "{}a", MinSize: 1, MaxSize: 1, P: "a"},
						},
					},
				},
				Globs: map[string]bool{"{}a": true},
			},
			matchPaths: map[string][]string{"a": {"{}a"}},
			missPaths:  []string{"", "b", "ab"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestGlobMatcher(t, tt)
		})
	}
}

var (
	targetsBatchList = []string{
		"DB.*.Cassandra.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount.DownEndpointCount",
		"DB.*.Cassandra.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.Status",
		"DB.*.Postgresql.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.*.queries_time.p99",
		"DB.*.Postgresql.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.Status",
		"DB.*.MSSQL.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount.DownEndpointCount",
		"DB.*.MSSQL.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.Status",
	}
	pathsBatchList = []string{
		"DB.Sales.Cassandra.SalesCluster.node1.DownEndpointCount.DownEndpointCount",
		"DB.Sales.Cassandra.NoSalesCluster.node2.DownEndpointCount.DownEndpointCount",
		"DB.Sales.Cassandra.SalesCluster.node3.DownEndpointCount.DownEndpointCount",
		"DB.Sales.Cassandra.SalesCluster.node4.DownEndpointCount.DownEndpointCount",
		"DB.Sales.Cassandra.SalesCluster.node5.DownEndpointCount.DownEndpointCount",
		"DB.Sales.Cassandra.SalesCluster.node6.DownEndpointCount.DownEndpointCount",
		"DB.Sales.Cassandra.SalesCluster.node7.DownEndpointCount.DownEndpointCount",
		"DB.Sales.Cassandra.SalesCluster.node8.DownEndpointCount.DownEndpointCount",
		"DB.Sales.Cassandra.SalesCluster.node9.DownEndpointCount.DownEndpointCount",
		"DB.Sales.Cassandra.SalesCluster.node10.DownEndpointCount.DownEndpointCount",
		"DB.Sales.Cassandra.SalesCluster.Status",
		"DB.Sales.Cassandra.NoSalesCluster.Status",
		"DB.Sales.Postgresql.SalesCluster.Status",
		"DB.Sales.Postgresql.NoSalesCluster.Status",
		"DB.MSSQL.Postgresql.SalesCluster.Status",
		"DB.MSSQL.Postgresql.NoSalesCluster.Status",
	}
)

// becnmark for suffix optimization
func BenchmarkList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		err := w.Add(targetsBatchList[0])
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathsBatchList[0])
		if len(globs) != 1 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkList_ByParts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		err := w.Add(targetsBatchList[0])
		if err != nil {
			b.Fatal(err)
		}
		parts := items.PathSplit(pathsBatchList[0])
		globs := w.MatchByParts(parts)
		if len(globs) != 1 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkList_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetsBatchList[0])
		if !w.MatchString(pathsBatchList[0]) {
			b.Fatal(pathsBatchList[0])
		}
	}
}

func BenchmarkList_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetsBatchList[0])
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathsBatchList[0])
		if len(globs) != 1 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkList_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetsBatchList[0])
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathsBatchList[0], &globs)
		if len(globs) != 1 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkList_Prealloc_ByParts(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetsBatchList[0])
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parts := items.PathSplit(pathsBatchList[0])
		globs = globs[:0]
		w.MatchByPartsB(parts, &globs)
		if len(globs) != 1 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkList_Prealloc_ByParts2(b *testing.B) {
	parts := items.PathSplit(pathsBatchList[0])
	w := NewGlobMatcher()
	err := w.Add(targetsBatchList[0])
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchByPartsB(parts, &globs)
		if len(globs) != 1 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkList_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetsBatchList[0])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !w.MatchString(pathsBatchList[0]) {
			b.Fatal(pathsBatchList[0])
		}
	}
}

func BenchmarkList_Batch(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Adds(targetsBatchList)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			_ = w.Match(path)
		}
	}
}

func BenchmarkList_Batch_ByParts(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Adds(targetsBatchList)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			parts := items.PathSplit(path)
			_ = w.MatchByParts(parts)
		}
	}
}

func BenchmarkList_Batch_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetsBatchList[0])
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			globs = globs[:0]
			w.MatchB(path, &globs)
		}
	}
}

func BenchmarkList_Batch_Prealloc_ByParts(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetsBatchList[0])
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			parts := items.PathSplit(path)
			globs = globs[:0]
			w.MatchByPartsB(parts, &globs)
		}
	}
}

func BenchmarkList_Batch_Prealloc_ByParts2(b *testing.B) {
	w := NewGlobMatcher()
	err := w.Add(targetsBatchList[0])
	if err != nil {
		b.Fatal(err)
	}
	partsBatchList := make([][]string, len(pathsBatchList))
	for i, path := range pathsBatchList {
		partsBatchList[i] = items.PathSplit(path)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, parts := range partsBatchList {
			globs = globs[:0]
			w.MatchByPartsB(parts, &globs)
		}
	}
}

func BenchmarkList_Batch_Regex(b *testing.B) {
	w := buildGlobRegexp(targetsBatchList[0])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			w.MatchString(path)
		}
	}
}
