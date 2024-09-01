package template_test

import (
	"errors"
	"kumquat/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationErrorsWithNoErrors(t *testing.T) {
	err := template.NewValidationErrors(t.Name())
	assert.Equal(t, "Template 'TestValidationErrorsWithNoErrors': unspecified validation error", err.Error())
	assert.Equal(t, t.Name(), err.Template())
}

func TestValidationErrorsWithErrors(t *testing.T) {
	err := template.NewValidationErrors(t.Name())
	err.Append(errors.New("first error"))
	err.Append(errors.New("second error"))
	assert.Equal(t, "Invalid Template 'TestValidationErrorsWithErrors':\nfirst error\nsecond error", err.Error())
	assert.True(t, err.HasErrors())
}
