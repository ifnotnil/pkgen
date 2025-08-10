package templates

import (
	"embed"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/moukoublen/pick"
	"github.com/moukoublen/pick/iter"
	"sigs.k8s.io/yaml"
)

//go:embed *.tmpl
var templatesFS embed.FS

var ErrTemplateNotFound = errors.New("template not found")

func Get(name string) (*template.Template, error) {
	b, err := templatesFS.ReadFile(name + ".tmpl")
	if err != nil {
		return nil, errors.Join(ErrTemplateNotFound, err)
	}

	return template.New(name).Parse(string(b))
}

func customTemplate(filePath string) (*template.Template, error) {
	b, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, errors.Join(ErrTemplateNotFound, err)
	}

	name := filepath.Base(filePath)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	return template.New(name).Parse(string(b))
}

func GetTemplates(c TemplateConfigs) ([]*template.Template, error) {
	if len(c) == 0 {
		return []*template.Template{}, nil
	}

	sl := make([]*template.Template, 0, len(c))
	for _, cnf := range c {
		switch {
		case cnf.Name != "":
			t, err := Get(cnf.Name)
			if err != nil {
				return nil, err
			}
			sl = append(sl, t)
		case cnf.CustomTemplateFile != "":
			t, err := customTemplate(cnf.CustomTemplateFile)
			if err != nil {
				return nil, err
			}
			sl = append(sl, t)
		}
	}

	return sl, nil
}

type TemplateConfig struct {
	Name               string `yaml:"name"`
	CustomTemplateFile string `yaml:"template_file"`
}

type TemplateConfigs []TemplateConfig

func (tc *TemplateConfigs) UnmarshalYAML(value []byte) error {
	var raw any
	if err := yaml.Unmarshal(value, &raw); err != nil {
		return err
	}

	// can be a single string
	if i, _ := iter.Len(raw); i > 0 {
		p := pick.Wrap(raw).Relaxed()
		*tc = pick.RelaxedMap(p, "", func(rp pick.RelaxedAPI) (TemplateConfig, error) {
			// can be a string
			s := rp.String("")
			if s != "" {
				return TemplateConfig{Name: s}, nil
			}

			// or an object
			tc := TemplateConfig{
				Name:               rp.String("name"),
				CustomTemplateFile: rp.String("template_file"),
			}

			if tc.Name == "" && tc.CustomTemplateFile == "" {
				return tc, ErrMalformedTemplateConfigs
			}

			return tc, nil
		})

		if len(*tc) != i {
			return ErrMalformedTemplateConfigs
		}
	}

	s, err := pick.Convert[string](raw)
	if err != nil {
		return ErrMalformedTemplateConfigs
	}

	(*tc) = []TemplateConfig{{Name: s}}

	return nil
}

var ErrMalformedTemplateConfigs = errors.New("malformed template configs")
