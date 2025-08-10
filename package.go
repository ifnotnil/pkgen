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

	Dir     string `yaml:"build_flags"`
	Pattern string `yaml:"build_flags"` // e.g. "./..."
}

func Packages(ctx context.Context, q PackagesQueryConfig) ([]*packages.Package, error) {
	if q.Pattern == "" {
		q.Pattern = "./..."
	}

	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedModule | packages.NeedFiles,
		Context:    ctx,
		Tests:      q.IncludeTests,
		Dir:        q.Dir,
		Env:        append(os.Environ(), q.Env...),
		BuildFlags: q.BuildFlags,
	}

	return packages.Load(cfg, q.Pattern)
}
