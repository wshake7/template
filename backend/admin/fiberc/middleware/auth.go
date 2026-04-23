package middleware

import (
	"admin/auth"
	"admin/config"
	"admin/domains"
	"admin/fiberc/handler"
	"admin/fiberc/res"
	"admin/services/redisc"
	"errors"
	"github.com/click33/sa-token-go/core/manager"
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
			switch {
			case errors.Is(err, manager.ErrAccountDisabled):
				return res.FailAccountDisabled
			case errors.Is(err, manager.ErrNotLogin):
				return res.FailNotLogin
			case errors.Is(err, manager.ErrTokenNotFound):
				return res.FailTokenNotFound
			case errors.Is(err, manager.ErrInvalidTokenData):
				return res.FailInvalidTokenData
			case errors.Is(err, manager.ErrLoginLimitExceeded):
				return res.FailLoginLimitExceeded
			case errors.Is(err, manager.ErrTokenKickout):
				return res.FailTokenKickout
			case errors.Is(err, manager.ErrTokenReplaced):
				return res.FailTokenReplaced
			default:
				ctx.L().Error("get session error", zap.Error(err))
				return res.FailDefault
			}
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
