package controller_test

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"kumquat/repository"
// 	"testing"

// 	controller "kumquat/internal/controller"

// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// )

// // TestSuite holds the shared state for the tests.
// type TestSuite struct {
// 	mockRepo *MockRepository
// 	ctx      context.Context
// }

// // Setup initializes the shared state before each test.
// func (s *TestSuite) Setup() {
// 	s.ctx = context.Background()
// 	s.mockRepo = &MockRepository{}
// }

// // Teardown cleans up the shared state after each test.
// func (s *TestSuite) Teardown() {
// 	// If there are resources to clean up, do it here.
// 	// For this example, we don't have any specific teardown steps.
// }

// // MockRepository is a mock implementation of repository.Repository
// type MockRepository struct {
// 	DeleteFunc func(namespace, name, tableName string) error
// 	UpsertFunc func(resource repository.Resource) error
// }

// // Delete mocks the Delete method of the repository
// func (m *MockRepository) Delete(namespace, name, tableName string) error {
// 	if m.DeleteFunc != nil {
// 		return m.DeleteFunc(namespace, name, tableName)
// 	}
// 	return nil
// }

// // Upsert mocks the Upsert method of the repository
// func (m *MockRepository) Upsert(resource repository.Resource) error {
// 	if m.UpsertFunc != nil {
// 		return m.UpsertFunc(resource)
// 	}
// 	return nil
// }

// // Implement other methods if required
// func (m *MockRepository) Query(query string) (repository.ResultSet, error) {
// 	return repository.ResultSet{}, nil
// }

// func (m *MockRepository) Close() error {
// 	return nil
// }

// func (m *MockRepository) ExtractTableNamesFromQuery(query string) []string {
// 	return nil
// }

// func (m *MockRepository) DropTable(table string) error {
// 	return nil
// }

// const (
// 	kind = "TestKind"
// )

// // TestDeleteResourceFromDatabaseByNameAndNameSpace_Success tests the success case
// func TestDeleteResourceFromDatabaseByNameAndNameSpace_Success(t *testing.T) {
// 	// Arrange
// 	suite := &TestSuite{}
// 	suite.Setup()
// 	defer suite.Teardown()

// 	kind := kind
// 	group := "TestGroup"
// 	namespace := "TestNamespace"
// 	name := "TestName"
// 	tableName := kind + "." + group

// 	suite.mockRepo.DeleteFunc = func(ns, n, tbl string) error {
// 		if ns != namespace || n != name || tbl != tableName {
// 			return fmt.Errorf("unexpected parameters")
// 		}
// 		return nil
// 	}

// 	// Act
// 	err := controller.DeleteResourceFromDatabaseByNameAndNameSpace(suite.mockRepo, kind, group, namespace, name)

// 	// Assert
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// }

// // TestDeleteResourceFromDatabaseByNameAndNameSpace_DeleteError tests the error case
// func TestDeleteResourceFromDatabaseByNameAndNameSpace_DeleteError(t *testing.T) {
// 	// Arrange
// 	suite := &TestSuite{}
// 	suite.Setup()
// 	defer suite.Teardown()

// 	kind := kind
// 	group := "TestGroup"
// 	namespace := "TestNamespace"
// 	name := "TestName"
// 	tableName := kind + "." + group
// 	expectedError := errors.New("delete error")

// 	suite.mockRepo.DeleteFunc = func(ns, n, tbl string) error {
// 		if ns != namespace || n != name || tbl != tableName {
// 			return fmt.Errorf("unexpected parameters")
// 		}
// 		return expectedError
// 	}

// 	// Act
// 	err := controller.DeleteResourceFromDatabaseByNameAndNameSpace(suite.mockRepo, kind, group, namespace, name)

// 	// Assert
// 	if err != expectedError {
// 		t.Errorf("expected error %v, got %v", expectedError, err)
// 	}
// }

// // TestUpsertResourceToDatabase_Success tests the success case
// func TestUpsertResourceToDatabase_Success(t *testing.T) {
// 	// Arrange
// 	suite := &TestSuite{}
// 	suite.Setup()
// 	defer suite.Teardown()

// 	resource := &unstructured.Unstructured{
// 		Object: map[string]interface{}{
// 			"apiVersion": "v1",
// 			"kind":       "TestKind",
// 			"metadata": map[string]interface{}{
// 				"name":      "test-resource",
// 				"namespace": "test-namespace",
// 			},
// 		},
// 	}

// 	// In the success case, we provide a valid resource and expect no error
// 	suite.mockRepo.UpsertFunc = func(res repository.Resource) error {
// 		// Assert that the resource has the expected values
// 		if res.Group() != "core" || res.Version() != "v1" || res.Kind() != "TestKind" {
// 			return fmt.Errorf("unexpected resource details: %+v", res)
// 		}
// 		if res.Name() != "test-resource" || res.Namespace() != "test-namespace" {
// 			return fmt.Errorf("unexpected resource metadata: name=%s, namespace=%s", res.Name(), res.Namespace())
// 		}
// 		return nil
// 	}

// 	// Act
// 	err := controller.UpsertResourceToDatabase(suite.mockRepo, resource, suite.ctx)

// 	// Assert
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// }

// // TestUpsertResourceToDatabase_MakeResourceError tests error handling when MakeResource fails
// func TestUpsertResourceToDatabase_MakeResourceError(t *testing.T) {
// 	// Arrange
// 	suite := &TestSuite{}
// 	suite.Setup()
// 	defer suite.Teardown()

// 	// Provide an invalid resource object that will cause MakeResource to fail
// 	resource := &unstructured.Unstructured{
// 		Object: map[string]interface{}{
// 			"kind": "TestKind",
// 			// "apiVersion" is missing to induce an error
// 			"metadata": map[string]interface{}{
// 				"name":      "test-resource",
// 				"namespace": "test-namespace",
// 			},
// 		},
// 	}

// 	// Act
// 	err := controller.UpsertResourceToDatabase(suite.mockRepo, resource, suite.ctx)

// 	// Assert
// 	if err == nil || err.Error() != "error creating resource: missing apiVersion" {
// 		t.Errorf("expected error about missing apiVersion, got %v", err)
// 	}
// }

// // TestUpsertResourceToDatabase_UpsertError tests error handling when Upsert returns an error
// func TestUpsertResourceToDatabase_UpsertError(t *testing.T) {
// 	// Arrange
// 	suite := &TestSuite{}
// 	suite.Setup()
// 	defer suite.Teardown()

// 	resource := &unstructured.Unstructured{
// 		Object: map[string]interface{}{
// 			"apiVersion": "v1",
// 			"kind":       "TestKind",
// 			"metadata": map[string]interface{}{
// 				"name":      "test-resource",
// 				"namespace": "test-namespace",
// 			},
// 		},
// 	}

// 	expectedError := errors.New("upsert error")

// 	suite.mockRepo.UpsertFunc = func(res repository.Resource) error {
// 		return expectedError
// 	}

// 	// Act
// 	err := controller.UpsertResourceToDatabase(suite.mockRepo, resource, suite.ctx)

// 	// Assert
// 	if err != expectedError {
// 		t.Errorf("expected error %v, got %v", expectedError, err)
// 	}
// }
