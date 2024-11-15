package controller

import (
	"errors"
	"fmt"
	"kumquat/repository"
	"testing"
)

// MockRepository is a mock implementation of repository.Repository
type MockRepository struct {
	DeleteFunc func(namespace, name, tableName string) error
	UpsertFunc func(resource repository.Resource) error
}

// Delete mocks the Delete method of the repository
func (m *MockRepository) Delete(namespace, name, tableName string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(namespace, name, tableName)
	}
	return nil
}

// Upsert mocks the Upsert method of the repository
func (m *MockRepository) Upsert(resource repository.Resource) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(resource)
	}
	return nil
}

// Implement other methods of the Repository interface if needed
func (m *MockRepository) Query(query string) (repository.ResultSet, error) {
	return repository.ResultSet{}, nil
}

func (m *MockRepository) Close() error {
	return nil
}

func (m *MockRepository) ExtractTableNamesFromQuery(query string) []string {
	return nil
}

func (m *MockRepository) DropTable(table string) error {
	return nil
}

// Save the original MakeResource function to restore after tests
var originalMakeResource = repository.MakeResource

// restoreMakeResource restores the original MakeResource function

// TestDeleteResourceFromDatabaseByNameAndNameSpace_Success tests the success case
func TestDeleteResourceFromDatabaseByNameAndNameSpace_Success(t *testing.T) {
	// Arrange
	kind := "TestKind"
	group := "TestGroup"
	namespace := "TestNamespace"
	name := "TestName"
	tableName := kind + "." + group

	mockRepo := &MockRepository{
		DeleteFunc: func(ns, n, tbl string) error {
			if ns != namespace || n != name || tbl != tableName {
				return fmt.Errorf("unexpected parameters")
			}
			return nil
		},
	}

	// Act
	err := DeleteResourceFromDatabaseByNameAndNameSpace(mockRepo, kind, group, namespace, name)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestDeleteResourceFromDatabaseByNameAndNameSpace_DeleteError tests the error case
func TestDeleteResourceFromDatabaseByNameAndNameSpace_DeleteError(t *testing.T) {
	// Arrange
	kind := "TestKind"
	group := "TestGroup"
	namespace := "TestNamespace"
	name := "TestName"
	tableName := kind + "." + group
	expectedError := errors.New("delete error")

	mockRepo := &MockRepository{
		DeleteFunc: func(ns, n, tbl string) error {
			if ns != namespace || n != name || tbl != tableName {
				return fmt.Errorf("unexpected parameters")
			}
			return expectedError
		},
	}

	// Act
	err := DeleteResourceFromDatabaseByNameAndNameSpace(mockRepo, kind, group, namespace, name)

	// Assert
	if err != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

// TestUpsertResourceToDatabase_Success tests the success case
// func TestUpsertResourceToDatabase_Success(t *testing.T) {
// 	// Arrange
// 	defer restoreMakeResource()

// 	ctx := context.Background()
// 	resource := &unstructured.Unstructured{
// 		Object: map[string]interface{}{
// 			"metadata": map[string]interface{}{
// 				"name":      "test-resource",
// 				"namespace": "test-namespace",
// 			},
// 			"kind":       "TestKind",
// 			"apiVersion": "v1",
// 		},
// 	}

// 	makedResource := repository.Resource{}

// 	mockMakeResource(makedResource, nil)

// 	mockRepo := &MockRepository{
// 		UpsertFunc: func(res repository.Resource) error {
// 			if res != makedResource {
// 				return fmt.Errorf("unexpected resource")
// 			}
// 			return nil
// 		},
// 	}

// 	// Act
// 	err := UpsertResourceToDatabase(mockRepo, resource, ctx)

// 	// Assert
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// }

// TestUpsertResourceToDatabase_MakeResourceError tests the error case when MakeResource returns an error
// func TestUpsertResourceToDatabase_MakeResourceError(t *testing.T) {
// 	// Arrange
// 	defer restoreMakeResource()

// 	ctx := context.Background()
// 	resource := &unstructured.Unstructured{}
// 	expectedError := errors.New("make resource error")

// 	mockMakeResource(repository.Resource{}, expectedError)

// 	// Act
// 	err := UpsertResourceToDatabase(nil, resource, ctx)

// 	// Assert
// 	if err == nil || err.Error() != fmt.Sprintf("error creating resource: %v", expectedError) {
// 		t.Errorf("expected error %v, got %v", expectedError, err)
// 	}
// }

// // TestUpsertResourceToDatabase_UpsertError tests the error case when Upsert returns an error
// func TestUpsertResourceToDatabase_UpsertError(t *testing.T) {
// 	// Arrange
// 	defer restoreMakeResource()

// 	ctx := context.Background()
// 	resource := &unstructured.Unstructured{}

// 	makedResource := repository.Resource{}
// 	expectedError := errors.New("upsert error")

// 	mockMakeResource(makedResource, nil)

// 	mockRepo := &MockRepository{
// 		UpsertFunc: func(res repository.Resource) error {
// 			if res != makedResource {
// 				return fmt.Errorf("unexpected resource")
// 			}
// 			return expectedError
// 		},
// 	}

// 	// Act
// 	err := UpsertResourceToDatabase(mockRepo, resource, ctx)

// 	// Assert
// 	if err != expectedError {
// 		t.Errorf("expected error %v, got %v", expectedError, err)
// 	}
// }
