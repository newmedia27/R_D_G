package log

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"project/internal/logger"
)

func Middleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		start := time.Now()
		reqId := ctx.Get("X-Request-Id")
		if reqId == "" {
			reqId = uuid.New().String()
		}
		ctx.SetContext(logger.WithAttrs(ctx.Context(), slog.String("request_id", reqId)))

		logger.LogInfoContext(
			ctx.Context(),
			"request start",
			slog.String("method", ctx.Method()),
			slog.String("path", ctx.Path()),
			slog.String("startedAt", start.Format(time.DateTime)),
		)
		err := ctx.Next()

		if err != nil {
			logger.LogErrorContext(ctx.Context(), err, slog.String("error", err.Error()))
		} else {
			duration := time.Since(start)
			s := []interface{}{
				slog.String("duration", duration.String()),
				slog.String("path", ctx.Path()),
				slog.Int("status", ctx.Response().StatusCode()),
			}
			id := ctx.Locals("userId")
			if id != nil {
				s = append(s, slog.String("user_id", id.(string)))
			}

			logger.LogInfoContext(
				ctx.Context(),
				"request end",
				s...,
			)

		}
		return err
	}
}
