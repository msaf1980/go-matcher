package gglob

import (
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

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

func BenchmarkBatch_List_Tree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		for j := 0; j < len(targetsBatchList); j++ {
			_, _, err := w.Add(targetsBatchList[j], j)
			if err != nil {
				b.Fatal(err)
			}
		}
		var globs []string
		first := items.MinStore{-1}
		for j := 0; j < len(pathsBatchList); j++ {
			_ = w.Match(pathsBatchList[j], &globs, nil, &first)
		}
	}
}

func BenchmarkBatch_List_GGlob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := parseGGlobs(targetsBatchList)
		for j := 0; j < len(targetsBatchList); j++ {
			for k := 0; k < len(g); k++ {
				_ = g[k].Match(targetsBatchList[j])
			}
		}
	}
}

func BenchmarkBatch_List_Tree_Precompiled(b *testing.B) {
	g := parseGGlobs(targetsBatchList)

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
		for j := 0; j < len(pathsBatchList); j++ {
			_ = w.Match(pathsBatchList[j], &globs, nil, &first)
		}
	}
}

func BenchmarkBatch_List_Tree_Precompiled2(b *testing.B) {
	g := parseGGlobs(targetsBatchList)
	w := NewTree()
	for j := 0; j < len(g); j++ {
		_, _, err := w.AddGlob(g[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var globs []string
		first := items.MinStore{-1}
		for j := 0; j < len(pathsBatchList); j++ {
			_ = w.Match(pathsBatchList[j], &globs, nil, &first)
		}
	}
}

func BenchmarkBatch_List_GGlob_Precompiled(b *testing.B) {
	g := parseGGlobs(targetsBatchList)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(targetsBatchList); j++ {
			for k := 0; k < len(g); k++ {
				_ = g[k].Match(targetsBatchList[j])
			}
		}
	}
}

func BenchmarkBatch_List_Tree_Prealloc(b *testing.B) {
	g := parseGGlobs(targetsBatchList)
	w := NewTree()
	for j := 0; j < len(g); j++ {
		_, _, err := w.AddGlob(g[j], j)
		if err != nil {
			b.Fatal(err)
		}
	}
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(pathsBatchList); j++ {
			globs = globs[:0]
			first.Init()
			_ = w.Match(pathsBatchList[j], &globs, nil, &first)
		}
	}
}
