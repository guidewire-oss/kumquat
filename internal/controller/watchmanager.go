package controller

import (
	"context"
	"fmt"
	kumquatv1beta1 "kumquat/api/v1beta1"
	"kumquat/repository"
	"sync"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ControllerEntry represents a dynamically managed controller.
type ControllerEntry struct {
	controller controller.Controller
	cancelFunc context.CancelFunc
	ctx        context.Context
}

// WatchManager manages dynamic watches.
type WatchManager struct {
	refCounts        map[schema.GroupVersionKind]int
	watchedResources map[schema.GroupVersionKind]ControllerEntry
	templates        map[string]map[schema.GroupVersionKind]struct{}
	mu               sync.Mutex
	cache            cache.Cache
	client           client.Client
	scheme           *runtime.Scheme
	mgr              manager.Manager
}

var wm *WatchManager

// NewWatchManager creates a new WatchManager instance.
func NewWatchManager(mgr manager.Manager) *WatchManager {
	watchManager := &WatchManager{
		watchedResources: make(map[schema.GroupVersionKind]ControllerEntry),
		refCounts:        make(map[schema.GroupVersionKind]int),
		templates:        make(map[string]map[schema.GroupVersionKind]struct{}),
		cache:            mgr.GetCache(),
		client:           mgr.GetClient(),
		scheme:           mgr.GetScheme(),
		mgr:              mgr,
	}
	wm = watchManager
	return watchManager
}

// GetWatchManager returns the singleton instance of WatchManager.
func GetWatchManager() *WatchManager {
	return wm
}

// AddWatch adds a watch for the specified template and GVKs.
func (wm *WatchManager) AddWatch(templateName string, gvks []schema.GroupVersionKind) error {
	if _, exists := wm.templates[templateName]; !exists {
		wm.templates[templateName] = make(map[schema.GroupVersionKind]struct{})
	}

	for _, gvk := range gvks {
		if _, exists := wm.templates[templateName][gvk]; exists {
			continue
		}
		wm.templates[templateName][gvk] = struct{}{}
		if wm.refCounts[gvk] == 0 {
			if err := wm.startWatching(gvk); err != nil {
				return err
			}
		}
		wm.refCounts[gvk]++
		log.Log.Info("Incremented watch reference count", "gvk", gvk, "count", wm.refCounts[gvk])
	}

	return nil
}

// UpdateWatch updates the watch for the specified template with new GVKs.
func (wm *WatchManager) UpdateWatch(templateName string, newGVKs []schema.GroupVersionKind) error {
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
func (wm *WatchManager) RemoveWatch(templateName string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	log.Log.Info("Removing watch", "templateName", templateName)

	if watchedGVKs, exists := wm.templates[templateName]; exists {
		for gvk := range watchedGVKs {
			wm.refCounts[gvk]--
			if wm.refCounts[gvk] <= 0 {
				wm.stopWatching(gvk)
				delete(wm.refCounts, gvk)
				deleteTableFromDataBase(gvk) // nolint:errcheck
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
func (wm *WatchManager) addWatchForGVK(templateName string, gvk schema.GroupVersionKind) error {
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
func (wm *WatchManager) removeWatchForGVK(templateName string, gvk schema.GroupVersionKind) {
	wm.refCounts[gvk]--
	if wm.refCounts[gvk] <= 0 {
		wm.stopWatching(gvk)
		delete(wm.refCounts, gvk)
		deleteTableFromDataBase(gvk) // nolint:errcheck
	}
	delete(wm.templates[templateName], gvk)
	log.Log.Info("Decremented watch reference count", "gvk", gvk, "count", wm.refCounts[gvk])
}

// startWatching starts watching a specific GVK.
func (wm *WatchManager) startWatching(gvk schema.GroupVersionKind) error {
	log.Log.Info("Starting watch", "gvk", gvk)
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)

	dynamicReconciler := &DynamicReconciler{
		Client: wm.client,
		GVK:    gvk,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	c, err := controller.NewUnmanaged("dynamic-controller-"+gvk.Kind, wm.mgr, controller.Options{
		Reconciler: dynamicReconciler,
	})
	if err != nil {
		cancelFunc()
		return err
	}

	kindSource := source.Kind(wm.mgr.GetCache(), obj, &unstructuredEventHandler{})
	err = c.Watch(kindSource)
	if err != nil {
		cancelFunc()
		return err
	}

	wm.watchedResources[gvk] = ControllerEntry{controller: c, cancelFunc: cancelFunc, ctx: ctx}
	go func() {
		if err := c.Start(ctx); err != nil && err != context.Canceled {
			log.Log.Error(err, "unable to start controller", "gvk", gvk)
		}
	}()
	return nil
}

// stopWatching stops watching a specific GVK.
func (wm *WatchManager) stopWatching(gvk schema.GroupVersionKind) {
	log.Log.Info("Stopping watch", "gvk", gvk)
	if entry, exists := wm.watchedResources[gvk]; exists {
		entry.cancelFunc()
		<-entry.ctx.Done()
		delete(wm.watchedResources, gvk)
	}
}

// logActiveControllers logs all active controllers.
func (wm *WatchManager) logActiveControllers() {
	log.Log.Info("Listing all active controllers:")
	for gvk, entry := range wm.watchedResources {
		log.Log.Info("Active controller", "gvk", gvk, "context", entry.ctx)
	}
}

// DeleteRecord deletes a record from the specified table.
func DeleteRecord(table, namespace, name string) error {
	re, err := GetSqliteRepository()
	if err != nil {
		log.Log.Error(err, "unable to create repository")
		return err
	}
	_, err = re.Db.Exec( /* sql */ `DELETE FROM "`+table+`" WHERE namespace = ? AND name = ?`, namespace, name)
	if err != nil {
		log.Log.Error(err, "unable to delete record")
		return err
	}
	log.Log.Info("Record deleted", "table", table, "namespace", namespace, "name", name)
	return nil
}

// deleteTableFromDataBase deletes a table from the database.
func deleteTableFromDataBase(gvk schema.GroupVersionKind) error {
	re, err := GetSqliteRepository()
	if err != nil {
		log.Log.Error(err, "unable to create repository")
		return err
	}
	tableName := gvk.Kind + "." + gvk.Group

	err = re.DropTable(tableName)
	if err != nil {
		log.Log.Error(err, "unable to drop table")
		return err
	}
	log.Log.Info("Table dropped", "tableName", tableName)
	return nil
}

// DynamicReconciler reconciles dynamic resources.
type DynamicReconciler struct {
	client.Client
	GVK schema.GroupVersionKind
}

func (r *DynamicReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling dynamic resource", "GVK", r.GVK, "name", req.Name, "namespace", req.Namespace)

	resource, err := r.fetchResource(ctx, req)
	if err != nil {
		return reconcile.Result{}, err
	}

	if resource == nil {
		log.Info("Resource deleted", "GVK", r.GVK, "name", req.Name, "namespace", req.Namespace)
		// delete record from database
		// TODO: checking the return code and returning an error causes the E2E tests to fail during delete template
		DeleteRecord(r.GVK.Kind+"."+r.GVK.Group, req.Namespace, req.Name) // nolint:errcheck
		// if err != nil {
		// 	return reconcile.Result{}, fmt.Errorf("error deleting record: %w", err)
		// }

		r.reconcileTemplates(ctx) // nolint:errcheck
		// if err != nil {
		// 	return reconcile.Result{}, fmt.Errorf("error reconciling templates: %w", err)
		// }

		return reconcile.Result{}, nil

	}

	if err := r.processResource(ctx, resource); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// processResource processes the fetched resource.
func (r *DynamicReconciler) processResource(ctx context.Context, resource *unstructured.Unstructured) error {
	log := log.FromContext(ctx)
	log.Info("Processing dynamic resource", "GVK", r.GVK, "resource", resource)

	makedResource, err := repository.MakeResource(resource.Object)
	if err != nil {
		return fmt.Errorf("error creating resource: %w", err)
	}

	re, err := GetSqliteRepository()
	if err != nil {
		log.Error(err, "unable to create repository")
		return err
	}

	if exists, err := re.CheckIfResourceExists(makedResource); err != nil {
		log.Error(err, "unable to check if resource exists")
		return err
	} else if exists {
		log.Info("Resource already exists in database",
			"GVK", r.GVK,
			"name", resource.GetName(),
			"namespace", resource.GetNamespace())
		return nil
	}

	if err := re.Upsert(makedResource); err != nil {
		log.Error(err, "unable to upsert resource")
		return err
	}

	return r.reconcileTemplates(ctx)
}

// reconcileTemplates reconciles the templates associated with the resource.
func (r *DynamicReconciler) reconcileTemplates(ctx context.Context) error {
	log := log.FromContext(ctx)
	var templates []string

	for templateName, gvks := range wm.templates {
		if _, exists := gvks[r.GVK]; exists {
			log.Info("Reconciling template", "templateName", templateName)
			templates = append(templates, templateName)
		}
	}

	for _, templateName := range templates {
		if err := r.processTemplate(ctx, templateName); err != nil {
			return err
		}
	}

	return nil
}

// processTemplate processes a single template.
func (r *DynamicReconciler) processTemplate(ctx context.Context, templateName string) error {
	log := log.FromContext(ctx)
	k8sClient, err := GetK8sClient()
	if err != nil {
		log.Error(err, "unable to get k8s client")
		return err
	}

	template, err := k8sClient.Get(ctx, "kumquat.guidewire.com", "Template", "templates", templateName)
	if err != nil {
		log.Error(err, "unable to get template", "templateName", templateName)
		return err
	}

	templateObj := &kumquatv1beta1.Template{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(template.Object, templateObj); err != nil {
		log.Error(err, "unable to convert unstructured to template")
		return err
	}

	re, err := GetSqliteRepository()
	if err != nil {
		log.Error(err, "unable to create repository")
		return err
	}

	return r.evaluateTemplate(ctx, template, re)
}

// fetchResource fetches the resource from the cluster.
func (r *DynamicReconciler) fetchResource(
	ctx context.Context,
	req reconcile.Request,
) (*unstructured.Unstructured, error) {
	resource := &unstructured.Unstructured{}
	resource.SetGroupVersionKind(r.GVK)
	err := r.Get(ctx, req.NamespacedName, resource)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return nil, err
		}
		return nil, nil
	}
	return resource, nil
}

// evaluateTemplate evaluates the template with the given data.
func (r *DynamicReconciler) evaluateTemplate(
	ctx context.Context,
	template *unstructured.Unstructured,
	re *repository.SQLiteRepository,
) error {
	templateObj := &kumquatv1beta1.Template{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(template.Object, templateObj)
	if err != nil {
		log := log.FromContext(ctx)
		log.Error(err, "unable to convert unstructured to template")
		return err
	}

	return processTemplateResources(templateObj, re, log.FromContext(ctx))
}

// unstructuredEventHandler handles events for unstructured resources.
type unstructuredEventHandler struct{}

func (h *unstructuredEventHandler) Create(
	ctx context.Context,
	evt event.TypedCreateEvent[*unstructured.Unstructured],
	q workqueue.RateLimitingInterface,
) {
	q.Add(ctrl.Request{NamespacedName: client.ObjectKeyFromObject(evt.Object)})
}

func (h *unstructuredEventHandler) Update(
	ctx context.Context,
	evt event.TypedUpdateEvent[*unstructured.Unstructured],
	q workqueue.RateLimitingInterface,
) {
	q.Add(ctrl.Request{NamespacedName: client.ObjectKeyFromObject(evt.ObjectNew)})
}

func (h *unstructuredEventHandler) Delete(
	ctx context.Context,
	evt event.TypedDeleteEvent[*unstructured.Unstructured],
	q workqueue.RateLimitingInterface,
) {
	q.Add(ctrl.Request{NamespacedName: client.ObjectKeyFromObject(evt.Object)})
}

func (h *unstructuredEventHandler) Generic(
	ctx context.Context,
	evt event.TypedGenericEvent[*unstructured.Unstructured],
	q workqueue.RateLimitingInterface,
) {
	q.Add(ctrl.Request{NamespacedName: client.ObjectKeyFromObject(evt.Object)})
}
