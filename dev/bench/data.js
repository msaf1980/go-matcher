window.BENCHMARK_DATA = {
  "lastUpdate": 1677941045396,
  "repoUrl": "https://github.com/msaf1980/go-matcher",
  "entries": {
    "My Project Go Benchmark": [
      {
        "commit": {
          "author": {
            "name": "msaf1980",
            "username": "msaf1980"
          },
          "committer": {
            "name": "msaf1980",
            "username": "msaf1980"
          },
          "id": "17387b7b192d5e829e0fb0badadadbee36640541",
          "message": "tests: try to compare benchmarks in CI",
          "timestamp": "2023-02-21T10:50:01Z",
          "url": "https://github.com/msaf1980/go-matcher/pull/13/commits/17387b7b192d5e829e0fb0badadadbee36640541"
        },
        "date": 1677941044401,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkBatchLarge_List_Tree",
            "value": 209417,
            "unit": "ns/op\t    152799 match/s\t   92762 B/op\t     473 allocs/op",
            "extra": "5618 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_Tree_ByParts",
            "value": 215479,
            "unit": "ns/op\t    148497 match/s\t   95868 B/op\t     505 allocs/op",
            "extra": "5097 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_GGlob",
            "value": 183754,
            "unit": "ns/op\t    174138 match/s\t   77616 B/op\t     354 allocs/op",
            "extra": "5926 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_Tree_Precompiled",
            "value": 24869,
            "unit": "ns/op\t   1286689 match/s\t   13292 B/op\t     118 allocs/op",
            "extra": "46548 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_Tree_Precompiled2",
            "value": 3643,
            "unit": "ns/op\t   8784306 match/s\t     504 B/op\t       6 allocs/op",
            "extra": "313209 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_GGlob_Precompiled",
            "value": 11800,
            "unit": "ns/op\t   2711817 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "100437 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_GGlob_Prealloc_ByParts",
            "value": 3643,
            "unit": "ns/op\t   8783873 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "332308 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_Tree_Prealloc",
            "value": 6821,
            "unit": "ns/op\t   4691241 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "175909 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchLarge_List_Tree_Prealloc_ByParts",
            "value": 7634,
            "unit": "ns/op\t   4191848 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "160233 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree",
            "value": 85628,
            "unit": "ns/op\t    186846 match/s\t   40640 B/op\t     231 allocs/op",
            "extra": "13485 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree_ByParts",
            "value": 90979,
            "unit": "ns/op\t    175858 match/s\t   42243 B/op\t     247 allocs/op",
            "extra": "13663 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_GGlob",
            "value": 77114,
            "unit": "ns/op\t    207478 match/s\t   33787 B/op\t     158 allocs/op",
            "extra": "15912 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree_Precompiled",
            "value": 10340,
            "unit": "ns/op\t   1547406 match/s\t    6456 B/op\t      73 allocs/op",
            "extra": "115774 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree_Precompiled2",
            "value": 4029,
            "unit": "ns/op\t   3971385 match/s\t     504 B/op\t       6 allocs/op",
            "extra": "282658 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_GGlob_Precompiled",
            "value": 2736,
            "unit": "ns/op\t   5848487 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "443431 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_GGlob_Prealloc_ByParts",
            "value": 1147,
            "unit": "ns/op\t  13952445 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "981447 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree_Prealloc",
            "value": 3335,
            "unit": "ns/op\t   4798209 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "368188 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_List_Tree_Prealloc_ByParts",
            "value": 3789,
            "unit": "ns/op\t   4222783 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "325094 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree",
            "value": 26906,
            "unit": "ns/op\t    817624 match/s\t   16124 B/op\t     204 allocs/op",
            "extra": "40666 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree_ByParts",
            "value": 29397,
            "unit": "ns/op\t    748337 match/s\t   17099 B/op\t     226 allocs/op",
            "extra": "40682 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_GGlob",
            "value": 16188,
            "unit": "ns/op\t   1358987 match/s\t    7232 B/op\t     125 allocs/op",
            "extra": "76147 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree_Precompiled",
            "value": 13097,
            "unit": "ns/op\t   1679800 match/s\t    8077 B/op\t      79 allocs/op",
            "extra": "92774 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree_Precompiled2",
            "value": 3791,
            "unit": "ns/op\t   5803783 match/s\t    1016 B/op\t       7 allocs/op",
            "extra": "317281 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_GGlob_Precompiled",
            "value": 3233,
            "unit": "ns/op\t   6804387 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "376894 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_GGlob_Prealloc_ByParts",
            "value": 1543,
            "unit": "ns/op\t  14261172 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "752170 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree_Prealloc",
            "value": 3025,
            "unit": "ns/op\t   7271944 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "408222 times\n2 procs"
          },
          {
            "name": "BenchmarkBatch_Moira_Tree_Prealloc_ByParts",
            "value": 3138,
            "unit": "ns/op\t   7010407 match/s\t       0 B/op\t       0 allocs/op",
            "extra": "413257 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_GGlob",
            "value": 935.2,
            "unit": "ns/op\t     456 B/op\t      10 allocs/op",
            "extra": "1259179 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Tree",
            "value": 1417,
            "unit": "ns/op\t     680 B/op\t      15 allocs/op",
            "extra": "778608 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Regex",
            "value": 9963,
            "unit": "ns/op\t    7550 B/op\t      57 allocs/op",
            "extra": "119068 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Tree_Precompiled",
            "value": 158.7,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "7700164 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Tree_Precompiled2",
            "value": 141.1,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "8663341 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_GGlob_Prealloc",
            "value": 92.17,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13443387 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Tree_Prealloc",
            "value": 102.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11279304 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_StringMiss_Regex_Precompiled",
            "value": 2102,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "597123 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Tree",
            "value": 1373,
            "unit": "ns/op\t     592 B/op\t      13 allocs/op",
            "extra": "811326 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Regex",
            "value": 7836,
            "unit": "ns/op\t    4900 B/op\t      53 allocs/op",
            "extra": "148546 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Tree_Precompiled",
            "value": 136,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "8791876 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Tree_Precompiled2",
            "value": 118.8,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "10142474 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Tree_Prealloc",
            "value": 82.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15024825 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ByteMiss_Regex_Precompiled",
            "value": 1989,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "550581 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Tree",
            "value": 1906,
            "unit": "ns/op\t     640 B/op\t      13 allocs/op",
            "extra": "646904 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Regex",
            "value": 8094,
            "unit": "ns/op\t    4906 B/op\t      53 allocs/op",
            "extra": "147806 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Tree_Precompiled",
            "value": 201,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "5828103 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Tree_Precompiled2",
            "value": 188.5,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "6438122 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Tree_Prealloc",
            "value": 153.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8132158 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_RuneRangesMiss_ASCII_Regex_Precompiled",
            "value": 1923,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "610680 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Tree",
            "value": 8901,
            "unit": "ns/op\t    2680 B/op\t      14 allocs/op",
            "extra": "137896 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Regex",
            "value": 152914,
            "unit": "ns/op\t  125080 B/op\t     421 allocs/op",
            "extra": "7958 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Tree_Precompiled",
            "value": 1513,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "752127 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Tree_Precompiled2",
            "value": 1481,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "823015 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Tree_Prealloc",
            "value": 1480,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "816080 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListMiss_Regex_Precompiled",
            "value": 20985,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "58438 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Tree",
            "value": 8719,
            "unit": "ns/op\t    2680 B/op\t      14 allocs/op",
            "extra": "131170 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Regex",
            "value": 151561,
            "unit": "ns/op\t  125125 B/op\t     421 allocs/op",
            "extra": "7644 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Tree_Precompiled",
            "value": 1521,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "811764 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Tree_Precompiled2",
            "value": 1505,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "799194 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Tree_Prealloc",
            "value": 1452,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "844944 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_ListSkip_Regex_Precompiled",
            "value": 21769,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "54859 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Tree",
            "value": 1239,
            "unit": "ns/op\t     528 B/op\t      12 allocs/op",
            "extra": "815166 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Regex",
            "value": 9446,
            "unit": "ns/op\t    6883 B/op\t      54 allocs/op",
            "extra": "128578 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Tree_Precompiled",
            "value": 155.7,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "7927023 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Tree_Precompiled2",
            "value": 130.9,
            "unit": "ns/op\t       8 B/op\t       1 allocs/op",
            "extra": "9171368 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Tree_Prealloc",
            "value": 96.93,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12569415 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_OneSkip_Regex_Precompiled",
            "value": 2004,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "642345 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Tree",
            "value": 14481,
            "unit": "ns/op\t    7568 B/op\t      49 allocs/op",
            "extra": "73869 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Regex",
            "value": 146218,
            "unit": "ns/op\t  138509 B/op\t     400 allocs/op",
            "extra": "7477 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Tree_Precompiled",
            "value": 1906,
            "unit": "ns/op\t    1752 B/op\t      21 allocs/op",
            "extra": "558717 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Tree_Precompiled2",
            "value": 323.1,
            "unit": "ns/op\t      24 B/op\t       2 allocs/op",
            "extra": "3737786 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Tree_Prealloc",
            "value": 211.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5556692 times\n2 procs"
          },
          {
            "name": "BenchmarkGready_List_Regex_Precompiled",
            "value": 8586,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "142777 times\n2 procs"
          },
          {
            "name": "Benchmark_GGlob_ASCII",
            "value": 3028,
            "unit": "ns/op\t    1514 B/op\t      22 allocs/op",
            "extra": "394197 times\n2 procs"
          },
          {
            "name": "Benchmark_Regex_ASCII",
            "value": 32019,
            "unit": "ns/op\t   30550 B/op\t     116 allocs/op",
            "extra": "37470 times\n2 procs"
          },
          {
            "name": "Benchmark_GGlob_ASCII_Precompiled",
            "value": 101.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11969570 times\n2 procs"
          },
          {
            "name": "Benchmark_Regex_ASCII_Precompiled",
            "value": 2369,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "500347 times\n2 procs"
          },
          {
            "name": "Benchmark_PathSplit",
            "value": 141.3,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8191206 times\n2 procs"
          },
          {
            "name": "Benchmark_PathSplitB",
            "value": 143,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "8411362 times\n2 procs"
          },
          {
            "name": "Benchmark_PathSplitB_Prealloc",
            "value": 62.09,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19544143 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Any_Miss",
            "value": 630.8,
            "unit": "ns/op\t     224 B/op\t       5 allocs/op",
            "extra": "1995178 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Any_Miss_Regex",
            "value": 9271,
            "unit": "ns/op\t    6883 B/op\t      54 allocs/op",
            "extra": "132840 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Any_Miss_Precompiled",
            "value": 54.61,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21938821 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Any_Miss_Regex_Precompiled",
            "value": 1962,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "578925 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Byte_Miss",
            "value": 586.2,
            "unit": "ns/op\t     256 B/op\t       6 allocs/op",
            "extra": "2050069 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Byte_Miss_Regex",
            "value": 8038,
            "unit": "ns/op\t    3890 B/op\t      42 allocs/op",
            "extra": "142766 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Byte_Miss_Precompiled",
            "value": 50.93,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23377704 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Byte_Miss_Regex_Precompiled",
            "value": 3304,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "373042 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Rune_Miss",
            "value": 759.6,
            "unit": "ns/op\t     261 B/op\t       7 allocs/op",
            "extra": "1527198 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Rune_Miss_Regex",
            "value": 8043,
            "unit": "ns/op\t    3900 B/op\t      42 allocs/op",
            "extra": "145746 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Rune_Miss_Precompiled",
            "value": 126.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9168538 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Rune_Miss_Regex_Precompiled",
            "value": 3510,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "357972 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_ASCII",
            "value": 1432,
            "unit": "ns/op\t     344 B/op\t       6 allocs/op",
            "extra": "796896 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_ASCII_Regex",
            "value": 10654,
            "unit": "ns/op\t    6675 B/op\t      57 allocs/op",
            "extra": "110263 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_ASCII_Precompiled",
            "value": 288.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3942499 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_ASCII_Regex_Precompiled",
            "value": 3493,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "346626 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_Unicode",
            "value": 2758,
            "unit": "ns/op\t     560 B/op\t      12 allocs/op",
            "extra": "391254 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_Unicode_Regex",
            "value": 11954,
            "unit": "ns/op\t    6857 B/op\t      60 allocs/op",
            "extra": "106232 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_Unicode_Precompiled",
            "value": 1447,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "742890 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_RunesRanges_Unicode_Regex_Precompiled",
            "value": 3614,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "332997 times\n2 procs"
          },
          {
            "name": "Benchmark_Size_Max",
            "value": 13370,
            "unit": "ns/op\t    7952 B/op\t      96 allocs/op",
            "extra": "87490 times\n2 procs"
          },
          {
            "name": "Benchmark_Size_Max_Regex",
            "value": 124482,
            "unit": "ns/op\t  190952 B/op\t     430 allocs/op",
            "extra": "10000 times\n2 procs"
          },
          {
            "name": "Benchmark_Size_Max_Precompiled",
            "value": 3.181,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "389566645 times\n2 procs"
          },
          {
            "name": "Benchmark_Size_Max_Regex_Precompiled",
            "value": 4.914,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "246220863 times\n2 procs"
          },
          {
            "name": "Benchmark_Suffix_Miss_ASCII",
            "value": 3673,
            "unit": "ns/op\t    1968 B/op\t      26 allocs/op",
            "extra": "312476 times\n2 procs"
          },
          {
            "name": "Benchmark_Suffix_Miss_ASCII_Regex",
            "value": 101801,
            "unit": "ns/op\t   47810 B/op\t     144 allocs/op",
            "extra": "10000 times\n2 procs"
          },
          {
            "name": "Benchmark_Suffix_Miss_ASCII_Precompiled",
            "value": 14.71,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "80259945 times\n2 procs"
          },
          {
            "name": "Benchmark_Suffix_Miss_ASCII_Regex_Precompiled",
            "value": 66803,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17956 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_ASCII",
            "value": 4620,
            "unit": "ns/op\t    1936 B/op\t      26 allocs/op",
            "extra": "239923 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_ASCII_Regex",
            "value": 101718,
            "unit": "ns/op\t   47791 B/op\t     144 allocs/op",
            "extra": "10000 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_ASCII_Precompiled",
            "value": 882.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1425252 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_ASCII_Regex_Precompiled",
            "value": 66747,
            "unit": "ns/op\t       2 B/op\t       0 allocs/op",
            "extra": "17742 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_Unicode",
            "value": 5224,
            "unit": "ns/op\t    2416 B/op\t      26 allocs/op",
            "extra": "232129 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_Unicode_Regex",
            "value": 113764,
            "unit": "ns/op\t   50430 B/op\t     166 allocs/op",
            "extra": "10000 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_Unicode_Precompiled",
            "value": 799.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1482304 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_String_Miss_Unicode_Regex_Precompiled",
            "value": 75508,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15808 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Suffix_Miss_Unicode",
            "value": 4273,
            "unit": "ns/op\t    2416 B/op\t      26 allocs/op",
            "extra": "272649 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Suffix_Miss_Unicode_Regex",
            "value": 116167,
            "unit": "ns/op\t   50408 B/op\t     166 allocs/op",
            "extra": "8660 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Suffix_Miss_Unicode_Precompiled",
            "value": 14.48,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "79910936 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_Suffix_Miss_Unicode_Regex_Precompiled",
            "value": 73145,
            "unit": "ns/op\t       2 B/op\t       0 allocs/op",
            "extra": "15769 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_StringList_ASCII",
            "value": 16109,
            "unit": "ns/op\t    4104 B/op\t       9 allocs/op",
            "extra": "73410 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_StringList_ASCII_Regex",
            "value": 295214,
            "unit": "ns/op\t  126114 B/op\t     430 allocs/op",
            "extra": "3830 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_StringList_ASCII_Precompiled",
            "value": 9099,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "130214 times\n2 procs"
          },
          {
            "name": "Benchmark_Star_StringList_Regex_Precompiled",
            "value": 162037,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8545 times\n2 procs"
          },
          {
            "name": "Benchmark_StringList03_Precompiled",
            "value": 18.04,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "66554913 times\n2 procs"
          },
          {
            "name": "Benchmark_StringList05_Precompiled",
            "value": 19.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "61282263 times\n2 procs"
          },
          {
            "name": "Benchmark_StringList10_Precompiled",
            "value": 23.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "44878886 times\n2 procs"
          },
          {
            "name": "Benchmark_StringList50_Precompiled",
            "value": 24.47,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "44979818 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_GlobTree_Add",
            "value": 20740,
            "unit": "ns/op\t   10108 B/op\t     185 allocs/op",
            "extra": "54708 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_GlobTree_Add_Cached",
            "value": 9390,
            "unit": "ns/op\t    5424 B/op\t      93 allocs/op",
            "extra": "122794 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_GlobTree",
            "value": 25338,
            "unit": "ns/op\t    4104 B/op\t     160 allocs/op",
            "extra": "47016 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_GlobTree_Prealloc",
            "value": 15292,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "81238 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_Glob_Parse",
            "value": 11302,
            "unit": "ns/op\t    5296 B/op\t     104 allocs/op",
            "extra": "103945 times\n2 procs"
          },
          {
            "name": "Benchmark_Batch_Glob_Prealloc",
            "value": 12655,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "92059 times\n2 procs"
          },
          {
            "name": "BenchmarkPathTagsMap",
            "value": 1235,
            "unit": "ns/op\t     360 B/op\t       3 allocs/op",
            "extra": "822564 times\n2 procs"
          },
          {
            "name": "BenchmarkPathTags",
            "value": 946.2,
            "unit": "ns/op\t     280 B/op\t       2 allocs/op",
            "extra": "1253479 times\n2 procs"
          },
          {
            "name": "BenchmarkByte_Find_Unicode",
            "value": 20.24,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "57748322 times\n2 procs"
          },
          {
            "name": "BenchmarkByte_Find_ASCII",
            "value": 14.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "79485906 times\n2 procs"
          },
          {
            "name": "BenchmarkRune_Find_Unicode",
            "value": 41.09,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30420915 times\n2 procs"
          },
          {
            "name": "BenchmarkRune_Find_ASCII",
            "value": 15.48,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "77996437 times\n2 procs"
          },
          {
            "name": "BenchmarkRunesRanges_Find_Unicode",
            "value": 617.8,
            "unit": "ns/op\t     176 B/op\t       6 allocs/op",
            "extra": "1943067 times\n2 procs"
          },
          {
            "name": "BenchmarkRunesRanges_Find_ASCII",
            "value": 90.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13038927 times\n2 procs"
          },
          {
            "name": "BenchmarkString_Find_ASCII",
            "value": 26.08,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "45302628 times\n2 procs"
          },
          {
            "name": "BenchmarkString_Find_Unicode",
            "value": 31.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "38296156 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_ASCII",
            "value": 722.3,
            "unit": "ns/op\t     152 B/op\t       3 allocs/op",
            "extra": "1715584 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_ASCII_Prealloc",
            "value": 70.58,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17021064 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_ASCII_Skip",
            "value": 290.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4122772 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_ASCII",
            "value": 427.6,
            "unit": "ns/op\t     152 B/op\t       3 allocs/op",
            "extra": "2840170 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_ASCII_Prealloc",
            "value": 39.71,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "31162809 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_ASCII_Skip",
            "value": 40.58,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29241658 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_Unicode",
            "value": 564.2,
            "unit": "ns/op\t     168 B/op\t       3 allocs/op",
            "extra": "2171767 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_Unicode_Prealloc",
            "value": 96.71,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12386790 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Find_Unicode_Skip",
            "value": 100.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11845533 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_Unicode",
            "value": 492.5,
            "unit": "ns/op\t     168 B/op\t       3 allocs/op",
            "extra": "2398592 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_Unicode_Prealloc",
            "value": 46.01,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "25528738 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_Unicode_Skip",
            "value": 46.71,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "26001914 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_ASCII_Miss",
            "value": 36.14,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "34103371 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_ASCII_Miss_Skip",
            "value": 40.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30205974 times\n2 procs"
          },
          {
            "name": "BenchmarkStringList_Match_Unicode_Miss",
            "value": 48.12,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "24615256 times\n2 procs"
          },
          {
            "name": "Benchmark_Contains_ASCIISet",
            "value": 11.22,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n2 procs"
          },
          {
            "name": "Benchmark_Contains_ASCIISet_Prealloc",
            "value": 0.5673,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexASCII_ASCIISet",
            "value": 112.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10575890 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexASCII_ASCIISet_Prealloc",
            "value": 102.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11518093 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexASCII_StringsAny",
            "value": 101.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11422575 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_ASCII_RuneSet_Large",
            "value": 121.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10149525 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_ASCII_RuneSet_Prealloc",
            "value": 107.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11585694 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_ASCII_StringsAny",
            "value": 100.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11039769 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_Unicode_RuneSet",
            "value": 286,
            "unit": "ns/op\t      24 B/op\t       2 allocs/op",
            "extra": "4158285 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_Unicode_RuneSet_Prealloc",
            "value": 172.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6672642 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_Unicode_StringsAny",
            "value": 1031,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1208108 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_Large_RuneSet",
            "value": 490.8,
            "unit": "ns/op\t      24 B/op\t       2 allocs/op",
            "extra": "2443300 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_Large_RuneSet_Prealloc",
            "value": 173.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6517699 times\n2 procs"
          },
          {
            "name": "Benchmark_IndexUnicode_StringsAny_Large",
            "value": 839.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1421040 times\n2 procs"
          },
          {
            "name": "Benchmark_String_SkipRunes",
            "value": 18.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "61545038 times\n2 procs"
          },
          {
            "name": "Benchmark_String_SkipRunesEmpty",
            "value": 3.514,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "351471736 times\n2 procs"
          }
        ]
      }
    ]
  }
}