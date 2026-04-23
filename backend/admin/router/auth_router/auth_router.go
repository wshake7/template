package auth_router

import (
	"admin/fiberc/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterRouters(router fiber.Router) {
	group := router.Use(middleware.AuthMiddleware(), middleware.EncryptMiddleware())
	registerRoleRouters(group)
}
