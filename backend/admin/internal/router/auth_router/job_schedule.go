package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerJobScheduleRouters(router fiber.Router) {
	jobScheduleHandler := logic.JobScheduleHandler{}
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("job_schedule"))

	router.Post("/list", handler.CtxHandlerFunc(jobScheduleHandler.List))
	router.Post("/detail", handler.CtxHandlerFunc(jobScheduleHandler.Detail))
	router.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Create))
	router.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Del))
	router.Post("/switch", logMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Switch))
	router.Post("/sync", logMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Sync))
	router.Post("/trigger", logMiddleware, handler.CtxHandlerNilFunc(jobScheduleHandler.Trigger))
}
