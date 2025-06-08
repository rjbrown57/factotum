package namespacecontroller

import ctrl "sigs.k8s.io/controller-runtime"

var log = ctrl.Log.WithName(controllerName)
var debugLog = log.V(1)
var traceLog = log.V(2)
