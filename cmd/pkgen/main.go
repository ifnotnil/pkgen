package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/ifnotnil/pkgen"
	"github.com/lmittmann/tint"
	"golang.org/x/term"
)

func slogHandler(loggerLevel *slog.LevelVar) slog.Handler {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		return tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:   false,
			Level:       loggerLevel,
			ReplaceAttr: nil,
			TimeFormat:  time.Kitchen,
			NoColor:     false,
		})
	}

	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       loggerLevel,
		ReplaceAttr: nil,
	})
}

func main() {
	ctx := context.Background()

	loggerLevel := &slog.LevelVar{}
	loggerLevel.Set(slog.LevelInfo)
	slog.SetDefault(slog.New(slogHandler(loggerLevel)))
	logger := slog.Default()

	// config
	cnf, err := pkgen.NewConfig(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "error while parsing config", errAttr(err))
		os.Exit(1)
	}

	p := PKGen{
		pk: pkgen.Packages{},
		tm: pkgen.Templates{},
		gn: pkgen.Generator{
			FileWriter: nil,
		},
	}

	if err := p.Run(ctx, cnf); err != nil {
		os.Exit(1)
	}
}
