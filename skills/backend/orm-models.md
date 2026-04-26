# Backend ORM 数据库操作与 Models 维护

## 场景与目标
- 适用场景：在 `internal/router/logic` 中新增或调整数据库读写逻辑，或修改 `internal/services/orm` 下模型/仓储能力时
- 目标：统一 `router -> logic -> orm/repo -> query/models` 的数据库操作风格，减少 SQL 行为偏差和回归风险

## 目录/文件位置
- 路由逻辑入口：`backend/admin/internal/router/logic/*.go`
- ORM 服务初始化：`backend/admin/internal/services/orm.go`
- ORM 核心：`backend/admin/internal/services/orm/orm.go`
- 模型目录：`backend/admin/internal/services/orm/models/`
- 查询代码：`backend/admin/internal/services/orm/query/`
- 仓储封装：`backend/admin/internal/services/orm/repo/`
- 代码生成脚本：`backend/admin/cmd/scripts/orm/main.go`

## 当前数据库调用链（按现状）
1. 路由层通过 `handler.Ctx*` 包装器进入 `internal/router/logic/*.go`
2. `logic` 层使用 `orm.DB()` 获取数据库实例，调用 `repo.XxxRepo` 执行读写
3. 条件构造优先使用 `query.Xxx`（`gorm/gen` 生成对象）而不是手写字符串字段
4. 返回错误统一转换为业务错误（如 `res.FailDefault` 或可读提示）

## 常见操作模式
1. 分页查询：`repo.XxxRepo.ListWithPagination(ctx.Context(), orm.DB(), req)`
2. 新增记录：`repo.XxxRepo.Create(ctx.Context(), orm.DB(), &models.Xxx{...})`
3. 按条件更新：`repo.XxxRepo.UpdateMap(map[field.Expr]any{...}, conds...)`
4. 软删除：`repo.XxxRepo.SoftDelete(ctx.Context(), orm.DB().Where(...))`
5. 查询字段：通过 `query.Xxx.<Field>` 明确列，避免硬编码字段名

## UpdateMap 约定（重点）
1. `repo/*_repo.go` 中统一接收 `map[field.Expr]any`
2. 由仓储层将 `field.Expr` 转换为列名字符串后执行 `Updates`
3. `logic` 层只负责组织“字段表达式 + 条件”，不在上层拼接列名
4. 条件模型要与仓储模型一致，避免误用其它 `query` 对象

## Models 维护约定
1. 每个模型文件通过 `init()` 把模型 append 到 `models.Models`
2. `orm.New()` 按配置启用 `WithAutoMigrate(models.Models...)`
3. 需要生成查询代码时，执行 `go run ./cmd/scripts/orm`
4. 表名通过 `TableName()` 显式定义，字段尽量复用 `orm-crud/gorm/mixin`

## 修改步骤（推荐）
1. 先在 `internal/router/logic/xxx.go` 明确读写需求、请求参数和错误语义
2. 复用现有 `repo` 能力；缺失方法再补到 `internal/services/orm/repo/*.go`
3. 需要新字段/新表时同步修改 `models/*.go` 并确认已注册到 `models.Models`
4. 运行生成脚本更新 `query/*.gen.go`，再回归 `logic` 调用点
5. 启动服务或执行测试验证行为与返回结构

## 常用命令
```bash
cd backend/admin

# 启动服务（可联动自动迁移与 query 默认实例）
go run ./cmd/main.go -f ./etc/config.yaml

# 生成 orm query 代码
go run ./cmd/scripts/orm

# 基础修复与回归
go fix ./...
go test ./...
```

## 注意事项
1. `logic` 层保持“编排为主”，不要堆积底层 SQL 细节
2. 更新与删除务必带条件，避免全表更新/误删
3. 变更字段名或类型后，要同步检查 `query/*.gen.go` 与 `repo.UpdateMap` 调用点
4. 新增模型未注册到 `models.Models` 会导致迁移与生成遗漏
5. `sys_login_log.go` 当前仍是占位文件，启用前需补全模型并注册
