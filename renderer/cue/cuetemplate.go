package cue

import (
	"fmt"
	"kumquat/renderer"
	"strings"

	cuelang "cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cue_errors "cuelang.org/go/cue/errors"
	"gopkg.in/yaml.v3"
)

func init() {
	err := renderer.Register("cue", func(template, source string) (renderer.Renderer, error) {
		return NewCUERenderer(template, source)
	})
	if err != nil {
		panic(err)
	}
}

type CUERenderer struct {
	config string
	source string
}

func NewCUERenderer(template string, source string) (*CUERenderer, error) {
	tpl := CUERenderer{
		config: template,
		source: source,
	}

	return &tpl, nil
}

func (r *CUERenderer) Render(results any, output *renderer.Output) error {
	err := r.evaluate(results, output)

	if err != nil {
		return fmt.Errorf("error evaluating CUE template '%s': %w", r.source, err)
	}

	return nil
}

func (t *CUERenderer) evaluate(r any, o *renderer.Output) error {
	c := cuecontext.New()
	v2 := c.Encode(r)
	resultsObj, err := v2.Eval().MarshalJSON()
	if err != nil {
		return fmt.Errorf("error converting results to CUE: %w", err)
	}

	suffix := "\n\ndata: " + string(resultsObj) + "\n"
	v := c.CompileString(t.config + suffix).Eval()
	if v.Err() != nil {
		return newRendererError(v.Err(), 0)
	}

	return appendOutput(o, v)
}

func appendOutput(o *renderer.Output, v cuelang.Value) error {
	oPath := cuelang.ParsePath("out")

	switch t := v.LookupPath(oPath).Kind(); t {
	case cuelang.ListKind:
		var output []map[string]any

		err := v.LookupPath(oPath).Decode(&output)
		if err != nil {
			return fmt.Errorf("error decoding 'out': %w", err)
		}

		var outputs []string
		for i := 0; i < len(output); i++ {
			outputByteArray, err := yaml.Marshal(output[i])
			if err != nil {
				return fmt.Errorf("error decoding 'out': %w", err)
			}

			outputs = append(outputs, string(outputByteArray))
		}

		o.Append(strings.Join(outputs, "---\n"))

	case cuelang.StructKind:
		var output map[string]any
		err := v.LookupPath(oPath).Decode(&output)
		if err != nil {
			return fmt.Errorf("error decoding 'out': %w", err)
		}
		//convert output to string
		outputByteArray, err := yaml.Marshal(output)
		if err != nil {
			return fmt.Errorf("error decoding 'out': %w", err)
		}

		o.Append(string(outputByteArray))

	case cuelang.BottomKind:
		return fmt.Errorf("'out' is not set to anything concrete")
	default:
		return fmt.Errorf("'out' has unsupported output type '%v'", t)
	}

	return nil
}

func newRendererError(err error, prefixLength int) *renderer.Error {
	pos := cue_errors.Positions(err)

	if len(pos) > 0 {
		return renderer.NewError(err, pos[0].Line()-prefixLength, pos[0].Column())
	}

	return renderer.NewError(err, 0, 0)
}
