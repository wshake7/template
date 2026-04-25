package auth_router

import (
	"admin/fiberc/handler"
	"admin/fiberc/middleware"
	"admin/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysRoleRouters(router fiber.Router) {
	sysRoleHandler := logic.SysRoleHandler{}
	router.Get("/list", handler.CtxHandlerFunc(sysRoleHandler.List))
	logMiddleware := middleware.OperationLogMiddleware(middleware.WithModule("role"))
	router.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.Create))
	router.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.Update))
	router.Post("/switch", logMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.Switch))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.Del))
}
