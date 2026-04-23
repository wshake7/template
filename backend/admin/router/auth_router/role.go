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
	router.Post("/create", middleware.OperationLogMiddleware(middleware.WithModule("role")), handler.CtxHandlerNilFunc(roleHandler.Create))
	router.Post("/update", middleware.OperationLogMiddleware(middleware.WithModule("role")), handler.CtxHandlerNilFunc(roleHandler.Update))
	router.Post("/switchStatus", middleware.OperationLogMiddleware(middleware.WithModule("role")), handler.CtxHandlerNilFunc(roleHandler.SwitchStatus))
	router.Post("/delete", middleware.OperationLogMiddleware(middleware.WithModule("role")), handler.CtxHandlerNilFunc(roleHandler.Delete))
}
