package controller

import (
	"context"
	"fmt"
	"kumquat/repository"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type WatchManager interface {
	UpdateGeneratedResources(templateName string, resourceSet mapset.Set[ResourceIdentifier])
	UpdateWatch(templateName string, newGVKs []schema.GroupVersionKind) error
	RemoveWatch(templateName string)
	GetGeneratedResources(templateName string) mapset.Set[ResourceIdentifier]
	GetManagedTemplates() map[string]map[schema.GroupVersionKind]struct{}
}

// ControllerEntry represents a dynamically managed controller.
type ControllerEntry struct {
	controller controller.Controller
	cancelFunc context.CancelFunc
	ctx        context.Context
}

type ResourceIdentifier struct {
	Group     string
	Kind      string
	Namespace string
	Name      string
}

func (r *ResourceIdentifier) ToString() string {
	return r.Group + "/" + r.Kind + "/" + r.Namespace + "/" + r.Name
}

// WatchManager manages dynamic watches.
type watchManager struct {
	refCounts          map[schema.GroupVersionKind]int
	watchedResources   map[schema.GroupVersionKind]ControllerEntry
	templates          map[string]map[schema.GroupVersionKind]struct{}
	mu                 sync.Mutex
	cache              cache.Cache
	client             client.Client
	scheme             *runtime.Scheme
	mgr                manager.Manager
	K8sClient          K8sClient
	generatedResources map[string]mapset.Set[ResourceIdentifier]
	repository         repository.Repository
}

// NewWatchManager creates a new WatchManager instance.
func NewWatchManager(mgr manager.Manager, k8sClient K8sClient, repo repository.Repository) WatchManager {
	wm := &watchManager{
		watchedResources:   make(map[schema.GroupVersionKind]ControllerEntry),
		refCounts:          make(map[schema.GroupVersionKind]int),
		templates:          make(map[string]map[schema.GroupVersionKind]struct{}),
		generatedResources: make(map[string]mapset.Set[ResourceIdentifier]),
		cache:              mgr.GetCache(),
		scheme:             mgr.GetScheme(),
		mgr:                mgr,
		K8sClient:          k8sClient,
		client:             mgr.GetClient(),
		repository:         repo,
	}
	return wm
}

func (wm *watchManager) GetManagedTemplates() map[string]map[schema.GroupVersionKind]struct{} {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	return wm.templates
}

// AddWatch adds a watch for the specified template and GVKs.
func (wm *watchManager) AddWatch(templateName string, gvks []schema.GroupVersionKind) error {
	if _, exists := wm.templates[templateName]; !exists {
		wm.templates[templateName] = make(map[schema.GroupVersionKind]struct{})
	}

	for _, gvk := range gvks {
		if _, exists := wm.templates[templateName][gvk]; exists {
			continue
		}
		wm.templates[templateName][gvk] = struct{}{}
		if wm.refCounts[gvk] == 0 {
			tableName := gvk.Kind + "." + gvk.Group
			if err := deleteTableFromDataBase(wm.repository, tableName); err != nil {
				return err
			}
			if err := wm.startWatching(gvk); err != nil {
				return err
			}
		}
		wm.refCounts[gvk]++
		log.Log.Info("Incremented watch reference count", "gvk", gvk, "count", wm.refCounts[gvk])
	}

	return nil
}
func (wm *watchManager) UpdateGeneratedResources(templateName string, resources mapset.Set[ResourceIdentifier]) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.generatedResources[templateName] = resources
}

// UpdateWatch updates the watch for the specified template with new GVKs.
func (wm *watchManager) UpdateWatch(templateName string, newGVKs []schema.GroupVersionKind) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if _, exists := wm.templates[templateName]; !exists {
		log.Log.Info("Resource template not found", "templateName", templateName)
		return wm.AddWatch(templateName, newGVKs)
	}

	oldGVKs := wm.templates[templateName]
	removedGVKs := make(map[schema.GroupVersionKind]struct{})
	addedGVKs := make(map[schema.GroupVersionKind]struct{})

	for gvk := range oldGVKs {
		removedGVKs[gvk] = struct{}{}
	}
	for _, gvk := range newGVKs {
		if _, exists := removedGVKs[gvk]; exists {
			delete(removedGVKs, gvk)
		} else {
			addedGVKs[gvk] = struct{}{}
		}
	}

	for gvk := range removedGVKs {
		wm.removeWatchForGVK(templateName, gvk)
	}

	for gvk := range addedGVKs {
		if err := wm.addWatchForGVK(templateName, gvk); err != nil {
			return err
		}
	}

	return nil
}

// RemoveWatch removes the watch for the specified template.
func (wm *watchManager) RemoveWatch(templateName string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	log.Log.Info("Removing watch", "templateName", templateName)

	if watchedGVKs, exists := wm.templates[templateName]; exists {
		for gvk := range watchedGVKs {
			wm.refCounts[gvk]--
			if wm.refCounts[gvk] <= 0 {
				wm.stopWatching(gvk)
				delete(wm.refCounts, gvk)
			}
		}
		delete(wm.templates, templateName)
	}
	for gvk, count := range wm.refCounts {
		log.Log.Info("Reference count", "gvk", gvk, "count", count)
	}

	wm.logActiveControllers()
}

// addWatchForGVK adds a watch for a specific GVK.
func (wm *watchManager) addWatchForGVK(templateName string, gvk schema.GroupVersionKind) error {
	wm.templates[templateName][gvk] = struct{}{}
	if wm.refCounts[gvk] == 0 {
		if err := wm.startWatching(gvk); err != nil {
			log.Log.Error(err, "unable to start watching", "gvk", gvk)
			return err
		}
	}
	wm.refCounts[gvk]++
	log.Log.Info("Incremented watch reference count", "gvk", gvk, "count", wm.refCounts[gvk])
	return nil
}

// removeWatchForGVK removes a watch for a specific GVK.
func (wm *watchManager) removeWatchForGVK(templateName string, gvk schema.GroupVersionKind) {
	wm.refCounts[gvk]--
	if wm.refCounts[gvk] <= 0 {
		wm.stopWatching(gvk)
		delete(wm.refCounts, gvk)
	}
	delete(wm.templates[templateName], gvk)
	log.Log.Info("Decremented watch reference count", "gvk", gvk, "count", wm.refCounts[gvk])
}

// startWatching starts watching a specific GVK.
func (wm *watchManager) startWatching(gvk schema.GroupVersionKind) error {
	log.Log.Info("Starting watch", "gvk", gvk)
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)
	dynamicReconciler := NewDynamicReconciler(wm.client, gvk, wm.K8sClient, wm, wm.repository)

	c, err := controller.NewUnmanaged("dynamic-controller-"+gvk.Kind, wm.mgr, controller.Options{
		Reconciler: dynamicReconciler,

		// Skip the name check introduced in v0.19.0 of controller-runtime via
		// https://github.com/kubernetes-sigs/controller-runtime/pull/2902; we managed the controller lifecycle
		// ourselves and it is not necessary to have unique names.
		SkipNameValidation: ptr.To(true),
	})
	if err != nil {
		fmt.Printf("Error creating controller: %v\n", err)
		return err
	}

	kindSource := source.Kind(wm.mgr.GetCache(), obj, &unstructuredEventHandler{})
	err = c.Watch(kindSource)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	wm.watchedResources[gvk] = ControllerEntry{controller: c, cancelFunc: cancelFunc, ctx: ctx}
	go func() {
		if err := c.Start(ctx); err != nil && err != context.Canceled {
			log.Log.Error(err, "unable to start controller", "gvk", gvk)
		}
	}()
	return nil
}

// stopWatching stops watching a specific GVK.
func (wm *watchManager) stopWatching(gvk schema.GroupVersionKind) {
	log.Log.Info("Stopping watch", "gvk", gvk)
	if entry, exists := wm.watchedResources[gvk]; exists {
		entry.cancelFunc()
		<-entry.ctx.Done()
		delete(wm.watchedResources, gvk)
	}
}

// logs all active controllers.
func (wm *watchManager) logActiveControllers() {
	log.Log.Info("Listing all active controllers:")
	for gvk, entry := range wm.watchedResources {
		log.Log.Info("Active controller", "gvk", gvk, "context", entry.ctx)
	}
}

// DeleteRecord deletes a record from the specified table.
func DeleteRecord(table, namespace, name string, repo repository.Repository) error {
	err := repo.Delete(namespace, name, table)
	if err != nil {
		log.Log.Error(err, "unable to delete record")
		return err
	}
	log.Log.Info("Record deleted", "table", table, "namespace", namespace, "name", name)
	return nil
}

func (wm *watchManager) GetGeneratedResources(templateName string) mapset.Set[ResourceIdentifier] {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	return wm.generatedResources[templateName]
}

// unstructuredEventHandler handles events for unstructured resources.
type unstructuredEventHandler struct{}

func (h *unstructuredEventHandler) Create(
	ctx context.Context,
	evt event.TypedCreateEvent[*unstructured.Unstructured],
	q workqueue.TypedRateLimitingInterface[ctrl.Request],
) {
	q.Add(ctrl.Request{NamespacedName: client.ObjectKeyFromObject(evt.Object)})
}

func (h *unstructuredEventHandler) Update(
	ctx context.Context,
	evt event.TypedUpdateEvent[*unstructured.Unstructured],
	q workqueue.TypedRateLimitingInterface[ctrl.Request],
) {
	q.Add(ctrl.Request{NamespacedName: client.ObjectKeyFromObject(evt.ObjectNew)})
}

func (h *unstructuredEventHandler) Delete(
	ctx context.Context,
	evt event.TypedDeleteEvent[*unstructured.Unstructured],
	q workqueue.TypedRateLimitingInterface[ctrl.Request],
) {
	q.Add(ctrl.Request{NamespacedName: client.ObjectKeyFromObject(evt.Object)})
}

func (h *unstructuredEventHandler) Generic(
	ctx context.Context,
	evt event.TypedGenericEvent[*unstructured.Unstructured],
	q workqueue.TypedRateLimitingInterface[ctrl.Request],
) {
	q.Add(ctrl.Request{NamespacedName: client.ObjectKeyFromObject(evt.Object)})
}
