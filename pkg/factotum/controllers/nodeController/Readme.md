# NodeController

The NodeController is used to support nodeconfigs.factotum.io. NodeConfigs contain labels/annotations and taints to be applied to Kubernetes Node Objects. Users creates NodeConfigs and the desired configuration is applied to the nodes. 

The NodeConfigs are cached since The NodeController also watches for any changes to Nodes that might impact one of our configurations. If for example a label is removed that is present in a NodeConfig. It will be added back immediately.


# NodeConfig Flow

```mermaid
sequenceDiagram
    User->>NodeConfig: User creates NodeConfig with labels/taints/annotations
    NodeConfig->>Reconciler: Reconciler is notified about NodeConfig Event
    Reconciler-->>NodeConfig: Adds Finalizer if not present
    Reconciler->>NodeController: Notifies NodeController to process NodeConfig. Blocks via WG
    NodeController->>NodeObject: Adds configured labels/taints/annotations
    NodeController->>Reconciler: Sends Wg.Done() after completing appropriate operations
    Reconciler->>NodeConfig: Status is updated with current state
```

# Node Event Flow

## Events from NodeConfig Changes

Events that handle NodeConfig CRD creations follow this pattern.

```mermaid
sequenceDiagram
  K8sApi ->> Reconcilier: Notify Of NodeConfigEvent
  Reconcilier ->> NodeController: Pass Node Config
  NodeController ->> v1.Node: Filter Known Nodes and apply
  NodeController ->> NodeController: Update Cache
  NodeController ->> Reconcilier: Notify Complete
  Reconcilier ->> K8sApi: Update CRD status fields
  Reconcilier ->> K8sApi: Send Result
```

## Events from Node Changes

Node Changes are enforcement events. If something has removed a configuration this is how we put it back.

```mermaid
sequenceDiagram
  K8sApi ->> NodeController: Send Node Events
  NodeController ->> NodeController: Read From Cache of Existing NodeConfigs
  NodeController ->> v1.Node: Filter Configs and apply
  NodeController ->> NodeController: Update Cache  
```