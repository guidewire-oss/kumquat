window.BENCHMARK_DATA = {
  "lastUpdate": 1731716937746,
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
          "id": "625b694de22f59b2d00087aeb5b49751063cc860",
          "message": "👷 Fix workflow failure when coverage report doesn't change",
          "timestamp": "2024-10-07T16:39:11-03:00",
          "tree_id": "4e983ef8f43c818c206f058fbb58968a812583d2",
          "url": "https://github.com/guidewire-oss/kumquat/commit/625b694de22f59b2d00087aeb5b49751063cc860"
        },
        "date": 1728330049115,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30136,
            "unit": "ns/op",
            "extra": "399273 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29937,
            "unit": "ns/op",
            "extra": "399139 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 26003,
            "unit": "ns/op",
            "extra": "460124 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1433099382,
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
          "id": "7a4668f734afed737c296ffac250dea35a34d9b5",
          "message": "👷 Run integration tests from CI",
          "timestamp": "2024-10-07T17:13:43-03:00",
          "tree_id": "8cbf643539acf1c902dec850b8451c94d6a88429",
          "url": "https://github.com/guidewire-oss/kumquat/commit/7a4668f734afed737c296ffac250dea35a34d9b5"
        },
        "date": 1728332104958,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30343,
            "unit": "ns/op",
            "extra": "396595 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29817,
            "unit": "ns/op",
            "extra": "394113 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25972,
            "unit": "ns/op",
            "extra": "468151 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1423302064,
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
          "id": "678b5a7a920068830ecc0153406ce4aff32e6019",
          "message": "👷 Run lint faster",
          "timestamp": "2024-10-07T17:48:57-03:00",
          "tree_id": "024e9292ef012efb0b3454a5218d71588d6d6e13",
          "url": "https://github.com/guidewire-oss/kumquat/commit/678b5a7a920068830ecc0153406ce4aff32e6019"
        },
        "date": 1728334222638,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29873,
            "unit": "ns/op",
            "extra": "404583 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29726,
            "unit": "ns/op",
            "extra": "400814 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25850,
            "unit": "ns/op",
            "extra": "456492 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1414068067,
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
          "id": "b3eaa9a6ba4b3aae9d6b34ff434190502d8d1ab7",
          "message": "👷 Fix permissions to allow coverage publishing",
          "timestamp": "2024-10-08T00:53:14Z",
          "tree_id": "d9bee752a101458cfdd35bb84db9530fe67754af",
          "url": "https://github.com/guidewire-oss/kumquat/commit/b3eaa9a6ba4b3aae9d6b34ff434190502d8d1ab7"
        },
        "date": 1728348887017,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30487,
            "unit": "ns/op",
            "extra": "402652 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 30143,
            "unit": "ns/op",
            "extra": "397993 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 26035,
            "unit": "ns/op",
            "extra": "454442 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1435218458,
            "unit": "ns/op",
            "extra": "7 times\n4 procs"
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
          "id": "69ed1819637faa4d82a1274d0bbc07969bd88ca7",
          "message": "👷 Publish ADRs to GitHub Pages",
          "timestamp": "2024-10-07T22:40:08-03:00",
          "tree_id": "e57450f09f0b7e4b54ac918df94fb06779623184",
          "url": "https://github.com/guidewire-oss/kumquat/commit/69ed1819637faa4d82a1274d0bbc07969bd88ca7"
        },
        "date": 1728351686086,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30237,
            "unit": "ns/op",
            "extra": "398116 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 30214,
            "unit": "ns/op",
            "extra": "392419 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 26054,
            "unit": "ns/op",
            "extra": "464565 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1433099208,
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
          "id": "c6b107b79d8e3e0374e2d9a595aad7b693c17663",
          "message": "👷 Assign publishing permissions to Log4brains workflow",
          "timestamp": "2024-10-08T01:43:27Z",
          "tree_id": "a23c1a2932e73728f277dada0e67b09c3beb502d",
          "url": "https://github.com/guidewire-oss/kumquat/commit/c6b107b79d8e3e0374e2d9a595aad7b693c17663"
        },
        "date": 1728351896961,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29964,
            "unit": "ns/op",
            "extra": "403876 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29951,
            "unit": "ns/op",
            "extra": "391968 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25978,
            "unit": "ns/op",
            "extra": "467504 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1416871885,
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
          "id": "b16a8d76845899c51029fc6159f5444e0a28b376",
          "message": "👷 Set up base URL path for ADR site",
          "timestamp": "2024-10-08T13:40:19Z",
          "tree_id": "c2b15faf23143a74ef08526e2a3a329fb34a3d7f",
          "url": "https://github.com/guidewire-oss/kumquat/commit/b16a8d76845899c51029fc6159f5444e0a28b376"
        },
        "date": 1728394903522,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30177,
            "unit": "ns/op",
            "extra": "397723 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 30006,
            "unit": "ns/op",
            "extra": "401588 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25955,
            "unit": "ns/op",
            "extra": "448214 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1402965526,
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
          "id": "b450dc793e0daeca386f7d96ce333e637a62edbf",
          "message": "👷 Serve coverage badge JSON through GitHub Pages",
          "timestamp": "2024-10-08T19:43:01Z",
          "tree_id": "02183bcef4ee4e82ed354c37242c992bdc033c26",
          "url": "https://github.com/guidewire-oss/kumquat/commit/b450dc793e0daeca386f7d96ce333e637a62edbf"
        },
        "date": 1728416674561,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29775,
            "unit": "ns/op",
            "extra": "406904 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29658,
            "unit": "ns/op",
            "extra": "398198 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25632,
            "unit": "ns/op",
            "extra": "469245 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1424987092,
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
          "id": "7b8312e9ca470c82a1da9693b40798df4cfa5166",
          "message": "💚 Fixed typo in coverage badge",
          "timestamp": "2024-10-08T19:46:48Z",
          "tree_id": "debd14c0cb37a1d8653bfbe446b46df29dc3c440",
          "url": "https://github.com/guidewire-oss/kumquat/commit/7b8312e9ca470c82a1da9693b40798df4cfa5166"
        },
        "date": 1728416897524,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30073,
            "unit": "ns/op",
            "extra": "404113 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29752,
            "unit": "ns/op",
            "extra": "400435 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25831,
            "unit": "ns/op",
            "extra": "468997 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1437100702,
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
          "id": "2ec990fcfcf9bef8006f44410a5ed45d3101d372",
          "message": "💚 Add coverage badge JSON to GitHub pages site",
          "timestamp": "2024-10-08T19:49:53Z",
          "tree_id": "a8357c07e28de67a02ea20156d7c41ee6d1a0f85",
          "url": "https://github.com/guidewire-oss/kumquat/commit/2ec990fcfcf9bef8006f44410a5ed45d3101d372"
        },
        "date": 1728417081090,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30366,
            "unit": "ns/op",
            "extra": "402963 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 30018,
            "unit": "ns/op",
            "extra": "394840 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25972,
            "unit": "ns/op",
            "extra": "458592 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1427440805,
            "unit": "ns/op",
            "extra": "8 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "56201313+amirbavand@users.noreply.github.com",
            "name": "Amir Bavand",
            "username": "amirbavand"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "6c8ee14641f5bd65776d16cb081bbc2fe3937781",
          "message": "🐛 Fixed a bug during deletion of one dependent resource\n\n* added initial implementation to resolve the bug\r\n\r\n* fixed some bugs\r\n\r\n* addd required files to implement integration test for deletion scenario\r\n\r\n* addd required files to implement integration test for deletion scenario\r\n\r\n* added required integration test\r\n\r\n* added new line to file\r\n\r\n* reset the file\r\n\r\n* delete some of log statements\r\n\r\n* fixed some lint issues\r\n\r\n* fixed some more lint issues\r\n\r\n* addressed review comments\r\n\r\n* used map instead of list to delete the resources\r\n\r\n* fixed a return issue with test",
          "timestamp": "2024-11-10T12:32:16-05:00",
          "tree_id": "e2fd24d5fa42128d0b1cdcc59aeb38c0471f8463",
          "url": "https://github.com/guidewire-oss/kumquat/commit/6c8ee14641f5bd65776d16cb081bbc2fe3937781"
        },
        "date": 1731260096808,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29750,
            "unit": "ns/op",
            "extra": "400173 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29722,
            "unit": "ns/op",
            "extra": "405238 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25604,
            "unit": "ns/op",
            "extra": "464215 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1401368523,
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
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "878cd32a3cf9b567de6553bd79e20f34c365e989",
          "message": "Merge pull request #18 from guidewire-oss/chore/update-dependencies\n\n⬆️ Bump some Kubernetes dependencies",
          "timestamp": "2024-11-10T23:35:17-05:00",
          "tree_id": "03c4a5959302f81b0a03388140ffeca71adc4d60",
          "url": "https://github.com/guidewire-oss/kumquat/commit/878cd32a3cf9b567de6553bd79e20f34c365e989"
        },
        "date": 1731299897756,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30243,
            "unit": "ns/op",
            "extra": "398774 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 30243,
            "unit": "ns/op",
            "extra": "396112 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25730,
            "unit": "ns/op",
            "extra": "467788 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1430332289,
            "unit": "ns/op",
            "extra": "7 times\n4 procs"
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
          "id": "da464f5aabab714701e35f1d5a0c7b4bb98cccdf",
          "message": ":arrow_up: Bump sigs.k8s.io/e2e-framework from 0.4.0 to 0.5.0\n\nBumps [sigs.k8s.io/e2e-framework](https://github.com/kubernetes-sigs/e2e-framework) from 0.4.0 to 0.5.0.\n- [Release notes](https://github.com/kubernetes-sigs/e2e-framework/releases)\n- [Changelog](https://github.com/kubernetes-sigs/e2e-framework/blob/main/RELEASE.md)\n- [Commits](https://github.com/kubernetes-sigs/e2e-framework/compare/v0.4.0...v0.5.0)\n\n---\nupdated-dependencies:\n- dependency-name: sigs.k8s.io/e2e-framework\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-11-10T23:45:01-05:00",
          "tree_id": "5eaa833eb4158d9d24c7ea1592708c8282ad0864",
          "url": "https://github.com/guidewire-oss/kumquat/commit/da464f5aabab714701e35f1d5a0c7b4bb98cccdf"
        },
        "date": 1731300478131,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29496,
            "unit": "ns/op",
            "extra": "399126 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29507,
            "unit": "ns/op",
            "extra": "407252 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25316,
            "unit": "ns/op",
            "extra": "473642 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1393948346,
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
          "id": "fe065d2207e2b0fc029ee5309373fc054937d7aa",
          "message": ":arrow_up: Bump sigs.k8s.io/controller-runtime from 0.18.5 to 0.19.1\n\nBumps [sigs.k8s.io/controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) from 0.18.5 to 0.19.1.\n- [Release notes](https://github.com/kubernetes-sigs/controller-runtime/releases)\n- [Changelog](https://github.com/kubernetes-sigs/controller-runtime/blob/main/RELEASE.md)\n- [Commits](https://github.com/kubernetes-sigs/controller-runtime/compare/v0.18.5...v0.19.1)\n\n---\nupdated-dependencies:\n- dependency-name: sigs.k8s.io/controller-runtime\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-11-10T23:52:48-05:00",
          "tree_id": "54bb4a34fffd48e363f2573eaa34d2aed85399ae",
          "url": "https://github.com/guidewire-oss/kumquat/commit/fe065d2207e2b0fc029ee5309373fc054937d7aa"
        },
        "date": 1731300946253,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29721,
            "unit": "ns/op",
            "extra": "408262 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29325,
            "unit": "ns/op",
            "extra": "395887 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25438,
            "unit": "ns/op",
            "extra": "472386 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1410669123,
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
          "id": "a8d7342258ae441e8d8e3bb09c3158d92d6c9b78",
          "message": ":arrow_up: Bump k8s.io/apimachinery from 0.30.3 to 0.31.2\n\nBumps [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) from 0.30.3 to 0.31.2.\n- [Commits](https://github.com/kubernetes/apimachinery/compare/v0.30.3...v0.31.2)\n\n---\nupdated-dependencies:\n- dependency-name: k8s.io/apimachinery\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-11-11T01:00:16-05:00",
          "tree_id": "13ba64d04a65f6f439796607fdb1564960cb7cf9",
          "url": "https://github.com/guidewire-oss/kumquat/commit/a8d7342258ae441e8d8e3bb09c3158d92d6c9b78"
        },
        "date": 1731305007582,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29660,
            "unit": "ns/op",
            "extra": "407241 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29470,
            "unit": "ns/op",
            "extra": "400644 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25493,
            "unit": "ns/op",
            "extra": "460437 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1408168882,
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
          "id": "6474fafa18564cbfba0b83b021eb85c2e7943190",
          "message": ":arrow_up: Bump k8s.io/client-go from 0.30.3 to 0.31.2\n\nBumps [k8s.io/client-go](https://github.com/kubernetes/client-go) from 0.30.3 to 0.31.2.\n- [Changelog](https://github.com/kubernetes/client-go/blob/master/CHANGELOG.md)\n- [Commits](https://github.com/kubernetes/client-go/compare/v0.30.3...v0.31.2)\n\n---\nupdated-dependencies:\n- dependency-name: k8s.io/client-go\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-11-11T01:07:33-05:00",
          "tree_id": "93fb0cef4f03985401f33c9acaede63bb0b29226",
          "url": "https://github.com/guidewire-oss/kumquat/commit/6474fafa18564cbfba0b83b021eb85c2e7943190"
        },
        "date": 1731305435153,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29937,
            "unit": "ns/op",
            "extra": "409462 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29320,
            "unit": "ns/op",
            "extra": "404432 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25406,
            "unit": "ns/op",
            "extra": "463424 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1474686419,
            "unit": "ns/op",
            "extra": "7 times\n4 procs"
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
          "id": "dcf8ea5a87ccc57e29a7dd137cfbf13afc95669e",
          "message": ":arrow_up: Bump cuelang.org/go from 0.10.0 to 0.10.1\n\nBumps cuelang.org/go from 0.10.0 to 0.10.1.\n\n---\nupdated-dependencies:\n- dependency-name: cuelang.org/go\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-11-11T09:53:26-05:00",
          "tree_id": "bb415f1f369e918f66248b55cfd630755d618c1a",
          "url": "https://github.com/guidewire-oss/kumquat/commit/dcf8ea5a87ccc57e29a7dd137cfbf13afc95669e"
        },
        "date": 1731336983595,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30109,
            "unit": "ns/op",
            "extra": "404769 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 30069,
            "unit": "ns/op",
            "extra": "392848 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25832,
            "unit": "ns/op",
            "extra": "458414 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1403374706,
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
          "id": "c6adad9f05029fd343c1afd9caa321b63990b13c",
          "message": ":arrow_up: Bump github.com/onsi/ginkgo/v2 from 2.20.2 to 2.21.0\n\nBumps [github.com/onsi/ginkgo/v2](https://github.com/onsi/ginkgo) from 2.20.2 to 2.21.0.\n- [Release notes](https://github.com/onsi/ginkgo/releases)\n- [Changelog](https://github.com/onsi/ginkgo/blob/master/CHANGELOG.md)\n- [Commits](https://github.com/onsi/ginkgo/compare/v2.20.2...v2.21.0)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/onsi/ginkgo/v2\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-11-11T09:53:36-05:00",
          "tree_id": "020d35ee819b585964f189edfbf9fcf865ec3e82",
          "url": "https://github.com/guidewire-oss/kumquat/commit/c6adad9f05029fd343c1afd9caa321b63990b13c"
        },
        "date": 1731336992943,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29697,
            "unit": "ns/op",
            "extra": "402199 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29922,
            "unit": "ns/op",
            "extra": "379426 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25558,
            "unit": "ns/op",
            "extra": "461241 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1428591849,
            "unit": "ns/op",
            "extra": "7 times\n4 procs"
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
          "id": "e16e7f7d86e2a6aec96815386907c58cc1d8f401",
          "message": ":arrow_up: Bump github.com/onsi/gomega from 1.34.2 to 1.35.1\n\nBumps [github.com/onsi/gomega](https://github.com/onsi/gomega) from 1.34.2 to 1.35.1.\n- [Release notes](https://github.com/onsi/gomega/releases)\n- [Changelog](https://github.com/onsi/gomega/blob/master/CHANGELOG.md)\n- [Commits](https://github.com/onsi/gomega/compare/v1.34.2...v1.35.1)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/onsi/gomega\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-11-11T10:03:45-05:00",
          "tree_id": "1395ffc63307211ab067266cfa94f67ff87dd564",
          "url": "https://github.com/guidewire-oss/kumquat/commit/e16e7f7d86e2a6aec96815386907c58cc1d8f401"
        },
        "date": 1731337607481,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30148,
            "unit": "ns/op",
            "extra": "399696 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29818,
            "unit": "ns/op",
            "extra": "400752 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25800,
            "unit": "ns/op",
            "extra": "457116 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1395665727,
            "unit": "ns/op",
            "extra": "8 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "d21it185@charusat.edu.in",
            "name": "D21IT185BapodraRajSatish",
            "username": "D21IT185BapodraRajSatish"
          },
          "committer": {
            "email": "dobson@softwarepunk.com",
            "name": "James Dobson",
            "username": "jamesdobson"
          },
          "distinct": true,
          "id": "200f66713bce4aa1d57215cdf777df74f21193e4",
          "message": "Refactor Delete method in SQLiteRepository for improved readability",
          "timestamp": "2024-11-11T10:29:22-05:00",
          "tree_id": "d2ad9e7b4164eecad7b104dfbf402f940e74f428",
          "url": "https://github.com/guidewire-oss/kumquat/commit/200f66713bce4aa1d57215cdf777df74f21193e4"
        },
        "date": 1731339077090,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30081,
            "unit": "ns/op",
            "extra": "396339 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29600,
            "unit": "ns/op",
            "extra": "395874 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25455,
            "unit": "ns/op",
            "extra": "469220 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1410525578,
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
          "id": "4013642a3bca2980e0107ead94b986062e90888b",
          "message": "✅  Determine architecture for TestController tests.",
          "timestamp": "2024-11-13T13:13:57-05:00",
          "tree_id": "50d3061e7c85e74869a2a9d5605a96287a4f493a",
          "url": "https://github.com/guidewire-oss/kumquat/commit/4013642a3bca2980e0107ead94b986062e90888b"
        },
        "date": 1731521736369,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 29715,
            "unit": "ns/op",
            "extra": "403309 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 29400,
            "unit": "ns/op",
            "extra": "402990 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25477,
            "unit": "ns/op",
            "extra": "464653 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1398228318,
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
          "id": "89ce80478776c66b36a703e7b085cd7b09997ee3",
          "message": "📝  Update example with realistic annotation using the '.' character.",
          "timestamp": "2024-11-15T19:26:52-05:00",
          "tree_id": "b3283faed79adbf1d5edb15a0ed6ce53620e1a17",
          "url": "https://github.com/guidewire-oss/kumquat/commit/89ce80478776c66b36a703e7b085cd7b09997ee3"
        },
        "date": 1731716937247,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkQueryPerformance/QueryFirst",
            "value": 30465,
            "unit": "ns/op",
            "extra": "392084 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryLast",
            "value": 30105,
            "unit": "ns/op",
            "extra": "394956 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryMissing",
            "value": 25594,
            "unit": "ns/op",
            "extra": "473449 times\n4 procs"
          },
          {
            "name": "BenchmarkQueryPerformance/QueryCartesianProduct",
            "value": 1426965514,
            "unit": "ns/op",
            "extra": "8 times\n4 procs"
          }
        ]
      }
    ]
  }
}