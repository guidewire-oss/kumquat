package template_test

import (
	"kumquat/repository"
	"kumquat/template"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestTemplate() map[string]any {
	return map[string]any{
		"apiVersion": "kumquat.guidewire.com/v1beta1",
		"kind":       "Template",
		"metadata": map[string]any{
			"name":      "TestTemplate",
			"namespace": "default",
		},
		"spec": map[string]any{
			"query": "SELECT * FROM test",
			"template": map[string]any{
				"data":     "content",
				"language": "gotemplate",
				"fileName": "test.out",
			},
		},
	}
}

func TestNewTemplate(t *testing.T) {
	r, err := repository.MakeResource(getTestTemplate())
	require.NoError(t, err)
	require.NotNil(t, r)

	tmpl, err := template.NewTemplate(r)
	require.NoError(t, err)
	require.NotNil(t, tmpl)
	assert.False(t, tmpl.BatchMode())
}

func TestNewTemplateWithBatchModeTrue(t *testing.T) {
	src := getTestTemplate()
	src["spec"].(map[string]any)["template"].(map[string]any)["batchModeProcessing"] = true

	r, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, r)

	tmpl, err := template.NewTemplate(r)
	require.NoError(t, err)
	require.NotNil(t, tmpl)
	assert.True(t, tmpl.BatchMode())
}

func TestNewTemplateBadAPIGroup(t *testing.T) {
	src := getTestTemplate()
	src["apiVersion"] = "guidewire.com/v1beta1"

	r, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, r)

	tmpl, err := template.NewTemplate(r)
	assert.Nil(t, tmpl)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "'TestTemplate'.apiVersion")
	assert.ErrorContains(t, err, "'guidewire.com' should be '"+template.TemplateAPIGroup+"'")
}

func TestNewTemplateBadAPIVersion(t *testing.T) {
	src := getTestTemplate()
	src["apiVersion"] = "kumquat.guidewire.com/v999"

	r, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, r)

	tmpl, err := template.NewTemplate(r)
	assert.Nil(t, tmpl)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "'TestTemplate'.apiVersion")
	assert.ErrorContains(t, err, "'v999' unsupported; supported versions are: 'v1beta1'")
}

func TestNewTemplateBadKind(t *testing.T) {
	src := getTestTemplate()
	src["kind"] = "Kumquat"

	r, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, r)

	tmpl, err := template.NewTemplate(r)
	assert.Nil(t, tmpl)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "'TestTemplate'.kind")
	assert.ErrorContains(t, err, "'Kumquat' should be '"+template.TemplateKind+"'")
}

func TestNewTemplateBadSpec(t *testing.T) {
	src := getTestTemplate()
	src["spec"] = 123

	r, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, r)

	tmpl, err := template.NewTemplate(r)
	assert.Nil(t, tmpl)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "'TestTemplate'.spec")
	assert.ErrorContains(t, err, "missing or not a map")
}

func TestNewTemplateBadQuery(t *testing.T) {
	src := getTestTemplate()
	src["spec"].(map[string]any)["query"] = 123

	r, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, r)

	tmpl, err := template.NewTemplate(r)
	assert.Nil(t, tmpl)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "'TestTemplate'.spec.query")
	assert.ErrorContains(t, err, "missing or not a string")
}

func TestNewTemplateBadTemplate(t *testing.T) {
	src := getTestTemplate()
	src["spec"].(map[string]any)["template"] = 123

	r, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, r)

	tmpl, err := template.NewTemplate(r)
	assert.Nil(t, tmpl)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "'TestTemplate'.spec.template")
	assert.ErrorContains(t, err, "missing or not a map")
}

// nolint:dupl
func TestNewTemplateBadFileName(t *testing.T) {
	src := getTestTemplate()

	t.Run("WrongTypeForFileName", func(t *testing.T) {
		src["spec"].(map[string]any)["template"].(map[string]any)["fileName"] = 123

		r, err := repository.MakeResource(src)
		require.NoError(t, err)
		require.NotNil(t, r)

		tmpl, err := template.NewTemplate(r)
		assert.Nil(t, tmpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'TestTemplate'.spec.template.fileName")
		assert.ErrorContains(t, err, "missing or not a string")
	})

	t.Run("BadValueForFileName", func(t *testing.T) {
		src["spec"].(map[string]any)["template"].(map[string]any)["fileName"] = "{{"

		r, err := repository.MakeResource(src)
		require.NoError(t, err)
		require.NotNil(t, r)

		tmpl, err := template.NewTemplate(r)
		assert.Nil(t, tmpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'TestTemplate'.spec.template.fileName")
		assert.ErrorContains(t, err, "error parsing Go template")
	})
}

// nolint:dupl
func TestNewTemplateBadLanguage(t *testing.T) {
	src := getTestTemplate()

	t.Run("WrongTypeForLanguage", func(t *testing.T) {
		src["spec"].(map[string]any)["template"].(map[string]any)["language"] = 123

		r, err := repository.MakeResource(src)
		require.NoError(t, err)
		require.NotNil(t, r)

		tmpl, err := template.NewTemplate(r)
		assert.Nil(t, tmpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'TestTemplate'.spec.template.language")
		assert.ErrorContains(t, err, "missing or not a string")
	})

	t.Run("BadValueForLanguage", func(t *testing.T) {
		src["spec"].(map[string]any)["template"].(map[string]any)["language"] = "english"

		r, err := repository.MakeResource(src)
		require.NoError(t, err)
		require.NotNil(t, r)

		tmpl, err := template.NewTemplate(r)
		assert.Nil(t, tmpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'TestTemplate'.spec.template.language")
		assert.ErrorContains(t, err, "unknown renderer 'english'")
	})
}

// nolint:dupl
func TestNewTemplateBadData(t *testing.T) {
	src := getTestTemplate()

	t.Run("WrongTypeForTemplate", func(t *testing.T) {
		src["spec"].(map[string]any)["template"].(map[string]any)["data"] = 123

		r, err := repository.MakeResource(src)
		require.NoError(t, err)
		require.NotNil(t, r)

		tmpl, err := template.NewTemplate(r)
		assert.Nil(t, tmpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'TestTemplate'.spec.template.data")
		assert.ErrorContains(t, err, "missing or not a string")
	})

	t.Run("BadValueForTemplate", func(t *testing.T) {
		src["spec"].(map[string]any)["template"].(map[string]any)["data"] = "{{"

		r, err := repository.MakeResource(src)
		require.NoError(t, err)
		require.NotNil(t, r)

		tmpl, err := template.NewTemplate(r)
		assert.Nil(t, tmpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'TestTemplate'.spec.template.data")
		assert.ErrorContains(t, err, "error parsing Go template")
	})
}

func TestMultipleValidationErrors(t *testing.T) {
	t.Run("KubernetesMetadataErrors", func(t *testing.T) {
		src := getTestTemplate()
		src["apiVersion"] = "guidewire.com/v999"
		src["kind"] = "Kumquat"

		r, err := repository.MakeResource(src)
		require.NoError(t, err)
		require.NotNil(t, r)

		tmpl, err := template.NewTemplate(r)
		assert.Nil(t, tmpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'TestTemplate'.apiVersion")
		assert.ErrorContains(t, err, "'guidewire.com' should be '"+template.TemplateAPIGroup+"'")
		assert.ErrorContains(t, err, "'v999' unsupported; supported versions are: 'v1beta1'")
		assert.ErrorContains(t, err, "'TestTemplate'.kind")
		assert.ErrorContains(t, err, "'Kumquat' should be '"+template.TemplateKind+"'")
		verr := err.(*template.ValidationErrors)
		assert.Equal(t, 3, len(verr.Unwrap()))
		assert.Equal(t, "TestTemplate", verr.Template())
	})

	t.Run("QueryAndTemplateErrors", func(t *testing.T) {
		src := getTestTemplate()
		src["spec"].(map[string]any)["query"] = 123
		src["spec"].(map[string]any)["template"] = nil

		r, err := repository.MakeResource(src)
		require.NoError(t, err)
		require.NotNil(t, r)

		tmpl, err := template.NewTemplate(r)
		assert.Nil(t, tmpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'TestTemplate'.spec.query")
		assert.ErrorContains(t, err, "missing or not a string")
		assert.ErrorContains(t, err, "'TestTemplate'.spec.template")
		assert.ErrorContains(t, err, "missing or not a map")
		verr := err.(*template.ValidationErrors)
		assert.Equal(t, 2, len(verr.Unwrap()))
		assert.Equal(t, "TestTemplate", verr.Template())
	})

	t.Run("QueryMissingAndTemplateErrors", func(t *testing.T) {
		src := getTestTemplate()
		delete(src["spec"].(map[string]any), "query")
		src["spec"].(map[string]any)["template"].(map[string]any)["fileName"] = "{{"
		src["spec"].(map[string]any)["template"].(map[string]any)["language"] = "english"
		src["spec"].(map[string]any)["template"].(map[string]any)["data"] = 123

		r, err := repository.MakeResource(src)
		require.NoError(t, err)
		require.NotNil(t, r)

		tmpl, err := template.NewTemplate(r)
		assert.Nil(t, tmpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'TestTemplate'.spec.query")
		assert.ErrorContains(t, err, "missing or not a string")
		assert.ErrorContains(t, err, "'TestTemplate'.spec.template.fileName")
		assert.ErrorContains(t, err, "error parsing Go template")
		assert.ErrorContains(t, err, "'TestTemplate'.spec.template.language")
		assert.ErrorContains(t, err, "unknown renderer 'english'")
		assert.ErrorContains(t, err, "'TestTemplate'.spec.template.data")
		assert.ErrorContains(t, err, "missing or not a string")
		verr := err.(*template.ValidationErrors)
		assert.Equal(t, 4, len(verr.Unwrap()))
		assert.Equal(t, "TestTemplate", verr.Template())
	})
}

type MockRepository struct {
	queryResult repository.ResultSet
	queryError  error
}

func (m *MockRepository) Query(query string) (repository.ResultSet, error) {
	return m.queryResult, m.queryError
}

func (m *MockRepository) Close() error {
	return nil
}

func (m *MockRepository) Upsert(r repository.Resource) error {
	return nil
}

func (m *MockRepository) Delete(namespace, name, table string) error {
	return nil
}
func (m *MockRepository) DropTable(tableName string) error {
	return nil
}
func (m *MockRepository) ExtractTableNamesFromQuery(query string) []string {
	return nil
}

func explodeOnErr[X any](x X, err error) X {
	if err != nil {
		panic(err)
	}

	return x
}

func TestEvaluateNonBatch(t *testing.T) {
	mockRepo := new(MockRepository)

	src := getTestTemplate()
	src["spec"].(map[string]any)["template"].(map[string]any)["fileName"] = "{{ .test.metadata.name }}.out"
	res, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, res)

	tmpl, err := template.NewTemplate(res)
	require.NoError(t, err)
	require.NotNil(t, tmpl)

	t.Run("QueryError", func(t *testing.T) {
		mockRepo.queryResult = repository.ResultSet{}
		mockRepo.queryError = assert.AnError
		output, err := tmpl.Evaluate(mockRepo)

		assert.Nil(t, output)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "query failed in Template 'TestTemplate'")
	})

	t.Run("QuerySuccess", func(t *testing.T) {
		mockRepo.queryError = nil
		mockRepo.queryResult = repository.ResultSet{
			Names: []string{"test"},
			Results: []map[string]repository.Resource{
				{"test": explodeOnErr(repository.MakeResource(map[string]any{
					"apiVersion": "v1",
					"kind":       "Whatever",
					"metadata": map[string]any{
						"name":      "abc",
						"namespace": "default",
					},
				}))},
				{"test": explodeOnErr(repository.MakeResource(map[string]any{
					"apiVersion": "v1",
					"kind":       "Whatever",
					"metadata": map[string]any{
						"name":      "def",
						"namespace": "default",
					},
				}))},
			},
		}

		output, err := tmpl.Evaluate(mockRepo)
		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, 2, output.Output.ResourceCount())

		s, err := output.Output.ResultString(0)
		require.NoError(t, err)
		assert.Equal(t, "content", s)
		assert.Equal(t, "abc.out", output.FileNames[0])

		s, err = output.Output.ResultString(1)
		require.NoError(t, err)
		assert.Equal(t, "content", s)
		assert.Equal(t, "def.out", output.FileNames[1])
	})
}

func TestEvaluateWithBatch(t *testing.T) {
	mockRepo := new(MockRepository)

	src := getTestTemplate()
	src["spec"].(map[string]any)["template"].(map[string]any)["batchModeProcessing"] = true
	src["spec"].(map[string]any)["template"].(map[string]any)["fileName"] = `{{ index . 0 "test" "metadata" "name" }}.out`
	res, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, res)

	tmpl, err := template.NewTemplate(res)
	require.NoError(t, err)
	require.NotNil(t, tmpl)

	t.Run("QueryError", func(t *testing.T) {
		mockRepo.queryResult = repository.ResultSet{}
		mockRepo.queryError = assert.AnError
		output, err := tmpl.Evaluate(mockRepo)

		assert.Nil(t, output)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "query failed in Template 'TestTemplate'")
	})

	t.Run("QuerySuccess", func(t *testing.T) {
		mockRepo.queryError = nil
		mockRepo.queryResult = repository.ResultSet{
			Names: []string{"test"},
			Results: []map[string]repository.Resource{
				{"test": explodeOnErr(repository.MakeResource(map[string]any{
					"apiVersion": "v1",
					"kind":       "Whatever",
					"metadata": map[string]any{
						"name":      "abc",
						"namespace": "default",
					},
				}))},
			},
		}

		output, err := tmpl.Evaluate(mockRepo)
		require.NoError(t, err)
		require.NotNil(t, output)

		assert.Equal(t, 1, output.Output.ResourceCount())
		s, err := output.Output.ResultString(0)
		require.NoError(t, err)
		assert.Equal(t, "content", s)

		assert.Len(t, output.FileNames, 1)
		assert.Equal(t, "abc.out", output.FileNames[0])
	})
}

func TestEvaluateWithRenderError(t *testing.T) {
	src := getTestTemplate()
	src["spec"].(map[string]any)["template"].(map[string]any)["data"] = "name is: {{ .test.metadata.name }}"

	res, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, res)

	tmpl, err := template.NewTemplate(res)
	require.NoError(t, err)
	require.NotNil(t, tmpl)

	mockRepo := new(MockRepository)
	mockRepo.queryError = nil
	mockRepo.queryResult = repository.ResultSet{
		Names: []string{"test"},
		Results: []map[string]repository.Resource{
			{"test": {}},
		},
	}

	output, err := tmpl.Evaluate(mockRepo)
	assert.Nil(t, output)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error rendering Template 'TestTemplate'")
}

func TestEvaluateWithFilenameRenderError(t *testing.T) {
	src := getTestTemplate()
	src["spec"].(map[string]any)["template"].(map[string]any)["fileName"] = "{{ .test.metadata.name }}.yaml"

	res, err := repository.MakeResource(src)
	require.NoError(t, err)
	require.NotNil(t, res)

	tmpl, err := template.NewTemplate(res)
	require.NoError(t, err)
	require.NotNil(t, tmpl)

	mockRepo := new(MockRepository)
	mockRepo.queryError = nil
	mockRepo.queryResult = repository.ResultSet{
		Names: []string{"test"},
		Results: []map[string]repository.Resource{
			{"test": {}},
		},
	}

	output, err := tmpl.Evaluate(mockRepo)
	assert.Nil(t, output)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "error rendering fileName in Template 'TestTemplate'")
}
