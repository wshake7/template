package auth_router

import (
	"admin/fiberc/handler"
	"admin/fiberc/middleware"
	"admin/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerDictRouters(router fiber.Router) {
	dictHandler := logic.DictHandler{}

	dictType := router.Group("/type")
	dictType.Post("/list", handler.CtxHandlerFunc(dictHandler.TypeList))
	logMiddleware := middleware.OperationLogMiddleware(middleware.WithModule("dict"))
	dictType.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(dictHandler.TypeCreate))
	dictType.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(dictHandler.TypeUpdate))
	dictType.Post("/switch", logMiddleware, handler.CtxHandlerNilFunc(dictHandler.TypeSwitch))
	dictType.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(dictHandler.TypeDel))

	dictEntry := router.Group("/entry")
	dictEntry.Post("/list", handler.CtxHandlerFunc(dictHandler.EntryList))
	dictEntry.Post("/create", logMiddleware, handler.CtxHandlerNilFunc(dictHandler.EntryCreate))
	dictEntry.Post("/update", logMiddleware, handler.CtxHandlerNilFunc(dictHandler.EntryUpdate))
	dictEntry.Post("/switch", logMiddleware, handler.CtxHandlerNilFunc(dictHandler.EntrySwitch))
	dictEntry.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(dictHandler.EntryDel))
}
