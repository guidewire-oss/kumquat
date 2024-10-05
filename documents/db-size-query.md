---
marp: true
theme: default
size: "16:9"
@auto-scaling code
---
# Kumquat Performance: Querying Large Kubernetes Clusters
September 2024

![bg right auto](kumquat-512.png)

---
## What is Kumquat?

- Kumquat watches for resources in the Kubernetes API
- Processes the watched resources with a template to make new resources
- The behaviour is specified declaratively and at run-time
- Any resource can be watched; any resource can be created


---
## Sample Kumquat Template

- Looks for `ConfigMap` resources in the `kube-system` namespace
- Selects only those where the `aggregate` annotation is set to `aws-auth`

```sql
SELECT cm.data AS cm
FROM "ConfigMap.core" AS cm
WHERE cm.namespace = 'kube-system'
AND json_extract(cm.data, '$.metadata.annotations.aggregate') = 'aws-auth'
```


---
## Sample Kumquat Template

- Aggregates the queried `ConfigMap` resources into a single, well-known `ConfigMap`:

```cue
import "strings"

mapRoles: strings.Join([for result in data {result.cm.data.mapRoles}], "")
out: {
  apiVersion: "v1"
  kind: "ConfigMap"
  metadata: {
    name: "aws-auth"
    namespace: "kube-system"
  }
  data: {
    "mapRoles": mapRoles
  }
}
```


---
## Observations About the Sample Template

- The query can specify any resource(s)
- The output template can create any resource(s)
- Kumquat is dynamic: the template itself is packaged as custom resource that Kumquat loads at run-time.


---
## Question: How will Kumquat Perform?

- Largest concern is the query performance
- Initial implementation uses SQLite in-memory database
- Will everything fit in memory?
- How quickly can we query it?


---
## Why Might it be Slow?

- Many of the queries we want will query mutiple resource types, so the domain
  through which SQLite must search could be quite large
- Queries that use `json_extract()` likely don't benefit from indexing, so this common case could be quite taxing on the query execution engine
- The amount of data loaded to memory is unknown: how much resource data do our clusters have that would need to be stored in Kumquat?


---
## Reasons for Optimism

- Processing data sets that are held in memory is *fast*
- It is likely Kumquat will be able to hold all resources it needs in memory: other
  controllers do so for all the resource kinds they watch. Consider:
  - Deployment controller
  - Crossplane AWS controller
- DB transactions aren't needed
- There are options to improve performance and/or memory consumption if needed:
  - Kumquat can parse queries itself to reduce the number of resources loaded
  - Kumquat can use SQLite's *Indexes on Expressions* feature to overcome `json_extract()` performance
  - Can use SQLite's `JSONB` format for JSON processing
  - Can switch to an external database if memory consumption is too high


---
## Test Configuration

- Use a Jupyter notebook with a Go kernel
- Download all resources in Kubernetes cluster to disk as YAML files
- Instantiate them in a `Repository` (Kumquat's internal in-memory storage)
- Try some queries to measure the performance
- Notbook is available at https://github.com/guidewire-oss/kumquat/blob/chore/repository-performance/documents/db-size-query.ipynb


---
## How Big Are Our Clusters?

- Go script in the Jupyter notebook that:
  - list all resources in the cluster via the Kubernetes API
  - save to disk as YAML
  - uses the default throttling behaviour of the Kubernetes client.

|  Small Cluster   |  Large Cluster    |
|------------------|-------------------|
| 36,766 resources | 110,948 resources |
| 337 MB of YAML   | 1.1 GB of YAML    |
| Download: 10m38s | Download: 32m47s  |


---
# Will Everything Fit in Memory?


---
## Load Everything and See What Happens...

- The worst case is Kumquat gets configured to watch all resources
- Measure the time to load
- Measure the memory consupmption


---
## First Attempt: Load the Small Cluster

- Loaded the directory with all the YAMLs for the small cluster
- Took 5 minutes, 49 seconds
- Go's `runtime.ReadMemStats()` showed 18 MiB of usage
- The `ps` command showed that RSS grew by 266.0 MiB
- Analysis:
  - `runtime.ReadMemStats()` gives Go heap usage, not total memory used
  - `go-sqlite` must be using C heap, not Go heap
  - We should be able to load 337 MB of YAML faster than 5 minutes, 49 seconds


---
## Second Attempt: Load a ZIP of the Small Cluster

- Hypothesis: loading many files takes a long time due to `fopen()` overhead
- ZIP the YAMLs into a single file instead
- Load time improved to 33.6 seconds
- RSS grew by 303.3 MiB
- Analysis:
  - In the previous attempt, performance was limited by `fopen()` not SQLite
  - Memory usage is roughly equal to disk size of 337 MB
  - Unsure why memory usage is not higher


---
## Load ZIP of the Large Cluster

- Loaded 1.1 GB of YAML from 110,948 resources (from a single ZIP file)
- Took 1 minute, 47 seconds to load
- Analysis:
  - Significantly less than the 32 minutes, 47 seconds to get all the resources from the Kubernetes API
  - Size of cluster (GB of YAML, count of resources) was roughtly 3x the small cluster
  - Took roughly 3x longer to load than the small cluster -- linear scaling?


---
## Will Everything Fit in Memory?

- Highly likely, even if we load the entire cluster to memory
- The in-memory size was not significantly different from the on-disk YAML size
- We should allocate enough memory to the controller pod to hold the resources of a large cluster
- For loading/monitoring/updating resources, SQLite will not be the bottleneck
- Prediction: the Kubernetes API will likely by the bottleneck 
- Future investigations:
  - Is there a larger cluster we can try?
  - Why was the in-memory size seemingly smaller than the on-disk size?


---
# How Quickly Can we Query?


---
## Query All the ConfigMap Resources

- This query gets all the `ConfigMap` resources:
```go
rs, err := repo.Query(`SELECT cm.data AS cm FROM "ConfigMap.core" AS cm`)
```
- It runs very quickly on the small and large clusters:

|           Small Cluster            |           Large Cluster            |
|------------------------------------|------------------------------------|
| 4438 `ConfigMap` resources         | 7784 `ConfigMap` resources         |
| Query duration: 317 ms             | Query duration: 894 ms             |
| Go heap consumed: 40 MiB           | Go heap consumed: 103 MiB          |
| `ConfigMap` YAML files size: 41 MB | `ConfigMap` YAML files size: 89 MB |


---
## What Happens During Query?

- SQLite executes the query and returns a result set
- Kumquat's `Repository` unmarshals the returned JSON strings
- The `Query` method returns `Resource` instances with the parsed JSON
- The `Resource` objects are on the Go heap
- All of that happened for 7784 resources in ~900 ms!


---
## Pathological Query Example

- Let's try to make a query that performs badly:
  - Generate a large domain by joining a large table with itself
  - Find ways to make the DB avoid using indexes
- This will make the DB engine take a long time to execute the query
- Additionaly, we can negatively impact the performance by:
  - Returning a large number of results
- This will eat up memory and spend time in Kumquat processing the results


---
## Pathological Query 1

- Find pairs of `ConfigMap` with the same name
```sql
SELECT cm1.data AS cm1, cm2.data AS cm2
FROM "ConfigMap.core" AS cm1
CROSS JOIN "ConfigMap.core" AS cm2
WHERE cm1.name = cm2.name AND cm1.namespace != cm2.namespace
```
- There are 4438 `ConfigMap` resources in the small cluster, and 7784 in the large cluster
- The `CROSS JOIN` creates a search domain of $4438^{2} = 19,695,844$ pairs for the small cluster
  - $7784^{2} = 60,590,656$ pairs for the large cluster
- The `WHERE` clause restricted the returned results, but by how much?


---
## Pathological Query 1: Results

- Ran it on the small cluster:
  - Executed in 22.8 seconds
  - Retrieved 356,524 pairs, e.g:
    - `acolhuan/istio-ca-root-cert`
    - `addons/istio-ca-root-cert`
  - Used 4.3 GiB of Go heap
- This is good performance; it's more than just SQLite being measured:
  - Kumquat's time to unmarshal every resource in those pairs


---
## Pathological Query 1: Results on a Larger Cluster

- Ran it on the large cluster:
  - The process terminated with `signal: killed`
  - It ran out of memory
  - Increased Rancher Desktop memory from 16 GiB to 24 GiB --> still OOM
- The problem is all the pairs of `ConfigMap` being returned as results


---
## Pathlogical Query 1: Half the Memory

- Modified the query to halve memory consumption
- Just return the first resource in each pair
```sql
SELECT cm1.data AS cm1 /* , cm2.data AS cm2 */
FROM "ConfigMap.core" AS cm1
CROSS JOIN "ConfigMap.core" AS cm2
WHERE cm1.name = cm2.name AND cm1.namespace != cm2.namespace
```
- Retrieved 1,903,728 pairs
- Took 7 minutes, 42 seconds
- Used 14 GiB


---
## What's Using Memory?

- The `Repository` interface returns the entire set of results:
![](https://kroki.io/plantuml/svg/eNq9kz9rwzAQxXd_ikNjwVDL6dKhpGQOhNCtdDjLF6Na_5DkqfS792xiiAlNBkMQCN3T0_tJSNqmjDEP1hSp1y5gRAsRXZ8ogHy-ENugoZIvF0pA1WNHHzobeje6c5ZcBkOnXBS--SaVQRwp-SEqEoAJdvsK4acA4Ij2FcTOu5Pu9hgEaw4tsaZT1r5UWEbvc6ko5nkyMW90YAiGRPH7P6R5BEQ9ACLXMxp03G5B6vUQZWi8e-1unmazHtQOUaObIBYDCO4-U47adV8L2KGaWMpW8FSWb9PLm2o51_J-hrzKaJYZ9f2M-ipDLTM2cwavG0xOk__sPlQ8Pjt5N8BVfa5G9JZcyx_3D5JMJbw=)
- The first element of the pair is repeated many times
- Each time it appears it is a different instance
- Not using:
  - Object pooling (`Resource` is treated as immutable, so it's possible)
  - Cursor / Iterator / Streaming (some template languages expect the full collection)


---
## Pathological Query 1: Avoid Duplicates

- Previously, it took 7 minutes, 42 seconds to retrieve 1.9 million pairs
  - Some of that time was query execution
  - Some of that time was unmarshaling the results
- How to evaluate time for just query execution?
- Consider using `SELECT DISTINCT`:
```sql
SELECT DISTINCT cm1.data AS cm1
FROM "ConfigMap.core" AS cm1
CROSS JOIN "ConfigMap.core" AS cm2
WHERE cm1.name = cm2.name AND cm1.namespace != cm2.namespace
```
- Retrieved 3587 `ConfigMap` rows in 7.8 seconds!


---
## Pathological Query 2

- Previously, we returned a lot of results to evaluate query performance
- Found that most of the time is spent outside SQLite
- The previous query likely benefited from indexes on `name` and `namespace`
- What about a query using `json_extract()` to circumvent indexing:
```sql
SELECT cm1.data AS cm1, cm2.data AS cm2
FROM "ConfigMap.core" AS cm1
CROSS JOIN "ConfigMap.core" AS cm2
WHERE json_extract(cm1.data,'$.metadata.name') = json_extract(cm2.data,'$.metadata.name')
AND json_extract(cm1.data,'$.metadata.namespace') != json_extract(cm2.data,'$.metadata.namespace')
```


---
## Pathological Query 2: Results

- Ran it on the small cluster:
  - Executed in ~199 seconds
  - Retrieved 356,524 pairs
- Slow compared to 22.8 seconds when using indexes
- Observation: it is still really fast!
  - A benchmark shows Go's `json.Unmarshal()` takes ~96.187 µs to unmarshal a 4155 byte `ConfigMap`
  - SQLite must do 4 unmarshals per pair in the search domain, so: $\frac{96.187}{10^6} * 4 * 4438^2 = 7577.9$ seconds
  - This is 38x the ~199 s time we measured!
  - Unknown how it's so fast


---
## Pathological Query 2: Results on a Larger Cluster

- Try with the large cluster
- Add `DISTINCT` and retrieve just the first value of each pair to avoid OOM
- This shouldn't cause a significantly different query plan in SQLite
- We can therefore predict the execution time:
  - For the small cluster, there were 4438 `ConfigMap` resources and the query took ~199 s.
  - The large cluster has 7784 `ConfigMap` resources.
  - The number of comparisons should scale quadratically with the number of `ConfigMap` resources:

$$
\frac{199 \,s}{4438^2} * 7784^2 \approx 612\,s
$$

- Prediction of 612 s is equivalent to 10 minutes and 12 seconds


---
## Pathological Query 2: Results on a Larger Cluster
- Executed in 10 minutes and 56 seconds, 7% longer than the prediction
- Can we over-specify the query to give the query optimizer enough information to use indexes?
```sql
SELECT DISTINCT cm1.data AS cm1
FROM "ConfigMap.core" AS cm1
CROSS JOIN "ConfigMap.core" AS cm2
WHERE json_extract(cm1.data,'$.metadata.name') = json_extract(cm2.data,'$.metadata.name')
AND json_extract(cm1.data,'$.metadata.namespace') != json_extract(cm2.data,'$.metadata.namespace')
AND cm1.name = cm2.name AND cm1.namespace != cm2.namespace
```
- Ran in 21 seconds
- Slower than using just the columns directly (7.8 seconds!)
- Faster than using only `json_extract()`


---
# Findings


---
## Conclusions

- **Will everything fit in memory?** Highly likely, even if Kumquat ends up watching the entire cluster.
- **How quickly can we query it?** Fast enough assuming rational queries and resource counts
  - Heavy use of `json_extract()` and similar functions will slow things down a lot
  - Users can optimize this by specifying constraints using indexed columns as much as possible
- Queries that return a lot of results can be slow or cause the process to be OOM Killed
  - We should protect the Kumquat controller from this


---
## Further Research

1. Try with the largest cluster we have; determine a good size for the Kumquat controller.
2. Why is the measured memory consumption less than expected?
3. How is SQLite able to evalute the query using `json_extract()` so quickly compared to our predictions?
4. Can we speed up `json_extract()` with *Indexes on Expressions* and/or JSONB?
5. Can Kumquat parse the queries to constrain the set of resources it watches?
6. How best to limit results to avoid OOM killed due to naive/malicious configurations?


---
# Questions?
