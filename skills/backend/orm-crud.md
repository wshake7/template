# Backend orm-crud 能力

## 场景与目标
- 适用场景：需要实现通用 CRUD、分页、过滤、排序或相关代码生成时
- 目标：复用 `orm-crud` 现有能力，减少重复手写 SQL/分页逻辑

## 目录/模块结构
- `backend/orm-crud/gorm`：GORM 客户端、repository、mixin、filter/sorting
- `backend/orm-crud/pagination`：分页器、查询参数转换器
- `backend/orm-crud/api`：pagination proto 与生成代码

## 关键能力
- `gorm/mixin/*`：通用模型字段组合（时间戳、租户、软删、版本等）
- `gorm/filter`：结构化过滤处理器
- `gorm/sorting`：结构化排序处理器
- `pagination/filter`：查询字符串到过滤条件转换
- `pagination/paginator`：offset/page/token 三类分页

## 操作步骤
1. 新增列表查询优先组合 `filter + sorting + paginator`
2. 新模型优先复用 `gorm/mixin`，避免重复字段定义
3. 协议层需要分页结构时，优先复用 `api/protos/pagination`
4. 改动后在对应子模块执行测试

## 常用命令
```bash
# gorm 模块测试
cd backend/orm-crud/gorm
go test ./...

# pagination 模块测试
cd ../pagination
go test ./...

# api 模块构建/测试
cd ../api
go test ./...
```

## 注意事项
1. `orm-crud` 属于基础设施层，避免耦合业务字段语义
2. 分页/过滤协议改动要关注兼容性
3. 生成代码文件改动前先确认对应脚本与来源文件
