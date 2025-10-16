package main

import (
	"context"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/ifnotnil/pkgen"
	"golang.org/x/tools/go/packages"
)

type PKGen struct {
	pk Packages
	tm Templates
	gn Generator
}

type Templates interface {
	Get(name string) (*template.Template, error)
	GetAll(c pkgen.TemplateConfigs) ([]*template.Template, error)
}

type Packages interface {
	Query(ctx context.Context, q pkgen.PackagesQueryConfig) ([]packages.Package, error)
}

type Generator interface {
	Generate(ctx context.Context, logger *slog.Logger, pkgs []packages.Package, tmps []*template.Template, cnf pkgen.GenerateConfig) error
}

func (p *PKGen) Run(ctx context.Context, cnf pkgen.Config) error {
	logger := slog.Default()

	logger.DebugContext(ctx, "config", slog.Any("config", cnf), slog.String("runnint_mode", pkgen.GetRunningMode().String()))

	// package query
	packages, err := p.pk.Query(ctx, cnf.PackagesQuery)
	if err != nil {
		logger.ErrorContext(ctx, "error while querying packages", errAttr(err))
		return err
	}
	debugLogPackages(ctx, packages)

	// templates
	tmps, err := p.tm.GetAll(cnf.Templates)
	if err != nil {
		logger.ErrorContext(ctx, "error while processing templates", errAttr(err))
		return err
	}

	// generate file
	err = p.gn.Generate(ctx, logger, packages, tmps, cnf.Generate)
	if err != nil {
		logger.ErrorContext(ctx, "error while generating files", errAttr(err))
		return err
	}

	return nil
}

func debugLogPackages(ctx context.Context, packages []packages.Package) {
	logger := slog.Default()
	if !logger.Enabled(ctx, slog.LevelDebug) {
		return
	}

	for _, p := range packages {
		logger.DebugContext(ctx, "queried package", slog.String("package", fmt.Sprintf("%#v", p)))
	}
}

func errAttr(err error) slog.Attr {
	return slog.String("err", err.Error())
}
