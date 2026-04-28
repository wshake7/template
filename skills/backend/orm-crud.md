# Backend orm-crud 基础设施

## 场景与目标

- 适用场景：维护通用分页、过滤、排序、GORM client、mixin 或分页 proto。
- 目标：复用 `orm-crud` 基础设施，避免业务层重复实现分页和查询解析。

## 模块结构

- `backend/orm-crud/gormc`：GORM client、repository、mixin、filter、sorting、分页适配。
- `backend/orm-crud/pagination`：分页器、查询字符串转换器。
- `backend/orm-crud/api`：pagination proto 与生成代码。

## 关键能力

- `gormc/mixin/*`：通用模型字段，如时间戳、软删、启用状态、审计字段、排序字段。
- `gormc/filter`：结构化过滤处理器。
- `gormc/sorting`：结构化排序处理器。
- `gormc/pagination`：page、offset、token 分页适配。
- `pagination/filter`：查询字符串到过滤条件转换。
- `pagination/paginator`：offset/page/token 三类分页器。
- `api/protos/pagination/v1/pagination.proto`：分页请求协议来源。

## 修改步骤

1. 新增列表能力时，优先复用 `PagingRequest`、filter、sorting、paginator。
2. 新模型字段优先复用 `gormc/mixin`。
3. 修改 proto 后，按 `backend/orm-crud/api` 的生成流程更新 `gen/go`。
4. 改基础能力后，在受影响子模块分别执行测试。

## 常用命令

```bash
cd backend/orm-crud/gormc
go test ./...

cd ../pagination
go test ./...

cd ../api
go test ./...
```

## 注意事项

1. `orm-crud` 是基础设施层，不写 admin 业务语义。
2. 分页、过滤协议改动要关注前后端兼容。
3. 生成代码改动前先确认来源文件和脚本。
4. `query` 字段通过 HTTP 传递时注意 JSON 字符串和 URL 编码。
