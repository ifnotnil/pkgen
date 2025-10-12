package pkgen

import (
	"context"
	"os"

	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
)

type Packages struct{}

func (Packages) Query(ctx context.Context, q PackagesQueryConfig) ([]packages.Package, error) {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedModule | packages.NeedFiles,
		Context:    ctx,
		Tests:      q.IncludeTests,
		Dir:        q.Dir,
		Env:        append(os.Environ(), q.Env...),
		BuildFlags: q.BuildFlags,
	}

	p, err := packages.Load(cfg, q.Patterns...)
	if err != nil {
		return nil, err
	}

	return lo.FilterMap(p, func(item *packages.Package, index int) (packages.Package, bool) {
		if item == nil {
			return packages.Package{}, false
		}
		return *item, true
	}), nil
}
