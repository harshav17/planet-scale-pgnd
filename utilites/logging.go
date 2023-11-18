package utilities

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func GetLogger() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, nil))
}
