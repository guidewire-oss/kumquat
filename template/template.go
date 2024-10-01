package template

import (
	"errors"
	"fmt"
	"kumquat/renderer"
	"kumquat/renderer/gotemplate"
	"kumquat/repository"
	"strings"
)

const TemplateAPIGroup = "kumquat.guidewire.com"
const TemplateKind = "Template"
const TemplateResourceType = TemplateKind + "." + TemplateAPIGroup

type Template struct {
	name             string
	query            string
	fileNameTemplate renderer.Renderer
	renderer         renderer.Renderer
	batchMode        bool
}

type TemplateOutput struct {
	Output    *renderer.Output
	FileNames []string
}

func NewTemplate(r repository.Resource) (*Template, error) {
	// name is guaranteed to be present since it comes from a Resource
	name := r.Content()["metadata"].(map[string]any)["name"].(string)
	var validationError = &ValidationErrors{template: name}

	// apiVersion is guaranteed to be a string with a "/" in it because it comes from a Resource
	apiVersion := r.Content()["apiVersion"].(string)
	s := strings.SplitN(apiVersion, "/", 2)

	if len(s) != 2 || s[0] != TemplateAPIGroup {
		validationError.Append(&FieldValidationError{
			template: name,
			field:    "apiVersion",
			err:      fmt.Errorf("'%s' should be '%s'", s[0], TemplateAPIGroup),
		})
	}

	if s[1] != "v1beta1" {
		validationError.Append(&FieldValidationError{
			template: name,
			field:    "apiVersion",
			err:      fmt.Errorf("'%s' unsupported; supported versions are: 'v1beta1'", s[1]),
		})
	}

	// kind is guaranteed to be a string because it comes from a Resource
	kind := r.Content()["kind"].(string)

	if kind != TemplateKind {
		validationError.Append(&FieldValidationError{
			template: name,
			field:    "kind",
			err:      fmt.Errorf("'%s' should be '%s'", kind, TemplateKind),
		})
	}

	spec, ok := r.Content()["spec"].(map[string]any)
	if !ok {
		validationError.Append(&FieldValidationError{
			template: name,
			field:    "spec",
			err:      errors.New("missing or not a map"),
		})
	}

	var query string
	var template map[string]any

	if len(spec) > 0 {
		query, _ = spec["query"].(string)
		if query == "" {
			validationError.Append(&FieldValidationError{
				template: name,
				field:    "spec.query",
				err:      errors.New("missing or not a string"),
			})
		}

		template, _ = spec["template"].(map[string]any)
		if len(template) == 0 {
			validationError.Append(&FieldValidationError{
				template: name,
				field:    "spec.template",
				err:      errors.New("missing or not a map"),
			})
		}
	}

	var lang, data, fileNameTemplate string

	if len(template) > 0 {
		data, _ = template["data"].(string)
		if data == "" {
			validationError.Append(&FieldValidationError{
				template: name,
				field:    "spec.template.data",
				err:      errors.New("missing or not a string"),
			})
		}

		lang, _ = template["language"].(string)
		if lang == "" {
			validationError.Append(&FieldValidationError{
				template: name,
				field:    "spec.template.language",
				err:      errors.New("missing or not a string"),
			})
		}

		fileNameTemplate, _ = template["fileName"].(string)
		if fileNameTemplate == "" {
			validationError.Append(&FieldValidationError{
				template: name,
				field:    "spec.template.fileName",
				err:      errors.New("missing or not a string"),
			})
		}
	}

	// batchModeProcessing false unless explicitly set to true in template
	batchModeProcessing := false
	if val, ok := template["batchModeProcessing"].(bool); ok {
		batchModeProcessing = val
	}

	t := &Template{
		name:      name,
		query:     query,
		batchMode: batchModeProcessing,
	}

	var err error

	if lang != "" {
		t.renderer, err = renderer.MakeRenderer(lang, data, name)
	}

	if err != nil {
		if _, ok := err.(*renderer.LookupError); ok {
			validationError.Append(&FieldValidationError{
				template: name,
				field:    "spec.template.language",
				err:      err,
			})
		} else {
			fmt.Printf("Error making renderer: %#v\n", err)
			validationError.Append(&FieldValidationError{
				template: name,
				field:    "spec.template.data",
				err:      err,
			})
		}
	}

	err = nil
	if fileNameTemplate != "" {
		t.fileNameTemplate, err = gotemplate.NewGoRenderer(fileNameTemplate, name)
	}

	if err != nil {
		validationError.Append(&FieldValidationError{
			template: name,
			field:    "spec.template.fileName",
			err:      err,
		})
	}

	if validationError.HasErrors() {
		return nil, validationError
	}

	return t, nil
}

// Evaluate runs the query and renders the results using the template.
func (t *Template) Evaluate(repo repository.Repository) (*TemplateOutput, error) {
	resultset, err := repo.Query(t.query)
	if err != nil {
		return nil, fmt.Errorf("query failed in Template '%s': %w", t.name, err)
	}

	output, err := renderer.Render(t.renderer, resultset.Results, t.batchMode)
	if err != nil {
		return nil, fmt.Errorf("error rendering Template '%s': %w", t.name, err)
	}

	fileNamesOutput, err := renderer.Render(t.fileNameTemplate, resultset.Results, t.batchMode)
	if err != nil {
		return nil, fmt.Errorf("error rendering fileName in Template '%s': %w", t.name, err)
	}

	// TODO revisit this when we start outputting resources to Kubernetes directly
	// When putting results in files, these invariants should hold:
	if t.batchMode {
		if fileNamesOutput.ResourceCount() != 1 {
			return nil, fmt.Errorf("expected one file name in batch mode; got %d", fileNamesOutput.ResourceCount())
		}

		if output.ResourceCount() != 1 {
			return nil, fmt.Errorf("expected one rendered resource in batch mode; got %d", output.ResourceCount())
		}
	} else {
		if fileNamesOutput.ResourceCount() != len(resultset.Results) {
			return nil, fmt.Errorf("expected %d file names; got %d", len(resultset.Results), fileNamesOutput.ResourceCount())
		}

		if output.ResourceCount() != len(resultset.Results) {
			return nil, fmt.Errorf("expected %d rendered resources; got %d", len(resultset.Results), output.ResourceCount())
		}
	}

	fileNames := make([]string, fileNamesOutput.ResourceCount())
	for i := 0; i < fileNamesOutput.ResourceCount(); i++ {
		fileName, err := fileNamesOutput.ResultString(i)
		if err != nil {
			return nil, fmt.Errorf("error getting filename: %w", err)
		}
		fileNames[i] = strings.TrimSpace(fileName)
	}

	return &TemplateOutput{Output: output, FileNames: fileNames}, nil
}

func (t *Template) BatchMode() bool {
	return t.batchMode
}

func (t *Template) Name() string {
	return t.name
}
