package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysLoginLogRouters(router fiber.Router) {
	sysLoginLogHandler := logic.SysLoginLogHandler{}

	router.Post("/list", handler.CtxHandlerFunc(sysLoginLogHandler.List))
	router.Post("/detail", handler.CtxHandlerFunc(sysLoginLogHandler.Detail))
}
