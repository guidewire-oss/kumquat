window.BENCHMARK_DATA = {
  "lastUpdate": 1731260097528,
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
      }
    ]
  }
}