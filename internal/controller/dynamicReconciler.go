package controller

import (
	"context"
	kumquatv1beta1 "kumquat/api/v1beta1"
	"kumquat/repository"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DynamicReconciler struct {
	client.Client
	GVK          schema.GroupVersionKind
	K8sClient    K8sClient
	WatchManager *WatchManager
	repository   repository.Repository
}

func NewDynamicReconciler(client client.Client, gvk schema.GroupVersionKind, k8sClient K8sClient, wm *WatchManager, repo repository.Repository) *DynamicReconciler {
	return &DynamicReconciler{
		Client:       client,
		GVK:          gvk,
		K8sClient:    k8sClient,
		WatchManager: wm,
		repository:   repo,
	}
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
		// set group to core if it is empty
		group := r.GVK.Group
		if r.GVK.Group == "" {
			group = "core"
		}

		err = DeleteResourceFromDatabaseByNameAndNameSpace(r.repository, r.GVK.Kind, group, req.Namespace, req.Name) // nolint:errcheck
		if err != nil {
			return reconcile.Result{}, err
		}
		r.findAndReProcessAffectedTemplates(ctx) // nolint:errcheck

		return reconcile.Result{}, nil

	}

	err = UpsertResourceToDatabase(r.repository, resource, ctx) // nolint:errcheck
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.findAndReProcessAffectedTemplates(ctx) // nolint:errcheck

	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// reconcileTemplates reconciles the templates associated with the resource.
func (r *DynamicReconciler) findAndReProcessAffectedTemplates(ctx context.Context) error {
	log := log.FromContext(ctx)
	var templates []string
	wm := r.WatchManager // Use the injected WatchManager

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

// fetchResource fetches the resource from the cluster.
func (r *DynamicReconciler) fetchResource(
	ctx context.Context,
	req reconcile.Request,
) (*unstructured.Unstructured, error) {
	resource := &unstructured.Unstructured{}
	resource.SetGroupVersionKind(r.GVK)
	err := r.Client.Get(ctx, req.NamespacedName, resource)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return nil, err
		}
		return nil, nil
	}
	return resource, nil
}

// processTemplate processes a single template.
func (r *DynamicReconciler) processTemplate(ctx context.Context, templateName string) error {
	log := log.FromContext(ctx)

	template, err := r.K8sClient.Get(ctx, "kumquat.guidewire.com", "Template", "templates", templateName)
	if err != nil {
		log.Error(err, "unable to get template", "templateName", templateName)
		return err
	}

	templateObj := &kumquatv1beta1.Template{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(template.Object, templateObj); err != nil {
		log.Error(err, "unable to convert unstructured to template")
		return err
	}
	return ProcessTemplateResources(templateObj, r.repository, log, r.K8sClient, r.WatchManager)

}
