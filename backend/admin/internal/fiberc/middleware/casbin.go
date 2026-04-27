package middleware

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/casbin"
	"fmt"
	"github.com/gofiber/fiber/v3"
)

func CasbinAPIMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		info := ctx.SessionInfo

		roles := info.Roles

		// 构建所有 subject（逐个尝试，any match = allow）
		subjects := []string{fmt.Sprintf("user:%d", info.Id)}
		for _, r := range roles {
			subjects = append(subjects, "role:"+r)
		}

		obj := ctx.Path()    // "/api/sys/dict/list"
		act := ctx.Matched() // "GET"

		for _, sub := range subjects {
			ok, _ := casbin.E.Enforce(sub, obj, act)
			if ok {
				return ctx.Next()
			}
		}
		return res.FailAccountDisabled
	})
}
