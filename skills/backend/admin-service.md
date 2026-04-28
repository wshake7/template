# Backend Admin 服务开发

## 场景与目标

- 适用场景：修改 `backend/admin` 的接口、配置、路由、中间件、ORM、Swagger 或服务初始化。
- 目标：按现有分层落地改动，保持启动链路、错误处理和文档生成一致。

## 目录/文件位置

- 启动入口：`backend/admin/cmd/main.go`
- 配置定义：`backend/admin/internal/config/`
- 配置文件：`backend/admin/etc/config.yaml`
- Fiber 初始化：`backend/admin/internal/fiberc/`
- 路由注册：`backend/admin/internal/router/`
- 业务逻辑：`backend/admin/internal/router/logic/`
- 服务初始化：`backend/admin/internal/services/init.go`
- ORM：`backend/admin/internal/services/orm/`
- Swagger 产物：`backend/admin/docs/`

## 当前启动链路

1. `cmd/main.go` 读取 `-f` 指定配置，默认 `./etc/config.yaml`。
2. 初始化日志、配置和 `services.New(conf)`。
3. `services.New` 组装 Httpc、Orm、Redis、Auth、Geo、Asynq、Casbin 等服务。
4. 创建 Fiber App，注册 `/api` 路由组和中间件。
5. 启动服务并等待优雅退出。

## 接口开发步骤

1. 在 `internal/router/logic/xxx.go` 定义请求结构体、返回类型和 handler 方法。
2. 公共接口在 `internal/router/xxx.go` 注册。
3. 鉴权接口在 `internal/router/auth_router/xxx.go` 注册。
4. 写操作补 `OperationLogMiddleware(middleware.WithModule("<module>"))`。
5. 涉及数据库时使用 `query.Xxx`、`models.Xxx`、`gormc.PagingResult`。
6. 涉及 Swagger 注释时执行 `make swagger`。
7. 最后执行 `go fix ./...`、`go test ./...`。

## Handler 风格

- 列表/查询：返回 `(*Data, error)`。
- 创建/更新/删除：返回 `error`。
- 入参统一使用 `json`、`binding`、`binding_msg`。
- Create 请求使用普通字段。
- Update 请求可选字段使用指针类型，保留 Patch 语义。
- 分页列表复用 `v1.PagingRequest`。
- 用户可感知错误返回 `res.FailMsg("...")`。
- 通用异常返回 `res.FailDefault`。
- 关键失败点用 `ctx.L().Error(...)` 记录日志，日志细节不直接暴露给用户。

## 数据库写法

- 列表：`query.Xxx.PageWithPaging(req)`。
- 创建：`query.Xxx.Create(&models.Xxx{...})`。
- Patch 更新：`[]field.AssignExpr` + `query.ExprAppendSelf` + `UpdateSimple`。
- 删除：按业务策略选择 `Delete()` 或模型软删除能力。
- 审计字段：从 `ctx.SessionInfo.Id` 写入 `CreatedBy` / `UpdatedBy`。
- 唯一键冲突：捕获 `gorm.ErrDuplicatedKey` 并转为明确业务提示。

## 常用命令

```bash
cd backend/admin
go run ./cmd/main.go -f ./etc/config.yaml
go run ./cmd/scripts/orm
make swagger
go fix ./...
go test ./...
```

## 注意事项

1. 路由文件只做映射和中间件编排，业务逻辑放在 `internal/router/logic`。
2. 新资源必须接入 `internal/router/router.go` 或 `auth_router.RegisterRouters`。
3. 配置新增字段要同步结构体和 `etc/config.yaml`。
4. ORM 模型或 query 生成见 `orm-models.md`。
