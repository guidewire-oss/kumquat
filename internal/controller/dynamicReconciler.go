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
	client.Client
	GVK       schema.GroupVersionKind
	K8sClient K8sClient
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
		// delete record from database
		// TODO: checking the return code and returning an error causes the E2E tests to fail during delete template

		DeleteRecord(r.GVK.Kind+"."+group, req.Namespace, req.Name) // nolint:errcheck

		r.reconcileTemplates(ctx) // nolint:errcheck

		return reconcile.Result{}, nil

	}

	if err := r.processResource(ctx, resource); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
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

	k8sClient := r.K8sClient
	if k8sClient == nil {
		err := fmt.Errorf("K8sClient is not initialized")
		log.Error(err, "K8sClient is not initialized")
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
	err := r.Client.Get(ctx, req.NamespacedName, resource)
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

	return processTemplateResources(templateObj, re, log.FromContext(ctx), r.K8sClient)
}
