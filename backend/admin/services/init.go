package services

import (
	"admin/config"
	"admin/services/orm"
	"admin/services/redisc"
)

func New(conf *config.Config) {
	repoService := NewRepo(conf.Repo)
	redisService := NewRedis(conf.Redis)
	conf.Fiber.Services = append(conf.Fiber.Services, NewHttpc(), repoService, redisService, NewAuth(conf.Auth, redisc.Client), NewGeo(), NewAsynq(conf.Redis), NewCasbin(orm.Client.DB))
}
