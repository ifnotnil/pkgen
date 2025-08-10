package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ifnotnil/pkgen"
	"github.com/ifnotnil/pkgen/templates"
)

func main() {
	ctx := context.Background()

	packages, err := pkgen.Packages(ctx, pkgen.PackagesQueryConfig{})
	if err != nil {
		fmt.Printf("error while querying packages: %s", err.Error())
		os.Exit(1)
	}

	tmp, err := templates.Get("pkgpath")
	if err != nil {
		fmt.Printf("error while querying packages: %s", err.Error())
		os.Exit(1)
	}

	for _, p := range packages {
		err = pkgen.GenerateInPackage(ctx, p, tmp, pkgen.GenerateConfig{})
		if err != nil {
			fmt.Printf("error while generating: %s", err.Error())
			os.Exit(1)
		}
	}
}
