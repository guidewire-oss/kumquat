package cue_test

import (
	"kumquat/renderer"
	"kumquat/renderer/cue"
	"kumquat/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCUERenderer(t *testing.T) {
	tpl, err := cue.NewCUERenderer("test", t.Name())
	require.NoError(t, err)
	require.NotNil(t, tpl)
}

func TestRenderCUEWithEmptyResources(t *testing.T) {
	tpl, err := cue.NewCUERenderer(`
	out: [
        for result in data {
			foo: "bar"
		}
	]`, t.Name())
	require.NoError(t, err)
	require.NotNil(t, tpl)

	results := []map[string]repository.Resource{}
	output, err := renderer.Render(tpl, results, false)
	require.NoError(t, err)
	assert.Equal(t, 0, output.ResourceCount())
}

func TestRenderCUEImportStrings(t *testing.T) {
	tpl, err := cue.NewCUERenderer(`
	import "strings"

	out: {
		foo: strings.Join(["a", "b"], "-")
	}`, t.Name())
	require.NoError(t, err)
	require.NotNil(t, tpl)

	results := []map[string]repository.Resource{}
	output, err := renderer.Render(tpl, results, true)
	require.NoError(t, err)
	assert.Equal(t, 1, output.ResourceCount())
	result, err := output.ResultString(0)
	require.NoError(t, err)
	assert.Equal(t, "foo: a-b\n", result)
}

func TestRenderCUEWithBatchModeProcessingOffWithSingleResourceOutput(t *testing.T) {
	tpl, err := cue.NewCUERenderer(`
	out: 
		foo: data.a.metadata.name
	`, t.Name())
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
	assert.Equal(t, 2, output.ResourceCount())

	out1, err := output.ResultString(0)
	assert.NoError(t, err)
	assert.YAMLEq(t, "foo: test1", out1)

	out2, err := output.ResultString(1)
	assert.NoError(t, err)
	assert.YAMLEq(t, "foo: test2", out2)
}

func TestRenderCUEWithBatchModeProcessingOffWithMultipleResourcesOutput(t *testing.T) {
	tpl, err := cue.NewCUERenderer(`
	out: [
		{
			foo: data.a.metadata.name
		},
		{
			bar: data.a.metadata.name
		}
	]
	`, t.Name())
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
	assert.Equal(t, 2, output.ResourceCount())

	out1, err := output.ResultString(0)
	assert.NoError(t, err)
	assert.YAMLEq(t, "foo: test1\n---\nbar: test1", out1)

	out2, err := output.ResultString(1)
	assert.NoError(t, err)
	assert.YAMLEq(t, "foo: test2\n---\nbar: test2", out2)
}

func TestRenderCUEWithBatchModeProcessingOnWithMultipleResourcesOutput(t *testing.T) {
	tpl, err := cue.NewCUERenderer(`
	out: [
		for result in data {
		foo: result.a.metadata.name
		}
	]`, t.Name())
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
	assert.Equal(t, 1, output.ResourceCount())

	out, err := output.ResultString(0)
	assert.NoError(t, err)
	assert.YAMLEq(t, "foo: test1\n---\nfoo: test2", out)
}

func TestRenderCUEWithBatchModeProcessingOnWithSingleResourceOutput(t *testing.T) {
	tpl, err := cue.NewCUERenderer(`
	out: {
		foos: [
			for result in data {
				foo: result.a.metadata.name
			}
		]
	}`, t.Name())
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
	assert.Equal(t, 1, output.ResourceCount())

	out, err := output.ResultString(0)
	assert.NoError(t, err)
	assert.YAMLEq(t, "foos:\n- foo: test1\n- foo: test2", out)
}

func TestErrorUnknownVariableReference(t *testing.T) {
	tpl, err := cue.NewCUERenderer(`
	out: {
		foo: errrrrr
	}`, t.Name())
	require.NoError(t, err)
	require.NotNil(t, tpl)

	t.Run("BatchModeOn", func(t *testing.T) {
		results := []map[string]repository.Resource{}
		_, err = renderer.Render(tpl, results, true)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error evaluating CUE template")
		assert.ErrorContains(t, err, "[line 3, column 8]")
		assert.ErrorContains(t, err, "out.foo: reference \"errrrrr\" not found")

		var rerr *renderer.Error
		require.ErrorAs(t, err, &rerr)
		assert.EqualValues(t, 3, rerr.Line())
		assert.EqualValues(t, 8, rerr.Column())
		assert.ErrorContains(t, rerr, "out.foo: reference \"errrrrr\" not found")
	})

	t.Run("BatchModeOff", func(t *testing.T) {
		results := []map[string]repository.Resource{
			{
				"x": repository.Resource{},
			},
		}
		_, err = renderer.Render(tpl, results, false)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error evaluating CUE template")
		assert.ErrorContains(t, err, "[line 3, column 8]")
		assert.ErrorContains(t, err, "out.foo: reference \"errrrrr\" not found")

		var rerr *renderer.Error
		require.ErrorAs(t, err, &rerr)
		assert.EqualValues(t, 3, rerr.Line())
		assert.EqualValues(t, 8, rerr.Column())
		assert.ErrorContains(t, rerr, "out.foo: reference \"errrrrr\" not found")
	})
}

func TestErrorBadlyFormed(t *testing.T) {
	tpl, err := cue.NewCUERenderer(`out: {`, t.Name())
	require.NoError(t, err)
	require.NotNil(t, tpl)

	results := []map[string]repository.Resource{}
	_, err = renderer.Render(tpl, results, true)
	assert.Error(t, err)

	var rerr *renderer.Error
	require.ErrorAs(t, err, &rerr)
	assert.EqualValues(t, 3, rerr.Line())
	assert.EqualValues(t, 10, rerr.Column())
	assert.ErrorContains(t, err, "expected '}'")
}

func TestErrorNoOutput(t *testing.T) {
	tpl, err := cue.NewCUERenderer(`x: {
	foo: "bar"
}`, t.Name())
	require.NoError(t, err)
	require.NotNil(t, tpl)

	results := []map[string]repository.Resource{
		{
			"x": repository.Resource{},
		},
	}

	t.Run("BatchModeOn", func(t *testing.T) {
		_, err = renderer.Render(tpl, results, true)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'out' is not set to anything concrete")
	})

	t.Run("BatchModeOff", func(t *testing.T) {
		_, err = renderer.Render(tpl, results, false)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'out' is not set to anything concrete")
	})
}

func TestErrorUnsupportedOutputType(t *testing.T) {
	tpl, err := cue.NewCUERenderer(`out: 42`, t.Name())
	require.NoError(t, err)
	require.NotNil(t, tpl)

	results := []map[string]repository.Resource{
		{
			"x": repository.Resource{},
		},
	}

	t.Run("BatchModeOn", func(t *testing.T) {
		_, err = renderer.Render(tpl, results, true)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'out' has unsupported output type 'int'")
	})

	t.Run("BatchModeOff", func(t *testing.T) {
		_, err = renderer.Render(tpl, results, false)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "'out' has unsupported output type 'int'")
	})
}

func TestErrorMismatchedOutputType(t *testing.T) {
	renderToList, err := cue.NewCUERenderer(`out: [42]`, t.Name())
	require.NoError(t, err)
	require.NotNil(t, renderToList)

	results := []map[string]repository.Resource{
		{
			"x": repository.Resource{},
		},
	}

	t.Run("BatchModeOn", func(t *testing.T) {
		_, err := renderer.Render(renderToList, results, true)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error decoding 'out'")
	})

	t.Run("BatchModeOff", func(t *testing.T) {
		_, err := renderer.Render(renderToList, results, false)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error decoding 'out'")
	})
}
