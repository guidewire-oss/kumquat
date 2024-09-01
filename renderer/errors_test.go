package renderer_test

import (
	"errors"
	"kumquat/renderer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorFormatting(t *testing.T) {
	t.Run("NoLineNoColumn", func(t *testing.T) {
		err := renderer.NewError(errors.New("test error"), 0, 0)
		assert.Equal(t, "[line ?, column ?] test error", err.Error())
	})

	t.Run("LineNoColumn", func(t *testing.T) {
		err := renderer.NewError(errors.New("test error"), 1, 0)
		assert.Equal(t, "[line 1, column ?] test error", err.Error())
	})

	t.Run("LineAndColumn", func(t *testing.T) {
		err := renderer.NewError(errors.New("test error"), 1, 2)
		assert.Equal(t, "[line 1, column 2] test error", err.Error())
	})
}

func TestErrorUnwrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := renderer.NewError(originalErr, 0, 0)
	assert.Equal(t, originalErr, err.Unwrap())
}

func TestErrorLine(t *testing.T) {
	err := renderer.NewError(errors.New("test error"), 3, 0)
	assert.Equal(t, 3, err.Line())
}

func TestErrorColumn(t *testing.T) {
	err := renderer.NewError(errors.New("test error"), 0, 4)
	assert.Equal(t, 0, err.Column())

	err = renderer.NewError(errors.New("test error"), 5, 4)
	assert.Equal(t, 4, err.Column())
}
