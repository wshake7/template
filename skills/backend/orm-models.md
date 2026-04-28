# Backend ORM Models 与 Query

## 场景与目标

- 适用场景：修改 `backend/admin` 的数据库模型、查询代码、分页列表、增删改逻辑。
- 目标：统一 `router -> logic -> query/models -> gormc` 的数据库操作风格。

## 目录/文件位置

- 路由逻辑：`backend/admin/internal/router/logic/*.go`
- ORM 服务：`backend/admin/internal/services/orm.go`
- ORM 核心：`backend/admin/internal/services/orm/orm.go`
- 模型目录：`backend/admin/internal/services/orm/models/`
- 查询代码：`backend/admin/internal/services/orm/query/`
- 生成脚本：`backend/admin/cmd/scripts/orm/main.go`
- 生成模板：`backend/admin/cmd/scripts/orm/templates/`

## 当前调用链

1. 路由通过 `handler.Ctx*` 包装器进入 `logic/*.go`。
2. `logic` 层优先使用 `query.Xxx` 生成对象执行查询和写入。
3. 分页列表使用 `query.Xxx.PageWithPaging(req)`，入参复用 `v1.PagingRequest`。
4. 模型复用 `orm-crud/gormc/mixin` 的通用字段组合。
5. 错误统一转为 `res.FailDefault` 或明确业务提示。

## 常见操作

```go
func (*SysResourceHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysResource], error) {
    pagination, err := query.SysResource.PageWithPaging(req)
    if err != nil {
        return nil, res.FailDefault
    }
    return pagination, nil
}
```

```go
exprs := []field.AssignExpr{sysResource.UpdatedBy.Value(operationID)}
query.ExprAppendSelf(&exprs, req.Name, sysResource.Name.Value)
_, err := sysResource.Where(sysResource.ID.Eq(req.ID)).UpdateSimple(exprs...)
```

## Models 维护约定

1. 每个模型文件通过 `init()` 把模型 append 到 `models.Models`。
2. `orm.New()` 按配置启用 `WithAutoMigrate(models.Models...)`。
3. 表名通过 `TableName()` 显式定义。
4. 通用字段优先复用 `orm-crud/gormc/mixin`。
5. 新增或修改模型后执行 ORM 生成脚本更新 `query/*.gen.go`。

## 修改步骤

1. 先在 `logic/xxx.go` 明确读写需求、请求参数和错误语义。
2. 需要新表或字段时修改 `models/*.go`，并确认已注册到 `models.Models`。
3. 执行 `go run ./cmd/scripts/orm` 更新 query 生成代码。
4. 回到 logic 层用 `query.Xxx` 字段表达式组织条件和更新。
5. 执行 `go fix ./...`、`go test ./...`。

## 注意事项

1. `logic` 层保持业务编排，不拼接裸 SQL 字段名。
2. 更新和删除必须带条件。
3. 变更字段名或类型后，同步检查前端 API 类型和页面表单。
4. 新增模型未注册到 `models.Models` 会导致迁移与生成遗漏。
