package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ifnotnil/pkgen"
	"github.com/ifnotnil/pkgen/templates"
	"sigs.k8s.io/yaml"
)

func main() {
	ctx := context.Background()

	var (
		config         string
		template       string
		templateCustom string
	)

	flag.StringVar(&config, "config", ".pkgen.yml", "configuration file to use")
	flag.StringVar(&config, "c", ".pkgen.yml", "configuration file to use")
	flag.StringVar(&template, "template", "", "template to generate")
	flag.StringVar(&templateCustom, "template-custom", "", "template to generate")
	flag.Parse()

	cnf := pkgen.Config{}

	if runningInsideGoGenerate() {
		// if it is running inside go:generate query only the local package.
		cnf.PackagesQuery.Patterns = []string{"."}
	} else {
		err := parseConfig(config, &cnf)
		if err != nil {
			fmt.Printf("error while parsing config: %s", err.Error())
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

	packages, err := pkgen.Packages(ctx, cnf.PackagesQuery)
	if err != nil {
		fmt.Printf("error while querying packages: %s", err.Error())
		os.Exit(1)
	}

	tmps, err := templates.GetTemplates(cnf.Templates)
	if err != nil {
		fmt.Printf("error while querying packages: %s", err.Error())
		os.Exit(1)
	}

	for _, p := range packages {
		for _, tmp := range tmps {
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

func parseConfig(filePath string, cnf *pkgen.Config) error {
	filePath = filepath.Clean(filePath)
	if fileExists(filePath) {
		if b, err := os.ReadFile(filePath); err == nil {
			err = yaml.Unmarshal(b, &cnf)
			if err != nil {
				fmt.Printf("error while parsing config yaml: %s", err.Error())
				os.Exit(1)
			}
		}
	}

	return nil
}
