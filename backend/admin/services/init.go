package services

import (
	"admin/config"
	"admin/services/orm"
	"admin/services/redisc"
)

func New(conf *config.Config) {
	ormService := NewOrm(conf.Orm)
	redisService := NewRedis(conf.Redis)
	conf.Fiber.Services = append(conf.Fiber.Services, NewHttpc(), ormService, redisService, NewAuth(conf.Auth, redisc.Client), NewGeo(), NewAsynq(conf.Redis), NewCasbin(orm.Client.DB))
}
