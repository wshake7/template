package router

import (
	"admin/fiberc/middleware"
	"github.com/gofiber/fiber/v3"
)

type Router struct{}

func (r *Router) RegisterRouters(group fiber.Router) {
	group = group.Group("/api")
	defaultGroup := group.Use(
		middleware.TimestampMiddleware(),
		//middleware.NonceMiddleware()
	)
	registerAccountRouters(defaultGroup.Group("account"))
	registerEncryptRouters(defaultGroup.Group("encrypt"))
}
