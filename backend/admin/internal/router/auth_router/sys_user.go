package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysUserRouters(router fiber.Router) {
	sysUserHandler := logic.SysUserHandler{}
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("user"))

	router.Post("/list", handler.CtxHandlerFunc(sysUserHandler.List))
	router.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysUserHandler.Create))
	router.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysUserHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysUserHandler.Del))
}
