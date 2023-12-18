package utilities

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type ENV string

const (
	Development ENV = "development"
	Production  ENV = "production"
)

func GetLogger() *slog.Logger {
	env, ok := os.LookupEnv("APP_ENV")
	if !ok {
		env = "development"
	}

	switch env {
	case "production":
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	default:
		return slog.New(tint.NewHandler(os.Stdout, nil))
	}
}
