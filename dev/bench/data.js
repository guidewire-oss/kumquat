window.BENCHMARK_DATA = {
  "lastUpdate": 1728327756610,
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
          "id": "0a9dfee4a3c07b206529e994c0a12178a8d05940",
          "message": "👷 Report code coverage in GitHub Pages",
          "timestamp": "2024-10-06T22:26:53-03:00",
          "tree_id": "374465c1245c2f75233be3bd3eaa4aafa318b6fc",
          "url": "https://github.com/guidewire-oss/kumquat/commit/0a9dfee4a3c07b206529e994c0a12178a8d05940"
        },
        "date": 1728264502322,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30388,
            "unit": "ns/op",
            "extra": "398895 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 30004,
            "unit": "ns/op",
            "extra": "394933 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25809,
            "unit": "ns/op",
            "extra": "466725 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1425577636,
            "unit": "ns/op",
            "extra": "8 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "dobson@softwarepunk.com",
            "name": "James Dobson",
            "username": "jamesdobson"
          },
          "distinct": true,
          "id": "b2a8186600815435102779c6a0390bcea4178ef6",
          "message": ":arrow_up: Bump github.com/mattn/go-sqlite3 from 1.14.23 to 1.14.24\n\nBumps [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) from 1.14.23 to 1.14.24.\n- [Release notes](https://github.com/mattn/go-sqlite3/releases)\n- [Commits](https://github.com/mattn/go-sqlite3/compare/v1.14.23...v1.14.24)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/mattn/go-sqlite3\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-10-07T15:59:53-03:00",
          "tree_id": "f02e27558f4e27e0f47ce410ab1c213e8562d68a",
          "url": "https://github.com/guidewire-oss/kumquat/commit/b2a8186600815435102779c6a0390bcea4178ef6"
        },
        "date": 1728327756111,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30463,
            "unit": "ns/op",
            "extra": "393195 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 30504,
            "unit": "ns/op",
            "extra": "390364 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25968,
            "unit": "ns/op",
            "extra": "463660 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1416133344,
            "unit": "ns/op",
            "extra": "8 times\n4 procs"
          }
        ]
      }
    ]
  }
}