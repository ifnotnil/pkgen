package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ifnotnil/pkgen"
	"github.com/ifnotnil/pkgen/templates"
	"go.yaml.in/yaml/v4"
)

func main() {
	ctx := context.Background()

	loggerLevel := &slog.LevelVar{}
	loggerLevel.Set(slog.LevelInfo)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: false, Level: loggerLevel}))
	slog.SetDefault(logger)

	var (
		config         string
		template       string
		templateCustom string
		verbose        bool
	)

	flag.StringVar(&config, "config", ".pkgen.yml", "configuration file to use")
	flag.StringVar(&config, "c", ".pkgen.yml", "configuration file to use")
	flag.StringVar(&template, "template", "", "template to generate")
	flag.StringVar(&templateCustom, "template-custom", "", "template to generate")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.Parse()

	if verbose {
		loggerLevel.Set(slog.LevelDebug)
	}

	logger.DebugContext(ctx, "cli arguments", slog.String("config", config), slog.String("template", template), slog.String("template-custom", templateCustom))

	cnf := pkgen.DefaultConfig

	if runningInsideGoGenerate() {
		// if it is running inside go:generate query only the local package.
		cnf.PackagesQuery.Patterns = []string{"."}
		logger.DebugContext(ctx, "running inside a go:generate")
	} else {
		err := parseConfig(ctx, config, &cnf)
		if err != nil {
			logger.ErrorContext(ctx, "error while parsing config", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}

	// cli arg overwrites config's one
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

	logger.DebugContext(ctx, "packages", slog.Any("packages", packages))

	tmps, err := templates.GetTemplates(cnf.Templates)
	if err != nil {
		logger.ErrorContext(ctx, "error while getting templates", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.DebugContext(ctx, "packages", slog.Any("packages", packages))

	logger.DebugContext(ctx, "generating", slog.Int("packages", len(packages)), slog.Int("templates", len(tmps)))
	for _, p := range packages {
		for _, tmp := range tmps {
			logger.DebugContext(ctx, "generating", slog.String("package", p.Name), slog.String("template", tmp.Name()))
			err = pkgen.GenerateInPackage(ctx, p, tmp, pkgen.GenerateConfig{})
			if err != nil {
				fmt.Printf("error while generating: %s", err.Error())
				os.Exit(1)
			}
		}
	}
}

func runningInsideGoGenerate() bool {
	_, exists := os.LookupEnv("GOFILE")

	return exists
}

func parseConfig(ctx context.Context, filePath string, cnf *pkgen.Config) error {
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
