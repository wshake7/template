package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysLanguageRouters(router fiber.Router) {
	sysLanguageHandler := logic.SysLanguageHandler{}

	langType := router.Group("/type")
	langType.Post("/list", handler.CtxHandlerFunc(sysLanguageHandler.TypeList))
	logMiddleware := middleware.OperationLogMiddleware(middleware.WithModule("language"))
	langType.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.TypeCreate))
	langType.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.TypeUpdate))
	langType.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.TypeDel))

	langEntry := router.Group("/entry")
	langEntry.Post("/list", handler.CtxHandlerFunc(sysLanguageHandler.EntryList))
	langEntry.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.EntryCreate))
	langEntry.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.EntryUpdate))
	langEntry.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.EntryDel))
	langEntry.Post("/batch/create", logMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.EntryBatchCreate))
}
