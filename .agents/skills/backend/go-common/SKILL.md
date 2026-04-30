# Skill: Backend Go Common

## 何时使用

当任务涉及 `backend/go-common` 的通用能力，包括配置、日志、结果包装、DTO、集合、mapper、加密、ID、字符串、转换、文件、重试、协程、IP/Geo、金额等工具时使用。

## 模块地图

```text
backend/go-common/
├── collection/        # bitmap、deque、sync map、tuple
├── dto/               # 通用 DTO，如分页
├── i18n/              # 字典/结构翻译辅助
├── log/               # Zap 日志初始化、encoder、writer
├── mapper/            # copier/enum 等映射转换
├── result/            # 标准结果结构
├── types/             # 通用类型
├── viperc/            # Viper 配置读取
└── utils/
    ├── encrypt/       # AES/RSA/加密服务
    ├── file/          # 文件工具
    ├── id/            # UUID、snowflake、sonyflake、machine id
    ├── ip_util/       # ip2region、GeoLite、qqwry
    ├── stringcase/    # snake/camel/kebab 等命名转换
    ├── trans/         # 类型和 JSON 转换
    └── ...            # retry、pool、promise、money、passwd 等
```

## 使用原则

- 这是跨模块基础库，修改时优先保持 API 小而稳定。
- 已有测试覆盖的包，改动时补充或更新对应测试。
- 工具函数优先保持无业务语义，不把 `backend/admin` 的业务规则下沉到这里。
- 对性能敏感或 unsafe 相关工具，先读现有测试和调用方，再改实现。
- 内置资产如 IP 数据库通过 `assets.go` 嵌入，不随意更换大文件。

## 常见任务流程

### 修改 mapper

1. 读 `mapper/interface.go`、`mapper/mapper.go`、`mapper/enum_converter.go`。
2. 搜索调用方：`rg "mapper\\." backend`。
3. 更新或新增 `mapper/*_test.go`。
4. 注意 pointer、nil、slice、enum converter 的边界。

### 修改 ID 或机器码

1. 读 `utils/id/README.md` 和相关实现。
2. 注意 Windows/Linux/Darwin/BSD/AIX 等平台文件的 build tags 或文件后缀。
3. 更新平台无关测试；平台特定测试只改能在当前环境验证的部分。

### 修改 stringcase 或 trans

1. 先读 README 和现有测试，保持现有命名转换语义。
2. 新增边界输入测试，例如空字符串、缩写、连续分隔符、数字、非 ASCII。
3. 只在确有必要时改变公开行为，并在日志或说明里记录原因。

### 修改配置或日志

1. 配置读取入口在 `viperc/viperc.go`，调用方常见于 `backend/admin/cmd/main.go`。
2. 日志入口在 `log/log.go`，writer/encoder 分别处理输出和编码。
3. 后端服务启动依赖这些基础能力，修改后至少在 `backend/admin` 做一次编译或测试验证。

## 命令

```bash
cd backend/go-common
go fix ./...
go vet ./...
go test ./...
```

## 验证

- 修改 `go-common` 后先跑本模块测试。
- 如果改动被 `backend/admin`、`orm-crud` 使用，继续在受影响模块跑 `go test ./...`。
- 涉及跨平台文件时，说明当前环境未覆盖的平台风险。
