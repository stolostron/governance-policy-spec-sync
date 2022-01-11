// NOTE: Boilerplate only.  Ignore this file.
// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

// Package v1 contains API Schema definitions for the policies v1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=policy.open-cluster-management.io
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "policy.open-cluster-management.io", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	// Kind Policy
	Kind = "Policy"
)
