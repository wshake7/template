package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysDictRouters(router fiber.Router) {
	sysDictHandler := logic.SysDictHandler{}

	dictType := router.Group("/type")
	dictType.Post("/list", handler.CtxHandlerFunc(sysDictHandler.TypeList))
	logMiddleware := middleware.OperationLogMiddleware(middleware.WithModule("dict"))
	dictType.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.TypeCreate))
	dictType.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.TypeUpdate))
	dictType.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.TypeDel))

	dictEntry := router.Group("/entry")
	dictEntry.Post("/list", handler.CtxHandlerFunc(sysDictHandler.EntryList))
	dictEntry.Post("/match", handler.CtxHandlerFunc(sysDictHandler.EntryMatch))
	dictEntry.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryCreate))
	dictEntry.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryUpdate))
	dictEntry.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryDel))
	dictEntry.Post("/batch/copy", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryBatchCopy))
}
