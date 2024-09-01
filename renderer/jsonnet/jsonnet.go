package jsonnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"kumquat/renderer"
	"reflect"

	js "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

func init() {
	err := renderer.Register("jsonnet", func(template, source string) (renderer.Renderer, error) {
		return NewJsonnetRenderer(template, source)
	})
	if err != nil {
		panic(err)
	}
}

type JsonnetRenderer struct {
	vm       *js.VM
	template ast.Node
}

func NewJsonnetRenderer(template string, source string) (*JsonnetRenderer, error) {
	parsed, err := js.SnippetToAST(source, template)

	if err != nil {
		return nil, fmt.Errorf("error parsing Jsonnet: %w", newRendererError(err))
	}

	return &JsonnetRenderer{template: parsed, vm: js.MakeVM()}, nil
}

func newRendererError(err error) *renderer.Error {
	re, ok := err.(js.RuntimeError)

	if ok {
		for _, frame := range re.StackTrace {
			if frame.Loc.Begin.Line > 0 {
				return renderer.NewError(err, frame.Loc.Begin.Line, frame.Loc.Begin.Column)
			}
		}
	} else {
		// Not a RuntimeError, maybe we can find a StaticError
		for e := err; e != nil; e = errors.Unwrap(err) {
			ev := reflect.ValueOf(e)
			mv := ev.MethodByName("Loc") // StaticError has method Loc() returning LocationRange

			if mv != reflect.ValueOf(nil) {
				rvs := mv.Call(nil)               // Call Loc() to get LocationRange
				bv := rvs[0].FieldByName("Begin") // LocationRange has a Begin field of type Location
				lv := bv.FieldByName("Line")      // Location has Line
				cv := bv.FieldByName("Column")    // Location has Column

				return renderer.NewError(err, int(lv.Int()), int(cv.Int()))
			}
		}
	}

	return renderer.NewError(err, 0, 0)
}

func (r *JsonnetRenderer) Render(results any, output *renderer.Output) error {
	b, err := json.Marshal(results)

	if err != nil {
		return fmt.Errorf("error converting results to Jsonnet: %w", err)
	}

	code := string(b)

	r.vm.ExtCode("data", code)
	js, err := r.vm.Evaluate(r.template)

	if err != nil {
		return fmt.Errorf("error rendering Jsonnet: %w", newRendererError(err))
	}

	output.Append(js)

	return nil
}
