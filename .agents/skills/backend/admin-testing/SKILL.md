# Skill: Backend Admin Testing

## 何时使用

当任务涉及 `backend/admin` 的单元测试、集成测试、覆盖率提升、测试架构选择、mock 策略、测试辅助设施或回归验证时使用。

特别适合以下场景：
- 给 `internal/router/logic/**` 新增测试
- 把难测 handler 改造成可测结构
- 决定该用 SQLite 还是 gomock
- 给 Temporal / Casbin / Redis / Auth / LoginLogger / DataPermission 相关分支补测试
- 跑覆盖率并定位低覆盖函数

## 当前测试总原则

### 1. 数据库逻辑统一测真实查询

`backend/admin` 当前约定是：
- **所有 DB 访问都走 `*query.Query`**
- **测试里优先用内存 SQLite 验证真实查询、事务和关联行为**

这意味着：
- `sys_user`
- `sys_language`
- `sys_api_log`
- `sys_login_log`
- `sys_resource_api`
- `sys_resource_menu`
- `sys_role`
- `sys_dict`
- `job_execution`
- `job_schedule`
- `account`

这些 logic handler 的数据库部分都不再通过 repo mock 来测。

### 2. 外部副作用统一 mock interface

以下能力不属于 DB 查询，保持 interface + gomock：
- `AuthService`
- `LoginLogger`
- `RedisCache`
- `CasbinService`
- `TemporalService`
- `DataPermissionService`

### 3. 不再使用的模式

- 不要为 admin logic 新增 repository mocks
- 不要回到“repo mock 测数据库行为”的老模式
- 不要把事务、多表 preload、权限拼装等 DB 行为抽象成一层测试用 repo 只为了 mock

## 关键测试文件

```text
backend/admin/internal/router/logic/
├── sqlite_test.go               # SQLite 建库与常用迁移 helper
├── testutil_test.go             # TestMain / 测试上下文辅助
├── account_test.go
├── encrypt_test.go
├── job_execution_test.go
├── job_schedule_test.go
├── job_schedule_pure_test.go
├── sys_api_log_test.go
├── sys_dict_test.go
├── sys_language_test.go
├── sys_login_log_test.go
├── sys_resource_api_test.go
├── sys_resource_menu_test.go
├── sys_role_test.go
└── sys_user_test.go
```

## 测试分层策略

### A. 纯 DB handler

例如：
- `SysUserHandler`
- `SysLanguageHandler`
- `SysApiLogHandler`
- `SysLoginLogHandler`
- `SysResourceMenuHandler`

写法：
1. 用 `setupSQLiteDB` 或 `mustMigrateXxx` 建库
2. seed 必要模型数据
3. `query.SetDefault(...)` 让全局 field builder 可用
4. 直接调用 handler 方法
5. 断言返回值与数据库状态

### B. DB + 外部副作用 handler

例如：
- `AccountHandler` = SQLite + `AuthService` + `LoginLogger`
- `SysResourceApiHandler` = SQLite + `CasbinService`
- `SysRoleHandler` = SQLite + `CasbinService`
- `SysDictHandler` = SQLite + `DataPermissionService`
- `JobExecutionHandler` = SQLite + `TemporalService`
- `JobScheduleHandler` = SQLite + `TemporalService`
- `EncryptHandler` = `RedisCache`（无 DB）

写法：
1. DB 部分照常走 SQLite
2. 只把外部依赖换成 gomock
3. 断言副作用调用次数、参数和错误分支

### C. 纯函数 / helper 测试

适合：
- normalize helper
- tree builder
- validation helper
- JSON 解析与 schedule spec 组装
- 类型转换 helper

优先把分支复杂但无外部依赖的逻辑拆成可直调 helper，并在 `*_pure_test.go` 或对应测试文件中直接覆盖。

## 现成测试基础设施

### SQLite 基础设施

`sqlite_test.go` 已提供：
- `setupSQLiteDB(t, tables...)`
- `mustMigrateResourceMenu(t)`
- `mustMigrateRole(t)`
- `mustMigrateDict(t)`
- `mustMigrateJob(t)`
- 以及各测试文件里的局部 `mustMigrateXxx(t)` helper

用途：
- 建立内存 SQLite
- AutoMigrate 当前 handler 需要的模型
- 返回 `*query.Query` 供 handler 注入

### TestMain / 默认 query

`testutil_test.go` 已处理全局 `query.*` field builder 初始化问题。

当你新增依赖 `query.<Model>` 全局表达式的测试时：
- 仍然要在具体测试里调用 `query.SetDefault(q.<Model>.UnderlyingDB())`
- 避免因为默认 DB 未设置导致字段表达式为 nil

### 测试上下文

复用 `newTestCtx(t)` 创建 handler 所需上下文。

如果测试依赖登录身份、操作人、语言或 session 数据，在该 helper 基础上最小化补字段，不要重新造一整套 Fiber 上下文。

## gomock 使用约定

### mock 来源

mock 文件位于 `backend/admin/internal/mock/**`，由 `internal/service/*.go` 上的 `go:generate` 维护。

### 常见写法

```go
ctrl := gomock.NewController(t)
defer ctrl.Finish()

mockTemporal := mock.NewMockTemporalService(ctrl)
mockTemporal.EXPECT().TriggerSchedule(gomock.Any(), "schedule-id", gomock.Any()).Return(nil)
```

### 什么时候需要重新生成 mock

只有在你修改了这些 interface 签名时才需要重新生成：
- `internal/service/auth_service.go`
- `internal/service/login_logger.go`
- `internal/service/redis_cache.go`
- `internal/service/casbin_service.go`
- `internal/service/temporal_service.go`
- `internal/service/data_permission_service.go`

生成方式示例：

```bash
cd backend/admin/internal/service
go generate ./...
```

如果本机没有 `mockgen`，先安装：

```bash
go install go.uber.org/mock/mockgen@latest
```

## 写新测试的推荐流程

### 给 logic handler 补测试

1. 先看对应 `*_test.go` 是否已有 `mustMigrateXxx` helper。
2. 确认这个分支属于：
   - DB 行为
   - 外部副作用
   - 纯 helper
3. DB 行为用 SQLite seed 数据复现。
4. 外部副作用用 gomock 打桩。
5. 尽量断言：
   - 返回值
   - 错误文本或错误类型
   - DB 最终状态
   - mock 是否被正确调用
6. 覆盖 happy path + 至少一个关键失败分支。

### 给事务逻辑补测试

优先验证事务结果，而不是 mock 事务本身。

例如：
- 创建失败是否回滚
- 关联表是否同步
- 删除时是否级联处理逻辑记录
- 更新父子关系后 tree path 是否正确

### 给权限/调度逻辑补测试

- 权限相关：优先测 `CasbinService` / `DataPermissionService` 调用与 DB 状态组合
- 调度相关：优先测 `TemporalService` 调用与 schedule / execution 状态变化组合

## 覆盖率工作流

### 快速看 logic 总覆盖率

```bash
cd backend/admin
go test ./internal/router/logic/... -cover
```

### 输出函数级覆盖率

```bash
cd backend/admin
go test ./internal/router/logic/... -coverprofile=/tmp/admin-logic.cover.out
go tool cover -func=/tmp/admin-logic.cover.out
```

### 全量回归

```bash
cd backend/admin
go test ./...
```

## 当前经验规则

- `account` 虽然有 `AuthService` / `LoginLogger`，但用户、角色、session 相关 DB 查询仍然直接测 SQLite。
- `encrypt` 没有 DB，直接 mock `RedisCache` 即可。
- `job_execution` / `job_schedule` 的可测性关键在 `TemporalService`，不要重新引入对底层 Temporal 客户端的直接依赖。
- `sys_dict` 的部分分页/过滤分支在 SQLite 下可能较脆，优先选择稳定断言，不要为了覆盖率写对 SQL 方言过度耦合的测试。
- `sys_resource_menu` / `sys_role` 的树结构、父子关系、关联集合校验适合做 helper + integration 双层测试。

## 常见坑

- 忘记 `query.SetDefault(...)`，导致全局 field builder 为空。
- 用错模型 ID 字段，很多 model 通过 `mixin.AutoIncrementID` 嵌入 `ID`。
- 只断言 `err != nil`，没有断言数据库状态，导致事务测试价值不足。
- 把 DB 行为 mock 掉后，测不出 preload、事务、唯一索引、软删除等真实问题。
- 修改 service interface 后忘记更新 `internal/mock/**`。

## 完成后的最低验证

```bash
cd backend/admin
go test ./internal/router/logic/...
go test ./...
```

如果任务目标是补覆盖率，再额外执行：

```bash
cd backend/admin
go test ./internal/router/logic/... -coverprofile=/tmp/admin-logic.cover.out
go tool cover -func=/tmp/admin-logic.cover.out
```
