package renderer

import (
	"fmt"
	"kumquat/repository"
)

type Renderer interface {
	// Render the results and add the output to the provided output object. The results may be a single result
	// of type `map[string]any` or a slice of such results (i.e. `[]map[string]any`).
	Render(results any, output *Output) error
}

// Render the results using the provided renderer. If batchMode is true, the renderer will be called once
// with all the results; otherwise, the renderer will be called once for each result.
func Render(renderer Renderer, results []map[string]*repository.Resource, batchMode bool) (*Output, error) {
	var o *Output
	resultsWithoutResources := StripResourcesFromResults(results)

	if batchMode {
		o = NewOutput(1)
		err := renderer.Render(resultsWithoutResources, o)
		if err != nil {
			return nil, err
		}
	} else {
		o = NewOutput(len(results))

		for _, r := range resultsWithoutResources {
			err := renderer.Render(r, o)
			if err != nil {
				return nil, err
			}
		}
	}

	return o, nil
}

type RendererMaker func(string, string) (Renderer, error)

var rendererRegistry map[string]RendererMaker

func Register(name string, f RendererMaker) error {
	if rendererRegistry == nil {
		rendererRegistry = make(map[string]RendererMaker)
	}

	if rendererRegistry[name] != nil {
		return fmt.Errorf("renderer '%s' already registered", name)
	}

	fmt.Printf("Renderer '%s' registered.\n", name)
	rendererRegistry[name] = f

	return nil
}

func MakeRenderer(name, template, source string) (Renderer, error) {
	f, ok := rendererRegistry[name]
	if !ok {
		return nil, &LookupError{rendererName: name}
	}

	return f(template, source)
}

type Output struct {
	output []string
}

func NewOutput(expectedSize int) *Output {
	return &Output{output: make([]string, 0, expectedSize)}
}

func NewOutputFromSlice(output []string) *Output {
	return &Output{output: output}
}

func (o *Output) Append(result string) {
	o.output = append(o.output, result)
}

func (o *Output) ResourceCount() int {
	return len(o.output)
}

func (o *Output) ResultString(resource int) (string, error) {
	if resource < 0 || resource >= len(o.output) {
		return "", fmt.Errorf("resource index out of range")
	}
	return o.output[resource], nil
}
