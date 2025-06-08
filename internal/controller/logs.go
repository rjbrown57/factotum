package controller

import ctrl "sigs.k8s.io/controller-runtime"

var controllerLog = ctrl.Log.WithName("factotum")
var debugLog = controllerLog.V(1)
var traceLog = controllerLog.V(2)
