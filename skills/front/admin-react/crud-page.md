# Admin CRUD 页面开发

## 场景与目标

- 适用场景：开发后台列表页，包含搜索、分页、新增、编辑、删除、批量操作。
- 目标：统一 `API -> usePagination -> ProTable -> ModalForm` 的实现方式。

## 文件位置

- 路由页面：`front/apps/admin-react/src/routes/`
- API 定义：`front/apps/admin-react/src/api/`
- 分页类型：`front/apps/admin-react/src/domains/page.ts`
- 表单工具：`front/apps/admin-react/src/utils/zod.ts`
- 提示工具：`front/apps/admin-react/src/utils/antd.ts`

## 推荐顺序

1. 读一个相近 CRUD 页面和对应 API 文件。
2. 定义或补齐 `src/api/<resource>.ts` 类型和 API 对象。
3. 新增路由页面并配置 `staticData.menu`。
4. 用 `usePagination` 接入列表接口。
5. 定义 zod schema、默认值和 `useZodForm`。
6. 配置 `ProTable` columns、分页、reload、toolbar。
7. 配置新增/编辑 `ModalForm`。
8. 补删除、批量删除、启停等操作。
9. 执行 `vp staged`。

## 分页模式

```typescript
const {
  data,
  total,
  page,
  pageSize,
  loading,
  send,
} = usePagination(
  (nextPage, nextPageSize) => {
    const params: Record<string, unknown> = {
      page: nextPage,
      pageSize: nextPageSize,
      orderBy: 'id desc',
    }
    if (searchText.trim()) {
      params.query = JSON.stringify({
        $or: [
          { code__icontains: searchText.trim() },
          { name__icontains: searchText.trim() },
        ],
      })
    }
    return ResourceApi.resourceList(params)
  },
  {
    initialData: { total: 0, items: [] },
    initialPage: 1,
    initialPageSize: DEFAULT_PAGE_SIZE,
    total: response => response.data?.total ?? 0,
    data: response => response.data?.items ?? [],
    watchingStates: [searchText],
    debounce: [500],
  },
)
```

## ProTable 约定

1. `rowKey` 使用稳定主键，如 `id`。
2. `search={false}`，搜索放 toolbar，避免两套查询状态。
3. 分页状态与 `usePagination` 绑定。
4. `options.reload` 调用 `send()`。
5. columns 使用 `useMemo`，依赖包含分页状态和操作回调。
6. 序号列手动计算：`(page - 1) * pageSize + index + 1`。
7. 删除等高危操作使用 `Popconfirm`。

## ModalForm + Zod 约定

```typescript
const ResourceSchema = z.object({
  code: z.string('请输入资源编码').min(1, '请输入资源编码'),
  name: z.string('请输入资源名称').min(1, '请输入资源名称'),
  isEnabled: z.boolean().default(true),
  remark: z.string().default(''),
})

const resourceDefaults = ResourceSchema.partial().parse({})
type ResourceFormValues = z.infer<typeof ResourceSchema>
```

1. 默认值从 zod schema 派生。
2. 关闭弹窗时清理 `editing` 和表单。
3. 保存成功后提示、关闭弹窗、刷新列表。
4. Update 请求可传部分字段时，与后端指针字段对齐。

## 搜索和批量操作

- 搜索使用 `Input.Search` 或 toolbar 内控件。
- 搜索文本加入 `watchingStates`，配合 `debounce`。
- 多字段模糊搜索用 `$or` 和 `__icontains`。
- 批量操作用 `rowSelection.selectedRowKeys` 维护本地选中状态。
- 批量成功后清空选中状态并刷新列表。

## 注意事项

1. 前端 `PagingRequest` 对应后端 `v1.PagingRequest`。
2. `query` JSON 语法见 `skills/backend/orm-query.md`。
3. 长文本列使用 `ellipsis: true`。
4. 成功提示保持简短：`保存成功`、`删除成功`、`操作成功`。
