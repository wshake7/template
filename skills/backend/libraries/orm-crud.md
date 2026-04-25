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

# api 模块测试
cd ../api
go test ./...
```

## PagingRequest query 过滤语法

`PagingRequest` 的 `query` 字段是一个 JSON 字符串，`ListWithPaging` 会自动解析为 SQL WHERE 条件。详细参考 `skills/backend/libraries/orm-query.md`。

### 基础语法
- 无显式操作符时默认执行 `=`（等于）
- 操作符通过双下划线 `__` 写在字段名后

```jsonc
// 等于
{"sysDictTypeId": 1}
// 等价于
{"sysDictTypeId__eq": 1}

// 大于等于
{"createTime__gte": "2024-01-01"}

// 包含（不区分大小写）
{"userName__icontains": "张"}
```

### 常见操作符
| 操作符 | 含义 | SQL |
|--------|------|-----|
| `__eq` | 等于 | `=` |
| `__ne` | 不等于 | `!=` |
| `__gt` / `__gte` | 大于 / 大于等于 | `>` / `>=` |
| `__lt` / `__lte` | 小于 / 小于等于 | `<` / `<=` |
| `__in` | 在集合中 | `IN (...)` |
| `__not_in` | 不在集合中 | `NOT IN (...)` |
| `__contains` | 包含 | `LIKE '%val%'` |
| `__icontains` | 包含（不区分大小写） | `ILIKE '%val%'` |
| `__startswith` | 前缀匹配 | `LIKE 'val%'` |
| `__endswith` | 后缀匹配 | `LIKE '%val'` |
| `__isnull` | 为空 | `IS NULL` |
| `__not_isnull` | 不为空 | `IS NOT NULL` |
| `__range` | 范围 | `BETWEEN a AND b` |

### 组合条件
- 顶层数组默认等价于 `$and`
- 支持 `$and` / `$or` 嵌套

```jsonc
// AND 组合：sysDictTypeId=1 且 isEnabled=true
{"$and": [{"sysDictTypeId": 1}, {"isEnabled": true}]}

// OR 组合：sysDictTypeId=1 或 sysDictTypeId=2
{"$or": [{"sysDictTypeId": 1}, {"sysDictTypeId": 2}]}

// AND 嵌套 OR：sysDictTypeId=1 且 (isEnabled=true 或 其他条件)
{"$and": [
  {"sysDictTypeId": 1},
  {"$or": [{"isEnabled": true}, {"remark__icontains": "test"}]}
]}
```

### 前端使用示例
```typescript
const payload = {
  page: 1,
  pageSize: 10,
  query: JSON.stringify({ sysDictTypeId: selectedType.id }),
  orderBy: 'sort_order asc,id desc',
}
// POST 请求 body 中 query 是 JSON 字符串
// 后端会自动解析为 WHERE sys_dict_type_id = {id}
```

## 注意事项
1. `orm-crud` 属于基础设施层，避免耦合业务字段语义
2. 分页/过滤协议改动要关注兼容性
3. 生成代码文件改动前先确认对应脚本与来源文件
4. `query` 字段在 GET 请求时注意 URL 编码
