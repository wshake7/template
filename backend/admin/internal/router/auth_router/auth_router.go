package auth_router

import (
	middleware2 "admin/internal/fiberc/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterRouters(router fiber.Router) {
	group := router.Use(middleware2.AuthMiddleware(), middleware2.EncryptMiddleware())
	registerSysRoleRouters(group.Group("/sys/role"))
	registerSysDictRouters(group.Group("/sys/dict"))
	registerSysOperationLogRouters(group.Group("/sys/operation/log"))
}
