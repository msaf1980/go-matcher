package gtags

import (
	"regexp"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
)

var (
	queriesBatchList = []string{
		`seriesByTag('name=DownEndpointCount', 'project=Sales', 'app=Cassandra', 'cluster={BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}')`,
		`seriesByTag('name=Status', 'project=Sales', 'app=Cassandra', 'cluster={BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}')`,
		`seriesByTag('name=DownEndpointCount', 'project=Sales', 'app=Postgresql', 'cluster={BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}')`,
		`seriesByTag('name=Status', 'project=Sales', 'app=Postgresql', 'cluster={BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}')`,
		`seriesByTag('name=DownEndpointCount', 'project=Sales', 'app=MSSQL', 'cluster={BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}')`,
		`seriesByTag('name=Status', 'project=Sales', 'app=MSSQL', 'cluster={BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}')`,
	}
	termsBatchList = taggedTermListList(queriesBatchList)

	pathsBatchList = []string{
		"DownEndpointCount?app=Cassandra&cluster=SalesCluster&node=node1&project=Sales",
		"DownEndpointCount?app=Cassandra&cluster=SalesCluster&node=node2&project=Sales",
		"DownEndpointCount?app=Cassandra&cluster=SalesCluster&node=node3&project=Sales",
		"DownEndpointCount?app=Cassandra&cluster=SalesCluster&node=node4&project=Sales",
		"DownEndpointCount?app=Cassandra&cluster=SalesCluster&node=node5&project=Sales",
		"DownEndpointCount?app=Cassandra&cluster=SalesCluster&node=node6&project=Sales",
		"DownEndpointCount?app=Cassandra&cluster=SalesCluster&node=node7&project=Sales",
		"DownEndpointCount?app=Cassandra&cluster=SalesCluster&node=node8&project=Sales",
		"DownEndpointCount?app=Cassandra&cluster=SalesCluster&node=node9&project=Sales",
		"DownEndpointCount?app=Cassandra&cluster=SalesCluster&node=node10&project=Sales",
		"Status?app=Cassandra&cluster=SalesCluster&project=Sales",
		"Status?app=Cassandra&cluster=NoSalesCluster&project=Sales",
		"Status?app=Postgresql&cluster=SalesCluster&project=Sales",
		"Status?app=Postgresql&cluster=NoSalesCluster&project=Sales",
		"Status?app=MSSQL&cluster=SalesCluster&project=Sales",
		"Status?app=MSSQL&cluster=NoSalesCluster&project=Sales",
	}
	tagsBatchList = storedTagsList(pathsBatchList)

	reBatchList = []*regexp.Regexp{
		regexp.MustCompile(`DownEndpointCount\?(.*&)?app=Cassandra(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
		regexp.MustCompile(`Status\?(.*&)?app=Cassandra(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
		regexp.MustCompile(`DownEndpointCount\?(.*&)?app=Postgresql(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
		regexp.MustCompile(`Status\?(.*&)?app=Postgresql(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
		regexp.MustCompile(`DownEndpointCount\?(.*&)?app=MSSQL(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
		regexp.MustCompile(`Status\?(.*&)?app=MSSQL(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
	}
)

func BenchmarkBatch_Terms_Precompiled(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			for _, terms := range termsBatchList {
				tags, err := PathTags(path)
				if err != nil {
					b.Fatal(err)
				}
				terms.MatchByTags(tags)
			}
		}
	}
}

func BenchmarkBatch_Terms_Prealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tags := range tagsBatchList {
			for _, terms := range termsBatchList {
				terms.MatchByTags(tags)
			}
		}
	}
}

func BenchmarkBatch_Tree_ByTags_Precompiled(b *testing.B) {
	w := NewTree()
	for i, q := range queriesBatchList {
		_, _, err := w.Add(q, i)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			tags, err := PathTags(path)
			if err != nil {
				b.Fatal(err)
			}
			var store items.AllStore
			store.Init()
			store.Grow(1)
			_ = w.MatchByTags(tags, &store)
		}
	}
}

func BenchmarkBatch_Tree_ByTags_Prealloc(b *testing.B) {
	w := NewTree()
	for i, q := range queriesBatchList {
		_, _, err := w.Add(q, i)
		if err != nil {
			b.Fatal(err)
		}
	}
	var store items.AllStore
	store.Init()
	store.Grow(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tags := range tagsBatchList {
			store.Init()
			_ = w.MatchByTags(tags, &store)
		}
	}
}

func _BenchmarkBatch_Regex_Precompiled(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			if !reBatchList[0].MatchString(path) {
				b.Fatalf("%s\n%s", path, reBatchList[0].String())
			}
		}
	}
}
