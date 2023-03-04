window.BENCHMARK_DATA = {
  "lastUpdate": 1677943611626,
  "repoUrl": "https://github.com/msaf1980/go-matcher",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "msaf1980@gmail.com",
            "name": "Michail Safronov",
            "username": "msaf1980"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "0f304decce252b5d6f61d3b0751d5f2b2f8a1d53",
          "message": "tests: try to compare benchmarks in CI (#14)",
          "timestamp": "2023-03-04T20:05:37+05:00",
          "tree_id": "e6c9480f37b5f6cc5e75af08f2512294ee4e53bc",
          "url": "https://github.com/msaf1980/go-matcher/commit/0f304decce252b5d6f61d3b0751d5f2b2f8a1d53"
        },
        "date": 1677943610597,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkBatchLarge_List_Tree",
            "value": 226281,
            "unit": "ns/op\t    141413 match/s\t   92753 B/op\t     473 allocs/op",
            "extra": "4875 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_Tree_ByParts",
            "value": 230930,
            "unit": "ns/op\t    138564 match/s\t   95892 B/op\t     505 allocs/op",
            "extra": "4930 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_GGlob",
            "value": 202025,
            "unit": "ns/op\t    158388 match/s\t   77629 B/op\t     354 allocs/op",
            "extra": "5588 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_Tree_Precompiled",
            "value": 26092,
            "unit": "ns/op\t   1226393 match/s\t   13293 B/op\t     118 allocs/op",
            "extra": "45577 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_Tree_Precompiled2",
            "value": 3986,
            "unit": "ns/op\t   8027060 match/s\t     504 B/op\t       6 allocs/op",
            "extra": "299070 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_GGlob_Precompiled",
            "value": 12807,
            "unit": "ns/op\t   2498507 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "90266 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_GGlob_Prealloc_ByParts",
            "value": 4075,
            "unit": "ns/op\t   7852181 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "307239 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_Tree_Prealloc",
            "value": 7785,
            "unit": "ns/op\t   4110572 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "155643 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_Tree_Prealloc_ByParts",
            "value": 8305,
            "unit": "ns/op\t   3853179 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "139360 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree",
            "value": 97947,
            "unit": "ns/op\t    163347 match/s\t   40644 B/op\t     231 allocs/op",
            "extra": "12388 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree_ByParts",
            "value": 97800,
            "unit": "ns/op\t    163591 match/s\t   42246 B/op\t     247 allocs/op",
            "extra": "12370 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_GGlob",
            "value": 84177,
            "unit": "ns/op\t    190067 match/s\t   33781 B/op\t     158 allocs/op",
            "extra": "14372 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree_Precompiled",
            "value": 11303,
            "unit": "ns/op\t   1415494 match/s\t    6456 B/op\t      73 allocs/op",
            "extra": "103952 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree_Precompiled2",
            "value": 4350,
            "unit": "ns/op\t   3678427 match/s\t     504 B/op\t       6 allocs/op",
            "extra": "263196 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_GGlob_Precompiled",
            "value": 2940,
            "unit": "ns/op\t   5441968 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "408115 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_GGlob_Prealloc_ByParts",
            "value": 1234,
            "unit": "ns/op\t  12968774 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "993475 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree_Prealloc",
            "value": 3648,
            "unit": "ns/op\t   4385968 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "324784 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree_Prealloc_ByParts",
            "value": 4001,
            "unit": "ns/op\t   3998487 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "302352 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree",
            "value": 31132,
            "unit": "ns/op\t    706630 match/s\t   16123 B/op\t     204 allocs/op",
            "extra": "38916 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree_ByParts",
            "value": 32141,
            "unit": "ns/op\t    684456 match/s\t   17100 B/op\t     226 allocs/op",
            "extra": "37293 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_GGlob",
            "value": 17767,
            "unit": "ns/op\t   1238184 match/s\t    7232 B/op\t     125 allocs/op",
            "extra": "67774 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree_Precompiled",
            "value": 13729,
            "unit": "ns/op\t   1602347 match/s\t    8078 B/op\t      79 allocs/op",
            "extra": "86984 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree_Precompiled2",
            "value": 4224,
            "unit": "ns/op\t   5208589 match/s\t    1016 B/op\t       7 allocs/op",
            "extra": "271208 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_GGlob_Precompiled",
            "value": 3595,
            "unit": "ns/op\t   6120151 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "329713 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_GGlob_Prealloc_ByParts",
            "value": 1727,
            "unit": "ns/op\t  12740178 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "700498 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree_Prealloc",
            "value": 3321,
            "unit": "ns/op\t   6624258 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "361567 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree_Prealloc_ByParts",
            "value": 3256,
            "unit": "ns/op\t   6755926 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "362618 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_GGlob",
            "value": 1041,
            "unit": "ns/op\t     456 B/op\t      10 allocs/op",
            "extra": "1153063 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Tree",
            "value": 1561,
            "unit": "ns/op\t     680 B/op\t      15 allocs/op",
            "extra": "738268 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Regex",
            "value": 10952,
            "unit": "ns/op\t    7549 B/op\t      57 allocs/op",
            "extra": "107776 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Tree_Precompiled",
            "value": 173,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "6948398 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Tree_Precompiled2",
            "value": 153.5,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "7816605 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_GGlob_Prealloc",
            "value": 99.24,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12108334 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Tree_Prealloc",
            "value": 112.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10575229 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Regex_Precompiled",
            "value": 2255,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "542691 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Tree",
            "value": 1384,
            "unit": "ns/op\t     592 B/op\t      13 allocs/op",
            "extra": "757750 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Regex",
            "value": 8720,
            "unit": "ns/op\t    4901 B/op\t      53 allocs/op",
            "extra": "120555 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Tree_Precompiled",
            "value": 150.3,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "7964032 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Tree_Precompiled2",
            "value": 128,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "9531772 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Tree_Prealloc",
            "value": 87.12,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13440889 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Regex_Precompiled",
            "value": 2035,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "591880 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Tree",
            "value": 1998,
            "unit": "ns/op\t     640 B/op\t      13 allocs/op",
            "extra": "537529 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Regex",
            "value": 8903,
            "unit": "ns/op\t    4908 B/op\t      53 allocs/op",
            "extra": "133395 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Tree_Precompiled",
            "value": 218.2,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "5444274 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Tree_Precompiled2",
            "value": 207.8,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "5815659 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Tree_Prealloc",
            "value": 165.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7237495 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Regex_Precompiled",
            "value": 2140,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "561742 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Tree",
            "value": 9531,
            "unit": "ns/op\t    2680 B/op\t      14 allocs/op",
            "extra": "123189 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Regex",
            "value": 168254,
            "unit": "ns/op\t  125136 B/op\t     421 allocs/op",
            "extra": "6160 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Tree_Precompiled",
            "value": 1683,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "713940 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Tree_Precompiled2",
            "value": 1634,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "731929 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Tree_Prealloc",
            "value": 1602,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "745502 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Regex_Precompiled",
            "value": 22617,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "52749 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Tree",
            "value": 9481,
            "unit": "ns/op\t    2680 B/op\t      14 allocs/op",
            "extra": "120858 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Regex",
            "value": 162064,
            "unit": "ns/op\t  125073 B/op\t     421 allocs/op",
            "extra": "6646 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Tree_Precompiled",
            "value": 1675,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "743083 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Tree_Precompiled2",
            "value": 1648,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "729352 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Tree_Prealloc",
            "value": 1589,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "731326 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Regex_Precompiled",
            "value": 23364,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "51097 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Tree",
            "value": 1334,
            "unit": "ns/op\t     528 B/op\t      12 allocs/op",
            "extra": "756470 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Regex",
            "value": 10427,
            "unit": "ns/op\t    6884 B/op\t      54 allocs/op",
            "extra": "115704 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Tree_Precompiled",
            "value": 164.8,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "6968479 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Tree_Precompiled2",
            "value": 145.6,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "7768576 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Tree_Prealloc",
            "value": 105.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11374250 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Regex_Precompiled",
            "value": 2135,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "553825 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Tree",
            "value": 16541,
            "unit": "ns/op\t    7567 B/op\t      49 allocs/op",
            "extra": "70167 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Regex",
            "value": 159778,
            "unit": "ns/op\t  138480 B/op\t     400 allocs/op",
            "extra": "6274 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Tree_Precompiled",
            "value": 2095,
            "unit": "ns/op\t    1752 B/op\t      21 allocs/op",
            "extra": "525020 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Tree_Precompiled2",
            "value": 358.2,
            "unit": "ns/op\t      24 B/op\t       2 allocs/op",
            "extra": "3314701 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Tree_Prealloc",
            "value": 239,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5134972 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Regex_Precompiled",
            "value": 9534,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "125773 times\n2 procs"
          },
          {
            "name": "Benchmark_GGlob_ASCII",
            "value": 3238,
            "unit": "ns/op\t    1514 B/op\t      22 allocs/op",
            "extra": "352627 times\n2 procs"
          },
          {
            "name": "Benchmark_Regex_ASCII",
            "value": 34708,
            "unit": "ns/op\t   30561 B/op\t     116 allocs/op",
            "extra": "34844 times\n2 procs"
          },
          {
            "name": "Benchmark_GGlob_ASCII_Precompiled",
            "value": 111.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10664649 times\n2 procs"
          },
          {
            "name": "Benchmark_Regex_ASCII_Precompiled",
            "value": 2656,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "453657 times\n2 procs"
          },
          {
            "name": "Benchmark_PathSplit",
            "value": 161.5,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7425078 times\n2 procs"
          },
          {
            "name": "Benchmark_PathSplitB",
            "value": 158.8,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "7575057 times\n2 procs"
          },
          {
            "name": "Benchmark_PathSplitB_Prealloc",
            "value": 69.91,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17392785 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Any_Miss",
            "value": 673,
            "unit": "ns/op\t     224 B/op\t       5 allocs/op",
            "extra": "1758018 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Any_Miss_Regex",
            "value": 9929,
            "unit": "ns/op\t    6883 B/op\t      54 allocs/op",
            "extra": "113440 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Any_Miss_Precompiled",
            "value": 61.85,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19169488 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Any_Miss_Regex_Precompiled",
            "value": 2145,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "546926 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Byte_Miss",
            "value": 681.7,
            "unit": "ns/op\t     256 B/op\t       6 allocs/op",
            "extra": "1861483 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Byte_Miss_Regex",
            "value": 8648,
            "unit": "ns/op\t    3892 B/op\t      42 allocs/op",
            "extra": "134206 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Byte_Miss_Precompiled",
            "value": 57.43,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21270552 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Byte_Miss_Regex_Precompiled",
            "value": 3661,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "335074 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Rune_Miss",
            "value": 832.1,
            "unit": "ns/op\t     261 B/op\t       7 allocs/op",
            "extra": "1461220 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Rune_Miss_Regex",
            "value": 8815,
            "unit": "ns/op\t    3897 B/op\t      42 allocs/op",
            "extra": "136051 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Rune_Miss_Precompiled",
            "value": 138,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8949243 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Rune_Miss_Regex_Precompiled",
            "value": 3755,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "314870 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_ASCII",
            "value": 1574,
            "unit": "ns/op\t     344 B/op\t       6 allocs/op",
            "extra": "673912 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_ASCII_Regex",
            "value": 11819,
            "unit": "ns/op\t    6674 B/op\t      57 allocs/op",
            "extra": "102912 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_ASCII_Precompiled",
            "value": 312.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3551629 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_ASCII_Regex_Precompiled",
            "value": 3803,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "327643 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_Unicode",
            "value": 3124,
            "unit": "ns/op\t     560 B/op\t      12 allocs/op",
            "extra": "359138 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_Unicode_Regex",
            "value": 12433,
            "unit": "ns/op\t    6863 B/op\t      60 allocs/op",
            "extra": "92964 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_Unicode_Precompiled",
            "value": 1593,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "710082 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_Unicode_Regex_Precompiled",
            "value": 3901,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "310730 times\n2 procs"
          },
          {
            "name": "Benchmark_Size_Max",
            "value": 15310,
            "unit": "ns/op\t    7952 B/op\t      96 allocs/op",
            "extra": "80083 times\n2 procs"
          },
          {
            "name": "Benchmark_Size_Max_Regex",
            "value": 135541,
            "unit": "ns/op\t  190952 B/op\t     430 allocs/op",
            "extra": "7471 times\n2 procs"
          },
          {
            "name": "Benchmark_Size_Max_Precompiled",
            "value": 3.512,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "332314334 times\n2 procs"
          },
          {
            "name": "Benchmark_Size_Max_Regex_Precompiled",
            "value": 5.496,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "224522536 times\n2 procs"
          },
          {
            "name": "Benchmark_Suffix_Miss_ASCII",
            "value": 3854,
            "unit": "ns/op\t    1968 B/op\t      26 allocs/op",
            "extra": "297630 times\n2 procs"
          },
          {
            "name": "Benchmark_Suffix_Miss_ASCII_Regex",
            "value": 112130,
            "unit": "ns/op\t   47787 B/op\t     144 allocs/op",
            "extra": "9798 times\n2 procs"
          },
          {
            "name": "Benchmark_Suffix_Miss_ASCII_Precompiled",
            "value": 15.69,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "73780645 times\n2 procs"
          },
          {
            "name": "Benchmark_Suffix_Miss_ASCII_Regex_Precompiled",
            "value": 73629,
            "unit": "ns/op\t       2 B/op\t       0 allocs/op",
            "extra": "16596 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_ASCII",
            "value": 4952,
            "unit": "ns/op\t    1936 B/op\t      26 allocs/op",
            "extra": "220440 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_ASCII_Regex",
            "value": 111987,
            "unit": "ns/op\t   47755 B/op\t     144 allocs/op",
            "extra": "9273 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_ASCII_Precompiled",
            "value": 915.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1300104 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_ASCII_Regex_Precompiled",
            "value": 72226,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16572 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_Unicode",
            "value": 5721,
            "unit": "ns/op\t    2416 B/op\t      26 allocs/op",
            "extra": "201810 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_Unicode_Regex",
            "value": 124779,
            "unit": "ns/op\t   50446 B/op\t     166 allocs/op",
            "extra": "9552 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_Unicode_Precompiled",
            "value": 876.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1356140 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_Unicode_Regex_Precompiled",
            "value": 80909,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14761 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Suffix_Miss_Unicode",
            "value": 4707,
            "unit": "ns/op\t    2416 B/op\t      26 allocs/op",
            "extra": "248582 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Suffix_Miss_Unicode_Regex",
            "value": 122480,
            "unit": "ns/op\t   50383 B/op\t     166 allocs/op",
            "extra": "8354 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Suffix_Miss_Unicode_Precompiled",
            "value": 15.48,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "78966303 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Suffix_Miss_Unicode_Regex_Precompiled",
            "value": 82496,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14409 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_StringList_ASCII",
            "value": 17609,
            "unit": "ns/op\t    4104 B/op\t       9 allocs/op",
            "extra": "64424 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_StringList_ASCII_Regex",
            "value": 317403,
            "unit": "ns/op\t  126165 B/op\t     431 allocs/op",
            "extra": "3291 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_StringList_ASCII_Precompiled",
            "value": 9650,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "125875 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_StringList_Regex_Precompiled",
            "value": 168411,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7318 times\n2 procs"
          },
          {
            "name": "Benchmark_StringList03_Precompiled",
            "value": 19.73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "60790599 times\n2 procs"
          },
          {
            "name": "Benchmark_StringList05_Precompiled",
            "value": 21.06,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "58726512 times\n2 procs"
          },
          {
            "name": "Benchmark_StringList10_Precompiled",
            "value": 27.61,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "46221302 times\n2 procs"
          },
          {
            "name": "Benchmark_StringList50_Precompiled",
            "value": 27.37,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "46762562 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_GlobTree_Add",
            "value": 23382,
            "unit": "ns/op\t   10109 B/op\t     185 allocs/op",
            "extra": "49664 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_GlobTree_Add_Cached",
            "value": 10911,
            "unit": "ns/op\t    5424 B/op\t      93 allocs/op",
            "extra": "116178 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_GlobTree",
            "value": 28084,
            "unit": "ns/op\t    4104 B/op\t     160 allocs/op",
            "extra": "41928 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_GlobTree_Prealloc",
            "value": 16442,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "73482 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_Glob_Parse",
            "value": 12487,
            "unit": "ns/op\t    5296 B/op\t     104 allocs/op",
            "extra": "94290 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_Glob_Prealloc",
            "value": 13719,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "85962 times\n2 procs"
          },
          {
            "name": "BenchmarkPathTagsMap",
            "value": 1306,
            "unit": "ns/op\t     360 B/op\t       3 allocs/op",
            "extra": "858591 times\n2 procs"
          },
          {
            "name": "BenchmarkPathTags",
            "value": 1050,
            "unit": "ns/op\t     280 B/op\t       2 allocs/op",
            "extra": "1000000 times\n2 procs"
          },
          {
            "name": "BenchmarkByte_Find_Unicode",
            "value": 22.54,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "53160799 times\n2 procs"
          },
          {
            "name": "BenchmarkByte_Find_ASCII",
            "value": 16.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "73727440 times\n2 procs"
          },
          {
            "name": "BenchmarkRune_Find_Unicode",
            "value": 44.07,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "27913249 times\n2 procs"
          },
          {
            "name": "BenchmarkRune_Find_ASCII",
            "value": 17.11,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "69095460 times\n2 procs"
          },
          {
            "name": "BenchmarkRunesRanges_Find_Unicode",
            "value": 679.3,
            "unit": "ns/op\t     176 B/op\t       6 allocs/op",
            "extra": "1772080 times\n2 procs"
          },
          {
            "name": "BenchmarkRunesRanges_Find_ASCII",
            "value": 101.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11914971 times\n2 procs"
          },
          {
            "name": "BenchmarkString_Find_ASCII",
            "value": 28.43,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "42044497 times\n2 procs"
          },
          {
            "name": "BenchmarkString_Find_Unicode",
            "value": 35.08,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "33903702 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_ASCII",
            "value": 780.8,
            "unit": "ns/op\t     152 B/op\t       3 allocs/op",
            "extra": "1533645 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_ASCII_Prealloc",
            "value": 77.61,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15009394 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_ASCII_Skip",
            "value": 319.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3766228 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_ASCII",
            "value": 468.4,
            "unit": "ns/op\t     152 B/op\t       3 allocs/op",
            "extra": "2493782 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_ASCII_Prealloc",
            "value": 43.35,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "28390852 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_ASCII_Skip",
            "value": 44.37,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "26931096 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_Unicode",
            "value": 618.2,
            "unit": "ns/op\t     168 B/op\t       3 allocs/op",
            "extra": "1929780 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_Unicode_Prealloc",
            "value": 107.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11267266 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_Unicode_Skip",
            "value": 110.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10823112 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_Unicode",
            "value": 543.8,
            "unit": "ns/op\t     168 B/op\t       3 allocs/op",
            "extra": "2134628 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_Unicode_Prealloc",
            "value": 51.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23227575 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_Unicode_Skip",
            "value": 51.84,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22818356 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_ASCII_Miss",
            "value": 40.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29505765 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_ASCII_Miss_Skip",
            "value": 45.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "26896206 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_Unicode_Miss",
            "value": 52.65,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22815709 times\n2 procs"
          },
          {
            "name": "Benchmark_Contains_ASCIISet",
            "value": 12.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "95616554 times\n2 procs"
          },
          {
            "name": "Benchmark_Contains_ASCIISet_Prealloc",
            "value": 0.6222,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexASCII_ASCIISet",
            "value": 125.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8993103 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexASCII_ASCIISet_Prealloc",
            "value": 114.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10337866 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexASCII_StringsAny",
            "value": 112.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10307383 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_ASCII_RuneSet_Large",
            "value": 135.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8867187 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_ASCII_RuneSet_Prealloc",
            "value": 118,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10135184 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_ASCII_StringsAny",
            "value": 114.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10665207 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_Unicode_RuneSet",
            "value": 315.3,
            "unit": "ns/op\t      24 B/op\t       2 allocs/op",
            "extra": "3794157 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_Unicode_RuneSet_Prealloc",
            "value": 191.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6224529 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_Unicode_StringsAny",
            "value": 1111,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_Large_RuneSet",
            "value": 537.4,
            "unit": "ns/op\t      24 B/op\t       2 allocs/op",
            "extra": "2270175 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_Large_RuneSet_Prealloc",
            "value": 195.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6215064 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_StringsAny_Large",
            "value": 959.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1274620 times\n2 procs"
          },
          {
            "name": "Benchmark_String_SkipRunes",
            "value": 22.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "54283428 times\n2 procs"
          },
          {
            "name": "Benchmark_String_SkipRunesEmpty",
            "value": 3.848,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "307685215 times\n2 procs"
          }
        ]
      }
    ]
  }
}