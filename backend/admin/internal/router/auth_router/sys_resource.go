package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysResourceRouters(router fiber.Router) {
	sysResourceHandler := logic.SysResourceHandler{}
	router.Post("/list", handler.CtxHandlerFunc(sysResourceHandler.List))
	logMiddleware := middleware.OperationLogMiddleware(middleware.WithModule("resource"))
	router.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysResourceHandler.Create))
	router.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysResourceHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysResourceHandler.Del))
}
