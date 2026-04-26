# Backend FiberC 核心结构

## 场景与目标
- 适用场景：需要修改 `fiberc` 初始化、全局中间件、错误处理、优雅关闭时
- 目标：在不破坏启动链路和全局行为的前提下进行扩展

## 目录/文件位置
- 主入口：`backend/admin/internal/fiberc/fiber.go`
- App 生命周期：`backend/admin/internal/fiberc/app.go`
- 自定义上下文：`backend/admin/internal/fiberc/handler/ctx.go`
- 中间件：`backend/admin/internal/fiberc/middleware/*.go`
- 统一响应：`backend/admin/internal/fiberc/res/res.go`

## 关键链路（当前实现）
1. `NewFiber(conf)` 调用 `initialize(conf)` 创建 Fiber App
2. `initializeHooks(app)`（代码里函数名为 `initializeHoos`）注册启动/关闭钩子
3. 挂载全局中间件：`recover -> cors -> pprof -> healthcheck -> metrics -> prometheus -> 日志 -> trace`
4. 设置 `app.done = gracefulShutdown(app)`，监听 `SIGINT/SIGTERM`
5. `App.Start()` 中 `Listen` 启动后阻塞等待 `<-app.done`

## Fiber 配置风格
- 通过 `config.Fiber` 映射到 `fiber.Config`，集中在 `initialize()`
- 使用 `fiber.NewWithCustomCtx` 注入自定义 `handler.Ctx`
- JSON 编解码统一使用 `sonic`
- `Services` 生命周期挂载到 Fiber 服务管理（Startup/Shutdown Context）

## 全局错误处理风格
- 默认响应码以 `200` 返回业务错误（Fiber 内置错误如 `404` 保持状态码）
- 若错误实现 `res.Response`，直接输出业务结构
- 其他未知错误统一返回 `domains.JsonErr` 并记录日志

## 操作步骤
1. 加全局中间件：在 `NewFiber()` 中按顺序追加 `app.Use`
2. 改错误处理：只在 `initialize()` 的 `ErrorHandler` 中统一改
3. 改优雅关闭：在 `gracefulShutdown()` 与 hooks 中扩展，不要散落到业务层
4. 新增上下文能力：优先扩展 `handler.Ctx`，保持 handler 风格一致

## 常用命令
```bash
cd backend/admin

# 启动服务观察中间件/错误处理行为
go run ./cmd/main.go -f ./etc/config.yaml

# 回归测试
go test ./...
```

## 注意事项
1. 中间件顺序是行为约束，调整顺序前要评估日志、追踪、恢复逻辑
2. `ErrorHandler` 属于全局语义，不能按单接口需求临时改
3. 若修改信号处理，需确认 `app.done` 一定会被关闭，避免进程悬挂
