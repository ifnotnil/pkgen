package templates

import (
	"embed"
	"errors"
	"text/template"
)

//go:embed *.tmpl
var templates embed.FS

var ErrTemplateNotFound = errors.New("template not found")

func Get(name string) (*template.Template, error) {
	b, err := templates.ReadFile(name + ".tmpl")
	if err != nil {
		return nil, errors.Join(ErrTemplateNotFound, err)
	}

	return template.New(name).Parse(string(b))
}
