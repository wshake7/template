# Backend Router 编写风格

## 场景与目标

- 适用场景：新增或调整 `backend/admin/internal/router` 下的接口路由。
- 目标：保持路由分层、handler 包装器和 middleware 串联风格一致。

## 目录/文件位置

- 路由总注册：`backend/admin/internal/router/router.go`
- 公共路由：`backend/admin/internal/router/account.go`、`backend/admin/internal/router/encrypt.go`
- 鉴权路由：`backend/admin/internal/router/auth_router/*.go`
- 业务逻辑：`backend/admin/internal/router/logic/*.go`
- 中间件：`backend/admin/internal/fiberc/middleware/*.go`

## 当前分层

1. `router.go`：挂载 `/api`，应用 `TimestampMiddleware()`，注册公共路由和鉴权路由入口。
2. `internal/router/*.go`：公共接口映射，如登录、加密公钥。
3. `internal/router/auth_router/*.go`：鉴权接口映射，统一追加 `AuthMiddleware()` 和 `EncryptMiddleware()`。
4. `internal/router/logic/*.go`：请求结构体、业务逻辑、数据库调用、错误语义。
5. `handler.Ctx*` 包装器：负责绑定请求、注入上下文、统一输出。

## Handler 包装器选择

- `handler.CtxFunc(fn)`：无请求结构体，仅依赖上下文。
- `handler.CtxHandlerFunc(fn)`：有请求结构体，返回数据和错误。
- `handler.CtxHandlerNilFunc(fn)`：有请求结构体，仅返回错误。

## 新增路由步骤

1. 在 `logic/xxx.go` 定义 `Req/Res` 和 `XxxHandler` 方法。
2. 公共接口放 `internal/router/xxx.go`。
3. 鉴权接口放 `internal/router/auth_router/xxx.go`。
4. 使用合适的 `handler.Ctx*` 包装器挂载逻辑方法。
5. 写操作补 `OperationLogMiddleware(middleware.WithModule("<module>"))`。
6. 在注册入口接入新资源分组。

## Middleware 顺序

- 公共接口：按资源需要使用 `PublicMiddleware()`、`EncryptMiddleware()`。
- 鉴权接口组：先 `AuthMiddleware()`，再 `EncryptMiddleware()`。
- 审计接口：在具体 create/update/delete 路由补 `OperationLogMiddleware`。
- 不随意交换顺序，避免鉴权、解密、审计语义变化。

## 常用命令

```bash
cd backend/admin
go fix ./...
go vet ./...
go test ./...
```

## 注意事项

1. 路由文件保持薄，不写业务查询和数据库逻辑。
2. `logic` 层请求结构体统一写 `binding` 与 `binding_msg`。
3. 新资源不接入注册入口时接口不会生效。
4. 修改 Swagger 注释后执行 `make swagger`。
