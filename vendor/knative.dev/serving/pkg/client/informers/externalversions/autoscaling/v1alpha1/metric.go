/*
Copyright 2019 The Knative Authors

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	autoscalingv1alpha1 "knative.dev/serving/pkg/apis/autoscaling/v1alpha1"
	versioned "knative.dev/serving/pkg/client/clientset/versioned"
	internalinterfaces "knative.dev/serving/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "knative.dev/serving/pkg/client/listers/autoscaling/v1alpha1"
)

// MetricInformer provides access to a shared informer and lister for
// Metrics.
type MetricInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.MetricLister
}

type metricInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewMetricInformer constructs a new informer for Metric type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewMetricInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredMetricInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredMetricInformer constructs a new informer for Metric type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredMetricInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AutoscalingV1alpha1().Metrics(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AutoscalingV1alpha1().Metrics(namespace).Watch(options)
			},
		},
		&autoscalingv1alpha1.Metric{},
		resyncPeriod,
		indexers,
	)
}

func (f *metricInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredMetricInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *metricInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&autoscalingv1alpha1.Metric{}, f.defaultInformer)
}

func (f *metricInformer) Lister() v1alpha1.MetricLister {
	return v1alpha1.NewMetricLister(f.Informer().GetIndexer())
}
