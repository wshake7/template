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

# 更新 Swagger 文档（修改 Swagger 注释后必须执行）
make swagger
```

## Handler 编写风格（以 router/logic/dict.go 为例）

### 1) 文件职责与组织
- 一个资源（如 `Dict`）集中在同一个 `router/logic/*.go` 文件内，按业务子段（类型/数据项）分段组织
- `Req*` 请求结构体与 Handler 方法就近定义，便于维护请求约束与业务逻辑的一致性
- 方法签名保持统一：读/列表操作返回 `(*Data, error)`，写操作仅返回 `error`

### 2) 请求结构体与校验风格
- 入参统一使用 `json` + `binding` + `binding_msg`，错误文案直接面向业务语义，不暴露底层
- Create 请求结构体使用普通字段（`string`、`bool`、`int32`），Update 请求的可选字段使用**指针类型**（`*string`、`*int32`、`*bool`），用于区分"未传"与"传零值"
- Switch/Delete 等简单动作使用最小请求体（如 `ID` + 状态位），避免冗余字段
- **批量操作**（批量删除/批量复制）使用 `IDs []uint64` 或 `EntryIds []uint64` 字段，`binding:"required,min=1"` 校验至少一项
- 分页列表统一复用 `v1.PagingRequest`，不重复定义分页 DTO

```go
type ReqDictTypeBatchDelete struct {
    IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择字典类型,min=至少选择一项"`
}

type ReqDictEntryBatchCopy struct {
    EntryIds     []uint64 `json:"entryIds" binding:"required,min=1" binding_msg:"required=请选择字典项,min=至少选择一项"`
    TargetTypeId uint64   `json:"targetTypeId" binding:"required" binding_msg:"required=目标字典类型不能为空"`
}
```

### 3) Repo 调用与更新模式
- 列表查询优先使用 `repo.XxxRepo.ListWithPaging()`，保持分页、排序、过滤能力一致
- Create/Update/Delete 时需手动填充审计字段（`CreatedBy`/`UpdatedBy`）：
  ```go
  operationID := ctx.SessionInfo.Id
  ```
- 更新优先使用 `UpdateNoNilMap`（Patch 语义），结合指针字段实现"仅更新传入字段"
- 状态切换使用 `UpdateMap`，显式指定 `UpdatedBy`  + 目标列，避免误写其他字段
- 软删除统一使用 `SoftDelete`，与项目数据生命周期策略保持一致

### 4) 业务校验与事务边界
- 创建/更新涉及外键时，先用 `repo.XxxRepo.Exists()` 做关联存在性校验
- **批量 Delete** 涉及子表联动时，使用 `orm.DB().Transaction()` 包裹，按"先删子表，再删主表"顺序执行：
  ```go
  err := orm.DB().Transaction(func(tx *gorm.DB) error {
      repo.SysDictEntryRepo.SoftDelete(ctx, tx.Where(entry.SysDictTypeId.In(ids...)))
      repo.SysDictTypeRepo.SoftDelete(ctx, tx.Where(dictType.ID.In(ids...)))
      return nil
  })
  ```
- **批量 Copy** 涉及"查询 + 创建"两步：
  1. `orm.DB().Where("id IN ?", ids).Find(&entries)` 查询源条目
  2. 遍历源条目，将关联 ID 替换为目标 ID，构造新 DTO 列表
  3. `repo.XxxRepo.BatchCreate()` 批量写入
  4. 自复制拦截：比较 `sourceEntries[0].SysDictTypeId == req.TargetTypeId`

### 5) 错误处理与返回风格
- 用户可感知错误使用 `res.FailMsg("...")`（如"类型编码已存在""目标字典类型不存在"）
- 通用异常统一返回 `res.FailDefault`，避免泄漏底层错误细节
- 对关键失败点补日志：`ctx.L().Error("描述", zap.Error(err), ...)`，日志与返回语义解耦
- `Switch/Del/Update` 等操作中若遇到 `gorm.ErrDuplicatedKey`，转译为 `res.FailMsg` 返回

### 6) Swagger 与注释风格
- 每个 Handler 前保持完整 Swagger 注释：`@Summary/@Description/@Tags/@Accept/@Produce/@Param/@Success/@Router`
- 分页列表接口 `@Success` 使用 `res.Response{data=gormc.PagingResult[models.Xxx]}`
- 简单写操作（create/update/del）使用 `res.Response` 无需 data
- 若修改了 Swagger 注释，需在 `backend/admin` 目录执行 `make swagger`

## 注意事项
1. 新增路由要挂到 `router/router.go` 或 `auth_router.RegisterRouters()` 下，否则不会生效
2. 业务逻辑尽量放在 `router/logic`，避免路由文件过重
3. router 层只做路由映射与 middleware 编排（如 `OperationLogMiddleware`）
4. 配置新增字段要"结构体 + yaml"双向同步
5. 涉及 DB 变更时，注意 `orm/models`、`orm/repo` 与生成查询代码一致性
