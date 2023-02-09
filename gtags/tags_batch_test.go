package gtags

import (
	"regexp"
	"testing"
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

	reBatchList = []*regexp.Regexp{
		regexp.MustCompile(`DownEndpointCount\?(.*&)?app=Cassandra(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
		regexp.MustCompile(`Status\?(.*&)?app=Cassandra(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
		regexp.MustCompile(`DownEndpointCount\?(.*&)?app=Postgresql(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
		regexp.MustCompile(`Status\?(.*&)?app=Postgresql(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
		regexp.MustCompile(`DownEndpointCount\?(.*&)?app=MSSQL(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
		regexp.MustCompile(`Status\?(.*&)?app=MSSQL(.*&)?cluster={BalanceCluster|BalanceStaging|Billing|BillingAutoTesting|BillingDocuments|BillingLoadTesting|BillingTesting|BlobStorageCluster|BusinessStatCluster|CashBoxCluster|CashLogCluster|CassandraClaims|CassandraClientTest|CassandraConnector|CassandraCore|CassandraDev|CassandraReliable|CassandraSentry|CassandraStats|CassandraTasks|CassandraTest|CassandraUsers|CassandraWeb|CoreCluster|CqlCoreCluster|CustomersCluster|LsaMetaindex|EventsCluster|QueueCluster|QueueTesting|LegacyCluster|ProductsCluster|ProductsTestingCluster|RemoteLockCluster|ReportCluster|ReviseCluster|ReviseTestingCluster|SalesCluster|SecondCluster|SecondTest|StoreCluster|UpProduction|UpTesting|WebCluster}(.*&)?project=Sales(&|$)`),
	}
)

func BenchmarkBatch_Precompiled_ByTags(b *testing.B) {
	w := NewTagsMatcher()
	err := w.Adds(queriesBatchList)
	if err != nil {
		b.Fatal(err)
	}
	queries := make([]string, 0, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			tags, err := PathTagsMap(path)
			if err != nil {
				b.Fatal(err)
			}
			w.MatchByTagsB(tags, &queries)
		}
	}
}

func BenchmarkBatch_Precompiled_ByPath(b *testing.B) {
	w := NewTagsMatcher()
	err := w.Adds(queriesBatchList)
	if err != nil {
		b.Fatal(err)
	}
	queries := make([]string, 0, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			w.MatchByPathB(path, &queries)
		}
	}
}

func BenchmarkBatch_Regex_Precompiled_ByTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, path := range pathsBatchList {
			if !reBatchList[0].MatchString(path) {
				b.Fatalf("%s\n%s", path, reBatchList[0].String())
			}
		}
	}
}
