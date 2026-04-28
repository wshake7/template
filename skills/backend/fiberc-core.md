# Backend FiberC 核心结构

## 场景与目标

- 适用场景：修改 `fiberc` 初始化、全局中间件、错误处理、上下文或优雅关闭。
- 目标：在不破坏全局行为的前提下扩展 FiberC。

## 目录/文件位置

- 主入口：`backend/admin/internal/fiberc/fiber.go`
- App 生命周期：`backend/admin/internal/fiberc/app.go`
- 自定义上下文：`backend/admin/internal/fiberc/handler/ctx.go`
- 中间件：`backend/admin/internal/fiberc/middleware/*.go`
- 统一响应：`backend/admin/internal/fiberc/res/res.go`

## 关键链路

1. `NewFiber(conf)` 调用初始化逻辑创建 Fiber App。
2. 注册启动和关闭 hooks。
3. 挂载全局中间件：recover、cors、pprof、healthcheck、metrics、prometheus、日志、trace 等。
4. 设置优雅关闭，监听 `SIGINT/SIGTERM`。
5. `App.Start()` 中启动监听并等待关闭信号。

## 全局错误处理

- 业务错误通常以 HTTP 200 返回统一响应结构。
- Fiber 内置错误如 404 保持对应状态码。
- 实现 `res.Response` 的错误直接输出业务结构。
- 未知错误统一输出 `domains.JsonErr` 并记录日志。

## 修改步骤

1. 新增全局中间件时，在 Fiber 初始化处按顺序追加。
2. 改错误语义时，只在全局 `ErrorHandler` 统一调整。
3. 改优雅关闭时，集中在关闭函数和 hooks 中处理。
4. 新增上下文能力时，优先扩展 `handler.Ctx`。
5. 修改完成后在 `backend/admin` 执行 `go fix ./...`、`go test ./...`。

## 注意事项

1. 中间件顺序是行为约束，调整前先评估日志、追踪、恢复、鉴权影响。
2. `ErrorHandler` 是全局语义，不为单接口临时改。
3. 信号处理必须保证关闭通道最终释放，避免进程悬挂。
