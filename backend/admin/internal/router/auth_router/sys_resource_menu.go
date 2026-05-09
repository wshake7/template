package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysResourceMenuRouters(router fiber.Router) {
	sysResourceMenuHandler := logic.SysResourceMenuHandler{}
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("resource_menu"))

	router.Post("/list", handler.CtxHandlerFunc(sysResourceMenuHandler.List))
	router.Get("/tree", handler.CtxFunc(sysResourceMenuHandler.Tree))
	router.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysResourceMenuHandler.Create))
	router.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysResourceMenuHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysResourceMenuHandler.Del))
}
