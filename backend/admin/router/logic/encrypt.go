package logic

import (
	"admin/fiberc/handler"
	"admin/fiberc/res"
	"admin/services/redisc"
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

func (r *EncryptHandler) PublicKey(ctx *handler.Ctx) (*ResPublicKey, error) {
	var keyPair redisc.DtoKeyPair
	err := redisc.Client.GetJson(ctx, redisc.KeyGlobalEncryptPublicKey, &keyPair)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			keyPair.PrivateKey, keyPair.PublicKey, err = rsa_util.GenerateKeyPair()
			if err != nil {
				ctx.L().Error("生成rsaKey错误", zap.Error(err))
				return nil, res.FailDefault
			}
			err = redisc.Client.Do(ctx, redisc.Client.B().Set().Key(redisc.KeyGlobalEncryptPublicKey).Value(rueidis.JSON(keyPair)).Build()).Error()
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
