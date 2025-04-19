package logger

import (
	"log/slog"
	"os"
	"person-project/config"
)

var Log *slog.Logger

func Init() {
	// Проверяем, что конфиг загружен
	if config.Cfg == nil {
		panic("Конфиг не загружен! Вызовите config.LoadConfig() перед logger.Init()")
	}

	// Уровень логирования
	var level slog.Level
	switch config.Cfg.App.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	if config.Cfg.App.Env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}
