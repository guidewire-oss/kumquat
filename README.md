![kumquat](documents/kumquat-256.png)

*Hi everybody, I'm kumquat! You can configure me to watch Kubernetes resources and generate new resources from them!*

![coverage-badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/jamesdobson/cb14f8ad568d175cf0ba9f7ac6d0a6ca/raw/kumquat-coverage.json)
[![performance-badge](https://img.shields.io/badge/benchmarks-olive)](https://guidewire-oss.github.io/kumquat/dev/bench/)

# Kumquat
Kumquat is **KU**bernetes **M**etacontroller **QU**ery **A**nd **T**emplate.

Kumquat:
1. Queries multiple Kubernetes resources.
2. Sends the query results to a template engine as the input data.
3. Uses a template to generate more Kubernetes resources.

The query and the template are configured in one or more `Template` resources that are loaded at run-time.

Kumquat supports the following template languages:
* cue
* gotemplate
* jsonnet (when built with `-tags jsonnet`)


## Example
Kubernetes RBAC is powerful, but sometimes it is a bit limited. Imagine you want a `ClusterRole`
that provides read access to all the Crossplane resources in your cluster that are provided by
the Crossplane AWS providers. Each provider creates resources under a multitude of API groups:
`s3.aws.upbound.io`, `dynamodb.aws.upbound.io`, etc. Since Kubernetes doesn't support wildcards
in the `apiGroups` field of the `ClusterRole`, you'll have to list them all manually, and keep
this list up-to-date as new AWS providers are installed.

Kumquat can help! In the following `Template`, kumquat queries for all the CRDs in the cluster, and
finds those that end in `.aws.upbound.io`. It passes those to a CUE program that finds the set of
unique API group names, and outputs a `ClusterRole` with those API groups.

```yaml
apiVersion: kumquat.guidewire.com/v1beta1
kind: Template
metadata:
  name: generate-role
  namespace: templates
spec:
  query: | #sql
    SELECT crd.data AS crd
    FROM "CustomResourceDefinition.apiextensions.k8s.io" AS crd
    WHERE crd.name LIKE "%.aws.upbound.io"
  template:
    language: cue
    batchModeProcessing: true
    data: | #cue
      unique_groups_map: {
        for result in data {
          "\(result.crd.spec.group)": result.crd.spec.group
        }
      }
      unique_groups: [ for g in unique_groups_map {g}]
      out: {
        apiVersion: "rbac.authorization.k8s.io/v1"
        kind: "ClusterRole"
        metadata: 
          name: "crossplane-aws-readers-role"
        rules: [
            {
              apiGroups: unique_groups
              resources: "*"
              verbs: ["get", "list", "watch"]
            }
        ]
      }
```

The expected output depends on which Crossplane AWS providers are installed, but should be roughly like this:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: crossplane-aws-readers-role
rules:
    - apiGroups:
        - dynamodb.aws.upbound.io
        - opensearch.aws.upbound.io
        - appautoscaling.aws.upbound.io
      resources: '*'
      verbs:
        - get
        - list
        - watch
```

Because kumquat is a controller, when new Crossplane AWS providers are installed it will automatically
detect the new CRDs and update the `ClusterRole` accordingly.


## Other Example Ideas

* PV needing a reference to an external volume and looking it up in a ConfigMap
* The example we use could be a toy example, e.g. hello world, and we could put better examples in an examples directory, with a README to motivate them


## Make All Examples
You can make all the examples in the `examples` folder by running:

```
go run mage.go examples
```

This will also compare each example output against the expected output, and acts as a kind of acceptance
test in a way.


## Add Experimental Jsonnet Support
```
go build -buildvcs=false -tags jsonnet .
```

## Development

Historic results for the benchmarks are available on the
[benchmarks page](https://guidewire-oss.github.io/kumquat/dev/bench/).


## References

1. [How to write a custom Kubernetes Controller](https://arunprasad86.medium.com/how-to-write-a-custom-kubernetes-controller-4904383cec4)


## Similar Projects / Inspiration
* Using SQL with Kubernetes:
  * [Dentrax / kubesql](https://github.com/Dentrax/kubesql)
  * [xuxinkun / kubesql](https://github.com/xuxinkun/kubesql)
* [Metacontroller](https://metacontroller.github.io/metacontroller/intro.html)
* [jsPolicy](https://www.jspolicy.com/); see the *Controller* policy type.

