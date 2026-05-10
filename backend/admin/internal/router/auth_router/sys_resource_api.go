package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysResourceApiRouters(router fiber.Router) {
	sysResourceApiHandler := logic.SysResourceApiHandler{}
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("resource_api"))

	router.Post("/list", handler.CtxHandlerFunc(sysResourceApiHandler.List))
	router.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysResourceApiHandler.Create))
	router.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysResourceApiHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysResourceApiHandler.Del))
}
