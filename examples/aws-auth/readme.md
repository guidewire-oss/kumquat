# Managing the EKS `aws-auth` `ConfigMap` with Multiple Owners

## Background
In AWS EKS, a `ConfigMap` called `aws-auth` is used to map IAM principals to
Kubernetes principals. The resource must be in the `kube-system` namespace, and
there is only one such resource.


## Problem
In some cases, principal mappings need to be added to `aws-auth` by various
systems. It's difficult to do this if those systems operate in a declarative
fashion, because each system may try to fully manage the resource. For example,
if `aws-auth` needs to be managed by multiple Terraform modules or environemnts
using the `kubernetes` provider, then the ownership is unclear and the resource
will never reconcile properly.

A use case might be if Terraform code is used to set up a cluster and configure
`aws-auth` with a basic admin role mapping, and separate Terraform code adds
role mappings for each tenant of the Kubernetes cluster. For example:

```yaml
apiVersion: v1
data:
    mapRoles: |
        - rolearn: arn:aws:iam::111122223333:role/admin-role
          username: admin-role
          groups:
            - system:masters
        - rolearn: arn:aws:iam::111122223333:role/tenant-acme-role
          username: acme
          groups:
            - tenant:acme
        - rolearn: arn:aws:iam::111122223333:role/tenant-umbrella-role
          username: umbrella
          groups:
            - tenant:umbrella
kind: ConfigMap
metadata:
    name: aws-auth
    namespace: kube-system
```

In the above, the `admin-role` is added by the cluster provisioner Terraform code,
and the `tenant-acme-role` and `tenant-umbrella-role` are added by Terraform code
for provisioning each tenant.

Both code bases cannot manage `aws-auth`!


## Solution
We use a kumquat template that looks for `ConfigMap` instances and merges them
to produce the final `aws-auth` `ConfigMap`. The Terraform code creates the
input `ConfigMap` instances, so each piece of Terraform code can own its own
resource. Kumquat is the only thing that owns `aws-auth`.

The following template does just that. It also check for the presence of the
`aggregate` annotation on the inputs, and ensures the value is `aws-auth`.
This ensure that the template only merges `ConfigMap` instances that belong
in `aws-auth`. The choice of aggregation name and value was arbitrary.

```yaml
apiVersion: kumquat.guidewire.com/v1beta1
kind: Template
metadata:
  name: aws-auth
  namespace: templates
spec:
  query: | #sql
    SELECT cm.data AS cm
    FROM "ConfigMap.core" AS cm
    WHERE cm.namespace = 'kube-system' AND
    json_extract(cm.data, '$.metadata.annotations.aggregate') = 'aws-auth'
  template:
    language: cue
    batchModeProcessing: true
    fileName: |
      ./output/out.yaml
    data: | #cue
      import "strings"

      _mapRoles: strings.Join([for result in data {result.cm.data.mapRoles}], "")
      {
        apiVersion: "v1"
        kind: "ConfigMap"
        metadata: {
          name: "aws-auth"
          namespace: "kube-system"
        }
        data: {
          "mapRoles": _mapRoles
        }
      }
```
