window.BENCHMARK_DATA = {
  "lastUpdate": 1728249829522,
  "repoUrl": "https://github.com/guidewire-oss/kumquat",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "name": "guidewire-oss",
            "username": "guidewire-oss"
          },
          "committer": {
            "name": "guidewire-oss",
            "username": "guidewire-oss"
          },
          "id": "354f6c33bcad6ba509fdfaa4ae1dfcc648201102",
          "message": "Simple benchmarks with CI support for continuous benchmarking",
          "timestamp": "2024-10-05T21:54:32Z",
          "url": "https://github.com/guidewire-oss/kumquat/pull/21/commits/354f6c33bcad6ba509fdfaa4ae1dfcc648201102"
        },
        "date": 1728230019977,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30378,
            "unit": "ns/op",
            "extra": "383454 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 30420,
            "unit": "ns/op",
            "extra": "395436 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 26053,
            "unit": "ns/op",
            "extra": "465858 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1431688576,
            "unit": "ns/op",
            "extra": "8 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "dobson@softwarepunk.com",
            "name": "James Dobson",
            "username": "jamesdobson"
          },
          "committer": {
            "email": "dobson@softwarepunk.com",
            "name": "James Dobson",
            "username": "jamesdobson"
          },
          "distinct": true,
          "id": "dc483d51cd0639c3ddffd272b759fccc96e7938d",
          "message": "👷 Add CI build for continuous benchmarking",
          "timestamp": "2024-10-06T18:20:58-03:00",
          "tree_id": "d9e0c27191aaede126fd3eefa5613023bfbbdad7",
          "url": "https://github.com/guidewire-oss/kumquat/commit/dc483d51cd0639c3ddffd272b759fccc96e7938d"
        },
        "date": 1728249828594,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29681,
            "unit": "ns/op",
            "extra": "405544 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29776,
            "unit": "ns/op",
            "extra": "398517 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25375,
            "unit": "ns/op",
            "extra": "470068 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1426820522,
            "unit": "ns/op",
            "extra": "8 times\n4 procs"
          }
        ]
      }
    ]
  }
}