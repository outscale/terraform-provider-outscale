package okshelpers

import (
	"context"
	"fmt"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/oks"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlSerializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

var decoder = yamlSerializer.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

func FromYAML(manifest string) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	_, _, err := decoder.Decode([]byte(manifest), nil, obj)
	if err != nil {
		return nil, fmt.Errorf("decode yaml: %w", err)
	}

	return obj, nil
}

func GetResourceInterfaceFromManifest(ctx context.Context, oksClient *oks.Client, clusterID string, manifest string, timeout time.Duration) (*unstructured.Unstructured, dynamic.ResourceInterface, error) {
	obj, err := FromYAML(manifest)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid manifest: %w", err)
	}

	dri, err := getResourceInterface(ctx, oksClient, clusterID, obj, timeout)
	if err != nil {
		return nil, nil, fmt.Errorf("get resource interface: %w", err)
	}

	return obj, dri, nil
}

func ToYAML(obj map[string]any) (string, error) {
	bytes, err := yaml.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("serialize yaml: %w", err)
	}
	return string(bytes), nil
}

func BuildK8sClient(ctx context.Context, oksClient *oks.Client, clusterID string, timeout time.Duration) (dynamic.Interface, meta.RESTMapper, error) {
	resp, err := oksClient.GetKubeconfig(ctx, clusterID, nil, options.WithRetryTimeout(timeout))
	if err != nil {
		return nil, nil, err
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(resp.Cluster.Data.Kubeconfig))
	if err != nil {
		return nil, nil, fmt.Errorf("parse kubeconfig: %w", err)
	}
	config.Timeout = timeout
	config.QPS = 10

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("create dynamic client: %w", err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("create discovery client: %w", err)
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))

	return client, mapper, nil
}

func getResourceInterface(ctx context.Context, oksClient *oks.Client, clusterID string, obj *unstructured.Unstructured, timeout time.Duration) (dynamic.ResourceInterface, error) {
	dynamicClient, mapper, err := BuildK8sClient(ctx, oksClient, clusterID, timeout)
	if err != nil {
		return nil, err
	}

	gvk := obj.GroupVersionKind()

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	resourceClient := dynamicClient.Resource(mapping.Resource)
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		namespace := obj.GetNamespace()
		if namespace == "" {
			namespace = metav1.NamespaceDefault
			obj.SetNamespace(namespace)
		}

		return resourceClient.Namespace(namespace), nil
	}

	return resourceClient, nil
}
