package gglob

import (
	"strings"
	"testing"
)

var (
	targetGready_StringMiss = "sys*tgicabcdERt*ltem"
	pathGready_StringMiss   = "sysSKIPSKIPSKIPSKIP_tgicabcdert_SKIPSKIPSKIPSKIPSKIP_tgicabcdeRt_gltem"
)

// becnmark for suffix optimization
func BenchmarkGready_StringMiss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(len(targetGready_StringMiss))
		_, err := w.Add(targetGready_StringMiss, &buf)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathGready_StringMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_StringMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetGready_StringMiss)
		if w.MatchString(pathGready_StringMiss) {
			b.Fatal(pathGready_StringMiss)
		}
	}
}

func BenchmarkGready_StringMiss_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_StringMiss))
	_, err := w.Add(targetGready_StringMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathGready_StringMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_StringMiss_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_StringMiss))
	_, err := w.Add(targetGready_StringMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathGready_StringMiss, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_StringMiss_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetGready_StringMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_StringMiss) {
			b.Fatal(pathGready_StringMiss)
		}
	}
}

var (
	targetGready_RuneMiss = "sys*A*ltem"
	pathGready_RuneMiss   = "sysSKIPSKIPSKIPSKIP_tgicabcdert_SKIPSKIPSKIPSKIPSKIP_tgicabcdeRt_gltem"
)

// becnmark for suffix optimization
func BenchmarkGready_RuneMiss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(len(targetGready_RuneMiss))
		_, err := w.Add(targetGready_RuneMiss, &buf)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathGready_RuneMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_RuneMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetGready_RuneMiss)
		if w.MatchString(pathGready_RuneMiss) {
			b.Fatal(pathGready_RuneMiss)
		}
	}
}

func BenchmarkGready_RuneMiss_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_RuneMiss))
	_, err := w.Add(targetGready_RuneMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathGready_RuneMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_RuneMiss_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_RuneMiss))
	_, err := w.Add(targetGready_RuneMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathGready_RuneMiss, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_RuneMiss_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetGready_RuneMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_RuneMiss) {
			b.Fatal(pathGready_RuneMiss)
		}
	}
}

var (
	targetGready_RuneRangesMiss = "sys*{A-E}*ltem"
	pathGready_RuneRangesMiss   = "sysSKIPSKIPSKIPSKIP_tgicabcdert_SKIPSKIPSKIPSKIPSKIP_tgicabcdeRt_gltem"
)

// becnmark for suffix optimization
func BenchmarkGready_RuneRangesMiss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(len(targetGready_RuneRangesMiss))
		_, err := w.Add(targetGready_RuneRangesMiss, &buf)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathGready_RuneRangesMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_RuneRangesMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetGready_RuneRangesMiss)
		if w.MatchString(pathGready_RuneRangesMiss) {
			b.Fatal(pathGready_RuneRangesMiss)
		}
	}
}

func BenchmarkGready_RuneRangesMiss_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_RuneRangesMiss))
	_, err := w.Add(targetGready_RuneRangesMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathGready_RuneRangesMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_RuneRangesMiss_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_RuneRangesMiss))
	_, err := w.Add(targetGready_RuneRangesMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathGready_RuneRangesMiss, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_RuneRangesMiss_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetGready_RuneRangesMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_RuneRangesMiss) {
			b.Fatal(pathGready_RuneRangesMiss)
		}
	}
}

var (
	targetGready_ListMiss = "DB*{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}_Status"
	pathGready_ListMiss   = "DBCassandraSalesSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIP_NoCluster_Status"
)

// becnmark for suffix optimization
func BenchmarkGready_ListMiss(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(len(targetGready_ListMiss))
		_, err := w.Add(targetGready_ListMiss, &buf)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathGready_ListMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ListMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetGready_ListMiss)
		if w.MatchString(pathGready_ListMiss) {
			b.Fatal(pathGready_ListMiss)
		}
	}
}

func BenchmarkGready_ListMiss_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_ListMiss))
	_, err := w.Add(targetGready_ListMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathGready_ListMiss)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ListMiss_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_ListMiss))
	_, err := w.Add(targetGready_ListMiss, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathGready_ListMiss, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ListMiss_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetGready_ListMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_ListMiss) {
			b.Fatal(pathGready_ListMiss)
		}
	}
}

var (
	targetGready_ListSkip = "DB*{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}_Status"
	pathGready_ListSkip   = "DBCassandraSalesSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIP_WebClusteAr_Status"
)

// becnmark for suffix optimization
func BenchmarkGready_ListSkip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(len(targetGready_ListSkip))
		_, err := w.Add(targetGready_ListSkip, &buf)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathGready_ListSkip)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ListSkip_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetGready_ListSkip)
		if w.MatchString(pathGready_ListSkip) {
			b.Fatal(pathGready_ListSkip)
		}
	}
}

func BenchmarkGready_ListSkip_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_ListSkip))
	_, err := w.Add(targetGready_ListSkip, &buf)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathGready_ListSkip)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ListSkip_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_ListSkip))
	_, err := w.Add(targetGready_ListSkip, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathGready_ListSkip, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ListSkip_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetGready_ListSkip)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_ListSkip) {
			b.Fatal(pathGready_ListSkip)
		}
	}
}

var (
	targetGready_OneSkip = "DB*?web*_Status"
	pathGready_OneSkip   = "DBCassandraSalesSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIP_we_Status"
)

// becnmark for suffix optimization
func BenchmarkGready_OneSkip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewGlobMatcher()
		var buf strings.Builder
		buf.Grow(len(targetGready_OneSkip))
		_, err := w.Add(targetGready_OneSkip, &buf)
		if err != nil {
			b.Fatal(err)
		}
		globs := w.Match(pathGready_OneSkip)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_OneSkip_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := buildGlobRegexp(targetGready_OneSkip)
		if w.MatchString(pathGready_OneSkip) {
			b.Fatal(pathGready_OneSkip)
		}
	}
}

func BenchmarkGready_OneSkip_Precompiled(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_OneSkip))
	_, err := w.Add(targetGready_OneSkip, &buf)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs := w.Match(pathGready_OneSkip)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_OneSkip_Prealloc(b *testing.B) {
	w := NewGlobMatcher()
	var buf strings.Builder
	buf.Grow(len(targetGready_OneSkip))
	_, err := w.Add(targetGready_OneSkip, &buf)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		w.MatchB(pathGready_OneSkip, &globs)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_OneSkip_Precompiled_Regex(b *testing.B) {
	w := buildGlobRegexp(targetGready_OneSkip)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_OneSkip) {
			b.Fatal(pathGready_OneSkip)
		}
	}
}
