package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

// initiating logger
func Init(level slog.Level) {
	Log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
}
