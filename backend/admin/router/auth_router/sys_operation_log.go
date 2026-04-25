package auth_router

import (
	"admin/fiberc/handler"
	"admin/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysOperationLogRouters(router fiber.Router) {
	sysOperationLogHandler := logic.SysOperationLogHandler{}

	router.Post("/list", handler.CtxHandlerFunc(sysOperationLogHandler.List))
	router.Post("/detail", handler.CtxHandlerFunc(sysOperationLogHandler.Detail))
}
