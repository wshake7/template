# Backend Admin 服务开发

## 场景与目标
- 适用场景：修改后台管理服务接口、路由、中间件、配置、基础服务时
- 目标：按项目既有结构快速落地改动，避免破坏启动链路

## 目录/文件位置
- 启动入口：`backend/admin/cmd/main.go`
- 配置定义：`backend/admin/config/`
- 配置文件：`backend/admin/etc/config.yaml`
- Fiber 初始化：`backend/admin/fiberc/`
- 路由注册：`backend/admin/router/router.go`
- 路由逻辑：`backend/admin/router/logic/`
- 服务初始化：`backend/admin/services/init.go`
- 数据层：`backend/admin/services/orm/`

## 核心依赖或组件
- `github.com/gofiber/fiber/v3`：HTTP 框架
- `go-common/viperc`：配置解析
- `go-common/log` + `zap`：日志
- `gorm` + `gorm/gen`：ORM 与查询代码
- `github.com/redis/rueidis`：Redis 客户端
- `github.com/click33/sa-token-go`：认证会话

## 启动链路
1. `cmd/main.go` 读取 `-f` 指定配置（默认 `./etc/config.yaml`）
2. 初始化日志与 services（Httpc/Orm/Redis/Auth/Geo/Asynq/Casbin 等）
3. 创建 Fiber App 并注册路由组（默认 `/api` 前缀）
4. 启动服务并等待优雅退出

## 操作步骤
1. 改接口逻辑：优先改 `router/logic/*.go`
2. 改路由映射：改 `router/*.go` 和 `router/router.go`
3. 改中间件：改 `fiberc/middleware/*.go`
4. 改配置项：补 `config/*.go` 结构体 + `etc/config.yaml` 字段
5. 改服务依赖：在 `services/init.go` 统一注入

## 常用命令
```bash
cd backend/admin

# 本地启动
go run ./cmd/main.go -f ./etc/config.yaml

# 执行测试
go test ./...

# 生成脚本（项目内置）
go run ./cmd/scripts/gen_imports
go run ./cmd/scripts/orm
```

## Handler 编写风格（以 router/logic/dict.go 为例）

### 1) 文件职责与组织
- 一个资源（如 `Dict`）集中在同一个 `router/logic/*.go` 文件内，按“类型/数据项”分段组织
- `Req*` 结构体与 Handler 方法就近定义，便于维护请求约束与业务逻辑的一致性
- 方法签名保持统一：读操作返回 `(*Data, error)`，写操作返回 `error`

### 2) 请求结构体与校验风格
- 入参统一使用 `json` + `binding` + `binding_msg`，错误文案直接面向业务语义
- `Update` 请求中的可选字段使用指针类型（如 `*string`、`*int32`），用于区分“未传”与“传零值”
- `Switch/Delete` 等简单动作使用最小请求体（如 `ID` + 状态位），避免冗余字段
- 分页列表统一复用 `v1.PagingRequest`，不重复定义分页 DTO

### 3) Repo 调用与更新模式
- 列表查询优先使用 `repo.XxxRepo.ListWithPaging(...)`，保持分页、排序、过滤能力一致
- 更新优先使用 `UpdateNoNilMap`（Patch 语义），结合指针字段实现“仅更新传入字段”
- 状态切换使用 `UpdateMap`，只更新目标列，避免误写其他字段
- 软删除统一使用 `SoftDelete`，与项目数据生命周期策略保持一致

### 4) 业务校验与事务边界
- 创建/更新涉及外键时，先用 `repo.XxxRepo.Exists(...)` 做关联存在性校验
- 需要联动删除或多表一致性时，使用 `orm.DB().Transaction(...)` 包裹
- 事务内步骤建议显式编号（先删子表，再删主表），降低维护成本

### 5) 错误处理与返回风格
- 用户可感知错误使用 `res.FailMsg("...")`（如“类型编码已存在”“字典类型不存在”）
- 通用异常统一返回 `res.FailDefault`，避免泄漏底层错误细节
- 对关键失败点补日志：`ctx.L().Error(..., zap.Error(err), ...)`，日志与返回语义解耦

### 6) Swagger 与注释风格
- 每个 Handler 前保持完整 Swagger 注释：`@Summary/@Description/@Tags/@Param/@Success/@Router`
- 分页接口 `@Success` 使用 `res.Response{data=gormc.PagingResult[models.Xxx]}`
- 若修改了 Swagger 注释，需在 `backend/admin` 执行 `make swagger`

## 注意事项
1. 新增路由要挂到 `router/router.go`，否则不会生效
2. 业务逻辑尽量放在 `router/logic`，避免路由文件过重
3. 配置新增字段要“结构体 + yaml”双向同步
4. 涉及 DB 变更时，注意 `orm/models`、`orm/repo` 与生成查询代码一致性
