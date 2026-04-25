package auth_router

import (
	"admin/fiberc/handler"
	"admin/fiberc/middleware"
	"admin/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerRoleRouters(router fiber.Router) {
	roleHandler := logic.RoleHandler{}
	router.Get("/list", handler.CtxHandlerFunc(roleHandler.List))
	logMiddleware := middleware.OperationLogMiddleware(middleware.WithModule("role"))
	router.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(roleHandler.Create))
	router.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(roleHandler.Update))
	router.Post("/switch", logMiddleware, handler.CtxHandlerNilFunc(roleHandler.Switch))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(roleHandler.Del))
}
