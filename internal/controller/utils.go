package controller

import (
	"context"
	"fmt"
	"kumquat/repository"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func DeleteResourceFromDatabaseByNameAndNameSpace(kind, group, namespace, name string) error {
	tableName := kind + "." + group

	re, err := GetSqliteRepository()
	if err != nil {
		log.Log.Error(err, "unable to create repository")
		return err
	}
	err = re.Delete(namespace, name, tableName)
	if err != nil {
		log.Log.Error(err, "unable to delete record")
		return err
	}
	log.Log.Info("Record deleted", "table", tableName, "namespace", namespace, "name", name)
	return nil
}

func UpsertResourceToDatabase(resource *unstructured.Unstructured, ctx context.Context) error {
	log := log.FromContext(ctx)
	fmt.Println("UpsertResourceTsdsdvcoDatabase")
	// log.Info("Processing dynamic resource", "GVK", r.GVK, "resource", resource)

	makedResource, err := repository.MakeResource(resource.Object)
	if err != nil {
		return fmt.Errorf("error creating resource: %w", err)
	}

	sr, err := GetSqliteRepository()
	if err != nil {
		log.Error(err, "unable to create repository")
		return err
	}

	if err := sr.Upsert(makedResource); err != nil {
		log.Error(err, "unable to upsert resource")
		return err
	}
	return nil

}
