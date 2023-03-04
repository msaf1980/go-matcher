package gglob

// var (
// 	globsBatchHugeMoira  = tests.LoadPatterns("plain_patterns.txt")
// 	gGlobsBatchHugeMoira = parseGGlobs(globsBatchHugeMoira)
// )

// func BenchmarkBatchHuge_Moira_Tree(b *testing.B) {
// 	start := time.Now()
// 	b.ResetTimer()
// 	pathsBatchHugeMoira := generatePaths(gGlobsBatchHugeMoira, b.N)
// 	for i := 0; i < b.N; i++ {
// 		w := NewTree()
// 		for j := 0; j < len(globsBatchHugeMoira); j++ {
// 			_, _, err := w.Add(globsBatchHugeMoira[j], j)
// 			if err != nil {
// 				b.Fatal(err)
// 			}
// 		}
// 		var globs []string
// 		first := items.MinStore{-1}
// 		_ = w.Match(pathsBatchHugeMoira[i], &globs, nil, &first)
// 	}
// 	b.StopTimer()
// 	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
// 	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
// }

// func BenchmarkBatchHuge_Moira_Tree_ByParts(b *testing.B) {
// 	start := time.Now()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		w := NewTree()
// 		for j := 0; j < len(globsBatchHugeMoira); j++ {
// 			_, _, err := w.Add(globsBatchHugeMoira[j], j)
// 			if err != nil {
// 				b.Fatal(err)
// 			}
// 		}
// 		var globs []string
// 		first := items.MinStore{-1}
// 		for j := 0; j < len(pathsBatchHugeMoira); j++ {
// 			parts := PathSplit(pathsBatchHugeMoira[j])
// 			_ = w.MatchByParts(parts, &globs, nil, &first)
// 		}
// 	}
// 	b.StopTimer()
// 	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
// 	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
// }

// func BenchmarkBatchHuge_Moira_GGlob(b *testing.B) {
// 	start := time.Now()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		g := parseGGlobs(globsBatchHugeMoira)
// 		for j := 0; j < len(globsBatchHugeMoira); j++ {
// 			for k := 0; k < len(g); k++ {
// 				_ = g[k].Match(globsBatchHugeMoira[j])
// 			}
// 		}
// 	}
// 	b.StopTimer()
// 	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
// 	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
// }

// func BenchmarkBatchHuge_Moira_Tree_Precompiled(b *testing.B) {
// 	g := parseGGlobs(globsBatchHugeMoira)

// 	start := time.Now()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		w := NewTree()
// 		for j := 0; j < len(g); j++ {
// 			_, _, err := w.AddGlob(g[j], j)
// 			if err != nil {
// 				b.Fatal(err)
// 			}
// 		}
// 		var globs []string
// 		first := items.MinStore{-1}
// 		for j := 0; j < len(pathsBatchHugeMoira); j++ {
// 			_ = w.Match(pathsBatchHugeMoira[j], &globs, nil, &first)
// 		}
// 	}
// 	b.StopTimer()
// 	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
// 	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
// }

// func BenchmarkBatchHuge_Moira_Tree_Precompiled2(b *testing.B) {
// 	g := parseGGlobs(globsBatchHugeMoira)
// 	w := NewTree()
// 	for j := 0; j < len(g); j++ {
// 		_, _, err := w.AddGlob(g[j], j)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 	}

// 	start := time.Now()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		var globs []string
// 		first := items.MinStore{-1}
// 		for j := 0; j < len(pathsBatchHugeMoira); j++ {
// 			_ = w.Match(pathsBatchHugeMoira[j], &globs, nil, &first)
// 		}
// 	}
// 	b.StopTimer()
// 	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
// 	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
// }

// func BenchmarkBatchHuge_Moira_GGlob_Precompiled(b *testing.B) {
// 	g := parseGGlobs(globsBatchHugeMoira)

// 	start := time.Now()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < len(globsBatchHugeMoira); j++ {
// 			for k := 0; k < len(g); k++ {
// 				_ = g[k].Match(globsBatchHugeMoira[j])
// 			}
// 		}
// 	}
// 	b.StopTimer()
// 	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
// 	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
// }

// func BenchmarkBatchHuge_Moira_GGlob_Prealloc_ByParts(b *testing.B) {
// 	g := parseGGlobs(globsBatchHugeMoira)
// 	parts := make([]string, 8)

// 	start := time.Now()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < len(globsBatchHugeMoira); j++ {
// 			PathSplitB(globsBatchHugeMoira[j], &parts)
// 			length := len(globsBatchHugeMoira[j])
// 			for k := 0; k < len(g); k++ {
// 				_ = g[k].MatchByParts(parts, length)
// 			}
// 		}
// 	}
// 	b.StopTimer()
// 	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
// 	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
// }

// func BenchmarkBatchHuge_Moira_Tree_Prealloc(b *testing.B) {
// 	g := parseGGlobs(globsBatchHugeMoira)
// 	w := NewTree()
// 	for j := 0; j < len(g); j++ {
// 		_, _, err := w.AddGlob(g[j], j)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 	}
// 	globs := make([]string, 0, 4)
// 	first := items.MinStore{-1}

// 	start := time.Now()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < len(pathsBatchHugeMoira); j++ {
// 			globs = globs[:0]
// 			first.Init()
// 			_ = w.Match(pathsBatchHugeMoira[j], &globs, nil, &first)
// 		}
// 	}
// 	b.StopTimer()
// 	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
// 	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
// }

// func BenchmarkBatchHuge_Moira_Tree_Prealloc_ByParts(b *testing.B) {
// 	g := parseGGlobs(globsBatchHugeMoira)
// 	w := NewTree()
// 	for j := 0; j < len(g); j++ {
// 		_, _, err := w.AddGlob(g[j], j)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 	}
// 	globs := make([]string, 0, 4)
// 	first := items.MinStore{-1}

// 	parts := make([]string, 10)

// 	start := time.Now()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < len(pathsBatchHugeMoira); j++ {
// 			globs = globs[:0]
// 			first.Init()
// 			_ = PathSplitB(pathsBatchHugeMoira[j], &parts)
// 			_ = w.MatchByParts(parts, &globs, nil, &first)
// 		}
// 	}
// 	b.StopTimer()
// 	d := time.Since(start) // TODO: Golang 1.20 has b.Elapsed() method
// 	b.ReportMetric(float64(b.N*len(pathsBatchHugeMoira))/d.Seconds(), "match/s")
// }
