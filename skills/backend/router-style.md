# Backend Router 编写风格

## 场景与目标
- 适用场景：新增/调整 `admin/router` 下接口路由与处理逻辑时
- 目标：保持当前路由分层、签名、middleware 串联风格一致

## 目录/文件位置
- 路由总注册：`backend/admin/internal/router/router.go`
- 公共资源路由：`backend/admin/internal/router/account.go`、`backend/admin/internal/router/encrypt.go`
- 鉴权资源路由：`backend/admin/internal/router/auth_router/*.go`
- 业务逻辑层：`backend/admin/internal/router/logic/*.go`

## 当前分层约定
1. `internal/router/*.go`：只做路由映射与 middleware 编排
2. `internal/router/auth_router/*.go`：承载"已鉴权资源"的路由映射
3. `internal/router/logic/*.go`：承载业务逻辑和请求/响应结构体
4. 业务入参出参通过 `handler.Ctx*` 包装器自动绑定和输出

## 路由注册风格（按现状）
- 统一在 `RegisterRouters()` 下挂载 `/api`
- 公共组 middleware 先走 `TimestampMiddleware()`
- 再按资源拆分子路由：`registerAccountRouters(defaultGroup.Group("account"))`
- 再按资源拆分子路由：`registerEncryptRouters(defaultGroup.Group("encrypt"))`
- 再按资源拆分子路由：`auth_router.RegisterRouters(defaultGroup.Group("role"))`
- `auth_router.RegisterRouters()` 内部统一追加鉴权与解密中间件，再分发到 `registerRoleRouters()`

## Handler 写法风格
- `handler.CtxFunc(fn)`：无请求结构体，仅依赖上下文，返回 `(*Res, error)`（如公钥查询）
- `handler.CtxHandlerFunc(fn)`：有请求结构体，返回 `(*Res, error)`（如分页查询）
- `handler.CtxHandlerNilFunc(fn)`：有请求结构体，仅返回 `error`（如 create/update/delete）
- 入参统一使用结构体标签校验（`binding` + `binding_msg`），避免在路由层手写解析

### 更新接口 (Patch 模式) 规范
- **结构体定义**：Update 请求结构体中的可选字段必须使用 **指针类型**（如 `*string`, `*int32`），以区分“未传值”与“传零值”。
- **逻辑实现**：使用 Repository 提供的 `UpdateNoNilMap` 方法进行更新，该方法会自动过滤掉 `nil` 指针，实现按需更新。
- **校验逻辑**：
    - 状态位校验：使用 `binding:"oneof=0 1"` 限制开关状态。
    - 关联性校验：在创建或更新包含外键（如 `SysDictTypeId`）的记录前，必须先调用 `repo.XxxRepo.Exists` 校验关联数据是否存在。

## Middleware 串联风格
- 公共接口（如登录）：`PublicMiddleware()` + `EncryptMiddleware()`
- 鉴权接口组：优先在分组层统一挂 `AuthMiddleware()` + `EncryptMiddleware()`
- 审计接口（create/update/delete）：在具体路由补 `OperationLogMiddleware(middleware.WithModule("<module>"))`
- 中间件顺序不要随意互换，先鉴权/解密，再做业务审计

## 新增路由步骤
1. 在 `internal/router/logic/xxx.go` 定义 `Req/Res` 与 `XxxHandler` 方法（校验标签与错误语义一并补齐）
2. 根据接口类型选择路由文件：
   - 公共接口：新增/修改 `internal/router/xxx.go`
   - 鉴权接口：新增/修改 `internal/router/auth_router/xxx.go`
3. 用 `handler.Ctx*` 包装器挂载逻辑方法，不在路由层堆业务代码
4. 在 `internal/router/router.go` 挂载资源分组；若是鉴权资源，走 `auth_router.RegisterRouters(...)` 统一收口
5. 对写操作接口补 `OperationLogMiddleware`，并使用 `WithModule("<module>")` 标识模块

## 常用命令
```bash
cd backend/admin

# 启动后验证路由
go run ./cmd/main.go -f ./etc/config.yaml

# 基础修复与回归
go fix ./...
go test ./...
```

## 注意事项
1. 路由文件不要写重业务逻辑，保持“薄路由、厚 logic”
2. middleware 顺序影响鉴权/加密/审计行为，不要随意交换
3. 新资源必须在 `internal/router/router.go` 或 `auth_router` 注册入口挂载，否则接口不会生效
4. `logic` 层错误信息需区分用户可见错误与日志细节（日志写 `ctx.L().Error(...)`）
