package renderer_test

import (
	"errors"
	"fmt"
	"kumquat/renderer"
	"kumquat/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResultStringValidIndex(t *testing.T) {
	output := []string{"resource1", "resource2", "resource3"}
	rendererOutput := renderer.NewOutputFromSlice(output)

	result, err := rendererOutput.ResultString(1)
	require.NoError(t, err)
	assert.Equal(t, "resource2", result)
}

func TestResultStringInvalidIndexNegative(t *testing.T) {
	output := []string{"resource1", "resource2", "resource3"}
	rendererOutput := renderer.NewOutputFromSlice(output)

	_, err := rendererOutput.ResultString(-1)
	assert.Error(t, err)
}

func TestResultStringInvalidIndexOutOfRange(t *testing.T) {
	output := []string{"resource1", "resource2", "resource3"}
	rendererOutput := renderer.NewOutputFromSlice(output)

	_, err := rendererOutput.ResultString(3)
	assert.Error(t, err)
}

func TestResultStringEmptyOutput(t *testing.T) {
	output := []string{}
	rendererOutput := renderer.NewOutputFromSlice(output)

	_, err := rendererOutput.ResultString(0)
	assert.Error(t, err)
}

func TestRegisterDuplicateRenderer(t *testing.T) {
	err := renderer.Register("test", func(s1, s2 string) (renderer.Renderer, error) {
		return nil, nil
	})
	assert.NoError(t, err)

	err = renderer.Register("test", func(s1, s2 string) (renderer.Renderer, error) {
		return nil, nil
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "'test' already registered")
}

type MockRenderer struct {
	render func(results any, output *renderer.Output) error
}

func (r *MockRenderer) Render(results any, output *renderer.Output) error {
	if r.render == nil {
		return errors.ErrUnsupported
	}

	return r.render(results, output)
}

func explodeOnErr[X any](x X, err error) X {
	if err != nil {
		panic(err)
	}

	return x
}

func getResourceContent() map[string]any {
	return map[string]any{
		"apiVersion": "guidewire.com/v1beta1",
		"kind":       "Test",
		"metadata": map[string]any{
			"name": "test",
		},
	}
}

func TestRenderBatchMode(t *testing.T) {
	res1 := getResourceContent()
	res1["metadata"].(map[string]any)["name"] = "test1"
	res2 := getResourceContent()
	res2["metadata"].(map[string]any)["name"] = "test2"
	r := new(MockRenderer)
	results := []map[string]repository.Resource{
		{"resource1": explodeOnErr(repository.MakeResource(res1))},
		{"resource1": explodeOnErr(repository.MakeResource(res2))},
	}

	t.Run("RendererGetsTwoRows", func(t *testing.T) {
		r.render = func(results any, output *renderer.Output) error {
			res := results.([]map[string]any)
			assert.Len(t, results, 2)
			assert.Equal(t, "test1", res[0]["resource1"].(map[string]any)["metadata"].(map[string]any)["name"])
			assert.Equal(t, "test2", res[1]["resource1"].(map[string]any)["metadata"].(map[string]any)["name"])
			output.Append("expected output")
			return nil
		}
		out, err := renderer.Render(r, results, true)
		require.NoError(t, err)
		require.NotNil(t, out)
		assert.Equal(t, 1, out.ResourceCount())
		assert.Equal(t, "expected output", explodeOnErr(out.ResultString(0)))
	})
}

func TestRenderNonBatchMode(t *testing.T) {
	res1 := getResourceContent()
	res1["metadata"].(map[string]any)["name"] = "test1"
	res2 := getResourceContent()
	res2["metadata"].(map[string]any)["name"] = "test2"
	r := new(MockRenderer)
	results := []map[string]repository.Resource{
		{"resource1": explodeOnErr(repository.MakeResource(res1))},
		{"resource1": explodeOnErr(repository.MakeResource(res2))},
	}

	t.Run("RendererGetsOneRowTwice", func(t *testing.T) {
		count := 0
		r.render = func(results any, output *renderer.Output) error {
			res := results.(map[string]any)
			count += 1
			assert.Equal(t,
				fmt.Sprintf("test%d", count),
				res["resource1"].(map[string]any)["metadata"].(map[string]any)["name"],
			)
			output.Append(fmt.Sprintf("expected output %d", count))
			return nil
		}
		out, err := renderer.Render(r, results, false)
		require.NoError(t, err)
		require.NotNil(t, out)
		assert.Equal(t, 2, out.ResourceCount())
		assert.Equal(t, "expected output 1", explodeOnErr(out.ResultString(0)))
		assert.Equal(t, "expected output 2", explodeOnErr(out.ResultString(1)))
	})
}
