package router

import (
	"admin/fiberc/handler"
	"admin/fiberc/middleware"
	"admin/router/logic"
	"github.com/gofiber/fiber/v3"
)

func registerAccountRouters(router fiber.Router) {
	accountHandler := logic.AccountHandler{}
	router.Post("/login/pwd", middleware.PublicMiddleware(), middleware.EncryptMiddleware(), handler.CtxHandlerFunc(func(ctx *handler.Ctx, req *logic.ReqAccountPwdLogin) (*logic.ResAccountPwdLogin, error) {
		res, err := accountHandler.PwdLogin(ctx, req)
		if err != nil {
			return res, err
		}
		ctx.Cookie(&fiber.Cookie{
			Name:     ctx.Config.Auth.TokenName,
			Value:    res.Token,
			SameSite: "Lax",
			Secure:   false,
			HTTPOnly: true,
		})
		return res, err
	}))
	router.Post("/changePwd", middleware.AuthMiddleware(), middleware.EncryptMiddleware(), handler.CtxHandlerNilFunc(accountHandler.ChangePwd))
	router.Get("/logout", middleware.AuthMiddleware(), handler.CtxHandlerNilFunc(accountHandler.Logout))
}
