# Skill: Backend Admin Service

## Generated Code Rule

- Do not edit generated files such as `internal/services/orm/query/*.gen.go` directly.
- If generated query behavior must change, update the source template in `cmd/scripts/orm/templates/` or the source model in `internal/services/orm/models/` first, then regenerate with `make script-orm` from `backend/admin`.
- Do not put table-specific business behavior into shared templates; keep shared templates generic and place per-model query behavior in non-generated logic, reusable scopes, or model definitions.
- Review the regenerated diff after `make script-orm`; keep manual business logic in non-generated files.

## 何时使用

当任务涉及 `backend/admin` 的服务启动、配置、路由、业务逻辑、中间件、认证权限、Swagger、GORM Gen 生成代码或后台 API 时使用。特别适合开发系统管理类资源（如角色、用户、菜单、API 资源、字典、语言、API 日志、登录日志等）。

## 核心路径

text
backend/admin/
├── cmd/main.go                         # 服务入口
├── cmd/scripts/init.sql                # 初始化种子数据
├── cmd/scripts/schema_hardening_postgres.sql  # 数据库结构与约束加固
├── cmd/scripts/orm/main.go             # ORM 代码生成（不再负责数据插入）
├── cmd/scripts/orm/templates/          # query 扩展模板
├── etc/config.yaml                     # 本地配置
├── internal/config/                    # 配置结构
├── internal/fiberc/                    # Fiber app、handler、middleware、response
├── internal/router/                    # 路由注册与业务 logic
├── internal/router/auth_router/        # 需要鉴权的系统资源路由
│   ├── sys_login_log.go               # 登录日志相关路由
│   ├── events.go                      # SSE 事件流路由
│   └── ...
├── internal/router/logic/              # 业务逻辑包（建议请求/响应定义集中于此）
│   ├── sys_resource_api.go
│   ├── sys_resource_menu.go
│   ├── sys_role.go
│   ├── sys_user.go
│   ├── sys_login_log.go               # 登录日志列表与详情
│   └── ...
├── internal/services/                  # ORM、Redis、Auth、Casbin、Asynq、HTTP、Geo
│   ├── orm/
│   │   ├── models/                     # GORM models
│   │   │   ├── sys_login_log.go       # 登录日志模型
│   │   │   └── ...
│   │   ├── query/                      # GORM Gen 生成代码与扩展
│   │   └── data_permission/            # 数据权限引擎（规划中）
│   └── casbin/
│       ├── casbin.go                   # 初始化 Casbin 执行器
│       ├── policy.go                   # 权限策略自动同步
│       └── pbac.conf                  # Casbin 模型定义
└── docs/                               # swag 生成的 Swagger 文件


## 启动链路

1. `cmd/main.go` 读取 `etc/config.yaml` 到 `config.Conf`。
2. 初始化日志：`go-common/log`。
3. `services.New(conf)` 初始化 ORM、Redis、HTTP client、Auth、Geo、Asynq、Casbin。
4. `fiberc.NewFiber(conf)` 创建 Fiber app。
5. `router.Router{}.RegisterRouters(group)` 注册 `/api/**` 路由。
6. `app.Start()` 启动服务，端口来自配置。

## API 响应码约定

- 所有接口返回 HTTP 200，业务结果通过 `code` 字段区分。
- 常见 code 及含义：

| Code | 说明 |
|------|------|
| 1    | 成功 |
| 2    | 服务繁忙 / 通用失败 |
| 3    | 请求超时 |
| 4    | 请求重放（Nonce 校验失败）|
| 5    | 请求错误（客户端错误）|
| 100  | 登录 / 权限相关失败 |
| 200  | 授权相关失败 |

- 使用 `internal/fiberc/res` 包中的标准错误构造函数，保证响应格式一致。

## 路由与 Handler 模式

- 路由前缀约定：
  - 公开接口（如账号、加密 key）使用 `/api/` 前缀，例如 `/api/account/login/pwd`、`/api/encrypt/public/key`。
  - 需鉴权的系统管理资源统一使用 `/api/sys/` 前缀，例如 `/api/sys/api/log/list`、`/api/sys/dict/entry/match`、`/api/sys/resource/menu/list`、`/api/sys/resource/api/list`、`/api/sys/role/list`、`/api/sys/user/list`、`/api/sys/login/log/list`、`/api/sys/login/log/detail` 等。
  - SSE 事件流端点使用 `/api/events`，通过 `auth_router/events.go` 注册。
- 普通路由聚合在 `internal/router/router.go`。
- 业务逻辑放在 `internal/router/logic/<resource>.go`，通过 `handler.CtxHandlerFunc` 等包装。
- 操作日志通过 `middleware.OperationLogMiddleware(middleware.WithModule("<module>"))` 注入。
- API 请求日志通过 `internal/fiberc/middleware/api_log.go` 中的中间件自动记录到 `sys_api_log` 表，该中间件使用 `github.com/mileusna/useragent` 解析 User-Agent 提取浏览器和操作系统信息，并对 request/response payload 中的敏感字段（如 token、password）进行脱敏处理。
- 登录日志由 `internal/router/logic/login_log_record.go` 在账号操作中记录，其字段与模型 `SysLoginLog` 保持一致。
- 响应错误优先返回 `internal/fiberc/res` 中的标准错误，不直接拼散乱响应。

## 新增后台资源的推荐流程

1. 在 `internal/services/orm/models` 新增 model，并在 `init()` 中追加到 `Models`（如 `sys_login_log` 模型）。
2. 如需特定查询扩展，修改 `cmd/scripts/orm/templates/` 下的模板。
3. 运行 `make script-orm` 重新生成 `internal/services/orm/query/` 下的代码。
4. 在 `internal/router/logic/<resource>.go` 定义请求/响应结构体（建议 `ReqXxx`、`RespXxx`）和业务方法（List/Create/Update/Delete/Switch 等）。响应结构体避免直接暴露 ORM Model。
5. 在 `internal/router/auth_router/<resource>.go` 注册路由，使用 `/api/sys/<resource>/*` 路径前缀。例如登录日志列表 `POST /api/sys/login/log/list`、详情 `POST /api/sys/login/log/detail`。
6. 在 `auth_router.RegisterRouters` 中汇总注册。
7. 如果需要种子数据（如超级管理员角色、默认菜单树、字典初始条目），在 `cmd/scripts/init.sql` 中添加 INSERT 语句，并确保幂等（常用 `ON CONFLICT` 或先删后插）。对于复杂数据库结构变更，编写迁移脚本（如 `schema_hardening_postgres.sql`）并在部署时执行。
8. 同步更新 Casbin 权限策略。若资源增删改会影响角色授权，需在业务逻辑中调用 `casbin` 包的同步方法（参考下一节）。
9. 运行 `make swagger` 更新 Swagger 文档；分页查询响应的 `data` 应引用 logic 包下的自定义 `Resp` 类型。
10. 前端同步新增 API 文件和页面时，切换到 frontend 技能。

## Casbin 权限自动同步

- Casbin 执行器已在启动时创建（`internal/services/casbin/casbin.go`），模型定义在 `pbac.conf` 中，当前使用简化的匹配器：
  [policy_definition]
  p = sub, obj, act

  [matchers]
  m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && r.act == p.act
  
  策略直接存储主体标识（如 `role:root`）、API 路径和方法，不再使用 eval 表达式。
- 所有权限变更均通过 `internal/services/casbin/policy.go` 中的函数自动同步，无需手动操作 Casbin API。
- 关键同步函数：
  - `AddRoleAPIPolicies` / `RemoveRoleAPIPolicies`：为角色添加或移除一组 API 权限。
  - `SyncRoleAPIPermissions`：当角色权限保存时，计算新旧 API ID 差异，增量更新策略。
  - `SyncRoleState`：当角色 code 或启用状态变更时，先移除旧策略再添加新策略。
  - `AddAPIResourcePolicies` / `RemoveAPIResourcePolicies`：当 API 资源被启用/禁用或路径/方法变更时，重新计算关联角色的策略。
  - `SyncAPIResourcePolicies`：API 资源更新时，移除旧策略并添加新策略。
- 在业务逻辑中（如 `sys_resource_api.go`、`sys_role.go`），当操作会影响权限时，必须调用对应的同步函数以确保 Casbin 策略与数据库一致。
- 种子数据（`cmd/scripts/init.sql`）已包含 root 角色的完整策略，部署后即可使用。

## 数据权限约定

- `sys_data_permission` 表新增 `action_key` 字段（varchar），存储去重后的 action 列表（如 `read,write`），用于唯一索引 `idx_sys_data_permission_subject_resource_action_active`，避免主体-资源-action 组合重复。
- 插入或更新数据权限记录时，务必保持 `action` JSON 数组与 `action_key` 字段同步。
- 权限过滤引擎仍基于 `action` 字段，无需改动。
- 系统内置资源的保护优先通过 `conditions` 配置（如 `id__not:1`），避免硬编码特殊 ID。
- 列表接口如需控制行级操作按钮，响应结构体应携带 `canWrite`、`canDelete` 等权限标记。

## 登录日志记录与查询

- 模型 `SysLoginLog` 记录每次登录/退出尝试，包含 `username`、`loginIP`、`userAgent`、`browserName`、`osName`、`statusCode`、`success`、`reason` 等字段，其中 `sys_user_id` 可为空（当用户不存在时）。
- 日志由 `login_log_record.go` 在账号相关操作时自动写入，业务代码无需显式调用。
- 前端可通过 `POST /api/sys/login/log/list` 分页查询日志，支持按用户名、IP、状态等筛选。
- 详情接口 `POST /api/sys/login/log/detail` 根据日志 ID 返回完整记录。
- SSE 事件流 `/api/events` 使用 `auth_router/events.go` 推送实时登录事件，前端通过 `api/eventStream.ts` 订阅。

## 配置与服务

- 本地配置在 `etc/config.yaml`，结构定义在 `internal/config/**`。
- ORM 服务在 `internal/services/orm.go` 与 `internal/services/orm/orm.go`。
- Redis 服务在 `internal/services/redis.go` 与 `internal/services/redisc/**`。
- Casbin 模型文件在 `internal/services/casbin/*.conf`。
- Header、HTTP code 等跨端约定需与前端 `src/domains/http.ts` 对齐。

## 命令

bash
cd backend/admin
go run .
make script-orm   # 生成 ORM query 代码（go run ./cmd/scripts/orm）
make swagger      # 更新 Swagger 文档
go fix ./...
go vet ./...
go test ./...

# 数据库迁移（根据环境执行加固脚本）
psql $DATABASE_URL -f cmd/scripts/schema_hardening_postgres.sql


## 验证

- 修改 `backend/admin` 后，在该模块执行 `go fix ./...`、`go vet ./...`、`go test ./...`。
- 修改 Swagger 注释后执行 `make swagger` 并检查 `docs/**`。
- 修改模型、请求/响应结构体或 query 生成模板后，先运行 `make script-orm` 重新生成，再确认生成文件与 Swagger 文档符合预期，避免手动修改生成产物后被覆盖。
- 若涉及 Casbin 权限同步逻辑，确保业务变更触发相应的同步函数，并验证种子数据执行后的策略正确性。
- 数据库迁移脚本需幂等大量使用 `ADD COLUMN IF NOT EXISTS`、`DROP CONSTRAINT IF EXISTS` 等安全操作，执行后检查索引和约束是否生效。
