package gglob

import (
	"testing"
	"time"

	"github.com/msaf1980/go-matcher/pkg/items"
)

var (
	globsBatchLargeList = []string{
		"DB.*.Cassandra.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount.DownEndpointCount",
		"DB.*.Cassandra.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.Status",
		"DB.*.Postgresql.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.*.queries_time.p99",
		"DB.*.Postgresql.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.Status",
		"DB.*.MSSQL.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.*.DownEndpointCount.DownEndpointCount",
		"DB.*.MSSQL.{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}.Status",
		"Sales.Backend.{Balance,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster}.*.cpu.load_avg",
		"Sales.Backend.{Balance,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster}.*.memory.free",
		"Sales.Backend.{Balance,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster}.*.memory.usage",
		"Sales.Backend.{Balance,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster}.*.memory.cache",
		"Sales.Front.{Balance,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster}.*.cpu.load_avg",
		"Sales.Front.{Balance,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster}.*.memory.free",
		"Sales.Front.{Balance,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster}.*.memory.usage",
		"Sales.Front.{Balance,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster}.*.memory.cache",
	}
	pathsBatchLargeList = []string{
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
		"DB.Store.Postgresql.SalesCluster.Status",
		"DB.Store.Postgresql.NoSalesCluster.Status",
		// node1
		"Sales.Backend.SecondTest.node1.cpu.load_avg",
		"Sales.Backend.SecondTest.node1.memory.free",
		"Sales.Backend.SecondTest.node1.memory.usage",
		"Sales.Backend.SecondTest.node1.memory.cache",
		"Sales.Front.SecondTest.node1.cpu.load_avg",
		"Sales.Front.SecondTest.node1.memory.free",
		"Sales.Front.SecondTest.node1.memory.usage",
		"Sales.Front.SecondTest.node1.memory.cache",
		// node25
		"Sales.Backend.SecondTest.node25.cpu.load_avg",
		"Sales.Backend.SecondTest.node25.memory.free",
		"Sales.Backend.SecondTest.node25.memory.usage",
		"Sales.Backend.SecondTest.node25.memory.cache",
		"Sales.Front.SecondTest.node25.cpu.load_avg",
		"Sales.Front.SecondTest.node25.memory.free",
		"Sales.Front.SecondTest.node25.memory.usage",
		"Sales.Front.SecondTest.node25.memory.cache",
	}
)

func BenchmarkBatchLarge_List_Tree(b *testing.B) {
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewTree()
		for j := 0; j < len(globsBatchLargeList); j++ {
			_, _, err := w.Add(globsBatchLargeList[j], j)
			if err != nil {
				b.Fatal(err)
			}
		}
		var globs []string
		first := items.MinStore{-1}
		for j := 0; j < len(pathsBatchLargeList); j++ {
			_ = w.Match(pathsBatchLargeList[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchLargeList))/d.Seconds(), "match/s")
}

func BenchmarkBatchLarge_List_Tree_ByParts(b *testing.B) {
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewTree()
		for j := 0; j < len(globsBatchLargeList); j++ {
			_, _, err := w.Add(globsBatchLargeList[j], j)
			if err != nil {
				b.Fatal(err)
			}
		}
		var globs []string
		first := items.MinStore{-1}
		for j := 0; j < len(pathsBatchLargeList); j++ {
			parts := PathSplit(pathsBatchLargeList[j])
			_ = w.MatchByParts(parts, &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchLargeList))/d.Seconds(), "match/s")
}

func BenchmarkBatchLarge_List_GGlob(b *testing.B) {
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g := parseGGlobs(globsBatchLargeList)
		for j := 0; j < len(globsBatchLargeList); j++ {
			for k := 0; k < len(g); k++ {
				_ = g[k].Match(globsBatchLargeList[j])
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchLargeList))/d.Seconds(), "match/s")
}

func BenchmarkBatchLarge_List_Tree_Precompiled(b *testing.B) {
	g := parseGGlobs(globsBatchLargeList)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewTree()
		for j := 0; j < len(g); j++ {
			_, _, err := w.AddGlob(g[j], j)
			if err != nil {
				b.Fatal(err)
			}
		}
		var globs []string
		first := items.MinStore{-1}
		for j := 0; j < len(pathsBatchLargeList); j++ {
			_ = w.Match(pathsBatchLargeList[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchLargeList))/d.Seconds(), "match/s")
}

func BenchmarkBatchLarge_List_Tree_Precompiled2(b *testing.B) {
	g := parseGGlobs(globsBatchLargeList)
	w := NewTree()
	for j := 0; j < len(g); j++ {
		_, _, err := w.AddGlob(g[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var globs []string
		first := items.MinStore{-1}
		for j := 0; j < len(globsBatchLargeList); j++ {
			_ = w.Match(pathsBatchLargeList[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchLargeList))/d.Seconds(), "match/s")
}

func BenchmarkBatchLarge_List_GGlob_Precompiled(b *testing.B) {
	g := parseGGlobs(globsBatchLargeList)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(globsBatchLargeList); j++ {
			for k := 0; k < len(g); k++ {
				_ = g[k].Match(globsBatchLargeList[j])
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchLargeList))/d.Seconds(), "match/s")
}

func BenchmarkBatchLarge_List_GGlob_Prealloc_ByParts(b *testing.B) {
	g := parseGGlobs(globsBatchLargeList)
	parts := make([]string, 8)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(globsBatchLargeList); j++ {
			PathSplitB(globsBatchLargeList[j], &parts)
			length := len(globsBatchLargeList[j])
			for k := 0; k < len(g); k++ {
				_ = g[k].MatchByParts(parts, length)
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchLargeList))/d.Seconds(), "match/s")
}

func BenchmarkBatchLarge_List_Tree_Prealloc(b *testing.B) {
	g := parseGGlobs(globsBatchLargeList)
	w := NewTree()
	for j := 0; j < len(g); j++ {
		_, _, err := w.AddGlob(g[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(pathsBatchLargeList); j++ {
			globs = globs[:0]
			first.Init()
			_ = w.Match(pathsBatchLargeList[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchLargeList))/d.Seconds(), "match/s")
}

func BenchmarkBatchLarge_List_Tree_Prealloc_ByParts(b *testing.B) {
	g := parseGGlobs(globsBatchLargeList)
	w := NewTree()
	for j := 0; j < len(g); j++ {
		_, _, err := w.AddGlob(g[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	parts := make([]string, 10)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(pathsBatchLargeList); j++ {
			globs = globs[:0]
			first.Init()
			_ = PathSplitB(pathsBatchLargeList[j], &parts)
			_ = w.MatchByParts(parts, &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchLargeList))/d.Seconds(), "match/s")
}
