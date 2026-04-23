# Backend go-common 公共库

## 场景与目标
- 适用场景：后端需要复用公共能力（日志、配置、加密、ID、工具函数）时
- 目标：优先复用 `go-common`，减少业务模块重复实现

## 目录/文件位置
- 模块目录：`backend/go-common`
- 配置解析：`backend/go-common/viperc/`
- 日志能力：`backend/go-common/log/`
- 工具集合：`backend/go-common/utils/`
- 返回结构：`backend/go-common/result/`
- 类型定义：`backend/go-common/types/`

## 常用能力
- `viperc`：配置文件/环境读取
- `log`：zap 封装与日志初始化
- `utils/encrypt/*`：AES/RSA 与加密服务
- `utils/id/*`：uuid/snowflake/机器码等
- `utils/retry`、`utils/pool`、`utils/promise`：基础并发与重试
- `dto/page.go`：分页 DTO

## 操作步骤
1. 开发前先在 `go-common` 搜索是否已有同类能力
2. 若已有，直接复用并补测试
3. 若新增公共能力，放在对应子目录并保持命名风格一致
4. 在使用方模块（如 `admin`）补充调用代码与验证

## 常用命令
```bash
cd backend/go-common

# 运行测试
go test ./...

# 仅测某个工具包（示例）
go test ./utils/encrypt/...
```

## 注意事项
1. `go-common` 变更影响面大，优先保持向后兼容
2. 新增工具需附最小可运行测试，避免回归
3. 不要把业务特定逻辑放进 `go-common`
