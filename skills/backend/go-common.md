# Backend go-common 公共库

## 场景与目标

- 适用场景：后端需要复用或新增公共能力，如日志、配置、加密、ID、工具函数。
- 目标：优先复用 `go-common`，减少业务模块重复实现，同时避免把业务语义下沉到公共库。

## 目录/文件位置

- 模块目录：`backend/go-common`
- 配置解析：`backend/go-common/viperc/`
- 日志能力：`backend/go-common/log/`
- 工具集合：`backend/go-common/utils/`
- 返回结构：`backend/go-common/result/`
- 类型定义：`backend/go-common/types/`

## 常用能力

- `viperc`：配置文件和环境读取。
- `log`：zap 封装与日志初始化。
- `utils/encrypt/*`：AES/RSA 与加密服务。
- `utils/id/*`：uuid、snowflake、机器码等。
- `utils/retry`、`utils/pool`、`utils/promise`：基础并发与重试。
- `dto/page.go`：分页 DTO。

## 修改步骤

1. 开发前先在 `go-common` 搜索是否已有同类能力。
2. 已有能力优先复用；缺少测试时补最小测试。
3. 新增公共能力时放在对应子目录，保持命名和包边界一致。
4. 使用方模块同步改调用代码。
5. 在 `backend/go-common` 和受影响使用方模块分别执行测试。

## 常用命令

```bash
cd backend/go-common
go test ./...

# 示例：只测加密工具
go test ./utils/encrypt/...
```

## 注意事项

1. `go-common` 影响面大，优先保持向后兼容。
2. 新增工具需附最小可运行测试。
3. 不把 admin、权限、菜单、具体业务流程放进公共库。
4. 改公共 API 时同步检查所有 import 使用点。
