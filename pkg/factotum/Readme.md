# Factotum

This package contains the interfaces that help drive re-use across the different controllers factotum implements.

# Factotum Controller

Each Factotum Controller implements the same general pattern. Each Controller is tied to a Config and provides a Watch method, and a Process method. Work flows from the operator-sdk provided reconciler to a FactotumController. The FactotumController also enforces desired statue by watching for object changes.

```mermaid
sequenceDiagram
    User->>Config: User creates FactotumConfig with desired state
    Config->>Reconciler: Reconciler is notified about Config Event
    Reconciler-->>Config: Adds Finalizer if not present
    Reconciler->>FactotumController: Notifies FactotumController to process Config. Blocks via WG
    FactotumController->>Object: Adds desired configuration
    FactotumController->>Reconciler: Sends Wg.Done() after completing appropriate operations
    Reconciler->>Config: Status is updated with current state
```

# Object Event Flow

## Events from Config Changes

Events that handle Config CRD creations follow this pattern.

```mermaid
sequenceDiagram
  K8sApi ->> Reconcilier: Notify Of ConfigEvent
  Reconcilier ->> FactotumController: Pass Config
  FactotumController ->> Object: Filter Known Objects and apply
  FactotumController ->> FactotumController: Update Cache
  FactotumController ->> Reconcilier: Notify Complete
  Reconcilier ->> K8sApi: Update CRD status fields
  Reconcilier ->> K8sApi: Send Result
```

## Events from Object Changes

Node Changes are enforcement events. If something has removed a configuration this is how we put it back.

```mermaid
sequenceDiagram
  K8sApi ->> FactotumController: Send Object Events
  FactotumController ->> FactotumController: Read From Cache of Existing Configs
  FactotumController ->> Object: Filter Configs and apply
  FactotumController ->> FactotumController: Update Cache  
```

# Factotum Handlers