// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project
// This contains ports of github.com/operator-framework/operator-sdk v0.19.4.
// This is required since operator-sdk is no longer a dependency after upgrading
// to operator-sdk v1.x.x for this project.

package tool

import (
	"fmt"
	"os"
)

type RunModeType string

const (
	ForceRunModeEnv             = "OSDK_FORCE_RUN_MODE"
	LocalRunMode    RunModeType = "local"
	ClusterRunMode  RunModeType = "cluster"
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	watchNamespaceEnvVar = "WATCH_NAMESPACE"
)

// ErrNoNamespace indicates that a namespace could not be found for the current
// environment
var ErrNoNamespace = fmt.Errorf("namespace not found for current environment")

// ErrRunLocal indicates that the operator is set to run in local mode (this error
// is returned by functions that only work on operators running in cluster mode)
var ErrRunLocal = fmt.Errorf("operator run mode forced to local")

// GetWatchNamespace returns the Namespace the operator should be watching for changes
func GetWatchNamespace() (string, error) {
	ns, found := os.LookupEnv(watchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", watchNamespaceEnvVar)
	}

	return ns, nil
}
