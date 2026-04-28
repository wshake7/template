# Backend 工作区总览

## 场景与目标

- 适用场景：开始任何 `backend/**` 任务前。
- 目标：先判断目标 Go module，避免在错误目录运行命令或误改跨模块依赖。

## 工作区结构

- 工作区入口：`backend/go.work`
- 当前模块：
  - `backend/admin`
  - `backend/go-common`
  - `backend/orm-crud/api`
  - `backend/orm-crud/gormc`
  - `backend/orm-crud/pagination`
  - `backend/sa-token/rueidis`

## 模块职责

- `admin`：后台服务主应用，包含 Fiber 路由、中间件、业务逻辑、ORM、Swagger。
- `go-common`：公共工具和基础能力，如日志、配置、加密、ID、集合、转换、i18n。
- `orm-crud/api`：分页协议 proto 与生成代码。
- `orm-crud/gormc`：GORM client、repository、mixin、filter、sorting、分页适配。
- `orm-crud/pagination`：分页器与查询字符串转换能力。
- `sa-token/rueidis`：sa-token 与 rueidis 的桥接封装。

## 定位步骤

1. 先读 `backend/go.work` 或用 `rg --files backend/<module>` 定位模块。
2. 后端命令进入具体模块目录执行，不在 `backend/` 根目录直接跑。
3. 跨模块改动时确认 import 路径、`go.work` 和被影响模块。
4. 只改文档或 `skills/**` 时不需要执行 Go 校验。

## 常用命令

```bash
cd backend/admin
go run ./cmd/main.go -f ./etc/config.yaml
go fix ./...
go vet ./...
go test ./...
```

## 注意事项

1. 这是 Go 多模块工作区，不是单一 `go.mod` 项目。
2. 修改多个后端模块时，每个受影响模块都要按 `AGENTS.md` 校验。
3. 修改 `backend/admin` Swagger 注释后必须在 `backend/admin` 执行 `make swagger`。
