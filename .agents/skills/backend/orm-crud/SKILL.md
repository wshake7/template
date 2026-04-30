# Skill: Backend ORM CRUD

## 何时使用

当任务涉及 `backend/orm-crud/**`，包括 GORM 客户端、通用 Repository、mixin、字段选择、过滤、排序、分页、pagination proto 或 OpenAPI/protobuf 辅助时使用。

## 模块地图

```text
backend/orm-crud/
├── gormc/                    # GORM CRUD 核心
│   ├── client.go             # GORM client 构建
│   ├── repository.go         # 通用 Repository
│   ├── options.go            # Option 配置
│   ├── mixin/                # 可复用模型字段
│   ├── filter/               # 结构化过滤
│   ├── sorting/              # 结构化排序
│   ├── pagination/           # GORM 层分页器
│   └── field/                # 字段选择
├── pagination/               # 分页请求表达式解析
│   ├── filter/               # JSON / AIP filter 转换器
│   ├── sorting/              # order_by 转换器
│   └── paginator/            # offset/page/token paginator
└── api/                      # pagination proto 与生成产物
    ├── protos/
    ├── buf.yaml
    └── gen/go/
```

## 设计边界

- `gormc` 关注 GORM 查询执行、模型 mixin、Repository、结构化过滤排序和分页落地。
- `pagination` 关注请求层分页参数、过滤表达式、排序表达式转换，不直接绑定具体业务 model。
- `api` 维护 protobuf 定义和生成代码，当前核心为 `pagination/v1/pagination.proto`。
- `backend/admin` 通过 `query.<Model>.PageWithPaging(req)` 使用这些能力。

## 过滤语法速查

- JSON 过滤支持 `field__operator`，无操作符默认等于。
- 支持 `$and`、`$or` 嵌套；顶层数组默认等价 `$and`。
- 常见操作符：`not`、`in`、`not_in`、`gte`、`gt`、`lte`、`lt`、`range`、`isnull`、`contains`、`icontains`、`startswith`、`endswith`、`exact`。
- JSON 字段可用 `field.nested__operator`。
- 详细规则见 `backend/orm-crud/pagination/filter/README.md`。

## 常见任务流程

### 修改分页或过滤行为

1. 先读 `pagination/filter/**` 的 converter 和测试。
2. 若是请求表达式解析问题，优先改 `pagination/filter` 或 `pagination/sorting`。
3. 若是 GORM 查询落地问题，再改 `gormc/filter`、`gormc/sorting` 或 `gormc/pagination`。
4. 为新增操作符或语法补充对应测试。

### 新增或修改 mixin

1. 在 `gormc/mixin` 新增字段结构或修改现有 mixin。
2. 检查 `backend/admin/internal/services/orm/models/**` 的复用方式。
3. 注意 GORM tag、JSON tag、软删除/唯一索引组合对生成代码的影响。

### 修改 proto/API

1. 修改 `api/protos/pagination/v1/pagination.proto`。
2. 检查 `api/buf.yaml`、`api/buf.gen.yaml`。
3. 运行对应生成命令或 Makefile。
4. 同步调整 `api/gen/go/**` 以及使用方。

## 命令

```bash
cd backend/orm-crud/gormc
go fix ./...
go vet ./...
go test ./...

cd ../pagination
go fix ./...
go vet ./...
go test ./...

cd ../api
go fix ./...
go vet ./...
go test ./...
```

## 验证

- 改哪个 Go module，就在对应目录执行 `go fix ./...`、`go vet ./...`、`go test ./...`。
- 涉及跨模块行为时，额外在 `backend/admin` 跑测试，确认业务调用没有断。
- 修改生成链路时要区分“源文件”和“生成产物”，不要只改会被覆盖的文件。
