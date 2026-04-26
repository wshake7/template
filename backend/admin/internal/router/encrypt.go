package router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/router/logic"
	"github.com/gofiber/fiber/v3"
)

func registerEncryptRouters(router fiber.Router) {
	encryptHandler := logic.EncryptHandler{}
	router.Get("/public/key", handler.CtxFunc(encryptHandler.PublicKey))
}
