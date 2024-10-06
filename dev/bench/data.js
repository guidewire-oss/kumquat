window.BENCHMARK_DATA = {
  "lastUpdate": 1728230020901,
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
      }
    ]
  }
}