package controller

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewDynamicK8sClient(client client.Client, restMapper meta.RESTMapper) (K8sClient, error) {

	return &DynamicK8sClient{
		client:     client,
		restMapper: restMapper,
	}, nil
}

// K8sClient interface remains the same
type K8sClient interface {
	Create(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error)
	CreateOrUpdate(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error)
	List(ctx context.Context, group, kind, namespace string) (*unstructured.UnstructuredList, error)
	Get(ctx context.Context, group, kind, namespace, name string) (*unstructured.Unstructured, error)
	Update(ctx context.Context, group, kind, namespace string, obj *unstructured.Unstructured) (
		*unstructured.Unstructured, error)
	Delete(ctx context.Context, group, kind, namespace, name string) error
	GetPreferredGVK(group, kind string) (schema.GroupVersionKind, error)
}

type DynamicK8sClient struct {
	client     client.Client
	restMapper meta.RESTMapper
}

// Implement the methods using client.Client

func (k *DynamicK8sClient) Create(ctx context.Context, obj *unstructured.Unstructured) (
	*unstructured.Unstructured, error) {

	err := k.client.Create(ctx, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (k *DynamicK8sClient) CreateOrUpdate(ctx context.Context, obj *unstructured.Unstructured) (
	*unstructured.Unstructured, error) {

	existing := &unstructured.Unstructured{}
	existing.SetGroupVersionKind(obj.GroupVersionKind())
	key := client.ObjectKey{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}

	err := k.client.Get(ctx, key, existing)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return nil, err
		}
		// Not found, create
		err = k.client.Create(ctx, obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}

	// Resource exists, update
	obj.SetResourceVersion(existing.GetResourceVersion())
	err = k.client.Update(ctx, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (k *DynamicK8sClient) List(ctx context.Context, group, kind, namespace string) (
	*unstructured.UnstructuredList, error) {

	gvk, err := k.GetPreferredGVK(group, kind)
	if err != nil {
		return nil, err
	}

	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(gvk)
	err = k.client.List(ctx, list, client.InNamespace(namespace))
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (k *DynamicK8sClient) Get(ctx context.Context, group, kind, namespace, name string) (
	*unstructured.Unstructured, error) {
	gvk, err := k.GetPreferredGVK(group, kind)
	if err != nil {
		return nil, err
	}

	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
	err = k.client.Get(ctx, key, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (k *DynamicK8sClient) Update(ctx context.Context, group, kind, namespace string, obj *unstructured.Unstructured) (
	*unstructured.Unstructured, error) {
	err := k.client.Update(ctx, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (k *DynamicK8sClient) Delete(ctx context.Context, group, kind, namespace, name string) error {
	gvk, err := k.GetPreferredGVK(group, kind)
	if err != nil {
		return err
	}

	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)
	obj.SetNamespace(namespace)
	obj.SetName(name)
	err = k.client.Delete(ctx, obj)
	if err != nil {
		return err
	}
	return nil
}

func (k *DynamicK8sClient) GetPreferredGVK(group, kind string) (schema.GroupVersionKind, error) {
	partialGVK := schema.GroupVersionKind{
		Group: group,
		Kind:  kind,
	}

	mapping, err := k.restMapper.RESTMapping(partialGVK.GroupKind())
	if err != nil {
		return schema.GroupVersionKind{}, fmt.Errorf("failed to get GVK from RESTMapper: %v", err)
	}

	return mapping.GroupVersionKind, nil

}
