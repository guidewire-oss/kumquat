package controller

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// create getk8sclient function
func GetK8sClient() (K8sClient, error) {
	return NewDynamicK8sClient()
}

// K8sClient defines the interface for Kubernetes clients
type K8sClient interface {
	Create(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error)
	CreateOrUpdate(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error)
	List(ctx context.Context, group, kind, namespace string) (*unstructured.UnstructuredList, error)
	Get(ctx context.Context, group, kind, namespace, name string) (*unstructured.Unstructured, error)
	Update(ctx context.Context, group, kind, namespace string, obj *unstructured.Unstructured,
	) (*unstructured.Unstructured, error)
	Delete(ctx context.Context, group, kind, namespace, name string) error
	GetPreferredGVK(group, kind string) (schema.GroupVersionKind, error)
}

type DynamicK8sClient struct {
	client dynamic.Interface
	dc     *discovery.DiscoveryClient
}

// NewDynamicK8sClient initializes and returns a new DynamicK8sClient
func NewDynamicK8sClient() (*DynamicK8sClient, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes config: %v", err)
	}

	dc, err := discovery.NewDiscoveryClientForConfig(cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %v", err)
	}

	return &DynamicK8sClient{
		client: dynamicClient,
		dc:     dc,
	}, nil
}

func (k *DynamicK8sClient) getGVRByGroupAndKind(group, kind string) (schema.GroupVersionResource, error) {
	// List all resources for the group
	apiResourceLists, err := k.dc.ServerPreferredResources()
	if err != nil {
		return schema.GroupVersionResource{}, fmt.Errorf("failed to get server resources: %v", err)
	}

	for _, apiResourceList := range apiResourceLists {
		gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			return schema.GroupVersionResource{}, fmt.Errorf("failed to parse GroupVersion: %v", err)
		}
		fmt.Println(gv.Group, gv.Version, "this is gv", group, kind)

		// Check if the group matches
		if gv.Group == group {
			fmt.Printf("this is happening")

			for _, apiResource := range apiResourceList.APIResources {
				fmt.Println("yessss", apiResource.Kind, kind)

				// Check if the kind matches
				if apiResource.Kind == kind {
					return schema.GroupVersionResource{
						Group:    gv.Group,
						Version:  gv.Version, // use the version returned by ServerPreferredResources
						Resource: apiResource.Name,
					}, nil
				}
			}
		}
	}

	return schema.GroupVersionResource{}, fmt.Errorf("resource not found for group: %v and kind: %v", group, kind)
}

// Create creates a resource based on group and kind
func (k *DynamicK8sClient) Create(
	ctx context.Context,
	obj *unstructured.Unstructured,
) (*unstructured.Unstructured, error) {
	gvk := obj.GroupVersionKind()
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: strings.ToLower(gvk.Kind) + "s", // Simple pluralization
	}

	return k.client.Resource(gvr).Namespace(obj.GetNamespace()).Create(ctx, obj, v1.CreateOptions{})
}

func (k *DynamicK8sClient) CreateOrUpdate(
	ctx context.Context,
	obj *unstructured.Unstructured,
) (*unstructured.Unstructured, error) {
	gvk := obj.GroupVersionKind()
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: strings.ToLower(gvk.Kind) + "s", // Simple pluralization
	}

	_, err := k.client.Resource(gvr).Namespace(obj.GetNamespace()).Get(ctx, obj.GetName(), v1.GetOptions{})
	if err != nil {
		return k.client.Resource(gvr).Namespace(obj.GetNamespace()).Create(ctx, obj, v1.CreateOptions{})
	}

	return k.client.Resource(gvr).Namespace(obj.GetNamespace()).Update(ctx, obj, v1.UpdateOptions{})
}

func (k *DynamicK8sClient) List(
	ctx context.Context,
	group, kind, namespace string,
) (*unstructured.UnstructuredList, error) {
	gvr, err := k.getGVRByGroupAndKind(group, kind)
	if err != nil {
		return nil, err
	}
	return k.client.Resource(gvr).Namespace(namespace).List(ctx, v1.ListOptions{})
}

// Get retrieves a resource based on group, kind, and name
func (k *DynamicK8sClient) Get(
	ctx context.Context,
	group, kind, namespace, name string,
) (*unstructured.Unstructured, error) {
	gvr, err := k.getGVRByGroupAndKind(group, kind)
	fmt.Println(gvr.Group, gvr.Resource, gvr.Version, "this is gvr")
	if err != nil {
		return nil, err
	}
	fmt.Println("this is namespaceeeeee", namespace)
	return k.client.Resource(gvr).Namespace(namespace).Get(ctx, name, v1.GetOptions{})
}

// Update updates a resource based on group and kind
func (k *DynamicK8sClient) Update(
	ctx context.Context,
	group, kind, namespace string,
	obj *unstructured.Unstructured,
) (*unstructured.Unstructured, error) {
	gvr, err := k.getGVRByGroupAndKind(group, kind)
	if err != nil {
		return nil, err
	}
	return k.client.Resource(gvr).Namespace(namespace).Update(ctx, obj, v1.UpdateOptions{})
}

// Delete deletes a resource based on group, kind, and name
func (k *DynamicK8sClient) Delete(ctx context.Context, group, kind, namespace, name string) error {
	gvr, err := k.getGVRByGroupAndKind(group, kind)
	if err != nil {
		return err
	}
	return k.client.Resource(gvr).Namespace(namespace).Delete(ctx, name, v1.DeleteOptions{})
}

// getPreferredGVK returns the preferred GroupVersionKind for a given group and kind
func (k *DynamicK8sClient) GetPreferredGVK(group, kind string) (schema.GroupVersionKind, error) {
	// List preferred resources available on the server
	apiResourceLists, err := k.dc.ServerPreferredResources()
	if err != nil {
		return schema.GroupVersionKind{}, fmt.Errorf("failed to get server resources: %v", err)
	}

	for _, apiResourceList := range apiResourceLists {
		gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			return schema.GroupVersionKind{}, fmt.Errorf("failed to parse GroupVersion: %v", err)
		}

		// Check if the group matches
		if gv.Group == group {
			for _, apiResource := range apiResourceList.APIResources {
				// Check if the kind matches
				if apiResource.Kind == kind {
					// Return the GVK with the preferred version
					return schema.GroupVersionKind{
						Group:   gv.Group,
						Version: gv.Version,
						Kind:    apiResource.Kind,
					}, nil
				}
			}
		}
	}

	return schema.GroupVersionKind{}, fmt.Errorf("preferred GVK not found for group: %v and kind: %v", group, kind)
}
