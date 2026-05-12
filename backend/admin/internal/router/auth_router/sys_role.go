package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysRoleRouters(router fiber.Router) {
	sysRoleHandler := logic.SysRoleHandler{}
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("role"))

	router.Post("/list", handler.CtxHandlerFunc(sysRoleHandler.List))
	router.Get("/tree", handler.CtxFunc(sysRoleHandler.Tree))
	router.Get("/:id/permissions", handler.CtxHandlerFunc(sysRoleHandler.Permissions))
	router.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.Create))
	router.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.Del))
	router.Post("/permissions", logMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.SavePermissions))
}
