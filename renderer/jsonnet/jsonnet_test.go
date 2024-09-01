package jsonnet_test

import (
	"kumquat/renderer"
	"kumquat/renderer/jsonnet"
	"kumquat/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJsonnetRenderer(t *testing.T) {
	t.Run("OkTemplate", func(t *testing.T) {
		r, err := jsonnet.NewJsonnetRenderer("1", t.Name())
		assert.NoError(t, err)
		assert.NotNil(t, r)
	})

	t.Run("BadTemplate", func(t *testing.T) {
		_, err := jsonnet.NewJsonnetRenderer("this template is erroneous", t.Name())
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error parsing Jsonnet: [line 1, column 6]")

		var rerr *renderer.Error
		require.ErrorAs(t, err, &rerr)
		assert.EqualValues(t, 1, rerr.Line())
		assert.EqualValues(t, 6, rerr.Column())
		assert.ErrorContains(t, rerr, "Did not expect")
	})

	t.Run("MalformedTemplate", func(t *testing.T) {
		_, err := jsonnet.NewJsonnetRenderer(`{"foo": "bar"`, t.Name())
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error parsing Jsonnet: [line 1, column 14]")

		var rerr *renderer.Error
		require.ErrorAs(t, err, &rerr)
		assert.EqualValues(t, 1, rerr.Line())
		assert.EqualValues(t, 14, rerr.Column())
		assert.ErrorContains(t, rerr, "Expected a comma")
	})
}

func TestJsonnetRendererWithNoResults(t *testing.T) {
	r, err := jsonnet.NewJsonnetRenderer(`{
		person1: {
			name: "Alice",
			welcome: "Hello " + self.name + "!",
		},
		person2: self.person1 { name: "Bob" },
	}`, t.Name())
	require.NoError(t, err)
	require.NotNil(t, r)

	res, err := renderer.Render(r, []map[string]repository.Resource{}, false)
	require.NoError(t, err)
	assert.Zero(t, res.ResourceCount())

	res, err = renderer.Render(r, []map[string]repository.Resource{}, true)
	require.NoError(t, err)
	assert.Equal(t, 1, res.ResourceCount())
	s, err := res.ResultString(0)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"person1": {
			"name": "Alice",
			"welcome": "Hello Alice!"
		},
		"person2": {
			"name": "Bob",
			"welcome": "Hello Bob!"
		}
	}`, s)
}

func TestJsonnetRendererWithOneResult(t *testing.T) {
	res1, err := repository.MakeResource(map[string]any{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]any{
			"name":      "my-config",
			"namespace": "default",
		},
		"data": map[string]any{
			"firstName": "Grace",
			"lastName":  "Hopper",
		},
	})
	require.NoError(t, err)

	t.Run("NonBatch", func(t *testing.T) {
		r, err := jsonnet.NewJsonnetRenderer(`{
			person1: std.extVar('data').x.data.firstName
		}`, t.Name())
		require.NoError(t, err)
		require.NotNil(t, r)

		rs := []map[string]repository.Resource{
			{"x": res1},
		}

		res, err := renderer.Render(r, rs, false)
		require.NoError(t, err)
		assert.Equal(t, 1, res.ResourceCount())
		s, err := res.ResultString(0)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"person1": "Grace"
		}`, s)
	})

	t.Run("Batch", func(t *testing.T) {
		r, err := jsonnet.NewJsonnetRenderer(`{
			['person' + i]: std.extVar('data')[i-1].x.data.firstName for i in std.range(1, std.length(std.extVar('data')))
		}`, t.Name())
		require.NoError(t, err)
		require.NotNil(t, r)

		rs := []map[string]repository.Resource{
			{"x": res1},
		}

		res, err := renderer.Render(r, rs, true)
		require.NoError(t, err)
		assert.Equal(t, 1, res.ResourceCount())
		s, err := res.ResultString(0)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"person1": "Grace"
		}`, s)
	})

	t.Run("NonBatchWithTemplateError", func(t *testing.T) {
		r, err := jsonnet.NewJsonnetRenderer(`{
			person1: std.extVar('notdefined').x.data.firstName
		}`, t.Name())
		require.NoError(t, err)
		require.NotNil(t, r)

		rs := []map[string]repository.Resource{
			{"x": res1},
		}

		_, err = renderer.Render(r, rs, false)
		var rerr *renderer.Error
		require.ErrorAs(t, err, &rerr)
		assert.EqualValues(t, 2, rerr.Line())
		assert.EqualValues(t, 13, rerr.Column())
		assert.ErrorContains(t, rerr, "Undefined external variable: notdefined")
	})
}

func TestJsonnetRendererWithTwoResults(t *testing.T) {
	rs1, err := repository.MakeResource(map[string]any{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]any{
			"name":      "my-config",
			"namespace": "default",
		},
		"data": map[string]any{
			"firstName": "Grace",
			"lastName":  "Hopper",
		},
	})
	require.NoError(t, err)

	rs2, err := repository.MakeResource(map[string]any{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]any{
			"name":      "my-config",
			"namespace": "default",
		},
		"data": map[string]any{
			"firstName": "Ada",
			"lastName":  "Lovelace",
		},
	})
	require.NoError(t, err)

	t.Run("NonBatch", func(t *testing.T) {
		r, err := jsonnet.NewJsonnetRenderer(`{
			person1: std.extVar('data').x.data.firstName
		}`, t.Name())
		require.NoError(t, err)
		require.NotNil(t, r)

		rs := []map[string]repository.Resource{
			{"x": rs1},
			{"x": rs2},
		}

		res, err := renderer.Render(r, rs, false)
		require.NoError(t, err)
		assert.Equal(t, 2, res.ResourceCount())
		s, err := res.ResultString(0)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"person1": "Grace"
		}`, s)
		s, err = res.ResultString(1)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"person1": "Ada"
		}`, s)
	})

	t.Run("Batch", func(t *testing.T) {
		r, err := jsonnet.NewJsonnetRenderer(`{
			['person' + i]: std.extVar('data')[i-1].x.data.firstName for i in std.range(1, std.length(std.extVar('data')))
		}`, t.Name())
		require.NoError(t, err)
		require.NotNil(t, r)

		rs := []map[string]repository.Resource{
			{"x": rs1},
			{"x": rs2},
		}

		res, err := renderer.Render(r, rs, true)
		require.NoError(t, err)
		assert.Equal(t, 1, res.ResourceCount())
		s, err := res.ResultString(0)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"person1": "Grace",
			"person2": "Ada"
		}`, s)
	})
}
