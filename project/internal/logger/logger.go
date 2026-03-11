package logger

import (
	"context"
	"log/slog"
	"os"

	"project/internal/config"
)

type CtxLogKey struct{}
type Logger struct {
	*slog.JSONHandler
}

func WithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	return context.WithValue(ctx, CtxLogKey{}, attrs)
}

func (l *Logger) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(getAttrs(ctx)...)
	return l.JSONHandler.Handle(ctx, record)
}

func getAttrs(ctx context.Context) []slog.Attr {
	attrs := ctx.Value(CtxLogKey{})
	if attrs == nil {
		return []slog.Attr{}
	}

	res, ok := attrs.([]slog.Attr)
	if !ok {
		return []slog.Attr{}
	}
	return res
}

func InitLogger(cfg *config.Config) {
	handler := &Logger{
		JSONHandler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: getLogLevel(cfg.LogLevel),
		}),
	}
	slog.SetDefault(slog.New(handler))
}

func LogWarnContext(ctx context.Context, msg string, args ...any) {
	slog.Default().WarnContext(ctx, msg, args...)
}

func LogInfoContext(ctx context.Context, msg string, args ...any) {
	slog.Default().InfoContext(ctx, msg, args...)
}

func LogErrorContext(ctx context.Context, err error, args ...any) {
	slog.Default().ErrorContext(ctx, err.Error(), args...)
}

func LogFatalContext(ctx context.Context, err error, args ...any) {
	LogErrorContext(ctx, err, args...)
	os.Exit(1)
}

func LogFatalContextWithCode(ctx context.Context, err error, code int, args ...any) {
	LogErrorContext(ctx, err, args...)
	os.Exit(code)
}

func LogPanicContext(ctx context.Context, err error, args ...any) {
	LogErrorContext(ctx, err, args...)
	panic(err)
}

func getLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
