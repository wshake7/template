# Backend Swagger 编写风格

## 场景与目标

- 适用场景：新增或修改 `backend/admin` API handler 的 Swagger 注释。
- 目标：保持接口文档完整、中文、与实际路由一致。

## 必填注释

每个 API 至少包含：

- `@Summary`
- `@Description`
- `@Tags`
- `@Accept json`
- `@Produce json`
- `@Param`（如有请求体、header、query、path）
- `@Success`
- `@Router`

## 响应格式

统一使用 `res.Response` 作为外层封装。

```go
// @Success 200 {object} res.Response "成功"
// @Success 200 {object} res.Response{data=ActualDataModel} "成功"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysRole]} "成功"
```

## 请求参数

```go
// @Param req body ReqStructName true "参数描述"
// @Param token header string true "登录 token"
```

## 路由格式

- 格式：`@Router /api/path [method]`
- 路径必须与实际注册路由一致。

```go
// @Router /api/sys/role/list [post]
```

## 修改步骤

1. 修改 handler 注释，摘要、描述、参数说明使用中文。
2. 检查 `@Tags` 是否对应业务模块。
3. 检查 `@Router` 路径和 method 是否与路由注册一致。
4. 在 `backend/admin` 执行 `make swagger`。
5. 再执行 `go fix ./...`、`go test ./...`。

## 注意事项

1. 只要改 Swagger 注释，就必须重新生成 `backend/admin/docs/*`。
2. 参数描述要说明用户能理解的业务含义。
3. `binding` 中复杂校验应在 `@Description` 中说明。
4. 引用泛型响应时确认包名与实际 import 一致。
