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

## 注意事项
1. 新增路由要挂到 `router/router.go`，否则不会生效
2. 业务逻辑尽量放在 `router/logic`，避免路由文件过重
3. 配置新增字段要“结构体 + yaml”双向同步
4. 涉及 DB 变更时，注意 `orm/models`、`orm/repo` 与生成查询代码一致性
