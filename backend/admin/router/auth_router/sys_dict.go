package auth_router

import (
	"admin/fiberc/handler"
	"admin/fiberc/middleware"
	"admin/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerSysDictRouters(router fiber.Router) {
	sysDictHandler := logic.SysDictHandler{}

	dictType := router.Group("/type")
	dictType.Post("/list", handler.CtxHandlerFunc(sysDictHandler.TypeList))
	logMiddleware := middleware.OperationLogMiddleware(middleware.WithModule("dict"))
	dictType.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.TypeCreate))
	dictType.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.TypeUpdate))
	dictType.Post("/switch", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.TypeSwitch))
	dictType.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.TypeDel))

	dictEntry := router.Group("/entry")
	dictEntry.Post("/list", handler.CtxHandlerFunc(sysDictHandler.EntryList))
	dictEntry.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryCreate))
	dictEntry.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryUpdate))
	dictEntry.Post("/switch", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntrySwitch))
	dictEntry.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryDel))
	dictEntry.Post("/batch/copy", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryBatchCopy))
}
