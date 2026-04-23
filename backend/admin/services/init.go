package services

import (
	"admin/config"
	"admin/services/redisc"
	"admin/services/repo"
)

func New(conf *config.Config) {
	repoService := NewRepo(conf.Repo)
	redisService := NewRedis(conf.Redis)
	conf.Fiber.Services = append(conf.Fiber.Services, NewHttpc(), repoService, redisService, NewAuth(conf.Auth, redisc.Client), NewGeo(), NewAsynq(conf.Redis), NewCasbin(repo.Client.DB))
}
