# Node Config

NodeConfig can be used to set labels, annotations, and taints on Node objects. If a selector is defined it will be used to limit application to matching nodes.

```
apiVersion: factotum.io/v1alpha1
kind: NodeConfig
metadata:
  labels:
    app.kubernetes.io/name: factotum
  name: nodeconfig-sample
spec:
  selector:
    kubernetes.io/hostname: hostname[1-9] # regexs are supported for label values
  labels:
    factotum: applied
  annotations:
    factotum: applied
  taints:
  - key: factotum
    value: tainted
    effect: NoSchedule
```