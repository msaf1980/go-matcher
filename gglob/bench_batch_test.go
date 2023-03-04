package gglob

import (
	"testing"
	"time"

	"github.com/msaf1980/go-matcher/pkg/items"
)

var (
	globsBatchList = []string{
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
		"DB.Store.Postgresql.SalesCluster.Status",
		"DB.Store.Postgresql.NoSalesCluster.Status",
	}
)

func BenchmarkBatch_List_Tree(b *testing.B) {
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewTree()
		for j := 0; j < len(globsBatchList); j++ {
			_, _, err := w.Add(globsBatchList[j], j)
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
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchList))/d.Seconds(), "match/s")
}

func BenchmarkBatch_List_Tree_ByParts(b *testing.B) {
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewTree()
		for j := 0; j < len(globsBatchList); j++ {
			_, _, err := w.Add(globsBatchList[j], j)
			if err != nil {
				b.Fatal(err)
			}
		}
		var globs []string
		first := items.MinStore{-1}
		for j := 0; j < len(pathsBatchList); j++ {
			parts := PathSplit(pathsBatchList[j])
			_ = w.MatchByParts(parts, &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchList))/d.Seconds(), "match/s")
}

func BenchmarkBatch_List_GGlob(b *testing.B) {
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g := parseGGlobs(globsBatchList)
		for j := 0; j < len(globsBatchList); j++ {
			for k := 0; k < len(g); k++ {
				_ = g[k].Match(globsBatchList[j])
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchList))/d.Seconds(), "match/s")
}

func BenchmarkBatch_List_Tree_Precompiled(b *testing.B) {
	g := parseGGlobs(globsBatchList)

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
		for j := 0; j < len(pathsBatchList); j++ {
			_ = w.Match(pathsBatchList[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchList))/d.Seconds(), "match/s")
}

func BenchmarkBatch_List_Tree_Precompiled2(b *testing.B) {
	g := parseGGlobs(globsBatchList)
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
		for j := 0; j < len(pathsBatchList); j++ {
			_ = w.Match(pathsBatchList[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchList))/d.Seconds(), "match/s")
}

func BenchmarkBatch_List_GGlob_Precompiled(b *testing.B) {
	g := parseGGlobs(globsBatchList)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(globsBatchList); j++ {
			for k := 0; k < len(g); k++ {
				_ = g[k].Match(globsBatchList[j])
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchList))/d.Seconds(), "match/s")
}

func BenchmarkBatch_List_GGlob_Prealloc_ByParts(b *testing.B) {
	g := parseGGlobs(globsBatchList)
	parts := make([]string, 8)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(globsBatchList); j++ {
			PathSplitB(globsBatchList[j], &parts)
			length := len(globsBatchList[j])
			for k := 0; k < len(g); k++ {
				_ = g[k].MatchByParts(parts, length)
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchList))/d.Seconds(), "match/s")
}

func BenchmarkBatch_List_Tree_Prealloc(b *testing.B) {
	g := parseGGlobs(globsBatchList)
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
		for j := 0; j < len(pathsBatchList); j++ {
			globs = globs[:0]
			first.Init()
			_ = w.Match(pathsBatchList[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchList))/d.Seconds(), "match/s")
}

func BenchmarkBatch_List_Tree_Prealloc_ByParts(b *testing.B) {
	g := parseGGlobs(globsBatchList)
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
		for j := 0; j < len(pathsBatchList); j++ {
			globs = globs[:0]
			first.Init()
			_ = PathSplitB(pathsBatchList[j], &parts)
			_ = w.MatchByParts(parts, &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchList))/d.Seconds(), "match/s")
}

var (
	globBatchMoira = []string{
		"Simple.matching.pattern",
		"Simple.matching.pattern.*",
		"Star.single.*",
		"Star.*.double.any*",
		"Bracket.{one,two,three}.pattern",
		"Bracket.pr{one,two,three}suf",
		"Complex.matching.pattern",
		"Complex.*.*",
		"Complex.*.",
		"Complex.*{one,two,three}suf*.pattern",
		"Question.?at_begin",
		"Question.at_the_end?",
	}

	pathsBatchMoira = []string{"Simple.matching.pattern",
		"Star.single.anything",
		"Star.anything.double.anything",
		"Bracket.one.pattern",
		"Bracket.two.pattern",
		"Bracket.three.pattern",
		"Bracket.pronesuf",
		"Bracket.prtwosuf",
		"Bracket.prthreesuf",
		"Complex.matching.pattern",
		"Complex.anything.pattern",
		"Complex.prefixonesuffix.pattern",
		"Complex.prefixtwofix.pattern",
		"Complex.anything.pattern",
		"Question.1at_begin",
		"Question.at_the_end2",
		"Two.dots..together",
		"Simple.notmatching.pattern",
		"Star.nothing",
		"Bracket.one.nothing",
		"Bracket.nothing.pattern",
		"Complex.prefixonesuffix",
	}
)

func BenchmarkBatch_Moira_Tree(b *testing.B) {
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewTree()
		for j := 0; j < len(globBatchMoira); j++ {
			_, _, err := w.Add(globBatchMoira[j], j)
			if err != nil {
				b.Fatal(err)
			}
		}
		var globs []string
		first := items.MinStore{-1}
		for j := 0; j < len(pathsBatchMoira); j++ {
			_ = w.Match(pathsBatchMoira[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatch_Moira_Tree_ByParts(b *testing.B) {
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := NewTree()
		for j := 0; j < len(globBatchMoira); j++ {
			_, _, err := w.Add(globBatchMoira[j], j)
			if err != nil {
				b.Fatal(err)
			}
		}
		var globs []string
		first := items.MinStore{-1}
		for j := 0; j < len(pathsBatchMoira); j++ {
			parts := PathSplit(pathsBatchMoira[j])
			_ = w.MatchByParts(parts, &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatch_Moira_GGlob(b *testing.B) {
	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g := parseGGlobs(globBatchMoira)
		for j := 0; j < len(globBatchMoira); j++ {
			for k := 0; k < len(g); k++ {
				_ = g[k].Match(globBatchMoira[j])
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatch_Moira_Tree_Precompiled(b *testing.B) {
	g := parseGGlobs(globBatchMoira)

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
		for j := 0; j < len(pathsBatchMoira); j++ {
			_ = w.Match(pathsBatchMoira[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatch_Moira_Tree_Precompiled2(b *testing.B) {
	g := parseGGlobs(globBatchMoira)
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
		for j := 0; j < len(pathsBatchMoira); j++ {
			_ = w.Match(pathsBatchMoira[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatch_Moira_GGlob_Precompiled(b *testing.B) {
	g := parseGGlobs(globBatchMoira)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(globBatchMoira); j++ {
			for k := 0; k < len(g); k++ {
				_ = g[k].Match(globBatchMoira[j])
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatch_Moira_GGlob_Prealloc_ByParts(b *testing.B) {
	g := parseGGlobs(globBatchMoira)
	parts := make([]string, 8)

	start := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(globBatchMoira); j++ {
			PathSplitB(globBatchMoira[j], &parts)
			length := len(globBatchMoira[j])
			for k := 0; k < len(g); k++ {
				_ = g[k].MatchByParts(parts, length)
			}
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatch_Moira_Tree_Prealloc(b *testing.B) {
	g := parseGGlobs(globBatchMoira)
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
		for j := 0; j < len(pathsBatchMoira); j++ {
			globs = globs[:0]
			first.Init()
			_ = w.Match(pathsBatchMoira[j], &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchMoira))/d.Seconds(), "match/s")
}

func BenchmarkBatch_Moira_Tree_Prealloc_ByParts(b *testing.B) {
	g := parseGGlobs(globBatchMoira)
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
		for j := 0; j < len(pathsBatchMoira); j++ {
			globs = globs[:0]
			first.Init()
			_ = PathSplitB(pathsBatchMoira[j], &parts)
			_ = w.MatchByParts(parts, &globs, nil, &first)
		}
	}
	b.StopTimer()
	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
	b.ReportMetric(float64(b.N*len(pathsBatchMoira))/d.Seconds(), "match/s")
}
