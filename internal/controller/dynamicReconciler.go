package controller

import (
	"context"
	"fmt"
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
	client       client.Client
	gvk          schema.GroupVersionKind
	k8sClient    K8sClient
	watchManager WatchManager
	repository   repository.Repository
}

func NewDynamicReconciler(client client.Client, gvk schema.GroupVersionKind, k8sClient K8sClient, wm WatchManager, repo repository.Repository) *DynamicReconciler {
	return &DynamicReconciler{
		client:       client,
		gvk:          gvk,
		k8sClient:    k8sClient,
		watchManager: wm,
		repository:   repo,
	}
}

func (r *DynamicReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling dynamic resource", "GVK", r.gvk, "name", req.Name, "namespace", req.Namespace)

	resource, err := r.fetchResource(ctx, req)
	if err != nil {
		return reconcile.Result{}, err
	}

	if resource == nil {
		log.Info("Resource deleted", "GVK", r.gvk, "name", req.Name, "namespace", req.Namespace)
		// set group to core if it is empty
		group := r.gvk.Group
		if r.gvk.Group == "" {
			group = "core"
		}

		err = DeleteResourceFromDatabaseByNameAndNameSpace(r.repository, r.gvk.Kind, group, req.Namespace, req.Name)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = r.findAndReProcessAffectedTemplates(ctx)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil

	}

	err = UpsertResourceToDatabase(r.repository, resource, ctx)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.findAndReProcessAffectedTemplates(ctx)

	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// reconcileTemplates reconciles the templates associated with the resource.
func (r *DynamicReconciler) findAndReProcessAffectedTemplates(ctx context.Context) error {
	log := log.FromContext(ctx)
	var templates []string

	templatesMap := r.watchManager.GetManagedTemplates()
	for templateName, gvks := range templatesMap {
		if _, exists := gvks[r.gvk]; exists {
			log.Info("Reconciling template", "templateName", templateName)
			templates = append(templates, templateName)
		}
	}

	for _, templateName := range templates {
		if err := r.processTemplate(ctx, templateName); err != nil {
			if err.Error() == "not found" {
				log.Info("Resource not found", "templateName", templateName)
			} else {
				fmt.Println("Error in processTemplate", err)
				return err

			}
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
	resource.SetGroupVersionKind(r.gvk)
	err := r.client.Get(ctx, req.NamespacedName, resource)
	if err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return resource, nil
}

// processTemplate processes a single template.
func (r *DynamicReconciler) processTemplate(ctx context.Context, templateName string) error {
	log := log.FromContext(ctx)

	template, err := r.k8sClient.Get(ctx, "kumquat.guidewire.com", "Template", "templates", templateName)
	if err != nil {
		log.Error(err, "unable to get template", "templateName", templateName)
		return err
	}

	templateObj := &kumquatv1beta1.Template{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(template.Object, templateObj); err != nil {
		log.Error(err, "unable to convert unstructured to template")
		return err
	}
	return ProcessTemplateResources(templateObj, r.repository, log, r.k8sClient, r.watchManager)

}
