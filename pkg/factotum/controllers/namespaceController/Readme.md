# NamespaceController

The NamespaceController is used to support objconfigs.factotum.io. NamespaceConfigs contain labels/annotations and taints to be applied to Kubernetes Node Objects. Users creates NamespaceConfigs and the desired configuration is applied to the objs. 

The NamespaceConfigs are cached since The NamespaceController also watches for any changes to Nodes that might impact one of our configurations. If for example a label is removed that is present in a NamespaceConfig. It will be added back immediately.


# NamespaceConfig Flow

```mermaid
sequenceDiagram
    User->>NamespaceConfig: User creates NamespaceConfig with labels/taints/annotations
    NamespaceConfig->>Reconciler: Reconciler is notified about NamespaceConfig Event
    Reconciler-->>NamespaceConfig: Adds Finalizer if not present
    Reconciler->>NamespaceController: Notifies NamespaceController to process NamespaceConfig. Blocks via WG
    NamespaceController->>NodeObject: Adds configured labels/taints/annotations
    NamespaceController->>Reconciler: Sends Wg.Done() after completing appropriate operations
    Reconciler->>NamespaceConfig: Status is updated with current state
```

# Node Event Flow

## Events from NamespaceConfig Changes

Events that handle NamespaceConfig CRD creations follow this pattern.

```mermaid
sequenceDiagram
  K8sApi ->> Reconcilier: Notify Of NamespaceConfigEvent
  Reconcilier ->> NamespaceController: Pass Node Config
  NamespaceController ->> v1.Node: Filter Known Nodes and apply
  NamespaceController ->> NamespaceController: Update Cache
  NamespaceController ->> Reconcilier: Notify Complete
  Reconcilier ->> K8sApi: Update CRD status fields
  Reconcilier ->> K8sApi: Send Result
```

## Events from Node Changes

Node Changes are enforcement events. If something has removed a configuration this is how we put it back.

```mermaid
sequenceDiagram
  K8sApi ->> NamespaceController: Send Node Events
  NamespaceController ->> NamespaceController: Read From Cache of Existing NamespaceConfigs
  NamespaceController ->> v1.Node: Filter Configs and apply
  NamespaceController ->> NamespaceController: Update Cache  
```