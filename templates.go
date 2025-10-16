package pkgen

import (
	"embed"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

var ErrTemplateNotFound = errors.New("template not found")

type Templates struct{}

func (t Templates) Get(name string) (*template.Template, error) {
	b, err := templatesFS.ReadFile(path.Join("templates", name+".tmpl"))
	if err != nil {
		return nil, errors.Join(ErrTemplateNotFound, err)
	}

	return template.New(name).Parse(string(b))
}

func (t Templates) customTemplate(filePath string) (*template.Template, error) {
	b, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, errors.Join(ErrTemplateNotFound, err)
	}

	name := filepath.Base(filePath)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	return template.New(name).Parse(string(b))
}

func (t Templates) GetAll(c TemplateConfigs) ([]*template.Template, error) {
	if len(c) == 0 {
		return []*template.Template{}, nil
	}

	sl := make([]*template.Template, 0, len(c))
	for _, cnf := range c {
		switch {
		case cnf.Name != "":
			t, err := t.Get(cnf.Name)
			if err != nil {
				return nil, err
			}
			sl = append(sl, t)
		case cnf.CustomTemplateFile != "":
			t, err := t.customTemplate(cnf.CustomTemplateFile)
			if err != nil {
				return nil, err
			}
			sl = append(sl, t)
		}
	}

	return sl, nil
}
