# Skill: Backend Admin Service

## 何时使用

当任务涉及 `backend/admin` 的服务启动、配置、路由、业务逻辑、中间件、认证权限、Swagger、GORM Gen 生成代码或后台 API 时使用。

## 核心路径

```text
backend/admin/
├── cmd/main.go                         # 服务入口
├── cmd/scripts/init.sql                # 初始化 SQL
├── cmd/scripts/orm/main.go             # ORM 代码生成与种子用户
├── cmd/scripts/orm/templates/          # query 扩展模板
├── etc/config.yaml                     # 本地配置
├── internal/config/                    # 配置结构
├── internal/fiberc/                    # Fiber app、handler、middleware、response
├── internal/router/                    # 路由注册与业务 logic
├── internal/services/                  # ORM、Redis、Auth、Casbin、Asynq、HTTP、Geo
├── internal/services/orm/models/       # GORM models
├── internal/services/orm/query/        # GORM Gen 生成代码与扩展
└── docs/                               # swag 生成的 Swagger 文件
```

## 启动链路

1. `cmd/main.go` 读取 `etc/config.yaml` 到 `config.Conf`。
2. 初始化日志：`go-common/log`。
3. `services.New(conf)` 初始化 ORM、Redis、HTTP client、Auth、Geo、Asynq、Casbin。
4. `fiberc.NewFiber(conf)` 创建 Fiber app。
5. `router.Router{}.RegisterRouters(group)` 注册 `/api/**` 路由。
6. `app.Start()` 启动服务，端口来自配置。

## 路由与 Handler 模式

- 普通路由入口在 `internal/router/router.go`。
- 公开接口如账号、加密 key 放在 `internal/router/account.go`、`encrypt.go`。
- 需鉴权的系统资源放在 `internal/router/auth_router/**`。
- 业务逻辑放在 `internal/router/logic/**`，通过 `handler.CtxHandlerFunc` 或 `CtxHandlerNilFunc` 包装。
- 操作日志通过 `middleware.OperationLogMiddleware(middleware.WithModule("<module>"))` 注入。
- 响应错误优先返回 `internal/fiberc/res` 中的标准错误，不直接拼散乱响应。

## 新增后台资源的推荐流程

1. 在 `internal/services/orm/models` 新增 model，并在 `init()` 中追加到 `Models`。
2. 如果需要通用字段，优先复用 `orm-crud/gormc/mixin`。
3. 运行 ORM 生成脚本，更新 `internal/services/orm/query/**`：

```bash
cd backend/admin
go run ./cmd/scripts/orm
```

4. 在 `internal/router/logic/<resource>.go` 编写 List/Create/Update/Switch/Delete 等方法。
5. 在 `internal/router/auth_router/<resource>.go` 注册鉴权路由。
6. 在 `auth_router.RegisterRouters` 聚合注册新资源。
7. 如果有 Swagger 注释，运行 `make swagger`。
8. 前端同步新增 `src/api/<resource>.ts` 和页面时，切换到 frontend 技能。

## 配置与服务

- 本地配置在 `etc/config.yaml`，结构定义在 `internal/config/**`。
- ORM 服务在 `internal/services/orm.go` 与 `internal/services/orm/orm.go`。
- Redis 服务在 `internal/services/redis.go` 与 `internal/services/redisc/**`。
- Casbin 模型文件在 `internal/services/casbin/*.conf`。
- Header、HTTP code 等跨端约定需与前端 `src/domains/http.ts` 对齐。

## 命令

```bash
cd backend/admin
go run ./cmd
go run ./cmd/scripts/orm
make swagger
go fix ./...
go vet ./...
go test ./...
```

## 验证

- 修改 `backend/admin` 后，在该模块执行 `go fix ./...`、`go vet ./...`、`go test ./...`。
- 修改 Swagger 注释后执行 `make swagger` 并检查 `docs/**`。
- 修改模型或 query 生成模板后，确认生成文件符合预期，避免手改生成产物后被脚本覆盖。

## 自动优化记录

<!-- ai-skill-optimizer:6eab4249bb74:1 -->
### 数据权限设计与列表响应携带操作标记

- 系统内置资源（如系统状态字典类型）的保护，应通过 sys_data_permission 表配置自定义条件（如 id__not:1），而非业务代码硬编码
- 角色权限可配置 "all" 动作，覆盖 read/write/delete，减少多条规则
- 列表接口返回自定义 Resp 结构体（如 RespDictType），携带 canWrite/canDelete 等权限标记，由数据权限引擎计算
- 分页查询响应 Swagger 注解中 data 类型应指向自定义 Resp，不再直接对外暴露 ORM Model
- ORM 模板新增 WithDBScopes 函数，便于在查询时注入动态 Scope（如数据权限过滤）

