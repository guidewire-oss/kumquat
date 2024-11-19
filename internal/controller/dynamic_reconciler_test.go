package controller

import (
	"context"
	"errors"

	kumquatv1beta1 "kumquat/api/v1beta1"
	"kumquat/repository"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("DynamicReconciler", func() {
	var (
		ctx                              context.Context
		scheme                           *runtime.Scheme
		fakeClient                       client.Client
		gvk                              schema.GroupVersionKind
		testResource                     *unstructured.Unstructured
		mockK8sClient                    *MockK8sClient
		mockWatchManager                 *MockWatchManager
		mockRepository                   *MockRepository
		dynamicReconciler                *DynamicReconciler
		originalProcessTemplateResources func(
			templateObj *kumquatv1beta1.Template,
			sr repository.Repository,
			log logr.Logger,
			k8sClient K8sClient,
			wm WatchManagerInterface,
		) error
	)

	BeforeEach(func() {
		ctx = context.TODO()

		// Set up the scheme
		scheme = runtime.NewScheme()
		Expect(kumquatv1beta1.AddToScheme(scheme)).To(Succeed())

		// Create mocks for dependencies
		mockK8sClient = &MockK8sClient{}
		mockWatchManager = &MockWatchManager{}
		mockRepository = &MockRepository{}

		// Save the original ProcessTemplateResources
		originalProcessTemplateResources = ProcessTemplateResources
	})

	AfterEach(func() {
		// Restore the original ProcessTemplateResources
		ProcessTemplateResources = originalProcessTemplateResources
	})

	Context("When reconciling resources", func() {
		When("the Template resource exists", func() {
			BeforeEach(func() {
				// Define GVK for Template
				gvk = schema.GroupVersionKind{
					Group:   "kumquat.guidewire.com",
					Version: "v1beta1",
					Kind:    "Template",
				}

				// Define test resource as Template
				testResource = &unstructured.Unstructured{}
				testResource.SetGroupVersionKind(gvk)
				testResource.SetName("test-template")
				testResource.SetNamespace("default")

				// Create a fake client with the testResource
				fakeClient = fake.NewClientBuilder().WithScheme(scheme).WithObjects(testResource.DeepCopy()).Build()

				// Create the DynamicReconciler with the fake client
				dynamicReconciler = NewDynamicReconciler(fakeClient, gvk, mockK8sClient, mockWatchManager, mockRepository)

				// Mock K8sClient.Get to return a Template
				mockK8sClient.GetFunc = func(ctx context.Context, group, kind, namespace, name string) (*unstructured.Unstructured, error) {
					if group == "kumquat.guidewire.com" && kind == "Template" {
						// Return a Template object
						template := &unstructured.Unstructured{}
						template.SetGroupVersionKind(schema.GroupVersionKind{
							Group:   "kumquat.guidewire.com",
							Version: "v1beta1",
							Kind:    "Template",
						})
						template.SetName(name)
						template.SetNamespace(namespace)
						template.Object["spec"] = map[string]interface{}{
							"query": "SELECT * FROM test_table",
						}
						// Do NOT set resourceVersion here
						return template, nil
					}
					return nil, errors.New("not found")
				}

				// Mock Repository.Upsert to succeed
				mockRepository.UpsertFunc = func(resource repository.Resource) error {
					return nil
				}

				// Mock WatchManager.GetManagedTemplates
				mockWatchManager.GetManagedTemplatesFunc = func() map[string]map[schema.GroupVersionKind]struct{} {
					return map[string]map[schema.GroupVersionKind]struct{}{
						"template1": {
							gvk: {},
						},
					}
				}

				// Mock ProcessTemplateResources to do nothing
				ProcessTemplateResources = func(
					templateObj *kumquatv1beta1.Template,
					sr repository.Repository,
					log logr.Logger,
					k8sClient K8sClient,
					wm WatchManagerInterface,
				) error {
					return nil
				}
			})

			It("should upsert the resource and process templates without error", func() {
				req := reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      testResource.GetName(),
						Namespace: testResource.GetNamespace(),
					},
				}
				_, err := dynamicReconciler.Reconcile(ctx, req)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("the Template resource does not exist (deleted)", func() {
			BeforeEach(func() {
				// Define GVK for Template
				gvk = schema.GroupVersionKind{
					Group:   "kumquat.guidewire.com",
					Version: "v1beta1",
					Kind:    "Template",
				}

				// Define test resource as Template (not adding to fake client)
				testResource = &unstructured.Unstructured{}
				testResource.SetGroupVersionKind(gvk)
				testResource.SetName("test-template")
				testResource.SetNamespace("default")

				// Create a fake client without the testResource
				fakeClient = fake.NewClientBuilder().WithScheme(scheme).WithObjects().Build()

				// Create the DynamicReconciler with the fake client
				dynamicReconciler = NewDynamicReconciler(fakeClient, gvk, mockK8sClient, mockWatchManager, mockRepository)

				// Mock K8sClient.Get to return not found
				mockK8sClient.GetFunc = func(ctx context.Context, group, kind, namespace, name string) (*unstructured.Unstructured, error) {
					if group == "kumquat.guidewire.com" && kind == "Template" {
						return nil, errors.New("not found")
					}
					return nil, errors.New("not found")
				}

				// Mock Repository.Delete to succeed
				mockRepository.DeleteFunc = func(namespace, name, table string) error {
					return nil
				}

				// Mock WatchManager.GetManagedTemplates
				mockWatchManager.GetManagedTemplatesFunc = func() map[string]map[schema.GroupVersionKind]struct{} {
					return map[string]map[schema.GroupVersionKind]struct{}{
						"template1": {
							gvk: {},
						},
					}
				}

				// Mock ProcessTemplateResources to do nothing
				ProcessTemplateResources = func(
					templateObj *kumquatv1beta1.Template,
					sr repository.Repository,
					log logr.Logger,
					k8sClient K8sClient,
					wm WatchManagerInterface,
				) error {
					return nil
				}
			})

			It("should delete the resource from the repository and process templates without error", func() {
				req := reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      testResource.GetName(),
						Namespace: testResource.GetNamespace(),
					},
				}
				_, err := dynamicReconciler.Reconcile(ctx, req)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

// MockK8sClient implements the K8sClient interface for testing
type MockK8sClient struct {
	CreateFunc          func(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error)
	CreateOrUpdateFunc  func(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error)
	ListFunc            func(ctx context.Context, group, kind, namespace string) (*unstructured.UnstructuredList, error)
	GetFunc             func(ctx context.Context, group, kind, namespace, name string) (*unstructured.Unstructured, error)
	UpdateFunc          func(ctx context.Context, group, kind, namespace string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error)
	DeleteFunc          func(ctx context.Context, group, kind, namespace, name string) error
	GetPreferredGVKFunc func(group, kind string) (schema.GroupVersionKind, error)
}

func (m *MockK8sClient) Create(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, obj)
	}
	return obj, nil
}

func (m *MockK8sClient) CreateOrUpdate(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	if m.CreateOrUpdateFunc != nil {
		return m.CreateOrUpdateFunc(ctx, obj)
	}
	return obj, nil
}

func (m *MockK8sClient) List(ctx context.Context, group, kind, namespace string) (*unstructured.UnstructuredList, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, group, kind, namespace)
	}
	return &unstructured.UnstructuredList{}, nil
}

func (m *MockK8sClient) Get(ctx context.Context, group, kind, namespace, name string) (*unstructured.Unstructured, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, group, kind, namespace, name)
	}
	return &unstructured.Unstructured{}, nil
}

func (m *MockK8sClient) Update(ctx context.Context, group, kind, namespace string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, group, kind, namespace, obj)
	}
	return obj, nil
}

func (m *MockK8sClient) Delete(ctx context.Context, group, kind, namespace, name string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, group, kind, namespace, name)
	}
	return nil
}

func (m *MockK8sClient) GetPreferredGVK(group, kind string) (schema.GroupVersionKind, error) {
	if m.GetPreferredGVKFunc != nil {
		return m.GetPreferredGVKFunc(group, kind)
	}
	return schema.GroupVersionKind{
		Group:   group,
		Version: "v1",
		Kind:    kind,
	}, nil
}

func (m *MockK8sClient) Reset() {
	m.CreateFunc = nil
	m.CreateOrUpdateFunc = nil
	m.ListFunc = nil
	m.GetFunc = nil
	m.UpdateFunc = nil
	m.DeleteFunc = nil
	m.GetPreferredGVKFunc = nil
}

// MockWatchManager implements the WatchManagerInterface for testing
type MockWatchManager struct {
	UpdateGeneratedResourcesFunc func(templateName string, resourceSet mapset.Set[ResourceIdentifier])
	UpdateWatchFunc              func(templateName string, newGVKs []schema.GroupVersionKind) error
	RemoveWatchFunc              func(templateName string)
	GetGeneratedResourcesFunc    func(templateName string) mapset.Set[ResourceIdentifier]
	GetManagedTemplatesFunc      func() map[string]map[schema.GroupVersionKind]struct{}
}

func (m *MockWatchManager) UpdateGeneratedResources(templateName string, resourceSet mapset.Set[ResourceIdentifier]) {
	if m.UpdateGeneratedResourcesFunc != nil {
		m.UpdateGeneratedResourcesFunc(templateName, resourceSet)
	}
}

func (m *MockWatchManager) UpdateWatch(templateName string, newGVKs []schema.GroupVersionKind) error {
	if m.UpdateWatchFunc != nil {
		return m.UpdateWatchFunc(templateName, newGVKs)
	}
	return nil
}

func (m *MockWatchManager) RemoveWatch(templateName string) {
	if m.RemoveWatchFunc != nil {
		m.RemoveWatchFunc(templateName)
	}
}

func (m *MockWatchManager) GetGeneratedResources(templateName string) mapset.Set[ResourceIdentifier] {
	if m.GetGeneratedResourcesFunc != nil {
		return m.GetGeneratedResourcesFunc(templateName)
	}
	return mapset.NewSet[ResourceIdentifier]()
}

func (m *MockWatchManager) GetManagedTemplates() map[string]map[schema.GroupVersionKind]struct{} {
	if m.GetManagedTemplatesFunc != nil {
		return m.GetManagedTemplatesFunc()
	}
	return map[string]map[schema.GroupVersionKind]struct{}{}
}

func (m *MockWatchManager) Reset() {
	m.UpdateGeneratedResourcesFunc = nil
	m.UpdateWatchFunc = nil
	m.RemoveWatchFunc = nil
	m.GetGeneratedResourcesFunc = nil
	m.GetManagedTemplatesFunc = nil
}

// MockRepository implements the repository.Repository interface for testing
type MockRepository struct {
	UpsertFunc func(resource repository.Resource) error
	DeleteFunc func(namespace, name, table string) error
}

func (m *MockRepository) Upsert(resource repository.Resource) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(resource)
	}
	return nil
}

func (m *MockRepository) Delete(namespace, name, table string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(namespace, name, table)
	}
	return nil
}

func (m *MockRepository) Query(query string) (repository.ResultSet, error) {
	return repository.ResultSet{}, nil
}

func (m *MockRepository) Close() error {
	return nil
}

func (m *MockRepository) ExtractTableNamesFromQuery(query string) []string {
	return []string{}
}

func (m *MockRepository) DropTable(table string) error {
	return nil
}

func (m *MockRepository) Reset() {
	m.UpsertFunc = nil
	m.DeleteFunc = nil
}
