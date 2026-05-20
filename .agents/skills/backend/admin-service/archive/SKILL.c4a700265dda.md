# Skill: Backend Admin Service

## Generated Code Rule

- Do not edit generated files such as `internal/services/orm/query/*.gen.go` directly.
- If generated query behavior must change, update the source template in `cmd/scripts/orm/templates/` or the source model in `internal/services/orm/models/` first, then regenerate with `make script-orm` from `backend/admin`.
- Keep business logic in `internal/router/logic/**` or non-generated helpers, not in generated query code.

## 何时使用

当任务涉及 `backend/admin` 的服务启动、配置、路由注册、业务逻辑、中间件、权限、Temporal 调度、Swagger 或后台 API 时使用。

如果任务重点是“怎么给 admin backend 写测试 / 补覆盖率 / 选 mock 还是 SQLite”，优先配合 `backend/admin-testing` 技能一起使用。

## 核心路径

```text
backend/admin/
├── cmd/main.go
├── cmd/scripts/init.sql
├── cmd/scripts/schema_hardening_postgres.sql
├── cmd/scripts/orm/main.go
├── cmd/scripts/orm/templates/
├── etc/config.yaml
├── internal/config/
├── internal/fiberc/
├── internal/lifecycle/
├── internal/mock/                    # gomock 生成物，面向 service interfaces
├── internal/router/
│   ├── account.go                    # 公共账号路由
│   ├── encrypt.go                    # 公共加密路由
│   ├── auth_router/                  # 需鉴权资源路由
│   └── logic/                        # 业务逻辑与测试
├── internal/service/                 # 业务侧外部依赖接口与实现
├── internal/services/
│   ├── orm/
│   │   ├── models/
│   │   └── query/
│   ├── casbin/
│   ├── redisc/
│   ├── temporalc/
│   └── temporaljob/
└── docs/
```

## 当前 handler 注入约定

### 统一原则

- **所有数据库访问统一通过 `*query.Query` 注入到 logic handler。**
- **所有外部副作用保留为 `internal/service/**` 下的 interface 注入。**
- 不再为 admin logic 维护单独的 repository 抽象层。

### 当前常见构造模式

- 纯 DB handler：
  - `logic.NewSysUserHandler(query.Q)`
  - `logic.NewSysLanguageHandler(query.Q)`
  - `logic.NewSysApiLogHandler(query.Q)`
  - `logic.NewSysLoginLogHandler(query.Q)`
  - `logic.NewSysResourceMenuHandler(query.Q)`
- DB + 外部服务：
  - `logic.NewAccountHandler(query.Q, service.NewAuthService(), service.NewLoginLogger())`
  - `logic.NewSysResourceApiHandler(query.Q, service.NewCasbinService())`
  - `logic.NewSysRoleHandler(query.Q, service.NewCasbinService())`
  - `logic.NewSysDictHandler(query.Q, service.NewDataPermissionService())`
  - `logic.NewJobExecutionHandler(query.Q, service.NewTemporalService())`
  - `logic.NewJobScheduleHandler(query.Q, service.NewTemporalService())`
- 无 DB 例外：
  - `logic.NewEncryptHandler(service.NewRedisCache())`

### 不要再做的事

- 不要新增 `internal/repository/**` 来包装 admin logic 的数据库 CRUD。
- 不要为 DB handler 新增 repo mocks。
- 不要把登录日志写回 `login_log_record.go`，该文件已删除；账号登录日志通过 `service.LoginLogger` 注入。
- 不要让 logic 直接依赖全局外部客户端来完成可替换副作用，优先走 service interface。

## 启动链路

1. `cmd/main.go` 读取 `etc/config.yaml` 到 `config.Conf`。
2. 初始化日志。
3. `services.New(conf)` 初始化 ORM、Redis、HTTP client、Auth、Geo、Asynq、Casbin、Temporal 等底层依赖。
4. `fiberc.NewFiber(conf)` 创建 Fiber app。
5. `router.Router{}.RegisterRouters(group)` 注册 `/api/**` 路由。
6. `app.Start()` 启动服务。

## 路由与 logic 模式

- 公开接口使用 `/api/` 前缀，例如：
  - `/api/account/login/pwd`
  - `/api/account/changePwd`
  - `/api/encrypt/public/key`
- 需鉴权资源统一使用 `/api/sys/` 前缀，例如：
  - `/api/sys/user/list`
  - `/api/sys/role/list`
  - `/api/sys/resource/menu/list`
  - `/api/sys/resource/api/list`
  - `/api/sys/dict/entry/match`
  - `/api/sys/job/schedule/list`
  - `/api/sys/job/execution/list`
- 业务逻辑放在 `internal/router/logic/<resource>.go`。
- 路由注册放在 `internal/router/*.go` 与 `internal/router/auth_router/*.go`。
- 请求/响应结构体优先与对应 handler 放在同一个 logic 文件中。
- 响应错误优先复用 `internal/fiberc/res` 中的标准错误。

## 外部依赖抽象

当前 admin backend 已沉淀的可 mock 服务接口位于 `internal/service/**`：

- `AuthService`
- `LoginLogger`
- `RedisCache`
- `CasbinService`
- `TemporalService`
- `DataPermissionService`

对应的 gomock 生成物位于 `internal/mock/**`，由各接口文件中的 `//go:generate mockgen ...` 维护。

## 新增后台资源的推荐流程

1. 在 `internal/services/orm/models` 新增或修改 model。
2. 如需 query 扩展，修改 `cmd/scripts/orm/templates/`。
3. 在 `backend/admin` 下运行 `make script-orm`。
4. 在 `internal/router/logic/<resource>.go` 编写 handler，请直接注入 `*query.Query`；只有外部副作用才新增 service interface 依赖。
5. 在 `internal/router/auth_router/<resource>.go` 或 `internal/router/<resource>.go` 注册路由，并传入 `query.Q` 与必要的 service 实现。
6. 如涉及种子数据或约束，同步修改 `cmd/scripts/init.sql` 或迁移脚本。
7. 如涉及权限、调度或缓存副作用，同步接入 `CasbinService`、`TemporalService`、`RedisCache` 等接口。
8. 运行 `make swagger` 更新 Swagger。
9. 按 `backend/admin-testing` 技能补齐测试。

## 安全与中间件

- 安全中间件位于 `internal/fiberc/middleware/security.go`。
- 常用中间件包括：时间戳校验、Nonce 防重放、请求解密/响应加密、签名验证。
- `EncryptMiddleware` 负责解密请求并加密响应。
- SSE 路由 `auth_router/events.go` 也会对事件 payload 做加密，并监听 `lifecycle.ShutdownDone()` 以便优雅关闭。

## 任务调度与 Temporal

- 业务逻辑不再直接依赖 `temporalc` 包完成可测试行为，而是通过 `service.TemporalService` 注入。
- `internal/services/temporalc/` 仍然是底层客户端封装位置。
- `internal/services/temporaljob/` 存放工作流定义。
- 对 `job_execution` / `job_schedule` 做功能修改时，优先保持：
  - DB 状态由 `*query.Query` 负责
  - Temporal 调用由 `TemporalService` 负责

## Casbin 权限自动同步

- Casbin 执行器初始化位于 `internal/services/casbin/`。
- 角色、API 资源等权限相关变更，应通过 `CasbinService` 触发同步，不要把 Casbin 调用散落进路由层。
- 涉及 `sys_role`、`sys_resource_api` 变更时，记得检查对应同步分支是否仍然成立。

## 推荐验证命令

```bash
cd backend/admin
go test ./...
go test ./internal/router/logic/... -cover
```

如修改了 ORM 生成链路，再补：

```bash
cd backend/admin
make script-orm
go test ./...
```
