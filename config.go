package pkgen

import "github.com/ifnotnil/pkgen/templates"

type Config struct {
	PackagesQuery PackagesQueryConfig       `yaml:"packages_query"`
	Templates     templates.TemplateConfigs `yaml:"templates"`
	Generate      GenerateConfig            `yaml:"generate"`
}

var DefaultConfig = Config{
	PackagesQuery: PackagesQueryConfig{
		IncludeTests: false,
		Env:          nil,
		BuildFlags:   nil,
		Dir:          "",
		Patterns:     []string{"./..."},
	},

	Templates: templates.TemplateConfigs{},

	Generate: GenerateConfig{
		OutputFile:    defaultNameFMT,
		OutputFileMod: 0o644,
	},
}
