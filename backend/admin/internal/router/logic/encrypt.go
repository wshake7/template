package logic

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	redisc2 "admin/internal/services/redisc"
	"errors"
	"go-common/utils/encrypt/rsa_util"

	"github.com/redis/go-redis/v9"
	"github.com/redis/rueidis"
	"go.uber.org/zap"
)

type EncryptHandler struct{}

type ResPublicKey struct {
	PublicKey string `json:"publicKey"`
}

// @Summary 获取加密公钥
// @Description 获取用于敏感数据加密的 RSA 公钥
// @Tags Encrypt
// @Produce json
// @Success 200 {object} res.Response{data=ResPublicKey} "成功"
// @Router /api/encrypt/public/key [get]
func (r *EncryptHandler) PublicKey(ctx *handler.Ctx) (*ResPublicKey, error) {
	var keyPair redisc2.DtoKeyPair
	err := redisc2.Client.GetJson(ctx, redisc2.KeyGlobalEncryptPublicKey, &keyPair)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			keyPair.PrivateKey, keyPair.PublicKey, err = rsa_util.GenerateKeyPair()
			if err != nil {
				ctx.L().Error("生成rsaKey错误", zap.Error(err))
				return nil, res.FailDefault
			}
			err = redisc2.Client.Do(ctx, redisc2.Client.B().Set().Key(redisc2.KeyGlobalEncryptPublicKey).Value(rueidis.JSON(keyPair)).Build()).Error()
			if err != nil {
				ctx.L().Error("保存rsaKey错误", zap.Error(err))
				return nil, res.FailDefault
			}
		} else {
			ctx.L().Error("获取全局Key错误", zap.Error(err))
			return nil, res.FailDefault
		}
	}
	return &ResPublicKey{PublicKey: keyPair.PublicKey}, nil
}
