package auth_router

import (
	middleware2 "admin/internal/fiberc/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterRouters(router fiber.Router) {
	group := router.Use(middleware2.AuthMiddleware(), middleware2.EncryptMiddleware(), middleware2.LanguageMiddleware())
	registerSysRoleRouters(group.Group("/sys/role"))
	registerSysResourceRouters(group.Group("/sys/resource"))
	registerSysDictRouters(group.Group("/sys/dict"))
	registerSysLanguageRouters(group.Group("/sys/language"))
	registerSysOperationLogRouters(group.Group("/sys/operation/log"))
}
