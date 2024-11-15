package controller

import (
	"context"
	"fmt"
	"kumquat/repository"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func DeleteResourceFromDatabaseByNameAndNameSpace(
	repo repository.Repository, kind, group, namespace, name string) error {
	tableName := kind + "." + group

	err := repo.Delete(namespace, name, tableName)
	if err != nil {
		log.Log.Error(err, "unable to delete record")
		return err
	}
	log.Log.Info("Record deleted", "table", tableName, "namespace", namespace, "name", name)
	return nil
}

func UpsertResourceToDatabase(
	repo repository.Repository, resource *unstructured.Unstructured, ctx context.Context) error {
	log := log.FromContext(ctx)
	fmt.Println("UpsertResourceToDatabase")

	makedResource, err := repository.MakeResource(resource.Object)
	if err != nil {
		return fmt.Errorf("error creating resource: %w", err)
	}

	if err := repo.Upsert(makedResource); err != nil {
		log.Error(err, "unable to upsert resource")
		return err
	}
	return nil
}
