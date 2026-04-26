# Backend Swagger 编写风格规范

本文档定义了后端 `admin` 服务中 API 处理程序（Handler）的 Swagger 注释编写规范，确保文档的一致性和可读性。

## 核心原则

1. **完整性**：每个 API 必须包含 `@Summary`, `@Description`, `@Tags`, `@Router`。
2. **一致性**：使用统一的响应封装格式和参数描述风格。
3. **语言**：所有摘要、描述和参数说明必须使用 **中文**。

## 注释项详解

### 1. 基础信息
- **@Summary**: 简短摘要，描述 API 的主要功能。
- **@Description**: 详细说明，描述 API 的具体行为、前置条件或副作用。
- **@Tags**: 分组标签，通常对应逻辑模块名（如 `Role`, `Account`, `Encrypt`）。

### 2. 请求与响应格式
- **@Accept**: 请求数据格式，统一使用 `json`。
- **@Produce**: 响应数据格式，统一使用 `json`。

### 3. 参数 (@Param)
- **Body 参数**:
  ```go
  // @Param req body ReqStructName true "参数描述"
  ```
- **Header 参数**:
  ```go
  // @Param token header string true "登录 token"
  ```
- **Query/Path 参数**: 根据实际情况使用 `query` 或 `path`。

### 4. 成功响应 (@Success)
统一使用 `res.Response` 作为外层封装。
- **带数据的响应**:
  ```go
  // @Success 200 {object} res.Response{data=ActualDataModel} "成功"
  ```
- **带分页数据的响应**:
  ```go
  // @Success 200 {object} res.Response{data=gorm.PagingResult[models.SysRole]} "成功"
  ```
- **无数据的响应**:
  ```go
  // @Success 200 {object} res.Response "成功"
  ```

### 5. 路由 (@Router)
- 格式：`@Router /api/path [method]`
- 示例：`@Router /api/role/list [get]`

## 示例代码

```go
// @Summary 获取角色分页列表
// @Description 分页查询角色信息
// @Tags Role
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gorm.PagingResult[models.SysRole]} "成功"
// @Router /api/role/list [get]
func (h *RoleHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gorm.PagingResult[models.SysRole], error) {
    // ...
}
```

## 文档更新命令

当 AI 或开发者修改了 API 处理程序中的 Swagger 注释后，**必须**运行以下命令以同步更新生成的文件：

```bash
# 在 backend/admin 目录下执行
make swagger
```

该命令会重新生成 `internal/docs/docs.go`, `internal/docs/swagger.json` 和 `internal/docs/swagger.yaml`。

## 注意事项
- 确保引用的 `models` 和 `res` 包路径正确。
- 参数描述中尽量明确字段含义。
- 路由路径应与实际注册的路由保持一致。
- **改动必更新**：只要修改了注释内容，必须执行 `make swagger` 并提交生成的文件。
- **校验反馈**：如果 `binding` 标签中包含复杂校验（如 `oneof`），应在 `@Description` 中同步说明，方便前端理解报错原因。
