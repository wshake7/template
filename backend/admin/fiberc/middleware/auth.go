package middleware

import (
	"admin/auth"
	"admin/config"
	"admin/domains"
	"admin/fiberc/handler"
	"admin/fiberc/res"
	"admin/services/redisc"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func AuthMiddleware() fiber.Handler {
	key := config.Conf.Auth.TokenName
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		token := ctx.Cookies(key)
		if token == "" {
			return res.FailNotLogin
		}
		session, err := auth.GetSessionByToken(token)
		if err != nil {
			ctx.L().Error("get session error", zap.Error(err))
			return auth.CheckLoginErr(err)
		}
		info, err := session.GetInfo()
		if err != nil || info.PrivateKey == "" {
			ctx.L().Error("get session info error", zap.Error(err), zap.String("key", info.PrivateKey))
			return res.FailDefault
		}
		ctx.PrivateKey = info.PrivateKey
		ctx.SessionInfo = &info
		ctx.AddResLogFields(zap.Any(domains.LogFieldSessionInfo, info))
		ctx.AddLogFields(zap.Any(domains.LogFieldSessionInfo, info))
		return ctx.Next()
	})
}

func PublicMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		var keyPair redisc.DtoKeyPair
		err := redisc.Client.GetJson(ctx, redisc.KeyGlobalEncryptPublicKey, &keyPair)
		if err != nil || keyPair.PrivateKey == "" {
			ctx.L().Error("get key error", zap.Error(err), zap.String("key", keyPair.PrivateKey))
			return res.FailRequestKey
		}
		ctx.PrivateKey = keyPair.PrivateKey
		return ctx.Next()
	})
}
