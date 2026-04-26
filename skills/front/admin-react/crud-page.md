# Admin CRUD 页面开发风格

## 场景与目标
- 适用场景：开发后台管理的 CRUD 页面，包含列表/搜索/新增/编辑/删除/批量操作等功能
- 目标：以 `dict` 模块为蓝本，统一 CRUD 页面的代码结构与交互模式

## 目录/文件位置
- 路由页面：`src/routes/`（遵循 TanStack Router 文件路由）
- API 定义：`src/api/`（一个资源一个文件，如 `dict.ts`）
- 领域模型：`src/domains/page.ts`（分页请求/响应类型 + 常量）
- 后端 Handler：`backend/admin/router/logic/*.go`

## 核心依赖或组件
- `@tanstack/react-router`：类型安全路由，`staticData.menu` 配置菜单
- `@ant-design/pro-components`：`ProTable`、`ModalForm` 等
- `alova` / `usePagination`：分页数据请求与状态管理
- `zod` + `useZodForm`：表单验证
- `Splitter`（Ant Design）：主从布局（左列表、右详情/子列表）

## 操作步骤

### 1) API 层定义
- 一个资源一个文件，如 `src/api/dict.ts`
- 接口命名：`typeList`、`typeCreate`、`typeUpdate`、`typeSwitch`、`typeDel`
- 请求/响应类型就近定义，分页复用 `PagingRequest`/`PagingResult`
- 使用 `API.Post` 统一发起请求，`cacheFor: 0` 禁用缓存
- 最后通过 `DictApi` 对象聚合导出

```typescript
async function typeList(req: PagingRequest) {
  return await API.Post<Res<PagingResult<DictType>>>('/api/sys/dict/type/list', req, {
    cacheFor: 0,
  }).send()
}

export const DictApi = {
  typeList, typeCreate, typeUpdate, typeSwitch, typeDel,
}
```

### 2) Route 页面配置
```typescript
export const Route = createFileRoute('/_app/system/dict')({
  staticData: {
    menu: {
      name: '数据字典',
      menuType: 'menu',
    },
  },
  staleTime: 1000 * 60 * 2,
  component: RouteComponent,
})
```

### 3) ProTable + usePagination 分页模式
- `usePagination` 的第一个参数是请求函数，接收 `(nextPage, nextPageSize)`
- 请求函数内部用 `Record<string, unknown>` 构造参数对象
- `orderBy: 'sort_order asc,id desc'` 传递排序
- 配置 `initialData`、`total`、`data` 映射
- `watchingStates` 监听状态变化自动刷新
- `initialPageSize` 使用 `DEFAULT_PAGE_SIZE`（`src/domains/page.ts` 导出）

```typescript
const {
  data, total, page, pageSize, loading, update, send,
} = usePagination(
  (nextPage, nextPageSize) => {
    const params: Record<string, unknown> = {
      page: nextPage,
      pageSize: nextPageSize,
      orderBy: 'sort_order asc,id desc',
    }
    return API.Post<Res<PagingResult<DictType>>>('/api/sys/dict/type/list', params, { cacheFor: 0 })
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

`ProTable` 的 `pagination` 与 `usePagination` 联动：
```typescript
<ProTable<DictType>
  rowKey="id"
  search={false}
  columns={columns}
  dataSource={data}
  loading={loading}
  pagination={{
    showSizeChanger: true,
    current: page,
    pageSize,
    total,
    onChange: (nextPage, nextPageSize) => {
      update({ page: nextPage, pageSize: nextPageSize })
    },
  }}
  options={{ reload: () => send() }}
/>
```

### 4) 搜索模式（Input.Search + OR 模糊搜索）
- 工具栏用 `Input.Search` + `allowClear` + `onSearch`
- 搜索文本通过 `watchingStates` 联动，带 `debounce` 防抖
- 后端过滤使用 JSON 格式，`$or` + `__icontains` 实现多字段模糊搜索

```typescript
// 前端请求
if (searchText.trim()) {
  params.query = JSON.stringify({
    $or: [
      { typeCode__icontains: searchText.trim() },
      { typeName__icontains: searchText.trim() },
      { description__icontains: searchText.trim() },
    ],
  })
}
```

详见 `skills/backend/orm-query.md` 的 JSON 格式语法。

### 5) ModalForm + useZodForm 表单模式
- Zod Schema 定义验证规则，`z.string('提示信息').min(1, '必填')`
- `dictTypeDefaults = DictTypeSchema.partial().parse({})` 提供默认空值
- `useZodForm` 封装了表单创建、验证提交流程
- 编辑时 `form.setFieldsValue(...)` 回填，新增时 `form.setFieldsValue(dictTypeDefaults)`
- 提交成功后 `form.resetFields()` + 关闭弹窗 + `send()` 刷新

```typescript
const DictTypeSchema = z.object({
  typeCode: z.string('请输入类型编码').min(1, '请输入类型编码'),
  typeName: z.string('请输入类型名称').min(1, '请输入类型名称'),
  isEnabled: z.boolean().default(true),
  sortOrder: z.number().default(0),
  description: z.string().default(''),
})
const dictTypeDefaults = DictTypeSchema.partial().parse({})
export type DictTypeFormValues = z.infer<typeof DictTypeSchema>

const { form, rules, onFinish } = useZodForm<DictTypeFormValues>({
  schema: DictTypeSchema,
  async onSubmit(values) {
    if (editing) { await DictApi.typeUpdate({ id: editing.id, ...values }) }
    else { await DictApi.typeCreate(values) }
    gMessage.success('保存成功')
    setEditing(undefined)
    form.resetFields()
    setFormOpen(false)
    await send()
  },
})
```

### 6) Splitter 主从布局
- 左侧是主列表（如字典类型），右侧是从列表（如字典项）
- 默认 `defaultSize="50%"`，设置 `min`/`max` 限制拖拽范围

```typescript
<Splitter>
  <Splitter.Panel defaultSize="50%" min="25%" max="75%">
    <DictTypePanel ... />
  </Splitter.Panel>
  <Splitter.Panel>
    <DictEntryPanel ... />
  </Splitter.Panel>
</Splitter>
```

### 7) Columns 定义规范
- `useMemo` 包裹，依赖包含 `page`、`pageSize` 等分页变量
- 序号列：`(page - 1) * pageSize + index + 1`，禁止使用 `valueType: 'indexBorder'`（与自定义 render 冲突）
- 操作列的 `valueType: 'option'`
- 状态列使用 `statusTag` 渲染函数
- 长文本使用 `ellipsis: true`
- 操作使用 `<Popconfirm>` 包裹删除/高危操作

### 8) Row Selection + Batch Delete 模式
- `rowSelection` 配置 `selectedRowKeys` + `onChange`
- 工具栏中条件渲染 `<Button danger>`，带选中计数
- 点击确认后调用批量删除 API，清空选中状态并刷新

```typescript
rowSelection={{
  selectedRowKeys: selectedTypeIds,
  onChange: (keys) => { setSelectedTypeIds(keys as number[]) },
}}
```

### 9) Refresh Key 模式（跨组件刷新）
- 当父组件需要触发表格刷新（如复制完成后刷新另一个面板）
- 在父组件维护 `refreshKey` 状态，复制成功后递增
- 子组件将 `refreshKey` 加入 `watchingStates`

```typescript
// RouteComponent
const [refreshKey, setRefreshKey] = useState(0)

const handleBatchCopyEntries = useCallback(async (entryIds, targetTypeId) => {
  await DictApi.entryBatchCopy({ entryIds, targetTypeId })
  setRefreshKey(k => k + 1)
}, [])

// DictEntryPanel - usePagination
watchingStates: [selectedType?.id, searchText, refreshKey],
```

### 10) 筛选清理交互
- headerTitle 拼接当前筛选信息 + `清除筛选` Tag
- 清除筛选时调用父组件回调，重置筛选状态

### 11) 跨面板交互 — 拖拽复制
- **源表**（如右侧从列表）：每行设 `draggable: true` + `onDragStart`
  - 有选中行时传递全部选中 ID，否则传递当前行 ID
  - `effectAllowed = 'copy'` 表示复制而非移动
- **目标表**（如左侧主列表）：通过 `onDragOver/onDragEnter/onDragLeave/onDrop` 接收
  - 不能复制到自身（如 `record.id === selectedType?.id` 时拦截）
  - `hoveredDropTypeId` 状态控制拖拽悬停高亮，`dropEffect` 在禁止时设为 `'none'`
- 数据传输通过 `dataTransfer.setData('text/plain', JSON.stringify(ids))` 序列化

```typescript
// 源表：配置 draggable 行
onRow={record => ({
  draggable: true,
  onDragStart: (e) => {
    const ids = selectedIds.length > 0 ? selectedIds : [record.id]
    e.dataTransfer.setData('text/plain', JSON.stringify(ids))
    e.dataTransfer.effectAllowed = 'copy'
  },
})}

// 目标表：接收拖放
onRow={record => ({
  onDragOver: (e) => {
    if (record.id === sourceId) { e.dataTransfer.dropEffect = 'none'; return }
    e.preventDefault()
    e.dataTransfer.dropEffect = 'copy'
  },
  onDragEnter: (e) => {
    if (record.id === sourceId) { return }
    e.preventDefault()
    setHoveredDropId(record.id)
  },
  onDragLeave: () => setHoveredDropId(undefined),
  onDrop: (e) => {
    e.preventDefault()
    if (record.id === sourceId) { return }
    const raw = e.dataTransfer.getData('text/plain')
    if (!raw) return
    const ids: number[] = JSON.parse(raw)
    onBatchCopy(ids, record.id)
  },
})}
```

### 12) 数据映射与稳健性
- 在 `usePagination` 的 `data` 适配阶段处理后端空值（如 `numericValue ?? 0`）
- 页面内只处理"展示所需的最小映射"，复杂转换下沉到 API/domain 层
- 成功提示文案统一简短："保存成功"、"删除成功"、"操作成功"

### 13) 后端 Handler 语法
- 列表：`repo.XxxRepo.ListWithPaging(ctx, orm.DB(), req)` 使用 `v1.PagingRequest`
- 创建：`repo.XxxRepo.Create(...)`，手动填充 `CreatedBy`/`UpdatedBy`
- 更新：使用指针字段 + `UpdateNoNilMap` 实现 Patch 语义
- 状态切换：`UpdateMap` 显式指定更新列
- 删除：`SoftDelete`
- 批量操作：`Transaction` 包裹，先删子表再删主表

详见 `skills/backend/admin-service.md`。

## 注意事项
1. `watchingStates` 监听多个状态时，任一变化都会触发重载（会重置到第 1 页）
2. `useMemo` 的 `columns` 必须包含 `page`、`pageSize` 等渲染依赖，否则翻页后序号不更新
3. `ModalForm` 的 `onOpenChange` 必须处理关闭清理，防止状态残留
4. 后端 `v1.PagingRequest` 的 `query` 字段传递 JSON 字符串，用 `$and`/`$or` + `__icontains` 语法
5. API 接口命名统一使用 `Post` 方法，`cacheFor: 0` 避免缓存干扰
