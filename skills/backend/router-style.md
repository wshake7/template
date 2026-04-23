# Backend Router 编写风格

## 场景与目标
- 适用场景：新增/调整 `admin/router` 下接口路由与处理逻辑时
- 目标：保持当前路由分层、签名、middleware 串联风格一致

## 目录/文件位置
- 路由总注册：`backend/admin/router/router.go`
- 资源路由文件：`backend/admin/router/account.go`、`backend/admin/router/encrypt.go`
- 业务逻辑层：`backend/admin/router/logic/*.go`

## 当前分层约定
1. `router/*.go`：只做路由映射与 middleware 编排
2. `router/logic/*.go`：承载业务逻辑和请求/响应结构体
3. 业务入参出参通过 `handler.Ctx*` 包装器自动绑定和输出

## 路由注册风格（按现状）
- 统一在 `RegisterRouters()` 下挂载 `/api`
- 公共组 middleware 先走 `TimestampMiddleware()`
- 再按资源拆分子路由：
  - `registerAccountRouters(defaultGroup.Group("account"))`
  - `registerEncryptRouters(defaultGroup.Group("encrypt"))`

## Handler 写法风格
- 带请求体并返回数据：`handler.CtxHandlerFunc(...)`
- 仅返回 error：`handler.CtxHandlerNilFunc(...)`
- 简单 GET 返回：`handler.CtxFunc(...)`

## Middleware 串联风格
- 公开接口（如登录）：
  - `PublicMiddleware()` + `EncryptMiddleware()`
- 鉴权接口：
  - `AuthMiddleware()` +（按需）`EncryptMiddleware()`
- 操作日志接口：
  - `OperationLogMiddleware(...)` 放在业务 handler 之前

## 新增路由步骤
1. 在 `router/logic/xxx.go` 定义 `Req/Res` 与 handler 方法
2. 在 `router/xxx.go` 新增 `registerXxxRouters(router fiber.Router)`
3. 使用 `handler.Ctx*` 包装器接入业务方法
4. 在 `router/router.go` 注册 `registerXxxRouters(defaultGroup.Group("xxx"))`

## 常用命令
```bash
cd backend/admin

# 启动后验证路由
go run ./cmd/main.go -f ./etc/config.yaml

# 回归测试
go test ./...
```

## 注意事项
1. 路由文件不要写重业务逻辑，保持“薄路由、厚 logic”
2. middleware 顺序影响鉴权/加密/审计行为，不要随意交换
3. 新资源必须在 `router/router.go` 注册，否则接口不会生效
4. `logic` 层错误信息需区分用户可见与日志细节
