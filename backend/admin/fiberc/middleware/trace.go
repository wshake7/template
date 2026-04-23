package middleware

import (
	"admin/fiberc/handler"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	TraceIDKey string = "trace_id"
)

func TraceMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		traceId := ctx.Get("X-Trace-Id")
		if traceId == "" {
			traceId = uuid.NewString()
		}
		ctx.Set("X-Trace-Id", traceId)
		ctx.TraceId = traceId
		ctx.AddLogFields(zap.String(TraceIDKey, traceId), zap.String("method", ctx.Method()), zap.String("path", ctx.Path()))
		ctx.AddResLogFields(zap.String(TraceIDKey, traceId))
		return ctx.Next()
	})
}
