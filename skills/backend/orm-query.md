# ORM 分页过滤 Query 语法

## 场景与目标

- 适用场景：前端构造 `PagingRequest.query`，或后端维护分页过滤解析能力。
- 目标：统一 JSON 过滤条件写法，保证前后端列表筛选一致。

## 使用位置

- 前端：`front/apps/admin-react/src/domains/page.ts` 的 `PagingRequest.query`
- 后端：`backend/orm-crud/gormc/filter`、`backend/orm-crud/pagination/filter`
- Admin 列表：`query.Xxx.PageWithPaging(req)`

## JSON 基础语法

`query` 传 JSON 字符串。字段名可追加 `__操作符`，不写操作符时默认等于。

```json
{ "deptId": 1 }
```

等价于：

```json
{ "deptId__eq": 1 }
```

常用示例：

```json
{ "entryTime__gte": "2024-01-01" }
```

```json
{ "userName__icontains": "张" }
```

## 逻辑组合

顶层数组默认是 `$and`：

```json
[
  { "deptId": 1 },
  { "entryTime__gte": "2024-01-01" },
  { "userName__icontains": "张" }
]
```

显式 `$and`：

```json
{
  "$and": [
    { "deptId": 1 },
    { "entryTime__gte": "2024-01-01" }
  ]
}
```

`$or`：

```json
{
  "$or": [
    { "deptId": 1 },
    { "deptId": 2 },
    { "userName__icontains": "张" }
  ]
}
```

嵌套组合：

```json
{
  "$and": [
    { "deptId": 1 },
    {
      "$or": [
        { "entryTime__gte": "2024-01-01" },
        { "userName__icontains": "张" }
      ]
    }
  ]
}
```

## 常用操作符

| 操作符 | 含义 | 示例 |
| --- | --- | --- |
| `eq` | 等于 | `{ "name__eq": "tom" }` |
| `not` | 不等于 | `{ "name__not": "tom" }` |
| `in` | 在集合内 | `{ "id__in": [1, 2] }` |
| `not_in` | 不在集合内 | `{ "id__not_in": [1, 2] }` |
| `gt` / `gte` | 大于 / 大于等于 | `{ "createdAt__gte": "2024-01-01" }` |
| `lt` / `lte` | 小于 / 小于等于 | `{ "createdAt__lt": "2025-01-01" }` |
| `range` | 区间 | `{ "createdAt__range": ["2024-01-01", "2024-12-31"] }` |
| `isnull` | 为空 | `{ "deletedAt__isnull": true }` |
| `not_isnull` | 不为空 | `{ "deletedAt__not_isnull": true }` |
| `contains` | 包含，区分大小写 | `{ "name__contains": "A" }` |
| `icontains` | 包含，不区分大小写 | `{ "name__icontains": "a" }` |
| `startswith` / `istartswith` | 前缀匹配 | `{ "code__startswith": "sys" }` |
| `endswith` / `iendswith` | 后缀匹配 | `{ "code__iendswith": "log" }` |
| `exact` / `iexact` | 精确 LIKE | `{ "name__iexact": "tom" }` |

## 日期提取操作符

日期时间字段支持按维度筛选：

- `date`
- `year`
- `month`
- `day`
- `week`
- `quarter`
- `hour`
- `minute`
- `second`

示例：

```json
{ "createdAt__year": "2026" }
```

## 前端构造示例

```typescript
params.query = JSON.stringify({
  $or: [
    { code__icontains: searchText.trim() },
    { name__icontains: searchText.trim() },
  ],
})
```

## 注意事项

1. 前端通过 HTTP 传递时，`query` 是 JSON 字符串，不是对象。
2. 日期区间用 `range` 时注意边界；同一天查询可能需要传完整时间范围。
3. 字段名要与后端 query/model 可识别字段一致。
4. 复杂逻辑优先写显式 `$and` / `$or`，避免阅读歧义。
