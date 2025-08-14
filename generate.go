package pkgen

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/go/packages"
)

const defaultNameFMT = "zz_generated.%s.go"

type GenerateConfig struct {
	OutputFile    string      `yaml:"output"` // the default pattern is zz_generated.{{template name}}.go
	OutputFileMod os.FileMode `yaml:"mod"`
}

func GenerateInPackage(ctx context.Context, pkg *packages.Package, tmp *template.Template, cnf GenerateConfig) error {
	if len(pkg.GoFiles) == 0 {
		return nil
	}

	// default values
	if cnf.OutputFile == "" {
		cnf.OutputFile = fmt.Sprintf(defaultNameFMT, tmp.Name())
	}
	if cnf.OutputFileMod == 0 {
		cnf.OutputFileMod = 0o644
	}

	// execute the template
	buf := bytes.Buffer{}
	err := tmp.Execute(&buf, pkg)
	if err != nil {
		return err
	}

	outPath := filepath.Join(filepath.Clean(pkg.Dir), cnf.OutputFile)
	return os.WriteFile(outPath, buf.Bytes(), cnf.OutputFileMod)
}

func Generate(ctx context.Context, logger *slog.Logger, pkgs []*packages.Package, tmps []*template.Template, cnf GenerateConfig) error {
	logger.DebugContext(ctx, "generating", slog.Int("packages", len(pkgs)), slog.Int("templates", len(tmps)))

	for _, p := range pkgs {
		for _, tmp := range tmps {
			logger.DebugContext(ctx, "generating", slog.String("package", p.Name), slog.String("dir", p.Dir), slog.String("template", tmp.Name()))
			if err := GenerateInPackage(ctx, p, tmp, cnf); err != nil {
				logger.ErrorContext(ctx, "error while rendering file", slog.String("package", p.Name), slog.String("dir", p.Dir), slog.String("template", tmp.Name()))
				return err
			}
		}
	}

	return nil
}
