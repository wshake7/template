package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerJobExecutionRouters(router fiber.Router) {
	jobExecutionHandler := logic.JobExecutionHandler{}
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("job_execution"))

	router.Post("/list", handler.CtxHandlerFunc(jobExecutionHandler.List))
	router.Post("/detail", handler.CtxHandlerFunc(jobExecutionHandler.Detail))
	router.Post("/cancel", logMiddleware, handler.CtxHandlerNilFunc(jobExecutionHandler.Cancel))
	router.Post("/retry", logMiddleware, handler.CtxHandlerNilFunc(jobExecutionHandler.Retry))
}
