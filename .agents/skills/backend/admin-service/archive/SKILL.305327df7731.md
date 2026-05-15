# Skill: Backend Admin Service

## Generated Code Rule

- Do not edit generated files such as `internal/services/orm/query/*.gen.go` directly.
- If generated query behavior must change, update the source template in `cmd/scripts/orm/templates/` or the source model in `internal/services/orm/models/` first, then regenerate with `make script-orm` from `backend/admin`.
- Do not put table-specific business behavior into shared templates; keep shared templates generic and place per-model query behavior in non-generated logic, reusable scopes, or model definitions.
- Review the regenerated diff after `make script-orm`; keep manual business logic in non-generated files.

## 何时使用

当任务涉及 `backend/admin` 的服务启动、配置、路由、业务逻辑、中间件、认证权限、Swagger、GORM Gen 生成代码、Temporal 任务调度或后台 API 时使用。特别适合开发系统管理类资源（如角色、用户、菜单、API 资源、字典、语言、API 日志、登录日志、任务调度等）。

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
├── internal/lifecycle/                 # 服务生命周期管理，优雅关闭信号广播
├── internal/router/                    # 路由注册与业务 logic
├── internal/router/auth_router/        # 需要鉴权的系统资源路由
│   ├── sys_login_log.go               # 登录日志相关路由
│   ├── events.go                      # SSE 事件流路由
│   ├── job_execution.go               # 任务执行记录路由
│   ├── job_schedule.go                # 任务调度路由
│   └── ...
├── internal/router/logic/              # 业务逻辑包（建议请求/响应定义集中于此）
│   ├── sys_resource_api.go
│   ├── sys_resource_menu.go
│   ├── sys_role.go
│   ├── sys_user.go
│   ├── sys_login_log.go               # 登录日志列表与详情
│   ├── job_execution.go               # 任务执行记录逻辑
│   ├── job_schedule.go                # 任务调度逻辑
│   └── ...
├── internal/services/                  # ORM、Redis、Auth、Casbin、Asynq、HTTP、Geo、Temporal
│   ├── orm/
│   │   ├── models/                     # GORM models
│   │   │   ├── sys_login_log.go       # 登录日志模型
│   │   │   ├── job_execution.go       # 任务执行记录模型
│   │   │   ├── job_schedule.go        # 任务调度模型
│   │   │   └── ...
│   │   ├── query/                      # GORM Gen 生成代码与扩展
│   │   │   ├── job_execution.gen.go   # 生成的基础查询
│   │   │   ├── job_execution_extend.gen.go   # 扩展查询
│   │   │   ├── job_schedule.gen.go
│   │   │   ├── job_schedule_extend.gen.go
│   │   └── data_permission/            # 数据权限引擎（规划中）
│   ├── temporal.go                     # Temporal 客户端初始化
│   ├── temporalc/                      # Temporal 客户端封装
│   │   └── temporalc.go
│   ├── temporaljob/                    # 任务工作流定义
│   │   ├── example.go
│   │   └── job.go
│   └── casbin/
│       ├── casbin.go                   # 初始化 Casbin 执行器
│       ├── policy.go                   # 权限策略自动同步
│       └── pbac.conf                  # Casbin 模型定义
└── docs/                               # swag 生成的 Swagger 文件


## 启动链路

1. `cmd/main.go` 读取 `etc/config.yaml` 到 `config.Conf`。
2. 初始化日志：`go-common/log`。
3. `services.New(conf)` 初始化 ORM、Redis、HTTP client、Auth、Geo、Asynq、Casbin，以及 Temporal 客户端（通过 `temporal.go` 和 `temporalc` 包）。
4. `fiberc.NewFiber(conf)` 创建 Fiber app。
5. `router.Router{}.RegisterRouters(group)` 注册 `/api/**` 路由。
6. `app.Start()` 启动服务，端口来自配置。

## 优雅关闭流程

- 服务通过 `fiberc.gracefulShutdown` 监听 OS 信号 (SIGINT, SIGTERM)，收到信号后调用 `app.ShutdownWithTimeout(5 * time.Second)` 启动超时关闭。
- 在 Fiber 的 `OnPreShutdown` 钩子中调用 `lifecycle.BeginShutdown()`，关闭 `lifecycle.shutdownDone` 通道，通知所有等待关闭的组件（例如 SSE 流）。
- `lifecycle` 包提供 `BeginShutdown()`（一次性关闭信号）和 `ShutdownDone()`（返回只读通道）。
- 长时间运行的任务（如 SSE 连接）应在循环中 select `lifecycle.ShutdownDone()` 主动退出，避免阻塞关闭过程。
- 关闭流程完成后的清理工作在 `OnPostShutdown` 钩子中处理（如日志同步、资源释放）。

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
  - 需鉴权的系统管理资源统一使用 `/api/sys/` 前缀，例如 `/api/sys/api/log/list`、`/api/sys/dict/entry/match`、`/api/sys/resource/menu/list`、`/api/sys/resource/api/list`、`/api/sys/role/list`、`/api/sys/user/list`、`/api/sys/login/log/list`、`/api/sys/login/log/detail`、`/api/sys/job/schedule/list`、`/api/sys/job/execution/list` 等。
  - SSE 事件流端点使用 `/api/events`，通过 `auth_router/events.go` 注册。
- 普通路由聚合在 `internal/router/router.go`。
- 业务逻辑放在 `internal/router/logic/<resource>.go`，通过 `handler.CtxHandlerFunc` 等包装。
- 操作日志通过 `middleware.OperationLogMiddleware(middleware.WithModule("<module>"))` 注入。
- API 请求日志通过 `internal/fiberc/middleware/api_log.go` 中的中间件自动记录到 `sys_api_log` 表，该中间件使用 `github.com/mileusna/useragent` 解析 User-Agent 提取浏览器和操作系统信息，并对 request/response payload 中的敏感字段（如 token、password）进行脱敏处理。
- 登录日志由 `internal/router/logic/login_log_record.go` 在账号操作中记录，其字段与模型 `SysLoginLog` 保持一致。
- 响应错误优先返回 `internal/fiberc/res` 中的标准错误，不直接拼散乱响应。

## 请求加密与安全中间件

安全中间件位于 `internal/fiberc/middleware/security.go`，提供时间戳校验、Nonce 重放防护、请求解密与响应加密、请求签名验证等功能。

- **TimestampMiddleware**：校验请求头 `X-Request-Timestamp` 与服务器时间的差值不能超过配置的过期时间，防止重放。
- **NonceMiddleware**：将请求头 `X-Request-ID` 存入 Redis，并设置过期时间，为后续 Nonce 唯一性检查提供依据。
- **EncryptMiddleware**：使用可复用的函数完成请求解密和响应加密。
  - 首先调用 `DecryptRequest(ctx)` 获取 AES 密钥和明文请求体。
  - `DecryptRequest` 内部：
    - `DecryptRequestAESKey` 从请求头 `X-Request-Encrypted-Key` 用服务端 RSA 私钥解密出 AES 密钥。
    - `DecryptRequestBody` 使用 AES 密钥和 `RequestAAD(ctx)` 生成的附加认证数据（包含 `X-Request-ID`、`X-Request-Timestamp` 和所有查询参数）解密请求体，解密成功后替换原始请求体。
  - 处理业务逻辑后，使用 `EncryptText`（`aes_util.Encrypt` 封装）加密响应体，并在响应头添加 `X-Response-Is-Encrypt: true`。
- **SignMiddleware**：用于不需要解密但需要验证请求完整性的场景（例如简单签名校验）。它构造与 `RequestAAD` 相似的参数组合，对请求体计算签名并与请求头 `X-Request-Signature` 对比。注意：该中间件不用于标准加密流程，仅用于特定签名接口。

SSE 事件流路由 `auth_router/events.go` 也应用了加密：
- 在处理连接时，调用 `middleware.DecryptRequest(c)` 获取 AES 密钥。
- 之后对每个推送的事件数据先用 `json.Marshal` 序列化原始数据，再用 `middleware.EncryptText` 加密，最后包装成 `{"payload": "加密后的字符串"}` 并通过 SSE 发送。
- 响应头同样设置 `X-Response-Is-Encrypt: true`。
- 事件循环现在同时监听 `lifecycle.ShutdownDone()` 信号，服务器关闭时会主动退出循环，避免阻塞 shutdown。

## 新增后台资源的推荐流程

1. 在 `internal/services/orm/models` 新增 model，并在 `init()` 中追加到 `Models`（如 `job_schedule`、`job_execution` 模型）。
2. 如需特定查询扩展，修改 `cmd/scripts/orm/templates/` 下的模板。
3. 运行 `make script-orm` 重新生成 `internal/services/orm/query/` 下的代码。
4. 在 `internal/router/logic/<resource>.go` 定义请求/响应结构体（建议 `ReqXxx`、`RespXxx`）和业务方法（List/Create/Update/Delete/Switch/Detail/Trigger/Sync/Cancel/Retry 等）。响应结构体避免直接暴露 ORM Model。
5. 在 `internal/router/auth_router/<resource>.go` 注册路由，使用 `/api/sys/<resource>/*` 路径前缀。例如任务调度列表 `POST /api/sys/job/schedule/list`、详情 `POST /api/sys/job/schedule/detail`；任务执行列表 `POST /api/sys/job/execution/list`、取消 `POST /api/sys/job/execution/cancel` 等。
6. 在 `auth_router.RegisterRouters` 中汇总注册。
7. 如果需要种子数据（如超级管理员角色、默认菜单树、字典初始条目），在 `cmd/scripts/init.sql` 中添加 INSERT 语句，并确保幂等（常用 `ON CONFLICT` 或先删后插）。对于复杂数据库结构变更，编写迁移脚本（如 `schema_hardening_postgres.sql`）并在部署时执行。
8. 同步更新 Casbin 权限策略。若资源增删改会影响角色授权，需在业务逻辑中调用 `casbin` 包的同步方法（参考下一节）。
9. 对于涉及 Temporal 的资源（如 `job_schedule`），需要确保 Temporal 客户端已初始化并在业务逻辑中通过 `temporalc` 包与 Temporal 通信。任务工作流定义应放在 `internal/services/temporaljob/` 目录下。
10. 运行 `make swagger` 更新 Swagger 文档；分页查询响应的 `data` 应引用 logic 包下的自定义 `Resp` 类型。
11. 前端同步新增 API 文件和页面时，切换到 frontend 技能。

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
- 种子数据（`cmd/scripts/init.sql`）已包含 root 角色的完整策略，包括新增的任务调度 API 权限，部署后即可使用。

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
- 详情接口 `POST /api/sys/login/log/detail` 返回完整信息。

## 任务调度模块约定

- 任务调度基于 Temporal 工作流引擎。
- 所有 Temporal 客户端连接池由 `internal/services/temporal.go` 初始化，配置项在 `etc/config.yaml` 中（Temporal Server 地址、namespace 等）。
- `internal/services/temporalc/temporalc.go` 封装了 Temporal 客户端的常用操作（如启动 Schedule、触发 Workflow、取消 Workflow 等）。
- 业务代码通过 `temporalc` 包与 Temporal 交互，不直接引用 Temporal SDK 的复杂接口。
- 任务调度配置表 `job_schedule` 存储定时任务定义，包括任务编码 (`job_code`)、工作流类型 (`workflow_type`)、调度策略 (`schedule_type`，如 INTERVAL/CRON)、状态 (`status`)、Temporal Schedule ID 等。
- 任务执行记录表 `job_execution` 记录每次执行，与 Temporal Workflow ID 关联。
- 新增任务调度时：
  1. 在 `internal/services/temporaljob/` 中实现对应的 Workflow 和 Activity（例如 `PrintCountWorkflow`）。
  2. 通过管理接口创建 `job_schedule` 记录，并调用 `Sync` 接口同步到 Temporal。
  3. `Trigger` 接口可用于立即触发一次任务执行，`Switch` 启用/停用定时。
  4. 任务执行记录由系统自动记录，前端可查看、重试或取消运行中的执行。
- 取消和重试操作通过 `POST /api/sys/job/execution/cancel` 和 `POST /api/sys/job/execution/retry` 实现。
