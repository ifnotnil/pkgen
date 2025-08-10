package pkgen

import "github.com/ifnotnil/pkgen/templates"

type Config struct {
	PackagesQuery PackagesQueryConfig       `yaml:"packages_query"`
	Templates     templates.TemplateConfigs `yaml:"templates"`
	Generate      GenerateConfig            `yaml:"generate"`
}
