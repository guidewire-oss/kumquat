package template

import (
	"fmt"
	"strings"
)

type FieldValidationError struct {
	template string
	field    string
	err      error
}

func (e *FieldValidationError) Error() string {
	return fmt.Sprintf("'%s'.%s: %v", e.template, e.field, e.Unwrap())
}

func (e *FieldValidationError) Unwrap() error {
	return e.err
}

type ValidationErrors struct {
	template string
	errors   []error
}

func NewValidationErrors(template string) *ValidationErrors {
	return &ValidationErrors{template: template}
}

func (e *ValidationErrors) Error() string {
	if !e.HasErrors() {
		return fmt.Sprintf("Template '%s': unspecified validation error", e.template)
	}

	var errs []string = make([]string, len(e.errors)+1)
	errs[0] = fmt.Sprintf("Invalid %s '%s':", TemplateKind, e.template)
	for i, err := range e.errors {
		errs[i+1] = err.Error()
	}
	return strings.Join(errs, "\n")
}

func (e *ValidationErrors) Unwrap() []error {
	return e.errors
}

func (e *ValidationErrors) Append(err error) {
	if err != nil {
		if e.errors == nil {
			e.errors = make([]error, 0, 1)
		}

		e.errors = append(e.errors, err)
	}
}

func (e *ValidationErrors) HasErrors() bool {
	return len(e.errors) > 0
}

func (e *ValidationErrors) Template() string {
	return e.template
}
