# Use SQL for the Query Language

- Status: accepted
- Deciders: James Dobson, Amir Bavand

## Context and Problem Statement

How can we represent queries over multiple Kubernetes resources in a compact fashion that allows us to express all the
kinds of relational queries we wish to perform? How can we ensure that this representation of the queries is compact
and easy to use, matching what the user might expect when querying Kubernetes resources?



## Decision Drivers

- Expressiveness: must allow users to write complex queries involving multiple Kubernetes resources.
- Ease of use: must not require users to learn a lot of new syntax.
- Must nicely fit the task at hand: querying Kubernetes resources.
- Performance: it should be easy to write queries that complete in a reasonable amount of time and that do not consume
  an undue amount of memory.
- Ease of implementation: we should be able to PoC a working system relatively quickly, without getting bogged down in
  the details like having to write query optimizers.



## Considered Options

- Option 1: Relational Algebra as AST in YAML
- Option 2: SQL
- Option 3: NoSQL Database


### Option 1: Relational Algebra as AST in YAML
One option is to use relational algebra to express the queries.

For example, the PersistentVolume referenced by each PersistentVolumeClaim in a namespace can be found with the
following relational algebra expression:

<p align="center">
𝜎<sub>*pvc.spec.volumeName=pv.metadata.name*</sub>(
𝜌<sub>*pvc*</sub>(𝜎<sub>*metadata.namespace="default"*</sub>(*PersistentVolumeClaim*)) ⨯
𝜌<sub>*pv*</sub>(𝜎(*PersistentVolume*)))
</p>


Such an expression could be represented by its AST in code:

```go
query := query.EqualQuery{
	Query: query.CrossProductQuery{
		Left: query.SelectionQuery{
			Kind:      "PersistentVolumeClaim",
			Namespace: "default",
			RenameAs:  "pvc",
		},
		Right: query.SelectionQuery{
			Kind:      "PersistentVolume",
			RenameAs:  "pv",
		},
	},
	Path1: []string{"pvc", "spec", "volumeName"},
	Path2: []string{"pv", "metadata", "name"},
}
```

The AST could be represented in the YAML template as:

```yaml
...
spec:
  query:
    type: eq
    path1:
      - pvc
      - spec
      - volumeName
    path2:
      - pv
      - metadata
      - name
    query:
      type: cross
      left:
        type: select
        kind: PersistentVolumeClaim
        namespace: default
        renameAs: pvc
      right:
        type: select
        kind: PersistentVolume
        renameAs: pv
...
```


### Option 2: SQL
Another option is to use SQL to represent the queries.

For example, the PersistentVolume referenced by each PersistentVolumeClaim can be found with the
following SQL expression:

```sql
SELECT persistentvolumeclaim.data AS pvc, persistentvolume.data AS pv
FROM "PersistentVolumeClaim.core" AS persistentvolumeclaim
JOIN "PersistentVolume.core" AS persistentvolume ON persistentvolumeclaim.data ->> '$.spec.volumeName' = persistentvolume.name
``` 

This assumes that each Kubernetes resource type is treated as a different table:


| namespace | name | data |
| --------- | ---- | :--- |
| default   | pvc1 | { "apiVersion": "core/v1", "kind": "PersistentVolumeClaim", ... |
| default   | pvc2 | { "apiVersion": "core/v1", "kind": "PersistentVolumeClaim", ... | 
**PersistentVolumeClaim.core**

| namespace | name | data |
| --------- | ---- | :--- |
| NULL      | pv1  | { "apiVersion": "core/v1", "kind": "PersistentVolume", ... |
| NULL      | pv2  | { "apiVersion": "core/v1", "kind": "PersistentVolume", ... | 
**PersistentVolume.core**


### Option 3: NoSQL Database
An option considered was to store the resources in a NoSQL document-based database. We could not find a one with
a query language that allowed sufficiently expressive queries across multiple object collections.



## Decision Outcome

Chosen option: "Option 2: SQL", because it is a compact, well-understood language that allows users of our system
to query across multiple Kubernetes resources.
 

### Positive Consequences 

- Less work to implement: we can use an off-the-shelf system that provides the SQL query capabilities we need.
- Ease of use: users who are already familiar with SQL will find that their knowledge ports easily. Users who are
  not already familiar with SQL can find a lot of reference and training material online. 
- Performance: our system will benefit from any query optimizations applied by the underlying DBMS.

### Negative Consequences

- Impedance mismatch: The SQL syntax does not map perfectly to our domain.
  - The attribute list on the select statement will need to be a specific format to work with our system as intended.
  - The Kubernetes resources are stored in a TEXT column as JSON. Most relations will have to operate on paths
    within the JSON, which requires the use of non-standard SQL extensions that are a bit clunky to use.
- Susceptible to DoS: users can intentionally construct queries that perform poorly and/or consume a lot of memory.



## Pros and Cons of the Options

### Option 1: Relational Algebra as AST in YAML

- Good, because we can develop a syntax that prevents mistakes because it is suited to the task of querying
  Kubernetes resources. 
- Good, because we can implement the syntax tree using reactive programming techniques so that it recalculates as
  Kubernetes resource create/update/delete events occur.
- Bad, because representing an AST in YAML is verbose and unweildy.

### Option 2: SQL

- Good, because SQL is compact to represent as a string in YAML.
- Good, because SQL is relatively well-known and there are many online resources to learn.
- Good, because we can make use of an off-the-shelf RDBMS to execute the queries.
- Bad, because SQL is not well-suited to working with JSON data.
- Bad, because SQL requires a certain amount of boiler-plate in the attribute list to work correctly with the rest
  of our system.

### Option 3: NoSQL Database
We could not find a viable system that allowed sufficiently-expressive queries.

We recognize that a viable option may be discovered in the future. Therefore, we will ensure that our system is
designed to work with multiple query engines.


