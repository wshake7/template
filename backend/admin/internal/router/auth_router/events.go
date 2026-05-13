package auth_router

import (
	"bufio"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func registerEventRouters(router fiber.Router) {
	router.Get("/events", func(ctx fiber.Ctx) error {
		ctx.Set("Content-Type", "text/event-stream")
		ctx.Set("Cache-Control", "no-cache")
		ctx.Set("Connection", "keep-alive")
		ctx.Set("Transfer-Encoding", "chunked")

		return ctx.SendStreamWriter(func(w *bufio.Writer) {
			ticker := time.NewTicker(time.Minute)
			defer ticker.Stop()
			count := 0
			for range ticker.C {
				count++
				if _, err := fmt.Fprintf(w, "event: count\ndata: {\"count\":%d}\n\n", count); err != nil {
					zap.L().Debug("write sse event failed", zap.Error(err))
					return
				}
				if err := w.Flush(); err != nil {
					zap.L().Debug("sse client disconnected", zap.Error(err))
					return
				}
			}
		})
	})
}
