package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
)

const (
	dev  string = "dev"
	prod string = "prod"
)

// New TODO возможно неэффективно
func New(environment string, filepath string) *slog.Logger {
	var logger *slog.Logger

	optsProd := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	optsDev := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	writer := &lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	if environment == prod {
		logger = slog.New(slog.NewJSONHandler(writer, optsProd))
	}

	if environment == dev {
		logger = slog.New(slog.NewTextHandler(os.Stdout, optsDev))
	}

	return logger
}
