package gotemplate

import (
	"bytes"
	"errors"
	"fmt"
	"kumquat/renderer"
	"regexp"
	"strconv"

	"text/template"
)

func init() {
	err := renderer.Register("gotemplate", func(template, source string) (renderer.Renderer, error) {
		return NewGoRenderer(template, source)
	})
	if err != nil {
		panic(err)
	}
}

type GoRenderer struct {
	template *template.Template
	source   string
}

func (r *GoRenderer) Render(results any, output *renderer.Output) error {
	var buffer bytes.Buffer
	err := r.template.Execute(&buffer, results)

	if err != nil {
		return fmt.Errorf("error executing Go template: %w", newRendererError(err, r.source))
	}

	output.Append(buffer.String())

	return nil
}

func NewGoRenderer(tmpl string, source string) (*GoRenderer, error) {
	t, err := template.New(source).Option("missingkey=zero").Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("error parsing Go template: %w", newRendererError(err, source))
	}

	return &GoRenderer{template: t, source: source}, nil
}

var lineRE = regexp.MustCompile(`^template: .*?:(\d+)(:(\d+))?: (.*)`)

func newRendererError(err error, source string) *renderer.Error {
	s := err.Error()

	// parse the line number from the error message
	var line, column int
	matches := lineRE.FindStringSubmatch(s)

	if matches != nil {
		// parse error returns 0; same as line unknown
		line, _ = strconv.Atoi(matches[1])

		// parse error returns 0; same as column unknown
		column, _ = strconv.Atoi(matches[3])

		s = matches[4]
	}

	// replace "sourcefile:line#" with "line line#" in the remainder of the error message
	re := regexp.MustCompile(source + `:(\d+)`)
	s = re.ReplaceAllString(s, "line $1")

	return renderer.NewError(errors.New(s), line, column)
}
