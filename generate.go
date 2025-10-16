package pkgen

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/go/packages"
)

type FileWriter interface {
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type Generator struct {
	FileWriter FileWriter
}

func (g Generator) GenerateInPackage(ctx context.Context, pkg packages.Package, tmp *template.Template, cnf GenerateConfig) error {
	if len(pkg.GoFiles) == 0 {
		return nil
	}

	// execute the template
	buf := bytes.Buffer{}
	err := tmp.Execute(&buf, pkg)
	if err != nil {
		return err
	}

	// get output filename
	outFileName, err := generateName(OutputName{TemplateName: tmp.Name()}, cnf.OutputFile)
	if err != nil {
		return err
	}

	outPath := filepath.Join(filepath.Clean(pkg.Dir), outFileName)

	var wf func(name string, data []byte, perm os.FileMode) error

	if g.FileWriter != nil {
		wf = g.FileWriter.WriteFile
	} else {
		wf = os.WriteFile
	}

	return wf(outPath, buf.Bytes(), cnf.OutputFileMod)
}

func (g Generator) Generate(ctx context.Context, logger *slog.Logger, pkgs []packages.Package, tmps []*template.Template, cnf GenerateConfig) error {
	logger.DebugContext(ctx, "generating", slog.Int("packages", len(pkgs)), slog.Int("templates", len(tmps)))

	for _, p := range pkgs {
		for _, tmp := range tmps {
			logger.DebugContext(ctx, "generating", slog.String("package", p.Name), slog.String("dir", p.Dir), slog.String("template", tmp.Name()))
			if err := g.GenerateInPackage(ctx, p, tmp, cnf); err != nil {
				logger.ErrorContext(ctx, "error while rendering file", slog.String("package", p.Name), slog.String("dir", p.Dir), slog.String("template", tmp.Name()))
				return err
			}
		}
	}

	return nil
}

const defaultOutputNameTemplate = `zz_generated.{{ .TemplateName }}.go`

type OutputName struct {
	TemplateName string
}

func generateName(n OutputName, tmpl string) (string, error) {
	t, err := template.New("outputName_" + n.TemplateName).Parse(tmpl)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	err = t.Execute(&buf, n)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
