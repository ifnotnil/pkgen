package pkgen

import (
	"bytes"
	"context"
	"fmt"
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

	dir := filepath.Dir(pkg.GoFiles[0])
	outPath := filepath.Join(dir, cnf.OutputFile)
	return os.WriteFile(outPath, buf.Bytes(), cnf.OutputFileMod)
}
