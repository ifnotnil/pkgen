package pkgen

import (
	"log/slog"
	"os"
	"testing"
)

func logger(t *testing.T) *slog.Logger {
	t.Helper()

	if testing.Verbose() {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelDebug,
			AddSource:   false,
			ReplaceAttr: nil,
		}))
	}

	return slog.New(slog.DiscardHandler)
}
