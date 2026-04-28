# Backend Services 生命周期

## 场景与目标

- 适用场景：修改 `backend/admin/internal/services` 下的服务初始化、健康检查、关闭流程或依赖顺序。
- 目标：确保服务按正确顺序注入 Fiber 服务容器，并在启动/终止阶段行为一致。

## 目录/文件位置

- 服务装配入口：`backend/admin/internal/services/init.go`
- 服务包装层：`backend/admin/internal/services/*.go`
- ORM 服务：`backend/admin/internal/services/orm/`
- Redis 服务：`backend/admin/internal/services/redisc/`
- HTTP 客户端：`backend/admin/internal/services/httpc/`
- Asynq 服务：`backend/admin/internal/services/asynq/`
- Casbin 服务：`backend/admin/internal/services/casbin/`

## 当前装配顺序

1. `NewHttpc()`
2. `NewOrm(conf.Orm)`
3. `NewRedis(conf.Redis)`
4. `NewAuth(conf.Auth, redisc.Client)`
5. `NewGeo()`
6. `NewAsynq(conf.Redis)`
7. `NewCasbin(orm.Client.DB)`

## 服务接口风格

每个服务封装为结构体，并实现：

- `Start(ctx context.Context) error`
- `String() string`
- `State(ctx context.Context) (string, error)`
- `Terminate(ctx context.Context) error`

## 修改步骤

1. 新增服务时，在 `internal/services/xxx.go` 实现统一四个方法。
2. 在 `internal/services/init.go` 按依赖顺序注入 `conf.Fiber.Services`。
3. 如果依赖包级客户端，如 `orm.Client` 或 `redisc.Client`，确认依赖已先初始化。
4. `State` 要能反映健康状态，`Terminate` 要释放资源。
5. 在 `backend/admin` 执行 `go fix ./...`、`go vet ./...`、`go test ./...`。

## 注意事项

1. `init.go` 顺序是关键，改顺序前先评估依赖链。
2. 包级客户端使用方便但耦合高，替换时要同步所有调用点。
3. 新服务必须实现关闭逻辑，避免“能启动不能退出”。
