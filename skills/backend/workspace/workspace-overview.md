# Backend 工作区总览

## 场景与目标
- 适用场景：AI 需要理解 backend 多模块结构、依赖关系、启动方式时
- 目标：快速定位应该改哪个模块，避免在错误模块改动

## 目录/模块结构
- 根目录：`backend/`
- 工作区文件：`backend/go.work`
- 当前工作区模块：
  - `backend/admin`
  - `backend/go-common`
  - `backend/orm-crud/api`
  - `backend/orm-crud/gorm`
  - `backend/orm-crud/pagination`
  - `backend/sa-token/rueidis`

## 模块职责
- `admin`：后台服务主应用（Fiber 路由、中间件、业务逻辑、服务初始化）
- `go-common`：公共工具与基础能力（日志、配置、加密、ID、工具函数）
- `orm-crud/*`：ORM CRUD 能力、分页/过滤/排序、proto API 定义
- `sa-token/rueidis`：sa-token 与 rueidis 的桥接封装

## 操作步骤
1. 先看 `backend/go.work` 明确本次任务落在哪个 module
2. 进入目标模块目录执行命令，避免在 `backend/` 根直接误跑
3. 涉及跨模块改动时，先确认 import 路径和 go.work 已包含对应模块

## 常用命令
```bash
# 查看 backend 工作区模块
type backend/go.work

# 进入 admin 模块开发
cd backend/admin
go run ./cmd/main.go -f ./etc/config.yaml

# 在目标模块运行测试（示例）
cd backend/go-common
go test ./...
```

## 注意事项
1. backend 是 Go 多模块工作区，不是单一 go.mod 项目
2. 跨模块改动要优先保持 `go.work` 与 import 一致
3. 启动/测试应在具体模块目录执行，减少环境误差
