package pkgen

import (
	"context"
	"os"

	"golang.org/x/tools/go/packages"
)

type PackagesQueryConfig struct {
	IncludeTests bool     `yaml:"include_tests"`
	Env          []string `yaml:"env"`
	BuildFlags   []string `yaml:"build_flags"`

	Dir      string   `yaml:"dir"`
	Patterns []string `yaml:"patterns"` // e.g. "./..."
}

func Packages(ctx context.Context, q PackagesQueryConfig) ([]*packages.Package, error) {
	if len(q.Patterns) == 0 {
		q.Patterns = []string{"./..."}
	}

	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedModule | packages.NeedFiles,
		Context:    ctx,
		Tests:      q.IncludeTests,
		Dir:        q.Dir,
		Env:        append(os.Environ(), q.Env...),
		BuildFlags: q.BuildFlags,
	}

	return packages.Load(cfg, q.Patterns...)
}
