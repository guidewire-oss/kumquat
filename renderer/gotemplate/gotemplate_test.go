package gotemplate_test

import (
	"kumquat/renderer"
	"kumquat/renderer/gotemplate"
	"kumquat/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGoRenderer(t *testing.T) {
	tpl, err := gotemplate.NewGoRenderer("test", t.Name())
	assert.NoError(t, err)

	r, err := renderer.Render(tpl, []map[string]repository.Resource{}, false)
	assert.NoError(t, err)
	assert.Zero(t, r.ResourceCount())

	r, err = renderer.Render(tpl, []map[string]repository.Resource{}, true)
	assert.NoError(t, err)
	assert.Equal(t, 1, r.ResourceCount())
	s, err := r.ResultString(0)
	assert.NoError(t, err)
	assert.Equal(t, "test", s)
}

func TestNewGoRendererWithEmptyTemplate(t *testing.T) {
	tpl, err := gotemplate.NewGoRenderer("", t.Name())
	require.NotNil(t, tpl)
	require.NoError(t, err)

	r, err := renderer.Render(tpl, []map[string]repository.Resource{}, false)
	require.NoError(t, err)
	assert.Zero(t, r.ResourceCount())

	r, err = renderer.Render(tpl, []map[string]repository.Resource{}, true)
	require.NoError(t, err)
	assert.Equal(t, 1, r.ResourceCount())
	s, err := r.ResultString(0)
	require.NoError(t, err)
	assert.Empty(t, s)
}

func TestNewGoRendererWithInvalidTemplate(t *testing.T) {
	t.Run("UnclosedAction", func(t *testing.T) {
		tpl, err := gotemplate.NewGoRenderer(`test
	
		{{`, t.Name())
		assert.Nil(t, tpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error parsing Go template: [line 3, column ?] unclosed action")

		var rerr *renderer.Error
		require.ErrorAs(t, err, &rerr)
		assert.EqualValues(t, 3, rerr.Line())
		assert.ErrorContains(t, rerr, "unclosed action")
	})

	t.Run("UnclosedAction2", func(t *testing.T) {
		tpl, err := gotemplate.NewGoRenderer(`test
	
		{{
		42`, t.Name())
		assert.Nil(t, tpl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error parsing Go template: [line 4, column ?] unclosed action")

		var rerr *renderer.Error
		require.ErrorAs(t, err, &rerr)
		assert.EqualValues(t, 4, rerr.Line())
		assert.ErrorContains(t, rerr, "unclosed action started at line 3")
	})
}

func TestRenderGoWithEmptyResourcesWithBatchModeOff(t *testing.T) {
	tpl, err := gotemplate.NewGoRenderer("test", t.Name())
	require.NoError(t, err)
	require.NotNil(t, tpl)

	results := []map[string]repository.Resource{}
	output, err := renderer.Render(tpl, results, false)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 0, output.ResourceCount())

	_, err = output.ResultString(0)
	assert.ErrorContains(t, err, "resource index out of range")
}

func TestRenderGoWithEmptyResourcesWithBatchModeOn(t *testing.T) {
	tpl, err := gotemplate.NewGoRenderer("test", t.Name())
	require.NoError(t, err)
	require.NotNil(t, tpl)

	results := []map[string]repository.Resource{}
	output, err := renderer.Render(tpl, results, true)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 1, output.ResourceCount())

	s, err := output.ResultString(0)
	require.NoError(t, err)
	assert.Equal(t, "test", s)

	_, err = output.ResultString(1)
	assert.ErrorContains(t, err, "resource index out of range")
}

func TestRenderGoWithTwoInputsBatchModeProcessingFalse(t *testing.T) {

	tpl, err := gotemplate.NewGoRenderer(
		`{{.a.apiVersion}},
{{.a.metadata.name}},
{{.a.metadata.namespace}}`, t.Name())

	require.NoError(t, err)
	require.NotNil(t, tpl)
	a1, err := repository.MakeResource(map[string]interface{}{
		"apiVersion": "kumquat.guidewire.com/v1beta1",
		"kind":       "Test",
		"metadata": map[string]interface{}{
			"name":      "test1",
			"namespace": "default",
		},
	})
	require.NoError(t, err)
	a2, err := repository.MakeResource(map[string]interface{}{
		"apiVersion": "kumquat.guidewire.com/v1beta1",
		"kind":       "Test",
		"metadata": map[string]interface{}{
			"name":      "test2",
			"namespace": "default",
		},
	})
	require.NoError(t, err)
	results := []map[string]repository.Resource{
		{"a": a1},
		{"a": a2},
	}
	output, err := renderer.Render(tpl, results, false)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 2, output.ResourceCount())
	out1, err := output.ResultString(0)
	assert.NoError(t, err)
	assert.Equal(t, "kumquat.guidewire.com/v1beta1,\ntest1,\ndefault", out1)
	out2, err := output.ResultString(1)
	assert.NoError(t, err)
	assert.Equal(t, "kumquat.guidewire.com/v1beta1,\ntest2,\ndefault", out2)

}

func TestRenderGoWithTwoInputsBatchModeProcessingOn(t *testing.T) {

	tpl, err := gotemplate.NewGoRenderer(
		`{{- range .}}
{{.a.apiVersion}},
{{.a.metadata.name}},
{{.a.metadata.namespace}}
{{- end}}`, t.Name())

	require.NoError(t, err)
	require.NotNil(t, tpl)
	a1, err := repository.MakeResource(map[string]interface{}{
		"apiVersion": "kumquat.guidewire.com/v1beta1",
		"kind":       "Test",
		"metadata": map[string]interface{}{
			"name":      "test1",
			"namespace": "default",
		},
	})
	require.NoError(t, err)
	a2, err := repository.MakeResource(map[string]interface{}{
		"apiVersion": "kumquat.guidewire.com/v1beta1",
		"kind":       "Test",
		"metadata": map[string]interface{}{
			"name":      "test2",
			"namespace": "default",
		},
	})
	require.NoError(t, err)
	results := []map[string]repository.Resource{
		{"a": a1},
		{"a": a2},
	}
	output, err := renderer.Render(tpl, results, true)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, 1, output.ResourceCount())
	out, err := output.ResultString(0)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"\nkumquat.guidewire.com/v1beta1,\ntest1,\ndefault\nkumquat.guidewire.com/v1beta1,\ntest2,\ndefault",
		out)

}

func TestRenderErrors(t *testing.T) {
	tpl, err := gotemplate.NewGoRenderer(
		`{{.a.apiVersion}},
{{.x.metadata.name}},
{{.x.metadata.namespace}}`, t.Name())

	require.NoError(t, err)
	require.NotNil(t, tpl)
	a1, err := repository.MakeResource(map[string]interface{}{
		"apiVersion": "kumquat.guidewire.com/v1beta1",
		"kind":       "Test",
		"metadata": map[string]interface{}{
			"name":      "test1",
			"namespace": "default",
		},
	})
	require.NoError(t, err)
	results := []map[string]repository.Resource{
		{"a": a1},
	}

	t.Run("BatchModeOff", func(t *testing.T) {
		_, err = renderer.Render(tpl, results, false)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error executing Go template")
		assert.ErrorContains(t, err, "[line 2, column 4]")

		var rerr *renderer.Error
		require.ErrorAs(t, err, &rerr)
		assert.EqualValues(t, 2, rerr.Line())
		assert.EqualValues(t, 4, rerr.Column())
		assert.ErrorContains(t, rerr, "nil pointer") // because "x" is not a key in ".", so map returns nil
	})

	t.Run("BatchModeOn", func(t *testing.T) {
		_, err = renderer.Render(tpl, results, true)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error executing Go template")
		assert.ErrorContains(t, err, "[line 1, column 4]")

		var rerr *renderer.Error
		require.ErrorAs(t, err, &rerr)
		assert.EqualValues(t, 1, rerr.Line())
		assert.EqualValues(t, 4, rerr.Column())
		assert.ErrorContains(t, rerr, "can't evaluate field a") // because "." is a slice
	})
}
