package controller

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	"k8s.io/apimachinery/pkg/runtime/schema"
// 	"k8s.io/client-go/discovery"
// 	dynamicfake "k8s.io/client-go/dynamic/fake"
// 	"k8s.io/client-go/rest"
// )

// // newTestDynamicK8sClient initializes a DynamicK8sClient with a fake dynamic client and a fake discovery client.
// func newTestDynamicK8sClient(objects ...runtime.Object) *DynamicK8sClient {
// 	// Create a fake dynamic client
// 	dynamicClient := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), objects...)

// 	// Define the fake responses for the discovery endpoints
// 	handler := http.NewServeMux()

// 	// Handle /api endpoint
// 	handler.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
// 		apiVersions := v1.APIVersions{
// 			Versions: []string{"v1"},
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		err := json.NewEncoder(w).Encode(apiVersions)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	})

// 	// Handle /apis endpoint
// 	handler.HandleFunc("/apis", func(w http.ResponseWriter, r *http.Request) {
// 		apiGroupList := v1.APIGroupList{
// 			Groups: []v1.APIGroup{
// 				{
// 					Name: "apps",
// 					Versions: []v1.GroupVersionForDiscovery{
// 						{GroupVersion: "apps/v1", Version: "v1"},
// 					},
// 					PreferredVersion: v1.GroupVersionForDiscovery{
// 						GroupVersion: "apps/v1",
// 						Version:      "v1",
// 					},
// 				},
// 			},
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		err := json.NewEncoder(w).Encode(apiGroupList)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}

// 	})

// 	// Handle /api/v1 endpoint
// 	handler.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
// 		apiResourceList := v1.APIResourceList{
// 			GroupVersion: "v1",
// 			APIResources: []v1.APIResource{
// 				{Name: "pods", Kind: "Pod", Namespaced: true},
// 				{Name: "services", Kind: "Service", Namespaced: true},
// 				{Name: "configmaps", Kind: "ConfigMap", Namespaced: true},
// 			},
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		err := json.NewEncoder(w).Encode(apiResourceList)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}

// 	})

// 	// Handle /apis/apps/v1 endpoint
// 	handler.HandleFunc("/apis/apps/v1", func(w http.ResponseWriter, r *http.Request) {
// 		apiResourceList := v1.APIResourceList{
// 			GroupVersion: "apps/v1",
// 			APIResources: []v1.APIResource{
// 				{Name: "deployments", Kind: "Deployment", Namespaced: true},
// 			},
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		err := json.NewEncoder(w).Encode(apiResourceList)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	})

// 	// Create the fake server
// 	server := httptest.NewServer(handler)

// 	// Ensure the server is closed when tests complete
// 	// You can also defer this in each test if needed
// 	// defer server.Close()

// 	// Create a rest.Config pointing to the fake server
// 	config := &rest.Config{
// 		Host: server.URL,
// 		// Disable TLS verification for the fake server
// 	}

// 	// Create a DiscoveryClient with the fake server's config
// 	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
// 	if err != nil {
// 		panic(err) // Handle error appropriately in production code
// 	}

// 	return &DynamicK8sClient{
// 		client: dynamicClient,
// 		dc:     discoveryClient,
// 	}
// }

// func TestCreate(t *testing.T) {
// 	client := newTestDynamicK8sClient()

// 	// Create a sample unstructured Pod object
// 	obj := &unstructured.Unstructured{}
// 	obj.SetUnstructuredContent(map[string]interface{}{
// 		"apiVersion": "v1",
// 		"kind":       "Pod",
// 		"metadata": map[string]interface{}{
// 			"name":      "test-pod",
// 			"namespace": "default",
// 		},
// 	})

// 	// Set the GroupVersionKind
// 	obj.SetGroupVersionKind(schema.GroupVersionKind{
// 		Group:   "",
// 		Version: "v1",
// 		Kind:    "Pod",
// 	})

// 	// Test the Create method
// 	createdObj, err := client.Create(context.Background(), obj)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "test-pod", createdObj.GetName())
// 	assert.Equal(t, "default", createdObj.GetNamespace())
// }

// func TestGet(t *testing.T) {
// 	// Create a sample unstructured Pod object to seed the fake client
// 	obj := &unstructured.Unstructured{}
// 	obj.SetUnstructuredContent(map[string]interface{}{
// 		"apiVersion": "v1",
// 		"kind":       "Pod",
// 		"metadata": map[string]interface{}{
// 			"name":      "test-pod",
// 			"namespace": "default",
// 		},
// 	})

// 	// Set the GroupVersionKind
// 	obj.SetGroupVersionKind(schema.GroupVersionKind{
// 		Group:   "",
// 		Version: "v1",
// 		Kind:    "Pod",
// 	})

// 	client := newTestDynamicK8sClient(obj)

// 	// Test the Get method
// 	fetchedObj, err := client.Get(context.Background(), "", "Pod", "default", "test-pod")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "test-pod", fetchedObj.GetName())
// 	assert.Equal(t, "default", fetchedObj.GetNamespace())
// }

// func TestUpdate(t *testing.T) {
// 	// Seed the fake client with an existing ConfigMap
// 	obj := &unstructured.Unstructured{}
// 	obj.SetUnstructuredContent(map[string]interface{}{
// 		"apiVersion": "v1",
// 		"kind":       "ConfigMap",
// 		"metadata": map[string]interface{}{
// 			"name":      "test-configmap",
// 			"namespace": "default",
// 		},
// 		"data": map[string]interface{}{
// 			"key": "value",
// 		},
// 	})

// 	// Set the GroupVersionKind
// 	obj.SetGroupVersionKind(schema.GroupVersionKind{
// 		Group:   "",
// 		Version: "v1",
// 		Kind:    "ConfigMap",
// 	})

// 	client := newTestDynamicK8sClient(obj)

// 	// Update the ConfigMap's data
// 	obj.Object["data"] = map[string]interface{}{
// 		"key": "new-value",
// 	}

// 	// Test the Update method
// 	updatedObj, err := client.Update(context.Background(), "", "ConfigMap", "default", obj)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "new-value", updatedObj.Object["data"].(map[string]interface{})["key"])
// }

// func TestDelete(t *testing.T) {
// 	// Seed the fake client with an existing Service
// 	obj := &unstructured.Unstructured{}
// 	obj.SetUnstructuredContent(map[string]interface{}{
// 		"apiVersion": "v1",
// 		"kind":       "Service",
// 		"metadata": map[string]interface{}{
// 			"name":      "test-service",
// 			"namespace": "default",
// 		},
// 	})

// 	// Set the GroupVersionKind
// 	obj.SetGroupVersionKind(schema.GroupVersionKind{
// 		Group:   "",
// 		Version: "v1",
// 		Kind:    "Service",
// 	})

// 	client := newTestDynamicK8sClient(obj)

// 	// Test the Delete method
// 	err := client.Delete(context.Background(), "", "Service", "default", "test-service")
// 	assert.NoError(t, err)

// 	// Verify the Service is deleted
// 	_, err = client.Get(context.Background(), "", "Service", "default", "test-service")
// 	assert.Error(t, err)
// }

// func TestList(t *testing.T) {
// 	// Seed the fake client with multiple Deployments
// 	obj1 := &unstructured.Unstructured{}
// 	obj1.SetUnstructuredContent(map[string]interface{}{
// 		"apiVersion": "apps/v1",
// 		"kind":       "Deployment",
// 		"metadata": map[string]interface{}{
// 			"name":      "deployment1",
// 			"namespace": "default",
// 		},
// 	})

// 	// Set the GroupVersionKind
// 	obj1.SetGroupVersionKind(schema.GroupVersionKind{
// 		Group:   "apps",
// 		Version: "v1",
// 		Kind:    "Deployment",
// 	})

// 	obj2 := obj1.DeepCopy()
// 	obj2.SetName("deployment2")

// 	client := newTestDynamicK8sClient(obj1, obj2)

// 	// Test the List method
// 	list, err := client.List(context.Background(), "apps", "Deployment", "default")
// 	assert.NoError(t, err)
// 	assert.Equal(t, 2, len(list.Items))
// }

// func TestGetPreferredGVK(t *testing.T) {
// 	client := newTestDynamicK8sClient()

// 	// Test the GetPreferredGVK method
// 	gvk, err := client.GetPreferredGVK("apps", "Deployment")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "apps", gvk.Group)
// 	assert.Equal(t, "v1", gvk.Version)
// 	assert.Equal(t, "Deployment", gvk.Kind)
// }
