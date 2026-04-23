package middleware

import (
	"admin/domains"
	"admin/fiberc/handler"
	"admin/fiberc/res"
	redisc2 "admin/services/redisc"
	"go-common/utils/encrypt"
	"go-common/utils/encrypt/aes_util"
	"go-common/utils/encrypt/rsa_util"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func TimestampMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		timestamp := fiber.GetReqHeader[int64](ctx, domains.HeaderXRequestTimestamp)
		if timestamp != 0 {
			now := time.Now().UnixMilli()
			if abs(now-timestamp) > encrypt.REQUEST_EXPIRE_TIME.Milliseconds() {
				return res.FailRequestExpired
			}
		}
		return ctx.Next()
	})
}

func NonceMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		if redisc2.Client != nil {
			nonce := ctx.RequestID()
			if nonce != "" {
				nonceKey := encrypt.NONCE_REDIS_KEY_PREFIX.Sprintf(nonce)

				err := redisc2.Client.Do(ctx.Context(), redisc2.Client.B().Set().Key(nonceKey).Value(nonce).Ex(encrypt.NONCE_EXPIRE_TIME).Build()).Error()
				if err != nil {
					zap.L().Error("redis error", zap.Error(err))
					return res.FailDefault
				}
			}
		}
		return ctx.Next()
	})
}

func EncryptMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		ctx.Set(domains.HeaderXResponseIsEncrypt, "false")
		privateKey := ctx.PrivateKey
		if privateKey == "" {
			ctx.L().Error("privateKey is empty")
			return res.FailRequest
		}
		encryptedKey := fiber.GetReqHeader[string](ctx, domains.HeaderXRequestEncryptedKey)
		if encryptedKey == "" {
			ctx.L().Error("encrypted key is empty")
			return res.FailRequest
		}
		aesKey, err := rsa_util.Decrypt(encryptedKey, privateKey)
		if err != nil {
			ctx.L().Error("encrypt error", zap.Error(err))
			return res.FailDefault
		}
		timestamp := fiber.GetReqHeader[int64](ctx, domains.HeaderXRequestTimestamp)
		nonce := ctx.RequestID()
		params := map[string]any{
			domains.HeaderXRequestID:        nonce,
			domains.HeaderXRequestTimestamp: timestamp,
		}

		for k, v := range ctx.Queries() {
			params[k] = v
		}
		add := encrypt.UriSort(params, func(key string) bool {
			return true
		})
		reqBody := ctx.Request().Body()
		sign := fiber.GetReqHeader[string](ctx, domains.HeaderXRequestSignature)
		var decrypt []byte
		if reqBody == nil {
			decrypt, err = aes_util.DecryptCiphertextAndTag("", sign, aesKey, add)
			if err != nil {
				ctx.L().Error("decrypt error", zap.Error(err))
				return res.FailDefault
			}
		} else {
			decrypt, err = aes_util.DecryptCiphertextAndTag(string(reqBody), sign, aesKey, add)
			if err != nil {
				ctx.L().Error("decrypt error", zap.Error(err))
				return res.FailDefault
			}
			ctx.Request().SetBody(decrypt)
		}
		err = ctx.Next()
		if err != nil {
			return err
		}
		resBody := ctx.Response().Body()
		resBodyStr := string(resBody)
		ctx.AddResLogFields(zap.String(domains.LogFieldDecryptResBody, resBodyStr))
		result, err := aes_util.Encrypt(resBodyStr, aesKey, "")
		if err != nil {
			ctx.L().Error("encrypt error", zap.Error(err))
			return res.FailDefault
		}
		ctx.Set(domains.HeaderXResponseIsEncrypt, "true")
		ctx.Response().SetBody([]byte(result.Combined))
		return err
	})
}

const SigData = "signData"

func SignMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		privateKey := ctx.PrivateKey
		if privateKey == "" {
			ctx.L().Error("privateKey is empty")
			return res.FailRequest
		}
		encryptedKey := fiber.GetReqHeader[string](ctx, domains.HeaderXRequestEncryptedKey)
		if encryptedKey == "" {
			ctx.L().Error("encrypted key is empty")
			return res.FailRequest
		}
		aesKey, err := rsa_util.Decrypt(encryptedKey, privateKey)
		if err != nil {
			ctx.L().Error("encrypt error", zap.Error(err))
			return res.FailDefault
		}
		timestamp := fiber.GetReqHeader[int64](ctx, domains.HeaderXRequestTimestamp)
		nonce := ctx.RequestID()
		params := map[string]any{
			domains.HeaderXRequestID:        nonce,
			domains.HeaderXRequestTimestamp: timestamp,
		}

		for k, v := range ctx.Queries() {
			params[k] = v
		}
		reqBody := ctx.Request().Body()

		params[SigData] = reqBody
		add := encrypt.UriSort(params, func(key string) bool {
			return true
		})
		sign := fiber.GetReqHeader[string](ctx, domains.HeaderXRequestSignature)
		_, err = aes_util.DecryptCiphertextAndTag("", sign, aesKey, add)
		if err != nil {
			ctx.L().Error("decrypt error", zap.Error(err))
			return res.FailDefault
		}
		return ctx.Next()
	})
}
