/*
Copyright 2015 The Kubernetes Authors.

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

package etcd

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/registry/event"
	"k8s.io/kubernetes/pkg/registry/generic"
	"k8s.io/kubernetes/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/runtime"
)

type REST struct {
	*registry.Store
}

// NewREST returns a RESTStorage object that will work against events.
func NewREST(opts generic.RESTOptions, ttl uint64) *REST {
	prefix := "/events"

	// We explicitly do NOT do any decoration here - switching on Cacher
	// for events will lead to too high memory consumption.
	storageInterface := opts.Storage

	store := &registry.Store{
		NewFunc:     func() runtime.Object { return &api.Event{} },
		NewListFunc: func() runtime.Object { return &api.EventList{} },
		KeyRootFunc: func(ctx api.Context) string {
			return registry.NamespaceKeyRootFunc(ctx, prefix)
		},
		KeyFunc: func(ctx api.Context, id string) (string, error) {
			return registry.NamespaceKeyFunc(ctx, prefix, id)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.Event).Name, nil
		},
		PredicateFunc: func(label labels.Selector, field fields.Selector) generic.Matcher {
			return event.MatchEvent(label, field)
		},
		TTLFunc: func(runtime.Object, uint64, bool) (uint64, error) {
			return ttl, nil
		},
		QualifiedResource:       api.Resource("events"),
		DeleteCollectionWorkers: opts.DeleteCollectionWorkers,

		CreateStrategy: event.Strategy,
		UpdateStrategy: event.Strategy,
		DeleteStrategy: event.Strategy,

		Storage: storageInterface,
	}
	return &REST{store}
}
