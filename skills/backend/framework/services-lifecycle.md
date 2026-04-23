# Backend Services 生命周期与依赖编排

## 场景与目标
- 适用场景：修改 `backend/admin/services` 下服务初始化、健康检查、关闭流程、服务依赖顺序时
- 目标：确保服务按正确顺序注入 Fiber 服务容器，并在启动/终止阶段行为一致

## 目录/文件位置
- 服务装配入口：`backend/admin/services/init.go`
- 服务包装层：`backend/admin/services/*.go`
- 各服务实现：
  - `backend/admin/services/orm/`
  - `backend/admin/services/redisc/`
  - `backend/admin/services/httpc/`
  - `backend/admin/services/asynq/`
  - `backend/admin/services/casbin/`

## 当前编排方式（按现状）
1. `services.New(conf)` 统一构建服务并 append 到 `conf.Fiber.Services`
2. 当前装配顺序：
   - `NewHttpc()`
   - `NewOrm(conf.Orm)`
   - `NewRedis(conf.Redis)`
   - `NewAuth(conf.Auth, redisc.Client)`
   - `NewGeo()`
   - `NewAsynq(conf.Redis)`
   - `NewCasbin(orm.Client.DB)`
3. 下游依赖上游的包级客户端：
   - `Auth` 依赖 `redisc.Client`
   - `Casbin` 依赖 `orm.Client.DB`

## 服务接口风格
- 每个服务封装为结构体，并实现统一方法：
  - `Start(ctx context.Context) error`
  - `String() string`
  - `State(ctx context.Context) (string, error)`
  - `Terminate(ctx context.Context) error`
- `State` 用于健康探针输出；`Terminate` 负责释放资源

## 各服务职责速览
- `orm.go`：初始化 DB 客户端，健康检查通过 `PingContext`
- `redis.go`：初始化 Redis 客户端，健康检查 `PING`
- `auth.go`：初始化 sa-token manager（基于 rueidis storage）
- `httpc.go`：初始化 HTTP 客户端并在关闭时释放
- `asynq.go`：初始化任务客户端，`Ping` 检查连通性
- `casbin.go`：初始化 casbin adapter/enforcer 并关闭 adapter
- `geo.go`：初始化 IP 地理库并设置全局 `ip_util.Client`

## 修改步骤（推荐）
1. 新增服务：先在 `services/xxx.go` 实现统一四个方法
2. 在 `services/init.go` 按依赖顺序注入 `conf.Fiber.Services`
3. 若依赖包级客户端（如 `orm.Client`），确认启动顺序在依赖之前
4. 增加必要的 `State` 与 `Terminate` 逻辑，避免“能启动不能关闭”

## 常用命令
```bash
cd backend/admin

# 启动服务观察 Start/Terminate 行为
go run ./cmd/main.go -f ./etc/config.yaml

# 回归测试
go test ./...
```

## 注意事项
1. `init.go` 顺序是关键，改顺序前先评估依赖链
2. 包级变量客户端（`orm.Client` / `redisc.Client`）使用方便但耦合高，改动要谨慎
3. 新服务必须实现 `Terminate`，避免资源泄漏
