package repository_test

import (
	"io/fs"
	"kumquat/repository"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeResource(t *testing.T) {
	r, err := repository.MakeResource(map[string]any{
		"apiVersion": "guidewire.com/v1beta1",
		"kind":       "Example",
		"metadata": map[string]any{
			"name":      "alpha",
			"namespace": "examples",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "Example", r.Kind())
	assert.Equal(t, "guidewire.com", r.Group())
	assert.Equal(t, "v1beta1", r.Version())
	assert.Equal(t, "examples", r.Namespace())
	assert.Equal(t, "alpha", r.Name())
}

func TestResourceMissingElements(t *testing.T) {
	t.Run("MissingAPIVersion", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"kind": "Example",
			"metadata": map[string]any{
				"name":      "alpha",
				"namespace": "examples",
			},
		})
		assert.ErrorContains(t, err, "missing apiVersion")
	})

	t.Run("MissingKind", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": "guidewire.com/v1beta1",
			"metadata": map[string]any{
				"name":      "alpha",
				"namespace": "examples",
			},
		})
		assert.ErrorContains(t, err, "missing kind")
	})

	t.Run("MissingMetadata", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": "guidewire.com/v1beta1",
			"kind":       "Example",
		})
		assert.ErrorContains(t, err, "missing metadata")
	})

	t.Run("MissingName", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": "guidewire.com/v1beta1",
			"kind":       "Example",
			"metadata": map[string]any{
				"namespace": "examples",
			},
		})
		assert.ErrorContains(t, err, "missing name")
	})
}

func TestResourceDataErrors(t *testing.T) {
	t.Run("APIVersionMustBeString", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": 123,
			"kind":       "Example",
			"metadata": map[string]any{
				"name":      "1.44.0",
				"namespace": "examples",
			},
		})
		assert.ErrorContains(t, err, "must be a string")
	})

	t.Run("KindMustBeString", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": "guidewire.com/v1beta1",
			"kind":       123,
			"metadata": map[string]any{
				"name":      "1.44.0",
				"namespace": "examples",
			},
		})
		assert.ErrorContains(t, err, "must be a string")
	})

	t.Run("MetadataMustBeMap", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": "guidewire.com/v1beta1",
			"kind":       "Example",
			"metadata":   "test",
		})
		assert.ErrorContains(t, err, "must be a map")
	})

	t.Run("NameMustBeString", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": "guidewire.com/v1beta1",
			"kind":       "Example",
			"metadata": map[string]any{
				"name":      123,
				"namespace": "examples",
			},
		})
		assert.ErrorContains(t, err, "must be a string")
	})

	t.Run("NamespaceMustBeString", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": "guidewire.com/v1beta1",
			"kind":       "Example",
			"metadata": map[string]any{
				"name":      "1.44.0",
				"namespace": 123,
			},
		})
		assert.ErrorContains(t, err, "must be a string")
	})

	t.Run("APIVersionMissingGroup", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": "/v1beta1",
			"kind":       "Example",
			"metadata": map[string]any{
				"name":      "1.44.0",
				"namespace": "examples",
			},
		})
		assert.ErrorContains(t, err, "missing group")
	})

	t.Run("APIVersionMissingVersion", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": "guidewire.com/",
			"kind":       "Example",
			"metadata": map[string]any{
				"name":      "1.44.0",
				"namespace": "examples",
			},
		})
		assert.ErrorContains(t, err, "missing version")
	})

	t.Run("APIVersionMissingSeparator", func(t *testing.T) {
		_, err := repository.MakeResource(map[string]any{
			"apiVersion": "guidewire.com",
			"kind":       "Example",
			"metadata": map[string]any{
				"name":      "1.44.0",
				"namespace": "examples",
			},
		})
		assert.ErrorContains(t, err, "missing '/' separator")
	})
}

func TestCoreAPIGroupWorks(t *testing.T) {
	_, err := repository.MakeResource(map[string]any{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]any{
			"name":      "test",
			"namespace": "default",
		},
	})
	assert.NoError(t, err)
}

type MockRepo struct {
	resources []repository.Resource
}

func (repo *MockRepo) Query(q string) (repository.ResultSet, error) {
	return repository.ResultSet{}, nil
}
func (repo *MockRepo) Close() error {
	return nil
}
func (repo *MockRepo) Upsert(r repository.Resource) error {
	repo.resources = append(repo.resources, r)

	return nil
}

func TestLoadYAMLResourcesFromDirectoryTreeWithErrors(t *testing.T) {
	repo := &MockRepo{}

	// Load a file with invalid YAML
	mockfs := fstest.MapFS{
		"invalid.yaml": {Data: []byte("hi")},
	}

	err := repository.LoadYAMLFromDirectoryTree(mockfs, ".", repo)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "yaml: unmarshal errors")
	assert.ErrorContains(t, err, "invalid.yaml")

	// Load a directory that doesn't exist
	err = repository.LoadYAMLFromDirectoryTree(mockfs, "unknown-directory", repo)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "file does not exist")
	assert.ErrorIs(t, err, fs.ErrNotExist)

	// Load a YAML file that isn't a Kubernetes resource
	mockfs = fstest.MapFS{
		"invalid.yaml": {Data: []byte("abc: 123\ndef: hello")},
	}
	err = repository.LoadYAMLFromDirectoryTree(mockfs, ".", repo)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "missing apiVersion")
	assert.ErrorContains(t, err, "invalid.yaml")
}

func TestLoadYAMLResourcesFromDirectoryTree(t *testing.T) {
	repo := &MockRepo{}

	mockfs := fstest.MapFS{
		"example1.yaml":   {Data: []byte("apiVersion: guidewire.com/v1beta1\nkind: Example\nmetadata:\n  name: 1.44.0\n  namespace: examples")},
		"example2.yaml":   {Data: []byte("apiVersion: guidewire.com/v1beta1\nkind: Example\nmetadata:\n  name: 1.45.0\n  namespace: examples")},
		"configmap1.yaml": {Data: []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\n  namespace: default")},
	}

	err := repository.LoadYAMLFromDirectoryTree(mockfs, ".", repo)
	assert.NoError(t, err)
	assert.Len(t, repo.resources, 3)
}
