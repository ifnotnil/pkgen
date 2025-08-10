package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ifnotnil/pkgen"
	"github.com/ifnotnil/pkgen/templates"
	"golang.org/x/tools/go/packages"
)

func main() {
	ctx := context.Background()

	loggerLevel := &slog.LevelVar{}
	loggerLevel.Set(slog.LevelInfo)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: false, Level: loggerLevel}))
	slog.SetDefault(logger)

	var (
		configPath     string
		template       string
		templateCustom string
		verbose        bool
	)

	flag.StringVar(&configPath, "config", ".pkgen.yml", "configuration file to use")
	flag.StringVar(&configPath, "c", ".pkgen.yml", "configuration file to use")
	flag.StringVar(&template, "template", "", "template to generate")
	flag.StringVar(&templateCustom, "template-custom", "", "template to generate")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.Parse()

	if verbose {
		loggerLevel.Set(slog.LevelDebug)
	}

	logger.DebugContext(ctx, "cli arguments", slog.String("config", configPath), slog.String("template", template), slog.String("template-custom", templateCustom))

	cnf, err := pkgen.ParseConfig(ctx, configPath)
	if err != nil {
		logger.ErrorContext(ctx, "error while parsing config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// cli arguments overwrite config's one
	if template != "" {
		cnf.Templates = []templates.TemplateConfig{{Name: template}}
	}
	if templateCustom != "" {
		cnf.Templates = append(cnf.Templates, templates.TemplateConfig{CustomTemplateFile: templateCustom})
	}

	logger.DebugContext(ctx, "config", slog.Any("config", cnf))

	packages, err := pkgen.Packages(ctx, cnf.PackagesQuery)
	if err != nil {
		logger.ErrorContext(ctx, "error while querying packages", slog.String("error", err.Error()))
		os.Exit(1)
	}
	debugLogPackages(ctx, packages)

	tmps, err := templates.GetTemplates(cnf.Templates)
	if err != nil {
		logger.ErrorContext(ctx, "error while getting templates", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.DebugContext(ctx, "generating", slog.Int("packages", len(packages)), slog.Int("templates", len(tmps)))
	for _, p := range packages {
		for _, tmp := range tmps {
			logger.DebugContext(ctx, "generating", slog.String("package", p.Name), slog.String("dir", p.Dir), slog.String("template", tmp.Name()))
			err = pkgen.GenerateInPackage(ctx, p, tmp, cnf.Generate)
			if err != nil {
				fmt.Printf("error while generating: %s", err.Error())
				os.Exit(1)
			}
		}
	}
}

func debugLogPackages(ctx context.Context, packages []*packages.Package) {
	logger := slog.Default()
	if !logger.Enabled(ctx, slog.LevelDebug) {
		return
	}

	for _, p := range packages {
		logger.DebugContext(ctx, "queried package", slog.String("package", fmt.Sprintf("%#v", p)))
	}
}
