package gglob

import (
	"strings"
	"testing"

	"github.com/msaf1980/go-matcher/pkg/items"
	"github.com/msaf1980/go-matcher/pkg/tests"
)

// benchmark for string miss
var (
	globGready_StringMiss = "sys*tgicabcdERt*ltem"
	pathGready_StringMiss = "sysSKIPSKIPSKIPSKIP_tgicabcdert_SKIPSKIPSKIPSKIPSKIP_tgicabcdeRt_gltem"
)

func BenchmarkGready_StringMiss_GGlob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ParseMust(globGready_StringMiss)

		if g.Match(pathGready_StringMiss) {
			b.Fatal(pathGready_StringMiss)
		}
	}
}

// becnmark for suffix optimization
func BenchmarkGready_StringMiss_Tree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		var buf strings.Builder
		buf.Grow(len(globGready_StringMiss))
		_, _, err := w.Add(globGready_StringMiss, 1)
		if err != nil {
			b.Fatal(err)
		}
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_StringMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_StringMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globGready_StringMiss)
		if w.MatchString(pathGready_StringMiss) {
			b.Fatal(pathGready_StringMiss)
		}
	}
}

func BenchmarkGready_StringMiss_Tree_Precompiled(b *testing.B) {
	w := NewTree()
	var buf strings.Builder
	buf.Grow(len(globGready_StringMiss))
	g := ParseMust(globGready_StringMiss)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.AddGlob(g, 1)

		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_StringMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_StringMiss_Tree_Precompiled2(b *testing.B) {
	w := NewTree()
	var buf strings.Builder
	buf.Grow(len(globGready_StringMiss))
	_, _, err := w.Add(globGready_StringMiss, 1)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_StringMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_StringMiss_GGlob_Prealloc(b *testing.B) {
	g := ParseMust(globGready_StringMiss)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if g.Match(pathGready_StringMiss) {
			b.Fatal(pathGready_StringMiss)
		}
	}
}

func BenchmarkGready_StringMiss_Tree_Prealloc(b *testing.B) {
	w := NewTree()
	var buf strings.Builder
	buf.Grow(len(globGready_StringMiss))
	_, _, err := w.Add(globGready_StringMiss, 1)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		first.Init()
		_ = w.Match(pathGready_StringMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_StringMiss_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globGready_StringMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_StringMiss) {
			b.Fatal(pathGready_StringMiss)
		}
	}
}

// benchmark for byte miss
var (
	globGready_ByteMiss = "sys*A*ltem"
	pathGready_ByteMiss = "sysSKIPSKIPSKIPSKIP_tgicabcdert_SKIPSKIPSKIPSKIPSKIP_tgicabcdeRt_gltem"
)

func BenchmarkGready_ByteMiss_Tree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		_, _, err := w.Add(globGready_ByteMiss, 1)
		if err != nil {
			b.Fatal(err)
		}
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_ByteMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_ByteMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globGready_ByteMiss)
		if w.MatchString(pathGready_ByteMiss) {
			b.Fatal(pathGready_ByteMiss)
		}
	}
}

func _BenchmarkGready_ByteMiss_Tree_Precompiled(b *testing.B) {
	w := NewTree()
	g := ParseMust(globGready_ByteMiss)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.AddGlob(g, 1)

		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_ByteMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ByteMiss_Tree_Precompiled2(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globGready_ByteMiss, 1)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_ByteMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ByteMiss_Tree_Prealloc(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globGready_ByteMiss, 1)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		first.Init()
		_ = w.Match(pathGready_ByteMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_ByteMiss_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globGready_ByteMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_ByteMiss) {
			b.Fatal(pathGready_ByteMiss)
		}
	}
}

var (
	globGready_RuneRangesMiss_ASCII = "sys*[A-E]*ltem"
	pathGready_RuneRangesMiss_ASCII = "sysSKIPSKIPSKIPSKIP_tgicabcdert_SKIPSKIPSKIPSKIPSKIP_tgicabcdeRt_gltem"
)

// becnmark for suffix optimization
func BenchmarkGready_RuneRangesMiss_ASCII_Tree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		_, _, err := w.Add(globGready_RuneRangesMiss_ASCII, 0)
		if err != nil {
			b.Fatal(err)
		}
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_RuneRangesMiss_ASCII, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_RuneRangesMiss_ASCII_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globGready_RuneRangesMiss_ASCII)
		if w.MatchString(pathGready_RuneRangesMiss_ASCII) {
			b.Fatal(pathGready_RuneRangesMiss_ASCII)
		}
	}
}

func BenchmarkGready_RuneRangesMiss_ASCII_Tree_Precompiled(b *testing.B) {
	w := NewTree()
	g := ParseMust(globGready_RuneRangesMiss_ASCII)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.AddGlob(g, 1)

		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_RuneRangesMiss_ASCII, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_RuneRangesMiss_ASCII_Tree_Precompiled2(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globGready_RuneRangesMiss_ASCII, 1)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_RuneRangesMiss_ASCII, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_RuneRangesMiss_ASCII_Tree_Prealloc(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globGready_RuneRangesMiss_ASCII, 1)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		first.Init()
		_ = w.Match(pathGready_RuneRangesMiss_ASCII, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_RuneRangesMiss_ASCII_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globGready_RuneRangesMiss_ASCII)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_RuneRangesMiss_ASCII) {
			b.Fatal(pathGready_RuneRangesMiss_ASCII)
		}
	}
}

var (
	globGready_ListMiss = "DB*{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}_Status"
	pathGready_ListMiss = "DBCassandraSalesSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIP_NoCluster_Status"
)

func BenchmarkGready_ListMiss_Tree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		_, _, err := w.Add(globGready_ListMiss, 1)
		if err != nil {
			b.Fatal(err)
		}
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_ListMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_ListMiss_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globGready_ListMiss)
		if w.MatchString(pathGready_ListMiss) {
			b.Fatal(pathGready_ListMiss)
		}
	}
}

func BenchmarkGready_ListMiss_Tree_Precompiled(b *testing.B) {
	w := NewTree()
	g := ParseMust(globGready_ListMiss)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.AddGlob(g, 1)

		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_ListMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ListMiss_Tree_Precompiled2(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globGready_ListMiss, 1)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_ListMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ListMiss_Tree_Prealloc(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globGready_ListMiss, 1)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		first.Init()
		_ = w.Match(pathGready_ListMiss, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_ListMiss_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globGready_ListMiss)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_ListMiss) {
			b.Fatal(pathGready_ListMiss)
		}
	}
}

var (
	globGready_ListSkip = "DB*{BalanceCluster,BalanceStaging,Billing,BillingAutoTesting,BillingDocuments,BillingLoadTesting,BillingTesting,BlobStorageCluster,BusinessStatCluster,CashBoxCluster,CashLogCluster,CassandraClaims,CassandraClientTest,CassandraConnector,CassandraCore,CassandraDev,CassandraReliable,CassandraSentry,CassandraStats,CassandraTasks,CassandraTest,CassandraUsers,CassandraWeb,CoreCluster,CqlCoreCluster,CustomersCluster,LsaMetaindex,EventsCluster,QueueCluster,QueueTesting,LegacyCluster,ProductsCluster,ProductsTestingCluster,RemoteLockCluster,ReportCluster,ReviseCluster,ReviseTestingCluster,SalesCluster,SecondCluster,SecondTest,StoreCluster,UpProduction,UpTesting,WebCluster}_Status"
	pathGready_ListSkip = "DBCassandraSalesSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIP_WebClusteAr_Status"
)

func BenchmarkGready_ListSkip_Tree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		_, _, err := w.Add(globGready_ListSkip, 1)
		if err != nil {
			b.Fatal(err)
		}
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_ListSkip, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_ListSkip_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globGready_ListSkip)
		if w.MatchString(pathGready_ListSkip) {
			b.Fatal(pathGready_ListSkip)
		}
	}
}

func BenchmarkGready_ListSkip_Tree_Precompiled(b *testing.B) {
	w := NewTree()
	g := ParseMust(globGready_ListSkip)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.AddGlob(g, 1)

		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_ListSkip, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ListSkip_Tree_Precompiled2(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globGready_ListSkip, 1)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_ListSkip, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_ListSkip_Tree_Prealloc(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globGready_ListSkip, 1)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		first.Init()
		_ = w.Match(pathGready_ListSkip, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_ListSkip_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globGready_ListSkip)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_ListSkip) {
			b.Fatal(pathGready_ListSkip)
		}
	}
}

var (
	globGready_OneSkip = "DB*?web*_Status"
	pathGready_OneSkip = "DBCassandraSalesSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIPSKIP_we_Status"
)

func BenchmarkGready_OneSkip_Tree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := NewTree()
		_, _, err := w.Add(globGready_OneSkip, 1)
		if err != nil {
			b.Fatal(err)
		}
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_OneSkip, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_OneSkip_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := tests.BuildGlobRegexp(globGready_OneSkip)
		if w.MatchString(pathGready_OneSkip) {
			b.Fatal(pathGready_OneSkip)
		}
	}
}

func BenchmarkGready_OneSkip_Tree_Precompiled(b *testing.B) {
	w := NewTree()
	g := ParseMust(globGready_OneSkip)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.AddGlob(g, 1)

		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_OneSkip, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_OneSkip_Tree_Precompiled2(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globGready_OneSkip, 1)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var globs []string
		first := items.MinStore{-1}
		_ = w.Match(pathGready_OneSkip, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func BenchmarkGready_OneSkip_Tree_Prealloc(b *testing.B) {
	w := NewTree()
	_, _, err := w.Add(globGready_OneSkip, 1)
	if err != nil {
		b.Fatal(err)
	}
	globs := make([]string, 0, 4)
	first := items.MinStore{-1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globs = globs[:0]
		first.Init()
		_ = w.Match(pathGready_OneSkip, &globs, nil, &first)
		if len(globs) > 0 {
			b.Fatal(globs)
		}
	}
}

func _BenchmarkGready_OneSkip_Regex_Precompiled(b *testing.B) {
	w := tests.BuildGlobRegexp(globGready_OneSkip)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if w.MatchString(pathGready_OneSkip) {
			b.Fatal(pathGready_OneSkip)
		}
	}
}
