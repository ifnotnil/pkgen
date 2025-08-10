package pkgen

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ifnotnil/pkgen/templates"
	"go.yaml.in/yaml/v4"
)

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
		OutputFile:    "",
		OutputFileMod: 0o644,
	},
}

func ParseYAMLConfig(filePath string, cnf *Config) error {
	filePath = filepath.Clean(filePath)

	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if b, err := os.ReadFile(filePath); err == nil {
		err = yaml.Unmarshal(b, &cnf)
		if err != nil {
			return err
		}
	}

	return nil
}

func ParseConfig(ctx context.Context, configPath string) (*Config, error) {
	cnf := DefaultConfig

	if runningInsideGoGenerate() {
		// if it is running inside go:generate query only the local package.
		cnf.PackagesQuery.Patterns = []string{"."}
		slog.Default().DebugContext(ctx, "running inside a go:generate")
	} else {
		err := ParseYAMLConfig(configPath, &cnf)
		return nil, err
	}

	return &cnf, nil
}

func runningInsideGoGenerate() bool {
	_, exists := os.LookupEnv("GOFILE")

	return exists
}
