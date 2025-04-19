![intro](logo/factotum.png)

factotum is a k8s operator for general cluster operations

## Description

factotum provides some useful crds for kubernetes cluster admins and users.

## NodeConfig

Creating a NodeConfig crd will cause factotum to apply the supplied options to the selected nodes. If no selectors are added it will be applied to all nodes

```yaml
apiVersion: factotum.io/v1alpha1
kind: NodeConfig
metadata:
  name: nodeconfig-sample
spec:
  selector:
    nodeSelector:
      kubernetes.io/hostname: node[1-3] 
  annotations:
    factotum: applied
  labels:
    factotum: applied
  taints:
  - key: factotum
    value:  tainted
    effect: NoSchedule
```
