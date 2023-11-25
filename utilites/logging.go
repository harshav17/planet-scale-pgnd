package utilities

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func GetLogger() *slog.Logger {
	// TODO - use JSON and separate pretty printing for local dev
	return slog.New(tint.NewHandler(os.Stdout, nil))
}
