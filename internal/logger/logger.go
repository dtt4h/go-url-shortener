package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New(w io.Writer, level slog.Level) *Logger {
	var handler slog.Handler
	if w == nil {
		w = os.Stdout
	}

	handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})

	return &Logger{
		logger: slog.New(handler),
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
	}
}

type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

func ParseLevel(s string) (Level, error) {
	switch s {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	default:
		return LevelInfo, nil
	}
}

func Default() *Logger {
	level := LevelInfo
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		if l, err := ParseLevel(v); err == nil {
			level = l
		}
	}
	return New(os.Stdout, level)
}
