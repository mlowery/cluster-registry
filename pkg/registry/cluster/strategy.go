/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	apistorage "k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/cluster-registry/pkg/apis/clusterregistry/install"
	clusterregistry "k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
	"k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1/validation"
)

type clusterStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var Strategy = clusterStrategy{install.Scheme, names.SimpleNameGenerator}

func (clusterStrategy) NamespaceScoped() bool {
	return false
}

func ClusterToSelectableFields(cluster *clusterregistry.Cluster) fields.Set {
	return generic.ObjectMetaFieldsSet(&cluster.ObjectMeta, false)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, bool, error) {
	cluster, ok := obj.(*clusterregistry.Cluster)
	if !ok {
		return nil, nil, false, fmt.Errorf("given object is not a cluster.")
	}
	return labels.Set(cluster.ObjectMeta.Labels), ClusterToSelectableFields(cluster), cluster.Initializers != nil, nil
}

func MatchCluster(label labels.Selector, field fields.Selector) apistorage.SelectionPredicate {
	return apistorage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
func (clusterStrategy) PrepareForCreate(ctx genericapirequest.Context, obj runtime.Object) {
}

// Validate validates a new cluster.
func (clusterStrategy) Validate(ctx genericapirequest.Context, obj runtime.Object) field.ErrorList {
	cluster := obj.(*clusterregistry.Cluster)
	return validation.ValidateCluster(cluster)
}

// Canonicalize normalizes the object after validation.
func (clusterStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for cluster.
func (clusterStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (clusterStrategy) PrepareForUpdate(ctx genericapirequest.Context, obj, old runtime.Object) {
}

// ValidateUpdate is the default update validation for an end user.
func (clusterStrategy) ValidateUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateClusterUpdate(obj.(*clusterregistry.Cluster), old.(*clusterregistry.Cluster))
}
func (clusterStrategy) AllowUnconditionalUpdate() bool {
	return true
}
